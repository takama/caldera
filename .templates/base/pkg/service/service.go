package service

import (
	"context"
	"fmt"
	"net/http"

	"{{[ .Project ]}}/pkg/config"
	{{[- if .Storage.Enabled ]}}
	"{{[ .Project ]}}/pkg/db"
	"{{[ .Project ]}}/pkg/db/migrations"
	{{[- if .Storage.Postgres ]}}
	"{{[ .Project ]}}/pkg/db/postgres"
	{{[- end ]}}
	{{[- if .Storage.MySQL ]}}
	"{{[ .Project ]}}/pkg/db/mysql"
	{{[- end ]}}
	"{{[ .Project ]}}/pkg/db/stub"
	{{[- end ]}}
	"{{[ .Project ]}}/pkg/info"
	"{{[ .Project ]}}/pkg/logger"
	{{[- if .API.Enabled ]}}
	"{{[ .Project ]}}/pkg/server"
	{{[- end ]}}
	"{{[ .Project ]}}/pkg/system"
	"{{[ .Project ]}}/pkg/version"

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
		database, err = postgres.New(&cfg.Database, log, migrations.New(&cfg.Migrations))
	{{[- end ]}}
	{{[- if .Storage.MySQL ]}}
	case mysql.Driver:
		database, err = mysql.New(&cfg.Database, log, migrations.New(&cfg.Migrations))
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

	{{[- if .API.Gateway ]}}

	// Create gateway server
	gw, err := server.NewGateway(context.Background(), &cfg.Server, log)
	if err != nil {
		return err
	}
	{{[- end ]}}
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
	{{[- if .API.Gateway ]}}
	is.RegisterLivenessProbe(gw.LivenessProbe)
	is.RegisterReadinessProbe(gw.ReadinessProbe)
	{{[- end ]}}
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
		{{[- if .API.Gateway ]}}
		gw,
		{{[- end ]}}
		{{[- end ]}}
		{{[- if .Storage.Enabled ]}}
		database,
		{{[- end ]}}
	)

	{{[- if .API.Enabled ]}}

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

	{{[- if .API.Gateway ]}}

	// Run gateway server
	go func() {
		if err := gw.Run(context.Background()); err != nil {
			// Check for known errors
			if err != http.ErrServerClosed {
				log.Fatal(err.Error())
			}
			log.Error(err.Error())
		}
	}()
	{{[- end ]}}
	{{[- end ]}}

	// Wait signals
	return system.NewSignals().Wait(log, operator)
}
