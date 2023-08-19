package data

import sq "github.com/Masterminds/squirrel"

const (
	OrderAsc  = "asc"
	OrderDesc = "desc"
)

const (
	NewsPost = "news"
)

var (
	BasicSqlizer = sq.Eq{"1": "1"}
)
