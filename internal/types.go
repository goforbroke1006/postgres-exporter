package internal

type StatUserTablesRow struct {
	SchemaName string
	RelName    string
	LiveTup    int64
	DeadTup    int64
}
