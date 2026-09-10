package adapter

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Freecode100Year/usc/usc-core/attestation"
	"github.com/Freecode100Year/usc/usc-core/capability"
)

// OpenClawAdapter adapts attested USC skills to OpenClaw runtime.
type OpenClawAdapter struct{}

func NewOpenClawAdapter() *OpenClawAdapter {
	return &OpenClawAdapter{}
}

func (a *OpenClawAdapter) Target() string {
	return TargetOpenClaw
}

func (a *OpenClawAdapter) DisplayName() string {
	return "OpenClaw AI Runtime"
}

func (a *OpenClawAdapter) DetectInstalled() (string, bool) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", false
	}
	skillPath := filepath.Join(home, ".openclaw", "skills")
	bins := []string{"openclaw", "openclaw.cmd", "openclaw.exe"}
	cfgs := []string{
		filepath.Join(home, ".openclaw", "openclaw.json"),
		filepath.Join(home, ".openclaw", "config.yaml"),
		filepath.Join(home, ".openclaw", "config.json"),
	}
	detected := CheckBinaryOrConfig(bins, cfgs)
	return skillPath, detected
}

func (a *OpenClawAdapter) GenerateBundle(att *attestation.Attestation, bp capability.Set, outDir string) (string, error) {
	skillName := CleanSkillName(att.Artifact.Name)
	if err := ValidateSkillName(skillName); err != nil {
		return "", err
	}
	targetDir := filepath.Join(outDir, "openclaw-"+skillName)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return "", err
	}
	skillMD := a.buildOpenClawSkillMD(att, skillName)
	if err := os.WriteFile(filepath.Join(targetDir, "SKILL.md"), []byte(skillMD), 0644); err != nil {
		return "", err
	}
	a.copySourceAssets(att.Artifact.SourceDir, targetDir)
	return targetDir, nil
}

func (a *OpenClawAdapter) buildOpenClawSkillMD(att *attestation.Attestation, skillName string) string {
	return fmt.Sprintf(`---
name: %s
description: %s
metadata:
  openclaw:
    sandbox: usc-cleanroom
    attestation_ref: %s
---

# %s (OpenClaw Verified Skill)

> Attested by USC Zero-Trust Clean-Room Compiler.
> - Attestation Digest: %s
> - Cleanroom Isolated: %v
> - Attack Surface Reduction: %.1f%%

This skill conforms to OpenClaw standard SKILL.md specification.
`, skillName, "Zero-trust verified skill compiled by USC", att.ProofReferences.AuditChainRoot,
		skillName, att.ProofReferences.AuditChainRoot, att.Claims.CleanroomIsolated, att.Claims.AttackSurfaceReduction*100)
}

func (a *OpenClawAdapter) copySourceAssets(srcDir, targetDir string) {
	if srcDir == "" {
		return
	}
	srcScripts := filepath.Join(srcDir, "scripts")
	if info, err := os.Stat(srcScripts); err == nil && info.IsDir() {
		_ = SafeCopyDirectory(srcScripts, filepath.Join(targetDir, "scripts"))
	}
}

func (a *OpenClawAdapter) Install(bundleDir string, destDir string) (string, error) {
	if destDir == "" {
		dest, _ := a.DetectInstalled()
		destDir = dest
	}
	skillName := CleanSkillName(filepath.Base(bundleDir))
	return SafeInstallDirectory(bundleDir, destDir, skillName)
}
