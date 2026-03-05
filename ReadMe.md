# Metamorphoun

[![Build Status](https://github.com/zodiac1913/Metamorphoun/workflows/Build%20and%20Release/badge.svg)](https://github.com/zodiac1913/Metamorphoun/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go)](https://golang.org/)

Metamorphoun is a flexible wallpaper changer written in Go. It automatically changes your desktop wallpaper using images from pre-configured websites or local folders, with support for inspirational quotes as overlays.

## ✨ Features

- 🖼️ **Multiple Image Sources**: Bing, NASA, Unsplash, Flickr, PicSum, or local folders
- 💬 **Quote Overlays**: Display inspirational quotes on your wallpaper
- 🎨 **Image Filters**: Blur, monochrome, oil painting, vortex, and more
- 🌐 **Web Interface**: Configure settings through an intuitive web UI at `http://localhost:8080`
- 🔄 **Auto-Change**: Configurable timer for automatic wallpaper rotation
- 💾 **History**: Keep track of previous wallpapers
- 🎯 **System Tray**: Easy access on Windows
- 📝 **Custom Quotes**: Add your own quotes via JSON
- 🖥️ **Cross-Platform**: Windows, Linux, and macOS support

Inspired by [Variety](https://github.com/peterlevi/variety) by Peter Levi.

## 📥 Installation

### Windows

**Option 1: Download Pre-built Binary**
1. Download `Metamorphoun-windows-amd64.exe` from [Releases](https://github.com/zodiac1913/Metamorphoun/releases)
2. Rename to `Metamorphoun.exe` and run

**Option 2: Build from Source**
```cmd
git clone https://github.com/zodiac1913/Metamorphoun.git
cd Metamorphoun
install-windows.bat
```

### Linux

```bash
git clone https://github.com/zodiac1913/Metamorphoun.git
cd Metamorphoun
chmod +x install-linux.sh
./install-linux.sh
```

Desktop integration may require additional packages depending on your environment (GNOME, KDE, XFCE).

### macOS

```bash
git clone https://github.com/zodiac1913/Metamorphoun.git
cd Metamorphoun
chmod +x install-macos.sh
./install-macos.sh
```

You may need to grant permissions in System Preferences > Security & Privacy:
- Full Disk Access
- Accessibility (if needed)

## 🚀 Quick Start

1. **Run the application**:
   - **Windows**: Double-click `Metamorphoun.exe` or find it in the system tray
   - **Linux/macOS**: `./metamorphoun`

2. **Access the web interface**: Opens automatically at `http://localhost:8080`

3. **Configure your preferences**:
   - Choose image sources (online or local folders)
   - Set wallpaper change interval
   - Enable/disable quote overlays
   - Apply image filters

On first run, configuration files and folders will be created in your user directory.

## 📝 Adding Custom Quotes

Create a JSON file with your quotes:

```json
[
  {
    "statement": "The only limit to our realization of tomorrow is our doubts of today.",
    "author": "Franklin D. Roosevelt"
  },
  {
    "statement": "In the middle of every difficulty lies opportunity.",
    "author": "Albert Einstein"
  }
]
```

Upload via the web interface under Quote Tools.

## 🛠️ Building from Source

### Prerequisites

- [Go](https://golang.org/dl/) 1.23 or newer
- Git

### Build Commands

**Windows (GUI mode - no console):**
```cmd
go build -ldflags="-H=windowsgui" -o Metamorphoun.exe
```

**Windows (with console for debugging):**
```cmd
go build -o Metamorphoun.exe
```

**Linux:**
```bash
go build -o metamorphoun
```

**macOS:**
```bash
go build -o metamorphoun
```

## 📂 Configuration

Configuration is stored in:
- **Windows**: `C:\Users\<username>\.Metamorphoun\`
- **Linux/macOS**: `~/.Metamorphoun/`

The `config.json` file contains all settings and is managed through the web interface.

## 🎨 Image Sources

### Built-in Online Sources
- **Bing**: Daily photo from Bing
- **NASA**: Astronomy Picture of the Day
- **Unsplash**: High-quality random photos
- **Flickr**: Curated collections
- **PicSum**: Random placeholder images

### Local Folders
Add any folder containing images through the web interface.

## 🖼️ Image Filters

- Original (no filter)
- Blur (soft/hard)
- Monochrome
- Oil painting effect
- Vortex distortion
- And more!

## 🤝 Contributing

Contributions are welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## 📜 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Peter Levi for the original [Variety](https://github.com/peterlevi/variety) wallpaper changer
- All the open-source libraries used in this project
- Contributors and users of Metamorphoun

## 📞 Support

- 🐛 [Report a Bug](https://github.com/zodiac1913/Metamorphoun/issues)
- 💡 [Request a Feature](https://github.com/zodiac1913/Metamorphoun/issues)

## 📋 Version History

**v2026.03.05** R
- Initial public release
- Multi-platform support (Windows, Linux, macOS)
- Web-based configuration interface
- Quote overlay system
- Image filters and effects

**v2024-10-14**
- Made changes to avoid distortions to Christian PD pics. Vortex and Dali are not respecting to Adonai.

---

**Made with ❤️ by the Metamorphoun community**
