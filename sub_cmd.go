package neotest

import (
	"fmt"
	"github.com/hzxiao/goutil"
	"os"
)

var _ Commander = new(SubCmd)
var _ ExprNode = new(SubCmd)
var _ Variate = new(SubCmd)
var _ Resultant = new(SubCmd)

type Resultant interface {
	ResultType() string
}

type SubCmd struct {
	*Cmd
	*varExpr
}

func (*SubCmd) Exec(vm *VM) error {
	return nil
}

func (*SubCmd) Run(vm *VM) (interface{}, error) {
	return nil, nil
}

func (*SubCmd) Type() ExprType {
	return SubCommand
}

func (*SubCmd) ResultType() string {
	return "string"
}

func (sc *SubCmd) Variables() []string {
	var variables []string
	for i := range sc.exprList {
		expr, ok := sc.exprList[i].(Variate)
		if ok {
			v := expr.Variables()
			if len(v) > 0 {
				variables = append(variables, v...)
			}
		}
	}
	return variables
}

type EnvSubCmd struct {
	SubCmd
}

func NewEnvSubCmd(line int) *EnvSubCmd {
	return &EnvSubCmd{
		SubCmd{
			Cmd:     NewCmd("env", "get env value from os", line),
			varExpr: &varExpr{},
		},
	}
}

func (env *EnvSubCmd) Run(vm *VM) (interface{}, error) {
	if len(env.exprList) != 1 {
		return nil, fmt.Errorf("num of expr must be 1, but it is %v", len(env.exprList))
	}

	v, err := env.exprList[0].Run(vm)
	if err != nil {
		return nil, err
	}

	return os.Getenv(goutil.String(v)), nil
}
