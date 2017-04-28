package dagflow

import (
	"github.com/pkg/errors"
	"fmt"
)

type Dag struct {
	root *Node
	nodes map[string]*Node
	state State
}

type DagFlowStatus int

const (
	PENDING DagFlowStatus = iota
	RUNNING
	SKIPPED
	SUCCESS
	FAILED
)


func NewDag(name string, state State) *Dag {
	root := NewNode(name, new(NoopOperator), true, true, 0)
	root.Status = SUCCESS
	dag := &Dag{
		root: root,
		nodes: make(map[string]*Node),
		state: state,
	}
	dag.nodes[name] = root
	return dag
}

func (dag *Dag) NumOfNodes() int {
	return len(dag.nodes)
}

func (dag *Dag) AddChild(parent string, child *Node) {
	pnode, ok := dag.nodes[parent]
	if !ok {
		panic(fmt.Sprintf("Dag doesn't contain any node named %s", parent))
	}
	dag.AddChildAt(parent, child, pnode.NumOfChildren())
}

func (dag *Dag) AddChildAt(parent string, child *Node, pos int) {
	dag.nodes[parent].addChildAt(child, pos)
	dag.nodes[child.Name] = child
}

func (dag *Dag) AddDag(parent string, child *Dag) {
	dag.AddChild(parent, child.root)
	for k, v:= range child.nodes {
		dag.nodes[k] = v
	}
}

func (dag *Dag) IsSolvable() bool {
	visited := make(map[string]bool)
	current := dag.root
	return !hasCirclularDep(current, visited)
}

func (dag *Dag) Solve() error {
	if (!dag.IsSolvable()) {
		return errors.New(
			fmt.Sprintf("%s dag is not solvable as it has cycles", dag.root.Name))
	}
	solver := newDagSolver(dag)
	return solver.Solve()
}

func hasCirclularDep(current *Node, visited map[string]bool) bool {
	visited[current.Name] = true
	for _, child := range current.Children {
		if _, ok := visited[child.Name]; ok {
		    return true
		}
		if (hasCirclularDep(child, visited)) {
			return true
		}
		delete(visited, child.Name)
	}
	delete(visited, current.Name)
	return false
}

type TaskUnit struct {
	node *Node
	completionStatus chan string
}

type dagSolver struct {
	dag *Dag
	status DagFlowStatus
	completed []*Node
	exit chan int
	completionStatus chan string
	queue chan *TaskUnit
}

func newDagSolver(dag *Dag) *dagSolver {
	solver := dagSolver{
		dag: dag, status: PENDING,
		completed: make([]*Node, 0, dag.NumOfNodes()),
		exit: make(chan int),
		completionStatus: make(chan string, dag.NumOfNodes()),
		queue: make(chan *TaskUnit, 10),
	}

	return &solver
}


func (solver *dagSolver) Solve() error {
	dag := solver.dag
	solver.completed = append(solver.completed, dag.root)
	go solver.work()
	go solver.solveChildren(dag.root)
	for len(solver.completed) < dag.NumOfNodes() {
		select {
		case <- solver.exit:
			solver.status = FAILED
			break

		case name := <- solver.completionStatus:
			node := solver.dag.nodes[name]
			if node.CanSolveChildren() {
				solver.completed = append(solver.completed, node)
				go solver.solveChildren(node)
			} else if node.Status == FAILED && node.Required {
				solver.exit <- 0
			}
			break
		}

		if solver.status == FAILED {
			break
		}
	}
	solver.exit <- 0
	return nil
}

func (solver *dagSolver) solveChildren(node *Node) {
	for _, child := range node.Children {
		if child.CanRun() {
			taskUnit := &TaskUnit{child, solver.completionStatus}
			solver.queue <- taskUnit
		}
	}
}

func (solver *dagSolver) work() {
	done := false
	for !done {
		select {
		case <- solver.exit:
			done = true
			break

		case task:= <- solver.queue:
			go func(){
				task.node.Solve(solver.dag.state)
				solver.completionStatus <- task.node.Name
			}()
			break
		}

	}
}





