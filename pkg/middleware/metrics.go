package middleware

import (
	"net/http"

	"github.com/MarkusAzer/products-service/pkg/metric"
)

type metricsHandler struct {
	handler  http.Handler
	mService metric.UseCase
}

func (h metricsHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	appMetric := metric.NewHTTP(req.URL.Path, req.Method)
	appMetric.Started()
	h.handler.ServeHTTP(w, req)
	appMetric.Finished()
	appMetric.StatusCode = w.Header().Get("status")
	h.mService.SaveHTTP(appMetric)
}

//Metrics to prometheus
func Metrics(mService metric.UseCase) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return metricsHandler{next, mService}
	}
}
