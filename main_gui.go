//go:build gui
// +build gui

package main

import (
	"database/sql"
	"os"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/getnf/embellish/internal/db"
	"github.com/getnf/embellish/internal/gui"
	ressources "github.com/getnf/embellish/internal/gui/resources"
	"github.com/getnf/embellish/internal/handlers"
	"github.com/getnf/embellish/internal/types"
)

func main() {
	var args types.Args
	arg.MustParse(&args)

	var database *sql.DB

	paths := types.NewPaths()
	downloadPath := paths.GetDownloadPath()
	extractPath := paths.GetInstallPath()
	dbPath := paths.GetDbPath()

	database = db.OpenDB(dbPath)

	db.CreateLastCheckedTable(database)

	lastChecked, _ := time.Parse(time.DateTime, db.GetLastChecked(database))
	DaysSinceLastChecked := int(time.Since(lastChecked).Hours() / 24)

	if db.TableIsEmpty(database, "lastChecked") || DaysSinceLastChecked > 5 || args.ForceCheck {
		remoteData, err := handlers.GetData()
		if err == nil {
			handlers.SetupDB(database, remoteData)
		}
		db.UpdateLastChecked(database)
	}

	var data types.NerdFonts

	data.Version = db.GetVersion(database)
	data.Fonts = db.GetAllFonts(database)

	params := types.GuiParams{Data: data, Database: database, DownloadPath: downloadPath, ExtractPath: extractPath}
	app := adw.NewApplication("io.github.getnf.embellish", 0)

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
	builder := gui.GetBuilder(ressources.MainUI)
	window := builder.GetObject("main-window").Cast().(*adw.ApplicationWindow)
	toastOverlay := builder.GetObject("toast-overlay").Cast().(*adw.ToastOverlay)
	mainPage := builder.GetObject("main-page").Cast().(*adw.StatusPage)

	mainPage.SetDescription("Install nerd font\nVersion: " + params.Data.Version)

	gui.SetupActions(app, builder)

	adwStyle := adw.StyleManagerGetDefault()
	adwStyle.SetColorScheme(adw.ColorSchemePreferLight)

	gui.HandleMainMenuActions(window, params)

	gui.HandleUpdateButton(builder, toastOverlay, params)

	gui.HandleFontsList(builder, params, toastOverlay)

	gui.HandleFontsSearch(builder, params, toastOverlay)

	app.AddWindow(&window.Window)
	window.Show()
}
