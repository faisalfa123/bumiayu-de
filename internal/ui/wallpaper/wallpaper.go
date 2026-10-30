package wallpaper

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bumiayu/bumiayu-de/internal/config"
)

type WallpaperType int

const (
	WallpaperPreset WallpaperType = iota // Preset Seni Bumiayu
	WallpaperDynamicTime                 // Berganti Otomatis Sesuai Waktu (Pagi/Siang/Sore/Malam)
	WallpaperCustom                      // Wallpaper Kustom Pengguna (File Teks/Hasil Generate)
)

type PresetScene int

const (
	SceneGunungSlamet PresetScene = iota // 🏔️ Gunung Slamet Megah
	SceneSawahHijau                      // 🌾 Sawah Terasering
	SceneLenteraSenja                    // 🏮 Lentera Senja Bumiayu
	SceneMalamBintang                    // 🌌 Langit Malam Berbintang
)

// GetCurrentWallpaper returns the active wallpaper string based on type and time
func GetCurrentWallpaper(wpType WallpaperType, preset PresetScene) string {
	switch wpType {
	case WallpaperDynamicTime:
		return getDynamicTimeWallpaper()
	case WallpaperCustom:
		customFile := filepath.Join(config.GetConfigDir(), "custom_wallpaper.txt")
		if data, err := os.ReadFile(customFile); err == nil && len(data) > 0 {
			return string(data)
		}
		// Fallback to preset if custom not generated yet
		return getPresetArtwork(preset)
	default:
		return getPresetArtwork(preset)
	}
}

func getDynamicTimeWallpaper() string {
	hour := time.Now().Hour()
	if hour >= 5 && hour < 11 {
		return getPresetArtwork(SceneGunungSlamet) // Pagi: Fajar Slamet
	} else if hour >= 11 && hour < 15 {
		return getPresetArtwork(SceneSawahHijau) // Siang: Sawah Hijau
	} else if hour >= 15 && hour < 18 {
		return getPresetArtwork(SceneLenteraSenja) // Sore: Lentera Senja
	}
	return getPresetArtwork(SceneMalamBintang) // Malam: Langit Malam
}

func getPresetArtwork(scene PresetScene) string {
	switch scene {
	case SceneSawahHijau:
		return `  🌾 HAMBARAN SAWAH TERASERING BUMIAYU 🌾
  ─────────────────────────────────────────────────────────────
  ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░
  ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░
  ▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒
  █████████████████████████████████████████████████████████████
  ░░ 🌾 Damai, Subur, dan Sejuk Menyejukkan Server Anda 🌾 ░░`

	case SceneLenteraSenja:
		return `  🏮 LENTERA SENJA BUMIAYU 🏮
  ─────────────────────────────────────────────────────────────
             .---.
            /     \       * Menerangi pekatnya konsol terminal
           | () () |      * Menuntun pengguna selayaknya ahli
            \  _  /
             '---'
            [====]
  ░░░░░░░░░░(____)░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░
  ░░ 🏮 Hangat Menyala, Memandu Perjalanan Linux Anda 🏮 ░░`

	case SceneMalamBintang:
		return `  🌌 LANGIT MALAM BERBINTANG BUMIAYU 🌌
  ─────────────────────────────────────────────────────────────
    *     .            .              .     *          .   *
              *               *                .           
        .             .                  *          .     .
  ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
  █████████████████████████████████████████████████████████████
  ░░ 🌌 Hening dan Tangguh Menjaga Server Latar Belakang 🌌 ░░`

	default: // SceneGunungSlamet
		return `  🏔️ GUNUNG SLAMET MEGAH BUMIAYU 🏔️
  ─────────────────────────────────────────────────────────────
                /\
               /  \       /\
              /    \     /  \
             /      \   /    \       /\
            /   /\   \ /      \     /  \
           /   /  \   \        \   /    \
  ________/___/____\___\________\_/______\_____________________
  ░░ 🏔️ Simbol Ketangguhan & Kedaulatan Server ODROID-C2 🏔️ ░░`
	}
}

// ConvertImageToPixelArt reads a local image file and converts it into ANSI TrueColor pixel art
func ConvertImageToPixelArt(imagePath string, targetWidth, targetHeight int) (string, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return "", fmt.Errorf("gagal membuka berkas gambar: %w", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return "", fmt.Errorf("format gambar tidak didukung (gunakan PNG/JPG): %w", err)
	}

	bounds := img.Bounds()
	srcWidth := bounds.Dx()
	srcHeight := bounds.Dy()

	if targetWidth <= 0 {
		targetWidth = 56
	}
	if targetHeight <= 0 {
		targetHeight = 12
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("  🖼️ WALLPAPER PIXEL ART KUSTOM: %s\n", filepath.Base(imagePath)))
	sb.WriteString("  " + strings.Repeat("─", targetWidth) + "\n")

	// Sample pixels using double-block character ▀ (upper pixel fg, lower pixel bg)
	for y := 0; y < targetHeight*2; y += 2 {
		sb.WriteString("  ")
		for x := 0; x < targetWidth; x++ {
			srcX := bounds.Min.X + (x * srcWidth / targetWidth)
			srcY1 := bounds.Min.Y + (y * srcHeight / (targetHeight * 2))
			srcY2 := bounds.Min.Y + ((y + 1) * srcHeight / (targetHeight * 2))

			r1, g1, b1, _ := img.At(srcX, srcY1).RGBA()
			r2, g2, b2, _ := img.At(srcX, srcY2).RGBA()

			// Scale from 16-bit to 8-bit color
			cr1, cg1, cb1 := uint8(r1>>8), uint8(g1>>8), uint8(b1>>8)
			cr2, cg2, cb2 := uint8(r2>>8), uint8(g2>>8), uint8(b2>>8)

			// ANSI 24-bit TrueColor: upper pixel foreground, lower pixel background
			sb.WriteString(fmt.Sprintf("\x1b[38;2;%d;%d;%dm\x1b[48;2;%d;%d;%dm▀\x1b[0m", cr1, cg1, cb1, cr2, cg2, cb2))
		}
		sb.WriteString("\n")
	}

	result := sb.String()

	// Save to custom wallpaper file
	customFile := filepath.Join(config.GetConfigDir(), "custom_wallpaper.txt")
	_ = os.MkdirAll(config.GetConfigDir(), 0755)
	_ = os.WriteFile(customFile, []byte(result), 0644)

	return result, nil
}
