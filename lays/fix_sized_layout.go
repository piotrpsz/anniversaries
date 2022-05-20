package lays

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type FixSizedLayout struct {
	mins      []float32
	isPadding bool
}

func NewFixSizedLayout(mins []float32, isPadding bool) *FixSizedLayout {
	return &FixSizedLayout{
		mins:      mins,
		isPadding: isPadding,
	}
}

func (fs *FixSizedLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	if len(fs.mins) != len(objects) {
		panic("wrong number of items")
	}

	w := float32(0)
	h := 2*theme.Padding() + 24.0
	if fs.isPadding {
		w = 2 * theme.Padding()
	}

	for idx, item := range objects {
		if item.Visible() {
			size := item.MinSize()
			if reqwidth := fs.mins[idx]; reqwidth > 0.0 {
				size.Width = reqwidth
			}
			w += size.Width
			if size.Height > h {
				h = size.Height
			}
		}
	}

	if fs.isPadding {
		w += theme.Padding()
	}
	return fyne.NewSize(w, h)
}

func (fs *FixSizedLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	topLeftX := float32(0)
	// topLeftY := float32(0)
	if fs.isPadding {
		topLeftX += 2 * theme.Padding()
	}

	for idx, item := range objects {
		if item.Visible() {
			itemSize := item.MinSize()
			if fs.mins[idx] > 0.0 {
				itemSize.Width = fs.mins[idx]
			}
			itemSize.Height = size.Height
			item.Resize(itemSize)

			{ // wypośrodkowanie w pionie
				x := topLeftX
				y := 0.0 + (size.Height-itemSize.Height)/2.0
				item.Move(fyne.NewPos(x, y))
			}
			topLeftX += itemSize.Width + theme.Padding()
		}
	}
}
