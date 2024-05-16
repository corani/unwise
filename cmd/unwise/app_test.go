package main

import (
	"testing"

	"github.com/corani/unwise/internal/config"
	"github.com/stretchr/testify/require"
)

func TestMain_newApp(t *testing.T) {
	rq := require.New(t)

	server := newServer(config.MustLoad(), nil)
	app := newApp(server)

	rq.NotNil(app)

	routes := make(map[string]struct{})

	for _, route := range app.GetRoutes(true) {
		routes[route.Path] = struct{}{}
	}

	rq.Contains(routes, "/api/v2/auth")
	rq.Contains(routes, "/api/v2/highlights")
	rq.Contains(routes, "/api/v2/books")
}
