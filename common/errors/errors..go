package errors

import "github.com/pkg/errors"

var (
	ErrAccessDenied = errors.New("access denied")
)
