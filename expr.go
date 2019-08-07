package neotest

import (
	"fmt"
	"github.com/hzxiao/goutil"
	"regexp"
	"strconv"
	"strings"
)

type ExprType uint

const (
	Invalid ExprType = iota
	Bool
	Float
	String
	Identity
	SubCommand
	InternalVar
)

func (t ExprType) String() string {
	switch t {
	case Invalid:
		return "invalid"
	case Bool:
		return "bool"
	case Float:
		return "number"
	case String:
		return "string"
	case Identity:
		return "identity"
	case SubCommand:
		return "subCmd"
	case InternalVar:
		return "internelVar"
	}

	return ""
}

type Variate interface {
	Variables() []string
}

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

var _ Variate = new(varExpr)
var _ Variate = new(boolExpr)
var _ Variate = new(stringExpr)
var _ Variate = new(floatExpr)

type varExpr struct {
	exprBackground
	val string
}

func (expr *varExpr) Variables() []string {
	all := regexp.MustCompile(`\$\([a-zA-Z_][a-zA-Z0-9_\\.]*\)`).FindAllString(expr.val, -1)

	var IDs []string
	for _, v := range all {
		IDs = append(IDs, v[2:len(v)-1])
	}

	return IDs
}

type stringExpr struct {
	varExpr
}

func newStringExpr(val string) ExprNode {
	return &stringExpr{varExpr{val: val}}
}

func (expr *stringExpr) Run(vm *VM) (interface{}, error) {
	val := expr.val

	//find contain var
	allIndex := regexp.MustCompile(`\$\([a-zA-Z_][a-zA-Z0-9_\\.]*\)`).FindAllStringIndex(expr.val, -1)
	for _, idx := range allIndex {
		ID, _ := isVar(expr.val[idx[0]:idx[1]])
		v, exist := vm.StringV(ID)
		if !exist {
			return nil, fmt.Errorf("%v: %v", ErrVariableUndefine.Error(), ID)
		}

		val = strings.Replace(val, expr.val[idx[0]:idx[1]], v, 1)
	}

	return val, nil
}

func (expr *stringExpr) Type() ExprType {
	return String
}

type boolExpr struct {
	varExpr
}

func newBoolExpr(val string) ExprNode {
	return &boolExpr{varExpr{val: val}}
}

func (expr *boolExpr) Run(vm *VM) (interface{}, error) {
	ID, yes := isVar(expr.val)
	if yes {
		return vm.VarByType(ID, "bool")
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

type floatExpr struct {
	varExpr
}

func newFloatExpr(val string) ExprNode {
	return &floatExpr{varExpr{val: val}}
}

func (expr *floatExpr) Run(vm *VM) (interface{}, error) {
	ID, yes := isVar(expr.val)
	if yes {
		return vm.VarByType(ID, "float64")
	}

	v, _ := strconv.ParseFloat(expr.val, 64)
	return v, nil
}

func (expr *floatExpr) Type() ExprType {
	return Float
}

type internalVarExpr struct {
	varExpr
}

func newInternalValExpr(val string) *internalVarExpr {
	return &internalVarExpr{varExpr{val: val}}
}

func (expr *internalVarExpr) Run(vm *VM) (interface{}, error) {
	ID, _ := isVar(expr.val)
	v, _ := vm.Var(ID)
	return v, nil
}

func (expr *internalVarExpr) Type() ExprType {
	return InternalVar
}

func isVar(v interface{}) (ID string, yes bool) {
	text, ok := v.(string)
	if !ok {
		return
	}
	if strings.HasPrefix(text, "$(") && strings.HasSuffix(text, ")") {
		l := len(text)
		ID = text[2 : l-1]
		yes = true
		return
	}

	return
}
