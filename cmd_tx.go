package neotest

import (
	"fmt"
	"github.com/hzxiao/goutil"
	"github.com/hzxiao/neotest/pkg/neo"
)

//TxCmd declare neo tx command
type TxCmd struct {
	*Cmd
}

func NewTxCmd(line int) *TxCmd {
	return &TxCmd{
		Cmd: NewCmd("tx", "tx <desc>", line),
	}
}

func (c *TxCmd) Exec(vm *VM) error {
	err := checkExprNumAndType(c.exprList, []int{0, 1}, String)
	if err != nil {
		return err
	}

	if vm.CurTx != nil {
		return fmt.Errorf("there is already a tx not be handled")
	}

	var name string
	if len(c.exprList) > 0 {
		name, err = toString(c.RunExprIndexOf(0, vm))
		if err != nil {
			return err
		}
	}

	vm.CurTx = neo.NewTx(name)

	return nil
}

func (c *TxCmd) CheckExpr(varType map[string]string) error {
	return checkExprNumAndType(c.exprList, []int{0, 1}, String)
}

//TxVCmd neo tx version command
type TxVCmd struct {
	*Cmd
}

func NewTxVCmd(line int) *TxVCmd {
	return &TxVCmd{
		Cmd: NewCmd("tx-v", "tx-v 0|1", line),
	}
}

func (c *TxVCmd) Exec(vm *VM) error {
	err := checkExprNumAndType(c.exprList, []int{1}, Float)
	if err != nil {
		return err
	}

	if vm.CurTx == nil {
		return fmt.Errorf("there is no declare a tx before")
	}

	v, err := toFloat64(c.RunExprIndexOf(0, vm))
	if err != nil {
		return err
	}
	if v != 0 && v != 1 {
		return fmt.Errorf("wrong tx version")
	}

	vm.CurTx.Version = uint8(v)
	return nil
}

func (c *TxVCmd) CheckExpr(varType map[string]string) error {
	return checkExprNumAndType(c.exprList, []int{1}, Float)
}

//TxTypeCmd neo tx type command
type TxTypeCmd struct {
	*Cmd
}

func NewTxTypeCmd(line int) *TxTypeCmd {
	return &TxTypeCmd{
		Cmd: NewCmd("tx-type", "tx-type <txType>", line),
	}
}

func (c *TxTypeCmd) Exec(vm *VM) error {
	err := checkExprNumAndType(c.exprList, []int{1}, String)
	if err != nil {
		return err
	}

	if vm.CurTx == nil {
		return fmt.Errorf("there is no declare a tx before")
	}

	Type, err := toString(c.RunExprIndexOf(0, vm))
	if err != nil {
		return err
	}

	return vm.CurTx.SetType(Type)
}

func (c *TxTypeCmd) CheckExpr(varType map[string]string) error {
	return checkExprNumAndType(c.exprList, []int{1}, String)
}

//TxFeeCmd neo tx fee command
type TxFeeCmd struct {
	*Cmd
}

func NewTxFeeCmd(line int) *TxFeeCmd {
	return &TxFeeCmd{
		Cmd: NewCmd("tx-type", "tx-fee <number>", line),
	}
}

func (c *TxFeeCmd) Exec(vm *VM) error {
	err := checkExprNumAndType(c.exprList, []int{1}, Float)
	if err != nil {
		return err
	}

	if vm.CurTx == nil {
		return fmt.Errorf("there is no declare a tx before")
	}

	fee, err := toFloat64(c.RunExprIndexOf(0, vm))
	if err != nil {
		return err
	}
	vm.CurTx.SetFee(fee)

	return nil
}

func (c *TxFeeCmd) CheckExpr(varType map[string]string) error {
	return checkExprNumAndType(c.exprList, []int{1}, Float)
}

//TxAttrCmd neo tx fee command
type TxAttrCmd struct {
	*Cmd
}

func NewTxAttrCmd(line int) *TxAttrCmd {
	return &TxAttrCmd{
		Cmd: NewCmd("tx-type", "tx-type <txType>", line),
	}
}

