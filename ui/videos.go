package ui

import (
	"fmt"

	"github.com/A1exander-liU/yt-playlists/controllers"
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
	app              *App
	controller       *controllers.VideosController
	view             *tview.List
	videos           []*youtube.PlaylistItem
	selectedVideos   map[int]*youtube.PlaylistItem
	selectedPlaylist *youtube.Playlist
}

// Creates a new Video component
func NewVideos(a *App, controller *controllers.VideosController) *Video {
	video := Video{
		app:            a,
		controller:     controller,
		view:           tview.NewList(),
		videos:         []*youtube.PlaylistItem{},
		selectedVideos: make(map[int]*youtube.PlaylistItem),
	}
	video.init()

	return &video
}

// Callback when playlist was selected
func (v *Video) PlaylistSelected(playlist *youtube.Playlist) {
	go func() {
		v.controller.SelectedPlaylist = playlist
		v.ClearSelected()
		v.app.QueueUpdateDraw(func() { v.refreshItems() })
	}()
}

func (v *Video) ToggleSelected(i int) {
	videos := v.controller.GetVideos()
	mainText := fmt.Sprintf("%s • %s", videos[i].Snippet.Title, videos[i].Snippet.VideoOwnerChannelTitle)
	v.controller.ToggleSelected(i)

	if v.controller.IsSelectedVideo(i) {
		mainText = fmt.Sprintf("[green]%s", mainText)
	} else {
		mainText = fmt.Sprintf("[white]%s", mainText)
	}

	v.view.SetItemText(i, mainText, "")
}

// removes all selected videos from the selected videos map
func (v *Video) ClearSelected() {
	v.controller.ClearSelectedVides()
}

// removes all selected videos from the map and unselects from ui
func (v *Video) ClearSelectedUI() {
	selectedVideos := v.controller.GetSelectedVideosMap()
	keys := make([]int, 0, len(selectedVideos))
	for k := range selectedVideos {
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
		videos, err := v.app.api.PlaylistItems.List(v.selectedPlaylist.Id, []string{"snippet"})
		if err != nil {
			return
		}
		v.videos = videos
		v.ClearSelected()

		v.app.QueueUpdateDraw(func() { v.refreshItems() })
	}()
}

// Adds the selected videos from current playlist to one specified by 'playlistId'
func (v *Video) AddVideos(playlistId string) {
	go func() {
		v.controller.AddVideos(playlistId)

		v.ClearSelected()
		v.app.QueueUpdateDraw(func() { v.refreshItems() })
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
		videos, err := v.app.api.PlaylistItems.List(v.selectedPlaylist.Id, []string{"snippet"})
		if err != nil {
			return
		}
		v.videos = videos
		v.ClearSelected()
		v.refreshItems()

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

	for _, video := range v.controller.GetVideos() {
		mainText := fmt.Sprintf("[white]%s • %s", video.Snippet.Title, video.Snippet.VideoOwnerChannelTitle)
		v.view.AddItem(mainText, "", 0, nil)
	}
}

func (v *Video) addVideosFlow() {
	playlists, err := v.app.api.Playlists.List([]string{"snippet"})
	if err != nil {
		return
	}
	filtered := v.app.playlistController.ExcludeFromPlaylists(playlists, v.controller.SelectedPlaylist.Id)

	sp := NewSelectPlaylist(v.app, "Add", filtered, func(p *youtube.Playlist) {
		v.AddVideos(p.Id)
		v.app.CloseModal("Add")
	})

	v.app.QueueUpdateDraw(func() {
		v.app.DisplayModal(sp.listModal, "Add")
	})
}

func (v *Video) moveVideosFlow() {
	playlists, err := v.app.api.Playlists.List([]string{"snippet"})
	if err != nil {
		return
	}

	filtered := v.app.playlistController.ExcludeFromPlaylists(playlists, v.selectedPlaylist.Id)

	sp := NewSelectPlaylist(v.app, "Move", filtered, func(p *youtube.Playlist) {
		v.MoveVideos(p.Id)
		v.app.CloseModal("Move")
	})

	v.app.QueueUpdateDraw(func() {
		v.app.DisplayModal(sp.listModal, "Move")
	})
}

// Message to display for dialogs confirming actions to add, move, or delete videos from playlists.
// Will display name of video if there is only one otherwise it will list the amount of videos.
func (v *Video) dialogActionMessage(verb string) string {
	message := fmt.Sprintf("%v %v videos", verb, len(v.selectedVideos))
	if len(v.selectedVideos) == 1 {
		oneVideo := v.firstSelectedVideo()
		message = fmt.Sprintf("%v %v", verb, oneVideo.Snippet.Title)
	}

	return message
}

func (v *Video) firstSelectedVideo() *youtube.PlaylistItem {
	for _, video := range v.selectedVideos {
		return video
	}
	return nil
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
		v.ClearSelectedUI()
	case 'a':
		if len(v.controller.GetSelectedVideos()) == 0 {
			return nil
		}
		go v.addVideosFlow()
	case 'm':
		if len(v.selectedVideos) == 0 {
			return nil
		}
		go v.moveVideosFlow()
	case 'd':
		if len(v.selectedVideos) == 0 {
			return nil
		}

		message := fmt.Sprintf("%v from %v", v.dialogActionMessage("Delete"), v.selectedPlaylist.Snippet.Title)
		dialog := Dialog(message, func() { v.DeleteVideos() }, func() { v.app.CloseModal("Delete") })
		v.app.DisplayModal(dialog, "Delete")

	}

	return event
}
