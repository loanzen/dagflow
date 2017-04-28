package dagflow

import (
	"os"
	"os/exec"
	"fmt"
	"strings"
)


type Operator interface {
	Run(State) error
}


type NoopOperator struct {}

func (op *NoopOperator) Run(state State) error { return nil }


type CmdOperator struct {
	cmd *exec.Cmd
}

func (op *CmdOperator) Run(state State) error {
	env := os.Environ()

	for tup := range state.Iter() {
		env = append(env,
			fmt.Sprintf("%s=%s", strings.ToUpper(tup.Key), tup.Val))
	}

	op.cmd.Env = env
	return op.cmd.Run()
}

type GoFuncOperator struct {
	Inputs []interface{}
	Func func(State, []interface{}) error
}

func (op *GoFuncOperator) Run(state State) error {
	return op.Func(state, op.Inputs)
}
