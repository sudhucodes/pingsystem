package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"

	"github.com/sudhucodes/pingsystem/internal/autostart"
	"github.com/sudhucodes/pingsystem/internal/config"
	"github.com/sudhucodes/pingsystem/internal/logger"
	"github.com/sudhucodes/pingsystem/internal/sysinfo"
	"github.com/sudhucodes/pingsystem/internal/telegram"
	"github.com/sudhucodes/pingsystem/internal/watcher"
)

const Version = "1.0.0"

func main() {
	installFlag := flag.Bool("install", false, "Enable Windows registry autostart")
	uninstallFlag := flag.Bool("uninstall", false, "Disable Windows registry autostart")
	statusFlag := flag.Bool("status", false, "Check autostart and configuration status")
	versionFlag := flag.Bool("version", false, "Display application version")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("PingSystem v%s (%s/%s)\n", Version, runtime.GOOS, runtime.GOARCH)
		os.Exit(0)
	}

	if *installFlag {
		if err := autostart.Enable(); err != nil {
			fmt.Printf("Error enabling autostart: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("PingSystem successfully registered for Windows autostart (HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Run).")
		os.Exit(0)
	}

	if *uninstallFlag {
		if err := autostart.Disable(); err != nil {
			fmt.Printf("Error disabling autostart: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("PingSystem successfully removed from Windows autostart.")
		os.Exit(0)
	}

	if *statusFlag {
		enabled, val, err := autostart.IsEnabled()
		fmt.Printf("PingSystem v%s Status:\n", Version)
		if err != nil {
			fmt.Printf("  Autostart Registry Status: Error checking (%v)\n", err)
		} else if enabled {
			fmt.Printf("  Autostart Registry Status: ENABLED (%s)\n", val)
		} else {
			fmt.Println("  Autostart Registry Status: DISABLED")
		}

		cfgPath, _ := config.GetConfigPath()
		fmt.Printf("  Config File Path: %s\n", cfgPath)

		logDir, _ := logger.GetLogDir()
		fmt.Printf("  Logs Directory Path: %s\n", logDir)

		cfg, err := config.Load()
		if err != nil {
			fmt.Printf("  Config Status: ERROR (%v)\n", err)
		} else if err := cfg.Validate(); err != nil {
			fmt.Printf("  Config Status: INCOMPLETE (%v)\n", err)
		} else {
			fmt.Println("  Config Status: VALID")
		}
		os.Exit(0)
	}

	// Step 1: Initialize logger
	if err := logger.Init(); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Info("Starting PingSystem v%s...", Version)

	// Step 2: Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Error("Failed to load configuration: %v", err)
		os.Exit(1)
	}

	cfgPath, _ := config.GetConfigPath()
	logger.Info("Configuration loaded from %s", cfgPath)

	if err := cfg.Validate(); err != nil {
		logger.Warn("Configuration incomplete: %v. Please edit %s to configure your Telegram botToken and chatId.", err, cfgPath)
	}

	// Step 3: Telegram Client
	tgClient := telegram.NewClient(cfg.BotToken, cfg.ChatID)

	// Collect system metadata
	info := sysinfo.GetInfo(cfg.DeviceAlias)

	// Step 4: Send Startup / User Login notification
	if err := cfg.Validate(); err == nil {
		logger.Info("Sending Startup / User Login notification via Telegram...")
		if err := tgClient.SendEvent(telegram.EventStartup, info); err != nil {
			logger.Error("Failed to send Startup notification: %v", err)
		} else {
			logger.Info("Startup notification sent successfully.")
		}
	}

	// Step 5: Event callbacks for Sleep, Wake, Lock, Unlock, Shutdown
	callbacks := watcher.EventCallbacks{
		OnSleep: func() {
			logger.Info("Detected event: Windows Sleep")
			if err := cfg.Validate(); err == nil {
				currentInfo := sysinfo.GetInfo(cfg.DeviceAlias)
				if err := tgClient.SendEvent(telegram.EventSleep, currentInfo); err != nil {
					logger.Error("Failed to send Sleep notification: %v", err)
				} else {
					logger.Info("Sleep notification sent successfully.")
				}
			}
		},
		OnWake: func() {
			logger.Info("Detected event: Windows Wake")
			if err := cfg.Validate(); err == nil {
				currentInfo := sysinfo.GetInfo(cfg.DeviceAlias)
				if err := tgClient.SendEvent(telegram.EventWake, currentInfo); err != nil {
					logger.Error("Failed to send Wake notification: %v", err)
				} else {
					logger.Info("Wake notification sent successfully.")
				}
			}
		},
		OnLock: func() {
			logger.Info("Detected event: Windows Lock")
			if err := cfg.Validate(); err == nil {
				currentInfo := sysinfo.GetInfo(cfg.DeviceAlias)
				if err := tgClient.SendEvent(telegram.EventLock, currentInfo); err != nil {
					logger.Error("Failed to send Lock notification: %v", err)
				} else {
					logger.Info("Lock notification sent successfully.")
				}
			}
		},
		OnUnlock: func() {
			logger.Info("Detected event: Windows Unlock")
			if err := cfg.Validate(); err == nil {
				currentInfo := sysinfo.GetInfo(cfg.DeviceAlias)
				if err := tgClient.SendEvent(telegram.EventUnlock, currentInfo); err != nil {
					logger.Error("Failed to send Unlock notification: %v", err)
				} else {
					logger.Info("Unlock notification sent successfully.")
				}
			}
		},
		OnShutdown: func() {
			logger.Info("Detected event: Windows Shutdown")
			if err := cfg.Validate(); err == nil {
				currentInfo := sysinfo.GetInfo(cfg.DeviceAlias)
				if err := tgClient.SendEvent(telegram.EventShutdown, currentInfo); err != nil {
					logger.Error("Failed to send Shutdown notification: %v", err)
				} else {
					logger.Info("Shutdown notification sent successfully.")
				}
			}
		},
	}

	// Handle process termination signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)

	eventWatcher := watcher.New(callbacks)

	go func() {
		sig := <-sigChan
		logger.Info("Received signal %v, shutting down...", sig)
		eventWatcher.Stop()
		os.Exit(0)
	}()

	// Step 6: Start event watcher loop
	if runtime.GOOS == "windows" {
		logger.Info("PingSystem background agent running. Listening for Windows events...")
		if err := eventWatcher.Start(); err != nil {
			logger.Error("Event watcher error: %v", err)
		}
	} else {
		logger.Info("PingSystem agent initialized on non-Windows OS (%s). (Win32 event watcher inactive)", runtime.GOOS)
		// On non-Windows OS, block on signal
		<-sigChan
	}
}

// Suppress unused imports
var _ = filepath.Join
