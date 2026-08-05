//go:build !windows

package sysinfo

import "runtime"

// GetWindowsVersion returns OS name when running on non-Windows platforms.
func GetWindowsVersion() string {
	return runtime.GOOS + " (" + runtime.GOARCH + ")"
}
