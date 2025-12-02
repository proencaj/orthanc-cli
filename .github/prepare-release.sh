#!/bin/bash
# prepare-release.sh - Helper script to prepare GitHub release notes

set -e

# Colors
BLUE='\033[0;34m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Get version from git tag or argument
VERSION="${1:-$(git describe --tags --abbrev=0 2>/dev/null)}"

if [ -z "$VERSION" ]; then
    echo -e "${RED}Error: No version specified and no git tag found${NC}"
    echo "Usage: $0 [version]"
    echo "Example: $0 v0.1.0"
    exit 1
fi

# Remove 'v' prefix if present for some operations
VERSION_NO_V="${VERSION#v}"

echo -e "${BLUE}Preparing release notes for ${VERSION}${NC}"
echo ""

# Find previous version
PREVIOUS_VERSION=$(git tag --sort=-version:refname | grep -A1 "^${VERSION}$" | tail -1 || echo "")

if [ -z "$PREVIOUS_VERSION" ]; then
    PREVIOUS_VERSION=$(git rev-list --max-parents=0 HEAD)
    echo -e "${YELLOW}No previous tag found, using first commit${NC}"
else
    echo -e "${GREEN}Previous version: ${PREVIOUS_VERSION}${NC}"
fi

# Create release notes from template
RELEASE_NOTES_FILE="bin/release/RELEASE_NOTES_${VERSION}.md"
mkdir -p bin/release

echo -e "${BLUE}Generating release notes...${NC}"

# Copy template and replace placeholders
cp .github/release-template.md "$RELEASE_NOTES_FILE"

# Replace version placeholders
if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS sed
    sed -i '' "s/{VERSION}/${VERSION}/g" "$RELEASE_NOTES_FILE"
    sed -i '' "s/{PREVIOUS_VERSION}/${PREVIOUS_VERSION}/g" "$RELEASE_NOTES_FILE"
else
    # Linux sed
    sed -i "s/{VERSION}/${VERSION}/g" "$RELEASE_NOTES_FILE"
    sed -i "s/{PREVIOUS_VERSION}/${PREVIOUS_VERSION}/g" "$RELEASE_NOTES_FILE"
fi

echo -e "${GREEN}✓ Release notes template created: ${RELEASE_NOTES_FILE}${NC}"
echo ""

# Extract relevant changelog section if CHANGELOG.md exists
if [ -f "CHANGELOG.md" ]; then
    echo -e "${BLUE}Extracting changelog section for ${VERSION}...${NC}"

    # Try to extract the section for this version
    CHANGELOG_SECTION=$(awk "/## \[${VERSION_NO_V}\]|## \[${VERSION}\]/,/## \[.*\]/{if (/## \[.*\]/ && !/${VERSION_NO_V}/ && !/${VERSION}/) exit; print}" CHANGELOG.md)

    if [ -n "$CHANGELOG_SECTION" ]; then
        echo "$CHANGELOG_SECTION" > "bin/release/CHANGELOG_SECTION_${VERSION}.txt"
        echo -e "${GREEN}✓ Changelog section extracted${NC}"
        echo ""
        echo -e "${YELLOW}Preview:${NC}"
        echo "----------------------------------------"
        head -20 "bin/release/CHANGELOG_SECTION_${VERSION}.txt"
        echo "----------------------------------------"
    else
        echo -e "${YELLOW}! Could not extract changelog section for ${VERSION}${NC}"
        echo -e "${YELLOW}  Make sure CHANGELOG.md has a section for this version${NC}"
    fi
fi

echo ""
echo -e "${BLUE}Next steps:${NC}"
echo -e "  1. Edit ${YELLOW}${RELEASE_NOTES_FILE}${NC}"
echo -e "     - Add 'What's New' highlights"
echo -e "     - Review and customize the content"
echo ""
echo -e "  2. Build release artifacts:"
echo -e "     ${GREEN}make release${NC}"
echo ""
echo -e "  3. Generate checksums:"
echo -e "     ${GREEN}make checksums${NC}"
echo ""
echo -e "  4. Copy SHA256SUMS.txt into the release notes"
echo ""
echo -e "  5. Create GitHub release:"
echo -e "     ${GREEN}gh release create ${VERSION} bin/release/*.tar.gz \\${NC}"
echo -e "     ${GREEN}  --title 'Orthanc CLI ${VERSION}' \\${NC}"
echo -e "     ${GREEN}  --notes-file ${RELEASE_NOTES_FILE}${NC}"
echo ""
echo -e "     Or manually at:"
echo -e "     ${BLUE}https://github.com/proencaj/orthanc-cli/releases/new${NC}"
