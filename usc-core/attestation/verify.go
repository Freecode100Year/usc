package attestation

import (
	"crypto/ed25519"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

var (
	ErrInvalidProofBundle = errors.New("USC-E7040: proof bundle is missing mandatory verification artifacts")
	ErrSignatureInvalid   = errors.New("USC-E7041: attestation signature mismatch or invalid")
)

// ProofBundlePaths lists required files within a machine proof bundle directory.
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
	attPath := filepath.Join(dirPath, "attestation.json")
	data, err := os.ReadFile(attPath)
	if err != nil {
		return nil, err
	}
	var att Attestation
	if err := json.Unmarshal(data, &att); err != nil {
		return nil, err
	}
	valid, err := VerifySignature(&att, pub)
	if err != nil || !valid {
		return nil, ErrSignatureInvalid
	}
	return &att, nil
}
