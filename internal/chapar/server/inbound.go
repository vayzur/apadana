package server

import (
	"github.com/gofiber/fiber/v3"
	zlog "github.com/rs/zerolog/log"
	metav1 "github.com/vayzur/apadana/pkg/apis/meta/v1"
	satrapv1 "github.com/vayzur/apadana/pkg/apis/satrap/v1"
	"github.com/vayzur/apadana/pkg/errs"
)

func (s *Server) GetInbound(c fiber.Ctx) error {
	params, err := s.requiredParams(c, "nodeName", "tag")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&errs.Error{
			Kind:    errs.KindInvalid,
			Reason:  errs.ReasonMissingParam,
			Message: err.Error(),
		})
	}

	nodeName := params["nodeName"]
	tag := params["tag"]

	inbound, err := s.inboundService.GetInbound(c.RequestCtx(), nodeName, tag)
	if err != nil {
		zlog.Error().Err(err).Str("component", "chapar").Str("resource", "inbound").Str("action", "get").Str("nodeName", nodeName).Str("tag", tag).Msg("failed")
		return errs.HandleAPIError(c, err)
	}

	zlog.Info().Str("component", "chapar").Str("resource", "inbound").Str("action", "get").Str("nodeName", nodeName).Str("tag", tag).Msg("retrieved")
	return c.Status(fiber.StatusOK).JSON(inbound)
}

func (s *Server) CreateInbound(c fiber.Ctx) error {
	nodeName := c.Params("nodeName")
	if nodeName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": errs.ErrInvalidNode,
			},
		)
	}

	inbound := &satrapv1.Inbound{}
	if err := c.Bind().JSON(inbound); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": &errs.Error{
					Kind:    errs.KindInvalid,
					Reason:  errs.ReasonUnmarshalFailed,
					Message: err.Error(),
				},
			},
		)
	}

	tag := inbound.Spec.Config.Tag

	if err := s.inboundService.CreateInbound(c.RequestCtx(), nodeName, inbound); err != nil {
		zlog.Error().Err(err).Str("component", "chapar").Str("resource", "inbound").Str("action", "create").Str("nodeName", nodeName).Str("tag", tag).Msg("failed")
		return errs.HandleAPIError(c, err)
	}

	zlog.Info().Str("component", "chapar").Str("resource", "inbound").Str("action", "create").Str("nodeName", nodeName).Str("tag", tag).Msg("created")
	return c.Status(fiber.StatusCreated).JSON(inbound)
}

func (s *Server) GetInbounds(c fiber.Ctx) error {
	nodeName := c.Params("nodeName")
	if nodeName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": errs.ErrInvalidNode,
			},
		)
	}

	inbounds, err := s.inboundService.GetInbounds(c.RequestCtx(), nodeName)
	if err != nil {
		zlog.Error().Err(err).Str("component", "chapar").Str("resource", "inbounds").Str("action", "list").Str("nodeName", nodeName).Int("count", len(inbounds)).Msg("failed")
		return errs.HandleAPIError(c, err)
	}

	zlog.Info().Str("component", "chapar").Str("resource", "inbounds").Str("action", "list").Str("nodeName", nodeName).Int("count", len(inbounds)).Msg("retrieved")
	return c.Status(fiber.StatusOK).JSON(inbounds)
}

func (s *Server) DeleteInbound(c fiber.Ctx) error {
	params, err := s.requiredParams(c, "nodeName", "tag")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&errs.Error{
			Kind:    errs.KindInvalid,
			Reason:  errs.ReasonMissingParam,
			Message: err.Error(),
		})
	}

	nodeName := params["nodeName"]
	tag := params["tag"]

	if err := s.inboundService.DeleteInbound(c.RequestCtx(), nodeName, tag); err != nil {
		zlog.Error().Err(err).Str("component", "chapar").Str("resource", "inbound").Str("action", "delete").Str("nodeName", nodeName).Str("tag", tag).Msg("failed")
		return errs.HandleAPIError(c, err)
	}

	zlog.Info().Str("component", "chapar").Str("resource", "inbound").Str("action", "delete").Str("nodeName", nodeName).Str("tag", tag).Msg("deleted")
	return c.SendStatus(fiber.StatusNoContent)
}

func (s *Server) CreateUser(c fiber.Ctx) error {
	params, err := s.requiredParams(c, "nodeName", "tag")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&errs.Error{
			Kind:    errs.KindInvalid,
			Reason:  errs.ReasonMissingParam,
			Message: err.Error(),
		})
	}

	nodeName := params["nodeName"]
	tag := params["tag"]

	user := &satrapv1.InboundUser{}
	if err := c.Bind().JSON(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	proto := user.Spec.Type
	email := user.Spec.Email

	if err := s.inboundService.CreateUser(c.RequestCtx(), nodeName, tag, user); err != nil {
		zlog.Error().Err(err).Str("component", "chapar").Str("resource", "inboundUser").Str("action", "create").Str("nodeName", nodeName).Str("tag", tag).Str("protocol", proto).Str("email", email).Msg("failed")
		return errs.HandleAPIError(c, err)
	}

	zlog.Info().Str("component", "chapar").Str("resource", "inboundUser").Str("action", "create").Str("nodeName", nodeName).Str("tag", tag).Str("protocol", proto).Str("email", email).Msg("created")
	return c.Status(fiber.StatusCreated).JSON(user)
}

