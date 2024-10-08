package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Form component for creating a new playlist. Contains a input field for the name,
// a dropdown to choose the privacy status and a text area to fill out the description.
type PlaylistForm struct {
	// an instance of the app
	app *App

	// the root view to display
	view tview.Primitive

	formView *tview.Form

	afterSubmit func(error)
}

// Information in playlist form
type formData struct {
	// name of the playlist
	name string

	// privacy of the playlist, could be one of: 'private', 'public', or 'unlisted'
	privacyStatus string

	// description of the playlist
	description string
}

// Creates a new playlist form component. Accepts an instance of the app and a function 'onSubmit' to be called
// when form was submitted.
func NewPlaylistForm(app *App) *PlaylistForm {
	form := PlaylistForm{
		app: app,
	}
	formView := form.init()
	form.view = Modal(formView, form.Close, 60, 20)
	form.formView = formView

	return &form
}

func (p *PlaylistForm) SetAfterSubmitFunc(afterSubmit func(error)) *PlaylistForm {
	p.afterSubmit = afterSubmit
	return p
}

// shows the form.
func (p *PlaylistForm) Show() {
	p.app.DisplayModal(p.view, "Form")
}

// closes the form.
func (p *PlaylistForm) Close() {
	p.app.CloseModal("Form")
}

// Submits the form and creates the playlist with the current form data.
func (p *PlaylistForm) Submit() error {
	formData := p.getFormData()

	if formData.name == "" {
		return fmt.Errorf("Name must be filled out")
	}

	p.app.playlistController.CreatePlaylist(formData.name, formData.description, formData.privacyStatus)
	return nil
}

// Initializes the component.
func (p *PlaylistForm) init() *tview.Form {
	form := tview.NewForm()

	dropdown := tview.NewDropDown().
		SetLabel("Privacy Status").
		SetOptions([]string{"private", "public", "unlisted"}, nil).
		SetCurrentOption(0)
	dropdown.SetInputCapture(p.keyboardDropdown)

	form.
		AddInputField("Name", "", 30, nil, nil).
		AddFormItem(dropdown).
		AddTextArea("Description", "", 40, 0, 0, nil).
		AddButton("Create", func() {
			err := p.Submit()

			if p.afterSubmit != nil {
				p.afterSubmit(err)
			}

			p.Close()
		}).
		AddButton("Cancel", func() { p.Close() })

	form.GetButton(0).SetInputCapture(p.keyboardButton)
	form.GetButton(1).SetInputCapture(p.keyboardButton)

	form.SetBorder(true).SetTitle("Create Playlist").SetBorderPadding(0, 0, 1, 1)

	return form
}

// Collects all the current form data.
func (p *PlaylistForm) getFormData() formData {
	name := p.formView.GetFormItemByLabel("Name").(*tview.InputField).GetText()
	_, privacyStatus := p.formView.GetFormItemByLabel("Privacy Status").(*tview.DropDown).GetCurrentOption()
	description := p.formView.GetFormItemByLabel("Description").(*tview.TextArea).GetText()

	return formData{
		name:          name,
		privacyStatus: privacyStatus,
		description:   description,
	}
}

// Keyboard input handler for buttons.
func (p *PlaylistForm) keyboardButton(event *tcell.EventKey) *tcell.EventKey {
	switch event.Rune() {
	case ' ':
		return tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone)
	}
	return event
}

// Keyboard input handler callback for dropdown.
func (p *PlaylistForm) keyboardDropdown(event *tcell.EventKey) *tcell.EventKey {
	switch event.Rune() {
	case 'j':
		return tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone)
	case 'k':
		return tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone)
	}

	return event
}
