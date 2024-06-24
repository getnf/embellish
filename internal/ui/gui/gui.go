package gui

import (
	"os"

	"github.com/getnf/getnf/internal/db"
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
	dialog.Show()
}

func myButton(style string, iconName string, tooltip string, sensitive bool) (*gtk.Button, *gtk.Spinner, *gtk.Image) {
	button := gtk.NewButton()
	button.AddCSSClass(style)
	button.SetSensitive(sensitive)
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

func activate(app *gtk.Application, params types.GuiParams) {
	builder := getBuilder(ressources.MainUI)
	window := builder.GetObject("main-window").Cast().(*adw.ApplicationWindow)
	toastOverlay := builder.GetObject("toast-overlay").Cast().(*adw.ToastOverlay)

	handleMainMenuActions(window)

	handleUpdateButton(builder, toastOverlay, params)

	HandleFontsList(builder, params, toastOverlay)

	app.AddWindow(&window.Window)
	window.Show()
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

func handleUpdateButton(builder *gtk.Builder, toastOverlay *adw.ToastOverlay, params types.GuiParams) {
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
		handlers.HandleUpdate(params.Args, params.Database, params.Data, params.DownloadPath, params.ExtractPath)

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

func HandleFontsList(builder *gtk.Builder, params types.GuiParams, toastOverlay *adw.ToastOverlay) {
	fontsList := builder.GetObject("fonts_list").Cast().(*gtk.ListBox)

	nerdFonts := params.Data.GetFonts()

	for _, font := range nerdFonts {
		row := fontRow(font.Name, fontRowButtons(font, params, toastOverlay))
		fontsList.Append(row)
	}
}

func fontRowButtons(font types.Font, params types.GuiParams, toastOverlay *adw.ToastOverlay) *gtk.Box {
	box := gtk.NewBox(0, 10)

	installButton, removeButton := handleFontRowButtons(toastOverlay, params, font)

	box.Append(installButton)
	box.Append(removeButton)
	return box
}

func fontRow(subtitle string, suffix *gtk.Box) *adw.ActionRow {
	row := adw.NewActionRow()
	row.AddCSSClass("property")
	row.SetTitle("Font")
	row.SetSubtitle(subtitle)
	row.AddSuffix(suffix)

	return row
}

func handleFontRowButtons(toastOverlay *adw.ToastOverlay, params types.GuiParams, font types.Font) (*gtk.Button, *gtk.Button) {
	installButton, installSpinner, InstallIcon := myButton("suggested-action", "folder-download-symbolic", "Install", true)
	removeButton, removeSpinner, removeIcon := myButton("destructive-action", "user-trash-symbolic", "Remove", false)

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

func handleInstallButtonAction(font types.Font, installButton *gtk.Button, removeButton *gtk.Button, spinner *gtk.Spinner, icon *gtk.Image, toastOverlay *adw.ToastOverlay, params types.GuiParams) {
	glib.IdleAdd(func() bool {
		icon.SetVisible(false)
		spinner.SetVisible(true)
		spinner.Start()
		return false
	})

	go func() {
		handlers.HandleGuiInstall(font, params.Database, params.Data, params.DownloadPath, params.ExtractPath)

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
		handlers.HandleGuiUninstall(font, params.Database, params.Data, params.ExtractPath)

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
