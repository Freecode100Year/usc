package pipeline

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Freecode100Year/usc/usc-core/attestation"
	"github.com/Freecode100Year/usc/usc-core/minimizer"
)

func TestCompilerPipelineEndToEnd(t *testing.T) {
	p := NewPipelineState("test-skill", "a1b2c3d4e5f6")
	if err := p.RunIngest("dummy/path"); err != nil {
		t.Fatalf("ingest failed: %v", err)
	}
	if err := p.RunDecontaminate(100.0); err != nil {
		t.Fatalf("decontaminate failed: %v", err)
	}
	target := minimizer.RuntimeSpec{TargetPlatform: "hermes", SupportedKinds: []string{"net.http"}}
	policy := minimizer.Policy{PolicyID: "allow_all"}
	if err := p.RunMinimize(target, policy); err != nil {
		t.Fatalf("minimize failed: %v", err)
	}
	if err := p.RunCleanRebuild(); err != nil {
		t.Fatalf("clean rebuild failed: %v", err)
	}
	if err := p.RunReAudit(); err != nil {
		t.Fatalf("re-audit failed: %v", err)
	}
	if err := p.RunSandbox(); err != nil {
		t.Fatalf("sandbox failed: %v", err)
	}
	pub, priv, _ := attestation.GenerateKeyPair()
	distDir := filepath.Join(t.TempDir(), "dist")
	if err := p.RunAttest(distDir, priv, "test-authority"); err != nil {
		t.Fatalf("attest failed: %v", err)
	}

	att, err := attestation.VerifyProofBundle(filepath.Join(distDir, "proof"), pub)
	if err != nil || att == nil {
		t.Fatalf("proof bundle verification failed: %v", err)
	}
	defer os.RemoveAll(distDir)
}
