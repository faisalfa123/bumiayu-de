package system

import (
	"bufio"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

// UserAccount represents a system user and their privileges
type UserAccount struct {
	Username       string
	UID            int
	GID            int
	HomeDir        string
	Shell          string
	IsSudoer       bool
	HasSSHKey      bool
	SSHKeyCount    int
	AccountLocked  bool
}

// GetUserAccounts lists users on the system (root and normal users UID >= 1000)
func GetUserAccounts() []UserAccount {
	if runtime.GOOS != "linux" {
		return []UserAccount{
			{Username: "root", UID: 0, GID: 0, HomeDir: "/root", Shell: "/bin/bash", IsSudoer: true, HasSSHKey: true, SSHKeyCount: 2, AccountLocked: false},
			{Username: "odroid", UID: 1000, GID: 1000, HomeDir: "/home/odroid", Shell: "/bin/bash", IsSudoer: true, HasSSHKey: true, SSHKeyCount: 1, AccountLocked: false},
			{Username: "operator", UID: 1001, GID: 1001, HomeDir: "/home/operator", Shell: "/bin/bash", IsSudoer: false, HasSSHKey: false, SSHKeyCount: 0, AccountLocked: false},
		}
	}

	sudoMembers := getSudoMembers()
	var users []UserAccount

	f, err := os.Open("/etc/passwd")
	if err != nil {
		return nil
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) >= 7 {
			username := parts[0]
			uid, _ := strconv.Atoi(parts[2])
			gid, _ := strconv.Atoi(parts[3])
			home := parts[5]
			shell := parts[6]

			// Only show root and human users (UID >= 1000), or users with login shells
			if uid == 0 || (uid >= 1000 && uid < 65534) {
				keyCount := countSSHKeys(home)
				users = append(users, UserAccount{
					Username:    username,
					UID:         uid,
					GID:         gid,
					HomeDir:     home,
					Shell:       shell,
					IsSudoer:    sudoMembers[username] || uid == 0,
					HasSSHKey:   keyCount > 0,
					SSHKeyCount: keyCount,
				})
			}
		}
	}

	return users
}

func getSudoMembers() map[string]bool {
	members := make(map[string]bool)
	f, err := os.Open("/etc/group")
	if err != nil {
		return members
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) >= 4 {
			groupName := parts[0]
			if groupName == "sudo" || groupName == "wheel" || groupName == "admin" {
				users := strings.Split(parts[3], ",")
				for _, u := range users {
					u = strings.TrimSpace(u)
					if u != "" {
						members[u] = true
					}
				}
			}
		}
	}
	return members
}

func countSSHKeys(homeDir string) int {
	authKeysPath := filepath.Join(homeDir, ".ssh", "authorized_keys")
	data, err := os.ReadFile(authKeysPath)
	if err != nil {
		return 0
	}
	count := 0
	lines := strings.Split(string(data), "\n")
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l != "" && !strings.HasPrefix(l, "#") {
			count++
		}
	}
	return count
}

// ModifyUserSudoRole grants or revokes sudo role for a user
func ModifyUserSudoRole(username string, grant bool) (string, error) {
	if runtime.GOOS != "linux" {
		return "Simulated: sudo change", nil
	}

	if grant {
		cmd := exec.Command("usermod", "-aG", "sudo", username)
		out, err := cmd.CombinedOutput()
		return string(out), err
	}
	
	cmd := exec.Command("gpasswd", "-d", username, "sudo")
	out, err := cmd.CombinedOutput()
	return string(out), err
}
