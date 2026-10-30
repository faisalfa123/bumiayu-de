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

// GetAnimatedPaddyFarmer returns rich animated pixel art of farmers in Bumiayu rice paddy
// Smooth frame interpolation: transitions every 12 ticks (~500ms @ 24fps) to prevent jittery flickering
func GetAnimatedPaddyFarmer(frame int) string {
	f := (frame / 12) % 6
	switch f {
	case 0:
		// Frame 0: Petani 1 bersiap, Kerbau membajak sawah di sisi kanan
		return "\x1b[1;32m  🌾 SAWAH TERASERING BUMIAYU [PIXEL ART ANIMASI HIDUP] 🌾\x1b[0m\n" +
			"\x1b[38;2;34;197;94m  " + strings.Repeat("━", 72) + "\x1b[0m\n" +
			"\x1b[38;2;254;240;138m    (o_  \x1b[38;2;251;191;36m__/\x1b[0m                       \x1b[38;2;161;98;7m(___)  \x1b[38;2;180;83;9m(___)\x1b[0m                  \x1b[38;2;253;224;71m☁️   Awan Berarak\x1b[0m\n" +
			"\x1b[38;2;251;146;60m    /|\\ \x1b[38;2;180;83;9m[\x1b[38;2;202;138;4m===|\x1b[0m                      \x1b[38;2;180;83;9m(o o)____/     \x1b[38;2;217;119;6m(o_ \x1b[38;2;251;191;36m__/\x1b[0m        \x1b[38;2;250;204;21m☀️  Matahari Pagi\x1b[0m\n" +
			"\x1b[38;2;217;119;6m    / \\                          \x1b[38;2;146;64;14m@@             \x1b[38;2;217;119;6m/|\\ \x1b[38;2;180;83;9m[\x1b[38;2;202;138;4m===|\x1b[0m\n" +
			"\x1b[38;2;134;239;172m  ,,,,, ,,,,, ,,,,, ,,,,, ,,,,, ,,,,, ,,,,, ,,,,, ,,,,, ,,,,, ,,,,, ,,,,, \x1b[0m\n" +
			"\x1b[38;2;74;222;128m _|||||_|||||_|||||_|||||_|||||_|||||_|||||_|||||_|||||_|||||_|||||_|||||_\x1b[0m\n" +
			"\x1b[38;2;22;163;74m ████████████████████████████████████████████████████████████████████████\x1b[0m\n" +
			"\x1b[1;33m ░░ 🌾 Petani Mengayun Cangkul & Kerbau Membajak Sawah Subur Bumiayu 🌾 ░░\x1b[0m"

	case 1:
		// Frame 1: Cangkul menancap ke tanah lumpur, burung terbang di angkasa
		return "\x1b[1;32m  🌾 SAWAH TERASERING BUMIAYU [PIXEL ART ANIMASI HIDUP] 🌾\x1b[0m\n" +
			"\x1b[38;2;34;197;94m  " + strings.Repeat("━", 72) + "\x1b[0m\n" +
			"\x1b[38;2;254;240;138m    (o_            \x1b[38;2;224;242;254m~v~                   \x1b[38;2;161;98;7m(___)  \x1b[38;2;180;83;9m(___)\x1b[0m            \x1b[38;2;253;224;71m☁️  Awan Berarak\x1b[0m\n" +
			"\x1b[38;2;251;146;60m    /|\\_                                \x1b[38;2;180;83;9m(o o)____/     \x1b[38;2;217;119;6m(o_\x1b[0m           \x1b[38;2;250;204;21m☀️  Matahari Pagi\x1b[0m\n" +
			"\x1b[38;2;217;119;6m    / \\ \x1b[38;2;180;83;9m\\                                \x1b[38;2;146;64;14m@@             \x1b[38;2;217;119;6m/|\\_\x1b[0m\n" +
			"\x1b[38;2;134;239;172m  ,,,,, ,\x1b[38;2;202;138;4m[===|\x1b[38;2;134;239;172m, ,,,,, ,,,,, ,,,,, ,,,,, ,,,,, ,,,,, ,,,,, ,\x1b[38;2;202;138;4m[===|\x1b[38;2;134;239;172m, ,,,,, ,,,,, \x1b[0m\n" +
			"\x1b[38;2;74;222;128m _|||||_|||||_|||||_|||||_|||||_|||||_|||||_|||||_|||||_|||||_|||||_|||||_\x1b[0m\n" +
			"\x1b[38;2;22;163;74m ████████████████████████████████████████████████████████████████████████\x1b[0m\n" +
			"\x1b[1;33m ░░ 🌾 Cangkul Menancap Kokoh — Burung Pipit Melintas di Pematang 🌾 ░░\x1b[0m"

	case 2:
		// Frame 2: Petani melangkah ke depan
		return "\x1b[1;32m  🌾 SAWAH TERASERING BUMIAYU [PIXEL ART ANIMASI HIDUP] 🌾\x1b[0m\n" +
			"\x1b[38;2;34;197;94m  " + strings.Repeat("━", 72) + "\x1b[0m\n" +
			"\x1b[38;2;254;240;138m      (o_                 \x1b[38;2;224;242;254m~v~                   \x1b[38;2;161;98;7m(___)  \x1b[38;2;180;83;9m(___)\x1b[0m        \x1b[38;2;253;224;71m☁️ Awan Berarak\x1b[0m\n" +
			"\x1b[38;2;251;146;60m     --|--\x1b[38;2;180;83;9m/                                    \x1b[38;2;180;83;9m(o o)____/    \x1b[38;2;217;119;6m--|--\x1b[38;2;180;83;9m/\x1b[0m     \x1b[38;2;250;204;21m☀️  Matahari Pagi\x1b[0m\n" +
			"\x1b[38;2;217;119;6m      / > \x1b[38;2;202;138;4m[===|                                \x1b[38;2;146;64;14m@@            \x1b[38;2;217;119;6m/ > \x1b[38;2;202;138;4m[===|\x1b[0m\n" +
			"\x1b[38;2;134;239;172m  ,,,,, ,,,,, ,,,,, ,,,,, ,,,,, ,,,,, ,,,,, ,,,,, ,,,,, ,,,,, ,,,,, ,,,,, \x1b[0m\n" +
			"\x1b[38;2;74;222;128m _|||||_|||||_|||||_|||||_|||||_|||||_|||||_|||||_|||||_|||||_|||||_|||||_\x1b[0m\n" +
			"\x1b[38;2;22;163;74m ████████████████████████████████████████████████████████████████████████\x1b[0m\n" +
			"\x1b[1;33m ░░ 🌾 Melangkah Menyusuri Pematang Menuju Lereng Gunung Slamet 🌾 ░░\x1b[0m"

	case 3:
		// Frame 3: Mengayunkan cangkul kembali ke atas
		return "\x1b[1;32m  🌾 SAWAH TERASERING BUMIAYU [PIXEL ART ANIMASI HIDUP] 🌾\x1b[0m\n" +
			"\x1b[38;2;34;197;94m  " + strings.Repeat("━", 72) + "\x1b[0m\n" +
			"\x1b[38;2;254;240;138m        (o_  \x1b[38;2;251;191;36m__/\x1b[0m                               \x1b[38;2;161;98;7m(___)  \x1b[38;2;180;83;9m(___)\x1b[0m      \x1b[38;2;253;224;71m☁️  Awan Berarak\x1b[0m\n" +
			"\x1b[38;2;251;146;60m        /|\\ \x1b[38;2;180;83;9m[\x1b[38;2;202;138;4m===|\x1b[0m                            \x1b[38;2;180;83;9m(o o)____/   \x1b[38;2;217;119;6m(o_  \x1b[38;2;251;191;36m__/\x1b[0m    \x1b[38;2;250;204;21m☀️  Matahari Pagi\x1b[0m\n" +
			"\x1b[38;2;217;119;6m        / \\                                \x1b[38;2;146;64;14m@@           \x1b[38;2;217;119;6m/|\\ \x1b[38;2;180;83;9m[\x1b[38;2;202;138;4m===|\x1b[0m\n" +
			"\x1b[38;2;134;239;172m  ,,,,, ,,,,, ,,,,, ,,,,, ,,,,, ,,,,, ,,,,, ,,,,, ,,,,, ,,,,, ,,,,, ,,,,, \x1b[0m\n" +
			"\x1b[38;2;74;222;128m _|||||_|||||_|||||_|||||_|||||_|||||_|||||_|||||_|||||_|||||_|||||_|||||_\x1b[0m\n" +
			"\x1b[38;2;22;163;74m ████████████████████████████████████████████████████████████████████████\x1b[0m\n" +
			"\x1b[1;33m ░░ 🌾 Angin Semilir Lembah Slamet Membawa Kesejukan Bagi Server 🌾 ░░\x1b[0m"

	case 4:
		// Frame 4: Mencangkul dalam tanah
		return "\x1b[1;32m  🌾 SAWAH TERASERING BUMIAYU [PIXEL ART ANIMASI HIDUP] 🌾\x1b[0m\n" +
			"\x1b[38;2;34;197;94m  " + strings.Repeat("━", 72) + "\x1b[0m\n" +
			"\x1b[38;2;254;240;138m        (o_                                     \x1b[38;2;161;98;7m(___)  \x1b[38;2;180;83;9m(___)\x1b[0m    \x1b[38;2;253;224;71m☁️   Awan Berarak\x1b[0m\n" +
			"\x1b[38;2;251;146;60m        /|\\_                                   \x1b[38;2;180;83;9m(o o)____/   \x1b[38;2;217;119;6m(o_\x1b[0m         \x1b[38;2;250;204;21m☀️  Matahari Pagi\x1b[0m\n" +
			"\x1b[38;2;217;119;6m        / \\ \x1b[38;2;180;83;9m\\                                   \x1b[38;2;146;64;14m@@           \x1b[38;2;217;119;6m/|\\_\x1b[0m\n" +
			"\x1b[38;2;134;239;172m  ,,,,, ,,,,, ,,\x1b[38;2;202;138;4m[===|\x1b[38;2;134;239;172m,,, ,,,,, ,,,,, ,,,,, ,,,,, ,,,,, ,,,,, ,,\x1b[38;2;202;138;4m[===|\x1b[38;2;134;239;172m,, \x1b[0m\n" +
			"\x1b[38;2;74;222;128m _|||||_|||||_|||||_|||||_|||||_|||||_|||||_|||||_|||||_|||||_|||||_|||||_\x1b[0m\n" +
			"\x1b[38;2;22;163;74m ████████████████████████████████████████████████████████████████████████\x1b[0m\n" +
			"\x1b[1;33m ░░ 🌾 Kesuburan Tanah Menghasilkan Pangan Terbaik untuk Nusantara 🌾 ░░\x1b[0m"

	default:
		// Frame 5: Istirahat sejenak sambil minum teh poci hangat
		return "\x1b[1;32m  🌾 SAWAH TERASERING BUMIAYU [PIXEL ART ANIMASI HIDUP] 🌾\x1b[0m\n" +
			"\x1b[38;2;34;197;94m  " + strings.Repeat("━", 72) + "\x1b[0m\n" +
			"\x1b[38;2;254;240;138m          (o_                                   \x1b[38;2;161;98;7m(___)  \x1b[38;2;180;83;9m(___)\x1b[0m    \x1b[38;2;253;224;71m☁️  Awan Berarak\x1b[0m\n" +
			"\x1b[38;2;251;146;60m         <|\\  \x1b[38;2;254;240;138m☕\x1b[0m                                \x1b[38;2;180;83;9m(o o)____/   \x1b[38;2;217;119;6m<|\\\x1b[0m         \x1b[38;2;250;204;21m☀️  Matahari Pagi\x1b[0m\n" +
			"\x1b[38;2;217;119;6m         < \\  \x1b[38;2;180;83;9m[\x1b[38;2;202;138;4m===|\x1b[0m                             \x1b[38;2;146;64;14m@@           \x1b[38;2;217;119;6m< \\ \x1b[38;2;180;83;9m[\x1b[38;2;202;138;4m===|\x1b[0m\n" +
			"\x1b[38;2;134;239;172m  ,,,,, ,,,,, ,,,,, ,,,,, ,,,,, ,,,,, ,,,,, ,,,,, ,,,,, ,,,,, ,,,,, ,,,,, \x1b[0m\n" +
			"\x1b[38;2;74;222;128m _|||||_|||||_|||||_|||||_|||||_|||||_|||||_|||||_|||||_|||||_|||||_|||||_\x1b[0m\n" +
			"\x1b[38;2;22;163;74m ████████████████████████████████████████████████████████████████████████\x1b[0m\n" +
			"\x1b[1;33m ░░ 🌾 Menikmati Teh Poci Gula Batu di Saung Pinggir Sawah Bumiayu 🌾 ░░\x1b[0m"
	}
}

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
		return "\x1b[1;32m  🌾 HAMBARAN SAWAH TERASERING BUMIAYU 🌾\x1b[0m\n" +
			"\x1b[38;2;34;197;94m  " + strings.Repeat("━", 60) + "\x1b[0m\n" +
			"\x1b[38;2;134;239;172m   ,,,,, ,,,,, ,,,,, ,,,,, ,,,,, ,,,,, ,,,,, ,,,,, ,,,,, ,,,,, \x1b[0m\n" +
			"\x1b[38;2;74;222;128m  _|||||_|||||_|||||_|||||_|||||_|||||_|||||_|||||_|||||_|||||_\x1b[0m\n" +
			"\x1b[38;2;34;197;94m  ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░\x1b[0m\n" +
			"\x1b[38;2;22;163;74m  ████████████████████████████████████████████████████████████\x1b[0m\n" +
			"\x1b[1;33m  ░░ 🌾 Damai, Subur, dan Sejuk Menyejukkan Server Anda 🌾 ░░\x1b[0m"

	case SceneLenteraSenja:
		return "\x1b[1;33m  🏮 LENTERA SENJA BUMIAYU 🏮\x1b[0m\n" +
			"\x1b[38;2;249;115;22m  " + strings.Repeat("━", 60) + "\x1b[0m\n" +
			"\x1b[38;2;251;146;60m             .---.            \x1b[38;2;254;240;138m✨ Menerangi pekatnya konsol terminal\x1b[0m\n" +
			"\x1b[38;2;249;115;22m            /     \\           \x1b[38;2;253;224;71m🔥 Menuntun pengguna selayaknya ahli\x1b[0m\n" +
			"\x1b[38;2;234;88;12m           | () () |          \x1b[38;2;250;204;21m🏮 Kehangatan suasana kota Bumiayu\x1b[0m\n" +
			"\x1b[38;2;194;65;12m            \\  _  /\n" +
			"             '---'\n" +
			"            [====]\n" +
			"\x1b[38;2;251;191;36m  ░░░░░░░░░░(____)░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░\x1b[0m\n" +
			"\x1b[1;31m  ░░ 🏮 Hangat Menyala, Memandu Perjalanan Linux Anda 🏮 ░░\x1b[0m"

	case SceneMalamBintang:
		return "\x1b[1;36m  🌌 LANGIT MALAM BERBINTANG BUMIAYU 🌌\x1b[0m\n" +
			"\x1b[38;2;56;189;248m  " + strings.Repeat("━", 60) + "\x1b[0m\n" +
			"\x1b[38;2;224;242;254m    ★     .            .      ✦       .     ★          .   ★\x1b[0m\n" +
			"\x1b[38;2;186;230;253m        .     ★               ★                .      ✦    \x1b[0m\n" +
			"\x1b[38;2;125;211;252m    ✦             .      ✦           .     ★          .    \x1b[0m\n" +
			"\x1b[38;2;14;165;233m  ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~\x1b[0m\n" +
			"\x1b[38;2;2;132;199m  ████████████████████████████████████████████████████████████\x1b[0m\n" +
			"\x1b[1;36m  ░░ 🌌 Hening dan Tangguh Menjaga Server Latar Belakang 🌌 ░░\x1b[0m"

	default: // SceneGunungSlamet
		return "\x1b[1;34m  🏔️ GUNUNG SLAMET MEGAH BUMIAYU 🏔️\x1b[0m\n" +
			"\x1b[38;2;99;102;241m  " + strings.Repeat("━", 60) + "\x1b[0m\n" +
			"\x1b[38;2;199;210;254m                /\\              \x1b[0m\n" +
			"\x1b[38;2;165;180;252m               /  \\       /\\    \x1b[0m\n" +
			"\x1b[38;2;129;140;248m              /    \\     /  \\   \x1b[0m\n" +
			"\x1b[38;2;99;102;241m             /      \\   /    \\       /\\\x1b[0m\n" +
			"\x1b[38;2;79;70;229m            /   /\\   \\ /      \\     /  \\\x1b[0m\n" +
			"\x1b[38;2;67;56;202m           /   /  \\   \\        \\   /    \\\x1b[0m\n" +
			"\x1b[38;2;55;48;163m  ________/___/____\\___\\________\\_/______\\_____________________\x1b[0m\n" +
			"\x1b[1;34m  ░░ 🏔️ Simbol Ketangguhan & Kedaulatan Server ODROID-C2 🏔️ ░░\x1b[0m"
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
