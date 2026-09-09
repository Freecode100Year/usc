package adapter

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
	cfgPath := filepath.Join(home, ".gemini", "config", "skills")
	if _, err := os.Stat(filepath.Dir(cfgPath)); err == nil {
		return cfgPath, true
	}
	cliPath := filepath.Join(home, ".gemini", "antigravity-cli", "skills")
	if _, err := os.Stat(filepath.Dir(cliPath)); err == nil {
		return cliPath, true
	}
	return cfgPath, false
}

func (a *AGYCLIAdapter) GenerateBundle(att *attestation.Attestation, bp capability.Set, outDir string) (string, error) {
	skillName := CleanSkillName(att.Artifact.Name)
	targetDir := filepath.Join(outDir, skillName)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return "", err
	}
	skillMd := buildAGYSkillMD(att, skillName)
	if err := os.WriteFile(filepath.Join(targetDir, "SKILL.md"), []byte(skillMd), 0644); err != nil {
		return "", err
	}
	copySourceAssets(att.Artifact.SourceDir, targetDir)
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

func buildAGYSkillMD(att *attestation.Attestation, skillName string) string {
	var baseDoc string
	if att.Artifact.SourceDir != "" {
		srcMD := filepath.Join(att.Artifact.SourceDir, "SKILL.md")
		if data, err := os.ReadFile(srcMD); err == nil {
			baseDoc = string(data)
		}
	}
	if baseDoc != "" {
		return injectAttestationHeader(baseDoc, att)
	}
	return defaultAGYDoc(att, skillName)
}

func injectAttestationHeader(doc string, att *attestation.Attestation) string {
	badge := fmt.Sprintf("\n> **[USC Zero-Trust Clean-Room Attested]**\n> - Attestation Root: `%s`\n> - Cleanroom Isolated: `%v`\n> - Attack Surface Reduction: `%.1f%%`\n\n",
		att.ProofReferences.AuditChainRoot, att.Claims.CleanroomIsolated, att.Claims.AttackSurfaceReduction*100)
	if idx := strings.Index(doc, "---"); idx != -1 {
		if endIdx := strings.Index(doc[idx+3:], "---"); endIdx != -1 {
			splitAt := idx + 3 + endIdx + 3
			return doc[:splitAt] + badge + doc[splitAt:]
		}
	}
	return badge + doc
}

func defaultAGYDoc(att *attestation.Attestation, skillName string) string {
	return fmt.Sprintf(`---
name: %s
description: Zero-trust verified skill compiled by USC (cleanroom-isolated, attested).
---

# %s (USC Attested)

> Attestation ID: %s
> Cleanroom Isolated: %v
> Attack Surface Reduction: %.1f%%

This skill runs inside the USC zero-trust confinement boundary.
`, skillName, skillName, att.ProofReferences.AuditChainRoot, att.Claims.CleanroomIsolated, att.Claims.AttackSurfaceReduction*100)
}

func copySourceAssets(srcDir, targetDir string) {
	if srcDir == "" {
		return
	}
	srcScripts := filepath.Join(srcDir, "scripts")
	if info, err := os.Stat(srcScripts); err == nil && info.IsDir() {
		_ = copyDirectory(srcScripts, filepath.Join(targetDir, "scripts"))
	}
}
