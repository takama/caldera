package logger

import (
	"bytes"
	"testing"

	"go.uber.org/zap"
)

const (
	customLevel     Level     = 17
	customFormatter Formatter = 17
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
		level Level
		str   string
	}{
		{
			level: LevelDebug,
			str:   "debug",
		},
		{
			level: LevelInfo,
			str:   "info",
		},
		{
			level: LevelWarn,
			str:   "warn",
		},
		{
			level: LevelError,
			str:   "error",
		},
		{
			level: LevelFatal,
			str:   "fatal",
		},
		{
			level: LevelPanic,
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
		if zapLevelConverter(l.level) > zap.FatalLevel || zapLevelConverter(l.level) < zap.DebugLevel {
			t.Errorf("Got incorrect data for %s log level", l.level.String())
		}
	}
	level := zapLevelConverter(customLevel)
	if level != zap.InfoLevel {
		t.Errorf("invalid log level:\ngot:  %s\nwant: %s", level, zap.InfoLevel)
	}
}

func TestFormatter(t *testing.T) {
	data := []struct {
		formatter Formatter
		str       string
	}{
		{
			formatter: JSONFormatter,
			str:       "json",
		},
		{
			formatter: TextFormatter,
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
	log := New(&Config{
		Format: TextFormatter.String(),
		Level:  LevelDebug,
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
	log = New(&Config{
		Format: JSONFormatter.String(),
		Level:  LevelWarn,
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
