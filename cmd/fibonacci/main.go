package main

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

var (
	port      = flag.Int("port", 8080, "the port to listen")
	develMode = flag.Bool("devel", false, "development mode")
)

func Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		numStr := r.URL.Query().Get("n")
		num, err := strconv.Atoi(numStr)
		if err != nil || num < 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		fmt.Fprint(w, getNumber(num))
	})
}

func main() {
	flag.Parse()

	logger := initLogger()

	handler := Handler()
	handler = LoggingMiddleware(logger, handler)
	handler = MetricsMiddleware(handler)

	http.Handle("/metrics", promhttp.Handler())
	http.Handle("/fibonacci", handler)

	// sugaredLogger := logger.Sugar()
	// sugaredLogger.Infow("starting http server", "port", *port)

	logger.Info("starting http server", zap.Int("port", *port))
	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		logger.Fatal("error starting http server", zap.Error(err))
	}
}
