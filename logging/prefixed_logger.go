package logging

type PrefixedLogger struct {
	Logger Logger
	Prefix string
}

// NewPrefixedLogger creates a logger that will prepend "prefix " to all log messages
func NewPrefixedLogger(logger Logger, prefix string) Logger {
	if _, isNil := logger.(NilLogger); logger == nil || isNil {
		return NilLogger{}
	}

	return &PrefixedLogger{
		Prefix: prefix,
		Logger: logger,
	}
}

// Tracef implements Logger.
func (l *PrefixedLogger) Tracef(fmtStr string, a ...any) {
	l.Logger.Tracef(l.Prefix+" "+fmtStr, a...)
}

// Debugf implements Logger.
func (l *PrefixedLogger) Debugf(fmtStr string, a ...any) {
	l.Logger.Debugf(l.Prefix+" "+fmtStr, a...)
}

// Infof implements Logger.
func (l *PrefixedLogger) Infof(fmtStr string, a ...any) {
	l.Logger.Infof(l.Prefix+" "+fmtStr, a...)
}

// Warnf implements Logger.
func (l *PrefixedLogger) Warnf(fmtStr string, a ...any) {
	l.Logger.Warnf(l.Prefix+" "+fmtStr, a...)
}

// Errorf implements Logger.
func (l *PrefixedLogger) Errorf(fmtStr string, a ...any) {
	l.Logger.Errorf(l.Prefix+" "+fmtStr, a...)
}

// Fatalf implements Logger.
func (l *PrefixedLogger) Fatalf(fmtStr string, a ...any) {
	l.Logger.Fatalf(l.Prefix+" "+fmtStr, a...)
}
