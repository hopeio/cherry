package db

const (
	TmFmtWithMS = "2006-01-02 15:04:05.999"
	NullStr     = "NULL"
)

const (
	ColumnDeletedAt = "deleted_at"
	ColumnId        = "id"
	ColumnName      = "name"
)

const (
	ExprEqual    = " = ?"
	ExprNotEqual = " != ?"
	ExprGreater  = " > ?"
)

const (
	Mysql    = "mysql"
	Postgres = "postgres"
	Sqlite   = "sqlite"
)
