package config

import (
	"encoding/json"
	"fmt"

	"github.com/go-apache-test/lib/filesystem"
)

type Config struct {
	Test struct {
		InsideDocker   bool
		FetchTimeoutMs int
	}
	Debug struct {
		LogLevel string
	}
	Httpd struct {
		Backends     []string
		BackendsPort string
		FormatCmd    string
	}
}

var DefaultConfig Config = Config{
	Test: struct {
		InsideDocker   bool
		FetchTimeoutMs int
	}{
		InsideDocker:   true,
		FetchTimeoutMs: 3000,
	},
	Debug: struct {
		LogLevel string
	}{
		LogLevel: "info",
	},
	Httpd: struct {
		Backends     []string
		BackendsPort string
		FormatCmd    string
	}{
		Backends:     []string{},
		BackendsPort: "3011",
		FormatCmd:    "",
	},
}

func (c Config) String() string {
	return fmt.Sprintf(`
		FetchTimeoutMs		%v
		Debug.LogLevel		%v
		httpd.Backends		%v
		httpd.backendsPort 	%v
		httpd.FormatCmd		%v
	`, c.Test.FetchTimeoutMs,
		c.Debug.LogLevel,
		c.Httpd.Backends,
		c.Httpd.BackendsPort,
		c.Httpd.FormatCmd)
}

func New(jsonFilePath string, filer filesystem.Filer) (Config, error) {
	var c Config

	content, err := filer.Read(jsonFilePath)

	if err != nil {
		return c, err
	}

	c = DefaultConfig
	err = json.Unmarshal(content, &c)

	return c, err
}
