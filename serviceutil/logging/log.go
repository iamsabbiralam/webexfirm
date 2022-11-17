package logging

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	isTerminal = terminal.IsTerminal(int(os.Stdout.Fd()))
	isTest     = strings.HasSuffix(os.Args[0], ".test")
	flagMu     = sync.Mutex{}
)

// NewLogger creates a new logrus logger with a fluentd formatter that will
// work natively with stackdriver logging
func NewLogger() *logrus.Logger {
	// Prepare a new logger
	logger := logrus.New()

	logger.Level = logrus.InfoLevel
	if isTest {
		// testing.Verbose can panic if Init hasn't been called yet,
		// e.g. if NewLogger is used as part of a global declaration.
		testing.Init()

		// Also parse the test flags, to be able to query -test.v. This
		// won't affect release binaries, because of the isTest check.
		// Do this behind a mutex to avoid data races, and don't parse
		// the flags twice.
		flagMu.Lock()
		if !flag.Parsed() {
			flag.Parse()
		}
		flagMu.Unlock()

		if !testing.Verbose() {
			// Keep the tests quiet, unless -test.v is used.
			logger.Level = logrus.FatalLevel
		}
	}

	logger.Out = os.Stdout
	if isTerminal {
		logger.Formatter = &logrus.TextFormatter{
			TimestampFormat: time.RFC3339,
			FullTimestamp:   true,
		}
	} else {
		logger.Formatter = &FluentdFormatter{
			TimestampFormat: time.RFC3339,
		}
	}

	logger.SetReportCaller(true)
	return logger
}

// WithError takes an error and logger and returns a standardised error logger
func WithError(err error, logger logrus.FieldLogger) *logrus.Entry {
	return logger.WithError(err).WithField("stacktrace", fmt.Sprintf("%+v", err))
}
