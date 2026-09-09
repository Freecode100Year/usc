package attestation

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestSigningAndVerification(t *testing.T) {
	pub, priv, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("key generation failed: %v", err)
	}
	att := Attestation{
		Schema:  "https://usc.dev/spec/v0.1/attestation.json",
		Version: "0.1.0",
		Artifact: ArtifactInfo{
			Name:           "weather.usc",
			SHA256:         "deadbeef",
			TargetPlatform: "hermes",
		},
		Claims: Claims{
			SourceTaintQuarantined:            true,
			CapabilityCountReduction:          0.714,
			AttackSurfaceReduction:            0.925,
			CleanroomIsolated:                 true,
			GeneratedAuthorityWithinBlueprint: true,
			DryRunValidation:                  "SMOKE_VALIDATED",
		},
		ProofReferences: ProofReferences{
			AuditChainRoot: "sha256:root",
		},
	}
	if err := SignAttestation(&att, priv, "test-key"); err != nil {
		t.Fatalf("sign failed: %v", err)
	}
	valid, err := VerifySignature(&att, pub)
	if err != nil || !valid {
		t.Fatalf("verification failed: %v", err)
	}
}

func TestVerifyProofBundle(t *testing.T) {
	dir := t.TempDir()
	pub, priv, _ := GenerateKeyPair()
	att := Attestation{Schema: "https://usc.dev/spec/v0.1/attestation.json", Version: "0.1.0"}
	if err := SignAttestation(&att, priv, "key1"); err != nil {
		t.Fatalf("sign error: %v", err)
	}

	for _, f := range requiredProofFiles {
		os.WriteFile(filepath.Join(dir, f), []byte("{}"), 0644)
	}
	attBytes, _ := json.Marshal(att)
	os.WriteFile(filepath.Join(dir, "attestation.json"), attBytes, 0644)

	if _, err := VerifyProofBundle(dir, pub); err != nil {
		t.Fatalf("expected proof bundle validation pass, got: %v", err)
	}
}
