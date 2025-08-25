package server

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
	"github.com/vayzur/apadana/internal/auth"
	"github.com/vayzur/apadana/pkg/service"
)

type Server struct {
	addr           string
	token          string
	prefork        bool
	app            *fiber.App
	inboundService *service.InboundService
	nodeService    *service.NodeService
}

func NewServer(addr, token string, prefork bool, inboundService *service.InboundService, nodeService *service.NodeService) *Server {
	app := fiber.New(fiber.Config{
		CaseSensitive: true,
		StrictRouting: true,
	})
	s := &Server{
		addr:           addr,
		token:          token,
		prefork:        prefork,
		app:            app,
		inboundService: inboundService,
		nodeService:    nodeService,
	}
	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	s.app.Use(s.authMiddleware)

	s.app.Get(healthcheck.LivenessEndpoint, healthcheck.New())
	s.app.Get(healthcheck.ReadinessEndpoint, healthcheck.New())

	api := s.app.Group("/api")
	v1 := api.Group("/v1")

	nodes := v1.Group("/nodes")
	nodes.Get("", s.GetNodes)
	nodes.Get("/active", s.GetActiveNodes)
	nodes.Get("/:nodeID", s.GetNode)
	nodes.Post("", s.CreateNode)
	nodes.Delete("/:nodeID", s.DeleteNode)
	nodes.Patch("/:nodeID/status", s.UpdateNodeStatus)

	inbounds := nodes.Group("/:nodeID/inbounds")
	inbounds.Get("", s.GetInbounds)
	inbounds.Post("", s.CreateInbound)
	inbounds.Get("/runtime/count", s.GetRuntimeInboundsCount)
	inbounds.Get("/:tag", s.GetInbound)
	inbounds.Delete("/:tag", s.DeleteInbound)
	inbounds.Patch("/:tag/renew", s.InboundRenew)

	inboundUsers := inbounds.Group("/:tag/users")
	inboundUsers.Get("", s.GetInboundUsers)
	inboundUsers.Post("", s.CreateUser)
	inboundUsers.Delete("/:email", s.DeleteUser)
	inboundUsers.Patch("/:email/renew", s.InboundUserRenew)
}

func (s *Server) StartTLS(certFilePath, keyFilePath string) error {
	return s.app.Listen(s.addr, fiber.ListenConfig{
		DisableStartupMessage: true,
		CertFile:              certFilePath,
		CertKeyFile:           keyFilePath,
		EnablePrefork:         s.prefork,
	})
}

func (s *Server) Start() error {
	return s.app.Listen(s.addr, fiber.ListenConfig{
		DisableStartupMessage: true,
		EnablePrefork:         s.prefork,
	})
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.app.ShutdownWithContext(ctx)
}

func (s *Server) authMiddleware(c fiber.Ctx) error {
	h := c.Get("Authorization")
	if h == "" {
		return fiber.ErrUnauthorized
	}

	if err := auth.VerifyHMAC(h, s.token); err != nil {
		return fiber.ErrUnauthorized
	}
	return c.Next()
}

func (s *Server) requiredParams(c fiber.Ctx, keys ...string) (map[string]string, error) {
	m := make(map[string]string)
	for _, k := range keys {
		v := c.Params(k)
		if v == "" {
			return nil, fmt.Errorf("%s parameter is required", k)
		}
		m[k] = v
	}
	return m, nil
}
