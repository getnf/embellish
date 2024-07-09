//go:build gui
// +build gui

package gui

import "github.com/diamondburned/gotk4/pkg/gtk/v4"

func getBuilder(file string) *gtk.Builder {
	return gtk.NewBuilderFromString(file, len(file))

}

func createButton(style string, iconName string, tooltip string, sensitive bool) (*gtk.Button, *gtk.Spinner, *gtk.Image) {
	button := gtk.NewButton()
	button.AddCSSClass(style)
	button.SetSensitive(sensitive)
	button.SetVAlign(3)
	button.SetTooltipText(tooltip)

	spinner := gtk.NewSpinner()
	spinner.SetVisible(false)

	icon := gtk.NewImage()
	icon.SetFromIconName(iconName)

	box := gtk.NewBox(0, 0)

	box.Append(spinner)
	box.Append(icon)

	button.SetChild(box)

	return button, spinner, icon
}
