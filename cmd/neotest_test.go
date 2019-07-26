package main

import (
	"github.com/hzxiao/goutil/assert"
	"testing"
)

func TestRun(t *testing.T)  {
	err := run([]string{"../testdata/echo.ntf"})
	assert.NoError(t, err)
}

func TestLetCmd(t *testing.T)  {
	err := run([]string{"../testdata/let.ntf"})
	assert.NoError(t, err)
}