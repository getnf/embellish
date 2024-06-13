package handlers

import (
	"archive/tar"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/briandowns/spinner"
	"github.com/getnf/getnf/internal/db"
	"github.com/getnf/getnf/internal/types"
	"github.com/getnf/getnf/internal/utils"
	"github.com/lithammer/fuzzysearch/fuzzy"

	"github.com/ulikunitz/xz"
)

func GetData() (types.NerdFonts, error) {
	url := "https://api.github.com/repos/ryanoasis/nerd-fonts/releases/latest"
	resp, err := http.Get(url)
	if err != nil {
		return types.NerdFonts{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return types.NerdFonts{}, err
	}

	var data types.NerdFonts
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatalln(err)
	}
	return data, nil
}

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

func downloadFont(fontURL string, path string, name string) (string, error) {
	fullPath := path + "/" + name + ".tar.xz"
	resp, err := http.Get(fontURL)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", err
	}

	// Make sure the path exists
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	// Create the file
	out, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}

	defer out.Close()
	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}

	return fullPath, nil
}

// extractTar extracts files from a tar archive provided in the reader
func extractFont(archivePath string, extractPath string, name string) ([]string, error) {
	var listOfInstalledFonts []string

	// Decompress the xz stream
	fontArchive, err := os.Open(archivePath)
	if err != nil {
		return []string{""}, err
	}
	xzReader, err := xz.NewReader(fontArchive)
	if err != nil {
		return []string{""}, err
	}

	defer fontArchive.Close()

	// Create a tar reader from the decompressed stream
	tarReader := tar.NewReader(xzReader)

	// Iterate over each file in the tar archive
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			// End of tar archive
			break
		}
		if err != nil {
			return []string{""}, err
		}

		// Extract the file name from the header
		fullPath := filepath.Join(extractPath, name, header.Name)
		extractPath := filepath.Join(extractPath, name)

		// Create directories if they don't exist, if the tar contains directories
		if header.Typeflag == tar.TypeDir {
			err := os.MkdirAll(fullPath, 0755)
			if err != nil {
				return []string{""}, err
			}
			continue
		}

		if _, err := os.Stat(extractPath); errors.Is(err, os.ErrNotExist) {
			err := os.Mkdir(extractPath, os.ModePerm)
			if err != nil {
				return []string{""}, err
			}
		}

		// Create file with same permissions as in the tar file
		file, err := os.OpenFile(fullPath, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
		if err != nil {
			return []string{""}, err
		}
		defer file.Close()

		// Write file content to disk
		_, err = io.Copy(file, tarReader)
		if err != nil {
			return []string{""}, err
		}

		listOfInstalledFonts = append(listOfInstalledFonts, header.Name)
	}

	return listOfInstalledFonts, nil
}

func deleteTar(tarPath string) error {
	if _, err := os.Stat(tarPath); os.IsNotExist(err) {
		return fmt.Errorf("tar file does not exist")
	} else {
		err = os.Remove(tarPath)
		if err != nil {
			return err
		}
	}
	return nil
}

func IsUpdateAvilable(remote string, local string) bool {
	remoteVersion, err := utils.StringToInt(remote)
	if err != nil {
		log.Fatalln(err)
	}

	localVersion, err := utils.StringToInt(local)
	if err != nil {
		log.Fatalln(err)
	}
	if remoteVersion > localVersion {
		return true
	} else {
		return false
	}
}

func InstallFont(font types.Font, downloadPath string, extractPath string, keepTar bool) {
	PlatformInstallFont(font, downloadPath, extractPath, keepTar)

}

func UninstallFont(path string, name string) error {
	err := PlatformUninstallFont(path, name)
	if err != nil {
		return err
	}
	return nil
}

func HandleInstall(args types.Args, database *sql.DB, data types.NerdFonts, downloadPath string, extractPath string) {
	var installedFonts []string
	var fontsToInstall []string
	for _, font := range args.Install.Fonts {
		if db.FontExists(database, font) {
			fontsToInstall = append(fontsToInstall, font)
		} else {
			fmt.Printf("%v is not a nerd font\n", font)
			fuzzySearchedFont, err := FuzzySearchFonts(font, data.GetFontsNames())
			if err != nil {
				fmt.Printf("did you mean: %v\n", err)
				os.Exit(0)
			}
			fmt.Printf("did you mean: %v\n", fuzzySearchedFont)
			return
		}
	}
	if len(fontsToInstall) > 0 {
		for _, font := range fontsToInstall {
			f := data.GetFont(font)
			InstallFont(f, downloadPath, extractPath, args.KeepTars)
			db.InsertIntoInstalledFonts(database, f, data.GetVersion())
			installedFonts = append(installedFonts, font)
		}
	}
	if len(installedFonts) > 0 {
		fmt.Printf("Installed font(s): %v\n", strings.Join(installedFonts, ", "))
	}
}

func HandleUninstall(args types.Args, database *sql.DB, data types.NerdFonts, extractPath string) {
	var fontsToUninstall []string
	for _, font := range args.Uninstall.Fonts {
		if db.IsFontInstalled(database, font) {
			fontsToUninstall = append(fontsToUninstall, font)
		} else {
			fmt.Printf("%v is either not installed or is not a nerd font\n", font)
			fuzzySearchedFont, err := FuzzySearchFonts(font, data.GetFontsNames())
			if err != nil {
				fmt.Printf("did you mean: %v\n", err)
				os.Exit(0)
			}
			fmt.Printf("did you mean: %v\n", fuzzySearchedFont)
			return
		}
	}
	if len(fontsToUninstall) > 0 {
		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		s.Suffix = " Uninstalling fonts"
		s.Color("red")
		s.Start()
		for _, font := range fontsToUninstall {
			UninstallFont(extractPath, font)
			db.DeleteInstalledFont(database, font)
		}
		s.FinalMSG = "uninstalled font(s): " + strings.Join(fontsToUninstall, ", ") + "\n"
		s.Stop()
	}
}

func IsFontUpdatAvilable(database *sql.DB, data types.NerdFonts) bool {
	updateCount := 0
	installedFonts := db.GetInstalledFonts(database)
	for _, font := range installedFonts {
		if IsUpdateAvilable(data.GetVersion(), font.InstalledVersion) {
			updateCount++
		}
	}

	return updateCount > 0
}

func HandleUpdate(args types.Args, database *sql.DB, data types.NerdFonts, downloadPath string, extractPath string) {
	if IsFontUpdatAvilable(database, data) {
		installedFonts := db.GetInstalledFonts(database)
		for _, font := range installedFonts {
			f := data.GetFont(font.Name)
			InstallFont(f, downloadPath, extractPath, args.KeepTars)
			db.UpdateInstalledFont(database, font.Name, data.GetVersion())
		}
	} else {
		fmt.Println("No updates are available")
	}
}

func IsAdmin() bool {
	isAdmine := PlatformIsAdmin()

	return isAdmine
}

func FuzzySearchFonts(font string, fonts []string) ([]string, error) {
	matches := fuzzy.RankFindFold(font, fonts)
	var match []string
	sort.Sort(matches)

	if len(matches) > 0 {
		var topMatches fuzzy.Ranks
		if len(matches) > 3 {
			topMatches = matches[0:3]
		} else {
			size := len(matches)
			topMatches = matches[0:size]
		}
		for _, font := range topMatches {
			match = append(match, font.Target)
		}
	} else {
		return []string{""}, fmt.Errorf("no match found")
	}
	return match, nil
}
