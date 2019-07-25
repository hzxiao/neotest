package neotest

import (
	"fmt"
	"github.com/hzxiao/goutil"
	"strings"
)

type ExprType uint

const (
	Invalid ExprType = iota
	Bool
	Float
	String
	Identity
	Variable
	SubCommand
)

type ExprNode interface {
	SetParent(ExprNode)
	Children() []ExprNode
	Parent() ExprNode
	Run(vm *VM) (interface{}, error)
	Type() ExprType
}

var _ ExprNode = new(exprBackground)

type exprBackground struct {
	parent   ExprNode
	children []ExprNode
}

func (eb *exprBackground) SetParent(e ExprNode) { eb.parent = e }

func (eb *exprBackground) Parent() ExprNode { return eb.parent }

func (eb *exprBackground) Children() []ExprNode { return eb.children }

func (*exprBackground) Run(vm *VM) (interface{}, error) { return nil, nil }

func (*exprBackground) Type() ExprType { return Invalid }

type stringExpr struct {
	exprBackground
	val interface{}
}

func (expr *stringExpr) Run(vm *VM) (interface{}, error) {
	return nil, nil
}

func (expr *stringExpr) Type() ExprType {
	return String
}

type boolExpr struct {
	exprBackground
	val interface{}
}

func (expr *boolExpr) Run(vm *VM) (interface{}, error) {
	ID, yes := isVar(expr.val)
	if yes {
		v, exist := vm.Var(ID)
		if !exist {
			return nil, fmt.Errorf("%v: %v", ErrVariableUndefine.Error(), ID)
		}
		if "bool" != fmt.Sprintf("%T", v) {
			return nil, fmt.Errorf("cannot use '%v' (type %T) as bool", ID, v)
		}

		return v, nil
	}
	return goutil.String(expr.val) == "true", nil
}

func (expr *boolExpr) Type() ExprType {
	return Bool
}

type IDExpr struct {
	exprBackground
	ID string
}

func (expr *IDExpr) Run(vm *VM) (interface{}, error) {
	return expr.ID, nil
}

func (expr *IDExpr) Type() ExprType {
	return Identity
}

func isVar(v interface{}) (ID string, yes bool) {
	text, ok := v.(string)
	if !ok {
		return
	}
	if strings.HasPrefix(text, "$(") && strings.HasSuffix(text, ")") {
		l := len(text)
		ID = text[2:l-1]
		yes = true
		return
	}

	return
}