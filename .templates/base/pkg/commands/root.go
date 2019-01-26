// Package commands process flags/environment variables/config file
// It contains global variables with configs and commands
// nolint: gochecknoglobals, gochecknoinits
package commands

import (
	"fmt"
	"os"
	"strings"

	"{{[ .Project ]}}/pkg/config"
	"{{[ .Project ]}}/pkg/helper"
	"{{[ .Project ]}}/pkg/logger"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "{{[ .Name ]}}",
	Short: "Service short description",
	Long:  `Service long description`,
}

// Run adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Run() {
	helper.LogF("Service bootstrap error", RootCmd.Execute())
}

func init() {
	viper.SetEnvPrefix(config.ServiceName)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	viper.SetConfigType("json")
	viper.SetConfigFile(config.DefaultConfigPath)
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is config/default.conf)")
	RootCmd.PersistentFlags().Int("log-level", int(config.DefaultLoggerLevel), "Logger level (0 - debug, 1 - info, ...)")
	RootCmd.PersistentFlags().String("log-format", logger.TextFormatter.String(), "Logger format: txt, json")
	helper.LogF("Flag error",
		viper.BindPFlag("logger.level", RootCmd.PersistentFlags().Lookup("log-level")))
	helper.LogF("Flag error",
		viper.BindPFlag("logger.format", RootCmd.PersistentFlags().Lookup("log-format")))
	helper.LogF("Env error", viper.BindEnv("logger.level"))
	helper.LogF("Env error", viper.BindEnv("logger.format"))

	{{[- if .Storage.Enabled ]}}

	RootCmd.PersistentFlags().StringP("database-driver", "d", "{{[ .Storage.Config.Driver ]}}", "A database driver")
	RootCmd.PersistentFlags().StringP("database-host", "H", "{{[ .Storage.Config.Host ]}}", "A database host")
	RootCmd.PersistentFlags().IntP("database-port", "p", {{[ .Storage.Config.Port ]}}, "A database post number")
	RootCmd.PersistentFlags().StringP("database-name", "n", "{{[ .Storage.Config.Name ]}}", "A database name")
	RootCmd.PersistentFlags().StringP("database-username", "U", "{{[ .Storage.Config.Username ]}}", "A database user name")
	RootCmd.PersistentFlags().StringP("database-password", "P", "{{[ .Storage.Config.Password ]}}", "A database user password")
	RootCmd.PersistentFlags().StringSlice("database-props", []string{"sslmode=disable"}, "A database properties")
	RootCmd.PersistentFlags().Int("max-conn", 10, "Maximum of database connections")
	RootCmd.PersistentFlags().Int("idle-conn", 1, "Count of idle database connections")
	RootCmd.PersistentFlags().StringP("fixtures-dir", "F", "fixtures", "A database fixtures directory")
	helper.LogF("Flag error",
		viper.BindPFlag("database.driver", RootCmd.PersistentFlags().Lookup("database-driver")))
	helper.LogF("Flag error",
		viper.BindPFlag("database.host", RootCmd.PersistentFlags().Lookup("database-host")))
	helper.LogF("Flag error",
		viper.BindPFlag("database.port", RootCmd.PersistentFlags().Lookup("database-port")))
	helper.LogF("Flag error",
		viper.BindPFlag("database.name", RootCmd.PersistentFlags().Lookup("database-name")))
	helper.LogF("Flag error",
		viper.BindPFlag("database.username", RootCmd.PersistentFlags().Lookup("database-username")))
	helper.LogF("Flag error",
		viper.BindPFlag("database.password", RootCmd.PersistentFlags().Lookup("database-password")))
	helper.LogF("Flag error",
		viper.BindPFlag("database.properties", RootCmd.PersistentFlags().Lookup("database-props")))
	helper.LogF("Flag error",
		viper.BindPFlag("database.connections.max", RootCmd.PersistentFlags().Lookup("max-conn")))
	helper.LogF("Flag error",
		viper.BindPFlag("database.connections.idle", RootCmd.PersistentFlags().Lookup("idle-conn")))
	helper.LogF("Flag error",
		viper.BindPFlag("database.fixtures.dir", RootCmd.PersistentFlags().Lookup("fixtures-dir")))

	helper.LogF("Env error",
		viper.BindEnv("database.driver", strings.ToUpper(config.ServiceName+".db.driver")))
	helper.LogF("Env error",
		viper.BindEnv("database.host", strings.ToUpper(config.ServiceName+".db.host")))
	helper.LogF("Env error",
		viper.BindEnv("database.port", strings.ToUpper(config.ServiceName+".db.port")))
	helper.LogF("Env error",
		viper.BindEnv("database.name", strings.ToUpper(config.ServiceName+".db.name")))
	helper.LogF("Env error",
		viper.BindEnv("database.username", strings.ToUpper(config.ServiceName+".db.username")))
	helper.LogF("Env error",
		viper.BindEnv("database.password", strings.ToUpper(config.ServiceName+".db.password")))
	helper.LogF("Env error",
		viper.BindEnv("database.properties", strings.ToUpper(config.ServiceName+".db.properties")))
	helper.LogF("Env error",
		viper.BindEnv("database.connections.max", strings.ToUpper(config.ServiceName+".db.connections.max")))
	helper.LogF("Env error",
		viper.BindEnv("database.connections.idle", strings.ToUpper(config.ServiceName+".db.connections.idle")))
	helper.LogF("Env error",
		viper.BindEnv("database.fixtures.dir", strings.ToUpper(config.ServiceName+".db.fixtures.dir")))
	{{[- end ]}}
}

// initConfig reads in config file
func initConfig() {
	// enable ability to specify config file via flag or via env
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else
	// Check for env variable with config path
	if cfgPath := os.Getenv(
		strings.ToUpper(strings.Replace(config.ServiceName, "-", "_", -1)) + "_CONFIG_PATH",
	); cfgPath != "" {
		viper.SetConfigFile(cfgPath)
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
