package neo

import (
	"encoding/hex"
	"github.com/CityOfZion/neo-go/pkg/util"
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
	return b.Value/math.Pow10(int(b.Decimals))
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