package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Help struct {
	view *tview.TextView
}

func NewHelp() *Help {
	help := Help{
		view: tview.NewTextView(),
	}
	help.init()

	return &help
}

// Hlper

func (h *Help) init() {
	h.view.SetText("Switch Tab: <tab> | Quit: q | Keybindings: ?")
	h.view.SetMouseCapture(h.mouse)
}

func (h *Help) mouse(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
	return action, nil
}

type HelpModal struct {
	app  *App
	view tview.Primitive
}

func NewHelpModal(app *App) *HelpModal {
	list := tview.NewList().SetTitle("Help").SetBorder(true)

	modal := HelpModal{
		app:  app,
		view: Modal(list, func() { app.CloseModal("Help") }, 40, 20),
	}

	return &modal
}
