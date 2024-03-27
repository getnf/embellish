package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/getnf/getnf/internal/handlers"
	"github.com/getnf/getnf/internal/types"
)

func getData() (types.NerdFonts, error) {
	url := "https://api.github.com/repos/ryanoasis/nerd-fonts/releases/latest"
	resp, err := http.Get(url)
	if err != nil {
		return types.NerdFonts{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return types.NerdFonts{}, err
	}

	var data types.NerdFonts
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatalln(err)
	}
	return data, nil
}

func main() {

	keepArchives := true

	data, err := getData()
	if err != nil {
		log.Fatalln(err)
	}

	homePath, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln(err)
	}

	downloadPath := homePath + "/.local/share/getnf"
	extractPath := homePath + "/.local/share/fonts"

	font := data.GetFont("MPlus.tar.xz")

	downloadedArchivePath, err := handlers.DownloadTar(font.BrowserDownloadUrl, downloadPath, font.Name)
	if err != nil {
		log.Fatalln(err)
	}

	err = handlers.ExtractTar(downloadedArchivePath, extractPath, font.Name)
	if err != nil {
		log.Fatalln(err)
	}

	if !keepArchives {
		err = handlers.CleanUpArchives(downloadedArchivePath)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
