package prometheus

import (
	"fmt"
	"github.com/dmarket/grpc-gateway/runtime"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"strconv"
	"time"
)

var metricsNamespace = "grpc_gateway"

func SetMetricNamespace(ns string) {
	metricsNamespace = ns
}

var (
	responseCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: metricsNamespace,
			Name:      "response_count",
			Help:      "Number of responses per endpoint.",
		},
		[]string{"path", "method", "code"},
	)
	responseTime = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: metricsNamespace,
			Name:      "http_server_response_time",
			Help:      "Histogram of response latency (seconds) of calls processed by the server.",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"path", "method", "code"},
	)
)

func Register() {
	prometheus.MustRegister(responseCount)
	prometheus.MustRegister(responseTime)
}

func wrapHandlerWithLogging(wrappedHandler http.Handler) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {

	}
}

func Handler(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	fmt.Printf("runtime.Handler_prometheus ctx: %+v", r.Context())
	now := time.Now()
	defer func() {
		skip := "0"
		responseCount.WithLabelValues(r.URL.Path, r.Method, skip).Inc()
		responseTime.WithLabelValues(r.URL.Path, r.Method, skip).
			Observe(time.Since(now).Seconds())
		//if d, ok := FromContext(r.Context()); ok {
		//	responseCount.WithLabelValues(d.Path, r.Method, strconv.Itoa(rec.status)).Inc()
		//	responseTime.WithLabelValues(d.Path, r.Method, strconv.Itoa(rec.status)).
		//		Observe(time.Since(now).Seconds())
		//	fmt.Println("hello from handler: end")
		//}
	}()
}

func PHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("\nhttp.Handler_prometheus ctx: %+v\n", r.Context())
		now := time.Now()
		ww := NewStatusRecorder(w)

		defer func() {
			//if d, ok := FromContext(r.Context()); ok {
			//	responseCount.WithLabelValues(d.Path, r.Method, strconv.Itoa(ww.status)).Inc()
			//	responseTime.WithLabelValues(d.Path, r.Method, strconv.Itoa(ww.status)).
			//		Observe(time.Since(now).Seconds())
			//}

			fmt.Printf("\ndefer http.Handler_prometheus header: %+v\n", r.Header)
			fmt.Printf("\ndefer http.Handler_prometheus ctx: %+v\n", r.Context())
			responseCount.WithLabelValues(r.URL.Path, r.Method, strconv.Itoa(ww.status)).Inc()
			responseTime.WithLabelValues(r.URL.Path, r.Method, strconv.Itoa(ww.status)).
				Observe(time.Since(now).Seconds())
		}()

		h.ServeHTTP(ww, r)
	})
}
