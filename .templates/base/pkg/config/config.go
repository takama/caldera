package config

import (
	"fmt"

{{[- if .Storage.Enabled ]}}
	"{{[ .Project ]}}/pkg/db"
	"{{[ .Project ]}}/pkg/db/migrations"
{{[- end ]}}
	"{{[ .Project ]}}/pkg/info"
	"{{[ .Project ]}}/pkg/logger"
{{[- if .API.Enabled ]}}
	"{{[ .Project ]}}/pkg/server"
{{[- end ]}}

	"github.com/spf13/viper"
)

// Default values: host, port, etc.
const (
	// ServiceName - default service name.
	ServiceName = "{{[ .Name ]}}"

	DefaultConfigPath = "config/default.conf"

	{{[- if .API.Enabled ]}}

	APIVersion = "{{[ .API.Version ]}}"

	DefaultServerPort     = {{[ .API.Config.Port ]}}
	{{[- if .API.Gateway ]}}
	DefaultGatewayPort    = {{[ .API.Config.Gateway.Port ]}}
	{{[- end ]}}
	{{[- end ]}}
	DefaultInfoPort       = 8080
	DefaultInfoStatistics = true
	DefaultLoggerLevel    = logger.LevelInfo
)

// Config -- Base config structure.
type Config struct {
	{{[- if .API.Enabled ]}}
	Server     server.Config
	{{[- end ]}}
	Info       info.Config
	{{[- if .Storage.Enabled ]}}
	Database   db.Config
	Migrations migrations.Config
	{{[- end ]}}
	Logger     logger.Config
}

// New - returns new config record initialized with default values.
func New() (*Config, error) {
	cfg := &Config{
		{{[- if .API.Enabled ]}}
		Server: server.Config{
			Port: DefaultServerPort,
			{{[- if .API.Gateway ]}}
			Gateway: server.Gateway{
				Port: DefaultGatewayPort,
			},
			{{[- end ]}}
		},
		{{[- end ]}}
		Info: info.Config{
			Port:       DefaultInfoPort,
			Statistics: DefaultInfoStatistics,
		},
		Logger: logger.Config{
			Level: DefaultLoggerLevel,
		},
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return cfg, nil
}
{{[- if .Storage.Enabled ]}}

func (cfg Config) Secure() *Config {
	c := cfg
	c.Database.Password = "***"

	return &c
}
{{[- end ]}}
