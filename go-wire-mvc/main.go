package main

import "go-wire-mvc/server"

func main() {
	app, err := server.NewAppServer(".")
	if err != nil {
		panic(err)
	}
	defer app.Stop()

	app.Init()
	app.Run()
}
