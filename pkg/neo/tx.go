package neo

import (
	"fmt"
	"github.com/CityOfZion/neo-go/pkg/core/transaction"
	"github.com/CityOfZion/neo-go/pkg/util"
	"github.com/hzxiao/goutil"
	"strconv"
)

type TxParam struct {
	Fee       util.Fixed8
	Attr      []goutil.Map
	Initiator string
	Vout      []goutil.Map
	Script    string
	Witness   []goutil.Map
}

type Tx struct {
	transaction.Transaction

	Param *TxParam
	Name  string
}

func NewTx(name string) *Tx {
	return &Tx{Transaction: transaction.Transaction{}, Name: name, Param: &TxParam{}}
}

func (tx *Tx) SetType(typ string) error {
	switch typ {
	case "contract":
		tx.Type = transaction.ContractType
	case "invocation":
		tx.Type = transaction.InvocationType
	default:
		return fmt.Errorf("unsupport tx type(%v)", typ)
	}

	return nil
}

func (tx *Tx) SetFee(fee float64) error {
	s := strconv.FormatFloat(fee, 'f', -1, 64)
	var err error
	tx.Param.Fee, err = util.Fixed8DecodeString(s)
	return err
}

func (tx *Tx) ParseScript(raw string) error {
	return nil
}

func (tx *Tx) Complete(node string) error {
	return nil
}

func (tx *Tx) ToMap() goutil.Map {
	return nil
}

// RelayTx relay tx to the neo node
func RelayTx(tx *Tx, node string) error {
	return nil
}
