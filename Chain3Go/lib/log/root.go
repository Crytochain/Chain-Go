package log

import (
	"os"
	"fmt"
)

var (
	root          = &logger{[]interface{}{}, new(swapHandler)}
	StdoutHandler = StreamHandler(os.Stdout, LogfmtFormat())
	StderrHandler = StreamHandler(os.Stderr, LogfmtFormat())
)

func init() {
	root.SetHandler(DiscardHandler())
}

// New returns a new logger with the given context.
// New is a convenient alias for Root().New
func New(ctx ...interface{}) Logger {
	return root.New(ctx...)
}

// Root returns the root logger
func Root() Logger {
	return root
}

// The following functions bypass the exported logger methods (logger.Debug,
// etc.) to keep the call depth the same for all paths to logger.write so
// runtime.Caller(2) always refers to the call site in client code.

// Trace is a convenient alias for Root().Trace
func Trace(msg string, ctx ...interface{}) {
	root.write(msg, LvlTrace, ctx)
}

// Debug is a convenient alias for Root().Debug
func Debug(msg string, ctx ...interface{}) {
	s := GoIDStr() + ":" + msg
	root.write(s, LvlDebug, ctx)
}

func Debugf(msg string, a ...interface{}) {
	s := GoIDStr() + ":"
	if a == nil {
		s += msg
	} else {
		msgString := fmt.Sprintf(msg, a ...)
		s += msgString
	}
	root.write(s, LvlDebug, nil)
}

// Info is a convenient alias for Root().Info
func Info(msg string, ctx ...interface{}) {
	s := GoIDStr() + ":" + msg
	root.write(s, LvlInfo, ctx)
}

func Infof(msg string, a ...interface{}) {
	s := GoIDStr() + ":"
	if a == nil {
		s += msg
	} else {
		msgString := fmt.Sprintf(msg, a ...)
		s += msgString
	}
	root.write(s, LvlInfo, nil)
}

// Warn is a convenient alias for Root().Warn
func Warn(msg string, ctx ...interface{}) {
	root.write(msg, LvlWarn, ctx)
}

// Error is a convenient alias for Root().Error
func Error(msg string, ctx ...interface{}) {
	root.write(msg, LvlError, ctx)
}

func Errorf(msg string, a ...interface{}) {
	s := GoIDStr() + ":"
	if a == nil {
		s += msg
	} else {
		msgString := fmt.Sprintf(msg, a ...)
		s += msgString
	}
	root.write(s, LvlError, nil)
}

// Crit is a convenient alias for Root().Crit
func Crit(msg string, ctx ...interface{}) {
	root.write(msg, LvlCrit, ctx)
	os.Exit(1)
}
