package server

import (
	"context"
	"errors"

	satrapv1 "github.com/vayzur/apadana/pkg/api/satrap/v1"

	"github.com/gofiber/fiber/v3"
	zlog "github.com/rs/zerolog/log"
	"github.com/vayzur/apadana/pkg/errs"
)

func (s *Server) InboundsCount(c fiber.Ctx) error {
	inbounds, err := s.xrayClient.ListInbounds(context.Background())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	count := satrapv1.Count{
		Value: int32(len(inbounds)),
	}

	return c.Status(fiber.StatusOK).JSON(count)
}

func (s *Server) AddInbound(c fiber.Ctx) error {
	b := c.Body()
	zlog.Info().Str("component", "satrap").Str("resource", "inbound").Str("action", "create").Str("body", string(b)).Msg("received")

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

func (s *Server) AddUser(c fiber.Ctx) error {
	tag := c.Params("tag")
	if tag == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "tag parameter is required"})
	}

	user := &satrapv1.InboundUser{}
	if err := c.Bind().JSON(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	account, err := user.ToAccount()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := s.xrayClient.AddUser(c.RequestCtx(), tag, user.Email, account); err != nil {
		if errors.Is(err, errs.ErrConflict) {
			return c.SendStatus(fiber.StatusConflict)
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (s *Server) RemoveUser(c fiber.Ctx) error {
	params, err := s.requiredParams(c, "tag", "email")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := s.xrayClient.RemoveUser(context.Background(), params["tag"], params["email"]); err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return c.SendStatus(fiber.StatusNotFound)
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
