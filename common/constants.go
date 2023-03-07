package common

import sq "github.com/Masterminds/squirrel"

const (
	OrderAsc  = "asc"
	OrderDesc = "desc"
)

var (
	BasicSqlizer = sq.Eq{"1": "1"}
)
