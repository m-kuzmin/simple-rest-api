package logging

import "fmt"

// NilLogger ignores all Logger function calls, except Logger.Fatalf() which panic()s
type NilLogger struct{}

// Tracef implements Logger.
func (l NilLogger) Tracef(string, ...any) {}

// Debugf implements Logger.
func (l NilLogger) Debugf(string, ...any) {}

// Infof implements Logger.
func (l NilLogger) Infof(string, ...any) {}

// Warnf implements Logger.
func (l NilLogger) Warnf(string, ...any) {}

// Errorf implements Logger.
func (l NilLogger) Errorf(string, ...any) {}

// Fatalf implements Logger.
func (l NilLogger) Fatalf(fmtStr string, a ...any) {
	panic(fmt.Sprintf(fmtStr, a...))
}
