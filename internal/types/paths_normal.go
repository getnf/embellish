//go:build normal
// +build normal

package types

import (
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
)

func NewPaths() *Paths {
	paths := &Paths{}

	paths.Download = filepath.Join(xdg.UserDirs.Download, "embellish")
	paths.Install = xdg.FontDirs[0]
	paths.Db = filepath.Join(xdg.DataHome, "embellish")

	os.MkdirAll(paths.Download, 0755)
	os.MkdirAll(paths.Install, 0755)
	os.MkdirAll(paths.Db, 0755)

	return paths
}
