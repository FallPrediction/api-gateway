package helper

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	GatewayRequestTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gateway_requests_total",
			Help: "Tracks the number of gateway requests.",
		}, []string{"upstream", "operation", "uri"},
	)
	RequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Tracks the number of HTTP requests.",
		}, []string{"upstream", "method", "code", "uri"},
	)
	RequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Tracks the latencies for HTTP requests.",
			Buckets: prometheus.ExponentialBuckets(0.1, 5, 5),
		},
		[]string{"upstream", "uri"},
	)
	CircuitBreakerOpen = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "circuit_breaker_open",
			Help: "Track the number of times the circuit breaker open.",
		}, []string{"upstream"},
	)
)

func NewRegistry() *prometheus.Registry {
	registry := prometheus.NewRegistry()
	registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		RequestDuration,
		RequestsTotal,
		GatewayRequestTotal,
		CircuitBreakerOpen,
	)
	return registry
}
