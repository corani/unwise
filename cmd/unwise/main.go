package main

import (
	"context"

	"github.com/corani/unwise/internal/config"
	"github.com/corani/unwise/internal/storage/sqlite"
)

func main() {
	conf := config.MustLoad()

	stor, err := sqlite.New(context.Background(), conf)
	if err != nil {
		conf.Logger.Errorf("sqlite: %v", err)
	}

	serv := newServer(conf, stor)
	app := newApp(serv)

	conf.PrintBanner()

	// default RestAddr=":3123"
	if err := app.Listen(conf.RestAddr); err != nil {
		conf.Logger.Errorf("listen: %v", err)
	}
}
