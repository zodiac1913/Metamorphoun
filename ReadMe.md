# Metamorphoun

Metamorphoun is a flexible wallpaper changer written in Go. It can automatically change your desktop wallpaper using images from pre-configured websites or local folders. The app also supports displaying inspirational or custom quotes, which you can add in JSON format.

## Features

- Change wallpapers from online sources or local folders
- Display quotes as overlays on your wallpaper
- Add your own quotes using a simple JSON file
- System tray integration (Windows)
- Highly configurable for different use cases
- Inspired by [Variety](https://github.com/peterlevi/variety) by Peter Levi

## Adding Quotes

You can add your own quotes by uploading a JSON file in the following format:

```json
[
  {"statement": "The only limit to our realization of tomorrow is our doubts of today.", "author": "Franklin D. Roosevelt"},
  {"statement": "In the middle of every difficulty lies opportunity.", "author": "Albert Einstein"}
]
```

## Building the Executable

### Prerequisites

- [Go](https://golang.org/dl/) 1.16 or newer installed
- Windows, Linux, or Mac OS development environment

---

### Build for Windows

#### Console Shown

```sh
go build -o Metamorphoun.exe
```

#### No Console (GUI/System Tray Only)

```sh
go build -ldflags="-H=windowsgui" -o Metamorphoun.exe
```

---

### Build for Linux

```sh
go build -o metamorphoun
```

- The resulting binary can be run from the terminal or added to startup scripts.
- Desktop integration (wallpaper changing, tray icon) may require additional packages depending on your desktop environment (e.g., GNOME, KDE, XFCE).

---

### Build for Mac OS

```sh
go build -o metamorphoun
```

- The resulting binary can be run from the terminal.
- Wallpaper changing and tray integration may require additional permissions or helper tools on Mac OS.
- You may need to grant the app access to "Full Disk Access" or "Accessibility" in System Preferences for some features.

---

## Running

Double-click the executable or run it from the command line:

**Windows:**
```sh
Metamorphoun.exe
```

**Linux/Mac OS:**
```sh
./metamorphoun
```

On first run, configuration files and folders will be created in your user directory.

---

## License

This project is based on the ideas of [Variety](https://github.com/peterlevi/variety) by Peter Levi.  
See LICENSE