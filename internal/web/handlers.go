package web

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) HandleRoot(c *fiber.Ctx) error {
	// NOTE(daniel): unauthenticated endpoint.

	return c.SendStatus(http.StatusNoContent)
}

func (s *Server) HandleAuth(c *fiber.Ctx) error {
	return c.SendStatus(http.StatusNoContent)
}

func (s *Server) HandleCreateHighlights(c *fiber.Ctx) error {
	ctx := c.Context()

	var req CreateHighlightRequest

	if err := c.BodyParser(&req); err != nil {
		return fmt.Errorf("%w: %v (raw=%q)", fiber.ErrBadRequest, err, string(c.Body()))
	}

	r := strings.NewReplacer("\t", "", "\n", "")

	s.conf.Logger.Info("create highlights",
		"raw", r.Replace(string(c.Body())),
		"content-type", c.Get("Content-Type"),
		"req", req)

	var list []CreateHighlightResponse

	// TODO(daniel): NumHighlights and LastHighlightAt may not be correct.
	for _, rh := range req.Highlights {
		b, err := s.stor.AddBook(ctx, rh.Title, rh.Author, rh.SourceURL)
		if err != nil {
			return err
		}

		h, err := s.stor.AddHighlight(ctx, b, rh.Text, rh.Note, rh.Chapter, rh.Location, rh.HighlightURL)
		if err != nil {
			return err
		}

		found := false

		for i, v := range list {
			if v.ID == b.ID {
				list[i].ModifiedHighlights = append(v.ModifiedHighlights, h.ID)
				list[i].NumHighlights++
				list[i].LastHighlightAt = h.Updated.Format(time.RFC3339)
				found = true

				break
			}
		}

		if !found {
			list = append(list, CreateHighlightResponse{
				ID:                 b.ID,
				Title:              b.Title,
				Author:             b.Author,
				SourceURL:          b.SourceURL,
				Category:           HighlightCategoryBooks,
				NumHighlights:      b.NumHighlights + 1,
				LastHighlightAt:    h.Updated.Format(time.RFC3339),
				UpdatedAt:          h.Updated.Format(time.RFC3339),
				ModifiedHighlights: []int{h.ID},
			})
		}
	}

	return c.JSON(list)
}

func (s *Server) HandleListHighlights(c *fiber.Ctx) error {
	ctx := c.Context()

	p, err := parseParams(c)
	if err != nil {
		return err
	}

	s.conf.Logger.Info("list highlights",
		"page_size", p.pageSize,
		"updated__lt", p.updatedLT,
		"updated__gt", p.updatedGT)

	var res ListHighlightsResponse

	hs, err := s.stor.ListHighlights(ctx, p.updatedLT, p.updatedGT)
	if err != nil {
		return err
	}

	for _, highlight := range hs {
		res.Results = append(res.Results, ListHighlight{
			ID:        highlight.ID,
			BookID:    highlight.BookID,
			Text:      highlight.Text,
			Note:      highlight.Note,
			Chapter:   highlight.Chapter,
			Location:  highlight.Location,
			URL:       highlight.URL,
			UpdatedAt: highlight.Updated.Format(time.RFC3339),
		})
	}

	return c.JSON(res)
}

func (s *Server) HandleListBooks(c *fiber.Ctx) error {
	ctx := c.Context()

	p, err := parseParams(c)
	if err != nil {
		return err
	}

	s.conf.Logger.Info("list books",
		"page_size", p.pageSize,
		"updated__lt", p.updatedLT,
		"updated__gt", p.updatedGT)

	var res ListBooksResponse

	bs, err := s.stor.ListBooks(ctx, p.updatedLT, p.updatedGT)
	if err != nil {
		return err
	}

	for _, book := range bs {
		res.Results = append(res.Results, ListBook{
			ID:            book.ID,
			Title:         book.Title,
			Author:        book.Author,
			SourceURL:     book.SourceURL,
			Category:      HighlightCategoryBooks,
			NumHighlights: book.NumHighlights,
			UpdatedAt:     book.Updated.Format(time.RFC3339),
		})
	}

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
