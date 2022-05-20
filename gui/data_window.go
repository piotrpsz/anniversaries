package gui

import (
	events "calendar/event"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func InitDataWindow(win fyne.Window) fyne.Window {
	splitter := container.NewHSplit(personListContainer(), RemindersView().Container())
	win.SetContent(splitter)
	splitter.SetOffset(0.25)
	e := events.NewEvent(events.PersonSelected, events.UserData{"id": int64(1)})
	events.Instance().Send(e)
	return win
}
