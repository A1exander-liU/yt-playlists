package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"google.golang.org/api/youtube/v3"
)

// Component to display videos of a playlist
type Video struct {
	app            *App
	view           *tview.List
	videos         []*youtube.PlaylistItem
	selectedVideos map[int]*youtube.PlaylistItem
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
		videos, err := v.app.api.PlaylistItems.List(playlistId, []string{"snippet"})
		if err != nil {
			return
		}
		v.videos = videos
		v.app.QueueUpdateDraw(func() { v.refreshItems() })
	}()
}

// Helpers

// Initializes the component
func (v *Video) init() {
	v.view.SetHighlightFullLine(true).ShowSecondaryText(true).SetWrapAround(false).SetBorder(true).SetTitle("Videos").SetBorderPadding(0, 0, 1, 1)
	v.view.SetInputCapture(v.keyboard)
}

// Redraws the video items
func (v *Video) refreshItems() {
	v.view.Clear()

	for _, video := range v.videos {
		v.view.AddItem(video.Snippet.Title, video.Snippet.VideoOwnerChannelTitle, 0, nil)
	}
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
	}

	return event
}
