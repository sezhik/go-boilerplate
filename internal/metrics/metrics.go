package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	httpRequests   *prometheus.CounterVec
	httpDuration   *prometheus.HistogramVec
	httpInFlight   prometheus.Gauge
	metricsHandler http.Handler
}

func New() *Metrics {
	registry := prometheus.NewRegistry()
	registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	httpRequests := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "app",
			Subsystem: "http",
			Name:      "requests_total",
			Help:      "Total number of HTTP requests.",
		},
		[]string{"method", "route", "status"},
	)

	httpDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "app",
			Subsystem: "http",
			Name:      "request_duration_seconds",
			Help:      "HTTP request latency in seconds.",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"method", "route", "status"},
	)

	httpInFlight := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "app",
			Subsystem: "http",
			Name:      "in_flight_requests",
			Help:      "Current number of in-flight HTTP requests.",
		},
	)

	registry.MustRegister(httpRequests, httpDuration, httpInFlight)

	return &Metrics{
		httpRequests:   httpRequests,
		httpDuration:   httpDuration,
		httpInFlight:   httpInFlight,
		metricsHandler: promhttp.HandlerFor(registry, promhttp.HandlerOpts{}),
	}
}

func (m *Metrics) Handler() http.Handler {
	return m.metricsHandler
}

func (m *Metrics) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/metrics" {
			c.Next()
			return
		}

		startedAt := time.Now()
		m.httpInFlight.Inc()
		defer m.httpInFlight.Dec()

		c.Next()

		route := c.FullPath()
		if route == "" {
			route = "unmatched"
		}

		status := strconv.Itoa(c.Writer.Status())
		method := c.Request.Method
		duration := time.Since(startedAt).Seconds()

		m.httpRequests.WithLabelValues(method, route, status).Inc()
		m.httpDuration.WithLabelValues(method, route, status).Observe(duration)
	}
}
