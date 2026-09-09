package adapter

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Freecode100Year/usc/usc-core/attestation"
	"github.com/Freecode100Year/usc/usc-core/capability"
)

const (
	TargetOpenClaw   = "openclaw"
	TargetHermes     = "hermes"
	TargetClaudeCode = "claudecode"
	TargetAGYCLI     = "agycli"
	TargetCodex      = "codex"
)

// SupportedTargets lists the 5 supported agent runtimes.
var SupportedTargets = []string{
	TargetOpenClaw,
	TargetHermes,
	TargetClaudeCode,
	TargetAGYCLI,
	TargetCodex,
}

// RuntimeAdapter translates an attested artifact into a platform-native skill.
type RuntimeAdapter interface {
	Target() string
	DisplayName() string
	DetectInstalled() (string, bool)
	GenerateBundle(att *attestation.Attestation, bp capability.Set, outDir string) (string, error)
	Install(bundleDir string, destDir string) (string, error)
}

// GetAdapter retrieves the adapter instance for a target name.
func GetAdapter(target string) (RuntimeAdapter, error) {
	norm := strings.ToLower(strings.TrimSpace(target))
	switch norm {
	case TargetOpenClaw:
		return NewOpenClawAdapter(), nil
	case TargetHermes:
		return NewHermesAdapter(), nil
	case TargetClaudeCode:
		return NewClaudeCodeAdapter(), nil
	case TargetAGYCLI:
		return NewAGYCLIAdapter(), nil
	case TargetCodex:
		return NewCodexAdapter(), nil
	default:
		return nil, fmt.Errorf("unsupported agent runtime: %s (supported: %s)", target, strings.Join(SupportedTargets, ", "))
	}
}

var ErrArtifactNotAttested = errors.New("USC-E7050: cannot export unverified artifact; target requires ATTESTED status")
