package repository

import (
	"context"

	"github.com/jackc/pgx/v4"

	"postgres-exporter/internal"
)

func NewStatUserTables(conn *pgx.Conn) *statUserTables {
	return &statUserTables{
		conn: conn,
	}
}

type statUserTables struct {
	conn *pgx.Conn
}

func (r statUserTables) FindTop() (result []internal.StatUserTablesRow, err error) {
	query := `
		SELECT 
			schemaname, 
			relname,
			n_live_tup,
			n_dead_tup 
		FROM pg_stat_user_tables
		ORDER BY n_dead_tup;
	`
	rows, err := r.conn.Query(context.TODO(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			schemaname string
			relname    string
			n_live_tup int64
			n_dead_tup int64
		)

		err = rows.Scan(&schemaname, &relname, &n_live_tup, &n_dead_tup)
		if err != nil {
			return nil, err
		}

		result = append(result, internal.StatUserTablesRow{
			SchemaName: schemaname,
			RelName:    relname,
			LiveTup:    n_live_tup,
			DeadTup:    n_dead_tup,
		})
	}

	return result, nil
}
