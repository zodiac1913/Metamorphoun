#!/bin/bash
# Metamorphoun Linux Installer

set -e

echo "========================================"
echo "Metamorphoun Linux Installer"
echo "========================================"
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "ERROR: Go is not installed"
    echo "Please install Go from https://golang.org/dl/"
    echo "Or use your package manager:"
    echo "  Ubuntu/Debian: sudo apt install golang-go"
    echo "  Fedora: sudo dnf install golang"
    echo "  Arch: sudo pacman -S go"
    exit 1
fi

echo "[1/5] Checking Go installation..."
go version
echo ""

echo "[2/5] Downloading dependencies..."
go mod download
echo ""

echo "[3/5] Building Metamorphoun..."
go build -o metamorphoun
echo ""

echo "[4/5] Making executable..."
chmod +x metamorphoun
echo ""

echo "[5/5] Installation options..."
echo ""
echo "Build complete! metamorphoun binary created."
echo ""
echo "Installation options:"
echo ""
echo "1. Run from current directory:"
echo "   ./metamorphoun"
echo ""
echo "2. Install system-wide (requires sudo):"
echo "   sudo cp metamorphoun /usr/local/bin/"
echo ""
echo "3. Add to autostart:"
echo "   mkdir -p ~/.config/autostart"
echo "   cat > ~/.config/autostart/metamorphoun.desktop << EOF"
echo "[Desktop Entry]"
echo "Type=Application"
echo "Name=Metamorphoun"
echo "Exec=$(pwd)/metamorphoun"
echo "Hidden=false"
echo "NoDisplay=false"
echo "X-GNOME-Autostart-enabled=true"
echo "EOF"
echo ""
echo "Installation complete!"
