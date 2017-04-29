package dagflow


type Node struct {
	Name string
	Operator Operator
	Status DagFlowStatus
	Required bool
	ContinueOnErr bool
	Retries int
	Err error
	Parents []*Node
	Children []*Node
}

func NewNode(name string, op Operator, required bool,
	cntErr bool, retries int) *Node {
	return &Node{
		Name: name, Operator: op, Required: required,
		ContinueOnErr: cntErr, Retries: retries,
		Parents: make([]*Node, 0, 1),
		Children: make([]*Node, 0, 0),
	}
}



func (node *Node) NumOfChildren() int {
	return len(node.Children)
}

func (node *Node) NumOfParents() int {
	return len(node.Parents)
}


func (node *Node) Skip () {
	if (!node.IsComplete()) {
		node.Status = SKIPPED
	}
}

func (node *Node) CanRun() bool {
	for _, parent := range node.Parents {
		if !parent.IsComplete() {
			return false
		}
	}
	return true
}

func (node* Node) CanSolveChildren() bool {
	return node.Status == SUCCESS ||
		(node.Status == FAILED && !node.Required && node.ContinueOnErr)
}

func (node *Node) Solve(state State, logger Logger) error {
	err := node.Operator.Run(state, logger)
	if err != nil {
		node.Status = FAILED
		node.Err = err
	} else {
		node.Status = SUCCESS
	}
	return err
}

func (node* Node) IsComplete() bool {
	return node.Status == SUCCESS ||
		(node.Status == FAILED && !node.Required)
}


func (node *Node) addChild(child *Node) {
	node.addChildAt(child, node.NumOfChildren())
}

func (node *Node) addChildAt(child *Node, pos int) {
	Nodes := append(node.Children, nil)
	copy(Nodes[pos + 1:], Nodes[:pos])
	Nodes[pos] = child
	node.Children = Nodes
	child.Parents = append(child.Parents, node)
}

