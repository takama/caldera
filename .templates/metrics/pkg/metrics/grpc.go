{{[- if .API.Enabled ]}}
// nolint: gochecknoglobals
{{[- end ]}}
package metrics

{{[- if .API.Enabled ]}}
import grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"

// Histogram buckets are used for requests' duration measuring. They
// should be closer to SLA / SLO, including max / min observable values.
var histogramBuckets = []float64{
	0.005, 0.01, 0.025, 0.05, 0.1, 0.15, 0.2, 0.25, 0.5, 1, 2.5, 5, 10,
}

func registerGRPCMetrics() {
	// for request / errors rate Prometheus gRPC interceptor (middleware) is used
	//
	// duration metrics
	grpc_prometheus.EnableHandlingTimeHistogram(
		grpc_prometheus.WithHistogramBuckets(histogramBuckets),
	)
}
{{[- else ]}}

// no gRPC API to expose metrics
{{[- end ]}}
