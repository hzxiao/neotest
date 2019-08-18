package main

import (
	"github.com/hzxiao/goutil/assert"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func init()  {
	verbose = true
}

func TestRun(t *testing.T)  {
	err := run([]string{"../testdata/echo.ntf"})
	assert.NoError(t, err)
}

func TestLetCmd(t *testing.T)  {
	err := run([]string{"../testdata/let.ntf"})
	assert.NoError(t, err)
}

func TestSubCmd(t *testing.T)  {
	err := run([]string{"../testdata/sub_cmd.ntf"})
	assert.NoError(t, err)
}

func TestEqualCmd(t *testing.T)  {
	err := run([]string{"../testdata/equal.ntf"})
	assert.NoError(t, err)
}

func TestReqCmd(t *testing.T)  {
	var hello = func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte("hi"))
	}

	var foo = func(w http.ResponseWriter, r *http.Request) {
		buf, err := ioutil.ReadAll(r.Body)
		assert.NoError(t, err)
		defer r.Body.Close()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(buf)
	}

	http.HandleFunc("/hello", hello)
	http.HandleFunc("/foo", foo)

	go func() {
		err := http.ListenAndServe(":10000", nil)
		assert.Error(t, err)
	}()

	time.Sleep(2*time.Second)
	err := run([]string{"../testdata/http.ntf"})
	assert.NoError(t, err)
}

func TestTx(t *testing.T)  {
	err := run([]string{"../testdata/tx.ntf"})
	assert.NoError(t, err)
}