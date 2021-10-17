package config_test

import (
	"testing"

	"{{[ .Project ]}}/pkg/config"
)

func TestConfig(t *testing.T) {
	t.Parallel()

	cfg, err := config.New()

	if err != nil {
		t.Error("Expected loading of config, got", err)
	}
	{{[- if .API.Enabled ]}}

	if cfg.Server.Port != config.DefaultServerPort {
		t.Errorf("Expected %d, got %d", config.DefaultServerPort, cfg.Server.Port)
	}
	{{[- end ]}}

	if cfg.Info.Port != config.DefaultInfoPort {
		t.Errorf("Expected %d, got %d", config.DefaultInfoPort, cfg.Info.Port)
	}

	if cfg.Logger.Level != config.DefaultLoggerLevel {
		t.Errorf("Expected %d, got %d", config.DefaultLoggerLevel, cfg.Logger.Level)
	}
}
