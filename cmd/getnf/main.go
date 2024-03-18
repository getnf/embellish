package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

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

	font := data.GetFont("Hack.tar.xz")
	fontTar, err := font.GetFontTar()
	if err != nil {
		log.Fatalln(err)
	}
	err = handlers.ExtractTar(extractPath, strings.Split(font.Name, ".")[0], fontTar)
	if err != nil {
		log.Fatalln(err)
	}
	err = handlers.SaveTar(downloadPath, font.Name, fontTar)
	if err != nil {
		log.Fatalln(err)
	}
}
