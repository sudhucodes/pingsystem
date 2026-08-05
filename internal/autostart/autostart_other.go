//go:build !windows

package autostart

import "errors"

func Enable() error {
	return errors.New("autostart registry configuration is only supported on Windows")
}

func Disable() error {
	return errors.New("autostart registry configuration is only supported on Windows")
}

func IsEnabled() (bool, string, error) {
	return false, "", nil
}
