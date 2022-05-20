package gui

import (
	"calendar/config"
	events "calendar/event"
	alingLayout "calendar/lays"
	"calendar/model"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/image/colornames"
)

var (
	data     []*model.Person
	idx      = -1
	dataList *widget.List
	enabled  = true
)

func personListContainer() fyne.CanvasObject {
	data = model.AllPersons()

	delButton := widget.NewButtonWithIcon("", theme.DeleteIcon(), removePerson)
	delButton.Importance = widget.LowImportance
	editButton := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), editPerson)
	editButton.Importance = widget.LowImportance
	addButton := widget.NewButtonWithIcon("", theme.ContentAddIcon(), addNewPerson)
	addButton.Importance = widget.LowImportance
	icons := container.New(layout.NewHBoxLayout(), delButton, editButton, addButton, layout.NewSpacer())

	dataList = widget.NewList(
		func() int {
			return len(data)
		},
		func() fyne.CanvasObject {
			firstName := canvas.NewText("", config.Config().TextColor())
			firstName.TextSize = float32(16.0)
			lastName := canvas.NewText("", colornames.Darkgray)
			lastName.TextSize = float32(12.0)
			return container.New(alingLayout.NewBaseLineAlignLaoyout(), firstName, lastName)
		},
		func(idx widget.ListItemID, template fyne.CanvasObject) {
			person := data[idx]
			content := template.(*fyne.Container)

			firstName := content.Objects[0].(*canvas.Text)
			firstName.Color = config.Config().TextColor()
			firstName.Text = person.FirstName

			lastName := content.Objects[1].(*canvas.Text)
			lastName.Color = colornames.Darkgray
			lastName.Text = person.LastName
			if person.SecondName != "" {
				lastName.Text = person.SecondName + " " + person.LastName
			}

		},
	)
	dataList.OnSelected = func(id widget.ListItemID) {
		if !enabled {
			dataList.Unselect(id)
			return
		}
		idx = id
		e := events.NewEvent(events.PersonSelected, map[string]any{"id": data[id].Id})
		events.Instance().Send(e)
	}

	go func(stream <-chan *events.Event) {
		for e := range stream {
			switch e.Id() {
			case events.RefreshRequest:
				data = model.AllPersons()
				enabled = true
				updateSelection(e.Data())
				dataList.Refresh()
			case events.KeyEvent:
				if len(data) > 0 && idx != -1 {
					if data := e.Data(); data != nil {
						if value, ok := data["event"]; ok {
							switch value.(*fyne.KeyEvent).Name {
							case fyne.KeyUp:
								keyUp()
							case fyne.KeyDown:
								keyDown()
							case fyne.KeyHome:
								keyHome()
							case fyne.KeyEnd:
								keyEnd()
							}
						}
					}
				}
			}
		}
	}(events.Instance().Register("person widget", events.RefreshRequest, events.KeyEvent))
	dataList.Select(0)

	spacer := layout.NewSpacer()
	return container.New(layout.NewBorderLayout(spacer, icons, nil, nil),
		spacer,
		dataList,
		icons,
	)
}

func updateSelection(data events.UserData) {
	if data != nil {
		if idv, ok := data["id"]; ok {
			if id, ok := idv.(int64); ok {
				dataList.Select(idxForId(id))
				return
			}
		}
	}
	dataList.Select(idx)
}

func idxForId(id int64) int {
	for idx, item := range data {
		if item.Id == id {
			return idx
		}
	}
	return -1
}

func saveState() {
	dataList.UnselectAll()
	enabled = false
}

func addNewPerson() {
	if enabled {
		saveState()
		e := events.NewEventWithID(events.AddNewPerson)
		events.Instance().Send(e)
	}
}

func editPerson() {
	if enabled {
		saveState()
		e := events.NewEvent(events.PersonEdit, events.UserData{"id": data[idx].Id})
		events.Instance().Send(e)
	}
}

func removePerson() {
	if enabled {
		n := len(data)
		if n == 0 {
			return
		}
		if idx >= n {
			idx = n - 1
		}

		model.DeleteWithId(data[idx].Id)
		data = model.AllPersons()
		dataList.Select(idx)
		dataList.Refresh()
	}
}

// ******************************************************************
// *                                                                *
// *        D A T A   L I S T 'S   K E Y   H A N D L E R S          *
// *                                                                *
// ******************************************************************

func keyUp() {
	if idx > 0 {
		dataList.Select(idx - 1)
	}
}

func keyDown() {
	if idx < len(data)-1 {
		dataList.Select(idx + 1)
	}
}

func keyHome() {
	if idx > 0 {
		idx = 0
		dataList.Select(idx)
	}
}

func keyEnd() {
	if idx < len(data)-1 {
		dataList.Select(len(data) - 1)
	}
}
