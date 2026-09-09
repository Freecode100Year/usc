package audit

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"sync"
	"time"
)

// Ledger maintains an append-only lossless audit log secured by a hash chain.
type Ledger struct {
	mu           sync.RWMutex
	genesisHash  string
	latestDigest string
	events       []Event
}

// NewLedger initializes the audit ledger with a cryptographic genesis digest.
func NewLedger(genesis string) *Ledger {
	return &Ledger{
		genesisHash:  genesis,
		latestDigest: genesis,
		events:       make([]Event, 0),
	}
}

// Append records a new event onto the ledger and computes the next hash.
func (l *Ledger) Append(stage Stage, rule string, decision Decision, payload string) (Event, string, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	evt := Event{
		Seq:            uint64(len(l.events) + 1),
		PreviousDigest: l.latestDigest,
		Stage:          stage,
		Rule:           rule,
		Decision:       decision,
		Payload:        payload,
		Timestamp:      time.Now().UnixNano(),
	}

	raw, err := json.Marshal(evt)
	if err != nil {
		return Event{}, "", err
	}
	canon, err := Canonicalize(raw)
	if err != nil {
		return Event{}, "", err
	}

	h := sha256.New()
	h.Write([]byte(l.latestDigest))
	h.Write(canon)
	nextDigest := hex.EncodeToString(h.Sum(nil))

	l.latestDigest = nextDigest
	l.events = append(l.events, evt)
	return evt, nextDigest, nil
}

// Events returns a copy of all ledger events.
func (l *Ledger) Events() []Event {
	l.mu.RLock()
	defer l.mu.RUnlock()
	copied := make([]Event, len(l.events))
	copy(copied, l.events)
	return copied
}

// LatestDigest returns the current root digest of the audit chain.
func (l *Ledger) LatestDigest() string {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.latestDigest
}
