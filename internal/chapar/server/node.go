package server

import (
	"context"
	"net/http"

	"github.com/gofiber/fiber/v3"
	zlog "github.com/rs/zerolog/log"
	corev1 "github.com/vayzur/apadana/pkg/apis/core/v1"
	metav1 "github.com/vayzur/apadana/pkg/apis/meta/v1"
	"github.com/vayzur/apadana/pkg/errs"
)

func (s *Server) GetNodes(c fiber.Ctx) error {
	nodes, err := s.nodeService.GetNodes(c.RequestCtx())
	if err != nil {
		return errs.HandleAPIError(c, err)
	}

	zlog.Info().Str("component", "chapar").Str("resource", "nodes").Str("action", "list").Int("count", len(nodes)).Msg("retrieved")
	return c.Status(http.StatusOK).JSON(nodes)
}

func (s *Server) GetActiveNodes(c fiber.Ctx) error {
	nodes, err := s.nodeService.GetActiveNodes(c.RequestCtx())
	if err != nil {
		return errs.HandleAPIError(c, err)
	}

	zlog.Info().Str("component", "chapar").Str("resource", "nodes").Str("action", "list").Int("count", len(nodes)).Msg("retrieved")
	return c.Status(http.StatusOK).JSON(nodes)
}

func (s *Server) GetNode(c fiber.Ctx) error {
	nodeName := c.Params("nodeName")
	if nodeName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": errs.ErrInvalidNode,
			},
		)
	}

	node, err := s.nodeService.GetNode(c.RequestCtx(), nodeName)
	if err != nil {
		return errs.HandleAPIError(c, err)
	}

	zlog.Info().Str("component", "chapar").Str("resource", "node").Str("action", "get").Str("nodeName", nodeName).Msg("retrieved")
	return c.Status(http.StatusOK).JSON(node)
}

func (s *Server) CreateNode(c fiber.Ctx) error {
	node := &corev1.Node{}
	if err := c.Bind().JSON(node); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&errs.Error{
			Kind:    errs.KindInvalid,
			Reason:  errs.ReasonUnmarshalFailed,
			Message: err.Error(),
		})
	}

	if err := s.nodeService.CreateNode(c.RequestCtx(), node); err != nil {
		return errs.HandleAPIError(c, err)
	}

	zlog.Info().Str("component", "chapar").Str("resource", "node").Str("action", "create").Str("nodeName", node.Metadata.Name).Msg("created")
	return c.Status(http.StatusCreated).JSON(node)
}

func (s *Server) DeleteNode(c fiber.Ctx) error {
	nodeName := c.Params("nodeName")
	if nodeName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": errs.ErrInvalidNode,
			},
		)
	}

	if err := s.nodeService.DeleteNode(context.Background(), nodeName); err != nil {
		return errs.HandleAPIError(c, err)
	}

	zlog.Info().Str("component", "chapar").Str("resource", "node").Str("action", "delete").Str("nodeName", nodeName).Msg("deleted")
	return c.SendStatus(fiber.StatusNoContent)
}

func (s *Server) UpdateNodeStatus(c fiber.Ctx) error {
	nodeName := c.Params("nodeName")
	if nodeName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": errs.ErrInvalidNode,
			},
		)
	}

	newStatus := &corev1.NodeStatus{}
	if err := c.Bind().JSON(newStatus); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	if err := s.nodeService.UpdateNodeStatus(c.RequestCtx(), nodeName, newStatus); err != nil {
		return errs.HandleAPIError(c, err)
	}

	zlog.Info().Str("component", "chapar").Str("resource", "node").Str("action", "update").Str("nodeName", nodeName).Msg("updated")
	return c.SendStatus(fiber.StatusOK)
}

func (s *Server) UpdateNodeMetadata(c fiber.Ctx) error {
	nodeName := c.Params("nodeName")
	if nodeName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": errs.ErrInvalidNode,
			},
		)
	}

	newMetadata := &metav1.ObjectMeta{}
	if err := c.Bind().JSON(newMetadata); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	if err := s.nodeService.UpdateNodeMetadata(c.RequestCtx(), nodeName, newMetadata); err != nil {
		return errs.HandleAPIError(c, err)
	}

	zlog.Info().Str("component", "chapar").Str("resource", "node").Str("action", "update").Str("nodeName", nodeName).Msg("updated")
	return c.SendStatus(fiber.StatusOK)
}
