package attestation

import (
	"crypto/ed25519"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var (
	ErrInvalidProofBundle = errors.New("USC-E7040: proof bundle is missing mandatory verification artifacts")
	ErrSignatureInvalid   = errors.New("USC-E7041: attestation signature mismatch or invalid")
	ErrPayloadCorrupted   = errors.New("USC-E7042: payload hash does not match attested digest")
)

var requiredProofFiles = []string{
	"attestation.json",
	"audit-chain.json",
	"provenance.json",
	"cleanroom-proof.json",
	"policy.json",
	"blueprint.json",
}

// VerifyProofBundle inspects the given proof directory for all required artifacts and signature.
func VerifyProofBundle(dirPath string, pub ed25519.PublicKey) (*Attestation, error) {
	for _, reqFile := range requiredProofFiles {
		p := filepath.Join(dirPath, reqFile)
		if _, err := os.Stat(p); err != nil {
			return nil, ErrInvalidProofBundle
		}
	}
	att, err := loadAttestation(filepath.Join(dirPath, "attestation.json"))
	if err != nil {
		return nil, err
	}
	valid, err := VerifySignature(att, pub)
	if err != nil || !valid {
		return nil, ErrSignatureInvalid
	}
	return att, nil
}

// VerificationResult contains the detailed outcomes of full artifact verification.
type VerificationResult struct {
	Attestation *Attestation
	ExtractedDir string
	PayloadDir  string
	PayloadHash string
	Valid       bool
}

// VerifyArtifactContainer opens, unpacks, and cryptographically audits a .usc artifact file.
func VerifyArtifactContainer(zipPath string, ts *TrustStore) (*VerificationResult, error) {
	if _, err := os.Stat(zipPath); err != nil {
		return nil, fmt.Errorf("artifact file not found: %s", zipPath)
	}
	tmpDir, err := os.MkdirTemp("", "usc-verify-*")
	if err != nil {
		return nil, err
	}
	att, err := UnpackArtifact(zipPath, tmpDir)
	if err != nil {
		_ = os.RemoveAll(tmpDir)
		return nil, err
	}
	return checkAttestationAndPayload(tmpDir, att, ts)
}

func checkAttestationAndPayload(tmpDir string, att *Attestation, ts *TrustStore) (*VerificationResult, error) {
	pub, err := ts.ResolvePublicKey(att.Signature.KeyID)
	if err != nil {
		_ = os.RemoveAll(tmpDir)
		return nil, err
	}
	valid, err := VerifySignature(att, pub)
	if err != nil || !valid {
		_ = os.RemoveAll(tmpDir)
		return nil, ErrSignatureInvalid
	}
	return auditPayloadDigest(tmpDir, att)
}

func auditPayloadDigest(tmpDir string, att *Attestation) (*VerificationResult, error) {
	payloadDir := filepath.Join(tmpDir, "payload")
	actualHash, err := ComputeDirSHA256(payloadDir)
	if err != nil {
		_ = os.RemoveAll(tmpDir)
		return nil, fmt.Errorf("failed to hash extracted payload: %w", err)
	}
	expectedHash := att.Artifact.SHA256
	if expectedHash != "sha256:"+actualHash {
		_ = os.RemoveAll(tmpDir)
		return nil, fmt.Errorf("%w: expected %s, got sha256:%s", ErrPayloadCorrupted, expectedHash, actualHash)
	}
	return &VerificationResult{
		Attestation:  att,
		ExtractedDir: tmpDir,
		PayloadDir:   payloadDir,
		PayloadHash:  actualHash,
		Valid:        true,
	}, nil
}
