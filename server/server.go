package server

import (
	"auth-service-api/storage/postgres"

	"github.com/gofiber/fiber/v2"
)

type Server struct {
	app     *fiber.App
	storage *postgres.Storage
}

func New(storage *postgres.Storage) *Server {
	app := fiber.New()

	return &Server{
		app:     app,
		storage: storage,
	}
}

func (s *Server) Start() error {
	apiv1 := s.app.Group("/api/v1")

	apiv1.Get("/tokens", s.handleGetToken)
	apiv1.Post("/refresh", s.handleRefresh)

	return s.app.Listen(":8080")
	// Gracefull shutdown
}
