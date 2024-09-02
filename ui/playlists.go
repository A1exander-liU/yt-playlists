package ui

import (
	"fmt"

	"github.com/A1exander-liU/yt-playlists/controllers"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"google.golang.org/api/youtube/v3"
)

// Listen for playlist being selected
type SelectedPlaylistListener interface {
	PlaylistSelected(*youtube.Playlist)
}

// Component to display playlists
type Playlist struct {
	// the application instance
	app *App

	// controller to manage data in the playlists view
	controller *controllers.PlaylistsController

	// root view of the component
	view *tview.List

	// list of current listeners for playlist selected event
	listeners []SelectedPlaylistListener
}

// Creates a new Playlist component
func NewPlaylists(a *App, p *controllers.PlaylistsController) *Playlist {
	playlist := Playlist{
		app:        a,
		controller: p,
		view:       tview.NewList().SetHighlightFullLine(true),
		listeners:  []SelectedPlaylistListener{},
	}
	playlist.init()

	return &playlist
}

// Adds a new listener
func (p *Playlist) AddListener(listener SelectedPlaylistListener) {
	p.listeners = append(p.listeners, listener)
}

// Notifie listeners when playlist was selected
func (p *Playlist) NotifySelected(playlist *youtube.Playlist) {
	for _, listener := range p.listeners {
		listener.PlaylistSelected(playlist)
	}
}

// Helpers

// Initializes the component
func (p *Playlist) init() {
	p.view.SetHighlightFullLine(true).ShowSecondaryText(false).SetWrapAround(false).SetBorder(true).SetTitle("Playlists").SetBorderPadding(0, 0, 1, 1)
	p.view.SetSelectedBackgroundColor(COLOR_HIGHLIGHT)

	p.view.SetSelectedFunc(p.selected)
	p.view.SetInputCapture(p.keyboard)

	go func() {
		p.controller.SyncPlaylists()
		p.app.QueueUpdateDraw(func() { p.SetPlaylists() })
	}()
}

func (p *Playlist) SetPlaylists() {
	var selectedPlaylistId string

	if p.controller.GetSelectedPlaylist() >= 0 {
		selectedPlaylistId = p.controller.GetPlaylists()[p.controller.GetSelectedPlaylist()].Id
	}

	p.view.Clear()

	for _, playlist := range p.controller.GetPlaylists() {
		mainText := fmt.Sprintf("[white]%s", playlist.Snippet.Title)
		if playlist.Id == selectedPlaylistId {
			mainText = fmt.Sprintf("[green]%s", playlist.Snippet.Title)
		}
		p.view.AddItem(mainText, "", 0, nil)
	}
}

func (p *Playlist) deletePlaylistFlow() {
	current := p.view.GetCurrentItem()

	confirm := func() {
		p.controller.DeletePlaylistId(current)

		if current == p.controller.GetSelectedPlaylist() {
			p.controller.SetSelectedPlaylist(-1)
		}
		p.app.QueueUpdateDraw(func() {
			p.SetPlaylists()
		})
	}

	message := fmt.Sprintf("Delete the playlist: %s ?", p.controller.GetPlaylists()[current].Snippet.Title)
	dialog := Dialog(
		message,
		func() {
			go confirm()
		},
		func() { p.app.CloseModal("Delete") })
	p.app.DisplayModal(dialog, "Delete")
}

// Handles keyboard input
func (p *Playlist) keyboard(event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyTAB {
		return nil
	}

	switch event.Rune() {
	case 'a':
		NewPlaylistForm(p.app).
			SetAfterSubmitFunc(func(err error) { p.SetPlaylists() }).
			Show()
	case 'd':
		p.deletePlaylistFlow()
	case 'j':
		return tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone)
	case 'k':
		return tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone)
	}

	return event
}

// Callback when item is selected (pressing <space> or <enter>) in the list.
func (p *Playlist) selected(i int, s1, s2 string, r rune) {
	playlists := p.controller.GetPlaylists()
	prevSelectedPlaylist := p.controller.GetSelectedPlaylist()

	if prevSelectedPlaylist >= 0 {
		prevSelected := fmt.Sprintf("[white]%v", playlists[prevSelectedPlaylist].Snippet.Title)
		p.view.SetItemText(prevSelectedPlaylist, prevSelected, "")
	}

	p.controller.SetSelectedPlaylist(i)
	newSelectedPlaylist := p.controller.GetSelectedPlaylist()
	newSelected := fmt.Sprintf("[green]%v", playlists[newSelectedPlaylist].Snippet.Title)
	p.view.SetItemText(newSelectedPlaylist, newSelected, "")

	p.NotifySelected(playlists[newSelectedPlaylist])
	p.app.SetFocus(p.app.views["Videos"])
}
