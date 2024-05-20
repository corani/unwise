package main

import (
	"github.com/corani/unwise/internal/config"
	"github.com/corani/unwise/internal/storage/mem"
	"github.com/corani/unwise/internal/storage/sqlite"
)

func main() {
	conf := config.MustLoad()
	stor := mem.New(conf)
	serv := newServer(conf, stor)
	app := newApp(serv)

	_, err := sqlite.New(conf)
	if err != nil {
		conf.Logger.Errorf("sqlite: %v", err)
	}

	conf.PrintBanner()

	// default RestAddr=":3123"
	if err := app.Listen(conf.RestAddr); err != nil {
		conf.Logger.Errorf("listen: %v", err)
	}
}
