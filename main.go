package main

import (
	"os"

	"example.com/demo/ui"
)

func main() {
	app := ui.New()
	app.Init()
	if err := app.Run(); err != nil {
		app.Stop()
		os.Exit(1)
	}
}
