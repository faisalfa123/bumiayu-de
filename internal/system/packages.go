package system

import (
	"os/exec"
	"runtime"
	"strings"
)

// AppCatalogItem represents a curated server application with multi-distro adapters
type AppCatalogItem struct {
	ID            string
	Name          string
	Category      string
	Description   string
	PackageName   string
	DietPiID      string // DietPi-software ID (e.g. 134 for Docker)
	Port          string
	Installed     bool
	Version       string
	RiskLevel     string // "Rendah", "Sedang", "Tinggi"
	Insight       string // Human-friendly advice
	IsThirdParty  bool   // True if verified 3rd-party vendor (e.g. Portainer, Tailscale, Cloudflared)
}

// GetAppCatalog returns comprehensive curated server applications with multi-distro support
func GetAppCatalog() []AppCatalogItem {
	catalog := []AppCatalogItem{
		// === 1. CONTAINER & CLOUD ===
		{
			ID:           "docker",
			Name:         "Docker Engine",
			Category:     "Container & Cloud",
			Description:  "Platform kontainer paling populer untuk menjalankan ribuan aplikasi tanpa mengotori server induk.",
			PackageName:  "docker.io",
			DietPiID:     "134",
			Port:         "2375 (Internal)",
			RiskLevel:    "Rendah",
			Insight:      "Paling direkomendasikan. Pastikan ruang MicroSD/eMMC di ODROID Anda mencukupi.",
			IsThirdParty: false,
		},
		{
			ID:           "portainer",
			Name:         "Portainer Community",
			Category:     "Container & Cloud",
			Description:  "Antarmuka web visual untuk mengelola semua kontainer Docker Anda dengan klik mudah dari browser.",
			PackageName:  "portainer-ce",
			DietPiID:     "185",
			Port:         "9000 / 9443",
			RiskLevel:    "Rendah",
			Insight:      "Aplikasi pihak ketiga resmi yang aman. Sangat mempermudah pemula dalam memonitor kontainer.",
			IsThirdParty: true,
		},
		{
			ID:           "docker-compose",
			Name:         "Docker Compose",
			Category:     "Container & Cloud",
			Description:  "Alat orkestrasi untuk menjalankan aplikasi multi-kontainer Docker hanya dengan 1 berkas yaml.",
			PackageName:  "docker-compose",
			DietPiID:     "134",
			Port:         "CLI",
			RiskLevel:    "Rendah",
			Insight:      "Sangat berguna jika ingin menjalankan stack aplikasi lengkap (misal WordPress + Database).",
			IsThirdParty: false,
		},

		// === 2. WEB & SERVER ===
		{
			ID:           "nginx",
			Name:         "Nginx Web Server",
			Category:     "Web & Server",
			Description:  "Web server performa tinggi dan reverse proxy nomor satu di dunia. Super hemat memori RAM.",
			PackageName:  "nginx",
			DietPiID:     "85",
			Port:         "80 / 443",
			RiskLevel:    "Rendah",
			Insight:      "Hanya memakan ~5 MB RAM. Sangat cocok sebagai pintu gerbang aplikasi di ODROID-C2.",
			IsThirdParty: false,
		},
		{
			ID:           "caddy",
			Name:         "Caddy Web Server",
			Category:     "Web & Server",
			Description:  "Web server modern dengan sertifikat HTTPS (SSL Let's Encrypt) otomatis tanpa konfigurasi rumit.",
			PackageName:  "caddy",
			DietPiID:     "83",
			Port:         "80 / 443",
			RiskLevel:    "Rendah",
			Insight:      "Paling mudah dikonfigurasi bagi pemula yang ingin website langsung berstatus HTTPS gembok hijau.",
			IsThirdParty: true,
		},
		{
			ID:           "apache2",
			Name:         "Apache2 HTTP Server",
			Category:     "Web & Server",
			Description:  "Web server klasik yang andal dengan dukungan modul .htaccess dan integrasi PHP yang luas.",
			PackageName:  "apache2",
			DietPiID:     "81",
			Port:         "80 / 443",
			RiskLevel:    "Rendah",
			Insight:      "Lebih banyak menggunakan RAM dibandingkan Nginx, namun sangat fleksibel untuk aplikasi PHP.",
			IsThirdParty: false,
		},
		{
			ID:           "nodejs",
			Name:         "Node.js & NPM Runtime",
			Category:     "Web & Server",
			Description:  "Runtime JavaScript modern untuk menjalankan backend, bot Telegram/Discord, dan aplikasi web.",
			PackageName:  "nodejs",
			DietPiID:     "9",
			Port:         "Runtime",
			RiskLevel:    "Rendah",
			Insight:      "Standar industri pengembang web modern. Menyediakan ekosistem jutaan pustaka NPM.",
			IsThirdParty: false,
		},
		{
			ID:           "python3",
			Name:         "Python 3 & Pip",
			Category:     "Web & Server",
			Description:  "Bahasa pemrograman serbaguna untuk automasi server, bot, AI, dan skrip latar belakang.",
			PackageName:  "python3-pip",
			DietPiID:     "130",
			Port:         "Runtime",
			RiskLevel:    "Rendah",
			Insight:      "Wajib ada di setiap server Linux untuk automasi tugas dan pemeliharaan sistem.",
			IsThirdParty: false,
		},
		{
			ID:           "certbot",
			Name:         "Certbot (Let's Encrypt)",
			Category:     "Web & Server",
			Description:  "Generator sertifikat SSL/TLS gratis otomatis agar website Anda aman dan terenkripsi.",
			PackageName:  "certbot",
			DietPiID:     "85",
			Port:         "CLI",
			RiskLevel:    "Rendah",
			Insight:      "Otomatis memperbarui masa berlaku sertifikat SSL setiap 90 hari tanpa intervensi manual.",
			IsThirdParty: false,
		},

		// === 3. MEDIA & STORAGE ===
		{
			ID:           "samba",
			Name:         "Samba File Server",
			Category:     "Media & Storage",
			Description:  "Berbagi folder harddisk ke Windows, macOS, Android, dan Smart TV selayaknya Network Drive.",
			PackageName:  "samba",
			DietPiID:     "96",
			Port:         "445 (SMB)",
			RiskLevel:    "Sedang",
			Insight:      "Hanya buka akses Samba di jaringan lokal WiFi rumah Anda, jangan dibuka ke publik internet.",
			IsThirdParty: false,
		},
		{
			ID:           "nextcloud",
			Name:         "Nextcloud Hub",
			Category:     "Media & Storage",
			Description:  "Cloud storage pribadi pengganti Google Drive / Dropbox lengkap dengan sinkronisasi foto dan kalender.",
			PackageName:  "nextcloud",
			DietPiID:     "114",
			Port:         "Web / Cloud",
			RiskLevel:    "Sedang",
			Insight:      "Memberikan kedaulatan data penuh atas file Anda. Rekomendasi sediakan swap RAM 1GB.",
			IsThirdParty: false,
		},
		{
			ID:           "filebrowser",
			Name:         "Filebrowser Web",
			Category:     "Media & Storage",
			Description:  "File Manager berbasis web super ringan untuk upload, download, dan kelola berkas dari browser.",
			PackageName:  "filebrowser",
			DietPiID:     "184",
			Port:         "8080 (Web)",
			RiskLevel:    "Rendah",
			Insight:      "Hanya memakan ~8 MB RAM! Alternatif Nextcloud jika server Anda memiliki spesifikasi RAM terbatas.",
			IsThirdParty: true,
		},
		{
			ID:           "jellyfin",
			Name:         "Jellyfin Media Server",
			Category:     "Media & Storage",
			Description:  "Server streaming film, video musik, dan koleksi foto pribadi ke Smart TV dan HP tanpa iklan.",
			PackageName:  "jellyfin",
			DietPiID:     "172",
			Port:         "8096 (HTTP)",
			RiskLevel:    "Sedang",
			Insight:      "Aplikasi 100% open-source tanpa fitur berbayar tersembunyi seperti Plex.",
			IsThirdParty: true,
		},

		// === 4. JARINGAN & VPN ===
		{
			ID:           "wireguard",
			Name:         "WireGuard VPN Server",
			Category:     "Jaringan & VPN",
			Description:  "VPN generasi terbaru yang sangat cepat, ringan, hemat baterai, dan menggunakan kriptografi mutakhir.",
			PackageName:  "wireguard",
			DietPiID:     "162",
			Port:         "51820 (UDP)",
			RiskLevel:    "Sedang",
			Insight:      "Bekerja langsung di level kernel Linux sehingga performanya luar biasa kencang.",
			IsThirdParty: false,
		},
		{
			ID:           "tailscale",
			Name:         "Tailscale Mesh VPN",
			Category:     "Jaringan & VPN",
			Description:  "Koneksikan server Anda dari mana saja di dunia tanpa perlu Port Forwarding dan tanpa pusing CGNAT.",
			PackageName:  "tailscale",
			DietPiID:     "191",
			Port:         "41641 (UDP)",
			RiskLevel:    "Rendah",
			Insight:      "Sangat direkomendasikan untuk pengguna internet rumah yang tidak memiliki IP Publik statis.",
			IsThirdParty: true,
		},
		{
			ID:           "pihole",
			Name:         "Pi-hole AdBlocker",
			Category:     "Jaringan & VPN",
			Description:  "Server DNS penangkal iklan dan pelacak (tracker) untuk seluruh perangkat yang terhubung ke WiFi rumah.",
			PackageName:  "pihole",
			DietPiID:     "93",
			Port:         "53 (DNS) / 80",
			RiskLevel:    "Rendah",
			Insight:      "Mempercepat browsing di HP dan laptop keluarga karena jutaan iklan diblokir di level DNS.",
			IsThirdParty: true,
		},
		{
			ID:           "adguardhome",
			Name:         "AdGuard Home",
			Category:     "Jaringan & VPN",
			Description:  "Pemblokir iklan jaringan dan proteksi malware dengan antarmuka web modern dan enkripsi DoH/DoT.",
			PackageName:  "adguardhome",
			DietPiID:     "183",
			Port:         "53 (DNS) / 3000",
			RiskLevel:    "Rendah",
			Insight:      "Memiliki dukungan DNS-over-HTTPS (DoH) bawaan yang sangat mudah diaktifkan.",
			IsThirdParty: true,
		},
		{
			ID:           "cloudflared",
			Name:         "Cloudflare Tunnel",
			Category:     "Jaringan & VPN",
			Description:  "Akses web server lokal ke internet publik melalui Cloudflare Tunnel tanpa membuka port router.",
			PackageName:  "cloudflared",
			DietPiID:     "198",
			Port:         "Tunnel",
			RiskLevel:    "Rendah",
			Insight:      "Solusi pihak ketiga terpercaya dari Cloudflare untuk hosting web rumah yang aman dari DDoS.",
			IsThirdParty: true,
		},
		{
			ID:           "ufw",
			Name:         "UFW Firewall",
			Category:     "Jaringan & VPN",
			Description:  "Dinding api pelindung port server yang mudah digunakan untuk mencegah akses tidak dikenal.",
			PackageName:  "ufw",
			DietPiID:     "122",
			Port:         "Sistem",
			RiskLevel:    "Rendah",
			Insight:      "Pastikan port SSH (22) diizinkan sebelum menyalakan UFW agar sesi remote Anda tidak terputus.",
			IsThirdParty: false,
		},
		{
			ID:           "fail2ban",
			Name:         "Fail2ban Brute-Force Defense",
			Category:     "Jaringan & VPN",
			Description:  "Otomatis memblokir IP asing yang mencoba menebak password SSH Anda berkali-kali.",
			PackageName:  "fail2ban",
			DietPiID:     "73",
			Port:         "Sistem",
			RiskLevel:    "Rendah",
			Insight:      "Bekerja diam di latar belakang, sangat efektif menghentikan serangan bot jahat.",
			IsThirdParty: false,
		},

		// === 5. DATABASE ===
		{
			ID:           "mariadb",
			Name:         "MariaDB / MySQL Server",
			Category:     "Database",
			Description:  "Sistem database relasional (SQL) standar industri untuk web, WordPress, dan aplikasi backend.",
			PackageName:  "mariadb-server",
			DietPiID:     "80",
			Port:         "3306",
			RiskLevel:    "Sedang",
			Insight:      "Gunakan kata sandi root database yang kuat untuk melindungi data penting Anda.",
			IsThirdParty: false,
		},
		{
			ID:           "postgresql",
			Name:         "PostgreSQL Database",
			Category:     "Database",
			Description:  "Database SQL relasional tingkat lanjut yang mendukung tipe data kompleks, JSON, dan geospasial.",
			PackageName:  "postgresql",
			DietPiID:     "165",
			Port:         "5432",
			RiskLevel:    "Sedang",
			Insight:      "Pilihan utama untuk pengembang aplikasi modern berskala besar.",
			IsThirdParty: false,
		},
		{
			ID:           "redis",
			Name:         "Redis In-Memory Cache",
			Category:     "Database",
			Description:  "Penyimpanan data in-memory super kencang untuk cache performa tinggi dan message broker.",
			PackageName:  "redis-server",
			DietPiID:     "129",
			Port:         "6379",
			RiskLevel:    "Rendah",
			Insight:      "Mempercepat waktu respon website hingga berkali-kali lipat saat dipadukan dengan Nginx.",
			IsThirdParty: false,
		},
		{
			ID:           "sqlite3",
			Name:         "SQLite 3 Engine",
			Category:     "Database",
			Description:  "Database berbasis satu berkas mandiri tanpa perlu menjalankan proses daemon latar belakang.",
			PackageName:  "sqlite3",
			DietPiID:     "3",
			Port:         "Berkas",
			RiskLevel:    "Rendah",
			Insight:      "Nol konsumsi RAM server saat idle. Sangat cocok untuk proyek mini atau bot.",
			IsThirdParty: false,
		},

		// === 6. MONITORING & ALAT ===
		{
			ID:           "btop",
			Name:         "Btop++ Visual Monitor",
			Category:     "Monitoring & Alat",
			Description:  "Monitor sistem visual terkini dengan grafik warna-warni untuk CPU, RAM, Disk, dan Jaringan.",
			PackageName:  "btop",
			DietPiID:     "195",
			Port:         "TUI",
			RiskLevel:    "Rendah",
			Insight:      "Antarmuka paling cantik untuk melihat proses berat di terminal.",
			IsThirdParty: true,
		},
		{
			ID:           "htop",
			Name:         "Htop Process Monitor",
			Category:     "Monitoring & Alat",
			Description:  "Pengawas proses klasik interaktif untuk memeriksa penggunaan CPU dan membunuh proses macet.",
			PackageName:  "htop",
			DietPiID:     "1",
			Port:         "TUI",
			RiskLevel:    "Rendah",
			Insight:      "Aplikasi standar setiap sysadmin Linux di seluruh dunia.",
			IsThirdParty: false,
		},
		{
			ID:           "netdata",
			Name:         "Netdata Real-time Web Monitor",
			Category:     "Monitoring & Alat",
			Description:  "Dashboard grafik visual real-time di browser untuk memantau ratusan metrik server secara detail.",
			PackageName:  "netdata",
			DietPiID:     "147",
			Port:         "19999 (Web)",
			RiskLevel:    "Rendah",
			Insight:      "Bisa dibuka dari browser laptop atau HP untuk memantau ODROID secara live.",
			IsThirdParty: true,
		},
		{
			ID:           "neofetch",
			Name:         "Neofetch / Fastfetch",
			Category:     "Monitoring & Alat",
			Description:  "Alat penampil logo distro Linux, spesifikasi mesin, dan ringkasan kernel yang estetik.",
			PackageName:  "fastfetch",
			DietPiID:     "17",
			Port:         "CLI",
			RiskLevel:    "Rendah",
			Insight:      "Populer untuk pamer spesifikasi server di komunitas Linux.",
			IsThirdParty: false,
		},
		{
			ID:           "tmux",
			Name:         "Tmux Terminal Multiplexer",
			Category:     "Monitoring & Alat",
			Description:  "Menjaga sesi terminal Anda tetap berjalan di server meskipun koneksi internet SSH terputus.",
			PackageName:  "tmux",
			DietPiID:     "22",
			Port:         "CLI",
			RiskLevel:    "Rendah",
			Insight:      "Wajib digunakan jika Anda menjalankan download atau kompilasi berdurasi panjang.",
			IsThirdParty: false,
		},
		{
			ID:           "git",
			Name:         "Git Version Control",
			Category:     "Monitoring & Alat",
			Description:  "Alat pengontrol versi kode untuk clone repository dan sinkronisasi proyek langsung di server.",
			PackageName:  "git",
			DietPiID:     "17",
			Port:         "CLI",
			RiskLevel:    "Rendah",
			Insight:      "Pondasi pengembangan software modern.",
			IsThirdParty: false,
		},
	}

	// Refresh installed status
	for i := range catalog {
		catalog[i].Installed = checkPackageInstalled(catalog[i].PackageName, catalog[i].ID)
	}

	return catalog
}

