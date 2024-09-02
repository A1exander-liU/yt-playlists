package ui

import (
	"fmt"
	"log"
	"slices"

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

	// root view of the component
	view *tview.List

	// list of playlists
	playlists []*youtube.Playlist

	// list of current listeners for playlist selected event
	listeners []SelectedPlaylistListener

	// currently selected playlist
	selectedPlaylist int
}

// Creates a new Playlist component
func NewPlaylists(a *App) *Playlist {
	playlist := Playlist{
		app:              a,
		view:             tview.NewList().SetHighlightFullLine(true),
		playlists:        []*youtube.Playlist{},
		listeners:        []SelectedPlaylistListener{},
		selectedPlaylist: -1,
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
		playlists, _ := p.app.api.Playlists.List([]string{"snippet"})
		p.app.QueueUpdateDraw(func() { p.SetPlaylists(playlists) })
	}()
}

func (p *Playlist) SetPlaylists(playlists []*youtube.Playlist) {
	var selectedPlaylistId string
	if p.selectedPlaylist >= 0 {
		selectedPlaylistId = p.playlists[p.selectedPlaylist].Id
	}

	p.playlists = playlists
	p.view.Clear()

	for _, playlist := range p.playlists {
		mainText := fmt.Sprintf("[white]%s", playlist.Snippet.Title)
		if playlist.Id == selectedPlaylistId {
			mainText = fmt.Sprintf("[green]%s", playlist.Snippet.Title)
		}
		p.view.AddItem(mainText, "", 0, nil)
	}
}

// Handles keyboard input
func (p *Playlist) keyboard(event *tcell.EventKey) *tcell.EventKey {
	log.Println(len(p.playlists))
	if event.Key() == tcell.KeyTAB {
		return nil
	}

	switch event.Rune() {
	case 'a':
		NewPlaylistForm(p.app).
			SetAfterSubmitFunc(func(playlist *youtube.Playlist, err error) {
				p.SetPlaylists(slices.Insert(p.playlists, 0, playlist))
			}).
			Show()
	case 'd':
		current := p.view.GetCurrentItem()

		confirm := func() {
			p.app.playlistController.DeletePlaylist(p.playlists[current].Id)
			p.app.QueueUpdateDraw(func() {
				p.SetPlaylists(slices.Delete(p.playlists, current, current+1))
			})
		}

		message := fmt.Sprintf("Delete the playlist: %s ?", p.playlists[current].Snippet.Title)
		dialog := Dialog(
			message,
			func() {
				go confirm()
			},
			func() { p.app.CloseModal("Delete") })
		p.app.DisplayModal(dialog, "Delete")
	case 'j':
		return tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone)
	case 'k':
		return tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone)
	}

	return event
}

// Callback when item is selected (pressing <space> or <enter>) in the list.
func (p *Playlist) selected(i int, s1, s2 string, r rune) {
	if i < 0 || i > len(p.playlists)-1 {
		return
	}

	if p.selectedPlaylist >= 0 {
		prevSelected := fmt.Sprintf("[white]%v", p.playlists[p.selectedPlaylist].Snippet.Title)
		p.view.SetItemText(p.selectedPlaylist, prevSelected, "")
	}

	p.selectedPlaylist = i

	newSelected := fmt.Sprintf("[green]%v", p.playlists[p.selectedPlaylist].Snippet.Title)
	p.view.SetItemText(i, newSelected, "")

	p.NotifySelected(p.playlists[i])
}
