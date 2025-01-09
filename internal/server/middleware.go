package server

import (
	"net/http"

	rmsMiddleware "github.com/RacoonMediaServer/rms-packages/pkg/middleware"
	rms_users "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-users"
	"github.com/didip/tollbooth"
	"github.com/go-openapi/runtime/middleware"
	"github.com/prometheus/client_golang/prometheus"
)

func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return rmsMiddleware.PanicHandler(rmsMiddleware.RequestsCountHandler(rmsMiddleware.UnauthorizedRequestsCountHandler(handler)))
}

func getSearchMiddleware(users rms_users.RmsUsersService, service string) middleware.Builder {
	return func(handler http.Handler) http.Handler {
		lm := tollbooth.NewLimiter(1, nil)
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			isAdmin := false
			token := request.Header.Get("X-Token")
			req := rms_users.CheckPermissionsRequest{
				Token: token,
				Perms: []rms_users.Permissions{rms_users.Permissions_AccountManagement},
			}
			resp, err := users.CheckPermissions(request.Context(), &req)
			if err == nil {
				isAdmin = resp.Allowed
			}

			// ограничиваем частоту поисковых запросов
			if !isAdmin && lm.LimitReached(token) {
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
