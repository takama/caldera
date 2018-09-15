
# {{[ .Name ]}}

{{[ .Description ]}}

## Run the service

```sh
make run
```

## Other make commands

* `all` - run default complete set of commands (build the service)
* `vendor` - import all vendors (using dep)
* `compile` - build the service binary
* `certs` - download latests certs from an alpine image and prepare it for service container
* `build` - build container image
* `push` - push an image in docker registry
{{[- if .Storage.Enabled  ]}}
* `db` - prepare local database for the service
* `migrate-up` - running up the database migration
* `migrate-down` - running down the database migration
{{[- end ]}}
* `run` - build and run the service
* `logs` - show service logs from container
* `deploy` - deployment of the service into Kubernetes environment
* `charts` - validate helm templates (charts)
* `test` - run unit tests
* `cover` - show testing coverage for packages
* `fmt` - format Go packages with go fmt
* `lint` - use set of linters ( ~ 20) to check the service code
* `stop` - stop running container
* `start` - start existing container (if it was stopped before)
* `rm` - remove stopped container
* `version` - add next major/minor/patch version
* `clean` - remove binary and running container
* `bootstrap` - check and setup if something from utilities is not exist

## Versioned changes

All changes in the project described in [changelog](docs/CHANGELOG.md)

_Generated using ([Caldera boilerplate](https://github.com/takama/caldera))_
