package ui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Help struct {
	app  *App
	view *tview.TextView
}

func NewHelp(a *App) *Help {
	help := Help{
		app:  a,
		view: tview.NewTextView(),
	}
	help.init()

	return &help
}

// Hlper
func (h *Help) init() {
	h.SetHelpText(h.app.keys.main, h.app.keys.global)
	h.view.SetMouseCapture(h.mouse)
}

// Updates the
func (h *Help) SetHelpText(keyGroups ...map[string]string) {
	keyTexts := make([]string, 0)
	for _, group := range keyGroups {
		for action, key := range group {
			keyTexts = append(keyTexts, fmt.Sprintf("%s: %s", action, key))
		}
	}

	h.view.SetText(strings.Join(keyTexts, " | "))
}

func (h *Help) mouse(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
	return action, nil
}

type HelpModal struct {
	app  *App
	view tview.Primitive
}

func NewHelpModal(app *App) *HelpModal {
	longestAction := 6
	addKeys := func(listView *tview.List, keys map[string]string) {
		for action, key := range keys {
			actionFormat := fmt.Sprintf("%%-%ds", longestAction)
			actionText := fmt.Sprintf(actionFormat, action)
			text := fmt.Sprintf("%s    %s", actionText, key)
			listView.AddItem(text, "", 0, nil)
		}
	}

	list := tview.NewList().
		ShowSecondaryText(false).
		SetWrapAround(false).
		SetHighlightFullLine(true).
		SetSelectedBackgroundColor(COLOR_HIGHLIGHT).
		SetSelectedTextColor(tcell.ColorWhite)
	list.SetTitle("Help").SetBorder(true)

	list.AddItem("[orange]Global", "", 0, nil)
	addKeys(list, app.keys.global)

	list.AddItem("[orange]Main", "", 0, nil)
	addKeys(list, app.keys.main)

	list.AddItem("[orange]Playlists", "", 0, nil)
	addKeys(list, app.keys.playlists)

	list.AddItem("[orange]Videos", "", 0, nil)
	addKeys(list, app.keys.videos)

	list.AddItem("[orange]Modal", "", 0, nil)
	addKeys(list, app.keys.modal)

	list.AddItem("[orange]Dialog", "", 0, nil)
	addKeys(list, app.keys.dialog)

	list.AddItem("[orange]Form", "", 0, nil)
	addKeys(list, app.keys.form)

	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'j':
			return tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone)
		case 'k':
			return tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone)
		}
		return event
	})

	modal := HelpModal{
		app:  app,
		view: Modal(list, func() { app.CloseModal("Help") }, 40, 20),
	}

	return &modal
}
