package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/getnf/getnf/internal/db"
	"github.com/getnf/getnf/internal/handlers"
	"github.com/getnf/getnf/internal/types"

	"github.com/alexflint/go-arg"
)

func setupDB(database *sql.DB, remoteData types.NerdFonts) {
	db.CreateVersionTable(database)
	db.CreateFontsTable(database)
	db.CreateInstalledFontsTable(database)

	if db.TableIsEmpty(database, "version") || handlers.IsUpdateAvilable(remoteData.TagName, db.GetVersion(database)) {
		db.InsertIntoVersion(database, remoteData.TagName)
		fmt.Println("Updated fonts version")
	}

	if db.TableIsEmpty(database, "fonts") || handlers.IsUpdateAvilable(remoteData.TagName, db.GetVersion(database)) {
		db.DeleteFontsTable(database)
		db.CreateFontsTable(database)
		db.InsertIntoFonts(database, remoteData.GetFonts())
		fmt.Println("Updating local fonts db")
	}
}

func main() {
	var args types.Args
	arg.MustParse(&args)

	remoteData, err := handlers.GetData()
	if err != nil {
		log.Fatalln(err)
	}

	var data types.NerdFonts

	database := db.OpenDB()

	setupDB(database, remoteData)

	data.TagName = db.GetVersion(database)
	data.Assets = db.GetAllFonts(database)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln(err)
	}

	downloadPath := homeDir + "/.local/share/getnf"
	extractPath := homeDir + "/.local/share/fonts"

	switch {
	case args.List != nil:
		if !args.List.Installed {
			handlers.ListFonts(handlers.FontsWithVersion(database, data.GetFonts(), data.GetVersion()), false)
		} else {
			handlers.ListFonts(handlers.FontsWithVersion(database, data.GetFonts(), data.GetVersion()), true)
		}
	case args.Install != nil:
		for _, font := range args.Install.Fonts {
			if db.FontExists(database, font) {
				f := data.GetFont(font)
				handlers.InstallFont(f, downloadPath, extractPath, args.Install.KeepTars)
				db.InsertIntoInstalledFonts(database, f, data.GetVersion())
			} else {
				fmt.Printf("%v is not a nerd font", font)
			}
		}
	case args.Uninstall != nil:
		for _, font := range args.Uninstall.Fonts {
			if db.IsFontInstalled(database, font) {
				handlers.UninstallFont(extractPath, font)
				db.DeleteInstalledFont(database, font)
			} else {
				fmt.Printf("%v is either not installed or is not a nerd font", font)
			}
		}
	case args.Update.Update:
		updateCount := 0
		for _, font := range db.GetInstalledFonts(database) {
			if handlers.IsUpdateAvilable(data.GetVersion(), font.InstalledVersion) {
				f := data.GetFont(font.Name)
				handlers.InstallFont(f, downloadPath, extractPath, args.Update.KeepTars)
				db.UpdateInstalledFont(database, font.Name, data.GetVersion())
				updateCount++
			}
		}
		if updateCount > 0 {
			if updateCount > 1 {
				fmt.Printf("%d fonts were updated\n", updateCount)
			} else {
				fmt.Printf("%d font was updated\n", updateCount)
			}
		} else {
			fmt.Println("no updates are available")
		}
	}
}
