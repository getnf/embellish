package db

import (
	"database/sql"
	"log"
	"sort"
	"strings"

	"github.com/getnf/embellish/internal/types"
	"github.com/getnf/embellish/internal/utils"
	_ "modernc.org/sqlite"
)

func OpenDB(path string) *sql.DB {
	db, err := sql.Open("sqlite", path+"/"+"embellish.sqlite3")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func TableIsEmpty(db *sql.DB, table string) bool {
	sqlstmt := "SELECT EXISTS (SELECT 1 FROM " + table + ")"
	statement, _ := db.Query(sqlstmt)
	defer statement.Close()

	var ver int

	for statement.Next() {
		statement.Scan(&ver)
	}

	if ver == 0 {
		return true
	} else {
		return false
	}
}

// last checked table for getting the time from now till the last time we check for updated

func CreateLastCheckedTable(db *sql.DB) {
	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS lastChecked (id INTEGER PRIMARY KEY, date TEXT)")
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer statement.Close()
	_, err = statement.Exec()
	if err != nil {
		log.Fatal(err)
	}
}

func UpdateLastChecked(db *sql.DB) {
	statement, err := db.Prepare("INSERT or REPLACE INTO lastChecked (id, date) VALUES (?, DateTime('now'))")
	if err != nil {
		log.Fatalln(err)
	}
	defer statement.Close()

	_, err = statement.Exec(1)
	if err != nil {
		log.Fatalln(err)
	}
}

func GetLastChecked(db *sql.DB) string {
	rows, err := db.Query("SELECT date FROM lastChecked")
	if err != nil {
		log.Fatalln(err)
	}
	var date string
	for rows.Next() {
		rows.Scan(&date)
	}
	return date
}

// Version table

func CreateVersionTable(db *sql.DB) {
	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS version (id INTEGER PRIMARY KEY, Version TEXT)")
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer statement.Close()
	_, err = statement.Exec()
	if err != nil {
		log.Fatal(err)
	}
}

func InsertIntoVersion(db *sql.DB, version string) {
	statement, err := db.Prepare("INSERT or REPLACE INTO version (id, Version) VALUES (?, ?)")
	if err != nil {
		log.Fatalln(err)
	}
	defer statement.Close()

	_, err = statement.Exec(1, version)
	if err != nil {
		log.Fatalln(err)
	}
}

func GetVersion(db *sql.DB) string {
	rows, err := db.Query("SELECT Version FROM version")
	if err != nil {
		log.Fatalln(err)
	}
	var version string
	for rows.Next() {
		rows.Scan(&version)
	}
	return version
}

// Fonts table

func CreateFontsTable(db *sql.DB) {
	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS fonts (id INTEGER PRIMARY KEY, Name TEXT, ContentType TEXT, BrowserDownloadUrl TEXT)")
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer statement.Close()
	_, err = statement.Exec()
	if err != nil {
		log.Fatal(err)
	}
}

func DeleteFontsTable(db *sql.DB) {
	statement, err := db.Prepare("DROP TABLE IF EXISTS fonts")
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer statement.Close()
	_, err = statement.Exec()
	if err != nil {
		log.Fatalln(err)
	}
}

func InsertIntoFonts(db *sql.DB, fonts []types.Font) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	statement, err := tx.Prepare("INSERT INTO fonts (Id, Name, ContentType, BrowserDownloadUrl) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer statement.Close()

	for _, font := range fonts {
		_, err = statement.Exec(font.Id, utils.FontNameWithoutExtention(font.Name), font.ContentType, font.BrowserDownloadUrl)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

func GetAllFonts(db *sql.DB) []types.Font {
	var fonts []types.Font
	var font types.Font
	rows, err := db.Query("SELECT Id, Name, ContentType, BrowserDownloadUrl FROM fonts")
	if err != nil {
		log.Fatalln(err)
	}
	for rows.Next() {
		rows.Scan(&font.Id, &font.Name, &font.ContentType, &font.BrowserDownloadUrl)
		fonts = append(fonts, font)
	}
	return fonts
}

func FontExists(db *sql.DB, font string) bool {
	var exists bool
	err := db.QueryRow("SELECT (Name == ?) From fonts WHERE Name = ?", font, font).Scan(&exists)

	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		return false
	}

	return exists
}

// Installed fonts table

func CreateInstalledFontsTable(db *sql.DB) {
	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS installedFonts (Id INTEGER PRIMARY KEY, Name TEXT, Version TEXT)")
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer statement.Close()
	_, err = statement.Exec()
	if err != nil {
		log.Fatal(err)
	}
}

func InsertIntoInstalledFonts(db *sql.DB, font types.Font, version string) {
	statement, err := db.Prepare("INSERT INTO installedFonts(Name, Version) VALUES (?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer statement.Close()

	_, err = statement.Exec(font.Name, version)
	if err != nil {
		log.Fatal(err)
	}
}

func GetInstalledFonts(db *sql.DB) []types.Font {
	var fonts []types.Font
	var font types.Font
	rows, err := db.Query("SELECT Id, Name, Version FROM installedFonts")
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&font.Id, &font.Name, &font.InstalledVersion)
		fonts = append(fonts, font)
	}
	sort.Slice(fonts, func(i, j int) bool { return strings.ToLower(fonts[i].Name) < strings.ToLower(fonts[j].Name) })
	return fonts
}

func IsFontInstalled(db *sql.DB, font string) bool {
	var isInstalled bool
	err := db.QueryRow("SELECT (Name == ?) From installedFonts WHERE Name = ?", font, font).Scan(&isInstalled)

	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		return false
	}

	return isInstalled
}

func GetInstalledFont(db *sql.DB, font types.Font) types.Font {
	var installedFont types.Font
	err := db.QueryRow("SELECT Id, Name, Version FROM installedFonts WHERE Name=?", font.Name).Scan(&installedFont.Id, &installedFont.Name, &installedFont.InstalledVersion)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Fatalln(err)
		}
		log.Fatalln(err)
	}

	return installedFont
}

func UpdateInstalledFont(db *sql.DB, name string, version string) {
	statement, err := db.Prepare("UPDATE installedFonts SET Version=? WHERE Name=?")
	if err != nil {
		log.Fatalln(err)
	}
	statement.Exec(version, name)
}

func DeleteInstalledFont(db *sql.DB, name string) {
	statement, err := db.Prepare("DELETE FROM installedFonts WHERE Name=?")
	if err != nil {
		log.Fatalln(err)
	}
	statement.Exec(name)
}
