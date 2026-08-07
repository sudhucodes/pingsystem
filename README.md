# PingSystem (Go Edition)

[![CI](https://github.com/sudhucodes/pingsystem/actions/workflows/ci.yml/badge.svg)](https://github.com/sudhucodes/pingsystem/actions/workflows/ci.yml)
[![License: Apache 2.0](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go)](go.mod)

**PingSystem** is a lightweight, production-ready, native Windows background agent written entirely in Go. It monitors system power and session events and dispatches notifications to Telegram via the Telegram Bot API.

It compiles into a single standalone executable (`PingSystem.exe`) with no runtime dependencies (no Node.js, no Python, no PowerShell polling).

---

## Features (v1)

-   🚀 **Windows Startup / User Login Notification**
-   🌙 **Windows Sleep Notification** (`WM_POWERBROADCAST` -> `PBT_APMSUSPEND`)
-   ☀️ **Windows Wake Notification** (`WM_POWERBROADCAST` -> `PBT_APMRESUMEAUTOMATIC` / `PBT_APMRESUMESUSPEND`)
-   🔒 **Windows Lock Notification** (`WM_WTSSESSION_CHANGE` -> `WTS_SESSION_LOCK`, e.g. `Win + L`)
-   🔓 **Windows Unlock Notification** (`WM_WTSSESSION_CHANGE` -> `WTS_SESSION_UNLOCK`)
-   🛑 **Windows Shutdown Notification** (`WM_QUERYENDSESSION` / `WM_ENDSESSION` - best effort)
-   ⚙️ **Autostart Helper Flags**: Simple command-line flags to enable/disable Windows registry autostart (`HKCU\Software\Microsoft\Windows\CurrentVersion\Run`).
-   📁 **Daily Logging**: Automatically maintains log files in `%APPDATA%\PingSystem\logs\YYYY-MM-DD.log`.
-   📦 **Single Binary Output**: Compiles into `PingSystem.exe` (~6MB background process).

---

## Configuration

Configuration is automatically generated on first launch at:

```text
%APPDATA%\PingSystem\config.json
```

Example `config.json`:

```json
{
    "botToken": "YOUR_TELEGRAM_BOT_TOKEN",
    "chatId": "YOUR_TELEGRAM_CHAT_ID",
    "deviceAlias": "My Laptop"
}
```

### Fields

-   `botToken`: Your Telegram Bot API token from [@BotFather](https://t.me/BotFather).
-   `chatId`: Your Telegram User or Group Chat ID.
-   `deviceAlias`: A custom friendly alias for the device (e.g. `"Workstation 01"`).

---

## CLI Options

```text
Usage of PingSystem.exe:
  -install
    	Enable Windows registry autostart
  -uninstall
    	Disable Windows registry autostart
  -status
    	Check autostart and configuration status
  -version
    	Display application version
```

---

## Telegram Notification Format

Each event sends a Telegram message formatted as follows:

```text
🚀 PingSystem Event: Startup / User Login

Device Alias: My Laptop
Username: Sudhu
Timestamp: 2026-08-05 10:30:00 PM MST
```

---

## Building from Source

### Requirements

-   Go 1.26+ installed

### Build Binary (Cross-compile on macOS / Linux)

```bash
./build.sh
```

Or manually:

```bash
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o PingSystem.exe .
```

### Build Binary (On Windows)

```cmd
build.bat
```

---

## Directory Structure

```text
.
├── main.go                       # Main entry point & CLI parser
├── go.mod                        # Go module definition
├── go.sum                        # Go module checksums
├── build.sh                      # Shell script to cross-compile for Windows
├── build.bat                     # Windows batch build script
├── LICENSE                       # Apache-2.0 Open Source License
├── CONTRIBUTING.md               # Contribution guidelines & release setup
├── CODE_OF_CONDUCT.md            # Contributor Covenant Code of Conduct
├── SECURITY.md                   # Security vulnerability reporting policy
├── .github/                      # Issue/PR templates & GitHub Actions CI/Release workflows
└── internal/
    ├── autostart/                # Registry autostart management (Enable, Disable, IsEnabled)
    ├── config/                   # Config loader & auto-creator (%APPDATA%\PingSystem\config.json)
    ├── logger/                   # Daily file logging (%APPDATA%\PingSystem\logs\)
    ├── sysinfo/                  # Native OS metadata (hostname, user, Windows build version)
    ├── telegram/                 # Telegram Bot API notification client
    └── watcher/                  # Win32 hidden message window event loop (WM_POWERBROADCAST, etc.)
```

---

## 🤝 Contributing

Contributions are welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) to get started with setup, coding guidelines, and submitting changesets.

Please also adhere to our [Code of Conduct](CODE_OF_CONDUCT.md).

---

## 🔒 Security

If you discover a security vulnerability, please refer to our [Security Policy](SECURITY.md) for instructions on how to report it confidentially.

---

## 📄 License

This project is licensed under the **Apache License 2.0 (Apache-2.0)** - see the [LICENSE](LICENSE) file for details.
