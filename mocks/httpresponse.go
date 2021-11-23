package mocks

import (
	"net/http"
)

type HttpResponse struct {
	Calls     map[string]int
	WriteArgs []string
	header    http.Header
}

func NewHttpResponse() HttpResponse {
	res := HttpResponse{}
	res.Reset()
	return res
}

func (res *HttpResponse) Reset() {
	res.Calls = map[string]int{
		"Header": 0,
		"Write":  0,
	}
	res.WriteArgs = []string{}
	res.header = http.Header{}
}

func (res *HttpResponse) Header() http.Header {
	res.Calls["Header"]++
	return res.header
}
func (res *HttpResponse) Write(data []byte) (int, error) {
	res.Calls["Write"]++
	res.WriteArgs = append(res.WriteArgs, string(data))
	return 0, nil
}
func (res HttpResponse) WriteHeader(statusCode int) {}
