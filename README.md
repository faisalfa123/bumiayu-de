# 🏝️ Bumiayu TUI

> **Desktop Environment Berbasis TUI (Terminal User Interface) Ultra-Ringan & Ramah SSH**  
> *Dirancang khusus untuk **ODROID-C2 (ARM64 Amlogic S905)** dan **Server Linux (AMD64 / x86_64)**.*

[![License: MIT](https://img.shields.io/badge/License-MIT-emerald.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.22%2B-blue.svg)](https://go.dev)
[![Architecture](https://img.shields.io/badge/Arch-ARM64%20|%20AMD64-cyan.svg)](https://github.com/bumiayu/bumiayu-de)
[![RAM Footprint](https://img.shields.io/badge/RAM%20Footprint-~12%20MB-success.svg)](https://github.com/bumiayu/bumiayu-de)
[![SSH Friendly](https://img.shields.io/badge/SSH-100%25%20Friendly-purple.svg)](https://github.com/bumiayu/bumiayu-de)

---

```
  ____                  _                    _____ _   _ ___ 
 | __ ) _   _ _ __ ___ (_) __ _ _   _ _   _ |_   _| | | |_ _|
 |  _ \| | | | '_ ` _ \| |/ _` | | | | | | |  | | | | | || | 
 | |_) | |_| | | | | | | | (_| | |_| | |_| |  | | | |_| || | 
 |____/ \__,_|_| |_| |_|_|\__,_|\__,_|\__,_|  |_|  \___/|___|
```

## ✨ Latar Belakang & Filosofi "Bumiayu"

**Bumiayu** melambangkan *Bumi yang Ayu* (indah, damai, dan bersahabat). 

Banyak pengguna server mini seperti **ODROID-C2** (quad-core ARM64 Cortex-A53 dengan 2GB DDR3 RAM) kesulitan ketika memasang desktop environment grafis konvensional (seperti XFCE, GNOME, atau LXDE) karena:
1. Memakan RAM 400 MB - 1.2 GB hanya untuk desktop, memicu **OOM (Out Of Memory) Freeze** saat menjalankan database atau kontainer.
2. Memerlukan forwarding X11 atau VNC yang lambat dan boros bandwidth internet saat diakses via remote.
3. Terminal murni tanpa antarmuka visual membuat pengguna pemula merasa cemas dan kesulitan mengaudit port, service, dan hak akses server.

**Bumiayu TUI** hadir sebagai jalan tengah yang sempurna: **Sensasi Desktop Environment penuh dalam jendela terminal**, hemat memori (< 15 MB RAM), responsif seketika via SSH, serta dilengkapi **"Senior Insight"** yang memandu pengguna awam selayaknya didampingi *Senior System Administrator*.

---

## 🖥️ 7 Ruang Kerja (Workspaces)

Bumiayu TUI membagi pengelolaan sistem ke dalam 7 panel terintegrasi:

| No | Ruang Kerja | Fungsi Utama |
|---|---|---|
| **1** | 📊 **Dashboard** | Pantau suhu SoC Amlogic ODROID-C2, clock CPU, RAM fisik & swap, load average, serta uptime. |
| **2** | 📦 **Aplikasi & Paket** | Audit jumlah paket terpasang (APT/dpkg) + 1-Click App Store (Docker, Nginx, WireGuard, Samba, Pi-hole, UFW, dll). |
| **3** | ⚙️ **Layanan Systemd** | Lihat status service real-time (`Running`, `Stopped`, `Failed`), Start, Stop, Restart, dan baca live logs Journalctl. |
| **4** | 🛡️ **Aturan & Firewall** | Pengendali Dinding Api (UFW), toggle proteksi, dan audit port terbuka (apakah publik atau lokal). |
| **5** | 👥 **Pengguna & Role** | Manajemen user, audit hak akses Sudo (Grant/Revoke Admin), dan inspeksi kunci SSH publik. |
| **6** | 📂 **Berkas & Konfigurasi** | Penjelajah direktori cepat (/etc, /var/log), inspeksi izin file (chmod visual), dan preview isi konfigurasi. |
| **7** | ⚡ **Task Manager (OOM)** | Pengawas proses dengan filter memori dan deteksi nilai `oom_score` Linux Kernel untuk mencegah freeze sistem. |

---

## 🚀 Panduan Instalasi Cepat (1-Baris)

### Di ODROID-C2 atau Server Linux (Ubuntu / Debian / Armbian):

Cukup jalankan satu perintah berikut di terminal Anda:

```bash
curl -fsSL https://raw.githubusercontent.com/bumiayu/bumiayu-de/main/scripts/install.sh | sudo bash
```

Skrip akan secara otomatis mendeteksi apakah mesin Anda adalah **ARM64** (ODROID-C2) atau **AMD64** (x86_64) dan memasang binary yang tepat ke `/usr/local/bin/bumiayu`.

---

## 🎮 Pintasan Keyboard (Hotkeys)

Bumiayu TUI dirancang untuk navigasi keyboard yang secepat kilat:

| Tombol | Tindakan |
|---|---|
| `1` - `7` | Langsung lompat ke Workspace 1 s.d. 7 |
| `Tab` / `Shift+Tab` | Pindah ke ruang kerja berikutnya / sebelumnya |
| `↑` / `↓` atau `k` / `j` | Gulir daftar item / layanan / proses |
| `Enter` / `Spasi` | Buka detail / lihat log Journalctl / buka folder |
| `L` | Ganti Bahasa Antarmuka secara instan (**Bahasa Indonesia** ⇄ **English**) |
| `S` | Nyalakan layanan (Start systemd unit) |
| `T` | Hentikan layanan (Stop systemd unit) |
| `R` | Mulai ulang layanan (Restart systemd unit) |
| `U` | Nyalakan / Matikan Dinding Api (Toggle UFW Firewall) |
| `G` / `V` | Berikan (Grant) / Cabut (Revoke) hak akses Sudo pada pengguna |
| `K` | Hentikan paksa proses (Kill Process via SIGTERM) |
| `Esc` | Tutup jendela log / kembali ke menu |
| `q` atau `Ctrl+C` | Keluar dari Bumiayu TUI |

---

## 🌐 Integrasi Otomatis Saat Login SSH

Agar server ODROID-C2 Anda langsung menyapa dengan tampilan Bumiayu TUI setiap kali Anda melakukan SSH, tambahkan baris berikut ke `~/.bashrc`:

```bash
# Tambahkan ke baris paling bawah ~/.bashrc
if [ -n "$SSH_CLIENT" ] || [ -n "$SSH_TTY" ]; then
    if [ -z "$BUMIAYU_LOADED" ]; then
        export BUMIAYU_LOADED=1
        exec bumiayu
    fi
fi
```

Ketika Anda keluar dari Bumiayu TUI (tekan `q`), sesi terminal SSH Anda akan tetap aktif seperti biasa.

---

## 🛠️ Kompilasi Mandiri (Cross-Compilation)

Jika Anda ingin membangun binary sendiri dari kode sumber:

```bash
# Clone repository
git clone https://github.com/bumiayu/bumiayu-de.git
cd bumiayu-de

# Jalankan pengujian unit
make test

# Kompilasi untuk ODROID-C2 (ARM64 Linux)
make build-arm64

# Kompilasi untuk Server AMD64 (x86_64 Linux)
make build-amd64

# Kompilasi seluruh arsitektur sekaligus
make build-all
```

Hasil binary statis siap pakai akan tersedia di folder `dist/`.

---

## 📄 Lisensi

Proyek ini dirilis di bawah lisensi terbuka [MIT License](LICENSE).
Bebas dimodifikasi, didistribusikan, dan digunakan untuk server pribadi maupun enterprise.
