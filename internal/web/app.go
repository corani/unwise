package web

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/keyauth"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func newApp(server *Server) *fiber.App {
	app := fiber.New(fiber.Config{
		EnablePrintRoutes: false,
		StreamRequestBody: true,
		StrictRouting:     false, // ignore trailing slashes
		ErrorHandler:      server.HandleError,
	})
	app.Use(
		logger.New(), // log each request
		helmet.New(), // secure headers
	)

	app.Get("/", server.HandleRoot)

	// UI
	ui := app.Group("/ui")
	ui.Use(basicauth.New(basicauth.Config{
		Users:      map[string]string{},
		Authorizer: server.CheckAuth,
	}))
	ui.Get("/", server.HandleUIIndex)
	ui.Get("/api/books", server.HandleUIListBooks)
	ui.Get("/api/books/:id/highlights", server.HandleUIListHighlights)
	ui.Static("/static", "./static")

	// API
	// default RestPath="/api/v2"
	api := app.Group(server.conf.RestPath)
	api.Use(keyauth.New(keyauth.Config{
		AuthScheme: "Token",
		Validator:  server.CheckToken,
	}))

	// check if token is valid
	api.Get("/auth", server.HandleAuth)

	// used by Moon+ Reader to create highlights
	api.Post("/highlights", server.HandleCreateHighlights)

	// Used by Obsidian
	api.Get("/highlights", server.HandleListHighlights)
	api.Get("/books", server.HandleListBooks)

	return app
}
