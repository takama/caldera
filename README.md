# Caldera

A command line utility Caldera allows you to create a boilerplate service that ready to run inside the container. This will save two or more days of developers working, who decided to create their first (micro) service.

## Features of the boilerplate service

- gRPC/REST API example using protobuf
- Implementation of the health checks
- Configuring the service using config file, environment variables or flags
- Processing of graceful shutdown for every registered component
- Database interfaces with migration features
- CI/CD pipelines integrated into Makefile
- Helm charts for deploying the service in Kubernetes environment
- Container SSL certificates integration for using a secure client
- Integration of the package manager
- Versioning automation

### Features of the command line utility Caldera

- Using of configuration file that contains your saved preferences for a new boilerplate service
- Interactive mode to select preferred features for a new service
- Using of CLI flags to create new service quickly

## Requirements

- [Go compiler](https://golang.org/dl/) v1.9 or newer
- [GNU make utility](https://en.wikipedia.org/wiki/Make_(software)) that probably already installed on your system

### Requirements for boilerplate service

- Docker service, version 18.03 or newer

## Setup

```sh
go get -u github.com/takama/caldera
cd $GOPATH/src/github.com/takama/caldera
make
```

## Usage of Caldera

### Interactive mode

In this mode, you'll be asked about the general properties associated with the new service. The configuration file will be used for all other data, such as the host, port, etc., if you have saved it before. Otherwise, the default settings will be used.

```txt
./caldera
Caldera boilerplate version: v0.0.1 build date: 2018-09-15T12:02:17+07

Provide name for your Github account (my-account):
Provide name for your service (my-service): new-service
Provide description for your service (New service): Very new service
Do you need API for the service? (y/n): y
Do you want gRPC (1) or gRPC+REST (2)?: 2
Do you want to terminate API with TLS? (y/n): n
Do you need storage driver? (y/n): y
Do you want postgres (1) or mysql (2)?: 1
Do you need Contract API example for the service? (y/n): y
Do you want to deploy your service to the Google Kubernetes Engine? (y/n): y
Provide ID of your project on the GCP (my-project-id):
Provide compute zone of your project on the GCP (europe-west4):
Provide cluster name in the GKE (my-cluster-name):
Templates directory (~/go/src/github.com/takama/caldera/.templates):
New service directory (~/go/src/github.com/my-account/my-service):
Do you want initialize service repository with git (y/n): y
```

### CLI mode

In this mode, you'll be not asked about everything. The configuration file will be used for all other data, such as the host, port, etc., if you have saved it before. Otherwise, the default settings will be used.

```sh
./caldera new [ --service <name> --description <description> --github <account> --grpc-client ]
```

### Save configuration for future use

For example of save a `storage` parameters in Caldera configuration file:

```sh
./caldera storage [flags]

Flags:
  -h, --help       help for storage
      --enabled    A Storage modules using
      --postgres   A postgres module using
      --mysql      A mysql module using
```

Save a `storage` parameters for database driver in Caldera configuration file:

```sh
./caldera storage driver [flags]

Flags:
  -h, --help              help for driver
      --host string       A host name (default "postgres")
      --port int          A port number (default 5432)
      --name string       A database name (default "postgres")
  -u, --username string   A name of database user (default "postgres")
  -p, --password string   An user password (default "postgres")
      --max-conn int      Maximum available connections (default 10)
      --idle-conn int     Count of idle connections (default 1)
```

Save an `API` parameters for `REST/gRPC` (REST always used gRCP gateway):

```sh
./caldera api [flags]

Flags:
  -h, --help           help for api
      --enabled        An API modules using
      --grpc           A gRPC module using
      --rest-gateway   A REST gateway module using
```

Save a common `API` parameters:

```sh
./caldera api config [flags]

Flags:
  -h, --help               help for config
      --port int           A service port number (default 8000)
      --gateway-port int   A service rest gateway port number (default 8001)
```

## gRPC/REST API example

This example contains a good approach to using the API with the code-generated Client/Server from the interfaces in the `.proto` definitions using the Go language. In addition, it contains a gRPC gateway that can be used to access the API via REST.

```proto
// Interface exported by the server
service Events {
  // Get the Event object by ID
  rpc GetEvent (request.ByID) returns (Event) {
    option (google.api.http).get = "/v1/events/id/{id}";
  }

  // Find the Event objects by name
  rpc FindEventsByName (request.ByName) returns (stream Event) {
    option (google.api.http).get = "/v1/events/name/{name}";
  }

  // List all Events
  rpc ListEvents (google.protobuf.Empty) returns (stream Event) {
    option (google.api.http).get = "/v1/events";
  }

  // Create a new Event object
  rpc CreateEvent (Event) returns (Event) {
    option (google.api.http) = {
      post: "/v1/events",
      body: "*"
    };
  }

  // Update the Event object
  rpc UpdateEvent (Event) returns (Event) {
    option (google.api.http) = {
      put: "/v1/events/id/{id}",
      body: "*"
    };
  }

  // Delete the Event object by ID
  rpc DeleteEvent (request.ByID) returns (google.protobuf.Empty) {
    option (google.api.http).delete = "/v1/events/id/{id}";
  }

  // Delete The Event objects by Event name
  rpc DeleteEventsByName (request.ByName) returns (google.protobuf.Empty) {
    option (google.api.http).delete = "/v1/events/name/{name}";
  }
}
```

## Health checks

Service should have [health checks](https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/) for successful execution in containers environment. It should helps with correct orchestration of the service.

## Configuring

The [twelve factors](https://12factor.net/config) service must be configured using environment variables. The service has a built-in library for automatically recognizing and allocating environment variables that are stored inside `struct` of different types. As additional methods, a configuration file and flags are used. All these methods of setting are directly linked to each other in using of configuration variables.

## System signals

The service has an ability to intercept system signals and transfer actions to special methods for graceful shutdown, maintenance mode, reload of configuration, etc.

```go
type Signals struct {
    shutdown    []os.Signal
    reload      []os.Signal
    maintenance []os.Signal
}
```

## Database integration

Boilerplate service is provided database drivers and optional `stub` driver for testing purposes.
The database drivers are using in local container environment as well. Corresponded `make` commands `db, migrate-up, migrate-down` allow to run database engine and migrate data into database.
Supported the following database drivers:

- Postgres
- MySQL

## Build automation

In the CI/CD pipeline, there is a series of commands for the static cross-compilation of the service for the specified OS. Build a docker image and push it into the container registry. Optimal and compact `docker` image `FROM SCRATCH`.

```Dockerfile
FROM scratch

ENV MY_SERVICE_SERVER_PORT 8000
ENV MY_SERVICE_INFO_PORT 8080
ENV MY_SERVICE_LOGGER_LEVEL 0

EXPOSE $MY_SERVICE_SERVER_PORT
EXPOSE $MY_SERVICE_INFO_PORT

COPY certs /etc/ssl/certs/
COPY migrations /migrations/
COPY bin/linux-amd64/service /

CMD ["/service", "serve"]
```

## SSL support

Certificates support for creating a secure SSL connection in the `Go` client. Attaching the certificate to the docker image.

## Testing

The command `make test` is running set of checks and tests:

- tool `go fmt` used on package sources
- set of linters used on package sources (20+ types of linters)
- tests used on package sources excluding vendor
- a testing coverage of new boilerplate service
- compile and check of Helm charts

## Helm charts and Continuous Delivery

A set of basic templates for the deployment of the service in Kubernetes has been prepared. Only one `make deploy` command loads the service into Kubernetes. Wait for the successful result, and the service will be ready to go.

## Package manager

To properly work with dependencies, we need to select a package manager. `go mod` is dependency management tools for Go.

## Versioning automation

Using a special script to increment the release version

```sh
make version
Current version v0.0.1.
Please enter new version [v0.0.2]:
```

## Contributing to the project

See the [contribution guidelines](docs/CONTRIBUTING.md) for information on how to participate in the Caldera project to submitting a pull request or creating a new issue.

## Versioned changes

All changes in the project described in [changelog](docs/CHANGELOG.md)

## License

[MIT Public License](https://github.com/takama/caldera/blob/master/LICENSE)
