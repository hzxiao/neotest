package neo

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/CityOfZion/neo-go/pkg/core/transaction"
	"github.com/CityOfZion/neo-go/pkg/crypto"
	"github.com/CityOfZion/neo-go/pkg/util"
	"github.com/CityOfZion/neo-go/pkg/wallet"
	"github.com/hzxiao/goutil"
	"math"
	"strconv"
)

type TxParam struct {
	Fee       util.Fixed8
	Attr      []goutil.Map
	Initiator *wallet.PrivateKey
	Vout      []goutil.Map
	Script    string
	Witness   []goutil.Map
}

func (p *TxParam) SetInitiator(initiator string) error {
	privateKey, err := wallet.NewPrivateKeyFromHex(initiator)
	if err != nil {
		return err
	}
	p.Initiator = privateKey
	return nil
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

func (tx *Tx) SetFee(fee float64) {
	tx.Param.Fee = Fixed8FromFloat64(fee)
}

func (tx *Tx) ParseScript(raw string) error {
	return nil
}

func (tx *Tx) Complete(node string) error {
	param := tx.Param
	if param == nil {
		return fmt.Errorf("tx param is nil")
	}
	if param.Initiator == nil {
		return fmt.Errorf("initiator is emptty")
	}

	address, err := param.Initiator.Address()
	if err != nil {
		return err
	}
	//fee
	if param.Fee > 0 {
		fee, _ := strconv.ParseFloat(param.Fee.String(), 64)
		all, inputs, err := getReference(GasAssetHash, address, fee, node)
		if err != nil {
			return err
		}
		tx.Inputs = append(tx.Inputs, inputs...)
		if all > fee { //redundant
			redundant := Fixed8FromFloat64(all).Sub(param.Fee)
			asset, _ := util.Uint256DecodeString(GasAssetHash)
			scripthash, _ := crypto.Uint160DecodeAddress(address)
			tx.Outputs = append(tx.Outputs, transaction.NewOutput(asset, redundant, scripthash))
		}
	}

	//vout
	for _, out := range param.Vout {
		value := out.GetFloat64("value")
		if value <= 0 {
			continue
		}
		all, inputs, err := getReference(out.GetString("asset"), out.GetString("address"), value, node)
		if err != nil {
			return err
		}
		tx.Inputs = append(tx.Inputs, inputs...)
		d, err := getAssetDecimals(node, out.GetString("asset"))
		if err != nil {
			return err
		}

		asset, err := util.Uint256DecodeString(out.GetString("asset"))
		if err != nil {
			return err
		}
		to, err := crypto.Uint160DecodeAddress(out.GetString("address"))
		if err != nil {
			return err
		}
		amount := BigDecimal{Value: value * math.Pow10(int(d)), Decimals: d}.ToFixed8()
		tx.Outputs = append(tx.Outputs, transaction.NewOutput(asset, amount, to))

		total := BigDecimal{Value: all * math.Pow10(int(d)), Decimals: d}.ToFixed8()
		if total.GreaterThan(amount) {
			initiator, _ := crypto.Uint160DecodeAddress(address)
			tx.Outputs = append(tx.Outputs, transaction.NewOutput(asset, total.Sub(amount), initiator))
		}
	}

	//attr
	return nil
}

func (tx *Tx) ToMap() goutil.Map {
	return nil
}

// EncodeHashableFields will only encode the fields that are not used for
// signing the transaction, which are all fields except the scripts.
func (tx *Tx) EncodeHashableFields() ([]byte, error) {
	w := &bytes.Buffer{}
	if err := binary.Write(w, binary.LittleEndian, tx.Type); err != nil {
		return nil, err
	}
	if err := binary.Write(w, binary.LittleEndian, tx.Version); err != nil {
		return nil, err
	}

	// Underlying TXer.
	if tx.Data != nil {
		if err := tx.Data.EncodeBinary(w); err != nil {
			return nil, err
		}
	}

	// Attributes
	lenAttrs := uint64(len(tx.Attributes))
	if err := util.WriteVarUint(w, lenAttrs); err != nil {
		return nil, err
	}
	for _, attr := range tx.Attributes {
		if err := attr.EncodeBinary(w); err != nil {
			return nil, err
		}
	}

	// Inputs
	if err := util.WriteVarUint(w, uint64(len(tx.Inputs))); err != nil {
		return nil, err
	}
	for _, in := range tx.Inputs {
		if err := in.EncodeBinary(w); err != nil {
			return nil, err
		}
	}

	// Outputs
	if err := util.WriteVarUint(w, uint64(len(tx.Outputs))); err != nil {
		return nil, err
	}
	for _, out := range tx.Outputs {
		if err := out.EncodeBinary(w); err != nil {
			return nil, err
		}
	}
	return w.Bytes(), nil
}

// RelayTx relay tx to the neo node
func RelayTx(tx *Tx, node string) error {
	return nil
}
