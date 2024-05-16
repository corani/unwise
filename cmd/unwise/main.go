package main

import (
	"github.com/corani/unwise/internal/config"
	"github.com/corani/unwise/internal/storage"
)

func main() {
	conf := config.MustLoad()
	stor := storage.New(conf)
	serv := newServer(conf, stor)
	app := newApp(serv)

	conf.PrintBanner()

	// default RestAddr=":3123"
	if err := app.Listen(conf.RestAddr); err != nil {
		conf.Logger.Errorf("listen: %v", err)
	}
}
