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
	if _, err := os.Stat(filepath.Dir(p)); err == nil {
		return p, true
	}
	return filepath.Join(home, ".codex", "tools"), false
}

func (a *CodexAdapter) GenerateBundle(att *attestation.Attestation, bp capability.Set, outDir string) (string, error) {
	targetDir := filepath.Join(outDir, "codex-"+att.Artifact.Name)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return "", err
	}
	fnDef := fmt.Sprintf(`{
  "type": "function",
  "function": {
    "name": "%s",
    "description": "USC zero-trust verified tool for %s",
    "parameters": { "type": "object", "properties": { "action": { "type": "string" } } }
  }
}`, att.Artifact.Name, att.Artifact.Name)
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
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return "", err
	}
	finalPath := filepath.Join(destDir, filepath.Base(bundleDir))
	return finalPath, copyDirectory(bundleDir, finalPath)
}
