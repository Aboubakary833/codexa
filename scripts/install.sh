#!/usr/bin/env bash
set -e

# Variables
VERSION=${1:-"latest"}
BIN_NAME="codexa"
INSTALL_DIR="${HOME}/.local/bin"
TMP_DIR=$(mktemp -d)

mkdir -p "$INSTALL_DIR"

echo "Downloading Codexa $VERSION..."

if [ "$VERSION" = "latest" ]; then
    OS=$(uname | tr '[:upper:]' '[:lower:]')
    URL=$(curl -s https://api.github.com/repos/aboubakary833/codexa/releases/latest \
        | grep "browser_download_url.*${OS}_amd64.tar.gz" \
        | cut -d '"' -f 4)
else
    OS=$(uname | tr '[:upper:]' '[:lower:]')
    URL="https://github.com/aboubakary833/codexa/releases/download/$VERSION/codexa_${OS}_amd64.tar.gz"
fi

curl -L "$URL" -o "$TMP_DIR/codexa.tar.gz"

echo "Extracting..."
tar -xzf "$TMP_DIR/codexa.tar.gz" -C "$TMP_DIR"

# Find extracted binary
BINARY_PATH=$(find "$TMP_DIR" -type f -name "codexa*" | head -n1)
if [ -z "$BINARY_PATH" ]; then
    echo "Error: no binary found after extraction"
    exit 1
fi

echo "Installing binary..."
mv "$BINARY_PATH" "$INSTALL_DIR/$BIN_NAME"
chmod +x "$INSTALL_DIR/$BIN_NAME"

# Add codexa to PATH if not already present
if ! echo "$PATH" | grep -q "$INSTALL_DIR"; then
    SHELL_RC=""
    case "$SHELL" in
        */bash) SHELL_RC="$HOME/.bashrc" ;;
        */zsh) SHELL_RC="$HOME/.zshrc" ;;
        */fish) SHELL_RC="$HOME/.config/fish/config.fish" ;;
        *) SHELL_RC="$HOME/.profile" ;;
    esac

    echo "Adding Codexa to PATH in $SHELL_RC"
    echo "export PATH=\"$INSTALL_DIR:\$PATH\"" >> "$SHELL_RC"
fi

# Generating shell completions
echo "Generating shell completions..."
case "$SHELL" in
    */bash)
        mkdir -p "$HOME/.local/share/bash-completion/completions"
        "$INSTALL_DIR/$BIN_NAME" completion bash > "$HOME/.local/share/bash-completion/completions/codexa"
        ;;
    */zsh)
        mkdir -p "$HOME/.zsh/completions"
        "$INSTALL_DIR/$BIN_NAME" completion zsh > "$HOME/.zsh/completions/_codexa"
        echo 'fpath=(~/.zsh/completions $fpath)' >> ~/.zshrc
        echo 'autoload -U compinit && compinit' >> ~/.zshrc
        ;;
    */fish)
        mkdir -p "$HOME/.config/fish/completions"
        "$INSTALL_DIR/$BIN_NAME" completion fish > "$HOME/.config/fish/completions/codexa.fish"
        ;;
esac

# Installing default categories of snippets
echo "Installing javascript and go snippets"
"$INSTALL_DIR/$BIN_NAME" sync js
"$INSTALL_DIR/$BIN_NAME" sync go

echo "Codexa installed successfully! Restart your shell."
