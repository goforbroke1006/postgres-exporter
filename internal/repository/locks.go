package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"

	"postgres-exporter/internal"
)

func NewLockRepository(conn *pgxpool.Pool) *lockRepository {
	return &lockRepository{
		conn: conn,
	}
}

var _ internal.LockRepository = &lockRepository{}

type lockRepository struct {
	conn *pgxpool.Pool
}

func (r lockRepository) Find(databaseName string) ([]internal.Lock, error) {
	query := fmt.Sprintf(`
		SELECT 
			db.datname		AS datname,
			cl.relname		AS relation,
			pl.locktype		AS locktype,
			psa.xact_start	AS tx_start,
			psa.query_start	AS q_start,
			psa.query		AS query
		FROM pg_locks pl
		INNER JOIN pg_database     db  ON pl.database = db.oid
		INNER JOIN pg_class        cl  ON pl.relation = cl.oid
		LEFT JOIN pg_stat_activity psa ON pl.pid      = psa.pid
		WHERE 
			db.datname = '%s'
			AND cl.relname NOT LIKE 'pg\_%%'
		;`, databaseName)

	rows, err := r.conn.Query(context.TODO(), query)
	if err != nil {
		return nil, err
	}

	var (
		datname     string
		relation    string
		lockType    string
		txStart     time.Time
		qStart      time.Time
		recentQuery string
	)

	var result []internal.Lock

	for rows.Next() {
		err := rows.Scan(&datname, &relation, &lockType, &txStart, &qStart, &recentQuery)
		if err != nil {
			return nil, err
		}

		result = append(result, internal.Lock{
			Database: datname,
			Relation: relation,

			LockType: lockType,

			TransactionStarted: txStart,
			QueryStarted:       qStart,

			Query: recentQuery,
		})
	}

	return result, err
}
