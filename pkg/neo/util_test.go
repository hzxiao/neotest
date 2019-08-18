package neo

import (
	"github.com/hzxiao/goutil/assert"
	"testing"
)

func TestBigDecimal_ChangeDecimals(t *testing.T) {
	b := BigDecimal{100, 2}
	cb := b.ChangeDecimals(5)
	assert.Equal(t, uint8(5), cb.Decimals)
	assert.Equal(t, float64(100000), cb.Value)
	assert.Equal(t, float64(1), cb.RealValue())
}

func TestFixed8FromFloat64(t *testing.T) {
	f := Fixed8FromFloat64(168612.94473278)
	assert.Equal(t, "16861294473278", f.String())
}