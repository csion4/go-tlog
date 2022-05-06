package tLog

import (
	"os"
	"time"
)

const (
	confConsole = "console"
)

func getStdOutLogger() (sol *stdOutLogger) {
	conf := getTlogConf()
	var level string
	if conf != nil {
		m := conf[confConsole]
		level = m[confLevel]
	}
	var iLevel int
	if level == "" {
		iLevel = getDefaultLevel()
	} else {
		iLevel = switchLevel(level)
	}
	return &stdOutLogger{
		&stdLogger{
			level: iLevel,
		},
	}
}


// ------  stdOutLogger  ------
type stdOutLogger struct {
	*stdLogger
}

func (do *stdOutLogger) trace(v string, t time.Time, h string) string  {
	if do.level <= Trace {
		if h == "" {
			h = formatHeader(t, "trace", "")
		}
		_, _ = os.Stdout.Write([]byte(h + v))
	}
	return h
}
func (do *stdOutLogger) debug(v string, t time.Time, h string) string  {
	if do.level <= Debug {
		if h == "" {
			h = formatHeader(t, "debug", "")
		}
		_, _ = os.Stdout.Write([]byte(h + v))
	}
	return h
}
func (do *stdOutLogger) info(v string, t time.Time, h string) string  {
	if do.level <= Info {
		if h == "" {
			h = formatHeader(t, "info", "")
		}
		_, _ = os.Stdout.Write([]byte(h + v))
	}
	return h
}
func (do *stdOutLogger) warn(v string, t time.Time, h string) string  {
	if do.level <= Warn {
		if h == "" {
			h = formatHeader(t, "warn", "")
		}
		_, _ = os.Stderr.Write([]byte(h + v))
	}
	return h
}
func (do *stdOutLogger) error(v string, t time.Time, h string) string  {
	if do.level <= Error {
		if h == "" {
			h = formatHeader(t, "error", "")
		}
		_, _ = os.Stderr.Write([]byte(h + v))
	}
	return h
}
