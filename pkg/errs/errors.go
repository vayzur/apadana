package errs

import (
	"fmt"
)

type ErrorKind string
type ErrorReason string

const (
	KindNotFound         ErrorKind = "NotFound"
	KindInvalid          ErrorKind = "Invalid"
	KindConflict         ErrorKind = "Conflict"
	KindInternal         ErrorKind = "Internal"
	KindCapacityExceeded ErrorKind = "CapacityExceeded"
)

const (
	ReasonMissingParam            ErrorReason = "MissingParam"
	ReasonUnknown                 ErrorReason = "Unknown"
	ReasonMarshalFailed           ErrorReason = "MarshalFailed"
	ReasonUnmarshalFailed         ErrorReason = "UnmarshalFailed"
	ReasonNodeNotFound            ErrorReason = "NodeNotFound"
	ReasonInboundConflict         ErrorReason = "InboundConflict"
	ReasonInboundNotFound         ErrorReason = "InboundNotFound"
	ReasonUserConflict            ErrorReason = "UserConflict"
	ReasonUserNotFound            ErrorReason = "UserNotFound"
	ReasonNodeCapacityExceeded    ErrorReason = "NodeCapacityExceeded"
	ReasonInboundCapacityExceeded ErrorReason = "InboundCapacityExceeded"
	ReasonResourceNotFound        ErrorReason = "ResourceNotFound"
)

type Error struct {
	Kind    ErrorKind         `json:"kind,omitempty"`
	Reason  ErrorReason       `json:"reason,omitempty"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
	Cause   error             `json:"-"`
}

var (
	ErrNodeNotFound            = &Error{Kind: KindNotFound, Reason: ReasonNodeNotFound, Message: "node not found"}
	ErrInboundConflict         = &Error{Kind: KindConflict, Reason: ReasonInboundConflict, Message: "inbound already exists"}
	ErrInboundNotFound         = &Error{Kind: KindNotFound, Reason: ReasonInboundNotFound, Message: "inbound not found"}
	ErrUserConflict            = &Error{Kind: KindConflict, Reason: ReasonUserConflict, Message: "user already exists"}
	ErrUserNotFound            = &Error{Kind: KindNotFound, Reason: ReasonUserNotFound, Message: "user not found"}
	ErrNodeCapacityExceeded    = &Error{Kind: KindCapacityExceeded, Reason: ReasonNodeCapacityExceeded, Message: "node capacity exceeded"}
	ErrInboundCapacityExceeded = &Error{Kind: KindCapacityExceeded, Reason: ReasonInboundCapacityExceeded, Message: "inbound capacity exceeded"}
	ErrInvalidNode             = &Error{Kind: KindInvalid, Reason: ReasonMissingParam, Message: "nodeName cannot be empty"}
	ErrInvalidInbound          = &Error{Kind: KindInvalid, Reason: ReasonMissingParam, Message: "tag cannot be empty"}
	ErrInvalidUser             = &Error{Kind: KindInvalid, Reason: ReasonMissingParam, Message: "email cannot be empty"}
	ErrResourceNotFound        = &Error{Kind: KindNotFound, Reason: ReasonResourceNotFound, Message: "resource not found"}
)

func (e *Error) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

func New(kind ErrorKind, reason ErrorReason, msg string, fields map[string]string, Cause error) *Error {
	return &Error{
		Kind:    kind,
		Reason:  reason,
		Message: msg,
		Fields:  fields,
		Cause:   Cause,
	}
}
