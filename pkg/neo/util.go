package neo

import (
	"encoding/hex"
	"fmt"
	"github.com/CityOfZion/neo-go/pkg/crypto"
	"github.com/CityOfZion/neo-go/pkg/util"
	"github.com/CityOfZion/neo-go/pkg/wallet"
	"math"
	"strconv"
	"strings"
)

type BigDecimal struct {
	Value    float64
	Decimals uint8
}

func (b BigDecimal) ToFixed8() util.Fixed8 {
	return util.Fixed8(b.ChangeDecimals(8).Value)
}

func (b BigDecimal) ChangeDecimals(decimals uint8) BigDecimal {
	if decimals == b.Decimals {
		return b
	}

	cb := BigDecimal{Decimals: decimals}
	if decimals > b.Decimals {
		cb.Value = math.Pow10(int(decimals-b.Decimals)) * b.Value
	} else {
		divisor := math.Pow10(int(b.Decimals - decimals))
		cb.Value = b.Value / divisor
	}
	return cb
}

func (b BigDecimal) RealValue() float64 {
	return b.Value / math.Pow10(int(b.Decimals))
}

func (b BigDecimal) String() string {
	v := strconv.FormatFloat(b.RealValue(), 'f', -1, 64)
	if strings.Contains(v, ".") {
		v = strings.TrimSuffix(strings.TrimSuffix(v, "0"), ".")
	}
	return v
}

func Fixed8FromFloat64(v float64) util.Fixed8 {
	s := strconv.FormatFloat(v, 'f', -1, 64)
	f, _ := util.Fixed8DecodeString(s)
	return f
}

func IsGlobalAsset(asset string) bool {
	bs, err := hex.DecodeString(asset)
	if err != nil {
		return false
	}
	return len(bs) == 32
}

func IsSmartContract(vScript []byte) bool {
	return len(vScript) == 20
}

//PublicKeyScriptFromPrivateKey
func PublicKeyScriptFromPrivateKey(privateKey *wallet.PrivateKey) ([]byte, error) {
	if privateKey == nil {
		return nil, fmt.Errorf("private key is nil")
	}

	pub, err := privateKey.PublicKey()
	if err != nil {
		return nil, err
	}
	return PublicKeyScript(pub), nil
}

//PublicKeyScript get public key script
func PublicKeyScript(pk *crypto.PublicKey) []byte {
	if pk == nil {
		return nil
	}

	b := pk.Bytes()
	b = append([]byte{0x21}, b...)
	b = append(b, 0xAC)
	return b
}
