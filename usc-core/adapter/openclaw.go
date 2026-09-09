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
	p := filepath.Join(home, ".openclaw", "skills")
	if _, err := os.Stat(filepath.Dir(p)); err == nil {
		return p, true
	}
	return filepath.Join(home, ".openclaw", "skills"), false
}

func (a *OpenClawAdapter) GenerateBundle(att *attestation.Attestation, bp capability.Set, outDir string) (string, error) {
	targetDir := filepath.Join(outDir, "openclaw-"+att.Artifact.Name)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return "", err
	}
	yamlContent := fmt.Sprintf("name: %s\nversion: %s\nsandbox: usc-cleanroom\nattestation_ref: %s\n",
		att.Artifact.Name, att.Version, att.ProofReferences.AuditChainRoot)
	if err := os.WriteFile(filepath.Join(targetDir, "skill.yaml"), []byte(yamlContent), 0644); err != nil {
		return "", err
	}
	runnerPy := "#!/usr/bin/env python3\n# OpenClaw Zero-Trust Skill Entrypoint\nprint('[OpenClaw] Executing zero-trust rebuilt skill via USC runtime')\n"
	if err := os.WriteFile(filepath.Join(targetDir, "runner.py"), []byte(runnerPy), 0644); err != nil {
		return "", err
	}
	return targetDir, nil
}

func (a *OpenClawAdapter) Install(bundleDir string, destDir string) (string, error) {
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
