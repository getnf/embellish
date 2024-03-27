package handlers

import (
	"archive/tar"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/ulikunitz/xz"
)

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

func CleanUpArchives(archivePath string) error {

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
