package system

import (
	"os/exec"
	"runtime"
	"strings"
)

// AppCatalogItem represents a curated server application
type AppCatalogItem struct {
	ID          string
	Name        string
	Category    string
	Description string
	DescriptionEN string
	PackageName string
	Port        string
	Installed   bool
	Version     string
	InstallCmd  string
	RemoveCmd   string
	RiskLevel   string // "Low", "Medium", "High"
	Insight     string // Human-friendly senior advice
}

// GetAppCatalog returns curated server applications suited for ODROID-C2 and Linux Servers
func GetAppCatalog() []AppCatalogItem {
	catalog := []AppCatalogItem{
		{
			ID:            "docker",
			Name:          "Docker Engine",
			Category:      "Containerization",
			Description:   "Menjalankan ribuan aplikasi siap pakai dalam kontainer terisolasi tanpa mengotori sistem.",
			DescriptionEN: "Run thousands of ready-made apps in isolated containers without cluttering the host.",
			PackageName:   "docker.io",
			Port:          "2375 (Internal)",
			InstallCmd:    "apt-get update && apt-get install -y docker.io",
			RemoveCmd:     "apt-get remove -y docker.io",
			RiskLevel:     "Low",
			Insight:       "Sangat fleksibel, namun perhatikan penggunaan disk di MicroSD / eMMC ODROID Anda.",
		},
		{
			ID:            "nginx",
			Name:          "Nginx Web Server",
			Category:      "Web / Proxy",
			Description:   "Web server super cepat dan reverse proxy untuk meneruskan traffic ke aplikasi Anda.",
			DescriptionEN: "High-performance web server and reverse proxy to route traffic to local apps.",
			PackageName:   "nginx",
			Port:          "80 / 443",
			InstallCmd:    "apt-get update && apt-get install -y nginx",
			RemoveCmd:     "apt-get remove -y nginx",
			RiskLevel:     "Low",
			Insight:       "Sangat hemat RAM (~5MB). Cocok dipadukan dengan SSL Let's Encrypt.",
		},
		{
			ID:            "wireguard",
			Name:          "WireGuard VPN",
			Category:      "Networking & VPN",
			Description:   "VPN generasi baru yang sangat cepat, ringan, dan aman untuk mengakses server dari luar rumah.",
			DescriptionEN: "Next-gen, lightweight, and ultra-secure VPN to access your home server remotely.",
			PackageName:   "wireguard",
			Port:          "51820 (UDP)",
			InstallCmd:    "apt-get update && apt-get install -y wireguard",
			RemoveCmd:     "apt-get remove -y wireguard",
			RiskLevel:     "Medium",
			Insight:       "Bekerja pada tingkat kernel Linux, meminimalkan latensi dan konsumsi CPU.",
		},
		{
			ID:            "samba",
			Name:          "Samba File Share",
			Category:      "Storage & Sharing",
			Description:   "Berbagi folder harddisk ke komputer Windows, Mac, dan Android selayaknya Network Drive.",
			DescriptionEN: "Share local folders with Windows, Mac, and mobile devices as a network drive.",
			PackageName:   "samba",
			Port:          "445 (SMB)",
			InstallCmd:    "apt-get update && apt-get install -y samba",
			RemoveCmd:     "apt-get remove -y samba",
			RiskLevel:     "Medium",
			Insight:       "Pastikan hanya membuka akses Samba di jaringan Wi-Fi lokal, jangan buka ke internet publik.",
		},
		{
			ID:            "ufw",
			Name:          "UFW Firewall",
			Category:      "Security",
			Description:   "Dinding api pelindung untuk menutup port yang tidak terpakai dari serangan luar.",
			DescriptionEN: "Intuitive firewall layer to close unused ports from unauthorized external probing.",
			PackageName:   "ufw",
			Port:          "System",
			InstallCmd:    "apt-get update && apt-get install -y ufw",
			RemoveCmd:     "apt-get remove -y ufw",
			RiskLevel:     "Low",
			Insight:       "Wajib diaktifkan. Pastikan port SSH (22) diizinkan sebelum menyalakan UFW agar tidak terputus.",
		},
		{
			ID:            "fail2ban",
			Name:          "Fail2ban Intrusion Defense",
			Category:      "Security",
			Description:   "Mendeteksi dan memblokir otomatis IP peretas yang mencoba menebak password SSH Anda berkali-kali.",
			DescriptionEN: "Automatically bans malicious IPs attempting brute-force password login over SSH.",
			PackageName:   "fail2ban",
			Port:          "System",
			InstallCmd:    "apt-get update && apt-get install -y fail2ban",
			RemoveCmd:     "apt-get remove -y fail2ban",
			RiskLevel:     "Low",
			Insight:       "Bekerja hening di latar belakang. Menghemat bandwidth dan mengamankan server.",
		},
		{
			ID:            "tailscale",
			Name:          "Tailscale Mesh VPN",
			Category:      "Remote Access",
			Description:   "Akses server ODROID dari mana saja tanpa perlu setting Port Forwarding di router.",
			DescriptionEN: "Connect to your ODROID from anywhere without port forwarding on your ISP router.",
			PackageName:   "tailscale",
			Port:          "41641 (UDP)",
			InstallCmd:    "curl -fsSL https://tailscale.com/install.sh | sh",
			RemoveCmd:     "apt-get remove -y tailscale",
			RiskLevel:     "Low",
			Insight:       "Sangat direkomendasikan jika ISP Anda menggunakan CGNAT (tidak punya IP publik statis).",
		},
		{
			ID:            "caddy",
			Name:          "Caddy Web Server",
			Category:      "Web / Proxy",
			Description:   "Web server otomatis mengurus sertifikat HTTPS tanpa repot konfigurasi certbot.",
			DescriptionEN: "Modern web server with automatic HTTPS and certificate renewal built-in.",
			PackageName:   "caddy",
			Port:          "80 / 443",
			InstallCmd:    "apt-get install -y debian-keyring debian-archive-keyring apt-transport-https && apt-get install -y caddy",
			RemoveCmd:     "apt-get remove -y caddy",
			RiskLevel:     "Low",
			Insight:       "Konfigurasi Caddyfile sangat singkat dan mudah dipahami bahkan untuk pemula.",
		},
		{
			ID:            "git",
			Name:          "Git Version Control",
			Category:      "Development",
			Description:   "Alat pengontrol versi kode untuk clone repository dan manajemen proyek langsung di server.",
			DescriptionEN: "Industry standard version control to clone repos and manage source trees.",
			PackageName:   "git",
			Port:          "CLI",
			InstallCmd:    "apt-get update && apt-get install -y git",
			RemoveCmd:     "apt-get remove -y git",
			RiskLevel:     "Low",
			Insight:       "Aplikasi wajib bagi setiap pengembang dan administrator server.",
		},
	}

	// Check installation status
	for i := range catalog {
		catalog[i].Installed = checkPackageInstalled(catalog[i].PackageName)
	}

	return catalog
}

// checkPackageInstalled checks if a package or binary exists
func checkPackageInstalled(pkg string) bool {
	if runtime.GOOS != "linux" {
		// Mock some as installed for preview
		return pkg == "git" || pkg == "ufw" || pkg == "nginx"
	}

	// 1. Check dpkg-query
	cmd := exec.Command("dpkg-query", "-W", "-f=${Status}", pkg)
	out, err := cmd.Output()
	if err == nil && strings.Contains(string(out), "install ok installed") {
		return true
	}

	// 2. Check binary path fallback
	if _, err := exec.LookPath(pkg); err == nil {
		return true
	}

	return false
}

// GetInstalledPackageCount returns total count of installed packages on the system
func GetInstalledPackageCount() int {
	if runtime.GOOS != "linux" {
		return 542 // Realistic mock count for display
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
