package cmd

import (
	"fmt"
	"os"

	"github.com/A1exander-liU/yt-playlists/ui"
	"github.com/spf13/cobra"
)

var rootCommand = &cobra.Command{
	Use:     ui.APP_NAME,
	Short:   "Simple CLI to manage your Youtube playlists",
	Long:    "Simple CLI to manage your Youtube playlists",
	Version: "1",
	Run: func(cmd *cobra.Command, args []string) {
		app := ui.New()
		app.Run()
	},
}

func Execute() {
	if err := rootCommand.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
