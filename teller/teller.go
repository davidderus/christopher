package teller

import (
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

// Teller is a logrus abstraction
type Teller struct {
	logger *log.Logger
}

// NewTeller introduces a Teller instance with a given log level and formatter
func NewTeller(logLevel, logFormatter string) *Teller {
	teller := &Teller{}

	teller.logger = log.New()
	teller.SetLogLevel(logLevel)
	teller.SetLogFormatter(logFormatter)
	teller.SetLogOutput(nil)

	return teller
}

// SetLogLevel defines the logger minimum level from a string
//
// Note that logrus default is log.InfoLevel
func (t *Teller) SetLogLevel(logLevel string) {
	switch logLevel {
	case "fatal":
		t.logger.Level = log.FatalLevel
	case "error":
		t.logger.Level = log.ErrorLevel
	case "warn":
		t.logger.Level = log.WarnLevel
	case "info":
		t.logger.Level = log.InfoLevel
	case "debug":
		t.logger.Level = log.DebugLevel
	}
}

// SetLogFormatter sets the log formatter
//
// Both text and json formatters are accepted
func (t *Teller) SetLogFormatter(logFormatter string) {
	switch logFormatter {
	case "json":
		t.logger.Formatter = &log.JSONFormatter{}
	default:
		t.logger.Formatter = &log.TextFormatter{}
	}
}

// SetLogOutput defines the logger output writer
//
// Default is os.Stdout
func (t *Teller) SetLogOutput(logOutput io.Writer) {
	if logOutput == nil {
		t.logger.Out = os.Stdout
	} else {
		t.logger.Out = logOutput
	}
}

// Log allows message loging
func (t *Teller) Log() *log.Logger {
	return t.logger
}

// LogWithFields allow logging some free fields associated to a message
func (t *Teller) LogWithFields(fields map[string]interface{}) *log.Entry {
	return t.logger.WithFields(fields)
}
