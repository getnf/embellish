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
	"text/tabwriter"
	"time"

	"github.com/getnf/getnf/internal/db"
	fontsTypes "github.com/getnf/getnf/internal/types/fonts"
	"github.com/getnf/getnf/internal/utils"

	"github.com/briandowns/spinner"
	"github.com/ulikunitz/xz"
)

func GetData() (fontsTypes.NerdFonts, error) {
	url := "https://api.github.com/repos/ryanoasis/nerd-fonts/releases/latest"
	resp, err := http.Get(url)
	if err != nil {
		return fontsTypes.NerdFonts{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fontsTypes.NerdFonts{}, err
	}

	var data fontsTypes.NerdFonts
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatalln(err)
	}
	return data, nil
}

func FontsWithVersion(database *sql.DB, fonts []fontsTypes.Font, version string) []fontsTypes.Font {
	var results []fontsTypes.Font
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

func ListFonts(fonts []fontsTypes.Font, onlyInstalled bool) {
	isInstalledFont := func(x fontsTypes.Font) bool { return x.InstalledVersion != "-" }
	if onlyInstalled {
		fonts = utils.Filter(fonts, isInstalledFont)
	}

	writer := tabwriter.NewWriter(os.Stdout, 0, 8, 4, '\t', tabwriter.AlignRight)

	fmt.Fprintln(writer, "Name:\tAvailable version:\tInstalledVersion:")

	if len(fonts) == 0 && onlyInstalled {
		fmt.Println("No fonts have been installed yet")
		return
	}
	for _, font := range fonts {
		fmt.Fprintln(writer, font.Name, "\t", font.AvailableVersion, "\t", font.InstalledVersion)
	}
	writer.Flush()
}

func DownloadTar(fontURL string, path string, name string) (string, error) {
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
func ExtractTar(archivePath string, extractPath string, name string) error {

	// Decompress the xz stream
	fontArchive, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	xzReader, err := xz.NewReader(fontArchive)
	if err != nil {
		return err
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
			return err
		}

		// Extract the file name from the header
		filename := filepath.Join(extractPath, name, header.Name)
		extractDir := filepath.Join(extractPath, name)

		// Create directories if they don't exist, if the tar contains directories
		if header.Typeflag == tar.TypeDir {
			err := os.MkdirAll(filename, 0755)
			if err != nil {
				return err
			}
			continue
		}

		if _, err := os.Stat(extractDir); errors.Is(err, os.ErrNotExist) {
			err := os.Mkdir(extractDir, os.ModePerm)
			if err != nil {
				return err
			}
		}

		// Create file with same permissions as in the tar file
		file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
		if err != nil {
			return err
		}
		defer file.Close()

		// Write file content to disk
		_, err = io.Copy(file, tarReader)
		if err != nil {
			return err
		}
	}

	return nil
}

func DeleteTar(tarPath string) error {
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

func InstallFont(font fontsTypes.Font, downloadPath string, extractPath string, keepTar bool) {
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	s.Color("yellow")
	s.Suffix = " Downloading font " + font.Name
	s.Start()
	downloadedTar, err := DownloadTar(font.BrowserDownloadUrl, downloadPath, font.Name)
	if err != nil {
		log.Fatalln(err)
	}
	s.Suffix = " Installing font " + font.Name
	s.Color("green")
	s.Restart()
	ExtractTar(downloadedTar, extractPath, font.Name)
	if !keepTar {
		DeleteTar(downloadedTar)
	}
	s.Stop()
}

func UninstallFont(path string, name string) error {
	fontPath := filepath.Join(path, name)
	if _, err := os.Stat(fontPath); os.IsNotExist(err) {
		return fmt.Errorf("font %v is not installed", name)
	} else {
		err = os.RemoveAll(fontPath)
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
