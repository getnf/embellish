//go:build flatpak
// +build flatpak

package types

import (
	"log"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
)

func NewPaths() *Paths {
	paths := &Paths{}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln(err)
	}

	paths.Download = filepath.Join(xdg.UserDirs.Download, "embellish")
	paths.Install = filepath.Join(homeDir, ".local", "share", "fonts")
	paths.Db = filepath.Join(xdg.DataHome, "embellish")

	os.MkdirAll(paths.Download, 0755)
	os.MkdirAll(paths.Install, 0755)
	os.MkdirAll(paths.Db, 0755)

	return paths
}
