package worker

import (
	"testing"

	"github.com/go-apache-test/lib/config"
	"github.com/go-apache-test/lib/cst"
	"github.com/go-apache-test/lib/logger"
	"github.com/go-apache-test/lib/vhost"
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

func Test_Init(t *testing.T) {
	fs := mocks.NewFilesystem()
	w := worker{
		filer:  &fs,
		config: config.DefaultConfig,
		vhosts: []vhost.Vhost{},
		logger: logger.NewLog("error"),
	}

	index := 0
	fs.SetStubWrite(func(filepath string, content string) error {
		index++
		if index == 1 {
			assert.Equal(t, cst.HttpdPortFile, filepath, "should write in expected filepath")
			assert.Equal(t, "Listen 80\n", content, "should write the expected content")
		}
		if index == 2 {
			assert.Equal(t, cst.DockerHostsFile, filepath, "should write in expected filepath")
			assert.Equal(t, "\n127.0.0.1 \n", content, "should write the expected content")
		}
		return nil
	})
	fs.SetStubAppend(func(filepath string, content string) error {
		assert.Equal(t, cst.HttpdConfFile, filepath, "should write in expected filepath")
		assert.Equal(t, "\nServerName http.front.com\n", content, "should write the expected content")
		return nil
	})

	err := w.Init()

	assert.Equal(t, nil, err, "should not return an error")
	assert.Equal(t, 1, fs.Calls["List"], "should call fs.List 1 time")
	assert.Equal(t, 1, fs.Calls["Append"], "should call fs.Append 1 times")
	assert.Equal(t, 2, fs.Calls["Write"], "should call fs.Write 2 times")
}

func Test_LoadVhosts(t *testing.T) {
	fs := mocks.NewFilesystem()
	conf := config.DefaultConfig
	conf.Test.InsideDocker = false
	w := worker{
		filer:  &fs,
		config: conf,
		vhosts: []vhost.Vhost{},
		logger: logger.NewLog("error"),
	}

	fs.StubList = []string{
		"www.mydomain1.com",
		"www.mydomain2.com",
		"www.mydomain3.com",
	}

	index := 0
	fs.StubList = []string{"mydomain.conf", "myotherdomain.conf"}
	fs.SetStubRead(func(filepath string) ([]byte, error) {
		index++
		// first call getAlias
		if index == 1 {
			return []byte(vhostContent), nil
		}
		// second call getSuites
		return []byte(testContent), nil
	})

	err := w.loadVhosts()

	assert.Equal(t, nil, err, "should not return an error")
	assert.Equal(t, 1, fs.Calls["List"], "should call fs.List 1 time")
	assert.Equal(t, 4, fs.Calls["Read"], "should call fs.Read 4 times")
	assert.Equal(t, 2, len(w.vhosts), "should load two vhosts")
	assert.Equal(t, 1, len(w.vhosts[0].Suites), "should load the suites")
	assert.Equal(t, 2, len(w.vhosts[0].Suites[0].Rules), "should load the rules")
}
