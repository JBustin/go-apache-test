package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewLog(t *testing.T) {
	log := NewLog("debug")
	assert.Equal(t, DEBUG, log.Level, "should initialize log level to debug level")

	log = NewLog("info")
	assert.Equal(t, INFO, log.Level, "should initialize log level to info level")

	log = NewLog("error")
	assert.Equal(t, ERROR, log.Level, "should initialize log level to error level")

	log = NewLog("unknown")
	assert.Equal(t, ERROR, log.Level, "should initialize log level to error level")
}
