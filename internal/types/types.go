package types

import (
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/getnf/getnf/internal/utils"
)

// Fonts

type NerdFonts struct {
	Version string `json:"tag_name"`
	Fonts   []Font `json:"assets"`
}

type Font struct {
	Id                 int    `json:"id"`
	Name               string `json:"name"`
	ContentType        string `json:"content_type"`
	BrowserDownloadUrl string `json:"browser_download_url"`
	AvailableVersion   string
	InstalledVersion   string
}

func (fs NerdFonts) GetVersion() string {
	return fs.Version
}

func (fs NerdFonts) GetFonts() []Font {
	isTar := func(f Font) bool { return f.ContentType == "application/x-xz" }
	fonts := utils.Filter(fs.Fonts, isTar)
	sort.Slice(fonts, func(i, j int) bool { return strings.ToLower(fonts[i].Name) < strings.ToLower(fonts[j].Name) })
	return fonts
}

func (fs NerdFonts) GetFont(f string) Font {
	isWantedFont := func(x Font) bool { return x.Name == f }
	font := utils.Filter(fs.Fonts, isWantedFont)
	return font[0]
}

func (fs NerdFonts) GetFontsNames() []string {
	fontNames := utils.Fold(fs.Fonts, func(f Font) string {
		return f.Name
	})
	sort.Slice(fontNames, func(i, j int) bool { return strings.ToLower(fontNames[i]) < strings.ToLower(fontNames[j]) })
	return fontNames
}

func (f *Font) AddInstalledVersion(ver string) {
	f.InstalledVersion = ver
}

func (f *Font) AddAvailableVersion(ver string) {
	f.AvailableVersion = ver
}

// Command line argumetns

type InstallCmd struct {
	Fonts []string `arg:"positional" help:"list of space separated fonts to install"`
}

type UninstallCmd struct {
	Fonts []string `arg:"positional" help:"list of space separated fonts to uninstall"`
}

type ListCmd struct {
	Installed bool `arg:"-i" help:"list only installed fonts"`
}

type UpdateCmd struct {
	Update bool `default:"true"`
}

type Args struct {
	Install    *InstallCmd   `arg:"subcommand:install" help:"install fonts"`
	Uninstall  *UninstallCmd `arg:"subcommand:uninstall" help:"uninstall fonts"`
	List       *ListCmd      `arg:"subcommand:list" help:"list fonts"`
	Update     *UpdateCmd    `arg:"subcommand:update" help:"update installed fonts"`
	KeepTars   bool          `arg:"-k" help:"Keep archives in the download location"`
	ForceCheck bool          `arg:"-f" help:"Force checking for updates"`
	Global     bool          `arg:"-g" help:"Do the operation globally, for all users"`
}

func (Args) Version() string {
	return "getnf v1.0.0"
}

// paths

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

	switch os := utils.OsType(); os {
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
		paths.Download.User = ""
		paths.Download.Root = filepath.Join(homeDir, "Downloads", "getnf")
		paths.Install.User = ""
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
