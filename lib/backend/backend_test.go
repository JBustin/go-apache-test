package backend

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/go-apache-test/mocks"
	"github.com/stretchr/testify/assert"
)

func Test_Handle(t *testing.T) {
	host := "www.gogole.com:3020"
	xForwardedHost := "www.myapp.com"
	path := "/foo/bar"
	req := http.Request{
		URL: &url.URL{
			Host: host,
			Path: path,
		},
		Host: host,
		Header: http.Header{
			"X-Forwarded-Host": []string{xForwardedHost},
		},
	}
	res := mocks.NewHttpResponse()
	handle(&res, &req)

	assert.Equal(t, 1, res.Calls["Header"], "Response header should be called 1 time")
	assert.Equal(t, 1, res.Calls["Write"], "Response writer should be called 1 time")
	assert.Equal(t, "application/json", res.Header().Get("Content-Type"))
	assert.Equal(t, "{\"Backend\":\"www.gogole.com\",\"Host\":\"www.myapp.com\",\"Pathname\":\"/foo/bar\",\"HttpStatusCode\":0,\"Body\":\"\"}", res.WriteArgs[0])
}
