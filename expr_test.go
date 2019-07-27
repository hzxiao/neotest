package neotest

import (
	"github.com/hzxiao/goutil/assert"
	"testing"
)

func TestVarExpr_Variables(t *testing.T) {
	v1 := &varExpr{val: "abc"}
	assert.Equal(t, []string(nil), v1.Variables())

	v2 := &varExpr{val: "a$(b)b"}
	assert.Equal(t, []string{"b"}, v2.Variables())

	v3 := &varExpr{val: "a$(b)c$(d)e$(12)f"}
	assert.Equal(t, []string{"b", "d"}, v3.Variables())
}

func TestStringExpr_Run(t *testing.T) {
	vm := NewVM(nil)
	vm.StoreVar("a", true)
	vm.StoreVar("b", 1.1)
	vm.StoreVar("c", "world")

	s1 := newStringExpr("$(c)")
	v1, err := s1.Run(vm)
	assert.NoError(t, err)
	assert.Equal(t, "world", v1)

	s2 := newStringExpr("hello,$(c),$(b),$(a)")
	v2, err := s2.Run(vm)
	assert.NoError(t, err)
	assert.Equal(t, "hello,world,1.1,true", v2)

	s3 := newStringExpr("hi")
	v3, err := s3.Run(vm)
	assert.NoError(t, err)
	assert.Equal(t, "hi", v3)

	s4 := newStringExpr("$(d)")
	_, err = s4.Run(vm)
	assert.Error(t, err)

	s5 := newStringExpr("hi,$(e)")
	_, err = s5.Run(vm)
	assert.Error(t, err)
}