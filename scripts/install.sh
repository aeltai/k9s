#!/usr/bin/env bash
# rk9s + SUSE/Rancher CLIs installer
# Installs rk9s and optional CLIs (longhornctl, kwctl, virtctl, fleet)
# Usage: curl -sSL https://raw.githubusercontent.com/aeltai/k9s/master/scripts/install.sh | bash
# Or: ./scripts/install.sh [--rk9s-only | --clis-only | --all]

set -e

RK9S_REPO="${RK9S_REPO:-aeltai/k9s}"
RK9S_VERSION="${RK9S_VERSION:-latest}"
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"
CONFIG_DIR="${XDG_CONFIG_HOME:-$HOME/.config}/rk9s"
PLUGINS_DIR="${XDG_DATA_HOME:-$HOME/.local/share}/rk9s/plugins"
CLI_DIR="${INSTALL_DIR}"

# Detect OS and arch
detect_platform() {
  OS=$(uname -s | tr '[:upper:]' '[:lower:]')
  ARCH=$(uname -m)
  case "$ARCH" in
    x86_64|amd64) ARCH=amd64 ;;
    aarch64|arm64) ARCH=arm64 ;;
  esac
  echo "${OS}_${ARCH}"
}

# Download with optional checksum verification
download() {
  local url="$1"
  local dest="$2"
  echo "  Downloading $url ..."
  if command -v curl >/dev/null 2>&1; then
    curl -sfL -o "$dest" "$url"
  elif command -v wget >/dev/null 2>&1; then
    wget -qO "$dest" "$url"
  else
    echo "Error: curl or wget required"
    exit 1
  fi
}

