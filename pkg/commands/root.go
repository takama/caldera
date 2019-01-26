// Package commands process flags/environment variables/config file
// It contains global variables with configs and commands
// nolint: gochecknoglobals, gochecknoinits, unparam
package commands

import (
	"fmt"

	"github.com/takama/caldera/pkg/config"
	"github.com/takama/caldera/pkg/generator"
	"github.com/takama/caldera/pkg/helper"
	"github.com/takama/caldera/pkg/input"
	"github.com/takama/caldera/pkg/version"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "caldera",
	Short: "A service boilerplate generator",
	Long: `In this mode, you'll be asked about the general
properties associated with the new service.
The configuration file will be used for all other data,
such as the host, port, etc., if you have saved it before.
Otherwise, the default settings will be used.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := new(config.Config)
		if err := viper.Unmarshal(&cfg); err != nil {
			fmt.Println("Error parsing of configuration, used default:", err)
		}
		cfg = input.Inquire(cfg)
		if cfg.Storage.MySQL &&
			cfg.Storage.Config.Port == config.DefaultPostgresPort {
			cfg.Storage.Config.Driver = config.StorageMySQL
			cfg.Storage.Config.Host = config.StorageMySQL
			cfg.Storage.Config.Port = config.DefaultMySQLPort
			cfg.Storage.Config.Username = config.StorageMySQL
			cfg.Storage.Config.Password = config.StorageMySQL
		}
		if cfg.Storage.Postgres &&
			cfg.Storage.Config.Port == config.DefaultMySQLPort {
			cfg.Storage.Config.Driver = config.StoragePostgres
			cfg.Storage.Config.Host = config.StoragePostgres
			cfg.Storage.Config.Port = config.DefaultPostgresPort
			cfg.Storage.Config.Username = config.StoragePostgres
			cfg.Storage.Config.Password = config.StoragePostgres
		}
		generator.Run(cfg)
	},
}

// Run adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Run() {
	helper.LogF("Bootstrap error", RootCmd.Execute())
}

func init() {
	fmt.Printf("%s %s\n\n", config.ServiceName, version.RELEASE)

	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./default.yaml)")
	RootCmd.PersistentFlags().String("templates", ".templates", "templates dir")
	RootCmd.PersistentFlags().String("service", "", "A boilerplate service repository dir")
	helper.LogF(
		"Flag error",
		viper.BindPFlag("directories.templates", RootCmd.PersistentFlags().Lookup("templates")),
	)
	helper.LogF(
		"Flag error",
		viper.BindPFlag("directories.service", RootCmd.PersistentFlags().Lookup("service")),
	)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName("default")        // name of config file (without extension)
	viper.AddConfigPath("/etc/caldera/")  // path to look for the config file in
	viper.AddConfigPath("$HOME/.caldera") // call multiple times to add many search paths
	viper.AddConfigPath(".")              // optionally look for config in the working directory
	viper.AutomaticEnv()                  // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		helper.LogF("Could not write config", viper.WriteConfigAs("default.yaml"))
	}
}
