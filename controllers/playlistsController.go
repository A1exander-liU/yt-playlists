package controllers

import (
	"slices"

	"github.com/A1exander-liU/yt-playlists/api"
	"google.golang.org/api/youtube/v3"
)

type PlaylistsController struct {
	api *api.ApiService

	// Index position of the selected playlist in PlaylistsController.playlists.
	selectedPlaylist int

	// All the playlists
	playlists []*youtube.Playlist
}

func NewPlaylistsController(api *api.ApiService) *PlaylistsController {
	return &PlaylistsController{
		api:              api,
		selectedPlaylist: -1,
		playlists:        make([]*youtube.Playlist, 0),
	}
}

func (p *PlaylistsController) GetSelectedPlaylist() int {
	return p.selectedPlaylist
}

func (p *PlaylistsController) SetSelectedPlaylist(i int) {
	p.selectedPlaylist = i
}

// Retrieves all the playlists of the authenticated user.
func (p *PlaylistsController) GetPlaylists() []*youtube.Playlist {
	return p.playlists
}

// Filters the playlists.
func (p *PlaylistsController) ExcludeFromPlaylists(playlists []*youtube.Playlist, playlistId string) []*youtube.Playlist {
	filtered := make([]*youtube.Playlist, 0)

	for _, playlist := range playlists {
		if playlist.Id == playlistId {
			continue
		}
		filtered = append(filtered, playlist)
	}

	return filtered
}

func (p *PlaylistsController) CreatePlaylist(name, privacyStatus, description string) {
	playlist, _ := p.api.Playlists.Insert(name, description, privacyStatus)
	slices.Insert(p.playlists, 0, playlist)
}

func (p *PlaylistsController) DeletePlaylist() {
	p.api.Playlists.Delete(p.playlists[p.selectedPlaylist].Id)
	p.SyncPlaylists()
}

func (p *PlaylistsController) DeletePlaylistId(i int) {
	p.api.Playlists.Delete(p.playlists[i].Id)
	p.SyncPlaylists()
}

func (p *PlaylistsController) SyncPlaylists() {
	playlists, _ := p.api.Playlists.List([]string{api.PART_SNIPPET})
	p.playlists = playlists
}
