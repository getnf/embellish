//go:build terminal
// +build terminal

package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/getnf/embellish/internal/db"
	"github.com/getnf/embellish/internal/handlers"
	"github.com/getnf/embellish/internal/tui"
	"github.com/getnf/embellish/internal/types"
	"github.com/getnf/embellish/internal/utils"
)

func main() {
	var args types.Args
	arg.MustParse(&args)

	var database *sql.DB

	paths := types.NewPaths()
	downloadPath := paths.GetDownloadPath()
	extractPath := paths.GetInstallPath()
	dbPath := paths.GetDbPath()
	isAdmin := handlers.IsAdmin()

	if utils.OsType() == "windows" && !isAdmin {
		log.Fatalln("embellish need admin rights to install fonts on windows, please run embellish as administrator")
	}

	if isAdmin && utils.OsType() != "windows" {
		log.Fatalln("Please don't run embellish with elevated privileges")
	} else {
		database = db.OpenDB(dbPath)
	}

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

	switch {
	case args.List != nil:
		if !args.List.Installed {
			handlers.ListFonts(handlers.FontsWithVersion(database, data.GetFonts(), data.GetVersion()), false)
		} else {
			handlers.ListFonts(handlers.FontsWithVersion(database, data.GetFonts(), data.GetVersion()), true)
		}
	case args.Install != nil:
		if len(args.Install.Fonts) == 0 {
			err := tui.SelectFontsToInstall(data, database, downloadPath, extractPath, args.KeepTars)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			err := handlers.HandleInstall(args, database, data, downloadPath, extractPath)
			if err != nil {
				fmt.Println(err)
			}
		}
	case args.Uninstall != nil:
		if len(args.Uninstall.Fonts) == 0 {
			err := tui.SelectFontsToUninstall(db.GetInstalledFonts(database), database, extractPath)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			err := handlers.HandleUninstall(args, database, data, extractPath)
			if err != nil {
				fmt.Println(err)
			}
		}
	case args.Update != nil:
		err := handlers.HandleUpdate(database, data, downloadPath, extractPath)
		if err != nil {
			fmt.Println(err)
		}
	default:
		err := tui.SelectFontsToInstall(data, database, downloadPath, extractPath, args.KeepTars)
		if err != nil {
			fmt.Println(err)
		}
	}
}
