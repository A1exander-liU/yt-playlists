package main

import (
	"log"
	"os"

	"github.com/A1exander-liU/yt-playlists/cmd"
)

func main() {
	if file, err := os.OpenFile("log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
		log.SetOutput(file)
	}
	cmd.Execute()
}
