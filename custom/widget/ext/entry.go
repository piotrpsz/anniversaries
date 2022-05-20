package ext

import (
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type EntryExt struct {
	widget.Entry
	enabled bool
}

func NexEntryExt() *EntryExt {
	ex := new(EntryExt)
	ex.ExtendBaseWidget(ex)
	return ex
}

func (ex *EntryExt) Enable() {
	ex.enabled = true
}
func (ex *EntryExt) Disable() {
	ex.enabled = false
}

func (ex *EntryExt) Cursor() desktop.Cursor {
	if ex.enabled {
		ex.Entry.Disable()
		return ex.Entry.Cursor()
	}
	ex.Entry.CursorColumn = 1000
	ex.Entry.Refresh()
	return ex.Entry.Cursor()
}
