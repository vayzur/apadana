package errs

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	satrapv1 "github.com/vayzur/apadana/pkg/api/satrap/v1"
	"google.golang.org/grpc/status"
)

func HandleXrayError(err error, resourceType satrapv1.Resource) error {
	if err == nil {
		return nil
	}

	s, _ := status.FromError(err)
	msg := s.Message()

	switch {
	case strings.Contains(msg, "existing tag"),
		strings.Contains(msg, "already exists"):
		switch resourceType {
		case satrapv1.ResourceInbound:
			return ErrInboundConflict
		case satrapv1.ResourceUser:
			return ErrUserConflict
		}
	case strings.Contains(msg, "not enough information for making a decision"),
		strings.Contains(msg, "handler not found"),
		strings.Contains(msg, "not found"):
		switch resourceType {
		case satrapv1.ResourceInbound:
			return ErrInboundNotFound
		case satrapv1.ResourceUser:
			return ErrUserNotFound
		}
	}

	return New(KindInternal, ReasonUnknown, "runtime operation failed", nil, err)
}

func HandleAPIError(c fiber.Ctx, err error) error {
	e, ok := err.(*Error)
	if !ok {
		e = &Error{
			Kind:    KindInternal,
			Reason:  ReasonUnknown,
			Message: err.Error(),
		}
	}

	var status int

	switch e.Kind {
	case KindNotFound:
		status = fiber.StatusNotFound
	case KindConflict:
		status = fiber.StatusConflict
	case KindCapacityExceeded:
		status = fiber.StatusTooManyRequests
	case KindInvalid:
		status = fiber.StatusBadRequest
	default:
		status = fiber.StatusInternalServerError
	}

	return c.Status(status).JSON(fiber.Map{
		"error": err,
	})
}
