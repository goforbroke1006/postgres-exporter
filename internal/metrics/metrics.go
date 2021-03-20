package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	LiveTuples = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "postgres_exporter_live_tuples",
		Help: "Count of dead tuples",
	}, []string{
		"location", "database", "schema", "table",
	})
	DeadTuples = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "postgres_exporter_dead_tuples",
		Help: "Count of dead tuples",
	}, []string{
		"location", "database", "schema", "table",
	})
)
