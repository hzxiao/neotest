package neotest

import (
	"fmt"
	"github.com/hzxiao/goutil"
	"github.com/k0kubun/pp"
	"github.com/spf13/cobra"
	"reflect"
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

func (echo *EchoCmd) CheckExpr(varType map[string]string) error {
	if len(echo.exprList) == 0 {
		return fmt.Errorf("no enough arguments")
	}
	return nil
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

func (let *LetCmd) CheckExpr(varType map[string]string) error {
	if len(let.exprList) != 2 {
		return fmt.Errorf("num of expr must be 2, but it is %v", len(let.exprList))
	}
	if let.exprList[0].Type() != Identity {
		return fmt.Errorf("invaild cmd syntax: the first argument should be @ID")
	}

	second := let.exprList[1]
	var vType string
	switch second.Type() {
	case Bool, Float, String:
		if second.Type() == Bool {
			vType = "bool"
		} else if second.Type() == Float {
			vType = "float"
		} else {
			vType = "string"
		}
	case SubCommand:
		vType = second.(Resultant).ResultType()
	case InternalVar:
		vType = "internal"
	default:
		return fmt.Errorf("invalid cmd syntax: invalid second argument type")
	}

	//check variable exist
	IDs := second.(Variate).Variables()
	for _, id := range IDs {
		err := CheckVar(id, varType)
		if err != nil {
			return err
		}
	}

	//record variable type on source-parsing stage
	ID := let.exprList[0].(*IDExpr).ID
	varType[ID] = vType
	return nil
}

type EqualCmd struct {
	*Cmd
}

func NewEqualCmd(line int) *EqualCmd {
	return &EqualCmd{
		Cmd: NewCmd("equal", "equal <object1> <object2>", line),
	}
}

func (eq *EqualCmd) Exec(vm *VM) error {
	if len(eq.exprList) != 2 {
		return fmt.Errorf("num of expr must be 2, but it is %v", len(eq.exprList))
	}

	first, err := eq.exprList[0].Run(vm)
	if err != nil {
		return err
	}
	second, err := eq.exprList[1].Run(vm)
	if err != nil {
		return err
	}

	ok := reflect.DeepEqual(first, second)
	if !ok {
		pp.Printf("line: %v: %v != %v\n", eq.line, first, second)
	}

	return nil
}

func (eq *EqualCmd) CheckExpr(varType map[string]string) error {
	if len(eq.exprList) != 2 {
		return fmt.Errorf("num of expr must be 2, but it is %v", len(eq.exprList))
	}
	return nil
}