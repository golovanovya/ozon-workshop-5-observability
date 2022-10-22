package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	InFlightRequests = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "ozon",
		Subsystem: "http",
		Name:      "in_flight_requests_total",
	})
)

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wrapper := NewResponseWrapper(w)

		next = promhttp.InstrumentHandlerInFlight(InFlightRequests, next)
		next.ServeHTTP(wrapper, r)
	})
}
