package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/bumiayu/bumiayu-de/internal/app"
	"github.com/bumiayu/bumiayu-de/internal/i18n"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	Version   = "1.0.0"
	BuildArch = "arm64/amd64"
)

func main() {
	showVersion := flag.Bool("v", false, "Tampilkan versi Bumiayu TUI")
	showVersionLong := flag.Bool("version", false, "Tampilkan versi Bumiayu TUI")
	langFlag := flag.String("lang", "id", "Pilihan bahasa awal (id atau en)")
	flag.Parse()

	if *showVersion || *showVersionLong {
		fmt.Printf("Bumiayu TUI v%s (%s)\nLightweight Desktop Environment for ODROID & Linux Servers\n", Version, BuildArch)
		os.Exit(0)
	}

	initialModel := app.InitialModel()
	if *langFlag == "en" {
		initialModel.Language = i18n.EN
	}

	p := tea.NewProgram(
		initialModel,
		tea.WithAltScreen(),       // Menggunakan alternate screen buffer agar terminal tetap bersih
		tea.WithMouseCellMotion(),  // Dukungan klik dan scroll mouse di terminal
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Kesalahan saat menjalankan Bumiayu TUI: %v\n", err)
		os.Exit(1)
	}
}
