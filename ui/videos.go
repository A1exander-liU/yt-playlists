package ui

import (
	"fmt"

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

// Helpers

// Initializes the component
func (v *Video) init() {
	v.view.SetHighlightFullLine(true).ShowSecondaryText(false).SetWrapAround(false).SetBorder(true).SetTitle("Videos").SetBorderPadding(0, 0, 1, 1)
	v.view.SetSelectedFunc(func(i int, s1, s2 string, r rune) { v.ToggleSelected(i) })
	v.view.SetSelectedBackgroundColor(tcell.NewRGBColor(40, 44, 52))
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
	}

	return event
}
