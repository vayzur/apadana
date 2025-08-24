package errs

import "errors"

var (
	ErrConflict     = errors.New("conflict")
	ErrNotFound     = errors.New("not found")
	ErrNodeCapacity = errors.New("node capacity exceeded")
)
