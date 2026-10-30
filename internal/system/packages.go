package system

import (
	"os/exec"
	"runtime"
	"strings"
)

// AppCatalogItem represents a curated server application with multi-distro adapters
type AppCatalogItem struct {
	ID           string
	Name         string
	Category     string
	Description  string
	PackageName  string // Debian/Ubuntu package name or binary
	DietPiID     string // DietPi-software ID (e.g. 215, 183, 134 or named slug)
	ArchPkg      string // Arch Linux pacman package
	Port         string
	Installed    bool
	Version      string
	RiskLevel    string // "Rendah", "Sedang", "Tinggi"
	Insight      string // Senior human-friendly advice
	IsThirdParty bool   // True if verified 3rd-party vendor (e.g. Tailscale, Portainer, Docker stack)
}

// GetAppCatalog returns comprehensive curated server applications with multi-distro support
func GetAppCatalog() []AppCatalogItem {
	catalog := []AppCatalogItem{
		// === 1. FOTO & CLOUD STORAGE ===
		{
			ID:           "immich",
			Name:         "Immich Photo & Video",
			Category:     "Foto & Cloud",
			Description:  "Solusi backup foto & video pribadi berkecepatan tinggi pengganti Google Photos dengan AI pencarian wajah.",
			PackageName:  "immich",
			DietPiID:     "215",
			ArchPkg:      "immich",
			Port:         "2283 (Web)",
			RiskLevel:    "Sedang",
			Insight:      "Paling dicari saat ini! Di DietPi otomatis terpasang sebagai service native (ID 215). Sediakan swap ~1GB untuk AI.",
			IsThirdParty: true,
		},
		{
			ID:           "nextcloud",
			Name:         "Nextcloud Hub",
			Category:     "Foto & Cloud",
			Description:  "Cloud storage pribadi pengganti Google Drive / Dropbox lengkap dengan sinkronisasi berkas, kontak, dan kalender.",
			PackageName:  "nextcloud",
			DietPiID:     "114",
			ArchPkg:      "nextcloud",
			Port:         "80 / 443 (Web)",
			RiskLevel:    "Sedang",
			Insight:      "Kedaulatan data penuh atas dokumen & foto Anda. Didukung aplikasi resmi Android & iOS.",
			IsThirdParty: false,
		},
		{
			ID:           "filebrowser",
			Name:         "Filebrowser Web",
			Category:     "Foto & Cloud",
			Description:  "File manager berbasis web super ringan untuk upload, unduh, dan kelola berkas langsung dari browser.",
			PackageName:  "filebrowser",
			DietPiID:     "184",
			ArchPkg:      "filebrowser",
			Port:         "8080 (Web)",
			RiskLevel:    "Rendah",
			Insight:      "Hanya memakan ~8 MB RAM! Pilihan utama jika server ODROID-C2 Anda ingin irit memori.",
			IsThirdParty: true,
		},
		{
			ID:           "syncthing",
			Name:         "Syncthing Peer-to-Peer",
			Category:     "Foto & Cloud",
			Description:  "Sinkronisasi berkas otomatis antar HP, laptop, dan server secara langsung terenkripsi tanpa server pihak ketiga.",
			PackageName:  "syncthing",
			DietPiID:     "syncthing",
			ArchPkg:      "syncthing",
			Port:         "8384 (Web UI)",
			RiskLevel:    "Rendah",
			Insight:      "Sangat hemat bandwidth & baterai HP. Folder otomatis sinkron saat terhubung ke WiFi rumah.",
			IsThirdParty: false,
		},
		{
			ID:           "samba",
			Name:         "Samba File Server (SMB)",
			Category:     "Foto & Cloud",
			Description:  "Berbagi folder penyimpanan server ke Windows Explorer, Android TV, Mac, dan iPhone sebagai Network Drive.",
			PackageName:  "samba",
			DietPiID:     "96",
			ArchPkg:      "samba",
			Port:         "445 (SMB)",
			RiskLevel:    "Sedang",
			Insight:      "Mudah diakses dari Windows File Explorer (\\\\ip-server\\dietpi). Buka hanya di LAN lokal.",
			IsThirdParty: false,
		},
		{
			ID:           "vaultwarden",
			Name:         "Vaultwarden (Bitwarden)",
			Category:     "Foto & Cloud",
			Description:  "Server pengelola kata sandi Bitwarden super ringan ditulis dalam Rust. Amankan ribuan password Anda.",
			PackageName:  "vaultwarden",
			DietPiID:     "183",
			ArchPkg:      "vaultwarden",
			Port:         "8001 (Web)",
			RiskLevel:    "Rendah",
			Insight:      "Kompatibel 100% dengan ekstensi browser Bitwarden & aplikasi Android. Membutuhkan koneksi HTTPS.",
			IsThirdParty: true,
		},

		// === 2. MEDIA & STREAMING ===
		{
			ID:           "jellyfin",
			Name:         "Jellyfin Media Server",
			Category:     "Media",
			Description:  "Server streaming film, serial, dan musik pribadi ke Smart TV, Android, dan browser tanpa langganan & tanpa iklan.",
			PackageName:  "jellyfin",
			DietPiID:     "178",
			ArchPkg:      "jellyfin-server",
			Port:         "8096 / 8097",
			RiskLevel:    "Sedang",
			Insight:      "Aplikasi 100% open-source murni tanpa paywall fitur premium layaknya Plex.",
			IsThirdParty: true,
		},
		{
			ID:           "transmission",
			Name:         "Transmission BitTorrent",
			Category:     "Media",
			Description:  "Downloader torrent server headless super ringan dengan kontrol antarmuka web yang bersih dan cepat.",
			PackageName:  "transmission-daemon",
			DietPiID:     "transmission",
			ArchPkg:      "transmission-cli",
			Port:         "9091 (Web)",
			RiskLevel:    "Rendah",
			Insight:      "Gunakan harddisk USB eksternal agar umur MicroSD/eMMC ODROID-C2 Anda tetap awet.",
			IsThirdParty: false,
		},
		{
			ID:           "qbittorrent",
			Name:         "qBittorrent Web UI",
			Category:     "Media",
			Description:  "Klien torrent bertenaga dengan antarmuka web lengkap, manajemen antrean, dan filter kategori.",
			PackageName:  "qbittorrent-nox",
			DietPiID:     "qbittorrent",
			ArchPkg:      "qbittorrent-nox",
			Port:         "8080 (Web)",
			RiskLevel:    "Rendah",
			Insight:      "Fitur sangat lengkap layaknya desktop qBittorrent namun dikontrol sepenuhnya lewat browser.",
			IsThirdParty: false,
		},
		{
			ID:           "minidlna",
			Name:         "MiniDLNA Media Streamer",
			Category:     "Media",
			Description:  "Server DLNA/UPnP sangat ringan agar Smart TV dan konsol game dapat memutar film dari server secara instan.",
			PackageName:  "minidlna",
			DietPiID:     "minidlna",
			ArchPkg:      "minidlna",
			Port:         "8200 (DLNA)",
			RiskLevel:    "Rendah",
			Insight:      "Konsumsi RAM di bawah 5 MB! Langsung terdeteksi otomatis di tab Network pada Smart TV Anda.",
			IsThirdParty: false,
		},

		// === 3. VPN, JARINGAN & KEAMANAN ===
		{
			ID:           "tailscale",
			Name:         "Tailscale Mesh VPN",
			Category:     "VPN & Jaringan",
			Description:  "Akses server dari mana saja di dunia tanpa buka port router, tanpa IP publik, dan tembus semua jenis CGNAT.",
			PackageName:  "tailscale",
			DietPiID:     "tailscale",
			ArchPkg:      "tailscale",
			Port:         "41641 (UDP)",
			RiskLevel:    "Rendah",
			Insight:      "Sangat direkomendasikan untuk WiFi Indihome/FirstMedia/Seluler yang tidak punya IP publik statis.",
			IsThirdParty: true,
		},
		{
			ID:           "wireguard",
			Name:         "WireGuard VPN Server",
			Category:     "VPN & Jaringan",
			Description:  "VPN generasi baru di level kernel Linux yang luar biasa kencang, hemat baterai, dan berenkripsi modern.",
			PackageName:  "wireguard",
			DietPiID:     "172",
			ArchPkg:      "wireguard-tools",
			Port:         "51820 (UDP)",
			RiskLevel:    "Sedang",
			Insight:      "Koneksi tunnel VPN paling kencang dengan ping terendah. Konfigurasi mudah dengan QR Code.",
			IsThirdParty: false,
		},
		{
			ID:           "cloudflared",
			Name:         "Cloudflare Tunnel",
			Category:     "VPN & Jaringan",
			Description:  "Publikasikan web server lokal ke internet dengan domain Anda sendiri tanpa port forwarding dan kebal DDoS.",
			PackageName:  "cloudflared",
			DietPiID:     "198",
			ArchPkg:      "cloudflared",
			Port:         "Tunnel",
			RiskLevel:    "Rendah",
			Insight:      "Solusi pihak ketiga resmi dari Cloudflare untuk menghubungkan server rumah ke dunia luar dengan aman.",
			IsThirdParty: true,
		},
		{
			ID:           "pihole",
			Name:         "Pi-hole AdBlocker",
			Category:     "VPN & Jaringan",
			Description:  "Server DNS penangkal iklan dan pelacak (tracker) untuk seluruh perangkat yang terhubung ke jaringan WiFi.",
			PackageName:  "pihole",
			DietPiID:     "93",
			ArchPkg:      "pi-hole-server",
			Port:         "53 (DNS) / 80",
			RiskLevel:    "Rendah",
			Insight:      "Membuat koneksi internet sekeluarga lebih ngebut dan hemat kuota karena iklan disaring di level DNS.",
			IsThirdParty: true,
		},
		{
			ID:           "adguardhome",
			Name:         "AdGuard Home DNS",
			Category:     "VPN & Jaringan",
			Description:  "Pemblokir iklan & situs berbahaya dengan antarmuka web modern, kontrol orang tua, dan enkripsi DoH/DoT.",
			PackageName:  "adguardhome",
			DietPiID:     "126",
			ArchPkg:      "adguardhome",
			Port:         "53 (DNS) / 3000",
			RiskLevel:    "Rendah",
			Insight:      "Mendukung DNS-over-HTTPS (DoH) bawaan dengan grafik statistik query yang sangat rapi & interaktif.",
			IsThirdParty: true,
		},
		{
			ID:           "ufw",
			Name:         "UFW Firewall",
			Category:     "VPN & Jaringan",
			Description:  "Dinding api pelindung port server yang mudah digunakan untuk mencegah akses berbahaya dari luar.",
			PackageName:  "ufw",
			DietPiID:     "ufw",
			ArchPkg:      "ufw",
			Port:         "Sistem",
			RiskLevel:    "Rendah",
			Insight:      "Pastikan port SSH (22) diizinkan sebelum menyalakan firewall agar tidak terkunci keluar.",
			IsThirdParty: false,
		},
		{
			ID:           "fail2ban",
			Name:         "Fail2ban Brute-Force Shield",
			Category:     "VPN & Jaringan",
			Description:  "Otomatis memblokir IP penyerang yang mencoba menebak kata sandi SSH server Anda berkali-kali.",
			PackageName:  "fail2ban",
			DietPiID:     "73",
			ArchPkg:      "fail2ban",
			Port:         "Sistem",
			RiskLevel:    "Rendah",
			Insight:      "Bekerja senyap di latar belakang. Sangat ampuh mengamankan server yang terhubung ke internet publik.",
			IsThirdParty: false,
		},

		// === 4. WEB, RUNTIME & DEVELOPER ===
		{
			ID:           "docker",
			Name:         "Docker Engine",
			Category:     "Web & Server",
			Description:  "Platform kontainer standar industri untuk menjalankan ribuan aplikasi mandiri tanpa merusak sistem induk.",
			PackageName:  "docker.io",
			DietPiID:     "134",
			ArchPkg:      "docker",
			Port:         "2375 (Daemon)",
			RiskLevel:    "Rendah",
			Insight:      "Paling direkomendasikan untuk sysadmin modern. Buka akses ke jutaan software di Docker Hub.",
			IsThirdParty: false,
		},
		{
			ID:           "portainer",
			Name:         "Portainer Community",
			Category:     "Web & Server",
			Description:  "Dashboard visual berbasis web untuk mengelola, membuat, dan memantau semua kontainer Docker dengan klik mudah.",
			PackageName:  "portainer-ce",
			DietPiID:     "185",
			ArchPkg:      "portainer-bin",
			Port:         "9000 / 9443",
			RiskLevel:    "Rendah",
			Insight:      "Sangat memanjakan pemula. Anda bisa deploy aplikasi hanya dengan copy-paste docker-compose yaml.",
			IsThirdParty: true,
		},
		{
			ID:           "docker-compose",
			Name:         "Docker Compose",
			Category:     "Web & Server",
			Description:  "Alat orkestrasi untuk menjalankan multi-kontainer aplikasi terpadu cukup dengan satu berkas konfigurasi.",
			PackageName:  "docker-compose",
			DietPiID:     "134",
			ArchPkg:      "docker-compose",
			Port:         "CLI",
			RiskLevel:    "Rendah",
			Insight:      "Wajib dipasang jika ingin menjalankan stack aplikasi rumit seperti WordPress + MariaDB.",
			IsThirdParty: false,
		},
		{
			ID:           "nginx",
			Name:         "Nginx Web Server",
			Category:     "Web & Server",
			Description:  "Web server berkecepatan tinggi dan reverse proxy nomor satu di dunia. Sangat hemat alokasi RAM.",
			PackageName:  "nginx",
			DietPiID:     "85",
			ArchPkg:      "nginx",
			Port:         "80 / 443",
			RiskLevel:    "Rendah",
			Insight:      "Konsumsi RAM hanya ~5 MB. Sangat cocok sebagai gerbang utama lalu lintas web server di ODROID-C2.",
			IsThirdParty: false,
		},
		{
			ID:           "caddy",
			Name:         "Caddy Web Server",
			Category:     "Web & Server",
			Description:  "Web server modern yang secara otomatis membuat dan memperbarui sertifikat HTTPS SSL Let's Encrypt.",
			PackageName:  "caddy",
			DietPiID:     "83",
			ArchPkg:      "caddy",
			Port:         "80 / 443",
			RiskLevel:    "Rendah",
			Insight:      "Konfigurasi Caddyfile paling manusiawi dan ringkas dibandingkan web server lainnya.",
			IsThirdParty: true,
		},
		{
			ID:           "apache2",
			Name:         "Apache2 HTTP Server",
			Category:     "Web & Server",
			Description:  "Web server klasik yang stabil dengan dukungan modular .htaccess dan integrasi PHP yang matang.",
			PackageName:  "apache2",
			DietPiID:     "81",
			ArchPkg:      "apache",
			Port:         "80 / 443",
			RiskLevel:    "Rendah",
			Insight:      "Sangat kompatibel untuk CMS warisan atau skrip PHP yang mengandalkan berkas .htaccess.",
			IsThirdParty: false,
		},
		{
			ID:           "nodejs",
			Name:         "Node.js & NPM Runtime",
			Category:     "Web & Server",
			Description:  "Runtime JavaScript modern untuk menjalankan aplikasi backend, bot Telegram/Discord, dan microservice.",
			PackageName:  "nodejs",
			DietPiID:     "9",
			ArchPkg:      "nodejs",
			Port:         "Runtime",
			RiskLevel:    "Rendah",
			Insight:      "Ekosistem pustaka terbesar di dunia melalui NPM. Siap untuk mengembangkan bot dan API.",
			IsThirdParty: false,
		},
		{
			ID:           "python3",
			Name:         "Python 3 & Pip Toolset",
			Category:     "Web & Server",
			Description:  "Bahasa pemrograman serbaguna untuk automasi server, AI, bot, web scraping, dan skrip pemeliharaan.",
			PackageName:  "python3-pip",
			DietPiID:     "130",
			ArchPkg:      "python-pip",
			Port:         "Runtime",
			RiskLevel:    "Rendah",
			Insight:      "Fondasi penting untuk banyak utilitas Linux dan automasi tugas berkala di server.",
			IsThirdParty: false,
		},
		{
			ID:           "certbot",
			Name:         "Certbot (Let's Encrypt SSL)",
			Category:     "Web & Server",
			Description:  "Generator sertifikat SSL/TLS gratis otomatis agar website Anda terenkripsi aman dengan simbol gembok hijau.",
			PackageName:  "certbot",
			DietPiID:     "85",
			ArchPkg:      "certbot",
			Port:         "CLI",
			RiskLevel:    "Rendah",
			Insight:      "Otomatis memperbarui masa berlaku sertifikat SSL setiap 90 hari melalui cron sistem.",
			IsThirdParty: false,
		},
		{
			ID:           "gitea",
			Name:         "Gitea / Forgejo Git Server",
			Category:     "Web & Server",
			Description:  "Layanan hosting repositori Git mandiri alternatif GitHub yang sangat ringan dan mudah digunakan.",
			PackageName:  "gitea",
			DietPiID:     "gitea",
			ArchPkg:      "gitea",
			Port:         "3000 (Web)",
			RiskLevel:    "Rendah",
			Insight:      "Hanya memakan ~40 MB RAM. Cocok untuk menyimpan kode proyek pribadi Anda secara mandiri di rumah.",
			IsThirdParty: true,
		},

		// === 5. DATABASE ===
		{
			ID:           "mariadb",
			Name:         "MariaDB / MySQL Server",
			Category:     "Database",
			Description:  "Database relasional (SQL) standar industri untuk website, blog WordPress, dan backend aplikasi.",
			PackageName:  "mariadb-server",
			DietPiID:     "mariadb",
			ArchPkg:      "mariadb",
			Port:         "3306",
			RiskLevel:    "Sedang",
			Insight:      "Pilihan paling umum dan didukung oleh hampir semua framework web di dunia.",
			IsThirdParty: false,
		},
		{
			ID:           "postgresql",
			Name:         "PostgreSQL Enterprise SQL",
			Category:     "Database",
			Description:  "Database objek-relasional tingkat tinggi yang terkenal dengan keandalan integritas data dan fitur JSON/Geospasial.",
			PackageName:  "postgresql",
			DietPiID:     "postgresql",
			ArchPkg:      "postgresql",
			Port:         "5432",
			RiskLevel:    "Sedang",
			Insight:      "Pilihan utama developer untuk aplikasi modern yang membutuhkan performa query data kompleks.",
			IsThirdParty: false,
		},
		{
			ID:           "redis",
			Name:         "Redis In-Memory Data Store",
			Category:     "Database",
			Description:  "Penyimpanan data in-memory super kencang untuk cache performa tinggi, manajemen session, dan message broker.",
			PackageName:  "redis-server",
			DietPiID:     "redis",
			ArchPkg:      "redis",
			Port:         "6379",
			RiskLevel:    "Rendah",
			Insight:      "Meningkatkan kecepatan respons web Anda hingga berkali-kali lipat dengan caching memori.",
			IsThirdParty: false,
		},
		{
			ID:           "sqlite3",
			Name:         "SQLite 3 Engine",
			Category:     "Database",
			Description:  "Mesin database berbasis satu berkas mandiri tanpa perlu menjalankan proses daemon latar belakang.",
			PackageName:  "sqlite3",
			DietPiID:     "sqlite3",
			ArchPkg:      "sqlite",
			Port:         "Berkas",
			RiskLevel:    "Rendah",
			Insight:      "Konsumsi RAM 0 saat tidak aktif! Sangat direkomendasikan untuk aplikasi mini atau database lokal.",
			IsThirdParty: false,
		},

		// === 6. MONITORING & ALAT ===
		{
			ID:           "dietpi-dashboard",
			Name:         "DietPi-Dashboard Web",
			Category:     "Monitoring & Alat",
			Description:  "Dashboard web resmi DietPi untuk memantau CPU, RAM, disk, suhu, dan proses dari browser perangkat apa pun.",
			PackageName:  "dietpi-dashboard",
			DietPiID:     "200",
			ArchPkg:      "dietpi-dashboard-bin",
			Port:         "5252 (Web)",
			RiskLevel:    "Rendah",
			Insight:      "Sangat ringan dan dibuat khusus untuk ekosistem DietPi (Software ID 200).",
			IsThirdParty: false,
		},
		{
			ID:           "uptimekuma",
			Name:         "Uptime Kuma Monitor",
			Category:     "Monitoring & Alat",
			Description:  "Alat pemantau status uptime website & server dengan antarmuka modern dan notifikasi Telegram/Discord.",
			PackageName:  "uptime-kuma",
			DietPiID:     "uptimekuma",
			ArchPkg:      "uptime-kuma",
			Port:         "3001 (Web)",
			RiskLevel:    "Rendah",
			Insight:      "Pantau apakah website atau koneksi internet rumah Anda sedang down secara real-time.",
			IsThirdParty: true,
		},
		{
			ID:           "btop",
			Name:         "Btop++ Visual Monitor",
			Category:     "Monitoring & Alat",
			Description:  "Monitor sistem visual terkini di terminal dengan grafik warna-warni untuk CPU, RAM, Disk, dan Jaringan.",
			PackageName:  "btop",
			DietPiID:     "195",
			ArchPkg:      "btop",
			Port:         "TUI",
			RiskLevel:    "Rendah",
			Insight:      "Antarmuka TUI paling estetik untuk melihat beban proses berat langsung di terminal SSH.",
			IsThirdParty: true,
		},
		{
			ID:           "htop",
			Name:         "Htop Process Monitor",
			Category:     "Monitoring & Alat",
			Description:  "Pengawas proses interaktif klasik untuk memeriksa penggunaan memori dan menghentikan proses yang macet.",
			PackageName:  "htop",
			DietPiID:     "1",
			ArchPkg:      "htop",
			Port:         "TUI",
			RiskLevel:    "Rendah",
			Insight:      "Alat standar wajib setiap sysadmin di seluruh dunia.",
			IsThirdParty: false,
		},
		{
			ID:           "netdata",
			Name:         "Netdata Real-time Web Monitor",
			Category:     "Monitoring & Alat",
			Description:  "Dashboard grafik visual real-time di browser untuk memeriksa performa server secara mendalam per detik.",
			PackageName:  "netdata",
			DietPiID:     "147",
			ArchPkg:      "netdata",
			Port:         "19999 (Web)",
			RiskLevel:    "Rendah",
			Insight:      "Menyajikan ribuan metrik perangkat keras tanpa membebani performa CPU.",
			IsThirdParty: true,
		},
		{
			ID:           "neofetch",
			Name:         "Fastfetch / Neofetch",
			Category:     "Monitoring & Alat",
			Description:  "Alat penampil logo distro Linux, spesifikasi mesin, dan ringkasan kernel yang estetik untuk terminal.",
			PackageName:  "fastfetch",
			DietPiID:     "17",
			ArchPkg:      "fastfetch",
			Port:         "CLI",
			RiskLevel:    "Rendah",
			Insight:      "Menampilkan ringkasan sistem yang indah saat pertama kali membuka terminal SSH.",
			IsThirdParty: false,
		},
		{
			ID:           "tmux",
			Name:         "Tmux Terminal Multiplexer",
			Category:     "Monitoring & Alat",
			Description:  "Menjaga sesi terminal Anda tetap berjalan di server meskipun koneksi SSH laptop terputus tiba-tiba.",
			PackageName:  "tmux",
			DietPiID:     "22",
			ArchPkg:      "tmux",
			Port:         "CLI",
			RiskLevel:    "Rendah",
			Insight:      "Wajib digunakan sebelum menjalankan perintah instalasi atau download berdurasi lama.",
			IsThirdParty: false,
		},
		{
			ID:           "git",
			Name:         "Git Version Control",
			Category:     "Monitoring & Alat",
			Description:  "Alat pengontrol versi kode untuk clone repositori dan sinkronisasi proyek langsung di server.",
			PackageName:  "git",
			DietPiID:     "17",
			ArchPkg:      "git",
			Port:         "CLI",
			RiskLevel:    "Rendah",
			Insight:      "Fondasi utama pengembangan perangkat lunak modern.",
			IsThirdParty: false,
		},

		// === 7. ENTERTAINMENT & MEDIA STREAMING TAMBAHAN ===
		{
			ID:           "plex",
			Name:         "Plex Media Server",
			Category:     "Media",
			Description:  "Pusat streaming film dan musik terpopuler di dunia dengan integrasi metadata poster otomatis.",
			PackageName:  "plexmediaserver",
			DietPiID:     "plex",
			ArchPkg:      "plex-media-server",
			Port:         "32400 (Web)",
			RiskLevel:    "Sedang",
			Insight:      "Mudah dihubungkan ke Smart TV, PlayStation, Xbox, dan HP Android/iPhone.",
			IsThirdParty: true,
		},
		{
			ID:           "emby",
			Name:         "Emby Media Server",
			Category:     "Media",
			Description:  "Media server alternatif Plex dengan kontrol hak akses pengguna dan folder keluarga yang kuat.",
			PackageName:  "emby-server",
			DietPiID:     "emby",
			ArchPkg:      "emby-server",
			Port:         "8096 (Web)",
			RiskLevel:    "Sedang",
			Insight:      "Pilihan seimbang antara fleksibilitas open-source dan aplikasi klien TV yang rapi.",
			IsThirdParty: true,
		},
		{
			ID:           "navidrome",
			Name:         "Navidrome Music Server",
			Category:     "Media",
			Description:  "Server streaming lagu MP3/FLAC pribadi pengganti Spotify yang super ringan (berbasis Go).",
			PackageName:  "navidrome",
			DietPiID:     "navidrome",
			ArchPkg:      "navidrome-bin",
			Port:         "4533 (Web)",
			RiskLevel:    "Rendah",
			Insight:      "Hanya memakan ~20 MB RAM! Kompatibel dengan aplikasi Subsonic/DSub di Android/iOS.",
			IsThirdParty: true,
		},
		{
			ID:           "audiobookshelf",
			Name:         "Audiobookshelf Server",
			Category:     "Media",
			Description:  "Server self-hosted untuk koleksi buku audio (audiobook) dan podcast dengan sinkronisasi posisi putar.",
			PackageName:  "audiobookshelf",
			DietPiID:     "audiobookshelf",
			ArchPkg:      "audiobookshelf-bin",
			Port:         "13378 (Web)",
			RiskLevel:    "Rendah",
			Insight:      "Aplikasi terlengkap untuk pecinta buku audio. Memiliki aplikasi resmi Android dan iOS.",
			IsThirdParty: true,
		},
		{
			ID:           "kodi",
			Name:         "Kodi Home Theater",
			Category:     "Media",
			Description:  "Media center pemutar video, musik, dan emulator game langsung ke layar TV melalui port HDMI.",
			PackageName:  "kodi",
			DietPiID:     "kodi",
			ArchPkg:      "kodi",
			Port:         "HDMI / 8080",
			RiskLevel:    "Sedang",
			Insight:      "Gunakan jika ODROID-C2 Anda terhubung langsung ke TV ruang keluarga dengan kabel HDMI.",
			IsThirdParty: false,
		},

		// === 8. AUTOMATION & ARR STACK (DOWNLOAD & MANAGEMENT) ===
		{
			ID:           "radarr",
			Name:         "Radarr Movie Manager",
			Category:     "Media",
			Description:  "Manajer koleksi film otomatis yang memantau rilis baru dan mengunduh via BitTorrent.",
			PackageName:  "radarr",
			DietPiID:     "radarr",
			ArchPkg:      "radarr",
			Port:         "7878 (Web)",
			RiskLevel:    "Sedang",
			Insight:      "Otomatis mencari film berkualitas terbaik sesuai preferensi resolusi Anda.",
			IsThirdParty: true,
		},
		{
			ID:           "sonarr",
			Name:         "Sonarr Series Tracker",
			Category:     "Media",
			Description:  "Pelacak dan pengunduh serial TV / Anime otomatis begitu episode baru dirilis di internet.",
			PackageName:  "sonarr",
			DietPiID:     "sonarr",
			ArchPkg:      "sonarr",
			Port:         "8989 (Web)",
			RiskLevel:    "Sedang",
			Insight:      "Terintegrasi sempurna dengan Transmission/qBittorrent dan Jellyfin/Plex.",
			IsThirdParty: true,
		},
		{
			ID:           "prowlarr",
			Name:         "Prowlarr Indexer Proxy",
			Category:     "Media",
			Description:  "Manajer integrasi torrent indexer terpadu untuk Sonarr, Radarr, dan Lidarr dalam satu pintu.",
			PackageName:  "prowlarr",
			DietPiID:     "prowlarr",
			ArchPkg:      "prowlarr",
			Port:         "9696 (Web)",
			RiskLevel:    "Rendah",
			Insight:      "Menghemat waktu konfigurasi indexer torrent di banyak aplikasi sekaligus.",
			IsThirdParty: true,
		},
		{
			ID:           "bazarr",
			Name:         "Bazarr Subtitle Finder",
			Category:     "Media",
			Description:  "Pengunduh teks terjemahan (subtitle) film dan serial otomatis (Subscene, OpenSubtitles).",
			PackageName:  "bazarr",
			DietPiID:     "bazarr",
			ArchPkg:      "bazarr",
			Port:         "6767 (Web)",
			RiskLevel:    "Rendah",
			Insight:      "Memastikan semua film yang baru diunduh langsung memiliki subtitle bahasa Indonesia.",
			IsThirdParty: true,
		},

		// === 9. SMART HOME, IOT & AUTOMATION ===
		{
			ID:           "homeassistant",
			Name:         "Home Assistant Core",
			Category:     "Web & Server",
			Description:  "Pusat otomasi rumah pintar (Smart Home) menghubungkan lampu WiFi, saklar pintar, dan sensor IoT.",
			PackageName:  "homeassistant",
			DietPiID:     "homeassistant",
			ArchPkg:      "home-assistant",
			Port:         "8123 (Web)",
			RiskLevel:    "Sedang",
			Insight:      "Solusi Smart Home paling populer dengan privasi lokal penuh tanpa ketergantungan cloud asing.",
			IsThirdParty: true,
		},
		{
			ID:           "nodered",
			Name:         "Node-RED Visual IoT",
			Category:     "Web & Server",
			Description:  "Alat pemrograman visual berbasis diagram alir (flow) untuk menghubungkan hardware IoT dan API.",
			PackageName:  "nodered",
			DietPiID:     "nodered",
			ArchPkg:      "node-red",
			Port:         "1880 (Web)",
			RiskLevel:    "Rendah",
			Insight:      "Sangat ramah pemula untuk membuat otomasi bot Telegram dan sensor tanpa banyak koding.",
			IsThirdParty: true,
		},
		{
			ID:           "mosquitto",
			Name:         "Mosquitto MQTT Broker",
			Category:     "Web & Server",
			Description:  "Broker protokol MQTT super ringan untuk komunikasi data antar perangkat mikrokontroler (ESP32/ESP8266).",
			PackageName:  "mosquitto",
			DietPiID:     "mosquitto",
			ArchPkg:      "mosquitto",
			Port:         "1883 (MQTT)",
			RiskLevel:    "Rendah",
			Insight:      "Konsumsi RAM di bawah 3 MB! Standar utama proyek IoT dunia.",
			IsThirdParty: false,
		},

		// === 10. SYSTEM MANAGEMENT & REMOTE CONTROL ===
		{
			ID:           "rustdesk",
			Name:         "RustDesk Relay Server",
			Category:     "VPN & Jaringan",
			Description:  "Server remote desktop mandiri alternatif AnyDesk / TeamViewer dengan enkripsi end-to-end milik sendiri.",
			PackageName:  "rustdesk-server",
			DietPiID:     "rustdeskserver",
			ArchPkg:      "rustdesk-server-bin",
			Port:         "21115 - 21119",
			RiskLevel:    "Rendah",
			Insight:      "Remote laptop kantor dan server tanpa batasan waktu dan bebas biaya lisensi komersial.",
			IsThirdParty: true,
		},
		{
			ID:           "webmin",
			Name:         "Webmin System Admin",
			Category:     "Monitoring & Alat",
			Description:  "Panel administrasi Linux berbasis web terlengkap untuk mengelola akun pengguna, cron job, dan disk.",
			PackageName:  "webmin",
			DietPiID:     "webmin",
			ArchPkg:      "webmin",
			Port:         "10000 (HTTPS)",
			RiskLevel:    "Sedang",
			Insight:      "Memudahkan konfigurasi server kompleks langsung lewat klik tombol di browser.",
			IsThirdParty: true,
		},
		{
			ID:           "cockpit",
			Name:         "Cockpit Web Console",
			Category:     "Monitoring & Alat",
			Description:  "Konsol server interaktif modern dari RedHat dengan monitor performa real-time dan terminal browser bawaan.",
			PackageName:  "cockpit",
			DietPiID:     "cockpit",
			ArchPkg:      "cockpit",
			Port:         "9090 (Web)",
			RiskLevel:    "Rendah",
			Insight:      "Antarmuka sangat modern, bersih, dan langsung terintegrasi dengan systemd Linux.",
			IsThirdParty: false,
		},
		{
			ID:           "glances",
			Name:         "Glances System Monitor",
			Category:     "Monitoring & Alat",
			Description:  "Monitor sistem berbasis Python dan Web yang menyajikan analisis beban I/O disk, sensor, dan jaringan.",
			PackageName:  "glances",
			DietPiID:     "glances",
			ArchPkg:      "glances",
			Port:         "61208 (Web)",
			RiskLevel:    "Rendah",
			Insight:      "Bisa dilihat langsung di terminal SSH atau diekspos sebagai antarmuka web HTML5.",
			IsThirdParty: true,
		},
	}

	// Refresh installed status
	for i := range catalog {
		catalog[i].Installed = checkPackageInstalled(catalog[i].PackageName, catalog[i].ID)
	}

	return catalog
}

