package adapter

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Freecode100Year/usc/usc-core/attestation"
	"github.com/Freecode100Year/usc/usc-core/capability"
)

// HermesAdapter adapts attested USC skills to Hermes Agent runtime.
type HermesAdapter struct{}

func NewHermesAdapter() *HermesAdapter {
	return &HermesAdapter{}
}

func (a *HermesAdapter) Target() string {
	return TargetHermes
}

func (a *HermesAdapter) DisplayName() string {
	return "Hermes Autonomous Agent"
}

func (a *HermesAdapter) DetectInstalled() (string, bool) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", false
	}
	p := filepath.Join(home, ".hermes", "skills")
	if _, err := os.Stat(filepath.Dir(p)); err == nil {
		return p, true
	}
	return filepath.Join(home, ".hermes", "skills"), false
}

func (a *HermesAdapter) GenerateBundle(att *attestation.Attestation, bp capability.Set, outDir string) (string, error) {
	targetDir := filepath.Join(outDir, "hermes-"+att.Artifact.Name)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return "", err
	}
	manifest := fmt.Sprintf(`{"name":"%s","version":"%s","engine":"usc-hermes","attested":true}`,
		att.Artifact.Name, att.Version)
	if err := os.WriteFile(filepath.Join(targetDir, "manifest.json"), []byte(manifest), 0644); err != nil {
		return "", err
	}
	jsContent := "// Hermes Agent USC Entrypoint\nconsole.log('[Hermes] Running clean-room rebuilt agent capability');\n"
	if err := os.WriteFile(filepath.Join(targetDir, "index.js"), []byte(jsContent), 0644); err != nil {
		return "", err
	}
	return targetDir, nil
}

func (a *HermesAdapter) Install(bundleDir string, destDir string) (string, error) {
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
