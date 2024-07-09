//go:build gui
// +build gui

package handlers

import (
	"database/sql"

	"github.com/getnf/embellish/internal/db"
	"github.com/getnf/embellish/internal/types"
)

func HandleInstall(font types.Font, database *sql.DB, data types.NerdFonts, downloadPath string, extractPath string) error {
	err := PlatformInstallFont(font, downloadPath, extractPath, false)
	if err != nil {
		return err
	}
	db.InsertIntoInstalledFonts(database, font, data.GetVersion())
	return nil
}

func HandleUninstall(font types.Font, database *sql.DB, extractPath string) error {
	err := PlatformUninstallFont(extractPath, font.Name)
	if err != nil {
		return err
	}
	db.DeleteInstalledFont(database, font.Name)
	return nil
}
