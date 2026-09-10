package adapter

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Freecode100Year/usc/usc-core/attestation"
	"github.com/Freecode100Year/usc/usc-core/capability"
)

// ClaudeCodeAdapter adapts attested USC skills to Claude Code CLI.
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
	skillPath := filepath.Join(home, ".claude", "skills")
	bins := []string{"claude", "claude.cmd", "claude.exe"}
	cfgs := []string{
		filepath.Join(home, ".claude.json"),
		filepath.Join(home, ".claude", "config.json"),
		filepath.Join(home, ".claude", "settings.json"),
	}
	detected := CheckBinaryOrConfig(bins, cfgs)
	return skillPath, detected
}

func (a *ClaudeCodeAdapter) GenerateBundle(att *attestation.Attestation, bp capability.Set, outDir string) (string, error) {
	skillName := CleanSkillName(att.Artifact.Name)
	if err := ValidateSkillName(skillName); err != nil {
		return "", err
	}
	targetDir := filepath.Join(outDir, "claudecode-"+skillName)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return "", err
	}
	if err := a.writeToolJSON(targetDir, skillName); err != nil {
		return "", err
	}
	return targetDir, a.writeSafeExecuteScript(targetDir, skillName)
}

func (a *ClaudeCodeAdapter) writeToolJSON(targetDir, skillName string) error {
	toolSchema := fmt.Sprintf(`{
  "name": %q,
  "description": "Zero-trust verified agent capability for %s",
  "input_schema": { "type": "object", "properties": { "query": { "type": "string" } } }
}`, skillName, skillName)
	return os.WriteFile(filepath.Join(targetDir, "tool.json"), []byte(toolSchema), 0644)
}

func (a *ClaudeCodeAdapter) writeSafeExecuteScript(targetDir, skillName string) error {
	// Defense against shell injection: strict quoted invocation without shell interpolation
	handler := fmt.Sprintf("#!/bin/sh\n# Claude Code Bridge to USC (Injection Safe)\nexec usc run %q \"$@\"\n", skillName)
	return os.WriteFile(filepath.Join(targetDir, "execute.sh"), []byte(handler), 0755)
}

func (a *ClaudeCodeAdapter) Install(bundleDir string, destDir string) (string, error) {
	if destDir == "" {
		dest, _ := a.DetectInstalled()
		destDir = dest
	}
	skillName := CleanSkillName(filepath.Base(bundleDir))
	return SafeInstallDirectory(bundleDir, destDir, skillName)
}
