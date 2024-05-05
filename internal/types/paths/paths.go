package paths

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
)

type Paths struct {
	Download string
	Install  InstallPaths
}

type InstallPaths struct {
	User string
	Root string
}

func NewPaths() *Paths {
	paths := &Paths{}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln(err)
	}

	switch os := osType(); os {
	case "linux":
		paths.Download = filepath.Join(homeDir, "Downloads", "getnf")
		paths.Install.User = filepath.Join(homeDir, ".local", "share", "fonts")
		paths.Install.Root = "/usr/share/fonts"
	case "darwin":
		paths.Download = filepath.Join(homeDir, "Downloads", "getnf")
		paths.Install.User = filepath.Join(homeDir, "Library", "Fonts")
		paths.Install.Root = "/Library/Fonts"
	case "windows":
		paths.Download = filepath.Join(homeDir, "Downloads")
		// User font directory for Windows
		paths.Install.User = filepath.Join(homeDir, "AppData", "Local", "Microsoft", "Windows", "Fonts")
		// System-wide font directory for Windows
		paths.Install.Root = filepath.Join("C:\\Windows", "Fonts")
	default:
		log.Fatalln("unsupported operating system")
	}

	return paths
}

func (p *Paths) GetDownloadPath() string {
	return p.Download
}

func (p *Paths) GetUserInstallPath() string {
	return p.Install.User
}

func (p *Paths) GetRootInstallPath() string {
	return p.Install.Root
}

func osType() string {
	switch os := runtime.GOOS; os {
	case "darwin":
		return "darwin"
	case "linux":
		return "linux"
	case "windows":
		return "windows"
	default:
		return "unsupported"
	}
}
