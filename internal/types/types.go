package types

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/getnf/getnf/internal/utils"
)

type NerdFonts struct {
	TagName string `json:"tag_name"`
	Assets  []Font `json:"assets"`
}

type Font struct {
	Name               string `json:"name"`
	ContentType        string `json:"content_type"`
	BrowserDownloadUrl string `json:"browser_download_url"`
}

func (fs NerdFonts) GetVersion() (int, error) {
	re := regexp.MustCompile("[0-9]+")
	versionCleaned := re.FindAllString(fs.TagName, -1)
	version, err := strconv.Atoi(strings.Join(versionCleaned[:], ""))
	if err != nil {
		return 0, err
	}
	return version, nil
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
