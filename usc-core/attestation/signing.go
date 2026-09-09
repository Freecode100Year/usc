package attestation

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"

	"github.com/Freecode100Year/usc/usc-core/audit"
)

// GenerateKeyPair creates a fresh Ed25519 key pair for the attestation authority.
func GenerateKeyPair() (ed25519.PublicKey, ed25519.PrivateKey, error) {
	return ed25519.GenerateKey(rand.Reader)
}

// SignAttestation signs the canonical JSON of the unsigned attestation payload.
func SignAttestation(att *Attestation, priv ed25519.PrivateKey, keyID string) error {
	att.Signature = Signature{KeyID: keyID, Algorithm: "ed25519"}
	bytesToSign, err := computeAttestationSignBytes(att)
	if err != nil {
		return err
	}
	sig := ed25519.Sign(priv, bytesToSign)
	att.Signature.Sig = hex.EncodeToString(sig)
	return nil
}

// VerifySignature checks the cryptographic signature using the provided public key.
func VerifySignature(att *Attestation, pub ed25519.PublicKey) (bool, error) {
	sigBytes, err := hex.DecodeString(att.Signature.Sig)
	if err != nil {
		return false, err
	}
	bytesToSign, err := computeAttestationSignBytes(att)
	if err != nil {
		return false, err
	}
	return ed25519.Verify(pub, bytesToSign, sigBytes), nil
}

func computeAttestationSignBytes(att *Attestation) ([]byte, error) {
	unsignedCopy := *att
	unsignedCopy.Signature.Sig = ""
	raw, err := json.Marshal(unsignedCopy)
	if err != nil {
		return nil, err
	}
	return audit.Canonicalize(raw)
}
