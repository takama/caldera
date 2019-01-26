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

var (
	databasePort   int
	databaseDriver string
)

// driverCmd represents the driver command
var driverCmd = &cobra.Command{
	Use:   "driver",
	Short: "Setup database driver settings",
	Run: func(cmd *cobra.Command, args []string) {
		err := viper.WriteConfig()
		if err != nil {
			fmt.Println("Error of writing storage driver configuration:", err)
		}
		fmt.Println("Storage driver configuration saved")
	},
}

func init() {
	storageCmd.AddCommand(driverCmd)

	if viper.GetBool("storage.mysql") {
		databasePort = config.DefaultMySQLPort
		databaseDriver = config.StorageMySQL
	} else {
		databasePort = config.DefaultPostgresPort
		databaseDriver = config.StoragePostgres
	}

	driverCmd.PersistentFlags().String("host", databaseDriver, "A host name")
	driverCmd.PersistentFlags().Int("port", databasePort, "A port number")
	driverCmd.PersistentFlags().String("name", "", "A database name")
	driverCmd.PersistentFlags().StringP("username", "u", databaseDriver, "A name of database user")
	driverCmd.PersistentFlags().StringP("password", "p", databaseDriver, "An user password")
	driverCmd.PersistentFlags().Int("max-conn", 10, "Maximum available connections")
	driverCmd.PersistentFlags().Int("idle-conn", 1, "Count of idle connections")
	helper.LogF(
		"Flag error",
		viper.BindPFlag("storage.config.host", driverCmd.PersistentFlags().Lookup("host")),
	)
	helper.LogF(
		"Flag error",
		viper.BindPFlag("storage.config.port", driverCmd.PersistentFlags().Lookup("port")),
	)
	helper.LogF(
		"Flag error",
		viper.BindPFlag("storage.config.name", driverCmd.PersistentFlags().Lookup("name")),
	)
	helper.LogF(
		"Flag error",
		viper.BindPFlag("storage.config.username", driverCmd.PersistentFlags().Lookup("username")),
	)
	helper.LogF(
		"Flag error",
		viper.BindPFlag("storage.config.password", driverCmd.PersistentFlags().Lookup("password")),
	)
	helper.LogF(
		"Flag error",
		viper.BindPFlag("storage.config.connections.max", driverCmd.PersistentFlags().Lookup("max-conn")),
	)
	helper.LogF(
		"Flag error",
		viper.BindPFlag("storage.config.connections.idle", driverCmd.PersistentFlags().Lookup("idle-conn")),
	)
}
