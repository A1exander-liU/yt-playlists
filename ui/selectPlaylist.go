package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"google.golang.org/api/youtube/v3"
)

type SelectPlaylist struct {
	app       *App
	listModal tview.Primitive
}

func NewSelectPlaylist(a *App, name string, playlists []*youtube.Playlist, selected func(*youtube.Playlist)) *SelectPlaylist {
	sp := SelectPlaylist{app: a}
	list := sp.setupList(playlists, selected)

	sp.listModal = Modal(list, func() { sp.app.CloseModal(name) }, 40, 20)

	return &sp
}

// Helper

func (sp *SelectPlaylist) setupList(playlists []*youtube.Playlist, selected func(*youtube.Playlist)) *tview.List {
	list := tview.NewList().ShowSecondaryText(false).SetHighlightFullLine(true)
	list.SetBorder(true)
	list.SetInputCapture(sp.keyboard)

	for _, playlist := range playlists {
		list.AddItem(playlist.Snippet.Title, "", 0, func() { selected(playlist) })
	}

	return list
}

func (sp *SelectPlaylist) keyboard(event *tcell.EventKey) *tcell.EventKey {
	switch event.Rune() {
	case 'j':
		return tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone)
	case 'k':
		return tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone)
	}

	return event
}
