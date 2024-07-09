//go:build gui
// +build gui

package gui

import (
	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/core/glib"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/getnf/embellish/internal/db"
	"github.com/getnf/embellish/internal/handlers"
	"github.com/getnf/embellish/internal/types"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

func HandleFontsList(builder *gtk.Builder, params types.GuiParams, toastOverlay *adw.ToastOverlay) {
	fontsList := builder.GetObject("fonts-list").Cast().(*gtk.ListBox)

	populateFontsList(fontsList, params, toastOverlay)
}

func handleFontsSearch(builder *gtk.Builder, params types.GuiParams, toastOverlay *adw.ToastOverlay) {
	searchButton := builder.GetObject("search-button").Cast().(*gtk.ToggleButton)
	searchBar := builder.GetObject("search-bar").Cast().(*gtk.SearchBar)
	searchEntry := builder.GetObject("search-entry").Cast().(*gtk.SearchEntry)
	searchList := builder.GetObject("search-list").Cast().(*gtk.ListBox)
	stack := builder.GetObject("stack").Cast().(*gtk.Stack)
	searchPage := builder.GetObject("search-page").Cast().(*gtk.ScrolledWindow)
	mainPage := builder.GetObject("main-page").Cast().(*adw.StatusPage)
	statusPage := builder.GetObject("status-page").Cast().(*adw.StatusPage)

	searchButton.ConnectClicked(func() {
		searchBar.SetSearchMode(!searchBar.SearchMode())
	})

	searchBar.Connect("notify::search-mode-enabled", func() {
		if searchBar.SearchMode() {
			stack.SetVisibleChild(searchPage)
		} else {
			stack.SetVisibleChild(mainPage)
		}
	})

	populateFontsList(searchList, params, toastOverlay)

	var resultsCount int

	searchList.SetFilterFunc(func(row *gtk.ListBoxRow) (ok bool) {
		match := fuzzy.MatchFold(searchEntry.Text(), row.Name())
		if match {
			resultsCount++
		}
		return match
	})

	searchEntry.Connect("search-changed", func() {
		resultsCount = 0
		searchList.InvalidateFilter()
		if resultsCount == 0 {
			stack.SetVisibleChild(statusPage)
		} else if searchBar.SearchMode() {
			stack.SetVisibleChild(searchPage)
		}
	})
}

func populateFontsList(list *gtk.ListBox, params types.GuiParams, toastOverlay *adw.ToastOverlay) {
	nerdFonts := params.Data.GetFonts()
	for _, font := range nerdFonts {
		row := createFontRow(font.Name, createFontRowSuffix(font, params, toastOverlay))
		list.Append(row)
	}
}

func createFontRowButtons(toastOverlay *adw.ToastOverlay, params types.GuiParams, font types.Font) (*gtk.Button, *gtk.Button) {
	installButton, installSpinner, InstallIcon := createButton("suggested-action", "folder-download-symbolic", "Install", true)
	removeButton, removeSpinner, removeIcon := createButton("destructive-action", "user-trash-symbolic", "Remove", false)

	if db.IsFontInstalled(params.Database, font.Name) {
		installButton.SetSensitive(false)
	}
	installButton.ConnectClicked(func() {
		handleInstallButtonAction(font, installButton, removeButton, installSpinner, InstallIcon, toastOverlay, params)
	})

	if db.IsFontInstalled(params.Database, font.Name) {
		removeButton.SetSensitive(true)
	}
	removeButton.ConnectClicked(func() {
		handleRemoveButtonAction(font, removeButton, installButton, removeSpinner, removeIcon, toastOverlay, params)
	})

	return installButton, removeButton
}

func createFontRowSuffix(font types.Font, params types.GuiParams, toastOverlay *adw.ToastOverlay) *gtk.Box {
	box := gtk.NewBox(0, 10)

	installButton, removeButton := createFontRowButtons(toastOverlay, params, font)

	box.Append(installButton)
	box.Append(removeButton)
	return box
}

func createFontRow(subtitle string, suffix *gtk.Box) *adw.ActionRow {
	row := adw.NewActionRow()
	row.SetName(subtitle)
	row.AddCSSClass("property")
	row.SetTitle("Font")
	row.SetSubtitle(subtitle)
	row.AddSuffix(suffix)

	return row
}

func handleInstallButtonAction(font types.Font, installButton *gtk.Button, removeButton *gtk.Button, spinner *gtk.Spinner, icon *gtk.Image, toastOverlay *adw.ToastOverlay, params types.GuiParams) {
	glib.IdleAdd(func() bool {
		icon.SetVisible(false)
		spinner.SetVisible(true)
		spinner.Start()
		return false
	})

	go func() {
		err := handlers.HandleInstall(font, params.Database, params.Data, params.DownloadPath, params.ExtractPath)

		if err != nil {
			glib.IdleAdd(func() bool {
				spinner.Stop()
				spinner.SetVisible(false)
				icon.SetVisible(true)
				return false
			})

			glib.IdleAdd(func() bool {
				toastOverlay.AddToast(adw.NewToast(err.Error()))
				return false
			})

			return
		}

		glib.IdleAdd(func() bool {
			spinner.Stop()
			spinner.SetVisible(false)
			icon.SetVisible(true)
			return false
		})

		glib.IdleAdd(func() bool {
			toastOverlay.AddToast(adw.NewToast("Install completed"))
			return false
		})

		glib.IdleAdd(func() bool {
			installButton.SetSensitive(false)
			removeButton.SetSensitive(true)
			return false
		})
	}()
}

func handleRemoveButtonAction(font types.Font, removeButton *gtk.Button, installButton *gtk.Button, spinner *gtk.Spinner, icon *gtk.Image, toastOverlay *adw.ToastOverlay, params types.GuiParams) {
	glib.IdleAdd(func() bool {
		icon.SetVisible(false)
		spinner.SetVisible(true)
		spinner.Start()
		return false
	})

	go func() {
		handlers.HandleUninstall(font, params.Database, params.ExtractPath)

		glib.IdleAdd(func() bool {
			spinner.Stop()
			spinner.SetVisible(false)
			icon.SetVisible(true)
			return false
		})

		glib.IdleAdd(func() bool {
			toastOverlay.AddToast(adw.NewToast("Uninstall completed"))
			return false
		})

		glib.IdleAdd(func() bool {
			removeButton.SetSensitive(false)
			installButton.SetSensitive(true)
			return false
		})
	}()
}
