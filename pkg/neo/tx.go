package neo

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
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
	Script    []byte
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

func (tx *Tx) Label() string {
	if tx.Name != "" {
		return tx.Name
	}
	return tx.Hash().String()
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

func (tx *Tx) ParseInvoke(raw string) error {
	if len(tx.Param.Script) > 0 {
		return fmt.Errorf("script is already exitsted")
	}

	var items []goutil.Map
	err := json.Unmarshal([]byte(raw), &items)
	if err != nil {
		return err
	}

	var params []*ScriptParam
	for i := len(items) - 1; i >= 0; i-- {
		item := items[i]
		typ, ok := ParamTypeLookup[item.GetString("type")]
		if !ok {
			return fmt.Errorf("unsupport param type: %v", item.GetString("type"))
		}

		params = append(params, &ScriptParam{
			Type:  typ,
			Value: item.Get("value"),
		})
	}
	sb := NewScriptBuilder()
	err = sb.EmitParams(params)
	if err != nil {
		return err
	}
	tx.Param.Script = sb.Bytes()
	return nil
}

func (tx *Tx) ParseInvokeFunc(raw string) error {
	if len(tx.Param.Script) > 0 {
		return fmt.Errorf("script is already exitsted")
	}

	var items []interface{}
	err := json.Unmarshal([]byte(raw), &items)
	if err != nil {
		return err
	}

	var params []*ScriptParam
	args := &ScriptParam{
		Type:  ArrayType,
		Value: nil,
	}
	if len(items) >= 3 {
		args.Value = goutil.MapArrayV(items[2])
	}

	if len(items) < 2 {
		return fmt.Errorf("not enough arguments to invoke function")
	}

	params = append(params, args, &ScriptParam{
		Type:  StringType,
		Value: items[1],
	}, &ScriptParam{
		Type:  AppCallType,
		Value: items[0],
	})
	sb := NewScriptBuilder()
	err = sb.EmitParams(params)
	if err != nil {
		return err
	}
	tx.Param.Script = sb.Bytes()
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
			redundant := util.Fixed8(Fixed8FromFloat64(all) - Fixed8FromFloat64(fee))
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
		all, inputs, err := getReference(out.GetString("asset"), address, value, node)
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
		if total > amount {
			initiator, _ := crypto.Uint160DecodeAddress(address)
			tx.Outputs = append(tx.Outputs, transaction.NewOutput(asset, util.Fixed8(total-amount), initiator))
		}
	}

	//invoke script
	if len(param.Script) > 0 {
		tx.Data = &transaction.InvocationTX{Script: param.Script}
		if tx.Type != transaction.InvocationType {
			return fmt.Errorf("wrong tx type, should be invocation")
		}
	}
	//attr
	for _, attr := range param.Attr {
		tx.Attributes = append(tx.Attributes, &transaction.Attribute{
			Usage: attrLookup[attr.GetString("usage")],
			Data:  attr.Get("data").([]byte),
		})
	}

	encode, err := tx.EncodeHashableFields()
	if err != nil {
		return err
	}

	//witness
	for _, witness := range param.Witness {
		vScript, err := hex.DecodeString(witness.GetString("witness"))
		if err != nil {
			return err
		}
		var iScript []byte
		if IsSmartContract(vScript) {
			if len(witness.GetString("v")) > 0 {
				iScript, err = hex.DecodeString(witness.GetString("v"))
				if err != nil {
					return err
				}
			}
		} else { //private key
			privateKey, err := wallet.NewPrivateKeyFromBytes(vScript)
			if err != nil {
				return err
			}
			vScript, _ = PublicKeyScriptFromPrivateKey(privateKey)

			//sign
			iScript, err = privateKey.Sign(encode)
			if err != nil {
				return err
			}
			sb := NewScriptBuilder()
			err = sb.EmitBytes(iScript)
			if err != nil {
				return err
			}
			iScript = sb.Bytes()
		}

		tx.Scripts = append(tx.Scripts, &transaction.Witness{
			VerificationScript: vScript,
			InvocationScript:   iScript,
		})
	}

	return nil
}

func (tx *Tx) ToMap() goutil.Map {
	m := goutil.Struct2Map(tx)
	if m == nil {
		return m
	}
	delete(m, "Param")
	delete(m, "Name")
	if tx.Type == transaction.InvocationType {
		var script string
		var fee float64
		inv, ok := tx.Data.(*transaction.InvocationTX)
		if ok && inv != nil {
			script = hex.EncodeToString(inv.Script)
			fee = float64(inv.Gas) / 8
		}
		m.Set("script", script)
		m.Set("gas", fee)
	}
	m.Set("txid", tx.Hash())
	m.Set("size", tx.Size())
	m.Set("net_fee", strconv.FormatFloat(float64(tx.Param.Fee)/8, 'f', -1, 64))
	return m
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
	if tx == nil {
		return fmt.Errorf("tx is nil")
	}

	tx.Hash()

	w := new(bytes.Buffer)
	err := tx.EncodeBinary(w)
	if err != nil {
		return err
	}

	raw := hex.EncodeToString(w.Bytes())
	var res bool
	err = Rpc(node, "sendrawtransaction", []string{raw}, &res)
	if err != nil {
		return err
	}

	return nil
}
