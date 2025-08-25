package errs

import (
	"errors"
	"fmt"
	"maps"
)

var (
	ErrConflict     = errors.New("conflict")
	ErrNotFound     = errors.New("not found")
	ErrNodeCapacity = errors.New("node capacity exceeded")
	ErrUnexpected   = errors.New("unexpected")
)

type Error struct {
	Kind   string
	Msg    string
	Fields map[string]string
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Kind, e.Msg)
}

func New(kind, msg string) *Error {
	return &Error{
		Kind:   kind,
		Msg:    msg,
		Fields: make(map[string]string),
	}
}

func (e *Error) WithField(key, value string) *Error {
	e.Fields[key] = value
	return e
}

func (e *Error) WithFields(fields map[string]string) *Error {
	maps.Copy(e.Fields, fields)
	return e
}
