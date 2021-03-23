package job

import (
	"time"

	"postgres-exporter/internal"
	"postgres-exporter/internal/metrics"
)

func CheckActiveLocks(repository internal.LockRepository, location string, database string) func() {
	return func() {
		locks, err := repository.Find(database)
		if err != nil {
			panic(err)
		}

		for _, l := range locks {
			metrics.ActiveLockCount.WithLabelValues(location, database, "", l.Relation).Inc()

			queryDuration := time.Since(l.TransactionStarted)
			metrics.ActiveLockDuration.WithLabelValues(location, database, "", l.Relation).Add(queryDuration.Seconds())
		}
	}
}
