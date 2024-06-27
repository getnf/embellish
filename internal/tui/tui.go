package tui

import (
	"database/sql"

	"github.com/charmbracelet/bubbles/key"
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

func ThemeGetnfUninstall() *huh.Theme {
	t := ThemeGetnfInstall()

	t.Focused.MultiSelectSelector = t.Focused.MultiSelectSelector.Foreground(lipgloss.Color("1"))
	t.Focused.SelectedOption = t.Focused.SelectedOption.Foreground(lipgloss.Color("1"))
	t.Focused.SelectedPrefix = t.Focused.SelectedPrefix.Foreground(lipgloss.Color("1"))

	return t
}

func myKeyBinds(submitMessage string) *huh.KeyMap {
	var binding huh.KeyMap

	binding.Quit = key.NewBinding(key.WithKeys("ctrl+c", "q"), key.WithHelp("ctrl+c / q", "quit"))
	binding.MultiSelect = huh.MultiSelectKeyMap{
		Up:          key.NewBinding(key.WithKeys("up", "k", "ctrl+p"), key.WithHelp("k / ↑ / C-p:", "Previous")),
		Down:        key.NewBinding(key.WithKeys("down", "j", "ctrl+n"), key.WithHelp("j / ↓ / C-n:", "Next")),
		GotoTop:     key.NewBinding(key.WithKeys("home"), key.WithHelp("Home:", "Go to the top")),
		GotoBottom:  key.NewBinding(key.WithKeys("end"), key.WithHelp("End:", "Go to the bottom")),
		Toggle:      key.NewBinding(key.WithKeys("tab"), key.WithHelp("Tab:", "Toggle")),
		Filter:      key.NewBinding(key.WithKeys(" ", "/"), key.WithHelp("space / /:", "Filter")),
		SetFilter:   key.NewBinding(key.WithKeys("enter"), key.WithHelp("⏎ :", "Set filter"), key.WithDisabled()),
		ClearFilter: key.NewBinding(key.WithKeys("esc"), key.WithHelp("Esc:", "Clear filter"), key.WithDisabled()),
		Submit:      key.NewBinding(key.WithKeys("enter"), key.WithHelp("⏎ :", submitMessage)),
	}

	return &binding
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

	form := huh.NewForm(
		huh.NewGroup(
			ms,
		),
	).WithTheme(
		ThemeGetnfInstall(),
	).WithKeyMap(myKeyBinds("Install fonts"))

	form.Run()

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

	form := huh.NewForm(
		huh.NewGroup(
			ms,
		),
	).WithTheme(
		ThemeGetnfUninstall(),
	).WithKeyMap(myKeyBinds("Uninstall fonts"))

	form.Run()

	for _, font := range selectedFonts {
		handlers.UninstallFont(font, extractPath)
		db.DeleteInstalledFont(database, font)
	}
}
