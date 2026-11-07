package launcher

import (
	"math"
	"strings"

	"github.com/bumiayu/bumiayu-de/internal/ui/styles"
	"github.com/charmbracelet/lipgloss"
)

// ActionType categorizes what an action does
type ActionType int

const (
	ActionService ActionType = iota
	ActionProcess
	ActionMacro
	ActionNavigate
	ActionApp
	ActionTheme
)

// LauncherItem represents an actionable command or search target
type LauncherItem struct {
	ID          string
	Title       string
	Description string
	Category    string
	Type        ActionType
	Command     string
	Keywords    []string
}

// SpotlightState keeps track of launcher query, cursor, ghost-text, and results
type SpotlightState struct {
	Active     bool
	Query      string
	Cursor     int
	GhostText  string // Ghost text completion ahead of cursor
	Results    []LauncherItem
	AllActions []LauncherItem
}

// NewSpotlightState creates initial launcher state with pre-populated macros and catalog
func NewSpotlightState() SpotlightState {
	s := SpotlightState{
		AllActions: getDefaultActions(),
	}
	s.Filter()
	return s
}

func getDefaultActions() []LauncherItem {
	return []LauncherItem{
		{
			ID:          "macro_clean",
			Title:       ">clean (Bersihkan RAM Cache)",
			Description: "Jalankan drop_caches untuk bebaskan memori RAM yang tertahan",
			Category:    "Sistem",
			Type:        ActionMacro,
			Command:     "clean_ram",
			Keywords:    []string{"clean", "bersih", "ram", "cache", "drop_caches", "memori"},
		},
		{
			ID:          "macro_ports",
			Title:       ">ports (Audit Port Terbuka)",
			Description: "Inspeksi port lokal & publik yang sedang mendengarkan (ss/netstat)",
			Category:    "Jaringan",
			Type:        ActionMacro,
			Command:     "audit_ports",
			Keywords:    []string{"ports", "port", "jaringan", "open", "audit", "listen"},
		},
		{
			ID:          "macro_ufw",
			Title:       ">ufw (Toggle Firewall Dinding Api)",
			Description: "Aktifkan atau matikan proteksi UFW Firewall secara instan",
			Category:    "Keamanan",
			Type:        ActionMacro,
			Command:     "toggle_ufw",
			Keywords:    []string{"ufw", "firewall", "dinding api", "security", "toggle"},
		},
		{
			ID:          "macro_kill_oom",
			Title:       ">kill-oom (Hentikan Proses Terberat)",
			Description: "Temukan dan kirim SIGTERM ke proses dengan skor OOM tertinggi",
			Category:    "Proses",
			Type:        ActionMacro,
			Command:     "kill_highest_oom",
			Keywords:    []string{"kill", "oom", "bunuh", "heavy", "freeze", "proses"},
		},
		{
			ID:          "macro_theme",
			Title:       ">theme (Ganti Tema Retro)",
			Description: "Beralih palet warna: 16-Bit Cyber ⇄ DOS Norton ⇄ GameBoy",
			Category:    "Tampilan",
			Type:        ActionTheme,
			Command:     "next_theme",
			Keywords:    []string{"theme", "tema", "warna", "color", "cyber", "dos", "gameboy"},
		},
		{
			ID:          "macro_layout_tiling",
			Title:       ">layout tiling (Tampilan Grid 2x2)",
			Description: "Ganti tampilan ke mode grid 4 panel seimbang (Default)",
			Category:    "Tampilan",
			Type:        ActionTheme,
			Command:     "layout_tiling",
			Keywords:    []string{"layout", "tiling", "grid", "2x2", "tampilan", "model"},
		},
		{
			ID:          "macro_layout_master",
			Title:       ">layout master (Tampilan Master-Stack)",
			Description: "Ganti tampilan: 1 panel telemetri besar kiri + 3 kartu di kanan",
			Category:    "Tampilan",
			Type:        ActionTheme,
			Command:     "layout_master",
			Keywords:    []string{"layout", "master", "stack", "kolom", "tampilan", "model"},
		},
		{
			ID:          "macro_layout_split",
			Title:       ">layout split (Tampilan Horizontal Split)",
			Description: "Ganti tampilan ke mode belah atas-bawah (Telemetri atas, Proses bawah)",
			Category:    "Tampilan",
			Type:        ActionTheme,
			Command:     "layout_split",
			Keywords:    []string{"layout", "split", "horizontal", "tampilan", "model"},
		},
		{
			ID:          "macro_fps_adaptive",
			Title:       ">fps adaptive (Mode Adaptif 24 FPS / Hemat Daya)",
			Description: "24 FPS saat navigasi aktif, otomatis hemat daya (2 FPS) saat idle",
			Category:    "Performa",
			Type:        ActionMacro,
			Command:     "fps_adaptive",
			Keywords:    []string{"fps", "adaptive", "hemat", "eco", "performa", "24"},
		},
		{
			ID:          "macro_fps_60",
			Title:       ">fps 60 (Mode Ultra Smooth 60 FPS)",
			Description: "Kunci frame rate tinggi 60 FPS terus menerus (Untuk PC / Server bertenaga)",
			Category:    "Performa",
			Type:        ActionMacro,
			Command:     "fps_60",
			Keywords:    []string{"fps", "60", "smooth", "ultra", "cepat"},
		},
		{
			ID:          "macro_fps_eco",
			Title:       ">fps eco (Mode Statis Super Hemat 5 FPS)",
			Description: "Kunci frame rate rendah 5 FPS untuk server hemat CPU & SSH lambat",
			Category:    "Performa",
			Type:        ActionMacro,
			Command:     "fps_eco",
			Keywords:    []string{"fps", "eco", "hemat", "low", "battery", "sbc"},
		},
		{
			ID:          "macro_autoclean_15m",
			Title:       ">autoclean 15m (Auto Clean RAM Tiap 15 Menit)",
			Description: "Otomatis bersihkan cache RAM dan dentries setiap 15 menit",
			Category:    "Sistem",
			Type:        ActionMacro,
			Command:     "autoclean_15m",
			Keywords:    []string{"autoclean", "15m", "clean", "ram", "otomatis", "15"},
		},
		{
			ID:          "macro_autoclean_30m",
			Title:       ">autoclean 30m (Auto Clean RAM Tiap 30 Menit)",
			Description: "Otomatis bersihkan cache RAM dan dentries setiap 30 menit",
			Category:    "Sistem",
			Type:        ActionMacro,
			Command:     "autoclean_30m",
			Keywords:    []string{"autoclean", "30m", "clean", "ram", "otomatis", "30"},
		},
		{
			ID:          "macro_autoclean_1h",
			Title:       ">autoclean 1h (Auto Clean RAM Tiap 1 Jam)",
			Description: "Otomatis bersihkan cache RAM dan dentries setiap 1 jam",
			Category:    "Sistem",
			Type:        ActionMacro,
			Command:     "autoclean_1h",
			Keywords:    []string{"autoclean", "1h", "clean", "ram", "otomatis", "jam"},
		},
		{
			ID:          "macro_autoclean_off",
			Title:       ">autoclean off (Matikan Auto Clean RAM)",
			Description: "Hentikan pembersihan RAM berkala (hanya bersihkan manual)",
			Category:    "Sistem",
			Type:        ActionMacro,
			Command:     "autoclean_off",
			Keywords:    []string{"autoclean", "off", "mati", "disable", "stop"},
		},
		{
			ID:          "macro_bg_auto",
			Title:       ">bg auto (Animasi Dinamis Jam Server)",
			Description: "Ganti animasi otomatis: Pagi Slamet, Siang Sawah, Sore Senja, Malam Bintang",
			Category:    "Tampilan",
			Type:        ActionTheme,
			Command:     "bg_auto",
			Keywords:    []string{"bg", "background", "animasi", "auto", "waktu", "jam", "dinamis", "wallpaper"},
		},
		{
			ID:          "macro_bg_stars",
			Title:       ">bg stars (Langit Bintang Hening & Minimalis)",
			Description: "Titik bintang kosmik halus yang melayang sangat lambat dan tenang (Sangat nyaman untuk mata)",
			Category:    "Tampilan",
			Type:        ActionTheme,
			Command:     "bg_stars",
			Keywords:    []string{"bg", "stars", "bintang", "tenang", "hening", "minimalis", "relax"},
		},
		{
			ID:          "macro_bg_matrix",
			Title:       ">bg matrix (Digital Matrix Stream Hijau Gelap)",
			Description: "Aliran karakter digital hijau halus dan gelap di latar desktop tanpa mengganggu teks",
			Category:    "Tampilan",
			Type:        ActionTheme,
			Command:     "bg_matrix",
			Keywords:    []string{"bg", "matrix", "hacker", "kode", "digital", "stream", "hijau"},
		},
		{
			ID:          "macro_bg_sawah",
			Title:       ">bg sawah (Sawah & Kerbau Berjalan)",
			Description: "Animasi pixel 24 FPS: Petani membajak sawah terasering subur & pematang",
			Category:    "Tampilan",
			Type:        ActionTheme,
			Command:     "bg_sawah",
			Keywords:    []string{"bg", "background", "sawah", "paddy", "kerbau", "petani", "animasi", "runner"},
		},
		{
			ID:          "macro_bg_gunung",
			Title:       ">bg gunung (Gunung Slamet & Awan)",
			Description: "Animasi parallax 24 FPS: Siluet Gunung Slamet megah dengan awan berarak & hutan pinus",
			Category:    "Tampilan",
			Type:        ActionTheme,
			Command:     "bg_gunung",
			Keywords:    []string{"bg", "background", "gunung", "slamet", "awan", "mountain", "animasi", "parallax"},
		},
		{
			ID:          "macro_bg_sakalibel",
			Title:       ">bg sakalibel (Kereta Api Jembatan Sakalibel)",
			Description: "Animasi 24 FPS: Kereta Api Kamandaka melintasi ikon jembatan bersejarah Bumiayu",
			Category:    "Tampilan",
			Type:        ActionTheme,
			Command:     "bg_sakalibel",
			Keywords:    []string{"bg", "background", "sakalibel", "kereta", "train", "jembatan", "bridge", "animasi"},
		},
		{
			ID:          "macro_bg_senja",
			Title:       ">bg senja (Lentera Senja Bumiayu)",
			Description: "Animasi 24 FPS: Lentera merah bergoyang di langit senja & warung teh poci gula batu",
			Category:    "Tampilan",
			Type:        ActionTheme,
			Command:     "bg_senja",
			Keywords:    []string{"bg", "background", "senja", "lentera", "dusk", "sunset", "teh poci", "animasi"},
		},
		{
			ID:          "macro_bg_malam",
			Title:       ">bg malam (Langit Malam & Meteor)",
			Description: "Animasi 24 FPS: Bintang berkelip, meteor melintas cepat, & kunang-kunang malam",
			Category:    "Tampilan",
			Type:        ActionTheme,
			Command:     "bg_malam",
			Keywords:    []string{"bg", "background", "malam", "bintang", "star", "night", "meteor", "animasi"},
		},
		{
			ID:          "macro_bg_curug",
			Title:       ">bg curug (Curug Cipendok Air Terjun)",
			Description: "Animasi 24 FPS: Percikan air terjun Curug Cipendok mengalir jernih di lembah",
			Category:    "Tampilan",
			Type:        ActionTheme,
			Command:     "bg_curug",
			Keywords:    []string{"bg", "background", "curug", "cipendok", "waterfall", "air terjun", "animasi"},
		},
		{
			ID:          "macro_bg_off",
			Title:       ">bg off (Matikan Animasi Background)",
			Description: "Nonaktifkan panorama animasi latar belakang (Kembali ke mode minimalis gelap murni)",
			Category:    "Tampilan",
			Type:        ActionTheme,
			Command:     "bg_off",
			Keywords:    []string{"bg", "background", "off", "mati", "disable", "nonaktif", "clean", "dark"},
		},
		{
			ID:          "macro_bg_dim_low",
			Title:       ">bg dim low (Opasitas Latar Ringan / Lebih Terang)",
			Description: "Animasi background lebih nyata terlihat tembus pandang di balik panel jendela",
			Category:    "Tampilan",
			Type:        ActionTheme,
			Command:     "bg_dim_low",
			Keywords:    []string{"bg", "dim", "low", "opasitas", "terang", "transparan"},
		},
		{
			ID:          "macro_bg_dim_med",
			Title:       ">bg dim med (Opasitas Latar Sedang / Seimbang)",
			Description: "Tingkat redup seimbang default: siluet estetik tanpa mengaburkan teks sistem",
			Category:    "Tampilan",
			Type:        ActionTheme,
			Command:     "bg_dim_med",
			Keywords:    []string{"bg", "dim", "med", "opasitas", "sedang", "seimbang", "default"},
		},
		{
			ID:          "macro_bg_dim_high",
			Title:       ">bg dim high (Opasitas Latar Kuat / Sangat Redup)",
			Description: "Latar animasi di dalam panel sangat samar dan redup untuk fokus maksimal pada data",
			Category:    "Tampilan",
			Type:        ActionTheme,
			Command:     "bg_dim_high",
			Keywords:    []string{"bg", "dim", "high", "opasitas", "redup", "gelap", "fokus"},
		},
		{
			ID:          "nav_telemetry",
			Title:       ">jump telemetry (Panel Hardware)",
			Description: "Fokus ke panel telemetri CPU, RAM, Suhu SoC, dan Load",
			Category:    "Navigasi",
			Type:        ActionNavigate,
			Command:     "jump_0",
			Keywords:    []string{"telemetry", "hardware", "cpu", "suhu", "temp", "panel 1"},
		},
		{
			ID:          "nav_services",
			Title:       ">jump services (Panel Layanan)",
			Description: "Fokus ke panel audit dan kontrol service systemd",
			Category:    "Navigasi",
			Type:        ActionNavigate,
			Command:     "jump_1",
			Keywords:    []string{"services", "layanan", "systemd", "panel 2"},
		},
		{
			ID:          "nav_processes",
			Title:       ">jump processes (Panel Proses OOM)",
			Description: "Fokus ke panel task manager dan pemantau proses",
			Category:    "Navigasi",
			Type:        ActionNavigate,
			Command:     "jump_2",
			Keywords:    []string{"processes", "proses", "task", "oom", "panel 3"},
		},
		{
			ID:          "nav_network",
			Title:       ">jump network (Panel Jaringan & Port)",
			Description: "Fokus ke panel audit port dan keamanan firewall",
			Category:    "Navigasi",
			Type:        ActionNavigate,
			Command:     "jump_3",
			Keywords:    []string{"network", "jaringan", "port", "firewall", "panel 4"},
		},
		{
			ID:          "svc_nginx",
			Title:       "nginx.service (Web Server)",
			Description: "Kelola web server Nginx: Start, Stop, atau Restart",
			Category:    "Systemd",
			Type:        ActionService,
			Command:     "service:nginx.service",
			Keywords:    []string{"nginx", "web", "http", "proxy", "server"},
		},
		{
			ID:          "svc_ssh",
			Title:       "ssh.service (SSH Remote Daemon)",
			Description: "Layanan remote secure shell server OpenSSH",
			Category:    "Systemd",
			Type:        ActionService,
			Command:     "service:ssh.service",
			Keywords:    []string{"ssh", "sshd", "remote", "terminal", "login"},
		},
		{
			ID:          "svc_docker",
			Title:       "docker.service (Container Engine)",
			Description: "Layanan runtime kontainer aplikasi Docker",
			Category:    "Systemd",
			Type:        ActionService,
			Command:     "service:docker.service",
			Keywords:    []string{"docker", "container", "dockerd", "kontainer"},
		},
		{
			ID:          "app_wireguard",
			Title:       "WireGuard VPN (App Store)",
			Description: "Protokol VPN modern berkecepatan tinggi & ramah memori",
			Category:    "Aplikasi",
			Type:        ActionApp,
			Command:     "install:wireguard",
			Keywords:    []string{"wireguard", "vpn", "tunnel", "jaringan", "secure"},
		},
		{
			ID:          "app_samba",
			Title:       "Samba File Sharing (App Store)",
			Description: "Berbagi berkas jaringan lokal kompatibel Windows & Mac",
			Category:    "Aplikasi",
			Type:        ActionApp,
			Command:     "install:samba",
			Keywords:    []string{"samba", "smb", "share", "nas", "berkas"},
		},
	}
}

