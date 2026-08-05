//go:build windows

package sysinfo

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
)

// GetWindowsVersion reads the Windows OS version from the Windows Registry.
func GetWindowsVersion() string {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion`, registry.QUERY_VALUE)
	if err != nil {
		return "Windows (Unknown Version)"
	}
	defer k.Close()

	productName, _, err := k.GetStringValue("ProductName")
	if err != nil || productName == "" {
		productName = "Windows"
	}

	displayVersion, _, err := k.GetStringValue("DisplayVersion")
	if err != nil || displayVersion == "" {
		displayVersion, _, _ = k.GetStringValue("ReleaseId")
	}

	buildNumber, _, err := k.GetStringValue("CurrentBuildNumber")
	if err != nil || buildNumber == "" {
		buildNumber, _, _ = k.GetStringValue("CurrentBuild")
	}

	verStr := productName
	if displayVersion != "" {
		verStr += fmt.Sprintf(" %s", displayVersion)
	}
	if buildNumber != "" {
		verStr += fmt.Sprintf(" (Build %s)", buildNumber)
	}

	return verStr
}
