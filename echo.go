package neotest

import (
	"fmt"
	"github.com/hzxiao/goutil"
	"github.com/spf13/cobra"
)

type EchoCmd struct {
	*Cmd
	exprList []ExprNode
}

func NewEchoCmd(line int) *EchoCmd {
	return &EchoCmd{
		Cmd: NewCmd("echo", "echo <object1> <object2>", line),
	}
}

func (echo *EchoCmd) Exec(vm *VM) error {
	echo.cmd.RunE = func(cmd *cobra.Command, args []string) error {
		for _, arg := range args {
			fmt.Printf("%v ", arg)
		}
		fmt.Println()
		return nil
	}

	var values []string
	for _, expr := range echo.exprList {
		result, err := expr.Run(vm)
		if err != nil {
			return err
		}

		values = append(values, goutil.String(result))
	}
	echo.cmd.SetArgs(values)
	return echo.cmd.Execute()
}