package internal

type StatUserTablesRepository interface {
	FindTopDeadTuples(limit uint) (result []StatUserTablesRow, err error)
	FindTopLiveTuples(limit uint) (result []StatUserTablesRow, err error)
}

type LockRepository interface {
	Find(databaseName string) ([]Lock, error)
}
