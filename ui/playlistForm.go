package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type PlaylistForm struct {
	app  *App
	view *tview.Form
}

type formData struct {
	name          string
	privacyStatus string
	description   string
}

// creates a new playlist form
func NewPlaylistForm(app *App, onSubmit func()) *PlaylistForm {
	form := PlaylistForm{
		app:  app,
		view: tview.NewForm(),
	}
	form.init(onSubmit)

	return &form
}

// shows the form
func (p *PlaylistForm) Show() {
	p.app.Display(p.view, "Form")
}

// closes the form
func (p *PlaylistForm) Close() {
	p.app.CloseModal("Form")
}

func (p *PlaylistForm) Submit() {
	formData := p.getFormData()

	if formData.name == "" {
		return
	}

	p.app.api.Playlists.Insert(formData.name, formData.description, formData.privacyStatus)
}

// initializes the component
func (p *PlaylistForm) init(onSubmit func()) {
	dropdown := tview.NewDropDown().
		SetLabel("Privacy Status").
		SetOptions([]string{"private", "public", "unlisted"}, nil).
		SetCurrentOption(0)
	dropdown.SetInputCapture(p.keyboardDropdown)

	p.view.
		AddInputField("Name", "", 30, nil, nil).
		AddFormItem(dropdown).
		AddTextArea("Description", "", 40, 0, 0, nil).
		AddButton("Create", func() {
			p.Submit()
			onSubmit()
			p.Close()
		}).
		AddButton("Cancel", func() { p.Close() })

	p.view.GetButton(0).SetInputCapture(p.keyboardButton)
	p.view.GetButton(1).SetInputCapture(p.keyboardButton)

	p.view.SetBorder(true).SetTitle("Create Playlist").SetBorderPadding(0, 0, 1, 1)
}

// returns a struct containing all current form data
func (p *PlaylistForm) getFormData() formData {
	name := p.view.GetFormItemByLabel("Name").(*tview.InputField).GetText()
	_, privacyStatus := p.view.GetFormItemByLabel("Privacy Status").(*tview.DropDown).GetCurrentOption()
	description := p.view.GetFormItemByLabel("Description").(*tview.TextArea).GetText()

	return formData{
		name:          name,
		privacyStatus: privacyStatus,
		description:   description,
	}
}

// keyboard input handler for buttons
func (p *PlaylistForm) keyboardButton(event *tcell.EventKey) *tcell.EventKey {
	switch event.Rune() {
	case ' ':
		return tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone)
	}
	return event
}

// keyboard input handler callback for dropdown
func (p *PlaylistForm) keyboardDropdown(event *tcell.EventKey) *tcell.EventKey {
	switch event.Rune() {
	case 'j':
		return tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone)
	case 'k':
		return tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone)
	}

	return event
}
