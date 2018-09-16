package logger

import "go.uber.org/zap/zapcore"

// Config contains params to setup logger
type Config struct {
	Format string
	Level  Level
	Time   bool
	Out    zapcore.WriteSyncer
	Err    zapcore.WriteSyncer
}
