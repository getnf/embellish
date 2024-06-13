package gui

import (
	"os"

	"github.com/getnf/getnf/internal/handlers"
	"github.com/getnf/getnf/internal/types"
	ressources "github.com/getnf/getnf/internal/ui/gui/resources"

	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/core/glib"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func RunGui(params types.GuiParams) {
	app := gtk.NewApplication("com.github.diamondburned.gotk4-examples.gtk4.simple", 0)
	app.ConnectActivate(func() {
		activate(app, params)
	})

	if code := app.Run(os.Args); code > 0 {
		os.Exit(code)
	}
}

func getBuilder(file string) *gtk.Builder {
	return gtk.NewBuilderFromString(file, len(file))

}

func openAboutDialog() {
	builder := getBuilder(ressources.AboutUI)
	dialog := builder.GetObject("about-dialog").Cast().(*adw.AboutWindow)
	dialog.Present()
}

func activate(app *gtk.Application, params types.GuiParams) {
	builder := getBuilder(ressources.MainUI)
	window := builder.GetObject("main-window").Cast().(*adw.ApplicationWindow)

	handleMainMenuActions(window)

	handleUpdateButton(builder, params)

	handleFontsList(builder, params.Data)

	app.AddWindow(&window.Window)
	window.Show()
}

func handleFontsList(builder *gtk.Builder, data types.NerdFonts) {
	// the list of fonts
	listView := builder.GetObject("list-view").Cast().(*gtk.ListView)
	stringModel := gtk.NewStringList(data.GetFontsNames())
	model := gtk.NewSingleSelection(stringModel)
	listView.SetModel(model)
}

func handleMainMenuActions(window *adw.ApplicationWindow) {
	appGroup := gio.NewSimpleActionGroup()
	window.InsertActionGroup("app", appGroup)
	about_action := gio.NewSimpleAction("about", nil)
	about_action.Connect("activate", func() {
		openAboutDialog()
	})
	appGroup.AddAction(about_action)
}

func handleUpdateButton(builder *gtk.Builder, params types.GuiParams) {
	updateButton := builder.GetObject("update-button").Cast().(*gtk.Button)

	updateButtonSpinner := builder.GetObject("update-button-spinner").Cast().(*gtk.Spinner)

	updateButtonIcon := builder.GetObject("update-button-icon").Cast().(*gtk.Image)

	if handlers.IsFontUpdatAvilable(params.Database, params.Data) {
		updateButton.SetVisible(true)
	}

	updateButton.ConnectClicked(
		func() {
			handleUpdateButtonAction(updateButtonSpinner, updateButton, updateButtonIcon, params)
		})
}

func handleUpdateButtonAction(spinner *gtk.Spinner, button *gtk.Button, icon *gtk.Image, params types.GuiParams) {
	glib.IdleAdd(func() bool {
		icon.SetVisible(false)
		spinner.SetVisible(true)
		spinner.Start()
		return false
	})

	go func() {
		handlers.HandleUpdate(params.Args, params.Database, params.Data, params.DownloadPath, params.ExtractPath)

		glib.IdleAdd(func() bool {
			spinner.Stop()
			spinner.SetVisible(false)
			button.SetVisible(false)
			return false
		})
	}()
}
