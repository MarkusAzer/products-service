package metric

import (
	"github.com/prometheus/client_golang/prometheus"
)

//Service implements UseCase interface
type Service struct {
	httpRequestHistogram *prometheus.HistogramVec
}

//NewPrometheusService create a new prometheus service
func NewPrometheusService() (*Service, error) {
	http := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "http",
		Name:      "request_duration_seconds",
		Help:      "The latency of the HTTP requests.",
		Buckets:   prometheus.DefBuckets,
	}, []string{"handler", "method", "code"})

	s := &Service{
		httpRequestHistogram: http,
	}

	err := prometheus.Register(s.httpRequestHistogram)
	if err != nil && err.Error() != "duplicate metrics collector registration attempted" {
		return nil, err
	}
	return s, nil
}

//SaveHTTP send metrics to server
func (s *Service) SaveHTTP(h *HTTP) {
	s.httpRequestHistogram.WithLabelValues(h.Handler, h.Method, h.StatusCode).Observe(h.Duration)
}
