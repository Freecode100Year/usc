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
	bins := []string{"hermes", "hermes.cmd", "hermes.exe"}
	cfgs := []string{
		filepath.Join(home, ".hermes", "config.yaml"),
		filepath.Join(home, ".hermes", "hermes.json"),
	}
	detected := CheckBinaryOrConfig(bins, cfgs)
	return p, detected
}

func (a *HermesAdapter) GenerateBundle(att *attestation.Attestation, bp capability.Set, outDir string) (string, error) {
	skillName := CleanSkillName(att.Artifact.Name)
	if err := ValidateSkillName(skillName); err != nil {
		return "", err
	}
	targetDir := filepath.Join(outDir, "hermes-"+skillName)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return "", err
	}
	manifest := fmt.Sprintf(`{"name":%q,"version":"%s","engine":"usc-hermes","attested":true}`,
		skillName, att.Version)
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
	skillName := CleanSkillName(filepath.Base(bundleDir))
	return SafeInstallDirectory(bundleDir, destDir, skillName)
}
