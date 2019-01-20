package input

import (
	"os"
	"path"
	"strings"

	"github.com/takama/caldera/pkg/config"
)

// Inquire for configuration
func Inquire(cfg *config.Config) *config.Config {
	cfg.Github = StringAnswer("Provide name for your Github account", cfg.Github)
	cfg.Name = StringAnswer("Provide name for your service", cfg.Name)
	cfg.Description = StringAnswer("Provide description for your service",
		strings.Title(strings.NewReplacer("-", " ", ".", " ", "_", " ").Replace(cfg.Name)))
	cfg.Project = StringAnswer("Provide project name", path.Join("github.com", cfg.Github, cfg.Name))
	cfg.Bin = StringAnswer("Provide binary file name", cfg.Name)
	apis := []string{config.APIGateway, config.APIgRPC}
	var count int
	question := "Do you need API for the service? "
	for len(apis) > 0 {
		if count > 0 {
			question = "Do you need one more API for the service?"
		}
		if BoolAnswer(question) {
			cfg.API.Enabled = true
			switch OptionAnswer("What kind of API do you need?", apis...) {
			case config.APIGateway:
				apis = delete(apis, config.APIGateway)
				apis = delete(apis, config.APIgRPC)
				cfg.API.Gateway = true
				cfg.API.GRPC = true
			case config.APIgRPC:
				apis = delete(apis, config.APIgRPC)
				cfg.API.GRPC = true
			}
		} else {
			if count == 0 {
				cfg.API.Enabled = false
			}
			break
		}
		count++
	}
	storages := []string{config.StoragePostgres, config.StorageMySQL}
	question = "Do you need storage driver?"
	if BoolAnswer(question) {
		cfg.Storage.Enabled = true
		switch OptionAnswer("What kind of storage driver do you need?", storages...) {
		case config.StoragePostgres:
			cfg.Storage.Postgres = true
			cfg.Storage.MySQL = false
		case config.StorageMySQL:
			cfg.Storage.MySQL = true
			cfg.Storage.Postgres = false
		}
	} else {
		cfg.Storage.Enabled = false
	}
	if cfg.API.Enabled && cfg.Storage.Enabled &&
		BoolAnswer("Do you need Contract API example for the service?") {
		cfg.Contract = true
	}
	if BoolAnswer("Do you want to deploy your service to the Google Kubernetes Engine?") {
		cfg.GKE.Enabled = true
		cfg.GKE.Project = StringAnswer("Provide ID of your project on the GCP", cfg.GKE.Project)
		cfg.GKE.Zone = StringAnswer("Provide compute zone of your project on the GCP", cfg.GKE.Zone)
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
	if BoolAnswer("Do you want initialize service repository with git") {
		cfg.GitInit = true
	} else {
		cfg.GitInit = false
	}

	return cfg
}

func delete(src []string, value string) (dst []string) {
	for i, v := range src {
		if v == value {
			// nolint: gocritic
			dst = append(src[:i], src[i+1:]...)
		}
	}
	return
}
