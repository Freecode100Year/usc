package adapter

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Freecode100Year/usc/usc-core/attestation"
	"github.com/Freecode100Year/usc/usc-core/capability"
)

// ClaudeCodeAdapter adapts attested USC skills to Anthropic Claude Code.
type ClaudeCodeAdapter struct{}

func NewClaudeCodeAdapter() *ClaudeCodeAdapter {
	return &ClaudeCodeAdapter{}
}

func (a *ClaudeCodeAdapter) Target() string {
	return TargetClaudeCode
}

func (a *ClaudeCodeAdapter) DisplayName() string {
	return "Claude Code CLI"
}

func (a *ClaudeCodeAdapter) DetectInstalled() (string, bool) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", false
	}
	p := filepath.Join(home, ".claude", "skills")
	if _, err := os.Stat(filepath.Dir(p)); err == nil {
		return p, true
	}
	return filepath.Join(home, ".claude", "skills"), false
}

func (a *ClaudeCodeAdapter) GenerateBundle(att *attestation.Attestation, bp capability.Set, outDir string) (string, error) {
	targetDir := filepath.Join(outDir, "claudecode-"+att.Artifact.Name)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return "", err
	}
	toolSchema := fmt.Sprintf(`{
  "name": "%s",
  "description": "Zero-trust verified agent capability for %s",
  "input_schema": { "type": "object", "properties": { "query": { "type": "string" } } }
}`, att.Artifact.Name, att.Artifact.Name)
	if err := os.WriteFile(filepath.Join(targetDir, "tool.json"), []byte(toolSchema), 0644); err != nil {
		return "", err
	}
	handler := "#!/bin/sh\n# Claude Code Bridge to USC\nusc run " + att.Artifact.Name + "\n"
	if err := os.WriteFile(filepath.Join(targetDir, "execute.sh"), []byte(handler), 0755); err != nil {
		return "", err
	}
	return targetDir, nil
}

func (a *ClaudeCodeAdapter) Install(bundleDir string, destDir string) (string, error) {
	if destDir == "" {
		dest, _ := a.DetectInstalled()
		destDir = dest
	}
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return "", err
	}
	finalPath := filepath.Join(destDir, filepath.Base(bundleDir))
	return finalPath, copyDirectory(bundleDir, finalPath)
}
