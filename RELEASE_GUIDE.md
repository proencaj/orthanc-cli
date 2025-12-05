# Release Guide

Complete guide for creating releases of Orthanc CLI using GoReleaser.

## Table of Contents

- [Overview](#overview)
- [Prerequisites](#prerequisites)
- [First-Time Setup](#first-time-setup)
- [Creating a Release](#creating-a-release)
- [Release Checklist](#release-checklist)
- [What Gets Released](#what-gets-released)
- [Package Distribution](#package-distribution)
- [Troubleshooting](#troubleshooting)

---

## Overview

Orthanc CLI uses [GoReleaser](https://goreleaser.com) for automated releases. When you push a git tag, GitHub Actions automatically:

- ✅ Runs all tests
- ✅ Builds binaries for multiple platforms (Linux/macOS, amd64/arm64)
- ✅ Creates distribution packages (DEB, RPM, APK)
- ✅ Generates checksums
- ✅ Creates a GitHub Release with all artifacts
- ✅ Updates the Homebrew tap

**Prerelease Policy**: All versions before v1.0.0 (v0.x.x) are automatically marked as prereleases.

---

## Prerequisites

### Required for All Releases

- Go 1.25 or later
- Git with all changes committed
- Updated CHANGELOG.md
- All tests passing

### Required for Automated Releases (GitHub Actions)

1. **Homebrew Tap Repository**
   - Create a public repository: `proencaj/homebrew-orthanc-cli`
   - GitHub: https://github.com/new

2. **GitHub Personal Access Token**
   - Go to: https://github.com/settings/tokens
   - Create token with `repo` scope
   - Name it: `HOMEBREW_TAP_GITHUB_TOKEN`

3. **Add Token to Repository Secrets**
   - Go to: https://github.com/proencaj/orthanc-cli/settings/secrets/actions
   - Add secret: `HOMEBREW_TAP_GITHUB_TOKEN`

### Optional: Package Signing

For signed DEB/RPM packages (recommended for production):

```bash
# Install GPG (macOS)
brew install gnupg

# Generate GPG key
gpg --full-generate-key
# Choose: RSA and RSA, 4096 bits, your name and email

# List keys to find KEY_ID
gpg --list-secret-keys --keyid-format=long

# Export private key
gpg --export-secret-keys YOUR_KEY_ID > gpg-private-key.asc

# Add to GitHub Secrets as GPG_KEY_FILE
# (Paste entire contents of gpg-private-key.asc)
```

**Note**: Packages will be unsigned if `GPG_KEY_FILE` is not set (fine for development).

---

## First-Time Setup

### Install GoReleaser (for local testing)

```bash
# macOS/Linux
brew install goreleaser

# Or use the Makefile
make install-goreleaser
```

### Verify Setup

```bash
# Test GoReleaser configuration
make goreleaser-test

# Build a local snapshot (no tag required)
make goreleaser-snapshot

# Check the dist/ directory
ls -lh dist/
```

---

## Creating a Release

### Quick Release (Automated)

```bash
# 1. Ensure all changes are committed
git status

# 2. Update CHANGELOG.md
vim CHANGELOG.md
git add CHANGELOG.md
git commit -m "docs: update CHANGELOG for v0.2.0"
git push origin main

# 3. Create and push tag
git tag -a v0.2.0 -m "Release v0.2.0"
git push origin v0.2.0

# 4. GitHub Actions handles the rest automatically!
# Monitor: https://github.com/proencaj/orthanc-cli/actions
```

### Step-by-Step Release Process

#### 1. Prepare the Release

Update CHANGELOG.md:
- Move items from `[Unreleased]` to new version section
- Add release date
- Update comparison links

```bash
git add CHANGELOG.md
git commit -m "docs: update CHANGELOG for v0.2.0"
git push origin main
```

#### 2. Validate Everything

```bash
# Run pre-release validation
make release-prepare

# This checks:
# - Working directory is clean
# - Required files exist
# - All tests pass
# - Code is formatted
# - No linting errors
```

#### 3. Create and Push Tag

```bash
# Set version
VERSION="v0.2.0"

# Create annotated tag
git tag -a $VERSION -m "Release $VERSION"

# Push tag (triggers GitHub Actions workflow)
git push origin $VERSION
```

#### 4. Monitor GitHub Actions

1. Go to: https://github.com/proencaj/orthanc-cli/actions
2. Watch the "Release" workflow (usually 2-5 minutes)
3. Wait for green checkmark ✅

#### 5. Verify Release

1. Check release page: https://github.com/proencaj/orthanc-cli/releases
2. Verify all artifacts are present:
   - Binary archives (linux/darwin, amd64/arm64)
   - DEB packages
   - RPM packages
   - APK packages
   - SHA256SUMS.txt
3. Check Homebrew tap: https://github.com/proencaj/homebrew-orthanc-cli

#### 6. Test Installation

```bash
# Test Homebrew (may take a few minutes to update)
brew upgrade orthanc-cli
orthanc --version

# Or test manual download
VERSION="v0.2.0"
curl -LO https://github.com/proencaj/orthanc-cli/releases/download/$VERSION/orthanc-cli-$VERSION-darwin-arm64.tar.gz
tar -xzf orthanc-cli-$VERSION-darwin-arm64.tar.gz
./orthanc --version
```

---

## Release Checklist

Use this checklist before creating a release:

### Pre-Release
- [ ] All changes committed and pushed to `main`
- [ ] Tests passing: `make test`
- [ ] Code formatted: `make fmt`
- [ ] Updated `CHANGELOG.md` with new version
- [ ] Version follows semantic versioning (e.g., v0.2.0)
- [ ] Run `make release-prepare` successfully

### Release
- [ ] Created annotated tag: `git tag -a v0.2.0 -m "Release v0.2.0"`
- [ ] Pushed tag: `git push origin v0.2.0`
- [ ] GitHub Actions workflow completed successfully
- [ ] All artifacts uploaded to GitHub Release

### Post-Release
- [ ] Verified release on GitHub
- [ ] Tested installation (Homebrew and/or manual)
- [ ] Updated documentation if needed
- [ ] Announced release (if applicable)
- [ ] Closed related issues/milestones

---

## What Gets Released

### Binary Archives (tar.gz)

- `orthanc-cli-VERSION-linux-amd64.tar.gz`
- `orthanc-cli-VERSION-linux-arm64.tar.gz`
- `orthanc-cli-VERSION-darwin-amd64.tar.gz`
- `orthanc-cli-VERSION-darwin-arm64.tar.gz`

Each archive contains:
- `orthanc` - Binary executable
- `README.md` - Documentation
- `LICENSE` - License file
- `CHANGELOG.md` - Changelog

### Linux Packages

**Debian/Ubuntu (DEB):**
- `orthanc-cli_VERSION_linux_amd64.deb`
- `orthanc-cli_VERSION_linux_arm64.deb`

**CentOS/RHEL/Fedora (RPM):**
- `orthanc-cli_VERSION_linux_amd64.rpm`
- `orthanc-cli_VERSION_linux_arm64.rpm`

**Alpine Linux (APK):**
- `orthanc-cli_VERSION_linux_amd64.apk`
- `orthanc-cli_VERSION_linux_arm64.apk`

### Other Artifacts

- `SHA256SUMS.txt` - Checksums for all artifacts
- Homebrew formula (auto-updated in `proencaj/homebrew-orthanc-cli`)

### Release Notes

GitHub Release includes:
- Installation instructions (Homebrew, DEB, RPM, APK, manual)
- Changelog from commits (auto-generated)
- Download links for all artifacts
- Checksum verification instructions

---

## Package Distribution

### Homebrew (macOS/Linux)

```bash
brew install proencaj/orthanc-cli
```

### Debian/Ubuntu

```bash
wget https://github.com/proencaj/orthanc-cli/releases/download/v0.2.0/orthanc-cli_0.2.0_linux_amd64.deb
sudo dpkg -i orthanc-cli_0.2.0_linux_amd64.deb
```

### CentOS/RHEL/Fedora

```bash
wget https://github.com/proencaj/orthanc-cli/releases/download/v0.2.0/orthanc-cli_0.2.0_linux_amd64.rpm
sudo rpm -i orthanc-cli_0.2.0_linux_amd64.rpm
# or
sudo dnf install orthanc-cli_0.2.0_linux_amd64.rpm
```

### Alpine Linux

```bash
wget https://github.com/proencaj/orthanc-cli/releases/download/v0.2.0/orthanc-cli_0.2.0_linux_amd64.apk
sudo apk add --allow-untrusted orthanc-cli_0.2.0_linux_amd64.apk
```

### Manual Installation

```bash
# Download
curl -LO https://github.com/proencaj/orthanc-cli/releases/latest/download/orthanc-cli-VERSION-linux-amd64.tar.gz

# Extract
tar -xzf orthanc-cli-VERSION-linux-amd64.tar.gz

# Install
sudo mv orthanc /usr/local/bin/
```

### Verify Downloads

```bash
# Download checksums
wget https://github.com/proencaj/orthanc-cli/releases/download/v0.2.0/SHA256SUMS.txt

# Verify package
sha256sum -c SHA256SUMS.txt --ignore-missing
```

---

## Troubleshooting

### Release Failed - Working Directory Not Clean

```bash
git status
git add .
git commit -m "fix: pre-release cleanup"
```

### Release Failed - Tests Failed

```bash
# Run tests locally to see what's failing
make test

# Fix issues and re-tag
git tag -d v0.2.0  # Delete local tag
git push origin :refs/tags/v0.2.0  # Delete remote tag
# Fix, commit, and try again
```

### GitHub Actions Failed - "resource not accessible"

Ensure workflow has write permissions in `.github/workflows/release.yml`:

```yaml
permissions:
  contents: write
  packages: write
```

### Homebrew Tap Not Updated

1. Check `HOMEBREW_TAP_GITHUB_TOKEN` secret exists
2. Verify token has `repo` scope
3. Confirm `homebrew-orthanc-cli` repository exists and is public
4. Check GitHub Actions logs for errors

### Want to Delete/Redo a Release

```bash
# Delete tag locally
git tag -d v0.2.0

# Delete tag remotely
git push origin :refs/tags/v0.2.0

# Delete GitHub Release (via web UI or gh CLI)
gh release delete v0.2.0

# Fix issues, then re-create tag and push
```

### Local Testing Without Creating Release

```bash
# Build snapshot (no tag required)
make goreleaser-snapshot

# Output goes to dist/ directory
ls -lh dist/
```

### GPG Signing Errors

If you get GPG errors during local testing:

```bash
# Option 1: Test without signing
goreleaser release --snapshot --clean --skip=publish

# Option 2: Set up GPG key locally
export GPG_KEY_FILE=/path/to/gpg-private-key.asc
```

### Version Shows as "dev" or Wrong Version

Ensure you've created and checked out the tag:

```bash
git tag -a v0.2.0 -m "Release v0.2.0"
git checkout v0.2.0

# Verify version in binary
make build
./bin/orthanc --version
```

---

## Quick Commands Reference

```bash
# Check what would be released
git log $(git describe --tags --abbrev=0)..HEAD --oneline

# View last release tag
git describe --tags --abbrev=0

# List all tags
git tag -l

# Delete local tag
git tag -d v0.2.0

# Delete remote tag
git push origin :refs/tags/v0.2.0

# Test GoReleaser config
make goreleaser-test

# Build local snapshot
make goreleaser-snapshot

# Run pre-release validation
make release-prepare
```

---

## Version Numbering

This project follows [Semantic Versioning](https://semver.org):

- **MAJOR** version (v1.0.0) - Incompatible API changes
- **MINOR** version (v0.2.0) - New functionality, backwards compatible
- **PATCH** version (v0.2.1) - Backwards compatible bug fixes

**Prerelease Policy**: All v0.x.x versions are marked as prereleases until v1.0.0.

---

## Additional Resources

- [GoReleaser Documentation](https://goreleaser.com)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Homebrew Tap Documentation](https://docs.brew.sh/How-to-Create-and-Maintain-a-Tap)
- [Semantic Versioning](https://semver.org)

---

## Support

If you encounter issues with the release process:

1. Check GitHub Actions logs for detailed errors
2. Test locally with `make goreleaser-snapshot`
3. Review this guide and [.goreleaser.yaml](.goreleaser.yaml)
4. Open an issue: https://github.com/proencaj/orthanc-cli/issues
