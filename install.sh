#!/bin/bash

# Branch Switcher Installation Script
# Usage: curl -sSL https://raw.githubusercontent.com/hoangbits/branch-switcher/main/install.sh | bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Detect OS and architecture
OS="$(uname -s)"
ARCH="$(uname -m)"

# Convert to lowercase
OS="$(echo "$OS" | tr '[:upper:]' '[:lower:]')"

# Map architecture names
case "$ARCH" in
    x86_64|amd64)
        ARCH="amd64"
        ;;
    arm64|aarch64)
        ARCH="arm64"
        ;;
    *)
        echo -e "${RED}Unsupported architecture: $ARCH${NC}"
        exit 1
        ;;
esac

# Map OS names
case "$OS" in
    linux)
        OS="linux"
        ;;
    darwin)
        OS="darwin"
        ;;
    *)
        echo -e "${RED}Unsupported OS: $OS${NC}"
        exit 1
        ;;
esac

BINARY_NAME="branch-switcher"
REPO="hoangbits/branch-switcher"
INSTALL_DIR="/usr/local/bin"

echo -e "${BLUE}🌿 Installing Branch Switcher...${NC}"
echo -e "${YELLOW}OS: $OS, Architecture: $ARCH${NC}"

# Check if running as root for system install
if [[ $EUID -eq 0 ]]; then
    INSTALL_DIR="/usr/local/bin"
else
    # Try to install to user's local bin if exists
    if [[ -d "$HOME/.local/bin" ]]; then
        INSTALL_DIR="$HOME/.local/bin"
    else
        mkdir -p "$HOME/.local/bin"
        INSTALL_DIR="$HOME/.local/bin"
        echo -e "${YELLOW}Created $HOME/.local/bin directory${NC}"
        echo -e "${YELLOW}Make sure $HOME/.local/bin is in your PATH${NC}"
    fi
fi

# Get the latest release
echo -e "${BLUE}Fetching latest release...${NC}"
LATEST_RELEASE=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep -o '"tag_name": *"[^"]*"' | grep -o '[^"]*$')

if [[ -z "$LATEST_RELEASE" ]]; then
    echo -e "${YELLOW}No releases found, building from source...${NC}"

    # Check if Go is installed
    if ! command -v go &> /dev/null; then
        echo -e "${RED}Go is not installed. Please install Go or download a binary release.${NC}"
        exit 1
    fi

    # Build from source
    TEMP_DIR=$(mktemp -d)
    cd "$TEMP_DIR"

    echo -e "${BLUE}Cloning repository...${NC}"
    git clone "https://github.com/$REPO.git"
    cd branch-switcher

    echo -e "${BLUE}Building binary...${NC}"
    go build -ldflags "-s -w" -o "$BINARY_NAME"

    # Install binary
    if [[ $EUID -eq 0 ]] || [[ "$INSTALL_DIR" == "$HOME/.local/bin" ]]; then
        mv "$BINARY_NAME" "$INSTALL_DIR/"
        # Create brs shortcut
        ln -sf "$INSTALL_DIR/$BINARY_NAME" "$INSTALL_DIR/brs"
    else
        sudo mv "$BINARY_NAME" "$INSTALL_DIR/"
        # Create brs shortcut
        sudo ln -sf "$INSTALL_DIR/$BINARY_NAME" "$INSTALL_DIR/brs"
    fi

    # Cleanup
    cd /
    rm -rf "$TEMP_DIR"
else
    # Download binary from release
    DOWNLOAD_URL="https://github.com/$REPO/releases/download/$LATEST_RELEASE/${BINARY_NAME}-${OS}-${ARCH}"

    echo -e "${BLUE}Downloading $BINARY_NAME $LATEST_RELEASE...${NC}"

    # Download to temp file
    TEMP_FILE=$(mktemp)
    if curl -sL "$DOWNLOAD_URL" -o "$TEMP_FILE"; then
        chmod +x "$TEMP_FILE"

        # Install binary
        if [[ $EUID -eq 0 ]] || [[ "$INSTALL_DIR" == "$HOME/.local/bin" ]]; then
            mv "$TEMP_FILE" "$INSTALL_DIR/$BINARY_NAME"
            # Create brs shortcut
            ln -sf "$INSTALL_DIR/$BINARY_NAME" "$INSTALL_DIR/brs"
        else
            sudo mv "$TEMP_FILE" "$INSTALL_DIR/$BINARY_NAME"
            # Create brs shortcut
            sudo ln -sf "$INSTALL_DIR/$BINARY_NAME" "$INSTALL_DIR/brs"
        fi
    else
        echo -e "${YELLOW}Binary not available for $OS-$ARCH, building from source...${NC}"

        # Fallback to building from source
        if ! command -v go &> /dev/null; then
            echo -e "${RED}Go is not installed. Please install Go or download a binary release manually.${NC}"
            exit 1
        fi

        TEMP_DIR=$(mktemp -d)
        cd "$TEMP_DIR"

        git clone "https://github.com/$REPO.git"
        cd branch-switcher
        go build -ldflags "-s -w" -o "$BINARY_NAME"

        if [[ $EUID -eq 0 ]] || [[ "$INSTALL_DIR" == "$HOME/.local/bin" ]]; then
            mv "$BINARY_NAME" "$INSTALL_DIR/"
            # Create brs shortcut
            ln -sf "$INSTALL_DIR/$BINARY_NAME" "$INSTALL_DIR/brs"
        else
            sudo mv "$BINARY_NAME" "$INSTALL_DIR/"
            # Create brs shortcut
            sudo ln -sf "$INSTALL_DIR/$BINARY_NAME" "$INSTALL_DIR/brs"
        fi

        cd /
        rm -rf "$TEMP_DIR"
    fi
fi

# Verify installation
if command -v "$BINARY_NAME" &> /dev/null; then
    echo -e "${GREEN}✅ Branch Switcher installed successfully!${NC}"
    echo -e "${GREEN}Location: $INSTALL_DIR/$BINARY_NAME${NC}"
    if command -v "brs" &> /dev/null; then
        echo -e "${GREEN}Shortcut: $INSTALL_DIR/brs${NC}"
    fi
    echo ""
    echo -e "${BLUE}Usage:${NC}"
    echo -e "  ${YELLOW}$BINARY_NAME${NC}  # Full name"
    echo -e "  ${YELLOW}brs${NC}           # Short name"
    echo ""
    echo -e "${BLUE}Try it now:${NC}"
    echo -e "  ${YELLOW}cd /path/to/your/repositories && brs${NC}"
else
    echo -e "${RED}❌ Installation failed${NC}"
    echo -e "${YELLOW}Please check that $INSTALL_DIR is in your PATH${NC}"
    exit 1
fi