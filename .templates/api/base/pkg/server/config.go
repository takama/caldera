package server

// Config contains params to setup server
type Config struct {
	{{[- if .API.Config.Insecure ]}}
	Port    int
	{{[- else ]}}
	Port         int
	Insecure     bool
	Certificates Certificates
	{{[- end ]}}
	Gateway Gateway
}

// Gateway contains params to setup gateway.
type Gateway struct {
	Port int
}

// Certificates contains path to certificates and key.
type Certificates struct {
	Crt string
	Key string
}
