package ui

import (
	"os"

	"example.com/demo/api"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type App struct {
	*tview.Application
	api    *api.ApiService
	pages  *tview.Pages
	views  map[string]tview.Primitive
	modals map[string]bool
}

func New() *App {
	app := App{
		Application: tview.NewApplication(),
		api:         api.New(),
		pages:       tview.NewPages(),
		views:       make(map[string]tview.Primitive),
		modals:      make(map[string]bool),
	}
	app.init()

	return &app
}

func (a *App) DisplayModal(p tview.Primitive, name string, onCancel func()) {
	a.pages.AddPage(name, Modal(p, onCancel, 40, 20), true, true)
	a.modals[name] = true
}

func (a *App) CloseModal(name string) {
	a.pages.RemovePage(name)
	delete(a.modals, name)
}

func (a *App) ModalActive() bool {
	return len(a.modals) > 0
}

func (a *App) Run() error {
	return a.Application.Run()
}

// Helper

func (a *App) init() {
	// create views
	playlists := NewPlaylists(a)
	videos := NewVideos(a)
	a.views["Playlists"] = playlists.view
	a.views["Videos"] = videos.view

	main := tview.NewFlex().
		AddItem(playlists.view, 0, 1, true).
		AddItem(videos.view, 0, 3, false)
	help := NewHelp()

	core := tview.NewFlex().SetDirection(tview.FlexRow)
	core.AddItem(main, 0, 10, true)
	core.AddItem(help.view, 1, 1, false)

	// set listeners
	playlists.AddListener(videos)

	// setup

	a.pages.AddPage("Main", core, true, true)

	main.SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
		if a.pages.HasPage("Modal") {
			return action, nil
		}
		return action, event
	})

	a.SetInputCapture(a.keyboard)
	a.Application.SetRoot(a.pages, true).EnableMouse(true)
}

func (a *App) keyboard(event *tcell.EventKey) *tcell.EventKey {
	if event.Rune() == 'q' {
		a.Stop()
		os.Exit(0)
	}

	if event.Rune() == '?' {
		a.toggleHelp()
	}

	if !a.ModalActive() && event.Key() == tcell.KeyTAB {
		if a.views["Playlists"].HasFocus() {
			a.SetFocus(a.views["Videos"])
		} else {
			a.SetFocus(a.views["Playlists"])
		}
	}

	return event
}

func (a *App) toggleHelp() {
	if a.pages.HasPage("Help") {
		a.CloseModal("Help")
	} else {
		view := ListModal("Help", 40, 20, func() { a.CloseModal("Help") })
		a.DisplayModal(view, "Help", func() { a.CloseModal("Help") })
	}
}
