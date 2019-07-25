package neotest

import (
	"bufio"
	"github.com/hzxiao/goutil/assert"
	"strings"
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

func TestSplitExpr(t *testing.T) {
	assert.Equal(t, []string(nil), splitRawExpr(""))
	assert.Equal(t, []string{"x"}, splitRawExpr("x"))
	assert.Equal(t, []string{"x"}, splitRawExpr(" x "))
	assert.Equal(t, []string{"echo", "1", "2"}, splitRawExpr("echo 1 2"))
	assert.Equal(t, []string{"let", "@host", "`env HOST`"}, splitRawExpr("let @host `env HOST`"))
}

func splitRawExpr(text string) []string {
	scanner := bufio.NewScanner(strings.NewReader(text))
	scanner.Split(splitExpr)

	var rawExprs []string
	for scanner.Scan() {
		rawExprs = append(rawExprs, scanner.Text())
	}

	return rawExprs
}