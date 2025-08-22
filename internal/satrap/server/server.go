package server

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
	"github.com/vayzur/apadana/internal/auth"
	xray "github.com/vayzur/apadana/pkg/satrap/xray/client"
)

type Server struct {
	addr       string
	token      string
	prefork    bool
	app        *fiber.App
	xrayClient *xray.Client
}

func NewServer(addr, token string, prefork bool, xrayClient *xray.Client) *Server {
	app := fiber.New(fiber.Config{
		CaseSensitive: true,
		StrictRouting: true,
	})
	s := &Server{
		addr:       addr,
		token:      token,
		prefork:    prefork,
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
	inbounds.Get("/count", s.InboundsCount)
	inbounds.Post("", s.requireJSON, s.AddInbound)
	inbounds.Delete("/:tag", s.RemoveInbound)

	users := inbounds.Group("/:tag/users")
	users.Post("", s.AddUser)
	users.Delete(":email", s.RemoveUser)
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

func (s *Server) requireJSON(c fiber.Ctx) error {
	ct := c.Get(fiber.HeaderContentType)
	if ct != fiber.MIMEApplicationJSON {
		return c.Status(fiber.StatusUnsupportedMediaType).
			JSON(fiber.Map{"error": "Content-Type must be application/json"})
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
