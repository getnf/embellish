package types

import (
	"database/sql"
	"sort"
	"strings"

	"github.com/getnf/embellish/internal/utils"
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
}

func (Args) Version() string {
	return "getnf v0.3.0"
}

// params

type GuiParams struct {
	Data         NerdFonts
	Database     *sql.DB
	DownloadPath string
	ExtractPath  string
}
