package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/getnf/getnf/internal/db"
	"github.com/getnf/getnf/internal/handlers"
	"github.com/getnf/getnf/internal/types"
	"github.com/getnf/getnf/internal/ui/gui"
	"github.com/getnf/getnf/internal/ui/tui"
	"github.com/getnf/getnf/internal/utils"

	"github.com/alexflint/go-arg"
)

func setupDB(database *sql.DB, remoteData types.NerdFonts) {
	db.CreateVersionTable(database)
	db.CreateFontsTable(database)
	db.CreateInstalledFontsTable(database)

	if db.TableIsEmpty(database, "version") || handlers.IsUpdateAvilable(remoteData.GetVersion(), db.GetVersion(database)) {
		db.InsertIntoVersion(database, remoteData.GetVersion())
		fmt.Println("Updated fonts version")
	}

	if db.TableIsEmpty(database, "fonts") || handlers.IsUpdateAvilable(remoteData.GetVersion(), db.GetVersion(database)) {
		db.DeleteFontsTable(database)
		db.CreateFontsTable(database)
		db.InsertIntoFonts(database, remoteData.GetFonts())
		fmt.Println("Updating local fonts db")
	}
}

func main() {
	var args types.Args
	arg.MustParse(&args)

	isAdmin := handlers.IsAdmin()

	var database *sql.DB

	if utils.OsType() == "windows" && !isAdmin {
		log.Fatalln("getnf need admin rights to install fonts on windows, please run getnf as administrator")
	}

	if isAdmin && utils.OsType() != "windows" {
		log.Fatalln("Please don't run getnf with elevated privileges")
	} else {
		database = db.OpenDB()
	}

	db.CreateLastCheckedTable(database)

	lastChecked, _ := time.Parse(time.DateTime, db.GetLastChecked(database))
	DaysSinceLastChecked := int(time.Since(lastChecked).Hours() / 24)

	if db.TableIsEmpty(database, "lastChecked") || DaysSinceLastChecked > 5 || args.ForceCheck {
		remoteData, err := handlers.GetData()
		if err == nil {
			setupDB(database, remoteData)
		}
		db.UpdateLastChecked(database)
	}

	var data types.NerdFonts

	data.Version = db.GetVersion(database)
	data.Fonts = db.GetAllFonts(database)

	paths := types.NewPaths()
	downloadPath := paths.GetDownloadPath()
	extractPath := paths.GetInstallPath()

	switch {
	case args.List != nil:
		if !args.List.Installed {
			handlers.ListFonts(handlers.FontsWithVersion(database, data.GetFonts(), data.GetVersion()), false)
		} else {
			handlers.ListFonts(handlers.FontsWithVersion(database, data.GetFonts(), data.GetVersion()), true)
		}
	case args.Install != nil:
		if len(args.Install.Fonts) == 0 {
			tui.SelectFontsToInstall(data, database, downloadPath, extractPath, args.KeepTars)
		} else {
			handlers.HandleInstall(args, database, data, downloadPath, extractPath)
		}
	case args.Uninstall != nil:
		if len(args.Uninstall.Fonts) == 0 {
			tui.SelectFontsToUninstall(db.GetInstalledFonts(database), database, extractPath)
		} else {
			handlers.HandleUninstall(args, database, data, extractPath)
		}
	case args.Update != nil:
		handlers.HandleUpdate(args, database, data, downloadPath, extractPath)
	default:
		gui.RunGui(types.GuiParams{Data: data, Database: database, Args: args, DownloadPath: downloadPath, ExtractPath: extractPath})
	}
}
