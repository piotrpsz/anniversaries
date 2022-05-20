package gui

import (
	events "calendar/event"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func InitDataWindow(win fyne.Window) fyne.Window {
	splitter := container.NewHSplit(personListContainer(), RemindersView().Container())
	splitter.SetOffset(0.35)
	win.SetContent(splitter)
	e := events.NewEvent(events.PersonSelected, map[string]any{"id": int64(3)})
	events.Instance().Send(e)
	return win
}
