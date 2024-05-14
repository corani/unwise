package main

import (
	"errors"
	"net/http"

	"github.com/corani/unwise/internal/config"
	"github.com/gofiber/fiber/v2"
)

type Server struct {
	conf *config.Config
}

type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

func newServer() *Server {
	return &Server{
		conf: config.MustLoad(),
	}
}

func (s *Server) CheckAuth(c *fiber.Ctx, key string) (bool, error) {
	return key == s.conf.Token, nil
}

func (s *Server) HandleAuth(c *fiber.Ctx) error {
	return c.SendStatus(http.StatusNoContent)
}

func (s *Server) HandleError(c *fiber.Ctx, err error) error {
	var e *fiber.Error

	code := http.StatusInternalServerError
	msg := err.Error()

	if errors.As(err, &e) {
		code = e.Code
		msg = e.Message
	}

	return c.
		Status(code).
		JSON(&ErrorResponse{
			Error: msg,
			Code:  code,
		})
}

func (s *Server) HandleCreateHighlights(c *fiber.Ctx) error {
	// ensure content type is json
	if c.Get("Content-Type") != fiber.MIMEApplicationJSON {
		return fiber.ErrUnsupportedMediaType
	}

	var req CreateHighlightRequest

	if err := c.BodyParser(&req); err != nil {
		return err
	}

	var res []CreateHighlightResponse

	return c.JSON(res)
}

func (s *Server) HandleListHighlights(c *fiber.Ctx) error {
	// optional, default 100, max 1000
	pageSize := c.QueryInt("page_size", 100)
	if pageSize < 0 || pageSize > 1000 {
		return fiber.ErrBadRequest
	}

	// optional, filter by last updated datetime (less than)
	updatedLT, err := parseISO8601Datetime(c.Query("updated__lt"))
	if err != nil {
		return err
	}

	// optional, filter by last updated datetime (greater than)
	updatedGT, err := parseISO8601Datetime(c.Query("updated__gt"))
	if err != nil {
		return err
	}

	_ = updatedLT
	_ = updatedGT

	var res ListHighlightsResponse

	return c.JSON(res)
}

func (s *Server) HandleListBooks(c *fiber.Ctx) error {
	// optional, default 100, max 1000
	pageSize := c.QueryInt("page_size", 100)
	if pageSize < 0 || pageSize > 1000 {
		return fiber.ErrBadRequest
	}

	// optional, filter by last updated datetime (less than)
	updatedLT, err := parseISO8601Datetime(c.Query("updated__lt"))
	if err != nil {
		return err
	}

	// optional, filter by last updated datetime (greater than)
	updatedGT, err := parseISO8601Datetime(c.Query("updated__gt"))
	if err != nil {
		return err
	}

	_ = updatedLT
	_ = updatedGT

	var res ListBooksResponse

	return c.JSON(res)
}
