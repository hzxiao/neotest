package neotest

import (
	"github.com/hzxiao/goutil/httputil"
	"io"
	"strings"
)

type HttpRequest struct {
	Method string
	URL    string
	Body   string
}

//Send send http request
func (r *HttpRequest) Send() (code int, header map[string]string, body interface{}, err error) {
	hds := make(map[string]string)
	var data io.Reader
	if r.Body != "" {
		hds["Content-Type"] = "application/json"
		data = strings.NewReader(r.Body)
	}
	var result string
	code, header, err = httputil.Request(r.Method, r.URL, hds, data, httputil.ReturnString, &result)
	if err != nil {
		return
	}
	body = result
	return
}
