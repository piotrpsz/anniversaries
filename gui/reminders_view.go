package gui

import (
	"strings"

	"calendar/event"
	"calendar/lays"
	"calendar/model"
	"calendar/shared"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	firstName     = "First name:"
	secondName    = "Second name:"
	lastName      = "Last name:"
	birdthday     = "Date of birth:"
	nameday       = "Name day:"
	oneDayBefore  = "day before reminder"
	twoDaysBefore = "two days before reminder"
	oneWeekBefore = "a reminder the week before"
	cancel        = "Cancel"
	save          = "Save"
)

const (
	BirthdayEntry   = "BirthdayEntry"
	BirthdayInfo    = "BirthdayInfo"
	NamedayEntry    = "NamedayEntry"
	NamedayInfo     = "NamedayInfo"
	FirstNameEntry  = "FirstNameEntry"
	FirstNameInfo   = "FirstNameInfo"
	SecondNameEntry = "SecondNameEntry"
	SecondNameInfo  = "SecondNameInfo"
	LastNameEntry   = "LastNameEntry"
	LastNameInfo    = "LastNameInfo"
	OneDayCheck     = "OneDayCheck"
	TwoDaysCheck    = "TwoDaysCheck"
	OneWeekCheck    = "OneWeekCheck"
	CancelButton    = "CancelButton"
	SaveButton      = "SaveButton"
)

type View struct {
	entries map[string]*widget.Entry
	infos   map[string]*widget.Label
	checks  map[string]*widget.Check
	buttons map[string]*widget.Button
}

func RemindersView() *View {
	v := &View{
		entries: make(map[string]*widget.Entry),
		infos:   make(map[string]*widget.Label),
		checks:  make(map[string]*widget.Check),
		buttons: make(map[string]*widget.Button),
	}

	v.entries[BirthdayEntry] = widget.NewEntry()
	v.entries[NamedayEntry] = widget.NewEntry()
	v.entries[FirstNameEntry] = widget.NewEntry()
	v.entries[SecondNameEntry] = widget.NewEntry()
	v.entries[LastNameEntry] = widget.NewEntry()

	v.infos[BirthdayInfo] = widget.NewLabel("")
	v.infos[NamedayInfo] = widget.NewLabel("")
	v.infos[FirstNameInfo] = widget.NewLabel("")
	v.infos[SecondNameInfo] = widget.NewLabel("")
	v.infos[LastNameInfo] = widget.NewLabel("")

	v.checks[OneDayCheck] = widget.NewCheck(oneDayBefore, func(bool) {})
	v.checks[TwoDaysCheck] = widget.NewCheck(twoDaysBefore, func(bool) {})
	v.checks[OneWeekCheck] = widget.NewCheck(oneWeekBefore, func(bool) {})

	v.buttons[CancelButton] = widget.NewButtonWithIcon(cancel, theme.CancelIcon(), func() {})
	v.buttons[SaveButton] = widget.NewButtonWithIcon(save, theme.DocumentSaveIcon(), func() {})

	go func(stream <-chan *events.Event) {
		for e := range stream {
			switch e.Id() {
			case events.PersonSelected:
				if data := e.Data(); data != nil {
					if value, ok := data["id"]; ok {
						if id, ok := value.(int64); ok {
							v.onlyDisplay(id)
						}
					}
				}
			case events.PersonEdit:
				if data := e.Data(); data != nil {
					if value, ok := data["id"]; ok {
						if id, ok := value.(int64); ok {
							v.editPerson(id)
						}
					}
				}
			case events.AddNewPerson:
				v.newPerson()
			}
		}
	}(events.Instance().Register("reminder view", events.PersonSelected, events.PersonEdit, events.AddNewPerson))
	v.emptyView()
	return v
}

func (v *View) Container() fyne.CanvasObject {
	checks := container.New(layout.NewVBoxLayout(),
		v.checks[OneDayCheck],
		v.checks[TwoDaysCheck],
		v.checks[OneWeekCheck],
	)
	birthdayContainer := container.New(lays.NewFixSizedLayout([]float32{110, 0}, false), v.entries[BirthdayEntry], v.infos[BirthdayInfo])
	namedayContainer := container.New(lays.NewFixSizedLayout([]float32{110, 0}, false), v.entries[NamedayEntry], v.infos[NamedayInfo])
	firstNameContainer := container.New(lays.NewFixSizedLayout([]float32{300, 0}, false), v.entries[FirstNameEntry], v.infos[FirstNameInfo])
	secondNameContainer := container.New(lays.NewFixSizedLayout([]float32{300, 0}, false), v.entries[SecondNameEntry], v.infos[SecondNameInfo])
	lastNameContainer := container.New(lays.NewFixSizedLayout([]float32{500, 500}, false), v.entries[LastNameEntry], v.infos[LastNameInfo])

	personalDataForm := widget.NewForm(
		widget.NewFormItem(firstName, firstNameContainer),
		widget.NewFormItem(secondName, secondNameContainer),
		widget.NewFormItem(lastName, lastNameContainer),
		widget.NewFormItem(birdthday, birthdayContainer),
		widget.NewFormItem(nameday, namedayContainer),
	)

	formAndChecks := container.New(layout.NewVBoxLayout(),
		shared.CenteredText("Personal data"),
		personalDataForm,
		widget.NewSeparator(),
		shared.CenteredText("Reminders"),
		checks,
	)
	buttons := container.New(layout.NewHBoxLayout(),
		layout.NewSpacer(),
		v.buttons[CancelButton],
		v.buttons[SaveButton],
	)
	spacer := layout.NewSpacer()
	view := container.New(layout.NewBorderLayout(formAndChecks, buttons, nil, nil), formAndChecks, spacer, buttons)
	return view
}

