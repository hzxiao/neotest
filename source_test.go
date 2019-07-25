package neotest

import (
	"github.com/hzxiao/goutil/assert"
	"testing"
)

func TestValidID(t *testing.T) {
	assert.True(t, ValidID("a"))
	assert.True(t, ValidID("a_b"))
	assert.True(t, ValidID("a1"))
	assert.True(t, ValidID("a_b_1"))
	assert.False(t, ValidID("@ab"))
	assert.False(t, ValidID("1ab"))
}