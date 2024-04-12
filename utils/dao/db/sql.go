package dbi

const (
	INSERT    = "INSERT"
	SELECT    = "SELECT"
	UPDATE    = "UPDATE"
	DELETE    = "DELETE"
	LEFTJOIN  = "LEFT JOIN"
	RIGHTJOIN = "RIGHT JOIN"
	INNERJOIN = "INNER JOIN"
)

const (
	NotDeleted     = ColumnDeletedAt + " IS " + NullStr
	WithNotDeleted = ` AND ` + NotDeleted
)
