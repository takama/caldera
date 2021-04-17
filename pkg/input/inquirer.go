package input

import (
	"os"
	"path"
	"strings"

	"github.com/takama/caldera/pkg/config"
)

// Inquire for configuration.
// nolint: funlen, gocognit
func Inquire(cfg *config.Config) *config.Config {
	cfg.Github = StringAnswer("Provide name for your Github account", cfg.Github)
	cfg.Name = StringAnswer("Provide a name for your service", cfg.Name)
	cfg.Description = StringAnswer("Provide description for your service",
		strings.Title(strings.NewReplacer("-", " ", ".", " ", "_", " ").Replace(cfg.Name)))
	cfg.Project = StringAnswer("Provide project name", path.Join("github.com", cfg.Github, cfg.Name))
	cfg.PrivateRepo = StringAnswer(
		"Provide private repositories for import if applicable",
		strings.Join([]string{
			path.Join("github.com", cfg.Github),
			cfg.PrivateRepo,
		}, ","),
	)
	cfg.Bin = StringAnswer("Provide binary file name", cfg.Name)

	if BoolAnswer("Do you need an API for the service?") {
		cfg.API.Enabled = true

		switch OptionAnswer("Do you need gRPC (1) or gRPC+REST (2)?", "1", "2") {
		case "2":
			cfg.API.Gateway = true

			fallthrough
		case "1":
			cfg.API.GRPC = true
		}
	}

	if cfg.API.Enabled {
		if BoolAnswer("Do you want to terminate API with TLS") {
			cfg.API.Config.Insecure = false
			cfg.API.Config.Certificates.Crt = StringAnswer(
				"Provide certificate file path",
				cfg.API.Config.Certificates.Crt,
			)
			cfg.API.Config.Certificates.Key = StringAnswer(
				"Provide certificate key file path",
				cfg.API.Config.Certificates.Key,
			)
		} else {
			cfg.API.Config.Insecure = true
		}
	}

	if BoolAnswer("Do you need storage driver?") {
		cfg.Storage.Enabled = true

		switch OptionAnswer("Do you want postgres (1) or mysql (2)?", "1", "2") {
		case "1":
			cfg.Storage.Postgres = true
			cfg.Storage.MySQL = false
		case "2":
			cfg.Storage.MySQL = true
			cfg.Storage.Postgres = false
		}
	} else {
		cfg.Storage.Enabled = false
	}

	if cfg.API.Enabled && cfg.Storage.Enabled &&
		BoolAnswer("Do you need Contract API example for the service?") {
		cfg.Example = true
	}

	if BoolAnswer("Do you want to deploy your service to the Google Kubernetes Engine?") {
		cfg.GKE.Enabled = true
		cfg.GKE.Project = StringAnswer("Provide ID of your project on the GCP", cfg.GKE.Project)
		cfg.GKE.Region = StringAnswer("Provide compute region of your project on the GCP", cfg.GKE.Region)
		cfg.GKE.Cluster = StringAnswer("Provide cluster name in the GKE", cfg.GKE.Cluster)
	}

	if !path.IsAbs(cfg.Directories.Templates) {
		if currentDir, err := os.Getwd(); err == nil {
			cfg.Directories.Templates = path.Join(currentDir, cfg.Directories.Templates)
		}
	}

	cfg.Directories.Templates = StringAnswer("Templates directory", cfg.Directories.Templates)

	if cfg.Directories.Service == "" {
		if goPath := os.Getenv("GOPATH"); goPath != "" {
			cfg.Directories.Service = path.Join(goPath, "src", cfg.Project)
		}
	}

	cfg.Directories.Service = StringAnswer("New service directory", cfg.Directories.Service)

	if BoolAnswer("Do you want initialize service repository with git?") {
		cfg.GitInit = true
	} else {
		cfg.GitInit = false
	}

	return cfg
}
