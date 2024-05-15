package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/corani/unwise/internal/config"
	"github.com/gofiber/fiber/v2"
)

type Server struct {
	conf *config.Config
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Details string `json:"details,omitempty"`
}

func newServer() *Server {
	return &Server{
		conf: config.MustLoad(),
	}
}

func (s *Server) HandleRoot(c *fiber.Ctx) error {
	// NOTE(daniel): unauthenticated endpoint.

	return c.SendStatus(http.StatusNoContent)
}

func (s *Server) HandleAuth(c *fiber.Ctx) error {
	return c.SendStatus(http.StatusNoContent)
}

func (s *Server) HandleCreateHighlights(c *fiber.Ctx) error {
	var req CreateHighlightRequest

	if err := c.BodyParser(&req); err != nil {
		return fmt.Errorf("%w: %v (raw=%q)", fiber.ErrBadRequest, err, string(c.Body()))
	}

	s.conf.Logger.Info("create highlights",
		"raw", string(c.Body()),
		"req", req)

	var res []CreateHighlightResponse

	return c.JSON(res)
}

func (s *Server) HandleListHighlights(c *fiber.Ctx) error {
	p, err := parseParams(c)
	if err != nil {
		return err
	}

	s.conf.Logger.Info("list highlights",
		"page_size", p.pageSize,
		"updated__lt", p.updatedLT,
		"updated__gt", p.updatedGT)

	var res ListHighlightsResponse

	return c.JSON(res)
}

func (s *Server) HandleListBooks(c *fiber.Ctx) error {
	p, err := parseParams(c)
	if err != nil {
		return err
	}

	s.conf.Logger.Info("list books",
		"page_size", p.pageSize,
		"updated__lt", p.updatedLT,
		"updated__gt", p.updatedGT)

	var res ListBooksResponse

	return c.JSON(res)
}

func (s *Server) CheckAuth(c *fiber.Ctx, key string) (bool, error) {
	return key == s.conf.Token, nil
}

func (s *Server) HandleError(c *fiber.Ctx, err error) error {
	var e *fiber.Error

	if errors.As(err, &e) {
		return c.
			Status(e.Code).
			JSON(&ErrorResponse{
				Error:   e.Message,
				Code:    e.Code,
				Details: err.Error(),
			})
	}

	return c.
		Status(http.StatusInternalServerError).
		JSON(&ErrorResponse{
			Error: err.Error(),
			Code:  http.StatusInternalServerError,
		})
}
