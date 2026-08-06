package sysinfo

import (
	"os"
	"os/user"
	"strings"
	"time"
)

// Info holds system details required for notifications.
type Info struct {
	DeviceAlias string
	Username    string
	Timestamp   string
}

// GetInfo collects system details for the current host.
func GetInfo(deviceAlias string) Info {
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

	now := time.Now().Format("2006-01-02 03:04:05 PM MST")

	return Info{
		DeviceAlias: deviceAlias,
		Username:    username,
		Timestamp:   now,
	}
}
