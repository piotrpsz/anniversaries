package lays

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
)

type BaseLineAlignLayout struct{}

func NewBaseLineAlignLaoyout() *BaseLineAlignLayout {
	return new(BaseLineAlignLayout)
}

func (l *BaseLineAlignLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	var size fyne.Size
	size.Width += theme.Padding()
	size.Height += theme.Padding()

	for _, item := range objects {
		if text, ok := item.(*canvas.Text); ok {
			ms := measureText(text.Text, text.TextSize, text.TextStyle)
			size.Width += theme.Padding()
			size.Width += ms.Width
			if ms.Height > size.Height {
				size.Height = ms.Height
			}
			size.Width += theme.Padding()
		}
	}
	size.Width += theme.Padding()
	size.Height += 2 * theme.Padding()
	return size
}

func (l *BaseLineAlignLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	maxHeight := float32(0)
	for _, item := range objects {
		if text, ok := item.(*canvas.Text); ok {
			ms := measureText(text.Text, text.TextSize, text.TextStyle)
			if ms.Height > maxHeight {
				maxHeight = ms.Height
			}
		}
	}

	y := float32(0)
	x := theme.Padding()
	for _, item := range objects {
		if text, ok := item.(*canvas.Text); ok {
			size := text.MinSize()

			x += theme.Padding()
			y = theme.Padding() + (maxHeight-size.Height)/2
			text.Move(fyne.NewPos(x, y))
			x += size.Width + theme.Padding()
		}
	}

	// baseY := 2*theme.Padding() + maxHeight
	// x := theme.Padding()
	// for _, item := range objects {
	// 	if text, ok := item.(*canvas.Text); ok {
	// 		x += theme.Padding()
	// 		size := text.MinSize()
	// 		text.Resize(size)
	// 		text.Move(custom.NewPos(x, baseY-size.Height))
	// 		x += size.Width + theme.Padding()
	// 	}
	// }
}

func measureText(text string, size float32, style fyne.TextStyle) fyne.Size {
	return fyne.MeasureText(text, size, style)
}
