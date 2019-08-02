package neotest

import (
	"encoding/json"
	"fmt"
	"github.com/hzxiao/goutil"
	"strings"
)

var (
	ErrVariableUndefine = fmt.Errorf("variable undefine")
)

var internalVarMap = goutil.Map{
	"neotest": goutil.Map{
		"version": "0.1",
		"author":  "hz",
	},
	"resp": nil,
}

type VM struct {
	variable goutil.Map
	commands []Commander

	CurHttpReq *HttpRequest
}

func NewVM(commands []Commander) *VM {
	vm := &VM{
		variable: goutil.Map{},
		commands: commands,
	}

	for k, v := range internalVarMap {
		vm.variable.Set(k, v)
	}
	return vm
}

func (vm *VM) Var(ID string) (interface{}, bool) {
	ID = strings.Replace(ID, ".", "/", -1)
	if strings.Contains(ID, "/") {
		v, _ := vm.variable.GetP(ID)
		return v, v != nil
	}
	v, ok := vm.variable[ID]
	if !ok {
		return nil, ok
	}
	return v, ok
}

func (vm *VM) StringV(ID string) (string, bool) {
	v, ok := vm.Var(ID)
	if !ok {
		return "", false
	}
	return goutil.String(v), true
}

func (vm *VM) FloatV(ID string) (float64, bool) {
	v, ok := vm.Var(ID)
	if !ok {
		return 0, false
	}
	return goutil.Float64(v), true
}

func (vm *VM) BoolV(ID string) (bool, bool) {
	v, ok := vm.Var(ID)
	if !ok {
		return false, false
	}
	return goutil.Bool(v), true
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

//SendHttp send current http request
func (vm *VM) SendHttp() error {
	code, header, body, err := vm.CurHttpReq.Send()
	if err != nil {
		return err
	}

	var m goutil.Map
	err = json.Unmarshal([]byte(body.(string)), &m)
	if err == nil {
		body = m
	}

	vm.StoreVar("resp", goutil.Map{
		"code":   code,
		"header": header,
		"body":   body,
	})

	//clear cur req
	vm.CurHttpReq = nil
	return nil
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

//CheckInternalVarID check whether ID is internal var and validation
func CheckInternalVarID(ID string) (bool, error) {
	if !strings.Contains(ID, ".") {
		return false, nil
	}

	fields := strings.Split(ID, ".")
	if _, ok := internalVarMap[fields[0]]; !ok {
		return true, fmt.Errorf("unknown internal variable: %v", fields[0])
	}

	return true, nil
}