func (c *TxAttrCmd) Exec(vm *VM) error {
	err := checkExprNumAndType(c.exprList, []int{2}, String, String)
	if err != nil {
		return err
	}

	if vm.CurTx == nil {
		return fmt.Errorf("there is no declare a tx before")
	}

	usage, err := toString(c.RunExprIndexOf(0, vm))
	if err != nil {
		return err
	}
	if !neo.ValidAttrUsage(usage) {
		return fmt.Errorf("invalid attr usage")
	}

	data, err := toString(c.RunExprIndexOf(0, vm))
	if err != nil {
		return err
	}

	vm.CurTx.Param.Attr = append(vm.CurTx.Param.Attr, goutil.Map{
		"usage": usage,
		"data":  data,
	})
	return nil
}

func (c *TxAttrCmd) CheckExpr(varType map[string]string) error {
	return checkExprNumAndType(c.exprList, []int{2}, String, String)
}

//TxInitiatorCmd neo tx fee command
type TxInitiatorCmd struct {
	*Cmd
}

func NewTxInitiatorCmd(line int) *TxInitiatorCmd {
	return &TxInitiatorCmd{
		Cmd: NewCmd("tx-initiator", "tx-initiator <privateKey>", line),
	}
}

func (c *TxInitiatorCmd) Exec(vm *VM) error {
	err := checkExprNumAndType(c.exprList, []int{1}, String)
	if err != nil {
		return err
	}

	if vm.CurTx == nil {
		return fmt.Errorf("there is no declare a tx before")
	}

	privateKey, err := toString(c.RunExprIndexOf(0, vm))
	if err != nil {
		return err
	}

	return vm.CurTx.Param.SetInitiator(privateKey)
}

func (c *TxInitiatorCmd) CheckExpr(varType map[string]string) error {
	return checkExprNumAndType(c.exprList, []int{1}, String)
}

//TxVoutCmd neo tx vout command
type TxVoutCmd struct {
	*Cmd
}

func NewTxVoutCmd(line int) *TxVoutCmd {
	return &TxVoutCmd{
		Cmd: NewCmd("tx-vout", "tx-vout <asset_hash> <address> <value>", line),
	}
}

func (c *TxVoutCmd) Exec(vm *VM) error {
	err := checkExprNumAndType(c.exprList, []int{3}, String, String, Float)
	if err != nil {
		return err
	}

	if vm.CurTx == nil {
		return fmt.Errorf("there is no declare a tx before")
	}

	asset, err := toString(c.RunExprIndexOf(0, vm))
	if err != nil {
		return err
	}
	address, err := toString(c.RunExprIndexOf(1, vm))
	if err != nil {
		return err
	}
	value, err := toFloat64(c.RunExprIndexOf(2, vm))
	if err != nil {
		return err
	}

	switch asset {
	case "gas":
		asset = neo.GasAssetHash
	case "neo":
		asset = neo.NeoAssetHash
	}
	vm.CurTx.Param.Vout = append(vm.CurTx.Param.Vout, goutil.Map{
		"asset":   asset,
		"address": address,
		"value":   value,
	})
	return nil
}

func (c *TxVoutCmd) CheckExpr(varType map[string]string) error {
	return checkExprNumAndType(c.exprList, []int{3}, String, String, Float)
}

//TxInvokeCmd neo tx invoke command
type TxInvokeCmd struct {
	*Cmd
}

func NewTxInvokeCmd(line int) *TxInvokeCmd {
	return &TxInvokeCmd{
		Cmd: NewCmd("tx-invoke", "tx-invoke '<json-data>'", line),
	}
}

func (c *TxInvokeCmd) Exec(vm *VM) error {
	err := checkExprNumAndType(c.exprList, []int{1}, String)
	if err != nil {
		return err
	}
	if vm.CurTx == nil {
		return fmt.Errorf("there is no declare a tx before")
	}

	invoke, err := toString(c.RunExprIndexOf(0, vm))
	if err != nil {
		return err
	}
	return vm.CurTx.ParseScript(invoke)
}

