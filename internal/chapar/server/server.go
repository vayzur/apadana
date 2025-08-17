package server

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
	"github.com/vayzur/apadana/internal/auth"
	xray "github.com/vayzur/apadana/pkg/chapar/xray/client"
)

type Server struct {
	addr       string
	token      string
	prefork    bool
	app        *fiber.App
	xrayClient *xray.Client
}

func NewServer(addr, token string, xrayClient *xray.Client) *Server {
	app := fiber.New(fiber.Config{
		CaseSensitive: true,
		StrictRouting: true,
	})
	s := &Server{
		addr:       addr,
		token:      token,
		app:        app,
		xrayClient: xrayClient,
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

	inbounds := v1.Group("/inbounds")
	inbounds.Post("", requireJSON, s.AddInbound)
	inbounds.Delete("/:tag", s.RemoveInbound)
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

func requireJSON(c fiber.Ctx) error {
	ct := c.Get(fiber.HeaderContentType)
	if ct != fiber.MIMEApplicationJSON {
		return c.Status(fiber.StatusUnsupportedMediaType).
			JSON(fiber.Map{"error": "Content-Type must be application/json"})
	}
	return c.Next()
}
