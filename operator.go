package dagflow

import (
	"os"
	"os/exec"
	"fmt"
	"strings"
)


type Operator interface {
	Run(State, Logger) error
}


type NoopOperator struct {}

func (op *NoopOperator) Run(state State, logger Logger) error { return nil }


type CmdOperator struct {
	cmd *exec.Cmd
}

func (op *CmdOperator) Run(state State, logger Logger) error {
	env := os.Environ()

	for tup := range state.Iter() {
		env = append(env,
			fmt.Sprintf("%s=%s", strings.ToUpper(tup.Key), tup.Val))
	}

	op.cmd.Env = env
	return op.cmd.Run()
}

type GoFuncOperator func(State, Logger) error

func (f GoFuncOperator) Run(state State, logger Logger) error {
	return f(state, logger)
}

