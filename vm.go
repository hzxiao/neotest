package neotest

import (
	"fmt"
	"github.com/hzxiao/goutil"
)

var (
	ErrVariableUndefine = fmt.Errorf("variable undefine")
)

type VM struct {
	variable goutil.Map
	commands []Commander
}

func NewVM(commands []Commander) *VM {
	return &VM{
		variable: goutil.Map{},
		commands: commands,
	}
}

func (vm *VM) Var(ID string) (interface{}, bool) {
	v, ok := vm.variable[ID]
	if !ok {
		return nil, ok
	}
	return v, ok
}

func (vm *VM) StringV(ID string) (string, bool) {
	if !vm.variable.Exist(ID) {
		return "", false
	}
	return vm.variable.GetString(ID), true
}

func (vm *VM) FloatV(ID string) (float64, bool) {
	if !vm.variable.Exist(ID) {
		return 0, false
	}
	return vm.variable.GetFloat64(ID), true
}

func (vm *VM) BoolV(ID string) (bool, bool) {
	if !vm.variable.Exist(ID) {
		return false, false
	}
	return vm.variable.GetBool(ID), true
}

func (vm *VM) StoreVar(ID string, v interface{}) error {
	vm.variable.Set(ID, v)
	return nil
}

func (vm *VM) VarByType(ID string, typ string) (interface{}, error) {
	v, exist := vm.Var(ID)
	if !exist {
		return nil, fmt.Errorf("%v: %v", ErrVariableUndefine.Error(), ID)
	}
	if typ != fmt.Sprintf("%T", v) {
		return nil, fmt.Errorf("cannot use '%v' (type %T) as %v", ID, v, typ)
	}

	return v, nil
}

func (vm *VM) Run() error {
	var err error
	for _, cmd := range vm.commands {
		err = cmd.Exec(vm)
		if err != nil {
			return fmt.Errorf("line %v: exec %v err: %v", cmd.Line(), cmd.Name(), err)
		}
	}
	return nil
}