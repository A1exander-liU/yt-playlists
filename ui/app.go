package ui

import (
	"example.com/demo/api"
	"github.com/gdamore/tcell/v2"
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
	state := appState{api: api.New(), playlists: []*youtube.Playlist{}, videos: []*youtube.PlaylistItem{}}

	root := tview.NewFlex()

	playlists := playlistsView(&state, app)
	videos := videosView()
	help := tview.NewBox().SetTitle("Help").SetBorder(true).SetBorderPadding(0, 0, 1, 1)

	root.AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(playlists, 0, 3, true).
		AddItem(help, 0, 1, false), 0, 1, false).
		AddItem(videos, 0, 3, false)

	pages := tview.NewPages()
	pages.AddPage("Main", root, true, true)

	views := []tview.Primitive{playlists, help, videos}
	currIndex := 0

	root.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// switch between the views
		switch event.Rune() {
		case 'h':
			if currIndex > 0 {
				currIndex--
				app.SetFocus(views[currIndex])
			}
		case 'l':
			if currIndex < len(views)-1 {
				currIndex++
				app.SetFocus(views[currIndex])
			}
		}

		return event
	})

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'q':
			app.Stop()
		}

		return event
	})

	app.SetRoot(pages, true).EnableMouse(true).SetFocus(playlists)
	return app
}
