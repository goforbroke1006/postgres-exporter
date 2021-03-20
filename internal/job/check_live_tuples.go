package job

import (
	"postgres-exporter/internal"
	"postgres-exporter/internal/metrics"
)

func CheckLiveTuples(repository internal.StatUserTablesRepository, location string, database string) func() {
	return func() {
		top, err := repository.FindTopLiveTuples(1000)
		if err != nil {
			panic(err)
		}

		for _, stat := range top {
			labels := []string{
				location,
				database,
				stat.SchemaName,
				stat.RelName,
			}
			metrics.LiveTuples.WithLabelValues(labels...).Add(float64(stat.LiveTup))
		}
	}
}
