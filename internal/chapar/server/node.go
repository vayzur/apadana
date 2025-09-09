package server

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	zlog "github.com/rs/zerolog/log"
	corev1 "github.com/vayzur/apadana/pkg/api/core/v1"
	"github.com/vayzur/apadana/pkg/errs"
)

func (s *Server) GetNodes(c fiber.Ctx) error {
	nodes, err := s.nodeService.GetNodes(c.RequestCtx())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	zlog.Info().Str("component", "chapar").Str("resource", "nodes").Str("action", "list").Int("count", len(nodes)).Msg("retrieved")
	return c.Status(http.StatusOK).JSON(nodes)
}

func (s *Server) GetActiveNodes(c fiber.Ctx) error {
	nodes, err := s.nodeService.GetActiveNodes(c.RequestCtx())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	zlog.Info().Str("component", "chapar").Str("resource", "nodes").Str("action", "list").Int("count", len(nodes)).Msg("retrieved")
	return c.Status(http.StatusOK).JSON(nodes)
}

func (s *Server) GetNode(c fiber.Ctx) error {
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

	zlog.Info().Str("component", "chapar").Str("resource", "node").Str("action", "get").Str("nodeID", nodeID).Msg("retrieved")
	return c.Status(http.StatusOK).JSON(node)
}

func (s *Server) CreateNode(c fiber.Ctx) error {
	node := &corev1.Node{}
	if err := c.Bind().JSON(node); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	node.Metadata.ID = uuid.NewString()
	node.Metadata.CreationTimestamp = time.Now()

	if err := s.nodeService.CreateNode(c.RequestCtx(), node); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	zlog.Info().Str("component", "chapar").Str("resource", "node").Str("action", "create").Str("nodeID", node.Metadata.ID).Msg("created")
	return c.Status(http.StatusCreated).JSON(node)
}

func (s *Server) DeleteNode(c fiber.Ctx) error {
	nodeID := c.Params("nodeID")
	if nodeID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": "nodeID parameter is required",
			},
		)
	}

	if err := s.nodeService.DeleteNode(context.Background(), nodeID); err != nil {
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

	zlog.Info().Str("component", "chapar").Str("resource", "node").Str("action", "delete").Str("nodeID", nodeID).Msg("deleted")
	return c.SendStatus(fiber.StatusNoContent)
}

func (s *Server) UpdateNodeStatus(c fiber.Ctx) error {
	nodeID := c.Params("nodeID")
	if nodeID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": "nodeID parameter is required",
			},
		)
	}

	nodeStatus := &corev1.NodeStatus{}
	if err := c.Bind().JSON(nodeStatus); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	if err := s.nodeService.UpdateNodeStatus(c.RequestCtx(), nodeID, nodeStatus); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	zlog.Info().Str("component", "chapar").Str("resource", "node").Str("action", "update").Str("nodeID", nodeID).Msg("updated")
	return c.SendStatus(fiber.StatusOK)
}

func (s *Server) UpdateNodeMetadata(c fiber.Ctx) error {
	nodeID := c.Params("nodeID")
	if nodeID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": "nodeID parameter is required",
			},
		)
	}

	metadata := &corev1.NodeMetadata{}
	if err := c.Bind().JSON(metadata); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	if err := s.nodeService.UpdateNodeMetadata(c.RequestCtx(), nodeID, metadata); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	zlog.Info().Str("component", "chapar").Str("resource", "node").Str("action", "update").Str("nodeID", nodeID).Msg("updated")
	return c.SendStatus(fiber.StatusOK)
}

func (s *Server) UpdateNodeSpec(c fiber.Ctx) error {
	nodeID := c.Params("nodeID")
	if nodeID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": "nodeID parameter is required",
			},
		)
	}

	spec := &corev1.NodeSpec{}
	if err := c.Bind().JSON(spec); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	if err := s.nodeService.UpdateNodeSpec(c.RequestCtx(), nodeID, spec); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	zlog.Info().Str("component", "chapar").Str("resource", "node").Str("action", "update").Str("nodeID", nodeID).Msg("updated")
	return c.SendStatus(fiber.StatusOK)
}