// checkPackageInstalled checks if package is installed via fast lookup
func checkPackageInstalled(pkg, id string) bool {
	if runtime.GOOS != "linux" {
		return id == "docker" || id == "nginx" || id == "git" || id == "htop" || id == "tailscale" || id == "samba"
	}

	// 1. Fast PATH lookup without spawning shell
	binaries := []string{pkg, id}
	for _, b := range binaries {
		if _, err := exec.LookPath(b); err == nil {
			return true
		}
	}

	return false
}

// ExecutePackageAction performs smart install or uninstall respecting DietPi, Arch Linux, Ubuntu, or Debian
func ExecutePackageAction(action string, item AppCatalogItem) (string, error) {
	if runtime.GOOS != "linux" {
		if action == "install" {
			return "Reading package lists... Done\nBuilding dependency tree... Done\nSetting up " + item.Name + " (" + item.Category + ")...\nStatus: Successfully installed " + item.Name + "!", nil
		}
		return "Removing " + item.Name + "...\nStatus: Successfully removed " + item.Name + "!", nil
	}

	distro := DetectDistro()

	// 1. If running on DietPi, prioritize native dietpi-software
	if distro == DistroDietPi && item.DietPiID != "" {
		if _, err := exec.LookPath("dietpi-software"); err == nil {
			var cmd *exec.Cmd
			if action == "install" {
				cmd = exec.Command("dietpi-software", "install", item.DietPiID)
			} else {
				cmd = exec.Command("dietpi-software", "uninstall", item.DietPiID)
			}
			cmd.Env = append(cmd.Environ(), "DEBIAN_FRONTEND=noninteractive", "DIETPI_UNATTENDED=1")
			out, err := cmd.CombinedOutput()
			if err == nil || len(out) > 0 {
				return string(out), err
			}
		}
	}

	// 2. Arch Linux pacman adapter
	if distro == DistroArch && item.ArchPkg != "" {
		if _, err := exec.LookPath("pacman"); err == nil {
			var cmd *exec.Cmd
			if action == "install" {
				cmd = exec.Command("pacman", "-Sy", "--noconfirm", item.ArchPkg)
			} else {
				cmd = exec.Command("pacman", "-R", "--noconfirm", item.ArchPkg)
			}
			out, err := cmd.CombinedOutput()
			if err == nil || len(out) > 0 {
				return string(out), err
			}
		}
	}

	// 3. Verified 3rd-party installers on Ubuntu/Debian/Generic Linux
	if item.IsThirdParty && action == "install" {
		switch item.ID {
		case "immich":
			// If not on DietPi, deploy Immich official docker compose stack
			cmd := exec.Command("sh", "-c", "mkdir -p /opt/immich-app && cd /opt/immich-app && wget -O docker-compose.yml https://github.com/immich-app/immich/releases/latest/download/docker-compose.yml && wget -O .env https://github.com/immich-app/immich/releases/latest/download/example.env && docker compose up -d")
			out, err := cmd.CombinedOutput()
			return string(out), err
		case "vaultwarden":
			cmd := exec.Command("sh", "-c", "docker run -d --name vaultwarden -v /vw-data/:/data/ --restart always -p 8001:80 vaultwarden/server:latest")
			out, err := cmd.CombinedOutput()
			return string(out), err
		case "uptimekuma":
			cmd := exec.Command("sh", "-c", "docker run -d --restart=always -p 3001:3001 -v uptime-kuma:/app/data --name uptime-kuma louislam/uptime-kuma:latest")
			out, err := cmd.CombinedOutput()
			return string(out), err
		case "tailscale":
			cmd := exec.Command("sh", "-c", "curl -fsSL https://tailscale.com/install.sh | sh")
			out, err := cmd.CombinedOutput()
			return string(out), err
		case "caddy":
			cmd := exec.Command("sh", "-c", "apt-get install -y debian-keyring debian-archive-keyring apt-transport-https && curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg && curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | tee /etc/apt/sources.list.d/caddy-stable.list && apt-get update && apt-get install -y caddy")
			out, err := cmd.CombinedOutput()
			return string(out), err
		case "cloudflared":
			cmd := exec.Command("sh", "-c", "curl -fsSL https://pkg.cloudflare.com/cloudflare-main.gpg | tee /usr/share/keyrings/cloudflare-main.gpg >/dev/null && echo 'deb [signed-by=/usr/share/keyrings/cloudflare-main.gpg] https://pkg.cloudflare.com/cloudflared bullseye main' | tee /etc/apt/sources.list.d/cloudflared.list && apt-get update && apt-get install -y cloudflared")
			out, err := cmd.CombinedOutput()
			return string(out), err
		case "portainer":
			cmd := exec.Command("sh", "-c", "docker volume create portainer_data && docker run -d -p 9000:9000 -p 9443:9443 --name portainer --restart=always -v /var/run/docker.sock:/var/run/docker.sock -v portainer_data:/data portainer/portainer-ce:latest")
			out, err := cmd.CombinedOutput()
			return string(out), err
		case "filebrowser":
			cmd := exec.Command("sh", "-c", "curl -fsSL https://raw.githubusercontent.com/filebrowser/get/master/get.sh | bash")
			out, err := cmd.CombinedOutput()
			return string(out), err
		}
	}

	// 4. Standard APT package manager fallback (Debian/Ubuntu noninteractive)
	var cmd *exec.Cmd
	if action == "install" {
		cmd = exec.Command("apt-get", "install", "-y", item.PackageName)
	} else {
		cmd = exec.Command("apt-get", "remove", "-y", item.PackageName)
	}
	cmd.Env = append(cmd.Environ(), "DEBIAN_FRONTEND=noninteractive")
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// GetInstalledPackageCount returns total count of installed packages on the system
func GetInstalledPackageCount() int {
	if runtime.GOOS != "linux" {
		return 590
	}

	cmd := exec.Command("dpkg", "-l")
	out, err := cmd.Output()
	if err == nil {
		lines := strings.Split(string(out), "\n")
		count := 0
		for _, l := range lines {
			if strings.HasPrefix(l, "ii ") {
				count++
			}
		}
		if count > 0 {
			return count
		}
	}
	return 0
}
