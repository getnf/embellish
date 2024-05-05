package paths

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
)

type Paths struct {
	Download downloadPaths
	Install  installPaths
}

type downloadPaths struct {
	User string
	Root string
}

type installPaths struct {
	User string
	Root string
}

func NewPaths() *Paths {
	paths := &Paths{}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln(err)
	}

	tempDir := os.TempDir()

	switch os := OsType(); os {
	case "linux":
		paths.Download.User = filepath.Join(homeDir, "Downloads", "getnf")
		paths.Download.Root = filepath.Join(tempDir, "getnf")
		paths.Install.User = filepath.Join(homeDir, ".local", "share", "fonts")
		paths.Install.Root = "/usr/share/fonts"
	case "darwin":
		paths.Download.User = filepath.Join(homeDir, "Downloads", "getnf")
		paths.Download.Root = filepath.Join(tempDir, "getnf")
		paths.Install.User = filepath.Join(homeDir, "Library", "Fonts")
		paths.Install.Root = "/Library/Fonts"
	case "windows":
		paths.Download.User = filepath.Join(homeDir, "Downloads", "getnf")
		paths.Download.Root = filepath.Join(tempDir, "getnf")
		paths.Install.User = filepath.Join(homeDir, "AppData", "Local", "Microsoft", "Windows", "Fonts")
		paths.Install.Root = filepath.Join("C:\\Windows", "Fonts")
	default:
		log.Fatalln("unsupported operating system")
	}

	return paths
}

func (p *Paths) GetUserDownloadPath() string {
	return p.Download.User
}

func (p *Paths) GetRootDownloadPath() string {
	return p.Download.Root
}

func (p *Paths) GetUserInstallPath() string {
	return p.Install.User
}

func (p *Paths) GetRootInstallPath() string {
	return p.Install.Root
}

func OsType() string {
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
