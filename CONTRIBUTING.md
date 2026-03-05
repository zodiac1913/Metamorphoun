# Contributing to Metamorphoun

Thank you for your interest in contributing to Metamorphoun! We welcome contributions from the community.

## How to Contribute

### Reporting Bugs

If you find a bug, please open an issue on GitHub with:
- A clear description of the problem
- Steps to reproduce the issue
- Your operating system and Go version
- Any relevant error messages or logs

### Suggesting Features

We love new ideas! Please open an issue to discuss:
- What problem the feature would solve
- How you envision it working
- Any implementation ideas you have

### Submitting Pull Requests

1. **Fork the repository** and create your branch from `main`
2. **Make your changes** with clear, descriptive commits
3. **Test your changes** on your platform (Windows/Linux/macOS)
4. **Update documentation** if you're changing functionality
5. **Submit a pull request** with a clear description of your changes

### Development Setup

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/Metamorphoun.git
cd Metamorphoun

# Install dependencies
go mod download

# Build
go build -o Metamorphoun

# Run
./Metamorphoun
```

### Code Style

- Follow standard Go conventions (use `go fmt`)
- Write clear, self-documenting code
- Add comments for complex logic
- Keep functions focused and concise

### Testing

Before submitting:
- Test on your target platform
- Verify wallpaper changing works
- Check that the web UI loads correctly
- Ensure quotes display properly

## Questions?

Feel free to open an issue for any questions about contributing!

## Code of Conduct

Be respectful, inclusive, and constructive in all interactions.
