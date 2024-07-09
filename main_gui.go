//go:build gui
// +build gui

package main

import (
	"database/sql"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/getnf/embellish/internal/db"
	"github.com/getnf/embellish/internal/gui"
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

	gui.RunGui(types.GuiParams{Data: data, Database: database, Args: args, DownloadPath: downloadPath, ExtractPath: extractPath})
}
