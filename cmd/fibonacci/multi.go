package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

func MultiHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		nums := r.URL.Query()["ns"]

		resText := strings.Builder{}
		if len(nums) == 0 {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "no numbers in request")
		}
		for _, num := range nums {
			res, err := requestAnotherService(ctx, num)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintln(w, err)
				return
			}

			fmt.Fprintf(&resText, "fib(%s) = %s\n", num, res)
		}
		w.Write([]byte(resText.String()))
	})
}

var errStatusCode = errors.New("wrong status code")

func requestAnotherService(ctx context.Context, n string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "requestAnotherService")
	span.SetTag("n", n)
	ext.SpanKindRPCClient.Set(span)
	defer span.Finish()

	query := url.Values{}
	query.Add("n", n)
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"http://127.0.0.1:8080/fibonacci?"+query.Encode(),
		nil,
	)
	if err != nil {
		return "", err
	}

	opentracing.GlobalTracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header),
	)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", errStatusCode
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(data)), nil
}
