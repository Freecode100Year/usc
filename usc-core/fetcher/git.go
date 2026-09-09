package fetcher

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// FetchGit clones a remote git repository into destDir.
func FetchGit(gitURL, destDir string) (*SkillMetadata, error) {
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return nil, err
	}
	repoURL := normalizeGitURL(gitURL)
	cmd := exec.Command("git", "clone", "--depth", "1", repoURL, destDir)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("git clone failed: %v, output: %s", err, string(out))
	}
	skillName := filepath.Base(destDir)
	return &SkillMetadata{
		Name:        skillName,
		Version:     "0.1.0",
		Description: fmt.Sprintf("Cloned git skill from %s", gitURL),
		SourceURL:   gitURL,
		LocalPath:   destDir,
	}, nil
}

func normalizeGitURL(in string) string {
	s := strings.TrimSpace(in)
	if !strings.HasPrefix(s, "http://") && !strings.HasPrefix(s, "https://") && !strings.HasPrefix(s, "git@") {
		return "https://" + s
	}
	return s
}
