package logging

// Tape records logs and dumps them only when a log of some level or higher is printed.
type Tape struct {
	Suppressed   Logger
	Unsuppressed Logger
	Logs         []LogEntry
	Dump         LogLevel
	Suppress     LogLevel
}

type LogEntry func(Logger)

/*
NewTape creates a special logger that allows you to have less cluttered logs if everything is fine. But if something
goes wrong you can see more details. If both suppressed and unsuppressed are nil or NilLogger then returns NilLogger.

- All logs below or equal to suppress level are stored for later.
- When a log of level higher than dump is logged all stored suppress logs are printed using suppressed logger.
- All logs above suppress level are immediately printed using unsuppressed logger.

Here is an example:

1. Debugf("Some small detail, not important")
2. ... more small logs, only useful when debugging
3. Errorf("Some error")
4. All debug logs are printed to assist in diagnosing the error.

If the Errorf was never called and you dont need to know those small debug statements then you can discard the tape.
*/
func NewTape(suppress LogLevel, suppressed Logger, dump LogLevel, unsuppressed Logger) Logger {
	if suppressed == nil {
		suppressed = NilLogger{}
	}

	if unsuppressed == nil {
		unsuppressed = NilLogger{}
	}

	return &Tape{
		Suppressed:   suppressed,
		Unsuppressed: unsuppressed,
		Logs:         []LogEntry{},
		Dump:         dump,
		Suppress:     suppress,
	}
}

func (t *Tape) decide(level LogLevel, levelFunc LogEntry) {
	if level <= t.Suppress {
		t.Logs = append(t.Logs, levelFunc)
		return
	}

	if level >= t.Dump {
		t.DumpTape(levelFunc)
		return
	}

	levelFunc(t.Unsuppressed)
}

func (t *Tape) DumpTape(header LogEntry) {
	header(t.Unsuppressed)

	for _, log := range t.Logs {
		log(t.Suppressed)
	}

	t.Logs = t.Logs[:0]
}

// Tracef implements Logger.
func (t *Tape) Tracef(fmtStr string, a ...any) {
	t.decide(TraceLevel, func(l Logger) { l.Tracef(fmtStr, a...) })
}

// Debugf implements Logger.
func (t *Tape) Debugf(fmtStr string, a ...any) {
	t.decide(DebugLevel, func(l Logger) { l.Debugf(fmtStr, a...) })
}

// Infof implements Logger.
func (t *Tape) Infof(fmtStr string, a ...any) {
	t.decide(InfoLevel, func(l Logger) { l.Infof(fmtStr, a...) })
}

// Warnf implements Logger.
func (t *Tape) Warnf(fmtStr string, a ...any) {
	t.decide(WarnLevel, func(l Logger) { l.Warnf(fmtStr, a...) })
}

// Errorf implements Logger.
func (t *Tape) Errorf(fmtStr string, a ...any) {
	t.decide(ErrorLevel, func(l Logger) { l.Errorf(fmtStr, a...) })
}

// Fatalf implements Logger.
func (t *Tape) Fatalf(fmtStr string, a ...any) {
	t.decide(FatalLevel, func(l Logger) { l.Fatalf(fmtStr, a...) })
}
