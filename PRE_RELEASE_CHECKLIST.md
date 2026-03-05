# Pre-Release Checklist for Metamorphoun

Use this checklist before making your repository public.

## 🔒 Security & Privacy

- [ ] Remove any hardcoded credentials, API keys, or tokens
- [ ] Review all files for personal information (paths, usernames, etc.)
- [ ] Check config.json is in .gitignore (user-specific data)
- [ ] Remove any sensitive test data or screenshots with personal info
- [ ] Review commit history for accidentally committed secrets

## 📁 Repository Cleanup

- [ ] Delete temporary files (on_exit_*.txt, debug files)
- [ ] Remove compiled binaries (*.exe, metamorphoun)
- [ ] Clean up test images and user-specific content
- [ ] Remove or document any hardcoded local paths in code
- [ ] Delete unnecessary files (rendered_page.html, etc.)

## 📝 Documentation

- [ ] Update README.md with correct GitHub username/repo URL
- [ ] Replace YOUR_USERNAME in README.md badges
- [ ] Add screenshots to docs/ folder (optional but recommended)
- [ ] Verify all installation instructions work
- [ ] Check that LICENSE file is present
- [ ] Review CONTRIBUTING.md guidelines
- [ ] Add CHANGELOG.md for version history (optional)

## 🧪 Testing

- [ ] Test build on Windows
- [ ] Test build on Linux (if possible)
- [ ] Test build on macOS (if possible)
- [ ] Verify installer scripts work
- [ ] Test web interface loads correctly
- [ ] Verify wallpaper changing works
- [ ] Test quote overlay functionality
- [ ] Check system tray integration (Windows)

## 🏗️ GitHub Setup

- [ ] Create repository on GitHub (keep private initially)
- [ ] Push code to GitHub
- [ ] Add repository description
- [ ] Add topics/tags (go, wallpaper, desktop, golang, etc.)
- [ ] Set up GitHub Pages (optional - for documentation)
- [ ] Enable Issues
- [ ] Enable Discussions (optional)
- [ ] Create initial release (v1.0.0)
- [ ] Test GitHub Actions workflow runs successfully

## 📦 Release Preparation

- [ ] Tag first release: `git tag v1.0.0`
- [ ] Push tags: `git push origin v1.0.0`
- [ ] Verify GitHub Actions creates release artifacts
- [ ] Download and test release binaries
- [ ] Write release notes
- [ ] Create release on GitHub with binaries attached

## 🌐 Going Public

- [ ] Make repository public in GitHub settings
- [ ] Share on relevant communities (r/golang, r/unixporn, etc.)
- [ ] Tweet/post about the release (optional)
- [ ] Add to awesome-go lists (optional)
- [ ] Submit to package managers (optional - Homebrew, AUR, etc.)

## 📋 Post-Release

- [ ] Monitor issues and respond to questions
- [ ] Set up project board for tracking features/bugs
- [ ] Consider adding CODE_OF_CONDUCT.md
- [ ] Set up branch protection rules
- [ ] Add SECURITY.md for vulnerability reporting

## 🔍 Final Review Commands

Run these before going public:

```bash
# Check for secrets
git log --all --full-history --source --find-object=<path-to-sensitive-file>

# Review all tracked files
git ls-files

# Check what will be committed
git status

# Review .gitignore is working
git check-ignore -v *

# Clean untracked files (BE CAREFUL!)
git clean -n  # dry run first
```

## ✅ Ready to Go Public?

Once all items are checked:

1. Go to GitHub repository Settings
2. Scroll to "Danger Zone"
3. Click "Change visibility"
4. Select "Make public"
5. Confirm the action

**Congratulations! Your project is now open source! 🎉**
