package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"go.uber.org/zap"
)

var (
	port = flag.Int("port", 8080, "the port to listen")
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

	// logger, err := zap.NewProduction()
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal("cannot init zap", err)
	}

	handler := Handler()

	http.Handle("/fibonacci", handler)

	logger.Info("starting http server", zap.Int("port", *port))
	err = http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		logger.Fatal("error starting http server", zap.Error(err))
	}
}
