package main

func main() {
	server := newServer()
	app := newApp(server)

	// default RestAddr=":3123"
	if err := app.Listen(server.conf.RestAddr); err != nil {
		server.conf.Logger.Errorf("listen: %v", err)
	}
}
