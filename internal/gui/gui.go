package gui

import (
	"os"

	ressources "github.com/getnf/getnf/internal/gui/resources"
	"github.com/getnf/getnf/internal/types"

	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
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
	toastOverlay := builder.GetObject("toast-overlay").Cast().(*adw.ToastOverlay)

	adwStyle := adw.StyleManagerGetDefault()
	adwStyle.SetColorScheme(adw.ColorSchemePreferLight)

	handleMainMenuActions(window, params)

	handleUpdateButton(builder, toastOverlay, params)

	HandleFontsList(builder, params, toastOverlay)

	handleFontsSearch(builder, params, toastOverlay)

	app.AddWindow(&window.Window)
	window.Show()
}
