// Package commands process flags/environment variables/config file
// It contains global variables with configs and commands
// nolint: gochecknoglobals, gochecknoinits
package commands

import (
	"fmt"
	"os"
	"path"

	"github.com/takama/caldera/pkg/config"
	"github.com/takama/caldera/pkg/generator"
	"github.com/takama/caldera/pkg/helper"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// newCmd represents the new command.
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Generates new service from templates using default/config settings",
	Long: `In this mode, you'll be not asked about everything.
The configuration file will be used for all other data,
such as the host, port, etc., if you have saved it before.
Otherwise, the default settings will be used.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := new(config.Config)
		if err := viper.Unmarshal(&cfg); err != nil {
			fmt.Println("Error parsing of configuration, used default:", err)
		}
		if !path.IsAbs(cfg.Directories.Templates) {
			if currentDir, err := os.Getwd(); err == nil {
				cfg.Directories.Templates = path.Join(currentDir, cfg.Directories.Templates)
			}
		}
		if cfg.Directories.Service == "" {
			if currentDir, err := os.Getwd(); err == nil {
				cfg.Directories.Service = path.Join(path.Dir(currentDir), cfg.Name)
			}
		}
		generator.Run(cfg)
	},
}

func init() {
	RootCmd.AddCommand(newCmd)

	newCmd.PersistentFlags().String("namespace", config.DefaultNamespace, "A name of your module or namespace")
	newCmd.PersistentFlags().String("name", "my-service", "A name of your new service")
	newCmd.PersistentFlags().String("description", "My service", "A description of your new service")
	RootCmd.PersistentFlags().String("github", "my-account", "A Github account name")
	RootCmd.PersistentFlags().String(
		"private-repo", "",
		"PrivateRepo contains list of private repositories for import",
	)
	RootCmd.PersistentFlags().Bool("git-init", false, "Initialize repository with git")
	RootCmd.PersistentFlags().Bool("contract-example", false, "A example of contract API using")
	helper.LogF("Flag error", viper.BindPFlag("namespace", newCmd.PersistentFlags().Lookup("namespace")))
	helper.LogF("Flag error", viper.BindPFlag("name", newCmd.PersistentFlags().Lookup("name")))
	helper.LogF("Flag error", viper.BindPFlag("description", newCmd.PersistentFlags().Lookup("description")))
	helper.LogF("Flag error", viper.BindPFlag("github", RootCmd.PersistentFlags().Lookup("github")))
	helper.LogF("Flag error", viper.BindPFlag("privaterepo", RootCmd.PersistentFlags().Lookup("private-repo")))
	helper.LogF("Flag error", viper.BindPFlag("gitinit", RootCmd.PersistentFlags().Lookup("git-init")))
	helper.LogF("Flag error", viper.BindPFlag("contract", RootCmd.PersistentFlags().Lookup("contract-example")))
}
