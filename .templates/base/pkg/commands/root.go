// Package commands process flags/environment variables/config file
// It contains global variables with configs and commands
// nolint: gochecknoglobals, gochecknoinits
package commands

import (
	"fmt"
	"log"
	"os"
	"strings"

	"{{[ .Project ]}}/pkg/config"
	"{{[ .Project ]}}/pkg/helper"
	"{{[ .Project ]}}/pkg/logger"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands.
var RootCmd = &cobra.Command{
	Use:   "{{[ .Name ]}}",
	Short: "Service short description",
	Long:  `Service long description`,
}

const (
	defaultDBPort          = {{[ .Storage.Config.Port ]}}
	defaultDBProperty      = "{{[ .Storage.Config.Property ]}}"
	defaultDBMaxConnectons = {{[ .Storage.Config.Connections.Max ]}}
	defaultDBIdleCount     = {{[ .Storage.Config.Connections.Idle.Count ]}}
	defaultDBIdleTaim      = {{[ .Storage.Config.Connections.Idle.Time ]}}
)

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

	var dbFlags = []bind{
		{"driver", "driver"},
		{"host", "host"},
		{"port", "port"},
		{"name", "name"},
		{"username", "username"},
		{"password", "password"},
		{"properties", "props"},
		{"connections.max", "max-conn"},
		{"connections.idle.count", "idle-count"},
		{"connections.idle.time", "idle-time"},
		{"fixtures.dir", "fixtures-dir"},
	}

	RootCmd.PersistentFlags().StringP("db-driver", "D", "{{[ .Storage.Config.Driver ]}}", "A database driver")
	RootCmd.PersistentFlags().StringP("db-host", "H", "{{[ .Storage.Config.Host ]}}", "A database host")
	RootCmd.PersistentFlags().IntP("db-port", "P", defaultDBPort, "A database post number")
	RootCmd.PersistentFlags().StringP("db-name", "d", "{{[ .Storage.Config.Name ]}}", "A database name")
	RootCmd.PersistentFlags().StringP("db-username", "u", "{{[ .Storage.Config.Username ]}}", "A database user name")
	RootCmd.PersistentFlags().StringP("db-password", "p", "{{[ .Storage.Config.Password ]}}", "A database user password")
	RootCmd.PersistentFlags().StringSlice("db-props", []string{defaultDBProperty}, "A database properties")
	RootCmd.PersistentFlags().Int("db-max-conn", defaultDBMaxConnectons, "Maximum of database connections")
	RootCmd.PersistentFlags().Int("db-idle-count", defaultDBIdleCount, "Count of idle database connections")
	RootCmd.PersistentFlags().Int("db-idle-time", defaultDBIdleTaim,
		"Maximum amount of time in seconds a connection may be idle")
	RootCmd.PersistentFlags().StringP("db-fixtures-dir", "F", "fixtures", "A database fixtures directory")

	bindFlags("database", "db", dbFlags, RootCmd)

	bindCustomEnvs("database", "db", dbFlags)
	{{[- end ]}}
}

// initConfig reads in config file.
func initConfig() {
	// enable ability to specify config file via flag or via env.
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else
	// Check for env variable with config path.
	if cfgPath := os.Getenv(
		strings.ToUpper(strings.ReplaceAll(config.ServiceName, "-", "_")) + "_CONFIG_PATH",
	); cfgPath != "" {
		viper.SetConfigFile(cfgPath)
	}


	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Cannot open config file: %s", err)
	}

	fmt.Println("Using config file:", viper.ConfigFileUsed())
}