func (v *View) refresh() {
	for _, v := range v.entries {
		v.Refresh()
	}
	for _, v := range v.infos {
		v.Refresh()
	}
	for _, v := range v.checks {
		v.Refresh()
	}
}

func (v *View) enableView() {
	for _, v := range v.entries {
		v.Hidden = false
	}
	for _, v := range v.infos {
		v.Hidden = true
	}
	for _, v := range v.checks {
		v.Enable()
	}
}

func (v *View) disableView() {
	for key, value := range v.entries {
		key = strings.Replace(key, "Entry", "Info", 1)
		value.Hidden = true
		info := v.infos[key]
		info.Hidden = false
		info.Text = value.Text
	}
	for _, v := range v.checks {
		v.Disable()
	}
}

// onlyDisplay tylko do wyświetlania danych.
// Używany gdy wyświetlane są dane gdy na liście
// użytkowników zmienio osobę.
func (v *View) onlyDisplay(id int64) {
	v.buttons[CancelButton].Hidden = true
	v.buttons[SaveButton].Hidden = true
	v.disableView()

	if person := model.PersonWithId(id); person != nil {
		v.entries[FirstNameEntry].Text = person.FirstName
		v.infos[FirstNameInfo].Text = person.FirstName
		v.entries[SecondNameEntry].Text = person.SecondName
		v.infos[SecondNameInfo].Text = person.SecondName
		v.entries[LastNameEntry].Text = person.LastName
		v.infos[LastNameInfo].Text = person.LastName

		v.entries[BirthdayEntry].Text = person.Birthday
		v.infos[BirthdayInfo].Text = person.Birthday
		v.entries[NamedayEntry].Text = person.Nameday
		v.infos[NamedayInfo].Text = person.Nameday

		v.checks[OneDayCheck].Checked = person.OneDayState()
		v.checks[TwoDaysCheck].Checked = person.TwoDaysState()
		v.checks[OneWeekCheck].Checked = person.OneWeekState()
	}

	v.refresh()
}

func (v *View) emptyView() {
	for _, item := range v.entries {
		item.Text = ""
	}
	for _, item := range v.infos {
		item.Text = ""
	}
	for _, item := range v.checks {
		item.Checked = true
	}
	// v.entries[FirstNameEntry].FocusGained()
	// Window.Canvas().Focus(myobj)
}

func (v *View) newPerson() {
	v.buttons[CancelButton].Hidden = false
	v.buttons[SaveButton].Hidden = false
	v.enableView()
	v.emptyView()
	v.refresh()

	v.buttons[CancelButton].OnTapped = func() {
		v.cancel()
	}

	v.buttons[SaveButton].OnTapped = func() {
		v.save(0)
	}
}

func (v *View) editPerson(id int64) {
	v.buttons[CancelButton].Hidden = false
	v.buttons[SaveButton].Hidden = false
	v.enableView()
	v.refresh()

	v.buttons[CancelButton].OnTapped = func() {
		v.cancel()
	}
	v.buttons[SaveButton].OnTapped = func() {
		v.save(id)
	}
}

func (v *View) cancel() {
	v.buttons[CancelButton].Hidden = false
	v.buttons[SaveButton].Hidden = false

	v.disableView()
	v.refresh()

	e := events.NewEventWithID(events.RefreshRequest)
	events.Instance().Send(e)
}

func (v *View) save(id int64) {
	v.buttons[CancelButton].Hidden = false
	v.buttons[SaveButton].Hidden = false

	v.disableView()
	v.refresh()

	p := &model.Person{
		Id:         id,
		FirstName:  v.infos[FirstNameInfo].Text,
		SecondName: v.infos[SecondNameInfo].Text,
		LastName:   v.infos[LastNameInfo].Text,
		Birthday:   v.infos[BirthdayInfo].Text,
		Nameday:    v.infos[NamedayInfo].Text,
	}
	p.SetOneDayState(v.checks[OneDayCheck].Checked)
	p.SetTwoDaysState(v.checks[TwoDaysCheck].Checked)
	p.SetOneWeekState(v.checks[OneWeekCheck].Checked)

	p.Save()
	e := events.NewEvent(events.RefreshRequest, events.UserData{"id": p.Id})
	events.Instance().Send(e)
}
