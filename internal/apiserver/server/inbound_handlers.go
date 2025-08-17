package server

import (
	"errors"

	"github.com/gofiber/fiber/v3"
	v1 "github.com/vayzur/apadana/pkg/api/v1"
	"github.com/vayzur/apadana/pkg/errs"
)

func (s *Server) GetInbound(c fiber.Ctx) error {
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

	inbound, err := s.inboundService.GetInbound(c.RequestCtx(), node, tag)
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

	inbound := new(v1.Inbound)
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

	return c.Status(fiber.StatusCreated).JSON(inbound)
}

func (s *Server) GetAllInbounds(c fiber.Ctx) error {
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

	return c.SendStatus(fiber.StatusNoContent)
}
