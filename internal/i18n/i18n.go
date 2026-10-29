package i18n

// Language represents the active language
type Language string

const (
	ID Language = "id"
	EN Language = "en"
)

// Translation holds localized strings
type Translation struct {
	AppTitle       string
	AppSubtitle    string
	CurrentHost    string
	Uptime         string
	CpuLoad        string
	Memory         string
	Temperature    string
	Architecture   string
	Storage        string
	
	// Navigation Tabs
	TabDashboard   string
	TabApps        string
	TabServices    string
	TabSecurity    string
	TabUsers       string
	TabFiles       string
	TabTasks       string

	// Actions & Status
	Active         string
	Inactive       string
	Failed         string
	Running        string
	Stopped        string
	Installed      string
	NotInstalled   string
	Yes            string
	No             string
	Confirm        string
	Cancel         string
	ActionStart    string
	ActionStop     string
	ActionRestart  string
	ActionEnable   string
	ActionDisable  string
	ActionInstall  string
	ActionRemove   string
	ActionKill     string
	ActionGrant    string
	ActionRevoke   string

	// Senior Insight Headers
	SeniorInsightTitle  string
	SeniorTipsTitle     string

	// Insights
	InsightDashboard string
	InsightServices  string
	InsightSecurity  string
	InsightUsers     string
	InsightTasks     string
	InsightApps      string
	InsightFiles     string

	// Footer Help
	HelpNav        string
	HelpAction     string
	HelpLangToggle string
	HelpQuit       string
}

