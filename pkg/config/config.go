package config

const (
	// ServiceName defines short service name
	ServiceName = "Caldera boilerplate generator"
	// DefaultPostgresPort defines default port for PostgreSQL
	DefaultPostgresPort = 5432
	// DefaultMySQLPort defines default port for MySQL
	DefaultMySQLPort = 3306
	// Base declared base templates
	Base = "base"
	// GKE declared GKE accounts/cluster/deployment
	GKE = "gke"
	// API declared type API
	API = "api"
	// APIGateway declared type API gateway: REST
	APIGateway = "rest"
	// APIgRPC declared type API: gRPC
	APIgRPC = "grpc"
	// Contract declared contract API example
	Contract = "contract"
	// Storage declared type Storage
	Storage = "storage"
	// StoragePostgres declared storage driver type: postgres
	StoragePostgres = "postgres"
	// StorageMySQL declared storage driver type: mysql
	StorageMySQL = "mysql"
)

// Config contains service configuration
type Config struct {
	Name        string
	Description string
	Github      string
	Project     string
	Bin         string
	GitInit     bool
	Contract    bool
	GKE         struct {
		Enabled bool
		Project string
		Zone    string
		Cluster string
	}
	Storage struct {
		Enabled  bool
		Postgres bool
		MySQL    bool
		Config   struct {
			Driver      string
			Host        string
			Port        int
			Name        string
			Username    string
			Password    string
			Connections struct {
				Max  int
				Idle int
			}
		}
	}
	API struct {
		Enabled bool
		GRPC    bool
		Gateway bool
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
}
