package controllers

import (
	"github.com/A1exander-liU/yt-playlists/api"
	"google.golang.org/api/youtube/v3"
)

// Handles all functionality to managing data on the Videos component.
type VideosController struct {
	api *api.ApiService

	// The currently selected playlist
	SelectedPlaylist *youtube.Playlist

	// To store all the current videos of the selected playlist
	videos []*youtube.PlaylistItem

	// Selected videos to perform operations on like: adding, moving, deleting to and from playlists. Stored as a map where
	// the key is the index into VideosController.videos.
	selectedVideos map[int]*youtube.PlaylistItem
}

// Creates a new VideosController
func NewVideosController(api *api.ApiService) *VideosController {
	controller := VideosController{
		api:            api,
		videos:         make([]*youtube.PlaylistItem, 0),
		selectedVideos: make(map[int]*youtube.PlaylistItem),
	}

	return &controller
}

// Retrieves the current videos in this object. Calls VideosController.syncVideos to maintain updated.
func (v *VideosController) GetVideos() []*youtube.PlaylistItem {
	v.syncVideos()
	return v.videos
}

func (v *VideosController) GetSelectedVideosMap() map[int]*youtube.PlaylistItem {
	return v.selectedVideos
}

// Retrieves the currently selected videos.
func (v *VideosController) GetSelectedVideos() []*youtube.PlaylistItem {
	videos := make([]*youtube.PlaylistItem, 0)

	for _, video := range v.selectedVideos {
		videos = append(videos, video)
	}

	return videos
}

// Toggles the selection of the video based on the index 'i' in VideosController.videos.
func (v *VideosController) ToggleSelected(i int) {
	if _, ok := v.selectedVideos[i]; ok {
		delete(v.selectedVideos, i)
	} else {
		v.selectedVideos[i] = v.videos[i]
	}
}

func (v *VideosController) IsSelectedVideo(i int) bool {
	_, ok := v.selectedVideos[i]
	return ok
}

// Removes all selected videos.
func (v *VideosController) ClearSelectedVides() {
	for k := range v.selectedVideos {
		delete(v.selectedVideos, k)
	}
}

// Adds the selected videos to the playlist specified by playlistId. The selected videos will be inserted into the playlist.
func (v *VideosController) AddVideos(playlistId string) {
	v.api.PlaylistItems.Add(playlistId, v.GetSelectedVideos())
	v.syncVideos()
}

// Moves the selected videos from the current playlist (VideosController.selectedPlaylist) to the playlist specified by playlistId.
// The selected videos will be deleted from the current playlist and then inserted in the new playlist.
func (v *VideosController) MoveVideos(playlistId string) {
	v.api.PlaylistItems.Move(playlistId, v.GetSelectedVideos())
	v.syncVideos()
}

// Deletes the selected videos from the current playlist (VideosController.selectedPlaylist).
func (v *VideosController) DeleteVideos() {
	ids := make([]string, 0)
	for _, video := range v.selectedVideos {
		ids = append(ids, video.Id)
	}

	v.api.PlaylistItems.Delete(ids)
	v.syncVideos()
}

// Retrieves the currennt videos of the selected playlist. Call this to sync the current videos with videos from the API server.
func (v *VideosController) syncVideos() {
	videos, _ := v.api.PlaylistItems.List(v.SelectedPlaylist.Id, []string{api.PART_SNIPPET})
	v.videos = videos
}
