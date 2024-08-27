package ui

import (
	"example.com/demo/api"
	"github.com/rivo/tview"
)

type App struct {
	*tview.Application
	api   *api.ApiService
	pages *tview.Pages
	views map[string]tview.Primitive
}

func New() *App {
	app := App{
		Application: tview.NewApplication(),
		api:         api.New(),
		pages:       tview.NewPages(),
		views:       make(map[string]tview.Primitive),
	}

	return &app
}

func (a *App) Init() {
	// create views
	playlists := NewPlaylists(a)
	videos := NewVideos(a)
	main := tview.NewFlex().
		AddItem(playlists.view, 0, 1, true).
		AddItem(videos.view, 0, 3, false)

	// set listeners
	playlists.AddListener(videos)

	// setup

	a.pages.AddPage("Main", main, true, true)

	a.Application.SetRoot(a.pages, true).EnableMouse(true)
}

func (a *App) Run() error {
	return a.Application.Run()
}
