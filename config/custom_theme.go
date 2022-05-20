package config

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type CustomTheme struct {
	fyne.Theme
	variant_ fyne.ThemeVariant
}

var _ fyne.Theme = (*CustomTheme)(nil)

func NewCustomTheme(dark bool) *CustomTheme {
	ct := new(CustomTheme)
	if dark {
		ct.variant_ = theme.VariantDark
	} else {
		ct.variant_ = theme.VariantLight
	}
	return ct
}

func (m CustomTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {

	switch m.variant_ {
	case theme.VariantDark:
		switch name {
		case theme.ColorNameBackground:
			return color.NRGBA{R: 0x30, G: 0x30, B: 0x30, A: 0xff}
		case theme.ColorNameButton:
			return color.Transparent
		case theme.ColorNameDisabled:
			return color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0x42}
			// return color.NRGBA{R: 0x10, G: 0xff, B: 0xff, A: 0x62}
		case theme.ColorNameDisabledButton:
			return color.NRGBA{R: 0x26, G: 0x26, B: 0x26, A: 0xff}
		case theme.ColorNameError:
			return color.NRGBA{R: 0xf4, G: 0x43, B: 0x36, A: 0xff}
		case theme.ColorNameForeground:
			return color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
		case theme.ColorNameHover:
			return color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0x0f}
		case theme.ColorNameInputBackground:
			return color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0x19}
		case theme.ColorNamePlaceHolder:
			return color.NRGBA{R: 0xb2, G: 0xb2, B: 0xb2, A: 0xff}
		case theme.ColorNamePressed:
			return color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0x66}
		case theme.ColorNameScrollBar:
			return color.NRGBA{A: 0x99}
		case theme.ColorNameShadow:
			return color.NRGBA{A: 0x66}
		}
	case theme.VariantLight:
		switch name {
		case theme.ColorNameBackground:
			return color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
		case theme.ColorNameButton:
			return color.Transparent
		case theme.ColorNameDisabled:
			return color.NRGBA{A: 0x42}
		case theme.ColorNameDisabledButton:
			return color.NRGBA{R: 0xe5, G: 0xe5, B: 0xe5, A: 0xff}
		case theme.ColorNameError:
			return color.NRGBA{R: 0xf4, G: 0x43, B: 0x36, A: 0xff}
		case theme.ColorNameForeground:
			return color.NRGBA{R: 0x21, G: 0x21, B: 0x21, A: 0xff}
		case theme.ColorNameHover:
			return color.NRGBA{A: 0x0f}
		case theme.ColorNameInputBackground:
			return color.NRGBA{A: 0x19}
		case theme.ColorNamePlaceHolder:
			return color.NRGBA{R: 0x88, G: 0x88, B: 0x88, A: 0xff}
		case theme.ColorNamePressed:
			return color.NRGBA{A: 0x19}
		case theme.ColorNameScrollBar:
			return color.NRGBA{A: 0x99}
		case theme.ColorNameShadow:
			return color.NRGBA{A: 0x33}
		}
	}

	return theme.DefaultTheme().Color(name, variant)
}

func (m CustomTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (m CustomTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (m CustomTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

func changeTheme(dark bool) {
	fyne.CurrentApp().Settings().SetTheme(NewCustomTheme(dark))
}

/*  ----- dark -----
ColorNameBackground:      color.NRGBA{0x30, 0x30, 0x30, 0xff},
ColorNameButton:          color.Transparent,
ColorNameDisabled:        color.NRGBA{0xff, 0xff, 0xff, 0x42},
ColorNameDisabledButton:  color.NRGBA{0x26, 0x26, 0x26, 0xff},
ColorNameError:           errorColor,
ColorNameForeground:      color.NRGBA{0xff, 0xff, 0xff, 0xff},
ColorNameHover:           color.NRGBA{0xff, 0xff, 0xff, 0x0f},
ColorNameInputBackground: color.NRGBA{0xff, 0xff, 0xff, 0x19},
ColorNamePlaceHolder:     color.NRGBA{0xb2, 0xb2, 0xb2, 0xff},
ColorNamePressed:         color.NRGBA{0xff, 0xff, 0xff, 0x66},
ColorNameScrollBar:       color.NRGBA{0x0, 0x0, 0x0, 0x99},
ColorNameShadow:          color.NRGBA{0x0, 0x0, 0x0, 0x66},
*/
/* ----- light -----
ColorNameBackground:      color.NRGBA{0xff, 0xff, 0xff, 0xff},
ColorNameButton:          color.Transparent,
ColorNameDisabled:        color.NRGBA{0x0, 0x0, 0x0, 0x42},
ColorNameDisabledButton:  color.NRGBA{0xe5, 0xe5, 0xe5, 0xff},
ColorNameError:           errorColor,
ColorNameForeground:      color.NRGBA{0x21, 0x21, 0x21, 0xff},
ColorNameHover:           color.NRGBA{0x0, 0x0, 0x0, 0x0f},
ColorNameInputBackground: color.NRGBA{0x0, 0x0, 0x0, 0x19},
ColorNamePlaceHolder:     color.NRGBA{0x88, 0x88, 0x88, 0xff},
ColorNamePressed:         color.NRGBA{0x0, 0x0, 0x0, 0x19},
ColorNameScrollBar:       color.NRGBA{0x0, 0x0, 0x0, 0x99},
ColorNameShadow:          color.NRGBA{0x0, 0x0, 0x0, 0x33},
*/
