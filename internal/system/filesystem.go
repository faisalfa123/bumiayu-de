package system

import (
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

// FileEntry represents a file or directory
type FileEntry struct {
	Name        string
	Path        string
	IsDir       bool
	Size        int64
	Permissions string
	Owner       string
	ModTime     string
}

// ListDirectory lists files in a directory safely
func ListDirectory(dirPath string) ([]FileEntry, error) {
	if dirPath == "" {
		if runtime.GOOS == "windows" {
			dirPath = "C:\\"
		} else {
			dirPath = "/etc"
		}
	}

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	var list []FileEntry

	// Parent directory entry if not at root
	parent := filepath.Dir(dirPath)
	if parent != dirPath {
		list = append(list, FileEntry{
			Name:        "..",
			Path:        parent,
			IsDir:       true,
			Permissions: "drwxr-xr-x",
		})
	}

	for _, e := range entries {
		info, err := e.Info()
		if err != nil {
			continue
		}

		perm := info.Mode().String()
		owner := getFileOwner(info)

		list = append(list, FileEntry{
			Name:        e.Name(),
			Path:        filepath.Join(dirPath, e.Name()),
			IsDir:       e.IsDir(),
			Size:        info.Size(),
			Permissions: perm,
			Owner:       owner,
			ModTime:     info.ModTime().Format("2006-01-02 15:04"),
		})
	}

	// Sort directories first, then files
	sort.Slice(list, func(i, j int) bool {
		if list[i].Name == ".." {
			return true
		}
		if list[j].Name == ".." {
			return false
		}
		if list[i].IsDir != list[j].IsDir {
			return list[i].IsDir
		}
		return strings.ToLower(list[i].Name) < strings.ToLower(list[j].Name)
	})

	return list, nil
}

func getFileOwner(info fs.FileInfo) string {
	// Simple portable owner representation
	if runtime.GOOS == "linux" {
		return "root:root"
	}
	return "user:staff"
}

// ReadFilePreview reads the first N lines of a text file safely
func ReadFilePreview(filePath string, maxLines int) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(data), "\n")
	if len(lines) > maxLines {
		lines = lines[:maxLines]
		lines = append(lines, "... [Truncated for TUI Preview] ...")
	}

	return strings.Join(lines, "\n"), nil
}
