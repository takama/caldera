package logger

import (
	"os"
	"strconv"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Level defines log levels
type Level int8

// Log levels
const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
	LevelPanic
)

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarn:
		return "warn"
	case LevelError:
		return "error"
	case LevelFatal:
		return "fatal"
	case LevelPanic:
		return "panic"
	default:
		return strconv.Itoa(int(l))
	}
}

// Formatter defines output formatter
type Formatter int

func (f Formatter) String() string {
	switch f {
	case JSONFormatter:
		return "json"
	case TextFormatter:
		fallthrough
	default:
		return "txt"
	}
}

// Format of outputs
const (
	TextFormatter Formatter = iota
	JSONFormatter
)

// New creates and configure new zap logger
func New(cfg *Config) *zap.Logger {
	// Define our level-handling logic.
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel && lvl >= zapLevelConverter(cfg.Level)
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel && lvl >= zapLevelConverter(cfg.Level)
	})

	// High-priority output should go to standard error, and low-priority
	// output should go to standard out.
	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)
	if cfg.Out != nil {
		consoleDebugging = zapcore.Lock(cfg.Out)
	}
	if cfg.Err != nil {
		consoleErrors = zapcore.Lock(cfg.Err)
	}

	var timeKey string
	if cfg.Time {
		timeKey = "ts"
	}
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        timeKey,
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Optimize console output for operators.
	var encoder zapcore.Encoder
	switch cfg.Format {
	case JSONFormatter.String():
		// Optimize console output for parsing.
		encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	case TextFormatter.String():
		fallthrough
	default:
		// Optimize console output for human operators.
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// Join the outputs, encoders, and level-handling functions into
	// zapcore.Cores, then tee the all cores together.
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, consoleErrors, highPriority),
		zapcore.NewCore(encoder, consoleDebugging, lowPriority),
	)

	// Use zapcore.Core to construct a Logger.
	return zap.New(core)
}

func zapLevelConverter(level Level) zapcore.Level {
	switch level {
	case LevelDebug:
		return zap.DebugLevel
	case LevelInfo:
		return zap.InfoLevel
	case LevelWarn:
		return zap.WarnLevel
	case LevelError:
		return zap.ErrorLevel
	case LevelFatal:
		return zap.FatalLevel
	case LevelPanic:
		return zap.PanicLevel
	default:
		return zap.InfoLevel
	}
}
