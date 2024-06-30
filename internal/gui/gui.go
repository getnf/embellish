package gui

import (
	"os"

	ressources "github.com/getnf/embellish/internal/gui/resources"
	"github.com/getnf/embellish/internal/types"

	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func RunGui(params types.GuiParams) {
	app := gtk.NewApplication("com.github.getnf.getnf", 0)
	app.ConnectActivate(func() {
		activate(app, params)
	})

	if code := app.Run(os.Args); code > 0 {
		os.Exit(code)
	}
}

func activate(app *gtk.Application, params types.GuiParams) {
	builder := getBuilder(ressources.MainUI)
	window := builder.GetObject("main-window").Cast().(*adw.ApplicationWindow)
	aboutBuilder := getBuilder(ressources.AboutUI)
	dialog := aboutBuilder.GetObject("about-dialog").Cast().(*adw.AboutWindow)
	dialog.SetTransientFor(&window.Window)
	dialog.SetDestroyWithParent(true)
	toastOverlay := builder.GetObject("toast-overlay").Cast().(*adw.ToastOverlay)

	quitAction := gio.NewSimpleAction("quit", nil)

	// TODO: fix quiting the app when dialog is open
	quitAction.Connect("activate", func() {
		app.Quit()
	})

	app.AddAction(quitAction)

	app.SetAccelsForAction("window.close", []string{"<Control>w"})
	app.SetAccelsForAction("app.quit", []string{"<Control>q"})

	adwStyle := adw.StyleManagerGetDefault()
	adwStyle.SetColorScheme(adw.ColorSchemePreferLight)

	handleMainMenuActions(window, dialog, params)

	handleUpdateButton(builder, toastOverlay, params)

	HandleFontsList(builder, params, toastOverlay)

	handleFontsSearch(builder, params, toastOverlay)

	app.AddWindow(&window.Window)
	window.Show()
}