// AddDynamicServices adds real detected systemd services to searchable catalog
func (s *SpotlightState) SyncDynamicItems(services []string, processes []string) {
	// Re-initialize with default actions
	base := getDefaultActions()

	for _, svc := range services {
		base = append(base, LauncherItem{
			ID:          "dyn_svc_" + svc,
			Title:       svc,
			Description: "Layanan sistem Linux (Systemd unit)",
			Category:    "Systemd",
			Type:        ActionService,
			Command:     "service:" + svc,
			Keywords:    []string{svc, "service", "systemd", strings.TrimSuffix(svc, ".service")},
		})
	}

	for _, proc := range processes {
		base = append(base, LauncherItem{
			ID:          "dyn_proc_" + proc,
			Title:       "Proses: " + proc,
			Description: "Proses aktif di sistem (Kirim sinyal / Kill)",
			Category:    "Proses",
			Type:        ActionProcess,
			Command:     "kill_process:" + proc,
			Keywords:    []string{proc, "process", "kill", "pid"},
		})
	}

	s.AllActions = base
	s.Filter()
}

// LevenshteinDistance computes edit distance for typo tolerance
func LevenshteinDistance(s1, s2 string) int {
	r1, r2 := []rune(strings.ToLower(s1)), []rune(strings.ToLower(s2))
	n, m := len(r1), len(r2)
	if n == 0 {
		return m
	}
	if m == 0 {
		return n
	}

	d := make([][]int, n+1)
	for i := range d {
		d[i] = make([]int, m+1)
		d[i][0] = i
	}
	for j := 0; j <= m; j++ {
		d[0][j] = j
	}

	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			cost := 1
			if r1[i-1] == r2[j-1] {
				cost = 0
			}
			d[i][j] = int(math.Min(
				float64(d[i-1][j]+1),
				math.Min(
					float64(d[i][j-1]+1),
					float64(d[i-1][j-1]+cost),
				),
			))
		}
	}
	return d[n][m]
}

