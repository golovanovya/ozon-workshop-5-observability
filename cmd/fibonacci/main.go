package main

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"
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

	handler := Handler()

	http.Handle("/fibonacci", handler)

	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}
