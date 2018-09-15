// Package commands process flags/environment variables/config file
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
	apiCmd.PersistentFlags().Bool("rest", false, "A REST module using")
	apiCmd.PersistentFlags().Bool("grpc", false, "A gRPC module using")
	helper.LogF(
		"Flag error",
		viper.BindPFlag("api.enabled", apiCmd.PersistentFlags().Lookup("enabled")),
	)
	helper.LogF(
		"Flag error",
		viper.BindPFlag("api.rest", apiCmd.PersistentFlags().Lookup("rest")),
	)
	helper.LogF(
		"Flag error",
		viper.BindPFlag("api.grpc", apiCmd.PersistentFlags().Lookup("grpc")),
	)
}
