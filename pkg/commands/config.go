// Package commands process flags/environment variables/config file
// It contains global variables with configs and commands
// nolint: gochecknoglobals, gochecknoinits, unparam
package commands

import (
	"fmt"

	"github.com/takama/caldera/pkg/helper"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configCmd represents API settings command
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

func init() {
	apiCmd.AddCommand(configCmd)

	configCmd.PersistentFlags().Int("port", 8000, "A service port number")
	configCmd.PersistentFlags().Int("gateway-port", 8480, "A service rest gateway port number")
	helper.LogF("Flag error", viper.BindPFlag("api.config.port", configCmd.PersistentFlags().Lookup("port")))
	helper.LogF(
		"Flag error",
		viper.BindPFlag("api.config.gateway.port", configCmd.PersistentFlags().Lookup("gateway-port")),
	)
}
