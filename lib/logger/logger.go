package logger

import (
	"bytes"
	"log"
	"os"
)

const (
	ERROR = iota
	INFO  = iota
	DEBUG = iota
)

var LEVELS = map[string]int{
	"error": ERROR,
	"info":  INFO,
	"debug": DEBUG,
}

type Logger struct {
	Level int
	log   log.Logger
}

func NewLog(logLevel string) Logger {
	var buf bytes.Buffer
	level, exist := LEVELS[logLevel]
	if !exist {
		level = ERROR
	}
	return Logger{
		Level: level,
		log:   *log.New(&buf, "", log.Lmicroseconds),
	}
}

func (l Logger) Debug(msg string) {
	if l.Level < DEBUG {
		return
	}
	l.log.SetOutput(os.Stdout)
	l.log.Printf("%v", msg)
}

func (l Logger) Info(msg string) {
	if l.Level < INFO {
		return
	}
	l.log.SetOutput(os.Stdout)
	l.log.Printf("%v", msg)
}

func (l Logger) Error(msg string) {
	l.log.SetOutput(os.Stderr)
	l.log.Printf("%v", msg)
}
