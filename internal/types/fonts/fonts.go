package fonts

import (
	"github.com/getnf/getnf/internal/utils"
)

// Fonts

type NerdFonts struct {
	TagName string `json:"tag_name"`
	Assets  []Font `json:"assets"`
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
	return fs.TagName
}

func (fs NerdFonts) GetFonts() []Font {
	isTar := func(f Font) bool { return f.ContentType == "application/x-xz" }
	fonts := utils.Filter(fs.Assets, isTar)
	return fonts
}

func (fs NerdFonts) GetFont(f string) Font {
	isWantedFont := func(x Font) bool { return x.Name == f }
	font := utils.Filter(fs.Assets, isWantedFont)
	return font[0]
}

func (f *Font) AddInstalledVersion(ver string) {
	f.InstalledVersion = ver
}

func (f *Font) AddAvailableVersion(ver string) {
	f.AvailableVersion = ver
}
