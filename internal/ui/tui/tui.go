package tui

import (
	"database/sql"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/getnf/getnf/internal/db"
	"github.com/getnf/getnf/internal/handlers"
	"github.com/getnf/getnf/internal/types"
	"github.com/getnf/getnf/internal/utils"
)

func ThemeGetnfInstall() *huh.Theme {
	t := huh.ThemeBase()

	t.Focused.Base = t.Focused.Base.BorderForeground(lipgloss.Color("7"))
	t.Focused.Title = t.Focused.Title.Foreground(lipgloss.Color("3"))
	t.Focused.SelectSelector = t.Focused.SelectSelector.Foreground(lipgloss.Color("3"))
	t.Focused.MultiSelectSelector = t.Focused.MultiSelectSelector.Foreground(lipgloss.Color("2"))
	t.Focused.SelectedOption = t.Focused.SelectedOption.Foreground(lipgloss.Color("2"))
	t.Focused.SelectedPrefix = t.Focused.SelectedPrefix.Foreground(lipgloss.Color("2"))
	t.Focused.UnselectedOption = t.Focused.UnselectedOption.Foreground(lipgloss.Color("0"))

	return t
}

func ThemeGetnfUnInstall() *huh.Theme {
	t := ThemeGetnfInstall()

	t.Focused.MultiSelectSelector = t.Focused.MultiSelectSelector.Foreground(lipgloss.Color("1"))
	t.Focused.SelectedOption = t.Focused.SelectedOption.Foreground(lipgloss.Color("1"))
	t.Focused.SelectedPrefix = t.Focused.SelectedPrefix.Foreground(lipgloss.Color("1"))

	return t
}

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
	ms.WithTheme(ThemeGetnfInstall()).Run()

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
	ms.WithTheme(ThemeGetnfUnInstall()).Run()

	for _, font := range selectedFonts {
		handlers.UninstallFont(font, extractPath)
		db.DeleteInstalledFont(database, font)
	}
}
