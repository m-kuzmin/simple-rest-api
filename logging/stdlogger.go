package logging

import (
	"fmt"
	"log"
)

// StdLogger prints the logs using log.Printf.
type StdLogger struct{}

// Tracef implements Logger.
func (StdLogger) Tracef(s string, a ...any) {
	log.Printf("[trace      ] %s", fmt.Sprintf(s, a...)) //nolint:forbidigo // Allowed here only
}

// Debugf implements Logger.
func (StdLogger) Debugf(s string, a ...any) {
	log.Printf("[debug      ] %s", fmt.Sprintf(s, a...)) //nolint:forbidigo // Allowed here only
}

// Infof implements Logger.
func (StdLogger) Infof(s string, a ...any) {
	log.Printf("[INFO      i] %s", fmt.Sprintf(s, a...)) //nolint:forbidigo // Allowed here only
}

// Warnf implements Logger.
func (StdLogger) Warnf(s string, a ...any) {
	log.Printf("[WARNING   i] %s", fmt.Sprintf(s, a...)) //nolint:forbidigo // Allowed here only
}

// Errorf implements Logger.
func (StdLogger) Errorf(s string, a ...any) {
	log.Printf("[ERROR     E] %s", fmt.Sprintf(s, a...)) //nolint:forbidigo // Allowed here only
}

// Fatalf implements Logger.
func (StdLogger) Fatalf(s string, a ...any) {
	log.Fatalf("[FATAL     E] %s", fmt.Sprintf(s, a...)) //nolint:forbidigo // Allowed here only
}
