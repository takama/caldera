package logger_test

import (
	"bytes"
	"testing"

	"{{[ .Project ]}}/pkg/logger"

	"go.uber.org/zap"
)

const (
	customLevel     logger.Level     = 17
	customFormatter logger.Formatter = 17
)

type syncBufferMock struct {
	buffer bytes.Buffer
}

func (sbm *syncBufferMock) Write(p []byte) (n int, err error) {
	return sbm.buffer.Write(p)
}

func (sbm syncBufferMock) String() string {
	return sbm.buffer.String()
}

func (sbm *syncBufferMock) Reset() {
	sbm.buffer.Reset()
}

func (sbm syncBufferMock) Sync() error {
	return nil
}

func TestLevel(t *testing.T) {
	data := []struct {
		level logger.Level
		str   string
	}{
		{
			level: logger.LevelDebug,
			str:   "debug",
		},
		{
			level: logger.LevelInfo,
			str:   "info",
		},
		{
			level: logger.LevelWarn,
			str:   "warn",
		},
		{
			level: logger.LevelError,
			str:   "error",
		},
		{
			level: logger.LevelFatal,
			str:   "fatal",
		},
		{
			level: logger.LevelPanic,
			str:   "panic",
		},
		{
			level: customLevel,
			str:   "17",
		},
	}

	for _, l := range data {
		if l.level.String() != l.str {
			t.Errorf("Expected level %s, got %s", l.str, l.level.String())
		}

		if logger.ZapLevelConverter(l.level) > zap.FatalLevel ||
			logger.ZapLevelConverter(l.level) < zap.DebugLevel {
			t.Errorf("Got incorrect data for %s log level", l.level.String())
		}
	}

	level := logger.ZapLevelConverter(customLevel)

	if level != zap.InfoLevel {
		t.Errorf("invalid log level:\ngot:  %s\nwant: %s", level, zap.InfoLevel)
	}
}

func TestFormatter(t *testing.T) {
	data := []struct {
		formatter logger.Formatter
		str       string
	}{
		{
			formatter: logger.JSONFormatter,
			str:       "json",
		},
		{
			formatter: logger.TextFormatter,
			str:       "txt",
		},
		{
			formatter: customFormatter,
			str:       "txt",
		},
	}

	for _, f := range data {
		if f.formatter.String() != f.str {
			t.Errorf("Expected formatter %s, got %s", f.str, f.formatter.String())
		}
	}
}

func TestLog(t *testing.T) {
	b := syncBufferMock{}
	log := logger.New(&logger.Config{
		Format: logger.TextFormatter.String(),
		Level:  logger.LevelDebug,
		Out:    &b,
		Err:    &b,
	})

	want := "\u001b[35mDEBUG\u001b[0m\tIs text correct?\n"

	log.Debug("Is text correct?")

	got := b.String()

	if got != want {
		t.Errorf("Incorrect logger text output:\ngot:  %s\nwant: %s", got, want)
	}

	b.Reset()

	log = logger.New(&logger.Config{
		Format: logger.JSONFormatter.String(),
		Level:  logger.LevelWarn,
		Out:    &b,
		Err:    &b,
	})
	log.Warn("Is JSON correct?", zap.String("text", "correct"))
	got = b.String()
	want = "{\"level\":\"warn\",\"msg\":\"Is JSON correct?\",\"text\":\"correct\"}\n"

	if got != want {
		t.Errorf("Incorrect logger JSON output:\ngot:  %s\nwant: %s", got, want)
	}
}
