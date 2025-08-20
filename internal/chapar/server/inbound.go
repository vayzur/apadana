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

	node, err := s.nodeService.GetNode(c.RequestCtx(), params["nodeID"])
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(
				fiber.Map{
					"error": err.Error(),
				},
			)
		}
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	inbound, err := s.inboundService.GetInbound(c.RequestCtx(), node, params["tag"])
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

	node, err := s.nodeService.GetNode(c.RequestCtx(), nodeID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(
				fiber.Map{
					"error": err.Error(),
				},
			)
		}
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	if err := s.inboundService.AddInbound(c.RequestCtx(), inbound, node); err != nil {
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

	node, err := s.nodeService.GetNode(c.RequestCtx(), nodeID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(
				fiber.Map{
					"error": err.Error(),
				},
			)
		}
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	inbounds, err := s.inboundService.ListInbounds(c.RequestCtx(), node)
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

	node, err := s.nodeService.GetNode(c.RequestCtx(), nodeID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(
				fiber.Map{
					"error": err.Error(),
				},
			)
		}
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	if err := s.inboundService.DelInbound(c.RequestCtx(), node, tag); err != nil {
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

	var req satrapv1.CreateUserRequest
	if err := c.Bind().JSON(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	node, err := s.nodeService.GetNode(c.RequestCtx(), params["nodeID"])
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(
				fiber.Map{
					"error": err.Error(),
				},
			)
		}
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	if err := s.inboundService.AddUser(c.RequestCtx(), node, params["tag"], req); err != nil {
		if errors.Is(err, errs.ErrConflict) {
			return c.SendStatus(fiber.StatusConflict)
		}
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	zlog.Info().Str("component", "chapar").Str("resource", "inboundUser").Str("action", "create").Str("nodeID", params["nodeID"]).Str("tag", params["tag"]).Str("protocol", req.Type).Str("email", req.Email).Str("account", string(req.Account)).Msg("created")
	return c.Status(fiber.StatusCreated).JSON(req)
}

func (s *Server) DeleteUser(c fiber.Ctx) error {
	params, err := s.requiredParams(c, "nodeID", "tag", "email")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	node, err := s.nodeService.GetNode(c.RequestCtx(), params["nodeID"])
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(
				fiber.Map{
					"error": err.Error(),
				},
			)
		}
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	if err := s.inboundService.DelUser(c.RequestCtx(), node, params["tag"], params["email"]); err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return c.SendStatus(fiber.StatusNotFound)
		}
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	zlog.Info().Str("component", "chapar").Str("resource", "inboundUser").Str("action", "delete").Str("nodeID", params["nodeID"]).Str("tag", params["tag"]).Msg("deleted")
	return c.SendStatus(fiber.StatusNoContent)
}
