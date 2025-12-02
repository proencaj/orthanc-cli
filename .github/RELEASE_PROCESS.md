# Release Process

This document describes the process for creating a new release of Orthanc CLI.

## Prerequisites

- All changes committed and pushed to `main` branch
- CHANGELOG.md updated with the new version
- All tests passing
- Go 1.25 or later installed
- `gh` CLI installed (optional, for automated GitHub releases)

## Release Checklist

### 1. Prepare the Release

Update the CHANGELOG.md:
```bash
# Edit CHANGELOG.md
# - Move items from [Unreleased] to a new version section
# - Add the release date
# - Update the comparison links at the bottom
```

Commit the changelog:
```bash
git add CHANGELOG.md
git commit -m "docs: update CHANGELOG for v0.1.0"
git push origin main
```

### 2. Create Git Tag

```bash
# Create an annotated tag
git tag -a v0.1.0 -m "Release v0.1.0 - Initial public release"

# Push the tag
git push origin v0.1.0
```

### 3. Run Release Validation

The `release-prepare` target will validate everything:

```bash
make release-prepare
```

This checks:
- ✓ Working directory is clean (no uncommitted changes)
- ✓ Required files exist (README.md, LICENSE, CHANGELOG.md)
- ✓ Version tag is set (not "dev" or "unknown")
- ✓ All tests pass
- ✓ Code is formatted
- ✓ Go vet passes

If validation fails, fix the issues and try again.

### 4. Build Release Artifacts

```bash
# This will:
# - Run release-prepare validation
# - Build binaries for all platforms
# - Create tar.gz archives with README, LICENSE, and CHANGELOG
make release
```

This creates files in `bin/release/`:
- `orthanc-v0.1.0-linux-amd64.tar.gz`
- `orthanc-v0.1.0-linux-arm64.tar.gz`
- `orthanc-v0.1.0-darwin-amd64.tar.gz`
- `orthanc-v0.1.0-darwin-arm64.tar.gz`

### 5. Generate Checksums

```bash
make checksums
```

This creates `bin/release/SHA256SUMS.txt` with checksums for all archives.

### 6. Prepare Release Notes

Run the helper script:

```bash
./.github/prepare-release.sh v0.1.0
```

This generates:
- `bin/release/RELEASE_NOTES_v0.1.0.md` - Pre-filled template
- `bin/release/CHANGELOG_SECTION_v0.1.0.txt` - Extracted changelog section

Edit the release notes file to:
1. Fill in the "What's New" section with highlights
2. Paste the contents of `SHA256SUMS.txt` into the checksums section
3. Customize any other sections as needed

### 7. Create GitHub Release

#### Option A: Using GitHub CLI (Recommended)

```bash
gh release create v0.1.0 \
  bin/release/*.tar.gz \
  bin/release/SHA256SUMS.txt \
  --title "Orthanc CLI v0.1.0" \
  --notes-file bin/release/RELEASE_NOTES_v0.1.0.md
```

#### Option B: Using GitHub Web UI

1. Go to https://github.com/proencaj/orthanc-cli/releases/new
2. Select the tag: `v0.1.0`
3. Set the title: `Orthanc CLI v0.1.0`
4. Copy the contents of `bin/release/RELEASE_NOTES_v0.1.0.md` into the description
5. Attach all files from `bin/release/`:
   - All `.tar.gz` archives
   - `SHA256SUMS.txt`
6. Click "Publish release"

### 8. Verify the Release

After publishing:

1. Check the release page: https://github.com/proencaj/orthanc-cli/releases/latest
2. Verify download links work
3. Test installation from a release archive:

```bash
# Download and test
curl -L https://github.com/proencaj/orthanc-cli/releases/download/v0.1.0/orthanc-v0.1.0-darwin-arm64.tar.gz | tar xz
./orthanc version
# Should show: Version: v0.1.0
```

### 9. Post-Release

1. Update README.md if installation instructions changed
2. Announce the release:
   - Social media
   - Mailing lists
   - Community forums
3. Start work on next version:
   - Create a new `[Unreleased]` section in CHANGELOG.md

## Troubleshooting

### "Working directory is not clean"

```bash
git status
git add .
git commit -m "fix: pre-release cleanup"
```

### "No version tag found"

```bash
# Make sure you created and pushed the tag
git tag -a v0.1.0 -m "Release v0.1.0"
git push origin v0.1.0
```

### "Tests failed"

```bash
# Run tests to see what's failing
go test ./...

# Fix issues and try again
```

### Release artifacts from previous build exist

```bash
# Clean everything
make clean

# Start fresh
make release
```

## Quick Reference

Complete release in one workflow:

```bash
# 1. Prepare
vim CHANGELOG.md  # Update changelog
git add CHANGELOG.md
git commit -m "docs: update CHANGELOG for v0.1.0"
git push origin main

# 2. Tag
git tag -a v0.1.0 -m "Release v0.1.0"
git push origin v0.1.0

# 3. Build
make release
make checksums

# 4. Prepare notes
./.github/prepare-release.sh v0.1.0
vim bin/release/RELEASE_NOTES_v0.1.0.md  # Customize

# 5. Release
gh release create v0.1.0 \
  bin/release/*.tar.gz \
  bin/release/SHA256SUMS.txt \
  --title "Orthanc CLI v0.1.0" \
  --notes-file bin/release/RELEASE_NOTES_v0.1.0.md
```

## Automation (Future)

Consider setting up GitHub Actions to automate:
- Building release artifacts when tags are pushed
- Running tests before release
- Creating draft releases automatically
- Publishing to package managers (Homebrew, etc.)
