package internal

import "time"

type StatUserTablesRow struct {
	SchemaName string
	RelName    string
	LiveTup    int64
	DeadTup    int64
}

type Lock struct {
	Database string
	Relation string

	LockType string

	TransactionStarted time.Time
	QueryStarted       time.Time

	Query string
}
