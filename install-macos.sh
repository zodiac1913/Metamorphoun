#!/bin/bash
# Metamorphoun macOS Installer

set -e

echo "========================================"
echo "Metamorphoun macOS Installer"
echo "========================================"
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "ERROR: Go is not installed"
    echo "Please install Go from https://golang.org/dl/"
    echo "Or use Homebrew: brew install go"
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
echo "2. Install to /usr/local/bin:"
echo "   sudo cp metamorphoun /usr/local/bin/"
echo ""
echo "3. Add to Login Items:"
echo "   - Open System Preferences > Users & Groups"
echo "   - Click Login Items tab"
echo "   - Click + and add metamorphoun"
echo ""
echo "Note: You may need to grant permissions in:"
echo "  System Preferences > Security & Privacy > Privacy"
echo "  - Full Disk Access"
echo "  - Accessibility (if needed)"
echo ""
echo "Installation complete!"
