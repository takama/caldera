package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/{{[ .Github ]}}/{{[ .Name ]}}/pkg/config"
	{{[- if .Storage.Enabled ]}}
	"github.com/{{[ .Github ]}}/{{[ .Name ]}}/pkg/db"
	"github.com/{{[ .Github ]}}/{{[ .Name ]}}/pkg/db/migrations"
	{{[- if .Storage.Postgres ]}}
	"github.com/{{[ .Github ]}}/{{[ .Name ]}}/pkg/db/postgres"
	{{[- end ]}}
	{{[- if .Storage.MySQL ]}}
	"github.com/{{[ .Github ]}}/{{[ .Name ]}}/pkg/db/mysql"
	{{[- end ]}}
	"github.com/{{[ .Github ]}}/{{[ .Name ]}}/pkg/db/stub"
	{{[- end ]}}
	"github.com/{{[ .Github ]}}/{{[ .Name ]}}/pkg/info"
	"github.com/{{[ .Github ]}}/{{[ .Name ]}}/pkg/logger"
	{{[- if .API.Enabled ]}}
	"github.com/{{[ .Github ]}}/{{[ .Name ]}}/pkg/server"
	{{[- end ]}}
	"github.com/{{[ .Github ]}}/{{[ .Name ]}}/pkg/system"
	"github.com/{{[ .Github ]}}/{{[ .Name ]}}/pkg/version"

	"go.uber.org/zap"
)

// Run the service
func Run(cfg *config.Config) error {
	// Setup zap logger
	log := logger.New(&cfg.Logger)
	defer func(*zap.Logger) {
		if err := log.Sync(); err != nil {
			log.Error(err.Error())
		}
	}(log)

	log.Info(
		config.ServiceName,
		zap.String("version", version.RELEASE+"-"+version.COMMIT+"-"+version.BRANCH),
	)

	{{[- if .Storage.Enabled ]}}

	// Connect to the database
	var database db.Store
	var err error
	switch cfg.Database.Driver {
	{{[- if .Storage.Postgres ]}}
	case postgres.Driver:
		database, err = postgres.New(&cfg.Database, migrations.New(&cfg.Migrations), log)
	{{[- end ]}}
	{{[- if .Storage.MySQL ]}}
	case mysql.Driver:
		database, err = mysql.New(&cfg.Database, migrations.New(&cfg.Migrations), log)
	{{[- end ]}}
	case stub.Driver:
		fallthrough
	default:
		database, err = stub.New(&cfg.Database, log)
	}
	if err != nil {
		return err
	}
	{{[- end ]}}

	{{[- if .API.Enabled ]}}

	// Create new core server
	srv, err := server.New(context.Background(), &cfg.Server, log)
	if err != nil {
		return err
	}
	{{[- end ]}}

	{{[- if .Contract ]}}

	// Register data store providers
	srv.RegisterEventsProvider(database.EventsProvider())
	{{[- end ]}}

	// Register info/health-check service
	is := info.NewService(log)
	{{[- if .API.Enabled ]}}
	is.RegisterLivenessProbe(srv.LivenessProbe)
	is.RegisterReadinessProbe(srv.ReadinessProbe)
	{{[- end ]}}
	{{[- if .Storage.Enabled ]}}
	is.RegisterReadinessProbe(database.Check)
	{{[- end ]}}

	// Run info/health-check service
	infoServer := is.Run(fmt.Sprintf(":%d", cfg.Info.Port))

	// Setup operator with info/health-check server, core server and data store
	operator := system.NewOperator(
		infoServer,
		{{[- if .API.Enabled ]}}
		srv,
		{{[- end ]}}
		{{[- if .Storage.Enabled ]}}
		database,
		{{[- end ]}}
	)

	// Run core server
	go func() {
		if err := srv.Run(context.Background()); err != nil {
			// Check for known errors
			if err != context.DeadlineExceeded &&
				err != context.Canceled &&
				err != http.ErrServerClosed {
				log.Fatal(err.Error())
			}
			log.Error(err.Error())
		}
	}()

	// Wait signals
	return system.NewSignals().Wait(log, operator)
}
