// Package commands process flags/environment variables/config file
// It contains global variables with configs and commands
// nolint: gochecknoglobals, gochecknoinits, unparam
package commands

import (
	"fmt"

	"github.com/takama/caldera/pkg/config"
	"github.com/takama/caldera/pkg/helper"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// storageCmd represents the storage command
var storageCmd = &cobra.Command{
	Use:   "storage",
	Short: "Setup your storage modules",
	Run: func(cmd *cobra.Command, args []string) {
		if viper.GetBool("storage.postgres") || viper.GetBool("storage.mysql") {
			viper.Set("storage.enabled", true)
			if viper.GetBool("storage.postgres") {
				viper.Set("storage.mysql", false)
			}
			if viper.GetBool("storage.mysql") &&
				viper.GetInt("storage.driver.port") == config.DefaultPostgresPort {
				viper.Set("storage.driver.host", config.StorageMySQL)
				viper.Set("storage.driver.port", config.DefaultMySQLPort)
				viper.Set("storage.driver.name", config.StorageMySQL)
				viper.Set("storage.driver.username", config.StorageMySQL)
				viper.Set("storage.driver.password", config.StorageMySQL)
			}
			if viper.GetBool("storage.postgres") &&
				viper.GetInt("storage.driver.port") == config.DefaultMySQLPort {
				viper.Set("storage.driver.host", config.StoragePostgres)
				viper.Set("storage.driver.port", config.DefaultPostgresPort)
				viper.Set("storage.driver.name", config.StoragePostgres)
				viper.Set("storage.driver.username", config.StoragePostgres)
				viper.Set("storage.driver.password", config.StoragePostgres)
			}
		} else {
			viper.Set("storage.enabled", false)
		}
		err := viper.WriteConfig()
		if err != nil {
			fmt.Println("Error of writing storage configuration:", err)
		}
		fmt.Println("Storage configuration saved")
	},
}

func init() {
	RootCmd.AddCommand(storageCmd)

	storageCmd.PersistentFlags().Bool("enabled", false, "A Storage modules using")
	storageCmd.PersistentFlags().Bool("postgres", false, "A postgres module using")
	storageCmd.PersistentFlags().Bool("mysql", false, "A mysql module using")
	helper.LogF(
		"Flag error",
		viper.BindPFlag("storage.enabled", storageCmd.PersistentFlags().Lookup("enabled")),
	)
	helper.LogF(
		"Flag error",
		viper.BindPFlag("storage.postgres", storageCmd.PersistentFlags().Lookup("postgres")),
	)
	helper.LogF(
		"Flag error",
		viper.BindPFlag("storage.mysql", storageCmd.PersistentFlags().Lookup("mysql")),
	)
}
