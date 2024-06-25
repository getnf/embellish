package gui

import (
	"os"

	"github.com/getnf/getnf/internal/types"
	ressources "github.com/getnf/getnf/internal/ui/gui/resources"

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

	handleMainMenuActions(window)

	handleUpdateButton(builder, toastOverlay, params)

	HandleFontsList(builder, params, toastOverlay)

	handleFontsSearch(builder, params, toastOverlay)

	app.AddWindow(&window.Window)
	window.Show()
}
