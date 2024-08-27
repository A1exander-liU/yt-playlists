package ui

import (
	"github.com/rivo/tview"
	"google.golang.org/api/youtube/v3"
)

type Video struct {
	app    *App
	view   *tview.List
	videos []*youtube.PlaylistItem
}

func NewVideos(a *App) *Video {
	video := Video{
		app:    a,
		view:   tview.NewList(),
		videos: []*youtube.PlaylistItem{},
	}
	video.Init()

	return &video
}

func (v *Video) Init() {
	v.view.SetHighlightFullLine(true).ShowSecondaryText(true).SetBorder(true).SetTitle("Videos").SetBorderPadding(0, 0, 1, 1)
}

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
func (p *Video) refreshItems() {
	p.view.Clear()

	for _, video := range p.videos {
		p.view.AddItem(video.Snippet.Title, video.Snippet.VideoOwnerChannelTitle, 0, nil)
	}
}
