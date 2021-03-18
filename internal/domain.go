package internal

type StatUserTablesRow struct {
	SchemaName string `db:"schemaname"`
	RelName    string `db:"relname"`
	LiveTup    int64  `db:"n_live_tup"`
	DeadTup    int64  `db:"n_dead_tup"`
}
