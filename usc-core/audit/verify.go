package audit

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
)

var ErrAuditChainBroken = errors.New("USC-E7005: audit hash chain broken or tampered")

// ComputeGenesisDigest computes H0 for the audit chain.
func ComputeGenesisDigest(source, blueprint, policy, compiler string) string {
	h := sha256.New()
	h.Write([]byte("USC-AUDIT-v1"))
	h.Write([]byte(source))
	h.Write([]byte(blueprint))
	h.Write([]byte(policy))
	h.Write([]byte(compiler))
	return hex.EncodeToString(h.Sum(nil))
}

// VerifyChain independently re-verifies the entire audit hash chain.
func VerifyChain(genesis string, events []Event, expectedLatest string) (bool, error) {
	curr := genesis
	for _, evt := range events {
		if evt.PreviousDigest != curr {
			return false, ErrAuditChainBroken
		}
		raw, err := json.Marshal(evt)
		if err != nil {
			return false, err
		}
		canon, err := Canonicalize(raw)
		if err != nil {
			return false, err
		}
		h := sha256.New()
		h.Write([]byte(curr))
		h.Write(canon)
		curr = hex.EncodeToString(h.Sum(nil))
	}
	if expectedLatest != "" && curr != expectedLatest {
		return false, ErrAuditChainBroken
	}
	return true, nil
}
