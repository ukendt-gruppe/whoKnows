package middleware

import (
	"net/http"
	"runtime"
	"time"

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

	requestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "whoknows_http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets, // Default buckets
		},
		[]string{"method", "endpoint"},
	)

	errorRates = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "whoknows_http_request_errors_total",
			Help: "Total number of HTTP request errors",
		},
		[]string{"method", "endpoint"},
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

		// Update CPU metrics
		var cpuStats runtime.MemStats
		runtime.ReadMemStats(&cpuStats)
		cpuUsage.Set(float64(cpuStats.Sys) / float64(1024*1024)) // Convert to MB
	})
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
