package server

import (
	"context"
	"errors"

	"github.com/gofiber/fiber/v3"
	"github.com/vayzur/apadana/pkg/errs"
)

func (s *Server) AddInbound(c fiber.Ctx) error {
	b := c.Body()

	if err := s.xrayClient.AddInbound(context.Background(), b); err != nil {
		if errors.Is(err, errs.ErrConflict) {
			return c.SendStatus(fiber.StatusConflict)
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusCreated)
}

func (s *Server) RemoveInbound(c fiber.Ctx) error {
	tag := c.Params("tag")

	if tag == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "tag parameter is required"})
	}

	if err := s.xrayClient.RemoveInbound(context.Background(), tag); err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return c.SendStatus(fiber.StatusNotFound)
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
