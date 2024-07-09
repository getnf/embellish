//go:build terminal
// +build terminal

package handlers

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/briandowns/spinner"
	"github.com/getnf/embellish/internal/db"
	"github.com/getnf/embellish/internal/types"
	"github.com/getnf/embellish/internal/utils"
)

func FontsWithVersion(database *sql.DB, fonts []types.Font, version string) []types.Font {
	var results []types.Font
	for _, font := range fonts {
		if db.IsFontInstalled(database, font.Name) {
			installedFont := db.GetInstalledFont(database, font)
			font.AddInstalledVersion(installedFont.InstalledVersion)
		} else {
			font.AddInstalledVersion("-")
		}
		font.AddAvailableVersion(version)
		results = append(results, font)
	}
	return results
}

func ListFonts(fonts []types.Font, onlyInstalled bool) {
	isInstalledFont := func(x types.Font) bool { return x.InstalledVersion != "-" }
	if onlyInstalled {
		fonts = utils.Filter(fonts, isInstalledFont)
	}

	writer := tabwriter.NewWriter(os.Stdout, 0, 8, 4, '\t', tabwriter.AlignRight)

	fmt.Fprintln(writer, "Name:\tAvailable Version:\tInstalled Version:")

	if len(fonts) == 0 && onlyInstalled {
		fmt.Println("No fonts have been installed yet")
		return
	}
	for _, font := range fonts {
		fmt.Fprintln(writer, font.Name, "\t", font.AvailableVersion, "\t", font.InstalledVersion)
	}
	writer.Flush()
}

func HandleInstall(args types.Args, database *sql.DB, data types.NerdFonts, downloadPath string, extractPath string) error {
	var installedFonts []string
	var fontsToInstall []string
	for _, font := range args.Install.Fonts {
		if db.FontExists(database, font) {
			fontsToInstall = append(fontsToInstall, font)
		} else {
			fmt.Printf("%v is not a nerd font\n", font)
			fuzzySearchedFont, err := FuzzySearchFonts(font, data.GetFontsNames())
			if err != nil {
				return fmt.Errorf("did you mean: %v: ", err)
			}
			// fmt.Printf("did you mean: %v\n", fuzzySearchedFont)
			return fmt.Errorf("did you mean: %v: ", fuzzySearchedFont)
		}
	}
	if len(fontsToInstall) > 0 {
		for _, font := range fontsToInstall {
			f := data.GetFont(font)
			err := PlatformInstallFont(f, downloadPath, extractPath, args.KeepTars)
			if err != nil {
				return err
			}
			db.InsertIntoInstalledFonts(database, f, data.GetVersion())
			installedFonts = append(installedFonts, font)
		}
	}
	if len(installedFonts) > 0 {
		fmt.Printf("Installed font(s): %v\n", strings.Join(installedFonts, ", "))
	}

	return nil
}

func HandleUninstall(args types.Args, database *sql.DB, data types.NerdFonts, extractPath string) error {
	var fontsToUninstall []string
	for _, font := range args.Uninstall.Fonts {
		if db.IsFontInstalled(database, font) {
			fontsToUninstall = append(fontsToUninstall, font)
		} else {
			fmt.Printf("%v is either not installed or is not a nerd font\n", font)
			fuzzySearchedFont, err := FuzzySearchFonts(font, data.GetFontsNames())
			if err != nil {
				return fmt.Errorf("did you mean: %v: ", err)
			}
			return fmt.Errorf("did you mean: %v: ", fuzzySearchedFont)
		}
	}
	if len(fontsToUninstall) > 0 {
		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		s.Suffix = " Uninstalling fonts"
		s.Color("red")
		s.Start()
		for _, font := range fontsToUninstall {
			err := PlatformUninstallFont(extractPath, font)
			if err != nil {
				s.Stop()
				return err
			}
			db.DeleteInstalledFont(database, font)
		}
		s.FinalMSG = "uninstalled font(s): " + strings.Join(fontsToUninstall, ", ") + "\n"
		s.Stop()
	}

	return nil
}

func IsAdmin() bool {
	isAdmine := PlatformIsAdmin()

	return isAdmine
}
