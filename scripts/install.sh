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
    URL=$(curl -s https://api.github.com/repos/aboubakary833/codexa/releases/latest \
        | grep "browser_download_url.*linux_amd64.tar.gz" \
        | cut -d '"' -f 4)
else
    URL="https://github.com/aboubakary833/codexa/releases/download/$VERSION/codexa_${VERSION}_linux_amd64.tar.gz"
fi

curl -L "$URL" -o "$TMP_DIR/codexa.tar.gz"

echo "Extracting..."
tar -xzf "$TMP_DIR/codexa.tar.gz" -C "$TMP_DIR"

echo "Installing binary..."
mv "$TMP_DIR/codexa" "$INSTALL_DIR/$BIN_NAME"
chmod +x "$INSTALL_DIR/$BIN_NAME"

# Add codexa to PATH
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

# Generating completions
echo "Generating shell completions..."
case "$SHELL" in
    */bash)
        "$INSTALL_DIR/$BIN_NAME" completion bash > "$HOME/.local/share/bash-completion/completions/codexa"
        ;;
    */zsh)
        mkdir -p "$HOME/.zsh/completions"
        "$INSTALL_DIR/$BIN_NAME" completion zsh > "$HOME/.zsh/completions/_codexa"
        ;;
    */fish)
        mkdir -p "$HOME/.config/fish/completions"
        "$INSTALL_DIR/$BIN_NAME" completion fish > "$HOME/.config/fish/completions/codexa.fish"
        ;;
esac

echo "Codexa installed successfully! Restart your shell."
