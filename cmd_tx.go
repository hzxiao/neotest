package neotest

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


	return nil
}

func (c *TxCmd) CheckExpr(varType map[string]string) error {
	err := checkExprNumAndType(c.exprList, []int{0, 1}, String)
	if err != nil {
		return err
	}

	return nil
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



	return nil
}

func (c *TxVCmd) CheckExpr(varType map[string]string) error {
	err := checkExprNumAndType(c.exprList, []int{1}, Float)
	if err != nil {
		return err
	}


	return nil
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



	return nil
}

func (c *TxTypeCmd) CheckExpr(varType map[string]string) error {
	err := checkExprNumAndType(c.exprList, []int{1}, String)
	if err != nil {
		return err
	}

	return nil
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


	return nil
}

func (c *TxFeeCmd) CheckExpr(varType map[string]string) error {
	err := checkExprNumAndType(c.exprList, []int{1}, Float)
	if err != nil {
		return err
	}


	return nil
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


	return nil
}

func (c *TxAttrCmd) CheckExpr(varType map[string]string) error {
	err := checkExprNumAndType(c.exprList, []int{2}, String, String)
	if err != nil {
		return err
	}


	return nil
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


	return nil
}

func (c *TxInitiatorCmd) CheckExpr(varType map[string]string) error {
	err := checkExprNumAndType(c.exprList, []int{1}, String)
	if err != nil {
		return err
	}

	return nil
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

	return nil
}

func (c *TxVoutCmd) CheckExpr(varType map[string]string) error {
	err := checkExprNumAndType(c.exprList, []int{3}, String, String, Float)
	if err != nil {
		return err
	}
	return nil
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

	return nil
}

func (c *TxInvokeCmd) CheckExpr(varType map[string]string) error {
	err := checkExprNumAndType(c.exprList, []int{1}, String)
	if err != nil {
		return err
	}
	return nil
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

	return nil
}

func (c *TxInvokeFuncCmd) CheckExpr(varType map[string]string) error {
	err := checkExprNumAndType(c.exprList, []int{1}, String)
	if err != nil {
		return err
	}
	return nil
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
	err := checkExprNumAndType(c.exprList, []int{1}, String)
	if err != nil {
		return err
	}
	return nil
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

	return nil
}

func (c *TxWitnessCmd) CheckExpr(varType map[string]string) error {
	err := checkExprNumAndType(c.exprList, []int{1, 2}, String, String)
	if err != nil {
		return err
	}
	return nil
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

	return nil
}

func (c *TxSendCmd) CheckExpr(varType map[string]string) error {
	err := checkExprNumAndType(c.exprList, []int{1}, String)
	if err != nil {
		return err
	}
	return nil
}
