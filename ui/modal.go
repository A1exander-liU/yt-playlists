package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

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
