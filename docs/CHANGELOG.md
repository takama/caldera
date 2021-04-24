# Version v0.2.4

## v0.2.4

- Added namespace for the service
- Used namespaced names in the service
## v0.2.3

- Used compute region instead of zone
- Used default driver to specify database certificates path
- Fixed ingress TLS secret name value
- Simplified initial environment data for Makefile

## v0.2.2

- Changed default name of a database credentials secret in the helm charts values
- Updated question about drivers

## v0.2.1

- Used default namespace and changed `NAMESPACE` to `CLUSTER`
- Fixed possible hidden issues with config file
- Fixed problem with port presented as the secret name
- Used database `ConfigMap` for host/port
- Added `nginx` annotations example
- Added default service account
- Extended TLS secret name block
- Added TLS hosts
- Changed ingress version

## v0.2.0

- Added sub-charts for Helm deployment.
- Added special vatiables environment to keep it outside of GIT.
- Added CORS handling.
- Managed system signals to be ignored.
- Added the possibility to use multiple database drivers.
- Added Prometheus monitoring package.
- Used gRPC gateway v2 and OpenApi v2 instead of swagger
- Redesigned dev/prod envs in Makefile.
- Updated Go version and packages.
- Fixed issue [#27](https://github.com/takama/caldera/issues/27) ([@takama](https://github.com/takama))
- Updated protoc [#37](https://github.com/takama/caldera/issues/37) ([@takama](https://github.com/takama))
- Added OpenAPI Swagger UI [#38](https://github.com/takama/caldera/issues/38) ([@takama](https://github.com/takama))
## v0.1.10

- Fixed label for service name ([@takama](https://github.com/takama))

## v0.1.9

- Updated helm templates & project region ([@takama](https://github.com/takama))

## v0.1.7

- Changed protobuf/gateway versions ([@takama](https://github.com/takama))
- Fixed some template generation bugs ([@takama](https://github.com/takama))
- Added private repositories for import path ([@takama](https://github.com/takama))
- Changed postgres driver version ([@takama](https://github.com/takama))
- Used standard health gRPC method ([@takama](https://github.com/takama))

## v0.1.6

- Changed protobuf/gateway versions ([@takama](https://github.com/takama))
- Fixed linter issues in CLI utility ([@takama](https://github.com/takama))
- Fixed linter issues in templates ([@takama](https://github.com/takama))

## v0.1.5

- Updated helm templates (Issue [32](https://github.com/takama/caldera/issues/32), [@takama](https://github.com/takama))

## v0.1.4

- Changed protobuf/gateway versions ([@takama](https://github.com/takama))
- Fixed linter issues in CLI utility ([@takama](https://github.com/takama))
- Fixed linter issues in templates ([@takama](https://github.com/takama))

## v0.1.3

- Changed protobuf/gateway versions ([@takama](https://github.com/takama))
- Added TLS support in the database ([@takama](https://github.com/takama))
- Added anti affinity attribute ([@takama](https://github.com/takama))
- Added support of Contour ingress routes ([@takama](https://github.com/takama))
- Added support of headless service ([@takama](https://github.com/takama))
- Other small fixes and changes ([@takama](https://github.com/takama))

## v0.1.1

### Codebase changes in v0.1.1

- Added info and health gRPC handlers ([#18](https://github.com/takama/k8sapp/pull/18), [@takama](https://github.com/takama))

## v0.1.0

### Codebase changes in v0.1.0

- Fixed linter issues for caldera utility and for templates - bug ([#16](https://github.com/takama/k8sapp/pull/16), [@takama](https://github.com/takama))
- Used Project in import path - enhancement ([#15](https://github.com/takama/k8sapp/pull/15), [@takama](https://github.com/takama))

## v0.0.1

### Codebase changes in v0.0.1

- Initial service ([@takama](https://github.com/takama))