var translations = map[Language]Translation{
	ID: {
		AppTitle:       "BUMIAYU TUI",
		AppSubtitle:    "Desktop Environment Ringan & Elegan untuk ODROID & Server Linux",
		CurrentHost:    "Perangkat",
		Uptime:         "Masa Aktif",
		CpuLoad:        "Beban CPU",
		Memory:         "Penggunaan RAM",
		Temperature:    "Suhu Inti",
		Architecture:   "Arsitektur",
		Storage:        "Penyimpanan",

		TabDashboard:   "1. Dashboard",
		TabApps:        "2. Aplikasi & Paket",
		TabServices:    "3. Layanan Systemd",
		TabSecurity:    "4. Aturan & Firewall",
		TabUsers:       "5. Pengguna & Role",
		TabFiles:       "6. Berkas & Konfigurasi",
		TabTasks:       "7. Task Manager (OOM)",

		Active:         "Aktif",
		Inactive:       "Mati",
		Failed:         "Gagal",
		Running:        "Berjalan",
		Stopped:        "Berhenti",
		Installed:      "Terpasang",
		NotInstalled:   "Belum Terpasang",
		Yes:            "Ya",
		No:             "Tidak",
		Confirm:        "Konfirmasi",
		Cancel:         "Batal",
		ActionStart:    "[S] Nyalakan",
		ActionStop:     "[T] Hentikan",
		ActionRestart:  "[R] Mulai Ulang",
		ActionEnable:   "[E] Pasang Autostart",
		ActionDisable:  "[D] Matikan Autostart",
		ActionInstall:  "[I] Pasang Aplikasi",
		ActionRemove:   "[X] Hapus Paket",
		ActionKill:     "[K] Hentikan Paksa (Kill)",
		ActionGrant:    "[G] Berikan Akses Admin",
		ActionRevoke:   "[V] Cabut Akses Admin",

		SeniorInsightTitle: "💡 PANDUAN SENIOR ENGINEER",
		SeniorTipsTitle:    "🛡️ Rekomendasi Praktik Terbaik",

		InsightDashboard: "ODROID-C2 menggunakan SoC Amlogic S905 dengan 2GB RAM. Suhu ideal berada di bawah 65°C. Selalu sediakan Swap minimal 1GB agar proses latar belakang tidak menyebabkan freeze sistem.",
		InsightServices:  "Service 'Failed' menandakan crash saat startup. Cek log dengan tombol [L] Journalctl sebelum menyalakan ulang service agar tahu akar masalahnya.",
		InsightSecurity:  "Port yang terbuka adalah pintu masuk ke server Anda. Pastikan hanya port penting (SSH:22, Web:80/443) yang diizinkan lewat Firewall (UFW). Selalu gunakan kata sandi kuat atau SSH Key.",
		InsightUsers:     "Prinsip 'Least Privilege': Berikan hak Sudo hanya kepada user yang Anda percayai penuh. Akun biasa tanpa sudo tidak akan bisa merusak berkas sistem jika terkompromi.",
		InsightTasks:     "Linux Kernel akan memicu OOM Killer jika sisa RAM di bawah 5%. Urutkan proses berdasarkan memori dan hentikan proses dengan oom_score tinggi jika sistem mulai melambat.",
		InsightApps:      "Paket yang dipasang melalui APT terintegrasi langsung dengan pembaruan keamanan resmi Debian/Ubuntu. Hindari menginstal paket tak terverifikasi di server produksi.",
		InsightFiles:     "Berkas konfigurasi penting berada di /etc/. Berkas kunci SSH privat harus memiliki izin chmod 600, dan direktori .ssh memiliki izin 700 demi keamanan.",

		HelpNav:        "Tab/1-7: Navigasi | ↑↓/jk: Gulir | Enter/Spasi: Pilih",
		HelpAction:     "Aksi Tersedia di Menu",
		HelpLangToggle: "L: Ganti Bahasa (ID/EN)",
		HelpQuit:       "q / Ctrl+C: Keluar",
	},
	EN: {
		AppTitle:       "BUMIAYU TUI",
		AppSubtitle:    "Lightweight & Elegant TUI Desktop Environment for ODROID & Linux Servers",
		CurrentHost:    "Device Host",
		Uptime:         "System Uptime",
		CpuLoad:        "CPU Load",
		Memory:         "RAM Usage",
		Temperature:    "Core Temp",
		Architecture:   "Arch",
		Storage:        "Disk Space",

		TabDashboard:   "1. Dashboard",
		TabApps:        "2. Apps & Packages",
		TabServices:    "3. Systemd Services",
		TabSecurity:    "4. Security & Rules",
		TabUsers:       "5. Users & Roles",
		TabFiles:       "6. Files & Configs",
		TabTasks:       "7. Task Manager (OOM)",

		Active:         "Active",
		Inactive:       "Inactive",
		Failed:         "Failed",
		Running:        "Running",
		Stopped:        "Stopped",
		Installed:      "Installed",
		NotInstalled:   "Not Installed",
		Yes:            "Yes",
		No:             "No",
		Confirm:        "Confirm",
		Cancel:         "Cancel",
		ActionStart:    "[S] Start",
		ActionStop:     "[T] Stop",
		ActionRestart:  "[R] Restart",
		ActionEnable:   "[E] Enable Autostart",
		ActionDisable:  "[D] Disable Autostart",
		ActionInstall:  "[I] Install App",
		ActionRemove:   "[X] Remove Package",
		ActionKill:     "[K] Kill Process",
		ActionGrant:    "[G] Grant Sudo",
		ActionRevoke:   "[V] Revoke Sudo",

		SeniorInsightTitle: "💡 SENIOR ENGINEER INSIGHT",
		SeniorTipsTitle:    "🛡️ Best Practice Recommendations",

		InsightDashboard: "ODROID-C2 features an Amlogic S905 SoC with 2GB DDR3 RAM. Ideal operational temperature is under 65°C. Ensure at least 1GB Swap file to prevent sudden OOM freezes under heavy loads.",
		InsightServices:  "A 'Failed' service indicates crash on startup or bad configuration. Always inspect live logs via [L] Journalctl before blindly restarting to pinpoint the root cause.",
		InsightSecurity:  "Every open port is a potential attack vector. Only expose necessary services (e.g. SSH:22, HTTP:80/443). Enable UFW default incoming deny, and disable root SSH login.",
		InsightUsers:     "Principle of Least Privilege: Only grant Sudo to trusted accounts. Standard users without administrative roles cannot alter critical system files even if compromised.",
		InsightTasks:     "The Linux kernel triggers the OOM (Out Of Memory) Killer when free memory drops under ~5%. Sort by memory and terminate high oom_score processes when resources run tight.",
		InsightApps:      "System packages managed via APT receive official security upstream patches. Only install essential services to minimize memory footprint and security surface.",
		InsightFiles:     "Critical server configuration files reside in /etc/. Private SSH keys must have 600 permissions, and the .ssh folder must have 700 to safeguard remote access.",

		HelpNav:        "Tab/1-7: Workspaces | ↑↓/jk: Navigate | Enter/Space: Select",
		HelpAction:     "Available Actions",
		HelpLangToggle: "L: Switch Language (ID/EN)",
		HelpQuit:       "q / Ctrl+C: Quit",
	},
}

// Get returns the translation struct for the given language
func Get(lang Language) Translation {
	if t, ok := translations[lang]; ok {
		return t
	}
	return translations[ID]
}
