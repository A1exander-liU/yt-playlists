package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func setPlaylistItems(a *appState, listView *tview.List) {
	listView.Clear()
	for _, p := range a.playlists {
		listView.AddItem(p.Snippet.Title, "", 0, nil)
	}
}

func playlistsView(a *appState, app *tview.Application) *tview.List {
	view := tview.NewList().ShowSecondaryText(false)
	view.SetTitle("Playlists").SetBorder(true).SetBorderPadding(0, 0, 1, 1)

	view.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// use k and j also to selected nex and prev list item
		switch event.Rune() {
		case 'k':
			return tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone)
		case 'j':
			return tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone)
		}
		return event
	})

	a.fetchPlaylists(func() {
		app.QueueUpdateDraw(func() {
			setPlaylistItems(a, view)
		})
	})

	return view
}