// Filter runs fuzzy matching + typo-tolerant ranking and calculates ghost-text
func (s *SpotlightState) Filter() {
	q := strings.TrimSpace(strings.ToLower(s.Query))
	s.GhostText = ""

	if q == "" {
		s.Results = s.AllActions
		if s.Cursor >= len(s.Results) {
			s.Cursor = 0
		}
		return
	}

	type ScoredItem struct {
		item  LauncherItem
		score int
	}

	var matched []ScoredItem

	for _, it := range s.AllActions {
		score := -1
		tLower := strings.ToLower(it.Title)
		cLower := strings.ToLower(it.Command)

		// 1. Exact prefix match
		if strings.HasPrefix(tLower, q) || strings.HasPrefix(cLower, q) {
			score = 100 - len(tLower)
		} else if strings.Contains(tLower, q) || strings.Contains(cLower, q) {
			score = 60
		} else {
			// Check keywords
			for _, kw := range it.Keywords {
				if strings.HasPrefix(kw, q) {
					score = 80
					break
				} else if strings.Contains(kw, q) {
					score = 50
					break
				} else if len(q) >= 3 && LevenshteinDistance(kw, q) <= 1 {
					// Typo tolerance: edit distance <= 1
					score = 40
					break
				}
			}
		}

		if score > 0 {
			matched = append(matched, ScoredItem{item: it, score: score})
		}
	}

	// Sort matched by score descending
	for i := 0; i < len(matched)-1; i++ {
		for j := i + 1; j < len(matched); j++ {
			if matched[j].score > matched[i].score {
				matched[i], matched[j] = matched[j], matched[i]
			}
		}
	}

	s.Results = nil
	for _, m := range matched {
		s.Results = append(s.Results, m.item)
	}

	if s.Cursor >= len(s.Results) {
		s.Cursor = 0
	}

	// Calculate Ghost Text for inline suggestion from top match
	if len(s.Results) > 0 {
		topTitle := s.Results[0].Title
		// Find clean keyword or title completion
		topLower := strings.ToLower(topTitle)
		if strings.HasPrefix(topLower, q) {
			s.GhostText = topTitle[len(q):]
		} else {
			// Check if any keyword prefixes q
			for _, kw := range s.Results[0].Keywords {
				if strings.HasPrefix(kw, q) {
					s.GhostText = kw[len(q):]
					break
				}
			}
		}
	}
}

