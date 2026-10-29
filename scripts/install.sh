#!/usr/bin/env bash
# ==============================================================================
# Bumiayu TUI - Installer Script
# Lightweight & Elegant TUI Desktop Environment for ODROID-C2 & Linux Servers
# ==============================================================================

set -e

REPO="bumiayu/bumiayu-de"
INSTALL_DIR="/usr/local/bin"
BIN_NAME="bumiayu"

RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

echo -e "${CYAN}"
echo "  ____                  _                    _____ _   _ ___ "
echo " | __ ) _   _ _ __ ___ (_) __ _ _   _ _   _ |_   _| | | |_ _|"
echo " |  _ \| | | | '_ \` _ \| |/ _\` | | | | | | |  | | | | | || | "
echo " | |_) | |_| | | | | | | | (_| | |_| | |_| |  | | | |_| || | "
echo " |____/ \__,_|_| |_| |_|_|\__,_|\__,_|\__,_|  |_|  \___/|___|"
echo -e "${NC}"
echo -e "${BLUE}Desktop Environment Berbasis TUI Ultra-Ringan & Ramah SSH${NC}"
echo -e "Dirancang khusus untuk ODROID-C2 (ARM64) dan Server Linux (AMD64)"
echo "------------------------------------------------------------------"

# 1. Check Root Privileges
if [ "$EUID" -ne 0 ]; then
  echo -e "${YELLOW}⚠️ Skrip membutuhkan hak akses root untuk memasang ke ${INSTALL_DIR}.${NC}"
  echo -e "Silakan jalankan kembali dengan: ${GREEN}sudo bash install.sh${NC}"
  exit 1
fi

# 2. Detect System Architecture
ARCH=$(uname -m)
case "$ARCH" in
  aarch64|arm64)
    TARGET_ARCH="arm64"
    echo -e "🔍 Terdeteksi arsitektur: ${GREEN}ARM64 (aarch64) - Cocok untuk ODROID-C2 / SBC${NC}"
    ;;
  x86_64|amd64)
    TARGET_ARCH="amd64"
    echo -e "🔍 Terdeteksi arsitektur: ${GREEN}AMD64 (x86_64) - Cocok untuk Server / VPS${NC}"
    ;;
  armv7l|armhf)
    TARGET_ARCH="armv7"
    echo -e "🔍 Terdeteksi arsitektur: ${GREEN}ARMv7 (32-bit)${NC}"
    ;;
  *)
    echo -e "${RED}❌ Arsitektur $ARCH belum didukung secara otomatis.${NC}"
    exit 1
    ;;
esac

# 3. Fetch latest release version from GitHub API
echo -e "🌐 Memeriksa versi rilis terbaru dari GitHub..."
LATEST_TAG=$(curl -sSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/' || true)

if [ -z "$LATEST_TAG" ]; then
  LATEST_TAG="v1.0.0"
  echo -e "${YELLOW}Menggunakan versi default: ${LATEST_TAG}${NC}"
else
  echo -e "🚀 Versi terbaru ditemukan: ${GREEN}${LATEST_TAG}${NC}"
fi

DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${LATEST_TAG}/bumiayu-linux-${TARGET_ARCH}.tar.gz"
TEMP_DIR=$(mktemp -d)

echo -e "📥 Mengunduh binary: ${CYAN}${DOWNLOAD_URL}${NC} ..."
if curl -fsSL "$DOWNLOAD_URL" -o "${TEMP_DIR}/bumiayu.tar.gz"; then
  echo -e "📦 Mengekstrak berkas..."
  tar -xzf "${TEMP_DIR}/bumiayu.tar.gz" -C "${TEMP_DIR}"
  install -m 755 "${TEMP_DIR}/bumiayu-linux-${TARGET_ARCH}" "${INSTALL_DIR}/${BIN_NAME}"
else
  echo -e "${YELLOW}Mengunduh binary langsung (fallback)...${NC}"
  curl -fsSL "https://github.com/${REPO}/releases/download/${LATEST_TAG}/bumiayu-linux-${TARGET_ARCH}" -o "${INSTALL_DIR}/${BIN_NAME}"
  chmod +x "${INSTALL_DIR}/${BIN_NAME}"
fi

# Cleanup
rm -rf "$TEMP_DIR"

echo "------------------------------------------------------------------"
echo -e "${GREEN}🎉 Selamat! Bumiayu TUI berhasil terpasang di ${INSTALL_DIR}/${BIN_NAME}${NC}"
echo ""
echo -e "Cara menjalankan:"
echo -e "  Ketik: ${CYAN}bumiayu${NC}"
echo ""
echo -e "Tips SSH Auto-start:"
echo -e "  Agar Bumiayu TUI otomatis terbuka setiap kali Anda login via SSH,"
echo -e "  tambahkan baris berikut ke ${YELLOW}~/.bashrc${NC} Anda:"
echo -e "  ${GREEN}[ -z \"\$BUMIAYU_LOADED\" ] && exec bumiayu${NC}"
echo "------------------------------------------------------------------"
