package server

import (
	"errors"

	"github.com/gofiber/fiber/v3"
	zlog "github.com/rs/zerolog/log"
	satrapv1 "github.com/vayzur/apadana/pkg/api/satrap/v1"
	"github.com/vayzur/apadana/pkg/errs"
)

func (s *Server) GetInbound(c fiber.Ctx) error {
	params, err := s.requiredParams(c, "nodeID", "tag")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	inbound, err := s.inboundService.GetInbound(c.RequestCtx(), params["nodeID"], params["tag"])
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return c.SendStatus(fiber.StatusNotFound)
		}
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	zlog.Info().Str("component", "chapar").Str("resource", "inbound").Str("action", "get").Str("nodeID", params["nodeID"]).Str("tag", params["tag"]).Msg("retrieved")
	return c.Status(fiber.StatusOK).JSON(inbound)
}

func (s *Server) CreateInbound(c fiber.Ctx) error {
	nodeID := c.Params("nodeID")
	if nodeID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": "nodeID parameter is required",
			},
		)
	}

	inbound := new(satrapv1.Inbound)
	if err := c.Bind().JSON(inbound); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	if err := s.inboundService.CreateInbound(c.RequestCtx(), nodeID, inbound); err != nil {
		if errors.Is(err, errs.ErrNodeCapacity) {
			return c.SendStatus(fiber.StatusTooManyRequests)
		}
		if errors.Is(err, errs.ErrConflict) {
			return c.SendStatus(fiber.StatusConflict)
		}
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	zlog.Info().Str("component", "chapar").Str("resource", "inbound").Str("action", "create").Str("nodeID", nodeID).Str("tag", inbound.Config.Tag).Msg("created")
	return c.Status(fiber.StatusCreated).JSON(inbound)
}

func (s *Server) GetInbounds(c fiber.Ctx) error {
	nodeID := c.Params("nodeID")
	if nodeID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": "nodeID parameter is required",
			},
		)
	}

	state := c.Query("state", "all") // default = "all"

	var inbounds []*satrapv1.Inbound
	var err error

	switch state {
	case "active":
		inbounds, err = s.inboundService.GetActiveInbounds(c.RequestCtx(), nodeID)
	case "expired":
		inbounds, err = s.inboundService.GetExpiredInbounds(c.RequestCtx(), nodeID)
	case "all":
		fallthrough
	default:
		inbounds, err = s.inboundService.GetInbounds(c.RequestCtx(), nodeID)
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	zlog.Info().Str("component", "chapar").Str("resource", "inbounds").Str("action", "list").Str("nodeID", nodeID).Int("count", len(inbounds)).Msg("retrieved")
	return c.Status(fiber.StatusOK).JSON(inbounds)
}

func (s *Server) DeleteInbound(c fiber.Ctx) error {
	nodeID := c.Params("nodeID")
	if nodeID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": "nodeID parameter is required",
			},
		)
	}

	tag := c.Params("tag")
	if tag == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": "tag parameter is required",
			},
		)
	}

	if err := s.inboundService.DeleteInbound(c.RequestCtx(), nodeID, tag); err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return c.SendStatus(fiber.StatusNotFound)
		}
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	zlog.Info().Str("component", "chapar").Str("resource", "inbound").Str("action", "delete").Str("nodeID", nodeID).Str("tag", tag).Msg("deleted")
	return c.SendStatus(fiber.StatusNoContent)
}

func (s *Server) CreateUser(c fiber.Ctx) error {
	params, err := s.requiredParams(c, "nodeID", "tag")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	user := &satrapv1.InboundUser{}
	if err := c.Bind().JSON(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	if err := s.inboundService.CreateUser(c.RequestCtx(), params["nodeID"], params["tag"], user); err != nil {
		if errors.Is(err, errs.ErrConflict) {
			return c.SendStatus(fiber.StatusConflict)
		}
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	zlog.Info().Str("component", "chapar").Str("resource", "inboundUser").Str("action", "create").Str("nodeID", params["nodeID"]).Str("tag", params["tag"]).Str("protocol", user.Type).Str("email", user.Email).Str("account", string(user.Account)).Msg("created")
	return c.Status(fiber.StatusCreated).JSON(user)
}

func (s *Server) DeleteUser(c fiber.Ctx) error {
	params, err := s.requiredParams(c, "nodeID", "tag", "email")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := s.inboundService.DeleteUser(c.RequestCtx(), params["nodeID"], params["tag"], params["email"]); err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return c.SendStatus(fiber.StatusNotFound)
		}
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	zlog.Info().Str("component", "chapar").Str("resource", "inboundUser").Str("action", "delete").Str("nodeID", params["nodeID"]).Str("tag", params["tag"]).Str("email", params["email"]).Msg("deleted")
	return c.SendStatus(fiber.StatusNoContent)
}

func (s *Server) GetInboundUsers(c fiber.Ctx) error {
	params, err := s.requiredParams(c, "nodeID", "tag")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	state := c.Query("state", "all") // default = "all"

	var users []*satrapv1.InboundUser

	switch state {
	case "active":
		users, err = s.inboundService.GetActiveUsers(c.RequestCtx(), params["nodeID"], params["tag"])
	case "expired":
		users, err = s.inboundService.GetExpiredUsers(c.RequestCtx(), params["nodeID"], params["tag"])
	case "all":
		fallthrough
	default:
		users, err = s.inboundService.GetUsers(c.RequestCtx(), params["nodeID"], params["tag"])
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	zlog.Info().Str("component", "chapar").Str("resource", "inboundUser").Str("action", "list").Str("nodeID", params["nodeID"]).Str("tag", params["tag"]).Int("count", len(users)).Msg("retrieved")
	return c.Status(fiber.StatusOK).JSON(users)
}

func (s *Server) InboundRenew(c fiber.Ctx) error {
	params, err := s.requiredParams(c, "nodeID", "tag")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	renew := &satrapv1.Renew{}
	if err := c.Bind().JSON(renew); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	if err := s.inboundService.InboundRenew(c.RequestCtx(), params["nodeID"], params["tag"], renew); err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return c.SendStatus(fiber.StatusNotFound)
		}
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	zlog.Info().Str("component", "chapar").Str("resource", "inbound").Str("action", "update").Str("nodeID", params["nodeID"]).Str("tag", params["tag"]).Msg("updated")
	return c.SendStatus(fiber.StatusOK)
}

func (s *Server) InboundUserRenew(c fiber.Ctx) error {
	params, err := s.requiredParams(c, "nodeID", "tag", "email")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	renew := &satrapv1.Renew{}
	if err := c.Bind().JSON(renew); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	if err := s.inboundService.InboundUserRenew(c.RequestCtx(), params["nodeID"], params["tag"], params["email"], renew); err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return c.SendStatus(fiber.StatusNotFound)
		}
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	zlog.Info().Str("component", "chapar").Str("resource", "inboundUser").Str("action", "update").Str("nodeID", params["nodeID"]).Str("tag", params["tag"]).Msg("updated")
	return c.SendStatus(fiber.StatusOK)
}

func (s *Server) GetRuntimeInboundsCount(c fiber.Ctx) error {
	nodeID := c.Params("nodeID")
	if nodeID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": "nodeID parameter is required",
			},
		)
	}

	count, err := s.inboundService.GetInboundsCount(c.RequestCtx(), nodeID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return c.SendStatus(fiber.StatusNotFound)
		}
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	zlog.Info().Str("component", "chapar").Str("resource", "inbound").Str("action", "list").Str("nodeID", nodeID).Int32("count", count.Value).Msg("retrieved")
	return c.Status(fiber.StatusOK).JSON(count)
}
