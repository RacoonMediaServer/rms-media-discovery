package server

import (
	"fmt"
	"github.com/didip/tollbooth"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/ztrue/tracerr"
	"net/http"
)

func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return setupPanicMiddleware(requestsCountMiddleware(unauthorizedReqCountMiddleware(handler)))
}

func requestsCountMiddleware(handler http.Handler) http.Handler {
	return promhttp.InstrumentHandlerCounter(totalRequestsCount, handler)
}

func unauthorizedReqCountMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Header.Get("X-Token") == "" {
			unauthorizedReqCnt.Inc()
		}
		handler.ServeHTTP(writer, request)
	})
}

func setupPanicMiddleware(handler http.Handler) http.Handler {
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

				panicCnt.Inc()
			}
		}()

		handler.ServeHTTP(w, r)
	})
}

func setupRateLimitMiddleware(handler http.Handler) http.Handler {
	lm := tollbooth.NewLimiter(1, nil)
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		token := request.Header.Get("X-Token")
		if lm.LimitReached(token) {
			limitReachedCnt.Inc()
			writer.WriteHeader(lm.GetStatusCode())
			return
		}
		handler.ServeHTTP(writer, request)
	})
}
