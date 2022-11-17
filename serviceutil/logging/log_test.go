package logging

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
)

func TestNewLogger(t *testing.T) {
	logger := NewLogger()

	if logger.Out != os.Stdout {
		t.Errorf("logger.Out is not os.Stdout")
	}
}

func TestWithError(t *testing.T) {
	logger := NewLogger().WithFields(logrus.Fields{
		"field1":      "value1",
		"field.key.2": "value2",
	})

	err := fmt.Errorf("This is an error")
	logger = WithError(err, logger)

	if len(logger.Data) != 4 {
		t.Errorf("len(logger.Data)(%d) does not equal 4", len(logger.Data))
	}

	v := logger.Data["field1"].(string)
	if v != "value1" {
		t.Errorf("logger.Data[field1](%s) does not equal 'value1'", v)
	}
	v = logger.Data["field.key.2"].(string)
	if v != "value2" {
		t.Errorf("logger.Data[field.key.2](%s) does not equal 'value2'", v)
	}
	if _, ok := logger.Data["error"].(error); !ok {
		t.Errorf("'error' key missing in logging data")
	}
	if _, ok := logger.Data["stacktrace"]; !ok {
		t.Errorf("'stacktrace' key missing in logging data")
	}
}

func TestContext(t *testing.T) {
	logr, hook := test.NewNullLogger()
	ctx := context.Background()
	tctx := WithLogger(ctx, logr)
	got := FromContext(tctx)
	if got == nil {
		t.Fatal("logger not loaded into context")
	}
	if got.Logger != logr {
		t.Error("retrieved logger does not match loaded")
	}
	got.Info("test")
	e := hook.LastEntry()
	if e.Data[correlationIDKey] == "" {
		t.Error("missing correlation id", e.Data)
	}
}
