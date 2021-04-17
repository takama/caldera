package generator

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/takama/caldera/pkg/config"
	"github.com/takama/caldera/pkg/helper"
)

// Run generator
// nolint: funlen
func Run(cfg *config.Config) {
	if cfg.Storage.Config.Name == "" {
		cfg.Storage.Config.Name = cfg.Name
	}

	if cfg.Storage.MySQL {
		cfg.Storage.Config.Driver = config.StorageMySQL
	}

	if cfg.Storage.Postgres {
		cfg.Storage.Config.Driver = config.StoragePostgres
	}

	helper.LogF("Base templates", copyTemplates(
		path.Join(cfg.Directories.Templates, config.Base),
		cfg.Directories.Service,
	))

	if cfg.API.Enabled {
		helper.LogF("Storage base templates", copyTemplates(
			path.Join(cfg.Directories.Templates, config.API, config.Base),
			cfg.Directories.Service,
		))

		if cfg.API.Gateway {
			helper.LogF("Gateway templates for API", copyTemplates(
				path.Join(cfg.Directories.Templates, config.API, config.APIGateway),
				cfg.Directories.Service,
			))
		}
	}

	if cfg.Storage.Enabled {
		helper.LogF("Storage base templates", copyTemplates(
			path.Join(cfg.Directories.Templates, config.Storage, config.Base),
			cfg.Directories.Service,
		))

		if cfg.Storage.Postgres {
			helper.LogF("Storage templates for postgres", copyTemplates(
				path.Join(cfg.Directories.Templates, config.Storage, config.StoragePostgres),
				cfg.Directories.Service,
			))
		}

		if cfg.Storage.MySQL {
			helper.LogF("Storage templates for mysql", copyTemplates(
				path.Join(cfg.Directories.Templates, config.Storage, config.StorageMySQL),
				cfg.Directories.Service,
			))
		}
	}

	if cfg.API.Enabled && cfg.Storage.Enabled && cfg.Example {
		helper.LogF("Contract example templates", copyTemplates(
			path.Join(cfg.Directories.Templates, config.Example, config.Base),
			cfg.Directories.Service,
		))

		if cfg.Storage.Postgres {
			helper.LogF("Contract templates for postgres", copyTemplates(
				path.Join(cfg.Directories.Templates, config.Example, config.StoragePostgres),
				cfg.Directories.Service,
			))
		}

		if cfg.Storage.MySQL {
			helper.LogF("Contract templates for mysql", copyTemplates(
				path.Join(cfg.Directories.Templates, config.Example, config.StorageMySQL),
				cfg.Directories.Service,
			))
		}
	}

	if cfg.Prometheus.Enabled {
		helper.LogF("Metrics templates for Prometheus", copyTemplates(
			path.Join(cfg.Directories.Templates, config.Metrics),
			cfg.Directories.Service,
		))
	}

	helper.LogF("Render templates", render(cfg))
	helper.LogF("Could not change directory", os.Chdir(cfg.Directories.Service))

	if cfg.API.Enabled {
		helper.LogF("Generate contracts", Exec("make", "contracts"))
	}

	log.Println("Initialize vendors:")
	helper.LogF("Init modules", Exec("make", "vendor"))

	helper.LogF("Tests", Exec("make", "check-all"))

	if cfg.GitInit {
		log.Println("Initialize Git repository:")
		helper.LogF("Init git", Exec("git", "init"))
		helper.LogF("Add repo files", Exec("git", "add", "--all"))
		helper.LogF("Initial commit", Exec("git", "commit", "-m", "'Initial commit'"))
	}

	fmt.Printf("New repository was created, use command 'cd %s'\n", cfg.Directories.Service)
}

// Exec runs the commands.
func Exec(commands ...string) error {
	execCmd := exec.Command(commands[0], commands[1:]...) // nolint: gosec
	execCmd.Stderr = os.Stderr
	execCmd.Stdout = os.Stdout
	execCmd.Stdin = os.Stdin

	return execCmd.Run()
}
