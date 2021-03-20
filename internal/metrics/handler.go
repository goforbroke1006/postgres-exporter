package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewHandler() http.Handler {
	r := prometheus.NewRegistry()
	r.MustRegister(
		LiveTuples,
		DeadTuples,
	)
	return promhttp.HandlerFor(r, promhttp.HandlerOpts{})
}
