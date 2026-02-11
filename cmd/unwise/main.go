package main

import (
	"context"

	"github.com/corani/unwise/internal/config"
	"github.com/corani/unwise/internal/storage/sqlite"
	"github.com/corani/unwise/internal/web"
	"github.com/gofiber/fiber/v3"
)

func main() {
	conf := config.MustLoad()

	stor, err := sqlite.New(context.Background(), conf)
	if err != nil {
		conf.Logger.Errorf("sqlite: %v", err)
	}

	listenConfig := fiber.ListenConfig{
		EnablePrintRoutes: false,
	}

	serv := web.New(conf, stor)
	app := serv.App()

	conf.PrintBanner()

	// default RestAddr=":3123"
	if err := app.Listen(conf.RestAddr, listenConfig); err != nil {
		conf.Logger.Errorf("listen: %v", err)
	}
}
