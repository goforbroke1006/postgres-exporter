package internal

type StatUserTablesRepository interface {
	FindTopDeadTuples(limit uint) (result []StatUserTablesRow, err error)
}
