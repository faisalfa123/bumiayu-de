package web

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/bumiayu/bumiayu-de/internal/system"
)

var (
	serverOnce sync.Once
	running    bool
	activePort int
)

// StartWebServer starts a lightweight embedded HTTP server in background
func StartWebServer(port int) {
	if port <= 0 {
		port = 8080
	}
	activePort = port

	serverOnce.Do(func() {
		go func() {
			mux := http.NewServeMux()
			mux.HandleFunc("/", handleHome)
			mux.HandleFunc("/api/telemetry", handleTelemetry)
			mux.HandleFunc("/api/immich/install", handleImmichInstall)

			addr := fmt.Sprintf(":%d", activePort)
			running = true
			_ = http.ListenAndServe(addr, mux)
		}()
	})
}

// GetWebStatus returns running status and web port
func GetWebStatus() (bool, int) {
	return running, activePort
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	info := system.GetHardwareInfo()
	catalog := system.GetAppCatalog()

	var immichStatus string
	if system.IsPackageInstalled("immich-sbc-tool", "immich-sbc") {
		immichStatus = `<span class="badge success">● Terpasang & Siap</span>`
	} else {
		immichStatus = `<span class="badge warning">○ Belum Terpasang</span>`
	}

	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="id">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Bumiayu DE - iPadOS Web Dashboard</title>
    <style>
        :root {
            --bg-canvas: #E5E7EB;
            --card-bg: #FFFFFF;
            --text-main: #0F172A;
            --text-muted: #64748B;
            --accent-blue: #2563EB;
            --accent-green: #16A34A;
            --accent-rose: #E11D48;
            --border-card: #E2E8F0;
        }
        body {
            margin: 0;
            padding: 24px;
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
            background: var(--bg-canvas);
            color: var(--text-main);
        }
        .header-bar {
            background: var(--card-bg);
            padding: 16px 24px;
            border-radius: 20px;
            display: flex;
            justify-content: space-between;
            align-items: center;
            box-shadow: 0 4px 6px -1px rgba(0,0,0,0.05);
            margin-bottom: 24px;
        }
        .header-title {
            font-size: 20px;
            font-weight: 700;
            display: flex;
            align-items: center;
            gap: 10px;
        }
        .top-nav {
            display: flex;
            gap: 8px;
            background: #F1F5F9;
            padding: 4px;
            border-radius: 12px;
        }
        .nav-item {
            padding: 8px 16px;
            border-radius: 8px;
            text-decoration: none;
            color: var(--text-muted);
            font-weight: 600;
            font-size: 14px;
            transition: all 0.2s;
        }
        .nav-item.active {
            background: #FFFFFF;
            color: var(--accent-blue);
            box-shadow: 0 2px 4px rgba(0,0,0,0.05);
        }
        .grid-cards {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
            gap: 20px;
            margin-bottom: 24px;
        }
        .card {
            background: var(--card-bg);
            padding: 20px;
            border-radius: 18px;
            box-shadow: 0 4px 6px -1px rgba(0,0,0,0.05);
            border: 1px solid var(--border-card);
        }
        .card-header {
            font-size: 13px;
            font-weight: 700;
            text-transform: uppercase;
            color: var(--text-muted);
            margin-bottom: 8px;
        }
        .card-value {
            font-size: 26px;
            font-weight: 800;
            color: var(--text-main);
            margin-bottom: 6px;
        }
        .badge {
            display: inline-block;
            padding: 4px 10px;
            border-radius: 9999px;
            font-size: 12px;
            font-weight: 700;
        }
        .badge.success { background: #DCFCE7; color: #15803D; }
        .badge.warning { background: #FEF3C7; color: #B45309; }
        .badge.blue { background: #DBEAFE; color: #1E40AF; }
        .btn-action {
            display: inline-block;
            background: var(--accent-blue);
            color: #FFFFFF;
            padding: 10px 20px;
            border-radius: 10px;
            text-decoration: none;
            font-weight: 700;
            border: none;
            cursor: pointer;
            transition: opacity 0.2s;
        }
        .btn-action:hover { opacity: 0.9; }
        .table-container {
            background: var(--card-bg);
            border-radius: 18px;
            padding: 20px;
            box-shadow: 0 4px 6px -1px rgba(0,0,0,0.05);
        }
        table {
            width: 100%%;
            border-collapse: collapse;
        }
        th, td {
            text-align: left;
            padding: 12px 14px;
            border-bottom: 1px solid #F1F5F9;
        }
        th { font-size: 13px; color: var(--text-muted); }
    </style>
</head>
<body>
    <div class="header-bar">
        <div class="header-title">
            <span>🏝️</span> Bumiayu DE Web Control Panel
        </div>
        <div class="top-nav">
            <a href="#" class="nav-item active">📊 Dasbor</a>
            <a href="#apps" class="nav-item">📦 Aplikasi</a>
            <a href="#immich" class="nav-item">📸 Immich C2</a>
            <a href="#system" class="nav-item">⚙️ Sistem</a>
        </div>
    </div>

    <div class="grid-cards">
        <div class="card">
            <div class="card-header">🧠 Penggunaan RAM</div>
            <div class="card-value">%d / %d MB</div>
            <div>Persentase: <b>%.0f%%</b></div>
        </div>

        <div class="card">
            <div class="card-header">🌡️ Suhu & Beban CPU</div>
            <div class="card-value">%.1f °C</div>
            <div>Load 1m: <b>%.2f</b></div>
        </div>

        <div class="card" id="immich">
            <div class="card-header">📸 Immich ODROID-C2</div>
            <div class="card-value">FaisalFa123</div>
            <div style="margin-bottom: 12px;">Status: %s</div>
            <button class="btn-action" onclick="installImmich()">⚡ Jalankan / Pasang</button>
        </div>
    </div>

    <div class="table-container" id="apps">
        <h3 style="margin-top: 0; font-size: 18px; font-weight: 700;">📦 Aplikasi Server Populer</h3>
        <table>
            <thead>
                <tr>
                    <th>Nama Aplikasi</th>
                    <th>Kategori</th>
                    <th>Port</th>
                    <th>Status</th>
                </tr>
            </thead>
            <tbody>
`, info.MemUsedMB, info.MemTotalMB, info.MemUsagePct, info.TemperatureC, info.Load1, immichStatus)

	for i, app := range catalog {
		if i >= 8 {
			break
		}
		statusBadge := `<span class="badge warning">Tersedia</span>`
		if app.Installed {
			statusBadge = `<span class="badge success">Terpasang</span>`
		}
		html += fmt.Sprintf(`
                <tr>
                    <td><b>%s</b></td>
                    <td>%s</td>
                    <td>%s</td>
                    <td>%s</td>
                </tr>`, app.Name, app.Category, app.Port, statusBadge)
	}

	html += `
            </tbody>
        </table>
    </div>

    <script>
        function installImmich() {
            if (confirm("Pasang atau perbarui Immich SBC Toolkit sekarang?")) {
                fetch("/api/immich/install", { method: "POST" })
                    .then(res => res.json())
                    .then(data => alert(data.message))
                    .catch(err => alert("Gagal menghubungi server"));
            }
        }
    </script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}

func handleTelemetry(w http.ResponseWriter, r *http.Request) {
	info := system.GetHardwareInfo()
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"ram_used_mb": %d, "ram_total_mb": %d, "temp_c": %.1f, "load_avg": %.2f}`,
		info.MemUsedMB, info.MemTotalMB, info.TemperatureC, info.Load1)
}

func handleImmichInstall(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	go func() {
		item := system.AppCatalogItem{
			ID:           "immich-sbc",
			Name:         "Immich SBC Optimizer",
			PackageName:  "immich-sbc-tool",
			IsThirdParty: true,
		}
		_, _ = system.ExecutePackageAction("install", item)
	}()
	fmt.Fprint(w, `{"status": "ok", "message": "Proses instalasi Immich SBC sedang berjalan di latar belakang server!"}`)
}
