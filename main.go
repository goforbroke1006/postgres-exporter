package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"

	"postgres-exporter/internal/job"
	"postgres-exporter/internal/metrics"
	"postgres-exporter/internal/repository"
	"postgres-exporter/pkg/periodic"
	"postgres-exporter/pkg/shutdowner"
)

var (
	httpAddr string
	target   string
	period   time.Duration
)

func main() {

	flag.StringVar(&httpAddr, "addr", "0.0.0.0:54380", "Handle HTTP request address")
	flag.StringVar(&target, "target", "", "Database connection string (like postgresql://user:password@localhost:5432/template1?sslmode=disable)")
	flag.DurationVar(&period, "period", 30*time.Second, "Collect data from DB period")
	flag.Parse()

	logrus.SetReportCaller(true)

	go func() {
		http.Handle("/metrics", metrics.NewHandler())
		if err := http.ListenAndServe(httpAddr, nil); err != nil {
			panic(err)
		}
	}()

	conn, err := pgxpool.Connect(context.Background(), target)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	var (
		targetLocation = fmt.Sprintf("%s:%d", conn.Config().ConnConfig.Host, conn.Config().ConnConfig.Port)
		targetDatabase = conn.Config().ConnConfig.Database
	)
	var (
		repo = repository.NewStatUserTablesRepository(conn)
	)

	jobRunner := periodic.NewJobRunner(
		periodic.Task{
			Name:      "live-tuples",
			Period:    30 * time.Second,
			RunOnInit: true,
			Function:  job.CheckLiveTuples(repo, targetLocation, targetDatabase),
		},
		periodic.Task{
			Name:      "dead-tuples",
			Period:    30 * time.Second,
			RunOnInit: true,
			Function:  job.CheckDeadTuples(repo, targetLocation, targetDatabase),
		},
	)
	go jobRunner.Run()
	defer jobRunner.Shutdown()

	shutdowner.WaitTermination()
}
