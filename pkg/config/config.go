package config

const (
	// ServiceName defines short service name.
	ServiceName = "Caldera boilerplate generator"
	// DefaultNamespace defines default namespace in Kubernetes environment.
	DefaultNamespace = "default"
	// DefaultPostgresPort defines default port for PostgreSQL.
	DefaultPostgresPort = 5432
	// DefaultMySQLPort defines default port for MySQL.
	DefaultMySQLPort = 3306
	// Base declared base templates.
	Base = "base"
	// GKE declared GKE accounts/cluster/deployment.
	GKE = "gke"
	// API declared type API.
	API = "api"
	// APIGateway declared type API gateway: REST.
	APIGateway = "rest"
	// APIgRPC declared type API: gRPC.
	APIgRPC = "grpc"
	// OpenAPI declares openapi templates.
	OpenAPI = "openapi"
	// Example declared contract API example.
	Example = "example"
	// Storage declared type Storage.
	Storage = "storage"
	// StoragePostgres declared storage driver type: postgres.
	StoragePostgres = "postgres"
	// StoragePostgresVersion declared storage version.
	StoragePostgresVersion = "12.6"
	// StoragePostgresDefaultProperty defined default postgres propery.
	StoragePostgresDefaultProperty = "sslmode=disable"
	// StorageMySQL declares storage driver type: mysql.
	StorageMySQL = "mysql"
	// StorageMySQLVersion declared storage version.
	StorageMySQLVersion = "8.0"
	// StorageMySQLDefaultProperty defined default Mysql propery.
	StorageMySQLDefaultProperty = "tls=false"
	// Metrics declared Prometheus common metrics for the service.
	Metrics = "metrics"
)

// Config contains service configuration.
type Config struct {
	Namespace   string
	Name        string
	Description string
	Github      string
	PrivateRepo string
	Project     string
	Bin         string
	GitInit     bool
	Example     bool
	Prometheus  struct {
		Enabled bool
	}
	GKE struct {
		Enabled bool
		Project string
		Region  string
		Cluster string
	}
	Storage struct {
		Enabled  bool
		Postgres bool
		MySQL    bool
		Config   struct {
			Driver      string
			Version     string
			Host        string
			Port        int
			Name        string
			Username    string
			Password    string
			Property    string
			Connections struct {
				Max  int
				Idle struct {
					Count int
					Time  int
				}
			}
		}
	}
	API struct {
		Enabled bool
		GRPC    bool
		Gateway bool
		CORS    struct {
			Enabled bool
		}
		UI      bool
		Version string
		Config  struct {
			Port    int
			Gateway struct {
				Port int
			}
		}
	}
	Directories struct {
		Templates string
		Service   string
	}
	Linter struct {
		Version string
	}
}
