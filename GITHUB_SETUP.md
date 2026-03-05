# GitHub Repository Setup Guide

Step-by-step guide to set up your Metamorphoun repository on GitHub.

## Step 1: Create GitHub Repository

1. Go to https://github.com/new
2. Repository name: `Metamorphoun`
3. Description: `A flexible wallpaper changer with quote overlays, written in Go`
4. **Keep it PRIVATE initially** (we'll make it public after review)
5. Do NOT initialize with README (we already have one)
6. Click "Create repository"

## Step 2: Prepare Local Repository

Before pushing, clean up your local repository:

```bash
# Navigate to your project
cd d:\Users\Dominic\Desktop\Dev\Metamorphoun

# Remove user-specific files
del on_exit_*.txt
del rendered_page.html
del *.exe
del __debug_bin*

# Review what will be committed
git status

# If git not initialized yet:
git init
git add .
git commit -m "Initial commit: Metamorphoun wallpaper changer"
```

## Step 3: Connect to GitHub

Replace `zodiac1913` with your actual GitHub username:

```bash
git remote add origin https://github.com/zodiac1913/Metamorphoun.git
git branch -M main
git push -u origin main
```

## Step 4: Configure Repository Settings

On GitHub, go to your repository settings:

### General Settings
- Add description: "A flexible wallpaper changer with quote overlays, written in Go"
- Add website: (your project website if you have one)
- Add topics: `go`, `golang`, `wallpaper`, `desktop`, `wallpaper-changer`, `quotes`, `cross-platform`

### Features
- ✅ Issues
- ✅ Discussions (optional but recommended)
- ✅ Projects (optional)
- ✅ Wiki (optional)

### Pull Requests
- ✅ Allow squash merging
- ✅ Allow rebase merging
- ✅ Automatically delete head branches

## Step 5: Create First Release

```bash
# Tag your first release
git tag -a v1.0.0 -m "Initial release of Metamorphoun"
git push origin v1.0.0
```

This will trigger the GitHub Actions workflow to build binaries for all platforms.

## Step 6: Verify GitHub Actions

1. Go to the "Actions" tab in your repository
2. Wait for the build to complete (usually 5-10 minutes)
3. Check that all three platforms (Windows, Linux, macOS) build successfully
4. If successful, binaries will be attached to the release

## Step 7: Create Release Notes

1. Go to "Releases" tab
2. Click on the v1.0.0 tag
3. Click "Edit release"
4. Add release notes:

```markdown
# Metamorphoun v1.0.0 - Initial Release

First public release of Metamorphoun, a flexible wallpaper changer written in Go!

## Features
- 🖼️ Multiple image sources (Bing, NASA, Unsplash, Flickr, local folders)
- 💬 Quote overlays with custom quotes support
- 🎨 Image filters (blur, monochrome, oil painting, vortex)
- 🌐 Web-based configuration interface
- 🔄 Automatic wallpaper rotation
- 💾 Wallpaper history
- 🎯 System tray integration (Windows)
- 🖥️ Cross-platform support

## Installation

Download the appropriate binary for your platform:
- Windows: `Metamorphoun-windows-amd64.exe`
- Linux: `metamorphoun-linux-amd64`
- macOS: `metamorphoun-macos-amd64`

See the [README](https://github.com/zodiac1913/Metamorphoun) for detailed installation instructions.

## Known Issues
- None yet! Please report any issues you find.

## Credits
Inspired by [Variety](https://github.com/peterlevi/variety) by Peter Levi.
```

5. Click "Publish release"

## Step 8: Review Before Going Public

Use the [PRE_RELEASE_CHECKLIST.md](PRE_RELEASE_CHECKLIST.md) to ensure everything is ready.

Key things to check:
- [ ] No personal information in code or commits
- [ ] No API keys or credentials
- [ ] All paths are relative or configurable
- [ ] README has correct URLs
- [ ] License is in place
- [ ] Builds work on all platforms

## Step 9: Make Repository Public

Once you've verified everything:

1. Go to Settings
2. Scroll to "Danger Zone"
3. Click "Change visibility"
4. Select "Make public"
5. Type the repository name to confirm
6. Click "I understand, make this repository public"

## Step 10: Share Your Project

Now that it's public, share it:

### Reddit
- r/golang - "Show & Tell" posts
- r/unixporn - Desktop customization
- r/Windows10 - Windows users
- r/linux - Linux users

### Social Media
- Twitter/X with hashtags: #golang #opensource #wallpaper
- LinkedIn
- Dev.to blog post

### Listings
- Submit to [Awesome Go](https://github.com/avelino/awesome-go)
- Add to [Go Projects](https://github.com/golang/go/wiki/Projects)

### Package Managers (Future)
- Homebrew (macOS)
- AUR (Arch Linux)
- Chocolatey (Windows)
- Snap/Flatpak (Linux)

## Maintenance Tips

### Responding to Issues
- Be friendly and welcoming
- Thank people for reporting issues
- Ask for more details if needed
- Close issues when resolved

### Pull Requests
- Review code carefully
- Test changes locally
- Provide constructive feedback
- Thank contributors

### Versioning
Use semantic versioning (MAJOR.MINOR.PATCH):
- MAJOR: Breaking changes
- MINOR: New features (backward compatible)
- PATCH: Bug fixes

Example:
```bash
git tag -a v1.1.0 -m "Added multi-monitor support"
git push origin v1.1.0
```

## Need Help?

- GitHub Docs: https://docs.github.com
- Open Source Guide: https://opensource.guide
- Choose a License: https://choosealicense.com

Good luck with your open source project! 🚀
