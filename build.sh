#!/usr/bin/env bash
set -e

echo "Building PingSystem.exe for Windows (amd64)..."
export PATH=$PATH:/usr/local/go/bin:/opt/homebrew/bin
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o PingSystem.exe .
echo "Build complete: PingSystem.exe ($(du -h PingSystem.exe | cut -f1))"
