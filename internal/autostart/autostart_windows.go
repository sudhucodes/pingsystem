//go:build windows

package autostart

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sys/windows/registry"
)

const (
	runKeyPath = `Software\Microsoft\Windows\CurrentVersion\Run`
	appName    = "PingSystem"
)

// Enable adds PingSystem to HKCU\Software\Microsoft\Windows\CurrentVersion\Run.
func Enable() error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	absPath, err := filepath.Abs(exePath)
	if err != nil {
		absPath = exePath
	}

	key, err := registry.OpenKey(registry.CURRENT_USER, runKeyPath, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("failed to open registry key %s: %w", runKeyPath, err)
	}
	defer key.Close()

	// Quote the path in case it contains spaces
	formattedValue := fmt.Sprintf(`"%s"`, absPath)

	if err := key.SetStringValue(appName, formattedValue); err != nil {
		return fmt.Errorf("failed to set registry value: %w", err)
	}

	return nil
}

// Disable removes PingSystem from HKCU\Software\Microsoft\Windows\CurrentVersion\Run.
func Disable() error {
	key, err := registry.OpenKey(registry.CURRENT_USER, runKeyPath, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("failed to open registry key %s: %w", runKeyPath, err)
	}
	defer key.Close()

	if err := key.DeleteValue(appName); err != nil && err != registry.ErrNotExist {
		return fmt.Errorf("failed to delete registry value: %w", err)
	}

	return nil
}

// IsEnabled checks if PingSystem is set to autostart in the registry.
func IsEnabled() (bool, string, error) {
	key, err := registry.OpenKey(registry.CURRENT_USER, runKeyPath, registry.QUERY_VALUE)
	if err != nil {
		if err == registry.ErrNotExist {
			return false, "", nil
		}
		return false, "", fmt.Errorf("failed to open registry key %s: %w", runKeyPath, err)
	}
	defer key.Close()

	val, _, err := key.GetStringValue(appName)
	if err != nil {
		if err == registry.ErrNotExist {
			return false, "", nil
		}
		return false, "", fmt.Errorf("failed to query registry value: %w", err)
	}

	exePath, _ := os.Executable()
	absPath, _ := filepath.Abs(exePath)

	cleanVal := strings.Trim(val, `"`)
	cleanPath := strings.Trim(absPath, `"`)

	enabled := strings.EqualFold(cleanVal, cleanPath)
	return enabled, val, nil
}
