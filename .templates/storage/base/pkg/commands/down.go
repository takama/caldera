// Package commands process flags/environment variables/config file
// It contains global variables with configs and commands
// nolint: gochecknoglobals, gochecknoinits, unparam
package commands

import (
	"github.com/spf13/cobra"
)

// downCmd represents the down command
var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Migrate database down to the specified version",
	Long: `This command scan for migrations files and roll out
these migrations down to the specified version or down to initial
state if no version specified.`,
	Run: func(cmd *cobra.Command, args []string) {
		doMigration(cmd, migrateDown)
	},
}

func init() {
	migrateCmd.AddCommand(downCmd)

	downCmd.PersistentFlags().IntP("version", "v", 0, "Migrate down to specified version")
}
