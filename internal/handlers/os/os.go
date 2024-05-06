package os

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/getnf/getnf/internal/types/paths"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

func IsAdmin() (bool, error) {
	if os.Geteuid() == 0 {
		return true, nil
	}

	if windows.GetCurrentProcessToken().IsElevated() {
		return true, nil
	}

	var message string
	if paths.OsType() == "linux" || paths.OsType() == "darwin" {
		message = "getnf has to be run with superuser privileges when using the -g flag"
	} else {
		message = "getnf has to be run as administrator when using the -g flag"
	}

	return false, fmt.Errorf(message)
}

func WriteToRegistry(path string, fontName string, fileName string) error {
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

func RemoveFromRegistry(name string) error {
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
