//go:build gui
// +build gui

package gui

import "github.com/diamondburned/gotk4/pkg/gtk/v4"

func GetBuilder(file string) *gtk.Builder {
	return gtk.NewBuilderFromString(file)

}

func createButton(style string, iconName string, tooltip string, visibility bool) (*gtk.Button, *gtk.Spinner, *gtk.Image) {
	button := gtk.NewButton()
	button.AddCSSClass(style)
	button.SetVisible(visibility)
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

func clearListBox(list *gtk.ListBox) {
	for {
		row := list.RowAtIndex(0)
		if row == nil {
			break
		}
		list.Remove(row)
	}
}
