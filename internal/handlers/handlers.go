package handlers

import (
	"archive/tar"
	"errors"
	"io"
	"os"

	"github.com/ulikunitz/xz"
)

func SaveTar(path string, name string, reader io.Reader) error {

	// create the directory
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			return err
		}
	}

	// Create the file
	out, err := os.Create(path + "/" + name)
	if err != nil {
		return err
	}
	defer out.Close()
	// Write the body to file
	_, err = io.Copy(out, reader)
	if err != nil {
		return err
	}
	return nil
}

// extractTar extracts files from a tar archive provided in the reader
func ExtractTar(path string, name string, reader io.Reader) error {
	// Decompress the xz stream
	xzReader, err := xz.NewReader(reader)
	if err != nil {
		return err
	}

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
		filename := path + "/" + name + "/" + header.Name

		// Create directories if they don't exist, if the tar contains directories
		if header.Typeflag == tar.TypeDir {
			err := os.MkdirAll(filename, 0755)
			if err != nil {
				return err
			}
			continue
		}

		if _, err := os.Stat(path + "/" + name); errors.Is(err, os.ErrNotExist) {
			err := os.Mkdir(path+"/"+name, os.ModePerm)
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
