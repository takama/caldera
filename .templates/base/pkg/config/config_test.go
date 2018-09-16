package config

import (
	"testing"
)

func TestConfig(t *testing.T) {
	config, err := New()
	if err != nil {
		t.Error("Expected loading of config, got", err)
	}
	{{[- if .API.Enabled ]}}
	if config.Server.Port != DefaultServerPort {
		t.Errorf("Expected %d, got %d", DefaultServerPort, config.Server.Port)
	}
	{{[- end ]}}
	if config.Info.Port != DefaultInfoPort {
		t.Errorf("Expected %d, got %d", DefaultInfoPort, config.Info.Port)
	}
	if config.Logger.Level != DefaultLoggerLevel {
		t.Errorf("Expected %d, got %d", DefaultLoggerLevel, config.Logger.Level)
	}
}
