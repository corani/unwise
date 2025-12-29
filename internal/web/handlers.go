package web

import (
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/corani/unwise/internal/storage"
	"github.com/corani/unwise/static"
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

	res := ListHighlightsResponse{
		Results: make([]ListHighlight, 0),
	}

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

	res := ListBooksResponse{
		Results: make([]ListBook, 0),
	}

	bs, err := s.stor.ListBooks(ctx, p.updatedLT, p.updatedGT)
	if err != nil {
		return err
	}

	for _, book := range bs {
		// NOTE(daniel): avoid "object null is not iterable" error in the Obsidian plugin.
		if book.NumHighlights == 0 {
			continue
		}

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

func (s *Server) CheckToken(c *fiber.Ctx, key string) (bool, error) {
	return key == s.conf.Token, nil
}

func (s *Server) CheckAuth(user, pass string) bool {
	success := user == s.conf.User && pass == s.conf.Token

	if !success {
		s.conf.Logger.Warn("failed basic auth attempt",
			"user", user)
	}

	return success
}

func (s *Server) HandleUIIndex(c *fiber.Ctx) error {
	data, err := static.FS.ReadFile("index.html")
	if err != nil {
		return err
	}

	c.Set("Content-Type", "text/html")
	return c.Send(data)
}

func (s *Server) HandleUIListBooks(c *fiber.Ctx) error {
	ctx := c.Context()

	res := ListBooksResponse{
		Results: make([]ListBook, 0),
	}

	// For UI, we want all books, so use zero times for filtering
	bs, err := s.stor.ListBooks(ctx, time.Time{}, time.Time{})
	if err != nil {
		return err
	}

	// Sort the books by Updated time descending
	sort.Slice(bs, func(i, j int) bool {
		return bs[i].Updated.After(bs[j].Updated)
	})

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

func (s *Server) HandleUIListHighlights(c *fiber.Ctx) error {
	ctx := c.Context()

	bookID, err := c.ParamsInt("id")
	if err != nil {
		return fmt.Errorf("%w: invalid book ID", fiber.ErrBadRequest)
	}

	res := ListHighlightsResponse{
		Results: make([]ListHighlight, 0),
	}

	hs, err := s.stor.ListHighlightsByBook(ctx, bookID)
	if err != nil {
		return err
	}

	// Sort the highlights by Location ascending
	sort.Slice(hs, func(i, j int) bool {
		return hs[i].Location < hs[j].Location
	})

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

func (s *Server) HandleUIUpdateHighlight(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get ID from URL
	id, err := c.ParamsInt("id")
	if err != nil {
		return fmt.Errorf("%w: invalid highlight ID", fiber.ErrBadRequest)
	}

	// Parse highlight from request body
	var req ListHighlight
	if err := c.BodyParser(&req); err != nil {
		return fmt.Errorf("%w: %v", fiber.ErrBadRequest, err)
	}

	// Validate required fields
	if req.Text == "" {
		return fmt.Errorf("%w: text is required", fiber.ErrBadRequest)
	}

	// Create full Highlight struct - storage will ignore unused fields
	h := storage.Highlight{
		ID:       id,
		BookID:   req.BookID,
		Text:     req.Text,
		Note:     req.Note,
		Chapter:  req.Chapter,
		Location: req.Location,
		URL:      req.URL,
	}

	updated, err := s.stor.UpdateHighlight(ctx, h)
	if err != nil {
		return err
	}

	return c.JSON(ListHighlight{
		ID:        updated.ID,
		BookID:    updated.BookID,
		Text:      updated.Text,
		Note:      updated.Note,
		Chapter:   updated.Chapter,
		Location:  updated.Location,
		URL:       updated.URL,
		UpdatedAt: updated.Updated.Format(time.RFC3339),
	})
}

func (s *Server) HandleUIDeleteHighlight(c *fiber.Ctx) error {
	ctx := c.Context()

	id, err := c.ParamsInt("id")
	if err != nil {
		return fmt.Errorf("%w: invalid highlight ID", fiber.ErrBadRequest)
	}

	if err := s.stor.DeleteHighlight(ctx, id); err != nil {
		return err
	}

	return c.SendStatus(http.StatusNoContent)
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
