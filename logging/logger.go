package logging

type LogLevel int

const (
	TraceLevel LogLevel = iota
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

type Logger interface {
	Tracef(string, ...any)
	Debugf(string, ...any)
	Infof(string, ...any)
	Warnf(string, ...any)
	Errorf(string, ...any)
	Fatalf(string, ...any) // Must always panic
}

/*
GlobalLogger is used in package level functions and serves as the default logging implementation. By default it is set
to NilLogger. The value should never be `nil`. If it is nil then package level logging funcs will panic. If you want to
silence this logger then you should set the value to NilLogger.
*/
var GlobalLogger Logger //nolint:gochecknoglobals // This is a stateless global logger for the entire app.

func init() { //nolint:gochecknoinits // For logging its ok to have global stuff
	GlobalLogger = NilLogger{}
}

// Tracef uses GlobalLogger to record a trace log.
func Tracef(s string, a ...any) {
	GlobalLogger.Tracef(s, a...)
}

// Debugf uses GlobalLogger to record a debug log.
func Debugf(s string, a ...any) {
	GlobalLogger.Debugf(s, a...)
}

// Infof uses GlobalLogger to record a debug log.
func Infof(s string, a ...any) {
	GlobalLogger.Infof(s, a...)
}

// Warnf uses GlobalLogger to record a debug log.
func Warnf(s string, a ...any) {
	GlobalLogger.Warnf(s, a...)
}

// Errorf uses GlobalLogger to record a debug log.
func Errorf(s string, a ...any) {
	GlobalLogger.Errorf(s, a...)
}

// Fatalf uses GlobalLogger to record a debug log.
func Fatalf(s string, a ...any) {
	GlobalLogger.Fatalf(s, a...)
}
