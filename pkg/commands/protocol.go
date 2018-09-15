package commands

import (
	"fmt"

	"github.com/takama/caldera/pkg/helper"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// protocolCmd represents the rest command
var protocolCmd = &cobra.Command{
	Use:   "rest",
	Short: "Setup API protocol settings",
	Run: func(cmd *cobra.Command, args []string) {
		err := viper.WriteConfig()
		if err != nil {
			fmt.Println("Error of writing API protocol configuration:", err)
		}
		fmt.Println("API protocol configuration saved")
	},
}

func init() {
	apiCmd.AddCommand(protocolCmd)

	protocolCmd.PersistentFlags().Int("port", 8000, "A service port number")
	helper.LogF("Flag error", viper.BindPFlag("api.config.port", protocolCmd.PersistentFlags().Lookup("port")))
}
