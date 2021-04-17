package metrics

const (
	// DefaultPath is default path to expose metrics (in Prometheus format).
	DefaultPath = "/metrics"
)

// Register all metrics from the package.
func Register() {
	registerConstMetrics()
	{{[- if .API.Enabled ]}}
	registerGRPCMetrics()
	{{[- end ]}}
	{{[- if .Storage.Enabled ]}}
	registerDatabaseMetrics()
	{{[- end ]}}
	registerBusinessMetrics()
}

func registerBusinessMetrics() {
	// add business metrics here
}
