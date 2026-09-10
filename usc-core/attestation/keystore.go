package attestation

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// TrustStore manages persistent authority keys and trusted public keys.
type TrustStore struct {
	KeysDir  string
	TrustDir string
}

func DefaultTrustStore() *TrustStore {
	home, _ := os.UserHomeDir()
	base := filepath.Join(home, ".usc")
	return &TrustStore{
		KeysDir:  filepath.Join(base, "keys"),
		TrustDir: filepath.Join(base, "trust"),
	}
}

func (ts *TrustStore) LoadOrCreateAuthority() (ed25519.PublicKey, ed25519.PrivateKey, string, error) {
	if err := os.MkdirAll(ts.KeysDir, 0700); err != nil {
		return nil, nil, "", err
	}
	privPath := filepath.Join(ts.KeysDir, "authority.priv")
	pubPath := filepath.Join(ts.KeysDir, "authority.pub")
	if data, err := os.ReadFile(privPath); err == nil && len(data) == ed25519.PrivateKeySize {
		pubData, _ := os.ReadFile(pubPath)
		pub := ed25519.PublicKey(pubData)
		keyID := "usc-local-" + hex.EncodeToString(pub[:4])
		return pub, ed25519.PrivateKey(data), keyID, nil
	}
	return ts.generateAndSaveAuthority(privPath, pubPath)
}

func (ts *TrustStore) generateAndSaveAuthority(privPath, pubPath string) (ed25519.PublicKey, ed25519.PrivateKey, string, error) {
	pub, priv, err := GenerateKeyPair()
	if err != nil {
		return nil, nil, "", err
	}
	if err := os.WriteFile(privPath, priv, 0600); err != nil {
		return nil, nil, "", err
	}
	if err := os.WriteFile(pubPath, pub, 0644); err != nil {
		return nil, nil, "", err
	}
	keyID := "usc-local-" + hex.EncodeToString(pub[:4])
	_ = ts.AddTrustedKey(keyID, pub)
	return pub, priv, keyID, nil
}

func (ts *TrustStore) AddTrustedKey(keyID string, pub ed25519.PublicKey) error {
	if err := os.MkdirAll(ts.TrustDir, 0755); err != nil {
		return err
	}
	m := ts.LoadTrustedKeys()
	m[keyID] = hex.EncodeToString(pub)
	data, _ := json.MarshalIndent(m, "", "  ")
	return os.WriteFile(filepath.Join(ts.TrustDir, "trusted_authorities.json"), data, 0644)
}

func (ts *TrustStore) LoadTrustedKeys() map[string]string {
	p := filepath.Join(ts.TrustDir, "trusted_authorities.json")
	data, err := os.ReadFile(p)
	if err != nil {
		return make(map[string]string)
	}
	var m map[string]string
	if err := json.Unmarshal(data, &m); err != nil {
		return make(map[string]string)
	}
	return m
}

func (ts *TrustStore) ResolvePublicKey(keyID string) (ed25519.PublicKey, error) {
	keys := ts.LoadTrustedKeys()
	hexPub, ok := keys[keyID]
	if !ok {
		return nil, fmt.Errorf("USC-E7045: untrusted authority key ID %s", keyID)
	}
	raw, err := hex.DecodeString(hexPub)
	if err != nil || len(raw) != ed25519.PublicKeySize {
		return nil, fmt.Errorf("USC-E7046: corrupted public key for %s", keyID)
	}
	return ed25519.PublicKey(raw), nil
}
