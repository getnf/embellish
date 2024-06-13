package ressources

import (
	_ "embed"
)

var (
	//go:embed main.ui
	MainUI string

	//go:embed about.ui
	AboutUI string

	//go:embed style.css
	StyleCSS string
)
