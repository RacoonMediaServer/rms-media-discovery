package server

import (
	rmsMiddleware "github.com/RacoonMediaServer/rms-packages/pkg/middleware"
	"github.com/didip/tollbooth"
	"github.com/go-openapi/runtime/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
)

func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return rmsMiddleware.PanicHandler(rmsMiddleware.RequestsCountHandler(rmsMiddleware.UnauthorizedRequestsCountHandler(handler)))
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
