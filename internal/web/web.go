package web

import (
	"github.com/corani/unwise/internal/config"
	"github.com/corani/unwise/internal/storage"
	"github.com/gofiber/fiber/v3"
)

type Server struct {
	conf *config.Config
	stor storage.Storage
}

func New(conf *config.Config, stor storage.Storage) *Server {
	return &Server{
		conf: conf,
		stor: stor,
	}
}

func (s *Server) App() *fiber.App {
	return newApp(s)
}
