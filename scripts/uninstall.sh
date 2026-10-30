#!/usr/bin/env bash
# ==============================================================================
# Bumiayu TUI - Uninstaller Script
# ==============================================================================

set -e

INSTALL_DIR="/usr/local/bin"
BIN_NAME="bumiayu"
UNINSTALLER_NAME="bumiayu-uninstall"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m'

echo -e "${YELLOW}🗑️ MEMULAI PENGHAPUSAN (UNINSTALL) BUMIAYU TUI...${NC}"

# Check root privileges
if [ "$EUID" -ne 0 ]; then
  echo -e "${RED}⚠️ Skrip membutuhkan hak akses root untuk menghapus instalasi.${NC}"
  echo -e "Silakan jalankan kembali dengan: ${GREEN}sudo bash $0${NC}"
  exit 1
fi

# 1. Remove binary
if [ -f "${INSTALL_DIR}/${BIN_NAME}" ]; then
  rm -f "${INSTALL_DIR}/${BIN_NAME}"
  echo -e "${GREEN}✓ Binary ${INSTALL_DIR}/${BIN_NAME} berhasil dihapus.${NC}"
else
  echo -e "${YELLOW}Binary ${INSTALL_DIR}/${BIN_NAME} sudah tidak ada di sistem.${NC}"
fi

# 2. Remove uninstaller script if installed
if [ -f "${INSTALL_DIR}/${UNINSTALLER_NAME}" ]; then
  rm -f "${INSTALL_DIR}/${UNINSTALLER_NAME}"
  echo -e "${GREEN}✓ Skrip uninstaller ${INSTALL_DIR}/${UNINSTALLER_NAME} dibersihkan.${NC}"
fi

# 3. Clean up .bashrc auto-start entries
echo -e "🧹 Memeriksa konfigurasi auto-start di ~/.bashrc..."
for bashrc in /root/.bashrc /home/*/.bashrc; do
  if [ -f "$bashrc" ] && grep -q "BUMIAYU_LOADED" "$bashrc"; then
    sed -i '/BUMIAYU_LOADED/d' "$bashrc"
    echo -e "${GREEN}✓ Auto-start dibersihkan dari ${bashrc}${NC}"
  fi
done

echo "------------------------------------------------------------------"
echo -e "${GREEN}🎉 Bumiayu TUI telah BERHASIL DIHAPUS BERSIH dari sistem Anda.${NC}"
echo -e "Jika di masa depan Anda ingin memasangnya kembali, cukup jalankan:"
echo -e "  ${CYAN}curl -fsSL https://raw.githubusercontent.com/faisalfa123/bumiayu-de/main/scripts/install.sh | sudo bash${NC}"
echo "------------------------------------------------------------------"
