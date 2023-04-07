package data

import "github.com/pkg/errors"

var (
	ErrNotFound        = errors.New("record not found")
	ErrDuplicateRecord = errors.New("this record is already present")
)
