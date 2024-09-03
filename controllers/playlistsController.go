package controllers

import (
	"log"
	"slices"

	"github.com/A1exander-liU/yt-playlists/api"
	"google.golang.org/api/youtube/v3"
)

// Holds all functionality to manage the data in the Playlists component.
type PlaylistsController struct {
	// Api instance for accessing the Youtube api
	api *api.ApiService

	// Index position of the selected playlist in PlaylistsController.playlists
	selectedPlaylist int

	// All the playlists
	playlists []*youtube.Playlist
}

// Creates a new PlaylistsController.
func NewPlaylistsController(api *api.ApiService) *PlaylistsController {
	return &PlaylistsController{
		api:              api,
		selectedPlaylist: -1,
		playlists:        make([]*youtube.Playlist, 0),
	}
}

// Retrieves the index of the currently selected playlistts. -1 will be returned if no playlist is selected.
func (p *PlaylistsController) GetSelectedPlaylist() int {
	return p.selectedPlaylist
}

// Updates the index of the currently selected playlist.
func (p *PlaylistsController) SetSelectedPlaylist(i int) {
	p.selectedPlaylist = i
}

// Retrieves all the playlists of the authenticated user.
func (p *PlaylistsController) GetPlaylists() []*youtube.Playlist {
	return p.playlists
}

// Filters the playlists, the playlist with the specified id will be ignored in the result.
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

// Creates a new playlist.
func (p *PlaylistsController) CreatePlaylist(name, description, privacyStatus string) {
	playlist, _ := p.api.Playlists.Insert(name, description, privacyStatus)
	slices.Insert(p.playlists, 0, playlist)
	log.Println("New", playlist.Snippet.Title)
}

// Deletes the currently selected playlist.
func (p *PlaylistsController) DeletePlaylist() {
	p.api.Playlists.Delete(p.playlists[p.selectedPlaylist].Id)
	p.SyncPlaylists()
}

// Deletes the playlist specified by its index position in PlaylistsController.playlists.
func (p *PlaylistsController) DeletePlaylistId(i int) {
	p.api.Playlists.Delete(p.playlists[i].Id)
	p.SyncPlaylists()
}

// Retrieves the currennt playlists. Call this to sync the current playlists with playlists from the API server.
func (p *PlaylistsController) SyncPlaylists() {
	playlists, _ := p.api.Playlists.List([]string{api.PART_SNIPPET})
	p.playlists = playlists
}
