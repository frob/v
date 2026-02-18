#!/usr/bin/env bash
set -euo pipefail

REPO="frob/v"
BIN_NAME="v"
INSTALL_DIR="/usr/local/bin"

# Detect OS
OS="$(uname -s)"
case "${OS}" in
  Linux)  OS="linux" ;;
  Darwin) OS="darwin" ;;
  *)
    echo "Unsupported OS: ${OS}" >&2
    exit 1
    ;;
esac

# Detect architecture
ARCH="$(uname -m)"
case "${ARCH}" in
  x86_64)          ARCH="amd64" ;;
  aarch64 | arm64) ARCH="arm64" ;;
  *)
    echo "Unsupported architecture: ${ARCH}" >&2
    exit 1
    ;;
esac

BASE_URL="https://github.com/${REPO}/releases/latest/download"

# Detect package manager and install via native package if on Linux
if [ "${OS}" = "linux" ]; then
  if command -v pacman &>/dev/null; then
    PKG="${BIN_NAME}_${OS}_${ARCH}.pkg.tar.zst"
    echo "Downloading ${PKG}..."
    curl -fsSL "${BASE_URL}/${PKG}" -o "/tmp/${PKG}"
    sudo pacman -U "/tmp/${PKG}"
    rm -f "/tmp/${PKG}"
    exit 0
  elif command -v dpkg &>/dev/null; then
    PKG="${BIN_NAME}_${OS}_${ARCH}.deb"
    echo "Downloading ${PKG}..."
    curl -fsSL "${BASE_URL}/${PKG}" -o "/tmp/${PKG}"
    sudo dpkg -i "/tmp/${PKG}"
    rm -f "/tmp/${PKG}"
    exit 0
  elif command -v rpm &>/dev/null; then
    PKG="${BIN_NAME}_${OS}_${ARCH}.rpm"
    echo "Downloading ${PKG}..."
    curl -fsSL "${BASE_URL}/${PKG}" -o "/tmp/${PKG}"
    sudo rpm -i "/tmp/${PKG}"
    rm -f "/tmp/${PKG}"
    exit 0
  fi
fi

# Fallback: extract binary from tar.gz
ARCHIVE="${BIN_NAME}_${OS}_${ARCH}.tar.gz"
echo "Downloading ${ARCHIVE}..."
curl -fsSL "${BASE_URL}/${ARCHIVE}" -o "/tmp/${ARCHIVE}"
tar -xzf "/tmp/${ARCHIVE}" -C /tmp "${BIN_NAME}"
chmod +x "/tmp/${BIN_NAME}"
sudo mv "/tmp/${BIN_NAME}" "${INSTALL_DIR}/${BIN_NAME}"
rm -f "/tmp/${ARCHIVE}"

echo "${BIN_NAME} installed to ${INSTALL_DIR}/${BIN_NAME}"
