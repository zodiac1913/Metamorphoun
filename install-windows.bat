@echo off
REM Metamorphoun Windows Installer
REM This script builds and installs Metamorphoun on Windows

echo ========================================
echo Metamorphoun Windows Installer
echo ========================================
echo.

REM Check if Go is installed
where go >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo ERROR: Go is not installed or not in PATH
    echo Please install Go from https://golang.org/dl/
    pause
    exit /b 1
)

echo [1/4] Checking Go installation...
go version
echo.

echo [2/4] Downloading dependencies...
go mod download
if %ERRORLEVEL% NEQ 0 (
    echo ERROR: Failed to download dependencies
    pause
    exit /b 1
)
echo.

echo [3/4] Building Metamorphoun...
go build -ldflags="-H=windowsgui" -o Metamorphoun.exe
if %ERRORLEVEL% NEQ 0 (
    echo ERROR: Build failed
    pause
    exit /b 1
)
echo.

echo [4/4] Build complete!
echo.
echo Metamorphoun.exe has been created in the current directory.
echo.
echo To run Metamorphoun:
echo   - Double-click Metamorphoun.exe
echo   - Or run from command line: Metamorphoun.exe
echo.
echo To add to startup (optional):
echo   1. Press Win+R
echo   2. Type: shell:startup
echo   3. Create a shortcut to Metamorphoun.exe in that folder
echo.
echo Installation complete!
pause
