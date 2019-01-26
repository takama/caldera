// Package commands process flags/environment variables/config file
// It contains global variables with configs and commands
// nolint: gochecknoglobals, gochecknoinits, unparam
package commands

import (
	"github.com/spf13/cobra"
)

// upCmd represents the up command
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Migrate database up to the specified version",
	Long: `This command scan for migrations files and apply
these migrations up to the specified version or up to last
version if no version specified.`,
	Run: func(cmd *cobra.Command, args []string) {
		doMigration(cmd, migrateUp)
	},
}

func init() {
	migrateCmd.AddCommand(upCmd)

	upCmd.PersistentFlags().Int64P("version", "v", 0, "Migrate up to specified version")
}
