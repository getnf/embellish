package gui

import (
	"os"

	ressources "github.com/getnf/embellish/internal/gui/resources"
	"github.com/getnf/embellish/internal/types"

	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func RunGui(params types.GuiParams) {
	app := adw.NewApplication("com.github.getnf.Embellish", 0)

	provider := gtk.NewCSSProvider()
	provider.LoadFromData(ressources.StylesCSS)

	app.ConnectActivate(func() {
		gtk.StyleContextAddProviderForDisplay(
			gdk.DisplayGetDefault(),
			provider,
			gtk.STYLE_PROVIDER_PRIORITY_APPLICATION,
		)
		activate(app, params)
	})

	if code := app.Run(os.Args); code > 0 {
		os.Exit(code)
	}
}

func activate(app *adw.Application, params types.GuiParams) {
	builder := getBuilder(ressources.MainUI)
	window := builder.GetObject("main-window").Cast().(*adw.ApplicationWindow)
	toastOverlay := builder.GetObject("toast-overlay").Cast().(*adw.ToastOverlay)

	setupActions(app, builder)

	adwStyle := adw.StyleManagerGetDefault()
	adwStyle.SetColorScheme(adw.ColorSchemePreferLight)

	handleMainMenuActions(window, params)

	handleUpdateButton(builder, toastOverlay, params)

	HandleFontsList(builder, params, toastOverlay)

	handleFontsSearch(builder, params, toastOverlay)

	app.AddWindow(&window.Window)
	window.Show()
}
