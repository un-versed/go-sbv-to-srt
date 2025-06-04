#!/bin/bash

# Shell completion setup script for go-sbv-to-srt
# This script helps users set up shell autocompletion

set -e

BINARY_NAME="go-sbv-to-srt"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "🔧 Setting up shell completion for ${BINARY_NAME}..."

# Detect shell
if [ -n "$ZSH_VERSION" ]; then
    SHELL_TYPE="zsh"
    COMPLETION_DIR="$HOME/.local/share/zsh/site-functions"
elif [ -n "$BASH_VERSION" ]; then
    SHELL_TYPE="bash"
    COMPLETION_DIR="$HOME/.local/share/bash-completion/completions"
else
    echo "❌ Unsupported shell. Please use bash or zsh."
    exit 1
fi

echo "📋 Detected shell: $SHELL_TYPE"

# Create completion directory if it doesn't exist
mkdir -p "$COMPLETION_DIR"

# Generate completion script
echo "📝 Generating completion script..."
if [ -f "${SCRIPT_DIR}/${BINARY_NAME}" ]; then
    "${SCRIPT_DIR}/${BINARY_NAME}" completion "$SHELL_TYPE" > "${COMPLETION_DIR}/${BINARY_NAME}"
    echo "✅ Completion script installed to: ${COMPLETION_DIR}/${BINARY_NAME}"
else
    echo "❌ Binary not found at ${SCRIPT_DIR}/${BINARY_NAME}"
    echo "Please build the project first: make build"
    exit 1
fi

# Provide instructions
echo
echo "📖 To enable completion, add the following to your shell configuration:"
case $SHELL_TYPE in
    "zsh")
        echo "   echo 'fpath=(~/.local/share/zsh/site-functions \$fpath)' >> ~/.zshrc"
        echo "   echo 'autoload -U compinit && compinit' >> ~/.zshrc"
        ;;
    "bash")
        echo "   echo 'source ~/.local/share/bash-completion/completions/${BINARY_NAME}' >> ~/.bashrc"
        ;;
esac
echo
echo "Then restart your shell or run: source ~/.${SHELL_TYPE}rc"
echo
echo "🎉 Shell completion setup complete!"
