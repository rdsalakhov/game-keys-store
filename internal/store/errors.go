package store

import "errors"

var (
	ErrRecordNotFound = errors.New("sql: no rows in result set")
)
