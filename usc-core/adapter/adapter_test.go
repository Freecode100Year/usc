package adapter

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Freecode100Year/usc/usc-core/attestation"
	"github.com/Freecode100Year/usc/usc-core/capability"
)

func TestAllAdaptersBundleGeneration(t *testing.T) {
	outDir := t.TempDir()
	defer os.RemoveAll(outDir)

	att := &attestation.Attestation{
		Version: "0.1.0",
		Artifact: attestation.ArtifactInfo{
			Name: "weather-tool",
		},
		Claims: attestation.Claims{
			CleanroomIsolated:      true,
			AttackSurfaceReduction: 0.925,
		},
		ProofReferences: attestation.ProofReferences{
			AuditChainRoot: "sha256:chainroot123",
		},
	}
	bp := capability.NewSet()

	for _, target := range SupportedTargets {
		ad, err := GetAdapter(target)
		if err != nil {
			t.Fatalf("failed to get adapter for %s: %v", target, err)
		}
		bundleDir, err := ad.GenerateBundle(att, bp, outDir)
		if err != nil {
			t.Fatalf("failed to generate bundle for %s: %v", target, err)
		}
		if _, err := os.Stat(bundleDir); err != nil {
			t.Fatalf("bundle dir missing for %s: %v", target, err)
		}
		instDir := filepath.Join(outDir, "installed-"+target)
		if _, err := ad.Install(bundleDir, instDir); err != nil {
			t.Fatalf("install failed for %s: %v", target, err)
		}
	}
}
