package job

import (
	"postgres-exporter/internal"
	"postgres-exporter/internal/metrics"
)

func CheckDeadTuples(repository internal.StatUserTablesRepository, location string, database string) func() {
	return func() {
		top, err := repository.FindTopDeadTuples(1000)
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
			metrics.DeadTuples.WithLabelValues(labels...).Add(float64(stat.DeadTup))
		}
	}
}
