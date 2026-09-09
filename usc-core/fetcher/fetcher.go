package fetcher

import (
	"fmt"
	"os"
	"path/filepath"
)

// Fetch fetches a skill source (remote URL, slug, or local directory) into destDir.
func Fetch(source, destDir string) (*SkillMetadata, error) {
	st := DetectSourceType(source)
	switch st {
	case SourceTypeClawHub:
		owner, skill, ok := ParseClawHubSlug(source)
		if !ok {
			return nil, fmt.Errorf("invalid clawhub format: %s", source)
		}
		rawURL := source
		if !IsRemoteSource(source) {
			rawURL = fmt.Sprintf("https://clawhub.ai/%s/skills/%s", owner, skill)
		}
		return FetchClawHub(owner, skill, rawURL, destDir)
	case SourceTypeGitHub:
		return FetchGit(source, destDir)
	case SourceTypeLocalDir:
		return inspectLocalDir(source)
	default:
		return nil, fmt.Errorf("unsupported source type: %s", source)
	}
}

func inspectLocalDir(srcDir string) (*SkillMetadata, error) {
	abs, err := filepath.Abs(srcDir)
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(abs); err != nil {
		return nil, err
	}
	name := filepath.Base(abs)
	return &SkillMetadata{
		Name:        name,
		Version:     "0.1.0",
		Description: fmt.Sprintf("Local skill directory at %s", abs),
		SourceURL:   abs,
		LocalPath:   abs,
	}, nil
}
