//go:build windows
// +build windows

package handlers

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/getnf/getnf/internal/types"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

// got inspiration for this from https://github.com/Crosse/font-install

func PlatformInstallFont(font types.Font, downloadPath string, extractPath string, keepTar bool) error {
	downloadedTar, err := downloadFont(font.BrowserDownloadUrl, downloadPath, font.Name)
	if err != nil {
		return fmt.Errorf("error downloading the tar file: %v", err)
	}
	extractedTar, err := extractFont(downloadedTar, extractPath, font.Name)
	if err != nil {
		return fmt.Errorf("error extracting the tar file: %v", err)
	}
	for _, fileName := range extractedTar {
		err = removeFromRegistry(fileName)
		if err != nil {
			log.Fatalln(err)
		}
		err = writeToRegistry(extractPath, font.Name, fileName)
		if err != nil {
			log.Fatalln(err)
		}
	}
	if !keepTar {
		deleteTar(downloadedTar)
	}

	return nil
}

func PlatformUninstallFont(path string, name string) error {
	fontPath := filepath.Join(path, name)
	fontFiles, err := os.ReadDir(fontPath)
	if err != nil {
		log.Fatalln(err)
	}

	var fileNames []string

	if _, err := os.Stat(fontPath); os.IsNotExist(err) {
		return fmt.Errorf("font %v is not installed", name)
	} else {
		for _, file := range fontFiles {
			fileNames = append(fileNames, file.Name())
		}

		err = os.RemoveAll(fontPath)
		if err != nil {
			return err
		}
		for _, file := range fileNames {
			removeFromRegistry(file)
		}
	}
	return nil
}

func writeToRegistry(path string, fontName string, fileName string) error {
	fullPath := filepath.Join(path, fontName, fileName)
	k, err := registry.OpenKey(
		registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows NT\CurrentVersion\Fonts`,
		registry.WRITE)
	if err != nil {
		os.Remove(fullPath)
		return fmt.Errorf("error opening registry key: %w", err)
	}
	defer k.Close()

	valueName := fmt.Sprintf("%s (TrueType)", fileName)
	err = k.SetStringValue(valueName, fullPath)
	if err != nil {
		os.Remove(fullPath)
		return fmt.Errorf("error writing to registry: %w", err)
	}

	return nil
}

func removeFromRegistry(name string) error {
	k, err := registry.OpenKey(
		registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows NT\CurrentVersion\Fonts`,
		registry.WRITE)
	if err != nil {
		return fmt.Errorf("error opening registry key: %w", err)
	}
	defer k.Close()

	valueName := fmt.Sprintf("%s (TrueType)", name)

	// Check if the value exists before attempting to remove it
	exists, err := valueExistsInRegistry(k, valueName)
	if err != nil {
		return fmt.Errorf("error checking if value exists: %w", err)
	}
	if !exists {
		return nil
	}

	err = k.DeleteValue(valueName)
	if err != nil {
		return fmt.Errorf("error deleting registry value: %w", err)
	}

	return nil
}

func valueExistsInRegistry(key registry.Key, name string) (bool, error) {
	k, err := registry.OpenKey(key, "", registry.QUERY_VALUE)
	if err != nil {
		return false, fmt.Errorf("error opening registry key: %w", err)
	}
	defer k.Close()
	_, _, err = k.GetStringValue(name)
	if err != nil {
		if err == registry.ErrNotExist {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func PlatformIsAdmin() bool {
	if windows.GetCurrentProcessToken().IsElevated() {
		return true
	}

	return false
}
