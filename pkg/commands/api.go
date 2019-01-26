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

// apiCmd represents the api command
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Select API modules which used in the service",
	Run: func(cmd *cobra.Command, args []string) {
		if viper.GetBool("api.rest") || viper.GetBool("api.grpc") {
			viper.Set("api.enabled", true)
			if viper.GetBool("api.rest") {
				viper.Set("api.grpc", true)
			}
		} else {
			viper.Set("api.enabled", false)
		}
		err := viper.WriteConfig()
		if err != nil {
			fmt.Println("Error of writing API configuration:", err)
		}
		fmt.Println("API configuration saved")
	},
}

func init() {
	RootCmd.AddCommand(apiCmd)

	apiCmd.PersistentFlags().Bool("enabled", false, "An API modules using")
	apiCmd.PersistentFlags().Bool("rest-gateway", false, "A REST gateway module using")
	apiCmd.PersistentFlags().Bool("grpc", false, "A gRPC module using")
	helper.LogF(
		"Flag error",
		viper.BindPFlag("api.enabled", apiCmd.PersistentFlags().Lookup("enabled")),
	)
	helper.LogF(
		"Flag error",
		viper.BindPFlag("api.gateway", apiCmd.PersistentFlags().Lookup("rest-gateway")),
	)
	helper.LogF(
		"Flag error",
		viper.BindPFlag("api.grpc", apiCmd.PersistentFlags().Lookup("grpc")),
	)
}
