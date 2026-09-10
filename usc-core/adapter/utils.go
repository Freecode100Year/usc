package adapter

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// CheckBinaryOrConfig checks if any candidate binary exists in PATH or candidate config exists.
func CheckBinaryOrConfig(binNames []string, configPaths []string) bool {
	for _, bin := range binNames {
		if _, err := exec.LookPath(bin); err == nil {
			return true
		}
	}
	for _, p := range configPaths {
		if _, err := os.Stat(p); err == nil {
			return true
		}
	}
	return false
}

// ValidateSkillName ensures the skill name contains only safe alphanumeric, hyphen, and underscore characters.
func ValidateSkillName(name string) error {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return errors.New("USC-E7060: skill name cannot be empty")
	}
	for _, r := range trimmed {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_') {
			return fmt.Errorf("USC-E7061: unsafe character '%c' in skill name '%s': only [a-zA-Z0-9_-] allowed", r, name)
		}
	}
	return nil
}

// CleanSkillName strips adapter prefixes and extensions to yield a canonical skill name.
func CleanSkillName(rawName string) string {
	s := strings.TrimSpace(rawName)
	s = strings.TrimSuffix(s, ".usc")
	for _, prefix := range []string{"agy-", "openclaw-", "hermes-", "claudecode-", "codex-"} {
		s = strings.TrimPrefix(s, prefix)
	}
	return s
}

// SafeCopyFile copies a file while refusing symlinks to avoid symlink traversal attacks.
func SafeCopyFile(src, dst string) error {
	srcInfo, err := os.Lstat(src)
	if err != nil {
		return err
	}
	if srcInfo.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("USC-E7062: refusing to copy symlink %s", src)
	}
	if dstInfo, err := os.Lstat(dst); err == nil {
		if dstInfo.Mode()&os.ModeSymlink != 0 {
			_ = os.Remove(dst) // Defend against pre-planted symlink overwrite
		}
	}
	return executeFileCopy(src, dst, srcInfo.Mode())
}

func executeFileCopy(src, dst string, mode os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}

// SafeCopyDirectory recursively copies directory entries safely.
func SafeCopyDirectory(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		lInfo, err := os.Lstat(path)
		if err != nil || lInfo.Mode()&os.ModeSymlink != 0 {
			return nil // Skip symlinks
		}
		if info.IsDir() {
			return os.MkdirAll(target, 0755)
		}
		return SafeCopyFile(path, target)
	})
}

// SafeInstallDirectory safely copies bundle into destDir with backup protection against overwrite.
func SafeInstallDirectory(bundleDir, destDir, skillName string) (string, error) {
	if err := ValidateSkillName(skillName); err != nil {
		return "", err
	}
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return "", err
	}
	finalPath := filepath.Join(destDir, skillName)
	if _, err := os.Stat(finalPath); err == nil {
		backupPath := fmt.Sprintf("%s.bak.%d", finalPath, time.Now().Unix())
		_ = os.Rename(finalPath, backupPath)
	}
	if err := SafeCopyDirectory(bundleDir, finalPath); err != nil {
		return "", err
	}
	return finalPath, nil
}
