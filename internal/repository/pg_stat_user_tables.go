package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"

	"postgres-exporter/internal"
)

func NewStatUserTablesRepository(conn *pgxpool.Pool) *statUserTablesRepository {
	return &statUserTablesRepository{
		conn: conn,
	}
}

var _ internal.StatUserTablesRepository = &statUserTablesRepository{}

type statUserTablesRepository struct {
	conn *pgxpool.Pool
}

func (r statUserTablesRepository) FindTopDeadTuples(limit uint) (result []internal.StatUserTablesRow, err error) {
	var limitChunk string
	if limit > 0 {
		limitChunk = fmt.Sprintf("LIMIT %d", limit)
	}

	query := fmt.Sprintf(`
		SELECT 
			schemaname, 
			relname,
			n_live_tup,
			n_dead_tup 
		FROM pg_stat_user_tables
		ORDER BY n_dead_tup DESC
		%s
		;
	`, limitChunk)

	rows, err := r.conn.Query(context.TODO(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		schemaName string
		relName    string
		nLiveTup   int64
		nDeadTup   int64
	)

	for rows.Next() {
		err = rows.Scan(&schemaName, &relName, &nLiveTup, &nDeadTup)
		if err != nil {
			return nil, err
		}

		result = append(result, internal.StatUserTablesRow{
			SchemaName: schemaName,
			RelName:    relName,
			LiveTup:    nLiveTup,
			DeadTup:    nDeadTup,
		})
	}

	return result, nil
}

func (r statUserTablesRepository) FindTopLiveTuples(limit uint) (result []internal.StatUserTablesRow, err error) {
	var limitChunk string
	if limit > 0 {
		limitChunk = fmt.Sprintf("LIMIT %d", limit)
	}

	query := fmt.Sprintf(`
		SELECT 
			schemaname, 
			relname,
			n_live_tup,
			n_dead_tup 
		FROM pg_stat_user_tables
		ORDER BY n_live_tup DESC
		%s
		;
	`, limitChunk)

	rows, err := r.conn.Query(context.TODO(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		schemaName string
		relName    string
		nLiveTup   int64
		nDeadTup   int64
	)

	for rows.Next() {
		err = rows.Scan(&schemaName, &relName, &nLiveTup, &nDeadTup)
		if err != nil {
			return nil, err
		}

		result = append(result, internal.StatUserTablesRow{
			SchemaName: schemaName,
			RelName:    relName,
			LiveTup:    nLiveTup,
			DeadTup:    nDeadTup,
		})
	}

	return result, nil
}
