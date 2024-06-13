//go:build linux || solaris || openbsd || freebsd || netbsd
// +build linux solaris openbsd freebsd netbsd

package handlers

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/getnf/getnf/internal/types"

	"github.com/briandowns/spinner"
)

func PlatformInstallFont(font types.Font, downloadPath string, extractPath string, keepTar bool) error {
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	s.Color("yellow")
	s.Suffix = " Downloading font " + font.Name
	s.Start()
	downloadedTar, err := downloadFont(font.BrowserDownloadUrl, downloadPath, font.Name)
	if err != nil {
		return fmt.Errorf("error downloading the tar file: %v", err)
	}
	s.Suffix = " Installing font " + font.Name
	s.Color("green")
	s.Restart()
	_, err = extractFont(downloadedTar, extractPath, font.Name)
	if err != nil {
		return fmt.Errorf("error extracting the tar file: %v", err)
	}
	if !keepTar {
		deleteTar(downloadedTar)
	}
	s.Stop()

	return nil
}

func PlatformUninstallFont(path string, name string) error {
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

func PlatformIsAdmin() bool {
	if os.Geteuid() == 0 {
		return true
	}

	return false
}
