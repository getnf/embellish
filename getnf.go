package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/getnf/getnf/internal/db"
	"github.com/getnf/getnf/internal/handlers"
	"github.com/getnf/getnf/internal/types"
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
		db.InsertIntoFonts(database, remoteData.Assets)
		fmt.Println("Updating local fonts db")
	}
}

func main() {
	// keepArchives := true

	remoteData, err := handlers.GetData()
	if err != nil {
		log.Fatalln(err)
	}

	var data types.NerdFonts

	database := db.OpenDB()

	setupDB(database, remoteData)

	data.TagName = db.GetVersion(database)
	data.Assets = db.GetAllFonts(database)

	fontsWithExtraInfo := handlers.FontsWithVersion(database, data.GetFonts())
	handlers.ListFonts(fontsWithExtraInfo)
}
