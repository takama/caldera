// Package commands process flags/environment variables/config file.
// It contains global variables with configs and commands.
// nolint: gochecknoglobals, gochecknoinits
package commands

import (
	"{{[ .Project ]}}/pkg/config"
	"{{[ .Project ]}}/pkg/helper"
	"{{[ .Project ]}}/pkg/service"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serveCmd represents the serve command.
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Listen and handle requests including health/ready checks",
	Long: `This command prepare the service for handling
of the requests to the service.
Also there are setup a health check and a readiness check
which should observe a liveness/readiness of registered modules`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.New()
		helper.LogF("Loading config error", err)
		// Runs the service
		helper.LogF("Service start error", service.Run(cfg))
	},
}

func init() {
	var infoFlags = []bind{
		{"port", "port"},
		{"statistics", "statistics"},
	}

	{{[- if .API.Enabled ]}}

	var serverFlags = []bind{
		{"port", "port"},
		{{[- if .API.Gateway ]}}
		{"gateway.port", "gw-port"},
		{{[- end ]}}
	}
	{{[- end ]}}

	RootCmd.AddCommand(serveCmd)

	serveCmd.PersistentFlags().Int("info-port", config.DefaultInfoPort, "Health port number")
	serveCmd.PersistentFlags().Bool("info-statistics", config.DefaultInfoStatistics, "Collect statistics information")
	{{[- if .API.Enabled ]}}
	serveCmd.PersistentFlags().Int("server-port", config.DefaultServerPort, "Service listening port number")
	{{[- if .API.Gateway ]}}
	serveCmd.PersistentFlags().Int("server-gw-port", config.DefaultGatewayPort, "Gateway listening port number")
	{{[- end ]}}
	{{[- end ]}}

	bindFlags("info", "info", infoFlags, serveCmd)
	{{[- if .API.Enabled ]}}
	bindFlags("server", "server", serverFlags, serveCmd)
	{{[- end ]}}
	bindEnvs("info", infoFlags)
	{{[- if .API.Enabled ]}}
	bindEnvs("server", serverFlags)
	{{[- end ]}}
}
