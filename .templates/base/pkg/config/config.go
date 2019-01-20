package config

import (
	"{{[ .Project ]}}/pkg/db"
	"{{[ .Project ]}}/pkg/db/migrations"
	"{{[ .Project ]}}/pkg/info"
	"{{[ .Project ]}}/pkg/logger"
	"{{[ .Project ]}}/pkg/server"

	"github.com/spf13/viper"
)

// Default values: host, port, etc
const (
	// ServiceName - default service name
	ServiceName = "{{[ .Name ]}}"

	APIVersion = "v1alpha"

	DefaultConfigPath = "config/default.conf"

	{{[- if .API.Enabled ]}}

	DefaultServerPort     = {{[ .API.Config.Port ]}}
	{{[- if .API.Gateway ]}}
	DefaultGatewayPort    = {{[ .API.Config.Gateway.Port ]}}
	{{[- end ]}}
	{{[- end ]}}
	DefaultInfoPort       = 8080
	DefaultInfoStatistics = true
	DefaultLoggerLevel    = logger.LevelInfo
)

// Config -- Base config structure
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

// New - returns new config record initialized with default values
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
		return nil, err
	}

	return cfg, nil
}
