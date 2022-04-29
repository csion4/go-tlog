package tLog

import (
	"io"
	"time"
)

const (
	Trace = iota + 1
	Debug
	Info
	Warn
	Error
)

type logger interface {
	trace(v string, t time.Time, h string) string
	debug(v string, t time.Time, h string) string
	info(v string, t time.Time, h string) string
	warn(v string, t time.Time, h string) string
	error(v string, t time.Time, h string) string
}

type stdLogger struct {
	out io.Writer
	level int
}
