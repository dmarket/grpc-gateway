package prometheus

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

var metricsNamespace = "grpc_gateway"

func SetMetricNamespace(ns string) {
	metricsNamespace = ns
}

var (
	successResponseCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: metricsNamespace,
			Name:      "success_response_count",
			Help:      "Number of success responses per endpoint.",
		},
		[]string{"path", "method"},
	)

	failedResponseCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: metricsNamespace,
			Name:      "failed_response_count",
			Help:      "Number of failed responses per endpoint",
		},
		[]string{"path", "method", "error"},
	)
	responseTime = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: metricsNamespace,
			Name:      "http_server_response_time",
			Help:      "Histogram of response latency (seconds) of calls processed by the server.",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"path", "method"},
	)
)

func Register() {
	prometheus.MustRegister(successResponseCount)
	prometheus.MustRegister(failedResponseCount)
	prometheus.MustRegister(responseTime)
}

func PushSuccessResponse(ctx context.Context) {
	if data, ok := FromContext(ctx); ok {
		successResponseCount.WithLabelValues(data.Path, data.Method).Inc()
		responseTime.WithLabelValues(data.Path, data.Method).
			Observe(time.Since(data.StartAt).Seconds())
	}
}

func PushFailedResponse(ctx context.Context, errCode string) {
	if data, ok := FromContext(ctx); ok {
		failedResponseCount.WithLabelValues(data.Path, data.Method, errCode).Inc()
		responseTime.WithLabelValues(data.Path, data.Method).
			Observe(time.Since(data.StartAt).Seconds())
	}
}
