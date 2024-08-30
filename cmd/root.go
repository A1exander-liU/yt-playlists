package cmd

import (
	"fmt"
	"os"

	"example.com/demo/ui"
	"github.com/spf13/cobra"
)

var rootCommand = &cobra.Command{
	Use:   "yt-playlists",
	Short: "Simple CLI to manage your Youtube playlists",
	Long:  "Simple CLI to manage your Youtube playlists",
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
