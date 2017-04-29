package main

import (
	"fmt"
	"github.com/loanzen/dagflow"
	"time"
)

func main() {
	state := dagflow.NewMemoryState()
	state.Set("Inputs", []interface{}{3 ,4})
	dag := dagflow.NewDag("test", state, nil)
	op1 := &dagflow.NoopOperator{}

	op2 := dagflow.GoFuncOperator(func (state dagflow.State, logger dagflow.Logger) error {
		time.Sleep(5 * time.Second)
		return nil
	})

	op3 := dagflow.GoFuncOperator(func (state dagflow.State, logger dagflow.Logger) error {
		time.Sleep(2 * time.Second)
		return nil
	})

	op4 := &dagflow.NoopOperator{}
	op5 := &dagflow.NoopOperator{}


	node1:= dagflow.NewNode("op1", op1, true, false, 5)
	node2:= dagflow.NewNode("op2", op2, true, false, 5)
	node3:= dagflow.NewNode("op3", op3, true, false, 5)
	node4:= dagflow.NewNode("op4", op4, true, false, 5)
	node5:= dagflow.NewNode("op5", op5, true, false, 5)
	dag.AddChild("test", node1)
	dag.AddChild("op1", node2)
	dag.AddChild("op2", node5)
	dag.AddChild("op1", node3)
	dag.AddChild("op3", node4)
	dag.AddChild("op4", node5)

	err := dag.Solve()
	if err != nil {
		fmt.Println("Failed Solving DAg", err)
	} else {
		fmt.Println("Done Solving DAg", err)
	}

}