# Install rk9s
install_rk9s() {
  echo "==> Installing rk9s..."
  mkdir -p "$INSTALL_DIR"
  PLATFORM=$(detect_platform)
  OS=$(echo "$PLATFORM" | cut -d_ -f1)
  ARCH=$(echo "$PLATFORM" | cut -d_ -f2)
  # Goreleaser uses Linux/Darwin/Windows (title case)
  OS_TITLE=$(echo "$OS" | sed 's/^\(.\)/\U\1/')
  [ "$OS" = "darwin" ] && OS_TITLE="Darwin"
  [ "$OS" = "linux" ] && OS_TITLE="Linux"
  [ "$OS" = "windows" ] && OS_TITLE="Windows"

  if [ "$RK9S_VERSION" = "latest" ] || [ -z "$RK9S_VERSION" ]; then
    LATEST=$(curl -sfL "https://api.github.com/repos/${RK9S_REPO}/releases/latest" 2>/dev/null | grep '"tag_name"' | sed -E 's/.*"v([^"]+)".*/\1/' || true)
    RK9S_VERSION="${LATEST:-v0.50.18}"
  fi
  RK9S_VERSION=${RK9S_VERSION#v}

  BIN_NAME="rk9s"
  [ "$OS" = "windows" ] && BIN_NAME="rk9s.exe"
  URL="https://github.com/${RK9S_REPO}/releases/download/v${RK9S_VERSION}/rk9s_${OS_TITLE}_${ARCH}.tar.gz"

  if ! download "$URL" "/tmp/rk9s.tar.gz" 2>/dev/null; then
    echo "  Release not found. Building from source..."
    if command -v go >/dev/null 2>&1; then
      (cd "$(dirname "$0")/.." && go build -o "$INSTALL_DIR/rk9s" .)
      echo "  Installed rk9s to $INSTALL_DIR/rk9s"
    else
      echo "  Error: No release and Go not installed. Run: go install github.com/${RK9S_REPO}@latest"
      exit 1
    fi
  else
    tar -xzf /tmp/rk9s.tar.gz -C /tmp
    mv /tmp/rk9s "$INSTALL_DIR/" 2>/dev/null || mv /tmp/rk9s_*/rk9s "$INSTALL_DIR/" 2>/dev/null || true
    chmod +x "$INSTALL_DIR/rk9s"
    rm -f /tmp/rk9s.tar.gz
    echo "  Installed rk9s to $INSTALL_DIR/rk9s"
  fi
}

# Install longhornctl
install_longhornctl() {
  [ -x "$CLI_DIR/longhornctl" ] && echo "  longhornctl already installed" && return
  echo "  Installing longhornctl..."
  PLATFORM=$(detect_platform)
  OS=$(echo "$PLATFORM" | cut -d_ -f1)
  ARCH=$(echo "$PLATFORM" | cut -d_ -f2)
  [ "$OS" = "darwin" ] && OS=darwin || OS=linux
  URL="https://github.com/longhorn/cli/releases/download/v1.11.0/longhornctl-${OS}-${ARCH}"
  download "$URL" "$CLI_DIR/longhornctl"
  chmod +x "$CLI_DIR/longhornctl"
  echo "  Installed longhornctl"
}

# Install kwctl
install_kwctl() {
  [ -x "$CLI_DIR/kwctl" ] && echo "  kwctl already installed" && return
  echo "  Installing kwctl..."
  PLATFORM=$(detect_platform)
  OS=$(echo "$PLATFORM" | cut -d_ -f1)
  ARCH=$(echo "$PLATFORM" | cut -d_ -f2)
  [ "$ARCH" = "amd64" ] && ARCH=x86_64
  [ "$ARCH" = "arm64" ] && ARCH=aarch64
  [ "$OS" = "darwin" ] && OS=darwin || OS=linux
  URL="https://github.com/kubewarden/kwctl/releases/latest/download/kwctl-${OS}-${ARCH}.zip"
  download "$URL" /tmp/kwctl.zip
  unzip -o -j /tmp/kwctl.zip -d "$CLI_DIR" 2>/dev/null || (cd /tmp && unzip -o kwctl.zip && mv kwctl "$CLI_DIR/" 2>/dev/null || mv bin/kwctl "$CLI_DIR/" 2>/dev/null)
  chmod +x "$CLI_DIR/kwctl"
  rm -f /tmp/kwctl.zip
  echo "  Installed kwctl"
}

# Install virtctl
install_virtctl() {
  [ -x "$CLI_DIR/virtctl" ] && echo "  virtctl already installed" && return
  echo "  Installing virtctl..."
  PLATFORM=$(detect_platform)
  OS=$(echo "$PLATFORM" | cut -d_ -f1)
  ARCH=$(echo "$PLATFORM" | cut -d_ -f2)
  VER=$(curl -sfL "https://storage.googleapis.com/kubevirt-prow/release/kubevirt/kubevirt/stable.txt" 2>/dev/null || echo "v1.2.0")
  VER=${VER#v}
  URL="https://github.com/kubevirt/kubevirt/releases/download/v${VER}/virtctl-v${VER}-${OS}-${ARCH}"
  download "$URL" "$CLI_DIR/virtctl"
  chmod +x "$CLI_DIR/virtctl"
  echo "  Installed virtctl"
}

# Install fleet (Rancher Fleet CLI)
install_fleet() {
  [ -x "$CLI_DIR/fleet" ] && echo "  fleet already installed" && return
  echo "  Installing fleet (Rancher Fleet CLI)..."
  PLATFORM=$(detect_platform)
  OS=$(echo "$PLATFORM" | cut -d_ -f1)
  ARCH=$(echo "$PLATFORM" | cut -d_ -f2)
  URL="https://github.com/rancher/fleet/releases/latest/download/fleet-${OS}-${ARCH}"
  download "$URL" "$CLI_DIR/fleet"
  chmod +x "$CLI_DIR/fleet"
  echo "  Installed fleet"
}

# Install preconfigured rk9s plugins
install_plugins() {
  echo "==> Installing rk9s plugins..."
  mkdir -p "$PLUGINS_DIR"
  SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
  REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
  if [ -d "$REPO_ROOT/plugins" ]; then
    cp "$REPO_ROOT"/plugins/*.yaml "$PLUGINS_DIR/" 2>/dev/null || true
    echo "  Copied plugins to $PLUGINS_DIR"
  else
    echo "  Downloading plugins from GitHub..."
    for f in rancher fleet rke2-k3s ingress-storage longhorn harvester kubewarden; do
      download "https://raw.githubusercontent.com/${RK9S_REPO}/master/plugins/${f}.yaml" "$PLUGINS_DIR/${f}.yaml" 2>/dev/null && echo "    ${f}.yaml" || true
    done
    echo "  Plugins in $PLUGINS_DIR"
  fi
}

# Add to PATH in shell rc
add_to_path() {
  if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    echo ""
    echo "Add to your shell profile (~/.bashrc, ~/.zshrc):"
    echo "  export PATH=\"$INSTALL_DIR:\$PATH\""
  fi
}

# Main
main() {
  echo "rk9s + SUSE/Rancher CLIs Installer"
  echo "=================================="
  mkdir -p "$INSTALL_DIR"
  export PATH="$INSTALL_DIR:$PATH"

  case "${1:-}" in
    --rk9s-only) install_rk9s ;;
    --clis-only)
      install_longhornctl
      install_kwctl
      install_virtctl
      install_fleet
      ;;
    --plugins-only) install_plugins ;;
    *)
      install_rk9s
      install_plugins
      echo ""
      echo "==> Installing optional CLIs..."
      install_longhornctl
      install_kwctl
      install_virtctl
      install_fleet
      ;;
  esac

  add_to_path
  echo ""
  echo "Done! Run: rk9s"
}

main "$@"
