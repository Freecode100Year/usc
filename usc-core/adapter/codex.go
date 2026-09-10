package adapter

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Freecode100Year/usc/usc-core/attestation"
	"github.com/Freecode100Year/usc/usc-core/capability"
)

// CodexAdapter adapts attested USC skills to OpenAI Codex / Assistants runtime.
type CodexAdapter struct{}

func NewCodexAdapter() *CodexAdapter {
	return &CodexAdapter{}
}

func (a *CodexAdapter) Target() string {
	return TargetCodex
}

func (a *CodexAdapter) DisplayName() string {
	return "OpenAI Codex / Assistants"
}

func (a *CodexAdapter) DetectInstalled() (string, bool) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", false
	}
	p := filepath.Join(home, ".codex", "tools")
	bins := []string{"codex", "codex.cmd", "codex.exe"}
	cfgs := []string{
		filepath.Join(home, ".codex", "config.json"),
		filepath.Join(home, ".codex", "codex.json"),
	}
	detected := CheckBinaryOrConfig(bins, cfgs)
	return p, detected
}

func (a *CodexAdapter) GenerateBundle(att *attestation.Attestation, bp capability.Set, outDir string) (string, error) {
	skillName := CleanSkillName(att.Artifact.Name)
	if err := ValidateSkillName(skillName); err != nil {
		return "", err
	}
	targetDir := filepath.Join(outDir, "codex-"+skillName)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return "", err
	}
	fnDef := fmt.Sprintf(`{
  "type": "function",
  "function": {
    "name": %q,
    "description": "USC zero-trust verified tool for %s",
    "parameters": { "type": "object", "properties": { "action": { "type": "string" } } }
  }
}`, skillName, skillName)
	if err := os.WriteFile(filepath.Join(targetDir, "function.json"), []byte(fnDef), 0644); err != nil {
		return "", err
	}
	wrapper := "// Codex Function Invocation Bridge\nmodule.exports = async function(args) { console.log('[Codex] Invoked USC verified tool'); };\n"
	if err := os.WriteFile(filepath.Join(targetDir, "index.js"), []byte(wrapper), 0644); err != nil {
		return "", err
	}
	return targetDir, nil
}

func (a *CodexAdapter) Install(bundleDir string, destDir string) (string, error) {
	if destDir == "" {
		dest, _ := a.DetectInstalled()
		destDir = dest
	}
	skillName := CleanSkillName(filepath.Base(bundleDir))
	return SafeInstallDirectory(bundleDir, destDir, skillName)
}
