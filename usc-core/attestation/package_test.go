package attestation

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPackageAndVerifyArtifact(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "usc-pkg-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	ts := &TrustStore{
		KeysDir:  filepath.Join(tmpDir, "keys"),
		TrustDir: filepath.Join(tmpDir, "trust"),
	}
	pub, priv, keyID, err := ts.LoadOrCreateAuthority()
	_ = pub
	if err != nil {
		t.Fatalf("failed to create authority: %v", err)
	}

	payloadDir := filepath.Join(tmpDir, "src_payload")
	_ = os.MkdirAll(payloadDir, 0755)
	_ = os.WriteFile(filepath.Join(payloadDir, "SKILL.md"), []byte("# Verified Skill"), 0644)

	proofDir := filepath.Join(tmpDir, "src_proof")
	_ = os.MkdirAll(proofDir, 0755)
	for _, f := range requiredProofFiles {
		_ = os.WriteFile(filepath.Join(proofDir, f), []byte(`{"test":true}`), 0644)
	}

	att := &Attestation{
		Schema:  "https://usc.dev/spec/v0.1/attestation.json",
		Version: "0.1.0",
		Artifact: ArtifactInfo{
			Name: "test-skill",
		},
		Claims: Claims{
			CleanroomIsolated: true,
		},
	}

	destZip := filepath.Join(tmpDir, "test-skill.usc")
	if err := SignAndPackageArtifact(destZip, att, priv, keyID, proofDir, payloadDir); err != nil {
		t.Fatalf("package failed: %v", err)
	}

	// 1. Verify valid container
	res, err := VerifyArtifactContainer(destZip, ts)
	if err != nil || !res.Valid {
		t.Fatalf("verification should succeed, got: %v", err)
	}
	defer os.RemoveAll(res.ExtractedDir)

	// 2. Verify non-existent container
	_, err = VerifyArtifactContainer(filepath.Join(tmpDir, "not-found.usc"), ts)
	if err == nil {
		t.Fatal("expected error on non-existent file, got nil")
	}

	// 3. Verify untrusted authority
	otherPub, otherPriv, _ := GenerateKeyPair()
	_ = otherPub
	attUntrusted := *att
	destUntrusted := filepath.Join(tmpDir, "untrusted.usc")
	_ = SignAndPackageArtifact(destUntrusted, &attUntrusted, otherPriv, "unknown-key-id", proofDir, payloadDir)
	_, err = VerifyArtifactContainer(destUntrusted, ts)
	if err == nil {
		t.Fatal("expected error on untrusted authority key, got nil")
	}
}
