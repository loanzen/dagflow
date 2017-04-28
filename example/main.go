package main

import (
	"fmt"
	"github.com/loanzen/dagflow"
)

func main() {
	state := dagflow.NewMemoryState()
	dag := dagflow.NewDag("test", state)
	op1 := &dagflow.GoFuncOperator{
		Inputs: []interface{}{2, 3},
		Func: func(state dagflow.State, inputs []interface{}) error {
			fmt.Println("Running First Operator")
			sum := 0
			for _, x := range inputs {
				sum += x.(int)
			}
			state.Set("sum", sum)
			return nil
		},
	}

	op2 := &dagflow.GoFuncOperator{
		Func: func(state dagflow.State, inputs []interface{}) error {
			sum, _  := state.Get("sum")
			fmt.Println("Running Second Operator: ", sum)
			return nil
		},
	}

	op3 := &dagflow.GoFuncOperator{
		Func: func(state dagflow.State, inputs []interface{}) error {
			fmt.Println("Running Third Operator: ")
			return nil
		},
	}


	add := dagflow.NewNode("add", op1, true, false, 5)
	dag.AddChild("test", add)
	dag.AddChild("add", dagflow.NewNode("print", op2, true, false, 5))
	dag.AddChild("print", dagflow.NewNode("third", op3, true, false, 5))
	dag.AddChild("third", add)
	err := dag.Solve()
	if err != nil {
		fmt.Println("Failed Solving DAg", err)
	} else {
		fmt.Println("Done Solving DAg", err)
	}

}

