package adapter

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

// CleanSkillName strips adapter prefixes and extensions to yield a canonical skill name.
func CleanSkillName(rawName string) string {
	s := strings.TrimSpace(rawName)
	s = strings.TrimSuffix(s, ".usc")
	for _, prefix := range []string{"agy-", "openclaw-", "hermes-", "claudecode-", "codex-"} {
		s = strings.TrimPrefix(s, prefix)
	}
	return s
}

func copyDirectory(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if info.IsDir() {
			return os.MkdirAll(target, 0755)
		}
		return copyFile(path, target)
	})
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}
