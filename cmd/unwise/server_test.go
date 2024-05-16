package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/corani/unwise/internal/config"
	"github.com/corani/unwise/internal/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServer_CheckAuth(t *testing.T) {
	tt := []struct {
		name  string
		key   string
		token string
		exp   bool
	}{
		{
			name:  "valid token",
			key:   "secret",
			token: "secret",
			exp:   true,
		},
		{
			name:  "invalid token",
			key:   "secret",
			token: "invalid",
			exp:   false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			rq := require.New(t)

			s := newServer(config.MustLoad(), nil)
			s.conf.Token = tc.token

			act, err := s.CheckAuth(nil, tc.key)
			rq.NoError(err)
			rq.Equal(tc.exp, act)
		})
	}
}

func TestServer_HandleRoot(t *testing.T) {
	rq := require.New(t)
	s := newServer(config.MustLoad(), nil)

	app := fiber.New()
	app.Get("/", s.HandleRoot)

	resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/", nil))
	rq.NoError(err)
	rq.Equal(http.StatusNoContent, resp.StatusCode)
}

func TestServer_HandleAuth(t *testing.T) {
	rq := require.New(t)
	s := newServer(config.MustLoad(), nil)

	app := fiber.New()
	app.Get("/auth", s.HandleAuth)

	resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/auth", nil))
	rq.NoError(err)
	rq.Equal(http.StatusNoContent, resp.StatusCode)
}

func TestServer_HandleError(t *testing.T) {
	rq := require.New(t)
	s := newServer(config.MustLoad(), nil)

	app := fiber.New(fiber.Config{
		ErrorHandler: s.HandleError,
	})
	app.Get("/custom", func(c *fiber.Ctx) error {
		return assert.AnError
	})

	t.Run("fiber error", func(t *testing.T) {
		resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/error", nil))
		rq.NoError(err)

		rq.Equal(http.StatusNotFound, resp.StatusCode)
		bodyJSONEq(t, `{
			"error":"Cannot GET /error",
			"code":404,
			"details":"Cannot GET /error"
		}`, resp.Body)
	})

	t.Run("custom error", func(t *testing.T) {
		resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/custom", nil))
		rq.NoError(err)

		rq.Equal(http.StatusInternalServerError, resp.StatusCode)
		bodyJSONEq(t, `{
			"error":"assert.AnError general error for testing",
			"code":500
		}`, resp.Body)
	})
}

func TestServer_HandleCreateHighlights(t *testing.T) {
	rq := require.New(t)
	conf := config.MustLoad()
	stor := storage.New(conf)
	serv := newServer(conf, stor)

	app := fiber.New(fiber.Config{
		ErrorHandler: serv.HandleError,
	})
	app.Post("/highlights", serv.HandleCreateHighlights)

	t.Run("invalid content type", func(t *testing.T) {
		resp, err := app.Test(httptest.NewRequest(http.MethodPost, "/highlights", nil))
		rq.NoError(err)
		rq.Equal(http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("invalid body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/highlights", strings.NewReader(`invalid`))
		req.Header.Set("Content-Type", fiber.MIMEApplicationJSON)

		resp, err := app.Test(req)
		rq.NoError(err)

		rq.Equal(http.StatusBadRequest, resp.StatusCode)
		bodyJSONEq(t, `{
			"error":"Bad Request",
			"code":400,
			"details":"Bad Request: invalid character 'i' looking for beginning of value (raw=\"invalid\")"
		}`, resp.Body)
	})

	t.Run("valid content type", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/highlights", strings.NewReader(`{}`))
		req.Header.Set("Content-Type", fiber.MIMEApplicationJSON)

		resp, err := app.Test(req)
		rq.NoError(err)

		rq.Equal(http.StatusOK, resp.StatusCode)
		bodyJSONEq(t, `null`, resp.Body)
	})
}

func TestServer_HandleListHighlights(t *testing.T) {
	rq := require.New(t)
	conf := config.MustLoad()
	stor := storage.New(conf)
	serv := newServer(conf, stor)

	app := fiber.New(fiber.Config{
		ErrorHandler: serv.HandleError,
	})
	app.Get("/highlights", serv.HandleListHighlights)

	t.Run("invalid page size", func(t *testing.T) {
		resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/highlights?page_size=-1", nil))
		rq.NoError(err)
		rq.Equal(http.StatusBadRequest, resp.StatusCode)

		bodyJSONEq(t, `{
			"error":"Bad Request",
			"code":400,
			"details":"Bad Request: invalid page_size -1"
		}`, resp.Body)
	})

	t.Run("invalid updated__lt", func(t *testing.T) {
		resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/highlights?updated__lt=invalid", nil))
		rq.NoError(err)
		rq.Equal(http.StatusBadRequest, resp.StatusCode)

		bodyJSONEq(t, `{
			"error":"Bad Request",
			"code":400,
			"details":"Bad Request: invalid updated__lt \"invalid\""
		}`, resp.Body)
	})

	t.Run("invalid updated__gt", func(t *testing.T) {
		resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/highlights?updated__gt=invalid", nil))
		rq.NoError(err)
		rq.Equal(http.StatusBadRequest, resp.StatusCode)

		bodyJSONEq(t, `{
			"error":"Bad Request",
			"code":400,
			"details":"Bad Request: invalid updated__gt \"invalid\""
		}`, resp.Body)
	})

	t.Run("valid body", func(t *testing.T) {
		resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/highlights", nil))
		rq.NoError(err)

		rq.Equal(http.StatusOK, resp.StatusCode)
		bodyJSONEq(t, `{}`, resp.Body)
	})
}

func TestServer_HandleListBooks(t *testing.T) {
	rq := require.New(t)
	conf := config.MustLoad()
	stor := storage.New(conf)
	serv := newServer(conf, stor)

	app := fiber.New(fiber.Config{
		ErrorHandler: serv.HandleError,
	})
	app.Get("/books", serv.HandleListBooks)

	t.Run("invalid page size", func(t *testing.T) {
		resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/books?page_size=-1", nil))
		rq.NoError(err)
		rq.Equal(http.StatusBadRequest, resp.StatusCode)

		bodyJSONEq(t, `{
			"error":"Bad Request",
			"code":400,
			"details": "Bad Request: invalid page_size -1"
		}`, resp.Body)
	})

	t.Run("invalid updated__lt", func(t *testing.T) {
		resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/books?updated__lt=invalid", nil))
		rq.NoError(err)
		rq.Equal(http.StatusBadRequest, resp.StatusCode)

		bodyJSONEq(t, `{
			"error":"Bad Request",
			"code":400,
			"details":"Bad Request: invalid updated__lt \"invalid\""
		}`, resp.Body)
	})

	t.Run("invalid updated__gt", func(t *testing.T) {
		resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/books?updated__gt=invalid", nil))
		rq.NoError(err)
		rq.Equal(http.StatusBadRequest, resp.StatusCode)

		bodyJSONEq(t, `{
			"error":"Bad Request",
			"code":400,
			"details":"Bad Request: invalid updated__gt \"invalid\""
		}`, resp.Body)
	})

	t.Run("valid body", func(t *testing.T) {
		resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/books", nil))
		rq.NoError(err)

		rq.Equal(http.StatusOK, resp.StatusCode)
		bodyJSONEq(t, `{}`, resp.Body)
	})
}

func bodyJSONEq(t *testing.T, exp string, act io.ReadCloser) {
	t.Helper()

	rq := require.New(t)

	bs, err := io.ReadAll(act)
	defer act.Close()

	rq.NoError(err)
	rq.JSONEq(exp, string(bs))
}
