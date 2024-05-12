package tui

import (
	"database/sql"

	"github.com/charmbracelet/huh"
	"github.com/getnf/getnf/internal/db"
	"github.com/getnf/getnf/internal/handlers"
	"github.com/getnf/getnf/internal/types"
	"github.com/getnf/getnf/internal/utils"
)

func SelectFontsToInstall(data types.NerdFonts, database *sql.DB, downloadPath string, extractPath string, keepTar bool) {
	var selectedFontsNames []string
	var selectedFonts []types.Font
	fontsNames := data.GetFontsNames()
	optionsFromFonts := huh.NewOptions(fontsNames...)

	ms := huh.NewMultiSelect[string]().
		Options(
			optionsFromFonts...,
		).
		Title("Select fonts to install").
		Value(&selectedFontsNames).
		Filterable(true)
	ms.Run()

	for _, fontName := range selectedFontsNames {
		selectedFontName := data.GetFont(fontName)
		selectedFonts = append(selectedFonts, selectedFontName)
	}

	for _, font := range selectedFonts {
		handlers.InstallFont(font, downloadPath, extractPath, keepTar)
		db.InsertIntoInstalledFonts(database, font, data.GetVersion())
	}
}

func SelectFontsToUninstall(installedFonts []types.Font, database *sql.DB, extractPath string) {
	var selectedFonts []string
	installedFontsNames := utils.Fold(installedFonts, func(f types.Font) string {
		return f.Name
	})
	optionsFromInstalledFonts := huh.NewOptions(installedFontsNames...)

	ms := huh.NewMultiSelect[string]().
		Options(
			optionsFromInstalledFonts...,
		).
		Title("Select fonts to uninstall").
		Value(&selectedFonts).
		Filterable(true)
	ms.Run()

	for _, font := range selectedFonts {
		handlers.UninstallFont(font, extractPath)
		db.DeleteInstalledFont(database, font)
	}
}
