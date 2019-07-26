package neotest

import (
	"fmt"
	"github.com/hzxiao/goutil"
	"github.com/spf13/cobra"
)

//EchoCmd 'echo' command
type EchoCmd struct {
	*Cmd
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

//LetCmd 'let' command
type LetCmd struct {
	*Cmd
}

func NewLetCmd(line int) *LetCmd {
	return &LetCmd{
		Cmd: NewCmd("let", "let @ID <object>", line),
	}
}

func (let *LetCmd) Exec(vm *VM) error {
	if len(let.exprList) != 2 {
		return fmt.Errorf("num of expr must be 2, but it is %v", len(let.exprList))
	}

	v, _ := let.exprList[0].Run(vm)
	ID := v.(string)

	rightValue, err := let.exprList[1].Run(vm)
	if err != nil {
		return err
	}

	old, exist := vm.Var(ID)
	if exist {
		oldType, newType := fmt.Sprintf("%T", old), fmt.Sprintf("%T", rightValue)
		if oldType != newType {
			return fmt.Errorf("cannot use '%v' (type %v) as %v", ID, newType, oldType)
		}
	}

	return vm.StoreVar(ID, rightValue)
}
