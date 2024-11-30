//go:generate fyne bundle -o bundled.go assets

package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type fysionTheme struct {
	fyne.Theme
}

func myCustomTheme() fyne.Theme {
	return &fysionTheme{Theme: theme.DefaultTheme()}
}

func (t *fysionTheme) Color(name fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	return t.Theme.Color(name, theme.VariantLight)
}

func (t *fysionTheme) Size(name fyne.ThemeSizeName) float32 {
	if name == theme.SizeNameText {
		return 10
	}
	return t.Theme.Size(name)
}
