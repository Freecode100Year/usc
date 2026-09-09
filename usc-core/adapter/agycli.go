package adapter

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Freecode100Year/usc/usc-core/attestation"
	"github.com/Freecode100Year/usc/usc-core/capability"
)

// AGYCLIAdapter adapts attested USC skills to Google Antigravity CLI (agy).
type AGYCLIAdapter struct{}

func NewAGYCLIAdapter() *AGYCLIAdapter {
	return &AGYCLIAdapter{}
}

func (a *AGYCLIAdapter) Target() string {
	return TargetAGYCLI
}

func (a *AGYCLIAdapter) DisplayName() string {
	return "Antigravity CLI (agy)"
}

func (a *AGYCLIAdapter) DetectInstalled() (string, bool) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", false
	}
	p := filepath.Join(home, ".gemini", "antigravity-cli")
	if _, err := os.Stat(p); err == nil {
		return filepath.Join(p, "skills"), true
	}
	return filepath.Join(home, ".gemini", "antigravity-cli", "skills"), false
}

func (a *AGYCLIAdapter) GenerateBundle(att *attestation.Attestation, bp capability.Set, outDir string) (string, error) {
	targetDir := filepath.Join(outDir, "agy-"+att.Artifact.Name)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return "", err
	}
	skillMd := fmt.Sprintf(`---
name: %s
description: Zero-trust verified skill compiled by USC (cleanroom-isolated, attested).
---

# %s (USC Attested)

> Attestation ID: %s
> Cleanroom Isolated: %v
> Attack Surface Reduction: %.1f%%

This skill runs inside the USC zero-trust confinement boundary.
`, att.Artifact.Name, att.Artifact.Name, att.ProofReferences.AuditChainRoot, att.Claims.CleanroomIsolated, att.Claims.AttackSurfaceReduction*100)
	if err := os.WriteFile(filepath.Join(targetDir, "SKILL.md"), []byte(skillMd), 0644); err != nil {
		return "", err
	}
	return targetDir, nil
}

func (a *AGYCLIAdapter) Install(bundleDir string, destDir string) (string, error) {
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
