package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"google.golang.org/api/youtube/v3"
)

// Listen for playlist being selected
type SelectedPlaylistListener interface {
	PlaylistSelected(string)
}

// Component to display playlists
type Playlist struct {
	// the application instance
	app *App

	// root view of the component
	view *tview.List

	// list of playlists
	playlists []*youtube.Playlist

	// list of current listeners for playlist selected event
	listeners []SelectedPlaylistListener
}

// Creates a new Playlist component
func NewPlaylists(a *App) *Playlist {
	playlist := Playlist{
		app:       a,
		view:      tview.NewList().SetHighlightFullLine(true),
		playlists: []*youtube.Playlist{},
		listeners: []SelectedPlaylistListener{},
	}
	playlist.init()

	return &playlist
}

// Adds a new listener
func (p *Playlist) AddListener(listener SelectedPlaylistListener) {
	p.listeners = append(p.listeners, listener)
}

// Notifie listeners when playlist was selected
func (p *Playlist) NotifySelected(playlistId string) {
	for _, listener := range p.listeners {
		listener.PlaylistSelected(playlistId)
	}
}

// Helpers

// Initializes the component
func (p *Playlist) init() {
	p.view.SetHighlightFullLine(true).ShowSecondaryText(false).SetWrapAround(false).SetBorder(true).SetTitle("Playlists").SetBorderPadding(0, 0, 1, 1)
	p.view.SetSelectedFunc(func(i int, s1, s2 string, r rune) {
		if i < 0 || i > len(p.playlists)-1 {
			return
		}
		p.NotifySelected(p.playlists[i].Id)
	})
	p.view.SetInputCapture(p.keyboard)

	go func() {
		playlists, _ := p.app.api.Playlists.List([]string{"snippet"})
		p.playlists = playlists
		p.app.QueueUpdateDraw(func() { p.refreshItems() })
	}()
}

// Redraws the playlist items
func (p *Playlist) refreshItems() {
	p.view.Clear()

	for _, playlist := range p.playlists {
		p.view.AddItem(playlist.Snippet.Title, "", 0, nil)
	}
}

// Handles keyboard input
func (p *Playlist) keyboard(event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyTAB {
		return nil
	}

	switch event.Rune() {
	case 'j':
		return tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone)
	case 'k':
		return tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone)
	}

	return event
}
