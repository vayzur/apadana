package server

import (
	"context"
	"errors"

	zlog "github.com/rs/zerolog/log"
	satrapv1 "github.com/vayzur/apadana/pkg/api/satrap/v1"
	"github.com/xtls/xray-core/infra/conf"

	"github.com/gofiber/fiber/v3"
	"github.com/vayzur/apadana/pkg/errs"
)

func (s *Server) AddInbound(c fiber.Ctx) error {
	inboundConfig := &conf.InboundDetourConfig{}
	if err := c.Bind().JSON(inboundConfig); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}
	if err := s.xrayClient.AddInbound(context.Background(), inboundConfig); err != nil {
		zlog.Error().Err(err).Str("component", "satrap").Str("resource", "inbound").Str("action", "create").Msg("failed")
		if errors.Is(err, errs.ErrInboundConflict) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": err,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err,
		})
	}
	return c.SendStatus(fiber.StatusCreated)
}

func (s *Server) RemoveInbound(c fiber.Ctx) error {
	tag := c.Params("tag")
	if tag == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "tag parameter is required"})
	}
	if err := s.xrayClient.RemoveInbound(context.Background(), tag); err != nil {
		zlog.Error().Err(err).Str("component", "satrap").Str("resource", "inbound").Str("action", "delete").Str("tag", tag).Msg("failed")
		if errors.Is(err, errs.ErrInboundNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err,
		})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (s *Server) AddUser(c fiber.Ctx) error {
	tag := c.Params("tag")
	if tag == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "tag parameter is required",
		})
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
		zlog.Error().Err(err).Str("component", "satrap").Str("resource", "user").Str("action", "create").Str("tag", tag).Msg("failed")
		if errors.Is(err, errs.ErrUserConflict) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": err,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err,
		})
	}
	return c.SendStatus(fiber.StatusCreated)
}

func (s *Server) RemoveUser(c fiber.Ctx) error {
	params, err := s.requiredParams(c, "tag", "email")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if err := s.xrayClient.RemoveUser(context.Background(), params["tag"], params["email"]); err != nil {
		zlog.Error().Err(err).Str("component", "satrap").Str("resource", "user").Str("action", "delete").Str("tag", params["tag"]).Str("email", params["email"]).Msg("failed")
		if errors.Is(err, errs.ErrUserNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err,
		})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
