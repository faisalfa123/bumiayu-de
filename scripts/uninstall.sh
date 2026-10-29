#!/usr/bin/env bash
# ==============================================================================
# Bumiayu TUI - Uninstaller Script
# ==============================================================================

set -e

INSTALL_DIR="/usr/local/bin"
BIN_NAME="bumiayu"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

if [ "$EUID" -ne 0 ]; then
  echo -e "${YELLOW}⚠️ Skrip membutuhkan hak akses root.${NC}"
  echo -e "Silakan jalankan dengan: ${GREEN}sudo bash uninstall.sh${NC}"
  exit 1
fi

if [ -f "${INSTALL_DIR}/${BIN_NAME}" ]; then
  rm -f "${INSTALL_DIR}/${BIN_NAME}"
  echo -e "${GREEN}✓ Binary ${INSTALL_DIR}/${BIN_NAME} berhasil dihapus.${NC}"
else
  echo -e "${YELLOW}Binary ${INSTALL_DIR}/${BIN_NAME} tidak ditemukan.${NC}"
fi

echo -e "${GREEN}Pembersihan selesai.${NC}"
