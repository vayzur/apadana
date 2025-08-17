package server

import (
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
	nodeService    *service.NodeSerivce
}

func NewServer(addr, token string, prefork bool, inboundService *service.InboundService, nodeSerivce *service.NodeSerivce) *Server {
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
		nodeService:    nodeSerivce,
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
	nodes.Get("", s.GetAllNodes)
	nodes.Get("/active", s.GetActiveNodes)
	nodes.Get("/:nodeID", s.GetNode)
	nodes.Post("", s.CreateNode)
	nodes.Delete("/:nodeID", s.DeleteNode)
	nodes.Patch("/:nodeID/status", s.UpdateNodeStatus)

	inbounds := nodes.Group("/:nodeID/inbounds")
	inbounds.Get("", s.GetAllInbounds)
	inbounds.Get("/:tag", s.GetInbound)
	inbounds.Post("", s.CreateInbound)
	inbounds.Delete("/:tag", s.DeleteInbound)
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

func (s *Server) Stop() error {
	return s.app.Shutdown()
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