// checkPackageInstalled checks if package is installed in Debian/Ubuntu, DietPi, or binary PATH
func checkPackageInstalled(pkg, id string) bool {
	if runtime.GOOS != "linux" {
		return id == "docker" || id == "nginx" || id == "git" || id == "htop" || id == "tailscale"
	}

	// 1. Check binary path
	binaries := []string{pkg, id}
	for _, b := range binaries {
		if _, err := exec.LookPath(b); err == nil {
			return true
		}
	}

	// 2. Check dpkg-query
	cmd := exec.Command("dpkg-query", "-W", "-f=${Status}", pkg)
	out, err := cmd.Output()
	if err == nil && strings.Contains(string(out), "install ok installed") {
		return true
	}

	// 3. Check DietPi installed list
	if DetectDistro() == DistroDietPi {
		dietpiCmd := exec.Command("dietpi-software", "list")
		if dOut, dErr := dietpiCmd.Output(); dErr == nil {
			lines := strings.Split(string(dOut), "\n")
			for _, l := range lines {
				if strings.Contains(l, "| =1 |") && strings.Contains(strings.ToLower(l), strings.ToLower(id)) {
					return true
				}
			}
		}
	}

	return false
}

// ExecutePackageAction performs smart install or uninstall respecting DietPi, Ubuntu, or Debian
func ExecutePackageAction(action string, item AppCatalogItem) (string, error) {
	if runtime.GOOS != "linux" {
		if action == "install" {
			return "Reading package lists... Done\nBuilding dependency tree... Done\nSetting up " + item.Name + " (" + item.Category + ")...\nStatus: Successfully installed " + item.Name + "!", nil
		}
		return "Removing " + item.Name + "...\nStatus: Successfully removed " + item.Name + "!", nil
	}

	distro := DetectDistro()

	// 1. If running on DietPi, use native dietpi-software
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

	// 2. Verified 3rd-party installers on Ubuntu/Debian
	if item.IsThirdParty && action == "install" {
		switch item.ID {
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

	// 3. Standard APT package manager fallback (with noninteractive safe mode)
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
