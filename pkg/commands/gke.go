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

// gkeCmd represents the GKE command
var gkeCmd = &cobra.Command{
	Use:   "GKE",
	Short: "Setup Google Kubernetes Engine properties to deploy the service",
	Run: func(cmd *cobra.Command, args []string) {
		err := viper.WriteConfig()
		if err != nil {
			fmt.Println("Error of writing GKE configuration:", err)
		}
		fmt.Println("GKE configuration saved")
	},
}

func init() {
	RootCmd.AddCommand(gkeCmd)

	gkeCmd.PersistentFlags().Bool("enabled", false, "A Google Kubernetes Engine enabled")
	gkeCmd.PersistentFlags().String("project", "my-project-id", "A project ID in GCP")
	gkeCmd.PersistentFlags().String("zone", "europe-west1-b", "A compute zone in GCP")
	gkeCmd.PersistentFlags().String("cluster", "my-cluster-name", "A cluster name in GKE")
	helper.LogF(
		"Flag error",
		viper.BindPFlag("gke.enabled", gkeCmd.PersistentFlags().Lookup("enabled")),
	)
	helper.LogF(
		"Flag error",
		viper.BindPFlag("gke.project", gkeCmd.PersistentFlags().Lookup("project")),
	)
	helper.LogF(
		"Flag error",
		viper.BindPFlag("gke.zone", gkeCmd.PersistentFlags().Lookup("zone")),
	)
	helper.LogF(
		"Flag error",
		viper.BindPFlag("gke.cluster", gkeCmd.PersistentFlags().Lookup("cluster")),
	)
}
