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
	assert.False(t, ValidID("a.b"))
}

func TestSplitExpr(t *testing.T) {
	assert.Equal(t, []string(nil), splitRawExpr(""))
	assert.Equal(t, []string{"x"}, splitRawExpr("x"))
	assert.Equal(t, []string{"x"}, splitRawExpr(" x "))
	assert.Equal(t, []string{"echo", "1", "2"}, splitRawExpr("echo 1 2"))
	assert.Equal(t, []string{"let", "@host", "`env HOST`"}, splitRawExpr("let @host `env HOST`"))
	assert.Equal(t, []string{"echo", "1", `"2 3"`}, splitRawExpr(`echo 1 "2 3"`))
	assert.Equal(t, []string{"let", "@a", "`base64 \"abc$(v)\"`"}, splitRawExpr("let @a `base64 \"abc$(v)\"`"))
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

func TestSplitCmd(t *testing.T) {
	s1 := `echo 1 2
echo 3 4
`
	assert.Equal(t, []string{"echo 1 2", "echo 3 4"}, splitRawCmd(s1))

	s2 := `echo '1234
56'`
	assert.Equal(t, []string{"echo '1234\n56'"}, splitRawCmd(s2))

	s3 := `echo '123' '123'`
	assert.Equal(t, []string{"echo '123' '123'"}, splitRawCmd(s3))

	s4 := `echo '123
456' '123
'`
	assert.Equal(t, []string{"echo '123\n456' '123\n'"}, splitRawCmd(s4))

	s5 := `echo '{
    "key": "value"
}'
`
	assert.Equal(t, []string{"echo '{\n    \"key\": \"value\"\n}'"}, splitRawCmd(s5))
}

func splitRawCmd(text string) []string {
	scanner := bufio.NewScanner(strings.NewReader(text))
	scanner.Split(splitCmd)

	var rawExprs []string
	for scanner.Scan() {
		rawExprs = append(rawExprs, scanner.Text())
	}

	return rawExprs
}