package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"google.golang.org/api/youtube/v3"
)

var (
	COLOR_HIGHLIGHT  = tcell.NewRGBColor(45, 49, 57)
	COLOR_BACKGROUND = tcell.NewRGBColor(40, 44, 52)
)

// Component to display videos of a playlist
type Video struct {
	app                *App
	view               *tview.List
	videos             []*youtube.PlaylistItem
	selectedVideos     map[int]*youtube.PlaylistItem
	selectedPlaylistId string
}

// Creates a new Video component
func NewVideos(a *App) *Video {
	video := Video{
		app:            a,
		view:           tview.NewList(),
		videos:         []*youtube.PlaylistItem{},
		selectedVideos: make(map[int]*youtube.PlaylistItem),
	}
	video.init()

	return &video
}

// Callback when playlist was selected
func (v *Video) PlaylistSelected(playlistId string) {
	go func() {
		v.selectedPlaylistId = playlistId
		videos, err := v.app.api.PlaylistItems.List(playlistId, []string{"snippet"})
		if err != nil {
			return
		}
		v.videos = videos
		v.app.QueueUpdateDraw(func() { v.refreshItems() })
	}()
}

func (v *Video) ToggleSelected(i int) {
	mainText := fmt.Sprintf("%s • %s", v.videos[i].Snippet.Title, v.videos[i].Snippet.VideoOwnerChannelTitle)

	if _, ok := v.selectedVideos[i]; ok {
		delete(v.selectedVideos, i)
		mainText = fmt.Sprintf("[white]%s", mainText)
	} else {
		v.selectedVideos[i] = v.videos[i]
		mainText = fmt.Sprintf("[green]%s", mainText)
	}
	v.view.SetItemText(i, mainText, "")
}

func (v *Video) ClearSelected() {
	keys := make([]int, 0, len(v.selectedVideos))
	for k := range v.selectedVideos {
		keys = append(keys, k)
	}

	for _, key := range keys {
		v.ToggleSelected(key)
	}
}

// Moves the selected videos from current playlist to one specified by 'playlistId'
//
// THe selected videos will be removed from the current playlist and will be added
// to the new playlist.
func (v *Video) MoveVideos(playlistId string) {
	go func() {
		videos := make([]*youtube.PlaylistItem, 0, len(v.selectedVideos))
		for _, video := range v.selectedVideos {
			videos = append(videos, video)
		}

		v.app.api.PlaylistItems.Move(playlistId, videos)
		videos, err := v.app.api.PlaylistItems.List(v.selectedPlaylistId, []string{"snippet"})
		if err != nil {
			return
		}
		v.videos = videos
		v.app.QueueUpdateDraw(func() { v.refreshItems() })
	}()
}

// Adds the selected videos from current playlist to one specified by 'playlistId'
func (v *Video) AddVideos(playlistId string) {
	go func() {
		videos := make([]*youtube.PlaylistItem, 0, len(v.selectedVideos))
		for _, video := range v.selectedVideos {
			videos = append(videos, video)
		}

		v.app.api.PlaylistItems.Add(playlistId, videos)
	}()
}

// Deletes the selected videos in the current playlist
func (v *Video) DeleteVideos() {
	go func() {
		ids := make([]string, 0, len(v.selectedVideos))
		for _, video := range v.selectedVideos {
			ids = append(ids, video.Id)
		}

		v.app.api.PlaylistItems.Delete(ids)
		videos, err := v.app.api.PlaylistItems.List(v.selectedPlaylistId, []string{"snippet"})
		if err != nil {
			return
		}
		v.videos = videos

		// clear selected videos
		for k := range v.selectedVideos {
			delete(v.selectedVideos, k)
		}

		v.app.QueueUpdateDraw(func() { v.refreshItems() })
	}()
}

// Helpers

// Initializes the component
func (v *Video) init() {
	v.view.SetHighlightFullLine(true).ShowSecondaryText(false).SetWrapAround(false).SetBorder(true).SetTitle("Videos").SetBorderPadding(0, 0, 1, 1)
	v.view.SetSelectedFunc(func(i int, s1, s2 string, r rune) { v.ToggleSelected(i) })
	v.view.SetSelectedBackgroundColor(COLOR_HIGHLIGHT)
	v.view.SetInputCapture(v.keyboard)
}

// Redraws the video items
func (v *Video) refreshItems() {
	v.view.Clear()

	for _, video := range v.videos {
		mainText := fmt.Sprintf("[white]%s • %s", video.Snippet.Title, video.Snippet.VideoOwnerChannelTitle)
		v.view.AddItem(mainText, "", 0, nil)
	}
}

func (v *Video) addVideosFlow() {
	playlists, err := v.app.api.Playlists.List([]string{"snippet"})
	if err != nil {
		return
	}

	filtered := make([]*youtube.Playlist, 0)
	for _, playlist := range playlists {
		if playlist.Id == v.selectedPlaylistId {
			continue
		}
		filtered = append(filtered, playlist)
	}

	sp := NewSelectPlaylist(v.app, "Add", filtered, func(p *youtube.Playlist) {
		v.AddVideos(p.Id)
		v.app.CloseModal("Add")
	})

	v.app.QueueUpdateDraw(func() {
		v.app.Display(sp.listModal, "Add")
	})
}

func (v *Video) moveVideosFlow() {
	playlists, err := v.app.api.Playlists.List([]string{"snippet"})
	if err != nil {
		return
	}

	filtered := make([]*youtube.Playlist, 0)
	for _, playlist := range playlists {
		if playlist.Id == v.selectedPlaylistId {
			continue
		}
		filtered = append(filtered, playlist)
	}

	sp := NewSelectPlaylist(v.app, "Move", filtered, func(p *youtube.Playlist) {
		v.MoveVideos(p.Id)
		v.app.CloseModal("Move")
	})

	v.app.QueueUpdateDraw(func() {
		v.app.Display(sp.listModal, "Move")
	})
}

// Handles keyboard input
func (v *Video) keyboard(event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyTAB {
		return nil
	}

	switch event.Rune() {
	case 'j':
		return tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone)
	case 'k':
		return tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone)
	case 'x':
		v.ClearSelected()
	case 'a':
		go v.addVideosFlow()
	case 'm':
		go v.moveVideosFlow()
	case 'd':
		videoCount := len(v.selectedVideos)
		if videoCount == 0 {
			return nil
		}

		var message string
		if videoCount == 1 {
			var selectedVideo *youtube.PlaylistItem
			for _, video := range v.selectedVideos {
				selectedVideo = video
			}

			message = fmt.Sprintf("Delete %v from %v", selectedVideo.Snippet.Title, v.selectedPlaylistId)
		} else {
			message = fmt.Sprintf("Delete %v videos from %v?", videoCount, v.selectedPlaylistId)
		}
		dialog := DeleteDialog(message, v.DeleteVideos, func() { v.app.CloseModal("Delete") })
		v.app.DisplayModal(dialog, "Delete", func() { v.app.CloseModal("Delete") })

	}

	return event
}
