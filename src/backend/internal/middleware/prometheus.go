package middleware

import (
	"net/http"
	"runtime"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	requestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "whoknows_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	cpuUsage = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "whoknows_cpu_usage_percent",
			Help: "Current CPU usage percentage",
		},
	)
)

func PrometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip metrics endpoint to avoid recursive counting
		if r.URL.Path != "/metrics" {
			// Record request metrics
			requestsTotal.WithLabelValues(r.Method, r.URL.Path, "200").Inc()

			// Update CPU metrics
			var cpuStats runtime.MemStats
			runtime.ReadMemStats(&cpuStats)
			cpuUsage.Set(float64(cpuStats.Sys) / float64(1024*1024)) // Convert to MB
		}

		next.ServeHTTP(w, r)
	})
}
