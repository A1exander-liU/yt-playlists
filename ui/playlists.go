package ui

import (
	"fmt"
	"os"

	"github.com/rivo/tview"
	"google.golang.org/api/youtube/v3"
)

type SelectedPlaylistListener interface {
	PlaylistSelected(string)
}

type Playlist struct {
	app       *App
	view      *tview.List
	playlists []*youtube.Playlist
	listeners []SelectedPlaylistListener
}

func NewPlaylists(a *App) *Playlist {
	playlist := Playlist{
		app:       a,
		view:      tview.NewList().SetHighlightFullLine(true),
		playlists: []*youtube.Playlist{},
		listeners: []SelectedPlaylistListener{},
	}
	playlist.Init()

	return &playlist
}

func (p *Playlist) AddListener(listener SelectedPlaylistListener) {
	p.listeners = append(p.listeners, listener)
}

func (p *Playlist) Init() {
	p.view.SetHighlightFullLine(true).ShowSecondaryText(false).SetBorder(true).SetTitle("Playlists").SetBorderPadding(0, 0, 1, 1)
	go func() {
		playlists, err := p.app.api.Playlists.List([]string{"snippet"})
		if err != nil {
			os.WriteFile("log.txt", []byte(fmt.Sprintf("%v", err)), 0700)
			return
		}
		p.playlists = playlists
		p.app.QueueUpdateDraw(func() { p.refreshItems() })
	}()
}

func (p *Playlist) NotifySelected(playlistId string) {
	for _, listener := range p.listeners {
		listener.PlaylistSelected(playlistId)
	}
}

// Helpers
func (p *Playlist) refreshItems() {
	p.view.Clear()

	for _, playlist := range p.playlists {
		p.view.AddItem(playlist.Snippet.Title, "", 0, func() {
			p.NotifySelected(playlist.Id)
		})
	}
}
