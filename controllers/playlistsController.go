package controllers

import (
	"github.com/A1exander-liU/yt-playlists/api"
	"google.golang.org/api/youtube/v3"
)

type PlaylistsController struct {
	api *api.ApiService
}

func NewPlaylistsController(api *api.ApiService) *PlaylistsController {
	return &PlaylistsController{
		api: api,
	}
}

// Retrieves all the playlists of the authenticated user.
func (p *PlaylistsController) GetPlaylists() []*youtube.Playlist {
	playlists, err := p.api.Playlists.List([]string{api.PART_SNIPPET})
	if err != nil {
		return make([]*youtube.Playlist, 0)
	}

	return playlists
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

func (p *PlaylistsController) CreatePlaylist(name, privacyStatus, description string) (*youtube.Playlist, error) {
	return p.api.Playlists.Insert(name, description, privacyStatus)
}

func (p *PlaylistsController) DeletePlaylist(playlistId string) {
	p.api.Playlists.Delete(playlistId)
}
