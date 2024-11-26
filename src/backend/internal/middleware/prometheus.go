package middleware

import (
	"net/http"
	"runtime"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Total HTTP requests
	requestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "whoknows_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	// HTTP request durations (response times)
	requestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "whoknows_http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets, // Default buckets
		},
		[]string{"method", "endpoint"},
	)

	// HTTP request errors
	errorRates = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "whoknows_http_request_errors_total",
			Help: "Total number of HTTP request errors",
		},
		[]string{"method", "endpoint"},
	)

	// CPU usage (in percentage)
	cpuUsage = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "whoknows_cpu_usage_percent",
			Help: "Current CPU usage percentage",
		},
	)

	// Visited endpoints (how many times a specific URL has been visited)
	visitedEndpoints = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "whoknows_visited_endpoints_total",
			Help: "Total number of visits to specific endpoints",
		},
		[]string{"endpoint"},
	)

	// System logs (a simple example of counting system events)
	systemLogs = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "whoknows_system_logs_total",
			Help: "Total number of system log entries",
		},

		[]string{"log_type", "message"},
	)

	// Memory usage (in MB)
	memoryUsage = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "whoknows_memory_usage_mb",
			Help: "Current memory usage in megabytes",
		},
	)
)

func PrometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now() // Start time for duration measurement

		// Create a ResponseWriter to capture the status code
		rec := &statusRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rec, r)

		// Record request metrics
		requestsTotal.WithLabelValues(r.Method, r.URL.Path, http.StatusText(rec.statusCode)).Inc()
		requestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(time.Since(start).Seconds())

		// Update error rates if the status code indicates an error
		if rec.statusCode >= 400 {
			errorRates.WithLabelValues(r.Method, r.URL.Path).Inc()
		}

		// Update CPU usage metrics
		var cpuStats runtime.MemStats
		runtime.ReadMemStats(&cpuStats)
		cpuUsage.Set(float64(cpuStats.Sys) / float64(1024*1024)) // Convert to MB

		// Update memory usage metrics
		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)
		memoryUsage.Set(float64(memStats.Alloc) / (1024 * 1024)) // Convert bytes to MB

		// Track visits to specific endpoints
		visitedEndpoints.WithLabelValues(r.URL.Path).Inc()

		// Log system event (example: log a message)
		logSystemEvent("INFO", "Request processed successfully")
	})
}

// logSystemEvent is a helper function to simulate logging system events to Prometheus
func logSystemEvent(logType, message string) {
	// For now, log system events as informational or errors
	systemLogs.WithLabelValues(logType, message).Inc()
}

// statusRecorder is a wrapper around http.ResponseWriter that captures the status code
type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.statusCode = code
	rec.ResponseWriter.WriteHeader(code)
}
