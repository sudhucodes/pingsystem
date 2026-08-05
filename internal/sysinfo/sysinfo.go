package sysinfo

import (
	"os"
	"os/user"
	"strings"
	"time"
)

// Info holds system details required for notifications.
type Info struct {
	DeviceAlias    string
	ComputerName   string
	Username       string
	Timestamp      string
	WindowsVersion string
}

// GetInfo collects system details for the current host.
func GetInfo(deviceAlias string) Info {
	host, err := os.Hostname()
	if err != nil || host == "" {
		host = os.Getenv("COMPUTERNAME")
		if host == "" {
			host = "Unknown-PC"
		}
	}

	username := os.Getenv("USERNAME")
	if username == "" {
		u, err := user.Current()
		if err == nil && u.Username != "" {
			// user.Current() on Windows may return "DOMAIN\Username", strip domain if present
			parts := strings.Split(u.Username, "\\")
			username = parts[len(parts)-1]
		} else {
			username = "Unknown-User"
		}
	}

	now := time.Now().Format("2006-01-02 15:04:05 MST")

	return Info{
		DeviceAlias:    deviceAlias,
		ComputerName:   host,
		Username:       username,
		Timestamp:      now,
		WindowsVersion: GetWindowsVersion(),
	}
}
