// Package commands process flags/environment variables/config file
// It contains global variables with configs and commands
// nolint: gochecknoglobals, gochecknoinits
package commands

import (
	"{{[ .Project ]}}/pkg/config"
	"{{[ .Project ]}}/pkg/db"
	"{{[ .Project ]}}/pkg/db/migrations"
	"{{[ .Project ]}}/pkg/helper"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type migrateAction uint8

const (
	migrateUp migrateAction = iota
	migrateDown
)

// migrateCmd represents the migrate command.
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate database up/down to the specified version",
	Long: `This command scan for migrations files and apply them or roll out
this migrations up/down to the specified version or up/down to
last version/initial state if no version specified.`,
}

func init() {
	RootCmd.AddCommand(migrateCmd)

	migrateCmd.PersistentFlags().String(
		"dir", "migrations/{{[ .Storage.Config.Driver ]}}", "A database migrations directory",
	)
	migrateCmd.PersistentFlags().Bool("active", true, "A database migrations are active")
	migrateCmd.PersistentFlags().String(
		"dialect", "{{[ .Storage.Config.Driver ]}}", "A database migrations dialect",
	)
	helper.LogF("Flag error",
		viper.BindPFlag("migrations.dir", migrateCmd.PersistentFlags().Lookup("dir")))
	helper.LogF("Flag error",
		viper.BindPFlag("migrations.active", migrateCmd.PersistentFlags().Lookup("active")))
	helper.LogF("Flag error",
		viper.BindPFlag("migrations.dialect", migrateCmd.PersistentFlags().Lookup("dialect")))
	helper.LogF("Env error", viper.BindEnv("migrations.dir"))
	helper.LogF("Env error", viper.BindEnv("migrations.active"))
	helper.LogF("Env error", viper.BindEnv("migrations.dialect"))
}

func doMigration(cmd *cobra.Command, action migrateAction) {
	cfg, err := config.New()
	helper.LogF("Loading config error", err)
	dbSQL, err := db.Connect(&cfg.Database)
	helper.LogF("Database connect error", err)
	// setup migration connection
	mig := migrations.New(&cfg.Migrations)
	err = mig.Setup(dbSQL)
	helper.LogF("Failed to setup migration", err)

	v, err := cmd.PersistentFlags().GetInt64("version")
	helper.LogE("Migrate version not used", err)

	switch action {
	case migrateDown:
		err = mig.MigrateDown(v)
	case migrateUp:
		fallthrough
	default:
		err = mig.MigrateUp(v)
	}

	helper.LogF("Migration error", err)
}