func (c *TxInvokeCmd) CheckExpr(varType map[string]string) error {
	return checkExprNumAndType(c.exprList, []int{1}, String)
}

//TxInvokeScriptCmd neo tx invoke function command
type TxInvokeFuncCmd struct {
	*Cmd
}

func NewTxInvokeFuncCmd(line int) *TxInvokeFuncCmd {
	return &TxInvokeFuncCmd{
		Cmd: NewCmd("tx-invokefunc", "tx-invokefunc '<json-data>'", line),
	}
}

func (c *TxInvokeFuncCmd) Exec(vm *VM) error {
	err := checkExprNumAndType(c.exprList, []int{1}, String)
	if err != nil {
		return err
	}
	if vm.CurTx == nil {
		return fmt.Errorf("there is no declare a tx before")
	}

	invoke, err := toString(c.RunExprIndexOf(0, vm))
	if err != nil {
		return err
	}
	return vm.CurTx.ParseScript(invoke)
}

func (c *TxInvokeFuncCmd) CheckExpr(varType map[string]string) error {
	return checkExprNumAndType(c.exprList, []int{1}, String)
}

//TxInvokeScriptCmd neo tx invoke script command
type TxInvokeScriptCmd struct {
	*Cmd
}

func NewTxInvokeScriptCmd(line int) *TxInvokeScriptCmd {
	return &TxInvokeScriptCmd{
		Cmd: NewCmd("tx-invokescript", "tx-invokescript <script>", line),
	}
}

func (c *TxInvokeScriptCmd) Exec(vm *VM) error {
	err := checkExprNumAndType(c.exprList, []int{1}, String)
	if err != nil {
		return err
	}

	return nil
}

func (c *TxInvokeScriptCmd) CheckExpr(varType map[string]string) error {
	return checkExprNumAndType(c.exprList, []int{1}, String)
}

//TxWitnessCmd neo tx witness command
type TxWitnessCmd struct {
	*Cmd
}

func NewTxWitnessCmd(line int) *TxWitnessCmd {
	return &TxWitnessCmd{
		Cmd: NewCmd("tx-witness", "tx-witness <witness> <invocation>", line),
	}
}

func (c *TxWitnessCmd) Exec(vm *VM) error {
	err := checkExprNumAndType(c.exprList, []int{1, 2}, String, String)
	if err != nil {
		return err
	}

	if vm.CurTx == nil {
		return fmt.Errorf("there is no declare a tx before")
	}

	witness, err := toString(c.RunExprIndexOf(0, vm))
	if err != nil {
		return err
	}

	var inv string
	if len(c.exprList) > 1 {
		inv, err = toString(c.RunExprIndexOf(1, vm))
		if err != nil {
			return err
		}
	}

	vm.CurTx.Param.Witness = append(vm.CurTx.Param.Witness, goutil.Map{
		"witness": witness,
		"inv":     inv,
	})
	return nil
}

func (c *TxWitnessCmd) CheckExpr(varType map[string]string) error {
	return checkExprNumAndType(c.exprList, []int{1, 2}, String, String)
}

//TxSendCmd send neo tx command
type TxSendCmd struct {
	*Cmd
}

func NewTxSendCmd(line int) *TxSendCmd {
	return &TxSendCmd{
		Cmd: NewCmd("tx-send", "tx-send <seed>", line),
	}
}

func (c *TxSendCmd) Exec(vm *VM) error {
	err := checkExprNumAndType(c.exprList, []int{1}, String)
	if err != nil {
		return err
	}

	if vm.CurTx == nil {
		return fmt.Errorf("there is no declare a tx before")
	}

	node, err := toString(c.RunExprIndexOf(0, vm))
	if err != nil {
		return err
	}

	err = vm.CurTx.Complete(node)
	if err != nil {
		return err
	}

	err = vm.SendTx(node)
	if err != nil {
		return err
	}

	err = vm.WaitTx()
	if err != nil {
		return err
	}
	return nil
}

func (c *TxSendCmd) CheckExpr(varType map[string]string) error {
	return checkExprNumAndType(c.exprList, []int{1}, String)
}
