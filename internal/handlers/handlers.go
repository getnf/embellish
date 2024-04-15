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
	"strings"

	"github.com/getnf/getnf/internal/db"
	"github.com/getnf/getnf/internal/types"
	"github.com/getnf/getnf/internal/utils"

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

func FontsWithVersion(database *sql.DB, fonts []types.Font) []types.Font {
	var results []types.Font
	for font := range fonts {
		var f types.Font = fonts[font]
		if db.IsFontInstalled(database, f) {
			installedFont := db.GetInstalledFont(database, f)
			f.AddVersion(installedFont.InstalledVersion)
		} else {
			f.AddVersion("-")
		}
		results = append(results, f)
	}
	return results
}

func ListFonts(fonts []types.Font) {
	for font := range fonts {
		var f types.Font = fonts[font]
		fmt.Printf("Name: %v, Installed version: %v\n", f.Name, f.InstalledVersion)
	}
}

func DownloadTar(fontURL string, path string, name string) (string, error) {
	fullPath := path + "/" + name
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

	fontNameWithExtention := strings.Split(name, ".")[0]

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
		filename := extractPath + "/" + fontNameWithExtention + "/" + header.Name

		// Create directories if they don't exist, if the tar contains directories
		if header.Typeflag == tar.TypeDir {
			err := os.MkdirAll(filename, 0755)
			if err != nil {
				return err
			}
			continue
		}

		if _, err := os.Stat(extractPath + "/" + fontNameWithExtention); errors.Is(err, os.ErrNotExist) {
			err := os.Mkdir(extractPath+"/"+fontNameWithExtention, os.ModePerm)
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

func CleanUpArchive(archivePath string) error {

	if _, err := os.Stat(archivePath); os.IsNotExist(err) {
		return fmt.Errorf("archive file does not exist")
	} else {
		err = os.Remove(archivePath)
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