// AcceptGhostText appends the ghost text to the query when Tab or Right Arrow is pressed
func (s *SpotlightState) AcceptGhostText() {
	if s.GhostText != "" {
		s.Query += s.GhostText
		s.GhostText = ""
		s.Filter()
	}
}

// RenderSpotlight renders the centered floating 16-bit spotlight modal
func RenderSpotlight(state SpotlightState, st styles.ThemeStyles, width, height int) string {
	modalWidth := 70
	if modalWidth > width-6 {
		modalWidth = width - 6
	}

	// Header Banner
	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#000000")).
		Background(st.PrimaryColor).
		Padding(0, 1).
		Render("🔍 UNIFIED SPOTLIGHT LAUNCHER")

	// Search Input Box with Ghost-Text
	queryPrompt := lipgloss.NewStyle().Foreground(st.PrimaryColor).Bold(true).Render(" ⚡ Ketik Perintah: ")
	queryText := lipgloss.NewStyle().Foreground(st.TextColor).Bold(true).Render(state.Query)
	cursorBlock := lipgloss.NewStyle().Foreground(st.PrimaryColor).Render("▌")
	ghostTextRender := lipgloss.NewStyle().Foreground(st.MutedColor).Italic(true).Render(state.GhostText)

	if state.Query == "" && state.GhostText == "" {
		ghostTextRender = lipgloss.NewStyle().Foreground(st.MutedColor).Italic(true).Render("Cari service, ketik '>clean', '>ufw', atau nama app...")
	}

	inputRow := lipgloss.JoinHorizontal(lipgloss.Left,
		queryPrompt,
		queryText,
		cursorBlock,
		ghostTextRender,
	)

	inputBox := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(st.GlowColor).
		Width(modalWidth - 6).
		Padding(0, 1).
		Render(inputRow)

	// Results List
	var resultRows []string
	maxDisplay := 6
	startIdx := 0
	if state.Cursor >= maxDisplay {
		startIdx = state.Cursor - maxDisplay + 1
	}

	endIdx := startIdx + maxDisplay
	if endIdx > len(state.Results) {
		endIdx = len(state.Results)
	}

	if len(state.Results) == 0 {
		emptyMsg := lipgloss.NewStyle().
			Foreground(st.DangerColor).
			Italic(true).
			Render("  ❌ Tidak ditemukan perintah yang cocok. Coba: >clean, >ports, >theme, atau nama service.")
		resultRows = append(resultRows, emptyMsg)
	} else {
		for i := startIdx; i < endIdx; i++ {
			it := state.Results[i]
			isSelected := (i == state.Cursor)

			tagStyle := lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#000000")).
				Background(st.AccentColor).
				Padding(0, 1).
				Render(it.Category)

			titleText := it.Title
			descText := " - " + it.Description
			if len(descText) > modalWidth-35 {
				descText = descText[:modalWidth-38] + "..."
			}

			if isSelected {
				arrow := lipgloss.NewStyle().Foreground(st.PrimaryColor).Bold(true).Render(" ► ")
				line := lipgloss.JoinHorizontal(lipgloss.Left,
					arrow,
					tagStyle,
					" ",
					lipgloss.NewStyle().Foreground(st.PrimaryColor).Bold(true).Render(titleText),
					lipgloss.NewStyle().Foreground(st.MutedColor).Render(descText),
				)
				resultRows = append(resultRows, lipgloss.NewStyle().Background(lipgloss.Color("#1E293B")).Width(modalWidth-6).Render(line))
			} else {
				line := lipgloss.JoinHorizontal(lipgloss.Left,
					"   ",
					tagStyle,
					" ",
					lipgloss.NewStyle().Foreground(st.TextColor).Render(titleText),
					lipgloss.NewStyle().Foreground(st.MutedColor).Render(descText),
				)
				resultRows = append(resultRows, line)
			}
		}
	}

	resultsBox := lipgloss.JoinVertical(lipgloss.Left, resultRows...)

	// Footer Help Hint
	hint := lipgloss.NewStyle().
		Foreground(st.MutedColor).
		Render("[Tab / →]: Lengkapi Teks | [↑/↓]: Pilih | [Enter]: Jalankan | [Esc]: Tutup")

	modalBody := lipgloss.JoinVertical(lipgloss.Center,
		header,
		"",
		inputBox,
		"",
		resultsBox,
		"",
		hint,
	)

	return st.Modal.Width(modalWidth).Render(modalBody)
}
