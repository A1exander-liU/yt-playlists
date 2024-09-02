package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"google.golang.org/api/youtube/v3"
)

// Renders the given primitive as a modal. onCancel function will be called when pressnig 'ESC' key
// which will cancel the modal.
func Modal(p tview.Primitive, onCancel func(), width, height int) tview.Primitive {
	root := tview.NewFlex()

	middleRow := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(p, height, 1, true).
		AddItem(nil, 0, 1, false)

	root.AddItem(nil, 0, 1, false)
	root.AddItem(middleRow, width, 1, true)
	root.AddItem(nil, 0, 1, false)

	root.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyESC {
			onCancel()
		}

		return event
	})

	return root
}

func ListModal(title string, width, height int, onCancel func()) tview.Primitive {
	list := tview.NewList().SetHighlightFullLine(true).SetTitle(title).SetBorder(true)
	return Modal(list, onCancel, width, height)
}

// A Dialog with 2 choices 'Yes' and 'No', pass the confirm function for callback when 'Yes' is selected,
// pass the cancel function when exiting the modal.
func Dialog(message string, confirm func(), cancel func()) *tview.Modal {
	dialog := tview.NewModal().
		SetText(message).
		AddButtons([]string{"No", "Yes"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			switch buttonLabel {
			case "No":
				cancel()
			case "Yes":
				confirm()
				cancel()
			}
		})

	dialog.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case ' ':
			return tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone)
		}

		return event
	})

	return dialog
}

// 2 modals on for list, one fir onfirm
func SelectPlaylistDialog(items []*youtube.Playlist, selected func(*youtube.Playlist), cancel func()) tview.Primitive {
	list := tview.NewList().
		SetHighlightFullLine(true).
		ShowSecondaryText(false)

	for _, item := range items {
		list.AddItem(item.Snippet.Title, "", 0, func() { selected(item) })
	}

	playlistModal := Modal(list, cancel, 40, 20)

	return playlistModal
}
