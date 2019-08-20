package neo

import (
	"bytes"
	"fmt"
	"github.com/CityOfZion/neo-go/pkg/crypto"
	"github.com/CityOfZion/neo-go/pkg/util"
	"github.com/CityOfZion/neo-go/pkg/vm"
	"github.com/hzxiao/goutil"
	"strings"
)

var ErrWrongValueOfParamType = fmt.Errorf("wrong value of the param type")

var ParamTypeLookup = map[string]ParamType{
	"Signature": SignatureType,
	"Boolean":   BoolType,
	"Integer":   IntegerType,
	"Hash160":   Hash160Type,
	"Hash256":   Hash256Type,
	"ByteArray": ByteArrayType,
	"PublicKey": PublicKeyType,
	"String":    StringType,
	"Array":     ArrayType,
	"AppCall":   AppCallType,
	"Address":   AddressType,
	"OpCode":    OpCodeType,
}

// ParamType represent the Type of the contract parameter
type ParamType int

// A list of supported smart contract parameter types.
const (
	SignatureType ParamType = iota
	BoolType
	IntegerType
	Hash160Type
	Hash256Type
	ByteArrayType
	PublicKeyType
	StringType
	ArrayType
	AppCallType
	AddressType
	OpCodeType
)

type ScriptParam struct {
	// Type of the parameter
	Type ParamType `json:"type"`
	// The actual value of the parameter.
	Value interface{} `json:"value"`
}

type ScriptBuilder struct {
	buf *bytes.Buffer
}

func NewScriptBuilder() *ScriptBuilder {
	return &ScriptBuilder{
		buf: &bytes.Buffer{},
	}
}

func (s *ScriptBuilder) EmitParams(params []*ScriptParam) error {
	for _, param := range params {
		err := s.EmitParam(*param)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *ScriptBuilder) EmitParam(params ScriptParam) (err error) {
	switch params.Type {
	case SignatureType:
	case BoolType:
		v, ok := params.Value.(bool)
		if !ok {
			return ErrWrongValueOfParamType
		}
		return s.EmitBool(v)
	case IntegerType:
		v, err := goutil.Int64E(params.Value)
		if err != nil {
			return err
		}
		return s.EmitInt(v)
	case Hash160Type:
		var hash util.Uint160
		switch v := params.Value.(type) {
		case string:
			hash, err = util.Uint160DecodeString(strings.TrimPrefix(v, "0x"))
		case []byte:
			hash, err = util.Uint160DecodeBytes(v)
		default:
			return ErrWrongValueOfParamType
		}
		if err != nil {
			return err
		}
		return s.EmitBytes(hash.Bytes())
	case Hash256Type:
		var hash util.Uint256
		switch v := params.Value.(type) {
		case string:
			hash, err = util.Uint256DecodeString(strings.TrimPrefix(v, "0x"))
		case []byte:
			hash, err = util.Uint256DecodeBytes(v)
		default:
			return ErrWrongValueOfParamType
		}
		if err != nil {
			return err
		}
		return s.EmitBytes(hash.Bytes())
	case ByteArrayType:

	case PublicKeyType:
	case StringType:
		v, ok := params.Value.(string)
		if !ok {
			return ErrWrongValueOfParamType
		}
		return s.EmitString(v)
	case ArrayType:
		var items []*ScriptParam
		var ok bool
		if items, ok = params.Value.([]*ScriptParam); !ok {
			arr := goutil.MapArrayV(params.Value)
			for _, a := range arr {
				typ, ok := ParamTypeLookup[a.GetString("type")]
				if !ok {
					return fmt.Errorf("unsupport param type: %v", a.GetString("type"))
				}
				var item = ScriptParam{
					Type:  typ,
					Value: a.Get("value"),
				}
				items = append(items, &item)
			}
		}
		for i := len(items) - 1; i >= 0; i-- {
			err = s.EmitParam(*items[i])
			if err != nil {
				return err
			}
		}
		s.EmitInt(int64(len(items)))
		s.EmitOpCode(vm.Opack)
	case AppCallType:
		v, ok := params.Value.(string)
		if !ok {
			return ErrWrongValueOfParamType
		}
		v = strings.TrimPrefix(v, "0x")
		var hash util.Uint160
		if len(v) == 40 {
			hash, err = util.Uint160DecodeString(v)
		} else {
			hash, err = crypto.Uint160DecodeAddress(v)
		}
		if err != nil {
			return err
		}
		return s.EmitAppCall(hash, false)
	case AddressType:
		v, ok := params.Value.(string)
		if !ok {
			return ErrWrongValueOfParamType
		}
		hash, err := crypto.Uint160DecodeAddress(v)
		if err != nil {
			return err
		}
		return s.EmitBytes(hash.Bytes())
	case OpCodeType:
		opcode := goutil.Int64(params.Value)
		return s.EmitOpCode(vm.Opcode(opcode))
	default:
		return fmt.Errorf("unknown param type")
	}
	return nil
}

func (s *ScriptBuilder) Emit(op vm.Opcode, b []byte) error {
	return vm.Emit(s.buf, op, b)
}

func (s *ScriptBuilder) EmitOpCode(op vm.Opcode) error {
	return vm.EmitOpcode(s.buf, op)
}

func (s *ScriptBuilder) EmitBool(b bool) error {
	return vm.EmitBool(s.buf, b)
}

func (s *ScriptBuilder) EmitInt(i int64) error {
	return vm.EmitInt(s.buf, i)
}

func (s *ScriptBuilder) EmitString(str string) error {
	return vm.EmitString(s.buf, str)
}

func (s *ScriptBuilder) EmitBytes(b []byte) error {
	return vm.EmitBytes(s.buf, b)
}

func (s *ScriptBuilder) EmitSyscall(api string) error {
	return vm.EmitSyscall(s.buf, api)
}

func (s *ScriptBuilder) EmitAppCall(scriptHash util.Uint160, tailCall bool) error {
	little, _ := util.Uint160DecodeBytes(scriptHash.BytesReverse())
	return vm.EmitAppCall(s.buf, little, tailCall)
}

func (s *ScriptBuilder) Bytes() []byte {
	return s.buf.Bytes()
}
