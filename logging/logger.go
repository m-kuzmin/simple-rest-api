package logging

type Logger interface {
	Tracef(string, ...any)
	Debugf(string, ...any)
	Infof(string, ...any)
	Warnf(string, ...any)
	Errorf(string, ...any)
	Fatalf(string, ...any)
}

/*
GlobalLogger is used in package level functions and serves as the default logging implementation. nil by
default meaning package level log functions discard the messages.
*/
var GlobalLogger Logger //nolint:gochecknoglobals // This is a stateless global logger for the entire app.

// Tracef uses GlobalLogger to record a trace log.
func Tracef(s string, a ...any) {
	if GlobalLogger != nil {
		GlobalLogger.Tracef(s, a...)
	}
}

// Debugf uses GlobalLogger to record a debug log.
func Debugf(s string, a ...any) {
	if GlobalLogger != nil {
		GlobalLogger.Debugf(s, a...)
	}
}

// Infof uses GlobalLogger to record a debug log.
func Infof(s string, a ...any) {
	if GlobalLogger != nil {
		GlobalLogger.Infof(s, a...)
	}
}

// Warnf uses GlobalLogger to record a debug log.
func Warnf(s string, a ...any) {
	if GlobalLogger != nil {
		GlobalLogger.Warnf(s, a...)
	}
}

// Errorf uses GlobalLogger to record a debug log.
func Errorf(s string, a ...any) {
	if GlobalLogger != nil {
		GlobalLogger.Errorf(s, a...)
	}
}

// Fatalf uses GlobalLogger to record a debug log.
func Fatalf(s string, a ...any) {
	if GlobalLogger != nil {
		GlobalLogger.Fatalf(s, a...)
	}
}
