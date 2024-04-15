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
	// keepArchives := true

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
			handlers.ListFonts(handlers.FontsWithVersion(database, data.GetFonts()), false)
		} else {
			handlers.ListFonts(handlers.FontsWithVersion(database, data.GetFonts()), true)
		}
	case args.Install != nil:
		for _, font := range args.Install.Fonts {
			if db.FontExists(database, font) {
				f := data.GetFont(font)
				downloadedTar, err := handlers.DownloadTar(f.BrowserDownloadUrl, downloadPath, f.Name)
				if err != nil {
					log.Fatalln(err)
				}
				handlers.ExtractTar(downloadedTar, extractPath, f.Name)
				if !args.Install.KeepArchives {
					handlers.CleanUpArchive(downloadedTar)
				}
			} else {
				fmt.Printf("Font: %v is not a nerd font", font)
			}
		}
	}
}
