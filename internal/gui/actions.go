package gui

import (
	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func setupActions(app *adw.Application, builder *gtk.Builder) {
	searchBar := builder.GetObject("search-bar").Cast().(*gtk.SearchBar)

	quitAction := gio.NewSimpleAction("quit", nil)
	searchAction := gio.NewSimpleAction("search", nil)

	// TODO: fix quiting the app when dialog is open
	quitAction.Connect("activate", func() {
		app.Quit()
	})

	searchAction.Connect("activate", func() {
		searchBar.SetSearchMode(!searchBar.SearchMode())
	})

	app.AddAction(quitAction)
	app.AddAction(searchAction)

	app.SetAccelsForAction("window.close", []string{"<Control>w"})
	app.SetAccelsForAction("app.quit", []string{"<Control>q"})
	app.SetAccelsForAction("app.search", []string{"<Control>f"})
}
