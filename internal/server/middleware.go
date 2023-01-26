package server

import (
	"fmt"
	"github.com/didip/tollbooth"
	"github.com/go-openapi/runtime/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/ztrue/tracerr"
	"net/http"
)

func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return panicMiddleware(requestsCountMiddleware(unauthorizedReqCountMiddleware(handler)))
}

func requestsCountMiddleware(handler http.Handler) http.Handler {
	return promhttp.InstrumentHandlerCounter(totalRequestsCounter, handler)
}

func unauthorizedReqCountMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Header.Get("X-Token") == "" {
			unauthorizedRequstsCounter.Inc()
		}
		handler.ServeHTTP(writer, request)
	})
}

func panicMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				err, ok := rec.(error)
				if !ok {
					err = fmt.Errorf("%v", rec)
				}

				e := tracerr.Wrap(err)
				frames := e.StackTrace()[2:]

				tracerr.PrintSourceColor(tracerr.CustomError(err, frames))

				w.WriteHeader(http.StatusInternalServerError)

				panicCounter.Inc()
			}
		}()

		handler.ServeHTTP(w, r)
	})
}

func getSearchMiddleware(service string) middleware.Builder {
	return func(handler http.Handler) http.Handler {
		lm := tollbooth.NewLimiter(1, nil)
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

			// ограничиваем частоту поисковых запросов
			token := request.Header.Get("X-Token")
			if lm.LimitReached(token) {
				limitReachedRequestsCounter.Inc()
				writer.WriteHeader(lm.GetStatusCode())
				return
			}

			// считаем длительность
			timer := prometheus.NewTimer(searchDurationMetric.WithLabelValues(service))
			handler.ServeHTTP(writer, request)
			timer.ObserveDuration()
		})
	}
}
