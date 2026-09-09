#!/usr/bin/env bash
# USC (Universal Skill Compiler) One-Click Automated Installer for Linux/macOS
# Usage: bash install.sh
# Or:    curl -fsSL https://raw.githubusercontent.com/Freecode100Year/usc/main/install.sh | bash

set -e

echo "================================================================"
echo "       USC (Universal Skill Compiler) One-Click Installer      "
echo "================================================================"

if ! command -v go >/dev/null 2>&1; then
    echo "Error: Go (golang) is required to build USC. Please install Go: https://golang.org/dl/"
    exit 1
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" 2>/dev/null && pwd || pwd)"
cd "$SCRIPT_DIR"

echo "[+] Building USC clean-room binary..."
mkdir -p bin
go build -o bin/usc ./cmd/usc

# Determine install location
TARGET_DIR="/usr/local/bin"
if [ ! -w "$TARGET_DIR" ]; then
    TARGET_DIR="$HOME/.local/bin"
    mkdir -p "$TARGET_DIR"
fi

cp bin/usc "$TARGET_DIR/usc"
chmod +x "$TARGET_DIR/usc"
echo "[✓] Installed binary to: $TARGET_DIR/usc"

# Register built-in skill into Antigravity or OpenClaw if present
AGY_SKILL_DIR="$HOME/.gemini/config/skills/usc"
if [ -d "$HOME/.gemini" ]; then
    mkdir -p "$AGY_SKILL_DIR"
    cat << 'EOF' > "$AGY_SKILL_DIR/SKILL.md"
---
name: usc
description: Universal Skill Compiler (USC) - Compile, verify, and install skills from ClawHub, GitHub, or local source into local AI Agents.
---

# Universal Skill Compiler (USC)

Use USC to compile and install any skill from ClawHub or GitHub into your local Agent:
- 一键编译并安装 ClawHub 技能: `usc install https://clawhub.ai/<author>/skills/<skill-name>`
- 一键编译并安装 GitHub 技能:  `usc install https://github.com/<owner>/<repo>`
- 仅执行零信任安全编译与审计:   `usc build <source-or-url>`
- 检查已支持的本地 Agent:        `usc targets`
EOF
    echo "[✓] Registered USC capability skill into Antigravity CLI ($AGY_SKILL_DIR)"
fi

echo ""
echo "================================================================"
echo " [✓] USC Installation Complete!"
echo " You or your Agent can now run:"
echo "   usc install https://clawhub.ai/thesentitrader/skills/us-stocks-analysis"
echo "================================================================"

"$TARGET_DIR/usc" targets
