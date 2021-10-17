// Package commands process flags/environment variables/config file.
// It contains global variables with configs and commands
// nolint: gochecknoglobals, gochecknoinits
package commands

import (
	"fmt"

	"github.com/takama/caldera/pkg/helper"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// linterCmd represents the linter settings command.
var linterCmd = &cobra.Command{
	Use:   "linter",
	Short: "Setup linter version and settings",
	Run: func(cmd *cobra.Command, args []string) {
		err := viper.WriteConfig()
		if err != nil {
			fmt.Println("Error of writing linter configuration:", err)
		}
		fmt.Println("Linter configuration saved")
	},
}

func init() {
	RootCmd.AddCommand(linterCmd)

	linterCmd.PersistentFlags().String("linter-version", "1.42.1", "Golang CI Linter default version")
	helper.LogF(
		"Flag error",
		viper.BindPFlag("linter.version", linterCmd.PersistentFlags().Lookup("linter-version")),
	)
}
