package capability

import "testing"

func TestPosetAxioms(t *testing.T) {
	c1 := Capability{
		CapabilityID: "cap1",
		Kind:         KindNetHTTP,
		Actions:      []string{"GET"},
		Resource: Resource{
			Scheme: "https",
			Host:   HostSpec{Type: "EXACT", Value: "api.github.com"},
			Path:   "/repos/*/pulls/*",
		},
	}
	c2 := Capability{
		CapabilityID: "cap2",
		Kind:         KindNetHTTP,
		Actions:      []string{"GET", "POST"},
		Resource: Resource{
			Scheme: "https",
			Host:   HostSpec{Type: "WILDCARD", Value: "*.github.com"},
			Path:   "/repos/*",
		},
	}
	if !VerifyReflexive(c1) {
		t.Fatalf("reflexivity failed for c1")
	}
	if !Subsumes(c2, c1) {
		t.Fatalf("expected c2 to subsume c1")
	}
	if !VerifyAntisymmetric(c1, c2) {
		t.Fatalf("antisymmetry verification failed")
	}
}

func TestMeetAndReduction(t *testing.T) {
	c1 := Capability{
		CapabilityID: "c1", Kind: KindNetHTTP, Actions: []string{"GET"},
		Resource: Resource{Host: HostSpec{Type: "EXACT", Value: "api.github.com"}},
	}
	c2 := Capability{
		CapabilityID: "c2", Kind: KindNetHTTP, Actions: []string{"GET", "POST"},
		Resource: Resource{Host: HostSpec{Type: "EXACT", Value: "api.github.com"}},
	}
	m, ok := Meet(c1, c2)
	if !ok || len(m.Actions) != 1 || m.Actions[0] != "GET" {
		t.Fatalf("unexpected meet result: %+v", m)
	}

	obs := NewSet(c1, c2)
	grant := NewSet(c1)
	ccr, asr := ComputeReduction(obs, grant, DefaultWeights())
	if ccr <= 0 || asr <= 0 {
		t.Fatalf("expected positive reduction, got ccr=%f, asr=%f", ccr, asr)
	}
}
