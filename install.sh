#!/bin/bash

# Color codes
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}==============================================${NC}"
echo -e "${BLUE} Installing mdview CLI globally...             ${NC}"
echo -e "${BLUE}==============================================${NC}"

# Find wails
WAILS_BIN="$(go env GOPATH)/bin/wails"
if [ ! -f "$WAILS_BIN" ]; then
    if command -v wails &> /dev/null; then
        WAILS_BIN="wails"
    else
        echo -e "${RED}Error: Wails CLI is not installed.${NC}"
        echo "Please install it by running:"
        echo "go install github.com/wailsapp/wails/v2/cmd/wails@latest"
        exit 1
    fi
fi

# 1. Build the application for local OS
echo -e "\n${YELLOW}[1/3] Building the application...${NC}"
"$WAILS_BIN" build
if [ $? -ne 0 ]; then
    echo -e "${RED}Error: Build failed.${NC}"
    exit 1
fi

# 2. Locate the compiled binary based on OS
OS_TYPE="$(uname -s)"
SRC_BIN=""
if [ "$OS_TYPE" = "Darwin" ]; then
    SRC_BIN="$(pwd)/build/bin/mdviewer.app/Contents/MacOS/mdviewer"
else
    SRC_BIN="$(pwd)/build/bin/mdviewer"
fi

if [ ! -f "$SRC_BIN" ]; then
    echo -e "${RED}Error: Compiled binary not found at $SRC_BIN${NC}"
    exit 1
fi

# 3. Create global symlink
TARGET_DIR="/usr/local/bin"
TARGET_LINK="$TARGET_DIR/mdview"

echo -e "\n${YELLOW}[2/3] Creating global symlink...${NC}"
echo "Target path: $TARGET_LINK"

# Write symlink. Require sudo if directory is not writable
if [ -w "$TARGET_DIR" ]; then
    ln -sf "$SRC_BIN" "$TARGET_LINK"
else
    echo "Asking for administrator (sudo) privileges to write to $TARGET_DIR..."
    sudo ln -sf "$SRC_BIN" "$TARGET_LINK"
fi

if [ $? -eq 0 ]; then
    echo -e "\n${GREEN}[3/3] Installation successful!${NC}"
    echo -e "You can now run: ${BLUE}mdview <file.md>${NC} from any terminal folder!"
else
    echo -e "\n${RED}Error: Failed to create symlink.${NC}"
    echo "Alternative: Add this to your shell profile (.zshrc or .bashrc):"
    echo "alias mdview=\"$SRC_BIN\""
fi
echo -e "${BLUE}==============================================${NC}"
