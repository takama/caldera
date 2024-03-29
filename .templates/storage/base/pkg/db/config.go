package db

// Config contains params to setup database connection.
type Config struct {
	Driver      string
	DSN         string
	Host        string
	Port        int
	Name        string
	Username    string
	Password    string
	Connections Connections
	Properties  []string
	Fixtures    Fixtures
}

// Connections configures DB connections state.
type Connections struct {
	Max  int
	Idle Idle
}

type Idle struct {
	Count int
	Time  int
}

// Fixtures attributes.
type Fixtures struct {
	Dir string
}