func (s *Server) DeleteUser(c fiber.Ctx) error {
	params, err := s.requiredParams(c, "nodeName", "tag", "email")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&errs.Error{
			Kind:    errs.KindInvalid,
			Reason:  errs.ReasonMissingParam,
			Message: err.Error(),
		})
	}

	nodeName := params["nodeName"]
	tag := params["tag"]
	email := params["email"]

	if err := s.inboundService.DeleteUser(c.RequestCtx(), nodeName, tag, params["email"]); err != nil {
		zlog.Error().Err(err).Str("component", "chapar").Str("resource", "inboundUser").Str("action", "delete").Str("nodeName", nodeName).Str("tag", tag).Str("email", email).Msg("failed")
		return errs.HandleAPIError(c, err)
	}

	zlog.Info().Str("component", "chapar").Str("resource", "inboundUser").Str("action", "delete").Str("nodeName", nodeName).Str("tag", tag).Str("email", email).Msg("deleted")
	return c.SendStatus(fiber.StatusNoContent)
}

func (s *Server) GetInboundUsers(c fiber.Ctx) error {
	params, err := s.requiredParams(c, "nodeName", "tag")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&errs.Error{
			Kind:    errs.KindInvalid,
			Reason:  errs.ReasonMissingParam,
			Message: err.Error(),
		})
	}

	nodeName := params["nodeName"]
	tag := params["tag"]

	users, err := s.inboundService.GetUsers(c.RequestCtx(), nodeName, tag)
	if err != nil {
		zlog.Error().Err(err).Str("component", "chapar").Str("resource", "inboundUser").Str("action", "list").Str("nodeName", nodeName).Str("tag", tag).Int("count", len(users)).Msg("failed")
		return errs.HandleAPIError(c, err)
	}

	zlog.Info().Str("component", "chapar").Str("resource", "inboundUser").Str("action", "list").Str("nodeName", nodeName).Str("tag", tag).Int("count", len(users)).Msg("retrieved")
	return c.Status(fiber.StatusOK).JSON(users)
}

func (s *Server) UpdateInboundMetadata(c fiber.Ctx) error {
	params, err := s.requiredParams(c, "nodeName", "tag")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&errs.Error{
			Kind:    errs.KindInvalid,
			Reason:  errs.ReasonMissingParam,
			Message: err.Error(),
		})
	}

	nodeName := params["nodeName"]
	tag := params["tag"]

	newMetadata := &metav1.ObjectMeta{}
	if err := c.Bind().JSON(newMetadata); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	if err := s.inboundService.UpdateInboundMetadata(c.RequestCtx(), nodeName, tag, newMetadata); err != nil {
		zlog.Error().Err(err).Str("component", "chapar").Str("resource", "inbound").Str("action", "update").Str("nodeName", nodeName).Str("tag", tag).Msg("failed")
		return errs.HandleAPIError(c, err)
	}

	zlog.Info().Str("component", "chapar").Str("resource", "inbound").Str("action", "update").Str("nodeName", nodeName).Str("tag", tag).Msg("updated")
	return c.SendStatus(fiber.StatusOK)
}

func (s *Server) UpdateInboundUserMetadata(c fiber.Ctx) error {
	params, err := s.requiredParams(c, "nodeName", "tag", "email")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&errs.Error{
			Kind:    errs.KindInvalid,
			Reason:  errs.ReasonMissingParam,
			Message: err.Error(),
		})
	}

	nodeName := params["nodeName"]
	tag := params["tag"]
	email := params["email"]

	newMetadata := &metav1.ObjectMeta{}
	if err := c.Bind().JSON(newMetadata); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	if err := s.inboundService.UpdateUserMetadata(c.RequestCtx(), nodeName, tag, email, newMetadata); err != nil {
		zlog.Error().Err(err).Str("component", "chapar").Str("resource", "inboundUser").Str("action", "update").Str("nodeName", nodeName).Str("tag", tag).Str("email", email).Msg("failed")
		return errs.HandleAPIError(c, err)
	}

	zlog.Info().Str("component", "chapar").Str("resource", "inboundUser").Str("action", "update").Str("nodeName", nodeName).Str("tag", tag).Str("email", email).Msg("updated")
	return c.SendStatus(fiber.StatusOK)
}

func (s *Server) CountInbounds(c fiber.Ctx) error {
	nodeName := c.Params("nodeName")
	if nodeName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": errs.ErrInvalidNode,
			},
		)
	}

	count, err := s.inboundService.CountInbounds(c.RequestCtx(), nodeName)
	if err != nil {
		zlog.Error().Err(err).Str("component", "chapar").Str("resource", "inbound").Str("action", "count").Str("nodeName", nodeName).Uint32("count", count).Msg("failed")
		return errs.HandleAPIError(c, err)
	}

	countResp := &satrapv1.Count{
		Value: count,
	}

	zlog.Info().Str("component", "chapar").Str("resource", "inbound").Str("action", "count").Str("nodeName", nodeName).Uint32("count", count).Msg("retrieved")
	return c.Status(fiber.StatusOK).JSON(countResp)
}

func (s *Server) CountInboundUsers(c fiber.Ctx) error {
	params, err := s.requiredParams(c, "nodeName", "tag")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&errs.Error{
			Kind:    errs.KindInvalid,
			Reason:  errs.ReasonMissingParam,
			Message: err.Error(),
		})
	}

	nodeName := params["nodeName"]
	tag := params["tag"]

	count, err := s.inboundService.CountUsers(c.RequestCtx(), nodeName, tag)
	if err != nil {
		zlog.Error().Err(err).Str("component", "chapar").Str("resource", "inboundUser").Str("action", "count").Str("nodeName", nodeName).Str("tag", tag).Uint32("count", count).Msg("failed")
		return errs.HandleAPIError(c, err)
	}

	countResp := &satrapv1.Count{
		Value: count,
	}

	zlog.Info().Str("component", "chapar").Str("resource", "inbound").Str("action", "count").Str("nodeName", nodeName).Str("tag", tag).Uint32("count", count).Msg("retrieved")
	return c.Status(fiber.StatusOK).JSON(countResp)
}
