package ext

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type ListItemID = int

type ListExt struct {
	widget.List
}

func NewListExt(length func() int, createItem func() fyne.CanvasObject, updateItem func(int, fyne.CanvasObject)) *ListExt {
	ex := new(ListExt)
	ex.BaseWidget = widget.BaseWidget{}
	ex.Length = length
	ex.CreateItem = createItem
	ex.UpdateItem = updateItem
	ex.ExtendBaseWidget(ex)
	return ex
}
