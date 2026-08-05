@echo off
echo Building PingSystem.exe for Windows...
go build -ldflags="-s -w" -o PingSystem.exe .
if %ERRORLEVEL% EQU 0 (
    echo Build successful: PingSystem.exe
) else (
    echo Build failed!
)
