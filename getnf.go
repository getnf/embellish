package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/getnf/getnf/internal/db"
	"github.com/getnf/getnf/internal/handlers"
	"github.com/getnf/getnf/internal/types/args"
	fontsTypes "github.com/getnf/getnf/internal/types/fonts"
	"github.com/getnf/getnf/internal/types/paths"

	"github.com/alexflint/go-arg"
	"github.com/briandowns/spinner"
)

func setupDB(database *sql.DB, remoteData fontsTypes.NerdFonts) {
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
	var args args.Args
	arg.MustParse(&args)

	isGlobal := args.Global
	isAdmine, _ := handlers.IsAdmin()

	var database *sql.DB

	if paths.OsType() == "windows" && !isAdmine {
		log.Fatalln("getnf can't install locally on windows, please run getnf as administrator")
	}

	if isGlobal && isAdmine {
		database = db.OpenGlobalDB()
	} else if isAdmine {
		log.Fatalln("only run getnf with elevated privileges if using the -g flag")
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

	var data fontsTypes.NerdFonts

	data.TagName = db.GetVersion(database)
	data.Assets = db.GetAllFonts(database)

	paths := paths.NewPaths()
	var extractPath string
	var downloadPath string
	if isGlobal {
		_, err := handlers.IsAdmin()
		if err != nil {
			log.Fatalln(err)
		}
		downloadPath = paths.GetRootDownloadPath()
		extractPath = paths.GetRootInstallPath()
	} else {
		downloadPath = paths.GetUserDownloadPath()
		extractPath = paths.GetUserInstallPath()
	}

	switch {
	case args.List != nil:
		if !args.List.Installed {
			handlers.ListFonts(handlers.FontsWithVersion(database, data.GetFonts(), data.GetVersion()), false)
		} else {
			handlers.ListFonts(handlers.FontsWithVersion(database, data.GetFonts(), data.GetVersion()), true)
		}
	case args.Install != nil:
		var installedFonts []string
		for _, font := range args.Install.Fonts {
			if db.FontExists(database, font) {
				f := data.GetFont(font)
				handlers.InstallFont(f, downloadPath, extractPath, args.Install.KeepTars)
				db.InsertIntoInstalledFonts(database, f, data.GetVersion())
				installedFonts = append(installedFonts, font)
			} else {
				fmt.Printf("%v is not a nerd font\n", font)
			}
		}
		if len(installedFonts) > 0 {
			fmt.Printf("Installed font(s): %v\n", strings.Join(installedFonts, ", "))
		}
	case args.Uninstall != nil:
		var fontsToUninstall []string
		for _, font := range args.Uninstall.Fonts {
			if db.IsFontInstalled(database, font) {
				fontsToUninstall = append(fontsToUninstall, font)
			} else {
				fmt.Printf("%v is either not installed or is not a nerd font\n", font)
			}
		}
		if len(fontsToUninstall) > 0 {
			s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
			s.Suffix = " Uninstalling fonts"
			s.Color("red")
			s.Start()
			for _, font := range fontsToUninstall {
				handlers.UninstallFont(extractPath, font)
				db.DeleteInstalledFont(database, font)
			}
			s.FinalMSG = "uninstalled font(s): " + strings.Join(fontsToUninstall, ", ") + "\n"
			s.Stop()
		}
	case args.Update != nil:
		updateCount := 0
		var installedFonts []fontsTypes.Font
		installedFonts = db.GetInstalledFonts(database)
		for _, font := range installedFonts {
			if handlers.IsUpdateAvilable(data.GetVersion(), font.InstalledVersion) {
				f := data.GetFont(font.Name)
				handlers.InstallFont(f, downloadPath, extractPath, args.Update.KeepTars)
				db.UpdateInstalledFont(database, font.Name, data.GetVersion())
				updateCount++
			}
		}
		if updateCount == 0 {
			fmt.Println("No updates are available")
		}
	default:
		fmt.Println(args.Version())
	}
}
