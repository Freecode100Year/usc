package replay

import (
	"encoding/json"
	"os"
	"sync"
)

// TraceEvent represents a single event recorded by the flight recorder.
type TraceEvent struct {
	Index     int    `json:"index"`
	Timestamp int64  `json:"timestamp"`
	Type      string `json:"type"` // e.g., "CAPABILITY_REQUEST"
	Target    string `json:"target"`
	Action    string `json:"action"`
	Decision  string `json:"decision"` // ALLOW / DENY / HARD_BLOCK
}

// TraceFile represents the serialized .usctrace format.
type TraceFile struct {
	Version      string       `json:"version"`
	SkillID      string       `json:"skill_id"`
	PolicyDigest string       `json:"policy_digest"`
	Events       []TraceEvent `json:"events"`
}

// FlightRecorder implements an ephemeral ring buffer for execution traces.
type FlightRecorder struct {
	mu       sync.RWMutex
	capacity int
	events   []TraceEvent
}

// NewFlightRecorder initializes a ring buffer flight recorder.
func NewFlightRecorder(capacity int) *FlightRecorder {
	return &FlightRecorder{
		capacity: capacity,
		events:   make([]TraceEvent, 0, capacity),
	}
}

// Record appends an event to the ring buffer, rotating oldest if capacity exceeded.
func (r *FlightRecorder) Record(evt TraceEvent) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.events) >= r.capacity {
		r.events = r.events[1:]
	}
	evt.Index = len(r.events)
	r.events = append(r.events, evt)
}

// Dump writes the recorded trace to a .usctrace file.
func (r *FlightRecorder) Dump(filePath, skillID, policyDigest string) error {
	r.mu.RLock()
	defer r.mu.RUnlock()
	tf := TraceFile{
		Version:      "0.1.0",
		SkillID:      skillID,
		PolicyDigest: policyDigest,
		Events:       r.events,
	}
	data, err := json.MarshalIndent(tf, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0644)
}

// LoadTrace reads and parses a .usctrace file.
func LoadTrace(filePath string) (*TraceFile, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var tf TraceFile
	if err := json.Unmarshal(data, &tf); err != nil {
		return nil, err
	}
	return &tf, nil
}
