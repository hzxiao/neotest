package neo

import (
	"encoding/hex"
	"github.com/hzxiao/goutil/assert"
	"testing"
)

func TestScriptBuilder_EmitParams(t *testing.T) {
	params := []*ScriptParam{
		{
			Type: OpCodeType,
			Value: 0xf1,
		},
		{
			Type:  AppCallType,
			Value: "0x02196f55f618cfb34e80bed272f2f3faaeba131e",
		},
		{
			Type:  StringType,
			Value: "transfer",
		},
		{
			Type: ArrayType,
			Value: []*ScriptParam{
				{
					Type:  AddressType,
					Value: "AWSuQXpjuY3v22gCbEFL2vHbSLMMVK1QD6",
				},
				{
					Type:  AddressType,
					Value: "AdP3gUNRXqg4EVVSQD4o1i3kfF9DmNQSw1",
				},
				{
					Type:  IntegerType,
					Value: 1000000,
				},
			},
		},
	}

	for i := 0; i < len(params)/2; i++ {
		j := len(params) - 1 - i
		params[i], params[j] = params[j], params[i]
	}
	sb := NewScriptBuilder()
	err := sb.EmitParams(params)
	assert.NoError(t, err)

	t.Log(hex.EncodeToString(sb.Bytes()))
}
