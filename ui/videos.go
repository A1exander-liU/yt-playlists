package ui

import "github.com/rivo/tview"

func videosView() *tview.List {
	view := tview.NewList()
	view.SetTitle("Videos").SetBorder(true).SetBorderPadding(0, 0, 1, 1)

	return view
}
