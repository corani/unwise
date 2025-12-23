package web

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	fake "github.com/corani/unwise/fakes/storage"
	"github.com/corani/unwise/internal/config"
	"github.com/corani/unwise/internal/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

			s := New(config.MustLoad(), nil)
			s.conf.Token = tc.token

			act, err := s.CheckAuth(nil, tc.key)
			rq.NoError(err)
			rq.Equal(tc.exp, act)
		})
	}
}

func TestServer_HandleRoot(t *testing.T) {
	rq := require.New(t)
	s := New(config.MustLoad(), nil)

	app := fiber.New()
	app.Get("/", s.HandleRoot)

	resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/", nil))
	rq.NoError(err)
	rq.Equal(http.StatusNoContent, resp.StatusCode)
}

func TestServer_HandleAuth(t *testing.T) {
	rq := require.New(t)
	s := New(config.MustLoad(), nil)

	app := fiber.New()
	app.Get("/auth", s.HandleAuth)

	resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/auth", nil))
	rq.NoError(err)
	rq.Equal(http.StatusNoContent, resp.StatusCode)
}

func TestServer_HandleError(t *testing.T) {
	tt := []struct {
		name     string
		endpoint string
		expCode  int
		expBody  string
	}{
		{
			name:     "fiber error",
			endpoint: "/error",
			expCode:  http.StatusNotFound,
			expBody: `{
				"error":"Cannot GET /error",
				"code":404,
				"details":"Cannot GET /error"
			}`,
		},
		{
			name:     "custom error",
			endpoint: "/custom",
			expCode:  http.StatusInternalServerError,
			expBody: `{
				"error":"assert.AnError general error for testing",
				"code":500
			}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			rq := require.New(t)
			s := New(config.MustLoad(), nil)

			app := fiber.New(fiber.Config{
				ErrorHandler: s.HandleError,
			})
			app.Get("/custom", func(c *fiber.Ctx) error {
				return assert.AnError
			})

			resp, err := app.Test(httptest.NewRequest(http.MethodGet, tc.endpoint, nil))
			rq.NoError(err)

			rq.Equal(tc.expCode, resp.StatusCode)
			bodyJSONEq(t, tc.expBody, resp.Body)
		})
	}
}

func TestServer_HandleCreateHighlights(t *testing.T) {
	tt := []struct {
		name        string
		content     string
		contentType string
		setup       func(*fake.Storage)
		expCode     int
		expBody     string
	}{
		{
			name:        "invalid content type",
			content:     "",
			contentType: "",
			expCode:     http.StatusBadRequest,
			expBody: `{
				"error":"Bad Request",
				"code":400,
				"details":"Bad Request: Unprocessable Entity (raw=\"\")"
			}`,
		},
		{
			name:        "invalid body",
			content:     "invalid",
			contentType: fiber.MIMEApplicationJSON,
			expCode:     http.StatusBadRequest,
			expBody: `{
				"error":"Bad Request",
				"code":400,
				"details":"Bad Request: invalid character 'i' looking for beginning of value (raw=\"invalid\")"
			}`,
		},
		{
			name: "valid request",
			content: `{
				"highlights": [
					{"title": "title1", "text": "text1"},
					{"title": "title1", "text": "text2"}
				]
			}`,
			contentType: fiber.MIMEApplicationJSON,
			setup: func(stor *fake.Storage) {
				stor.EXPECT().AddBook(mock.Anything, "title1", "", "").
					Return(storage.Book{ID: 1}, nil)
				stor.EXPECT().AddHighlight(mock.Anything, storage.Book{ID: 1}, "text1", "", "", 0, "").
					Return(storage.Highlight{ID: 1}, nil)
				stor.EXPECT().AddHighlight(mock.Anything, storage.Book{ID: 1}, "text2", "", "", 0, "").
					Return(storage.Highlight{ID: 2}, nil)
			},
			expCode: http.StatusOK,
			expBody: `[
				{
					"id": 1,
					"author": "",
					"title": "",
					"category": "books",
					"last_highlight_at": "0001-01-01T00:00:00Z",
					"modified_highlights": [1,2],
					"num_highlights": 2,
					"source_url": "",
					"updated": "0001-01-01T00:00:00Z"
				}
			]`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			rq := require.New(t)

			stor := fake.NewStorage(t)
			if tc.setup != nil {
				tc.setup(stor)
			}

			conf := config.MustLoad()
			serv := New(conf, stor)

			app := fiber.New(fiber.Config{
				ErrorHandler: serv.HandleError,
			})
			app.Post("/highlights", serv.HandleCreateHighlights)

			req := httptest.NewRequest(http.MethodPost, "/highlights", strings.NewReader(tc.content))

			if tc.contentType != "" {
				req.Header.Set("Content-Type", tc.contentType)
			}

			resp, err := app.Test(req)
			rq.NoError(err)

			rq.Equal(tc.expCode, resp.StatusCode)
			bodyJSONEq(t, tc.expBody, resp.Body)
		})
	}
}

func TestServer_HandleListHighlights(t *testing.T) {
	tt := []struct {
		name     string
		endpoint string
		setup    func(*fake.Storage)
		expCode  int
		expBody  string
	}{
		{
			name:     "invalid page size",
			endpoint: "/highlights?page_size=-1",
			expCode:  http.StatusBadRequest,
			expBody: `{
				"error":"Bad Request",
				"code":400,
				"details":"Bad Request: invalid page_size -1"
			}`,
		},
		{
			name:     "invalid updated__lt",
			endpoint: "/highlights?updated__lt=invalid",
			expCode:  http.StatusBadRequest,
			expBody: `{
				"error":"Bad Request",
				"code":400,
				"details":"Bad Request: invalid updated__lt \"invalid\""
			}`,
		},
		{
			name:     "invalid updated__gt",
			endpoint: "/highlights?updated__gt=invalid",
			expCode:  http.StatusBadRequest,
			expBody: `{
				"error":"Bad Request",
				"code":400,
				"details":"Bad Request: invalid updated__gt \"invalid\""
			}`,
		},
		{
			name:     "valid request",
			endpoint: "/highlights",
			setup: func(stor *fake.Storage) {
				stor.EXPECT().ListHighlights(mock.Anything, time.Time{}, time.Time{}).Return([]storage.Highlight{
					{ID: 1},
				}, nil)
			},
			expCode: http.StatusOK,
			expBody: `{"results":[
				{
					"id": 1,
					"book_id": 0,
					"chapter": "",
					"location": 0,
					"text": "",
					"note": "",
					"url": "",
					"updated": "0001-01-01T00:00:00Z"
				}
			]}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			rq := require.New(t)

			stor := fake.NewStorage(t)
			if tc.setup != nil {
				tc.setup(stor)
			}

			conf := config.MustLoad()
			serv := New(conf, stor)

			app := fiber.New(fiber.Config{
				ErrorHandler: serv.HandleError,
			})
			app.Get("/highlights", serv.HandleListHighlights)

			req := httptest.NewRequest(http.MethodGet, tc.endpoint, nil)

			resp, err := app.Test(req)
			rq.NoError(err)

			rq.Equal(tc.expCode, resp.StatusCode)
			bodyJSONEq(t, tc.expBody, resp.Body)
		})
	}
}

