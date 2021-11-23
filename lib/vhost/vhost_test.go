package vhost

import (
	"testing"

	"github.com/go-apache-test/mocks"
	"github.com/stretchr/testify/assert"
)

var vhostContent string = `<VirtualHost *:80>
	ServerName www.mydomain.com
	ServerAlias www.mydomain.org fr.mydomain.fr us.mydomain.com

	Alias /foo /bar
</VirtualHost>
`

var testContent string = `{
  "suites": [
    {
      "name": "200, balance to backend",
      "expectations": { "backend": "backend.com", "httpStatusCode": 200 },
      "rules": [
        { "host": "www.mydomain.com", "pathname": "/flags/foo" },
        { "host": "www.mydomain.com", "pathname": "/template/foo" }
	  ]
	}
  ]
}
`

func Test_GetAlias(t *testing.T) {
	v := Vhost{
		OriPath: "vhosts/mydomain.conf",
	}

	fs := mocks.NewFilesystem()
	fs.SetStubRead(func(filepath string) ([]byte, error) {
		return []byte(vhostContent), nil
	})
	err := v.GetAlias(&fs)

	assert.Equal(t, nil, err, "should not return an error")
	assert.Equal(t, 1, fs.Calls["Exists"], "should call fs.exists")
	assert.Equal(t, 1, fs.Calls["Read"], "should call fs.read")
	assert.Equal(t, []string{
		"www.mydomain.org",
		"fr.mydomain.fr",
		"us.mydomain.com",
		"www.mydomain.com",
	}, v.Alias, "should get all alias")
}

func Test_GetSuites(t *testing.T) {
	v := Vhost{
		TestPath: "tests/mydomain.conf.json",
	}

	fs := mocks.NewFilesystem()
	fs.SetStubRead(func(filepath string) ([]byte, error) {
		return []byte(testContent), nil
	})
	err := v.GetSuites(&fs)

	assert.Equal(t, nil, err, "should not return an error")
	assert.Equal(t, 1, fs.Calls["Exists"], "should call fs.exists")
	assert.Equal(t, 1, fs.Calls["Read"], "should call fs.read")
	assert.Equal(t, 1, len(v.Suites), "should populate the suites")
	assert.Equal(t, 2, len(v.Suites[0].Rules), "should populate the rules")
	assert.Equal(t, "backend.com", v.Suites[0].Expectations.Backend, "should populate the suite expectations")
	assert.Equal(t, 200, v.Suites[0].Expectations.HttpStatusCode, "should populate the suite expectations")
}
