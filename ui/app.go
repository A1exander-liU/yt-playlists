package ui

import (
	"example.com/demo/api"
	"github.com/rivo/tview"
	"google.golang.org/api/youtube/v3"
)

type appState struct {
	api       *api.ApiService
	playlists []*youtube.Playlist
	videos    []*youtube.PlaylistItem
}

func (a *appState) fetchPlaylists(onFetch func()) {
	go func() {
		if playlists, err := a.api.Playlists.List([]string{"snippet"}); err == nil {
			a.playlists = playlists
			if onFetch != nil {
				onFetch()
			}
		}
	}()
}

func (a *appState) fetchVideos(i int, onFetch func()) {
	go func() {
		if i < 0 || i > len(a.playlists)-1 {
			return
		}
		if videos, err := a.api.PlaylistItems.List(a.playlists[i].Id, []string{"snippet"}); err == nil {
			a.videos = videos
			if onFetch != nil {
				onFetch()
			}
		}
	}()
}

func App() *tview.Application {
	app := tview.NewApplication()

	return app
}
