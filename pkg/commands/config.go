// Package commands process flags/environment variables/config file
// It contains global variables with configs and commands
// nolint: gochecknoglobals, gochecknoinits
package commands

import (
	"fmt"

	"github.com/takama/caldera/pkg/helper"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configCmd represents API settings command.
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Setup API settings",
	Run: func(cmd *cobra.Command, args []string) {
		err := viper.WriteConfig()
		if err != nil {
			fmt.Println("Error of writing API settings:", err)
		}
		fmt.Println("API configuration saved")
	},
}

const (
	defaultServerPort  = 8000
	defaultGatewayPort = 8480
)

func init() {
	apiCmd.AddCommand(configCmd)

	configCmd.PersistentFlags().Int("port", defaultServerPort, "A service port number")
	configCmd.PersistentFlags().Int("gateway-port", defaultGatewayPort, "A service rest gateway port number")
	helper.LogF("Flag error", viper.BindPFlag("api.config.port", configCmd.PersistentFlags().Lookup("port")))
	helper.LogF(
		"Flag error",
		viper.BindPFlag("api.config.gateway.port", configCmd.PersistentFlags().Lookup("gateway-port")),
	)
}
