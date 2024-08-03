//go:build gui
// +build gui

package gui

import (
	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/core/glib"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	ressources "github.com/getnf/embellish/internal/gui/resources"
	"github.com/getnf/embellish/internal/handlers"
	"github.com/getnf/embellish/internal/types"
)

func openAboutDialog(window *adw.ApplicationWindow) {
	builder := GetBuilder(ressources.AboutUI)
	dialog := builder.GetObject("about-dialog").Cast().(*adw.AboutDialog)
	dialog.Present(window)
}

func HandleMainMenuActions(window *adw.ApplicationWindow, params types.GuiParams) {
	appGroup := gio.NewSimpleActionGroup()
	window.InsertActionGroup("app", appGroup)

	about_action := gio.NewSimpleAction("about", nil)
	about_action.Connect("activate", func() {
		openAboutDialog(window)
	})

	update_action := gio.NewSimpleAction("update", nil)
	update_action.Connect("activate", func() {
		handlers.SetupDB(params.Database, params.Data)
	})

	appGroup.AddAction(about_action)
	appGroup.AddAction(update_action)
}

func HandleUpdateButton(builder *gtk.Builder, toastOverlay *adw.ToastOverlay, params types.GuiParams) {
	updateButton := builder.GetObject("update-button").Cast().(*gtk.Button)

	updateButtonSpinner := builder.GetObject("update-button-spinner").Cast().(*gtk.Spinner)

	updateButtonIcon := builder.GetObject("update-button-icon").Cast().(*gtk.Image)

	if handlers.IsFontUpdatAvilable(params.Database, params.Data) {
		updateButton.SetVisible(true)
	}

	updateButton.ConnectClicked(
		func() {
			handleUpdateButtonAction(updateButton, updateButtonSpinner, updateButtonIcon, toastOverlay, params)
		})
}

func handleUpdateButtonAction(button *gtk.Button, spinner *gtk.Spinner, icon *gtk.Image, toastOverlay *adw.ToastOverlay, params types.GuiParams) {
	glib.IdleAdd(func() bool {
		icon.SetVisible(false)
		spinner.SetVisible(true)
		spinner.Start()
		return false
	})

	go func() {
		handlers.HandleUpdate(params.Database, params.Data, params.DownloadPath, params.ExtractPath)

		glib.IdleAdd(func() bool {
			spinner.Stop()
			spinner.SetVisible(false)
			return false
		})

		glib.IdleAdd(func() bool {
			toastOverlay.AddToast(adw.NewToast("updated completed"))
			return false
		})

		glib.IdleAdd(func() bool {
			button.SetVisible(false)
			return false
		})
	}()
}
