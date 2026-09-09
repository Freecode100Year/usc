package replay

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Freecode100Year/usc/usc-core/capability"
)

func TestFlightRecorderAndReplay(t *testing.T) {
	fr := NewFlightRecorder(5)
	fr.Record(TraceEvent{Type: "CAP", Target: "api.github.com", Action: "GET", Decision: "ALLOW"})
	fr.Record(TraceEvent{Type: "CAP", Target: "evil.com", Action: "POST", Decision: "DENY"})

	tmpFile := filepath.Join(t.TempDir(), "test.usctrace")
	if err := fr.Dump(tmpFile, "test-skill", "digest123"); err != nil {
		t.Fatalf("dump failed: %v", err)
	}
	defer os.Remove(tmpFile)

	tf, err := LoadTrace(tmpFile)
	if err != nil || len(tf.Events) != 2 {
		t.Fatalf("load failed: %v", err)
	}

	policy := capability.NewSet(capability.Capability{
		Kind:    capability.KindNetHTTP,
		Actions: []string{"GET"},
		Resource: capability.Resource{
			Host: capability.HostSpec{Type: "EXACT", Value: "api.github.com"},
		},
	})
	results, allMatch := ReplayDecisions(tf, policy)
	if !allMatch || len(results) != 2 {
		t.Fatalf("replay mismatch: %v", results)
	}
}