func TestServer_HandleListBooks(t *testing.T) {
	tt := []struct {
		name     string
		endpoint string
		setup    func(*fake.Storage)
		expCode  int
		expBody  string
	}{
		{
			name:     "invalid page size",
			endpoint: "/books?page_size=-1",
			expCode:  http.StatusBadRequest,
			expBody: `{
				"error":"Bad Request",
				"code":400,
				"details":"Bad Request: invalid page_size -1"
			}`,
		},
		{
			name:     "invalid updated__lt",
			endpoint: "/books?updated__lt=invalid",
			expCode:  http.StatusBadRequest,
			expBody: `{
				"error":"Bad Request",
				"code":400,
				"details":"Bad Request: invalid updated__lt \"invalid\""
			}`,
		},
		{
			name:     "invalid updated__gt",
			endpoint: "/books?updated__gt=invalid",
			expCode:  http.StatusBadRequest,
			expBody: `{
				"error":"Bad Request",
				"code":400,
				"details":"Bad Request: invalid updated__gt \"invalid\""
			}`,
		},
		{
			name:     "valid request",
			endpoint: "/books",
			setup: func(stor *fake.Storage) {
				stor.EXPECT().ListBooks(mock.Anything, time.Time{}, time.Time{}).Return([]storage.Book{
					{ID: 1},
				}, nil)
			},
			expCode: http.StatusOK,
			expBody: `{"results":[
				{
					"id": 1,
					"category": "books",
					"author": "",
					"title": "",
					"num_highlights": 0,
					"source_url": "",
					"updated": "0001-01-01T00:00:00Z"
				}
			]}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			rq := require.New(t)

			stor := fake.NewStorage(t)
			if tc.setup != nil {
				tc.setup(stor)
			}

			conf := config.MustLoad()
			serv := New(conf, stor)

			app := fiber.New(fiber.Config{
				ErrorHandler: serv.HandleError,
			})
			app.Get("/books", serv.HandleListBooks)

			req := httptest.NewRequest(http.MethodGet, tc.endpoint, nil)

			resp, err := app.Test(req)
			rq.NoError(err)

			rq.Equal(tc.expCode, resp.StatusCode)
			bodyJSONEq(t, tc.expBody, resp.Body)
		})
	}
}

func bodyJSONEq(t *testing.T, exp string, act io.ReadCloser) {
	t.Helper()

	rq := require.New(t)

	bs, err := io.ReadAll(act)
	defer act.Close()

	rq.NoError(err)
	rq.JSONEq(exp, string(bs))
}
