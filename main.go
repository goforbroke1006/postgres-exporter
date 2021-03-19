package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"postgres-exporter/internal/repository"
	"postgres-exporter/pkg/shutdowner"
)

var (
	httpAddr string
	target   string
	period   time.Duration
)

var (
	deadTuples = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "postgres_exporter_dead_tuples",
		Help: "Count of dead tuples",
	}, []string{
		"location", "database", "schema", "table",
	})
)

func main() {

	flag.StringVar(&httpAddr, "addr", "0.0.0.0:54380", "Handle HTTP request address")
	flag.StringVar(&target, "target", "", "Database connection string (like postgresql://user:password@localhost:5432/template1?sslmode=disable)")
	flag.DurationVar(&period, "period", 30*time.Second, "Collect data from DB period")
	flag.Parse()

	go func() {
		r := prometheus.NewRegistry()
		r.MustRegister(deadTuples)
		handler := promhttp.HandlerFor(r, promhttp.HandlerOpts{})

		http.Handle("/metrics", handler)
		if err := http.ListenAndServe(httpAddr, nil); err != nil {
			panic(err)
		}
	}()

	conn, err := pgx.Connect(context.Background(), target)
	if err != nil {
		panic(err)
	}
	defer conn.Close(context.Background())

	go func() {
		repo := repository.NewStatUserTablesRepository(conn)

		targetLocation := fmt.Sprintf("%s:%d", conn.Config().Host, conn.Config().Port)
		targetDatabase := conn.Config().Database

		for {
			startSpan := time.Now()

			top, err := repo.FindTopDeadTuples(1000)
			if err != nil {
				panic(err)
			}

			for _, stat := range top {
				labels := []string{
					targetLocation,
					targetDatabase,
					stat.SchemaName,
					stat.RelName,
				}
				deadTuples.WithLabelValues(labels...).Add(float64(stat.DeadTup))
			}

			sleep := period - time.Since(startSpan)
			if sleep > 0 {
				time.Sleep(sleep)
			}
		}
	}()

	shutdowner.WaitTermination()
}
