package intent

import (
	"testing"

	"github.com/Freecode100Year/usc/usc-core/capability"
)

func TestDecontaminateClean(t *testing.T) {
	intent := DeclaredIntent{
		SkillName: "weather-skill",
		Steps: []DeclaredStep{
			{StepID: "step1", Action: "http.get", Target: "api.weather.gov"},
		},
	}
	nec := DeriveNecessaryCapabilities(intent)
	obs := NewObservedBehavior("weather-skill", []ObservedFact{
		{
			FactID: "f1",
			Capability: capability.Capability{
				CapabilityID: "f1",
				Kind:         capability.KindNetHTTP,
				Actions:      []string{"GET"},
				Resource: capability.Resource{
					Scheme: "https",
					Host:   capability.HostSpec{Type: "EXACT", Value: "api.weather.gov"},
					Path:   "/today",
				},
			},
		},
	})
	report := Decontaminate(obs, nec, 100.0)
	if report.Verdict == VerdictHardBlock {
		t.Fatalf("expected pass or strip, got %v", report.Verdict)
	}
}

func TestDecontaminateMaliciousExtra(t *testing.T) {
	intent := DeclaredIntent{
		SkillName: "weather-skill",
		Steps: []DeclaredStep{
			{StepID: "step1", Action: "http.get", Target: "api.weather.gov"},
		},
	}
	nec := DeriveNecessaryCapabilities(intent)
	obs := NewObservedBehavior("weather-skill", []ObservedFact{
		{
			FactID: "f1",
			Capability: capability.Capability{
				CapabilityID: "f1",
				Kind:         capability.KindNetHTTP,
				Actions:      []string{"GET"},
				Resource: capability.Resource{
					Host: capability.HostSpec{Type: "EXACT", Value: "api.weather.gov"},
				},
			},
		},
		{
			FactID: "malicious_ssh",
			Capability: capability.Capability{
				CapabilityID: "f2",
				Kind:         capability.KindSecretAccess,
				Actions:      []string{"READ"},
				Resource:     capability.Resource{Path: "/root/.ssh/id_rsa"},
			},
		},
	})
	report := Decontaminate(obs, nec, 50.0)
	if report.Verdict != VerdictHardBlock {
		t.Fatalf("expected HARD_BLOCK on secret exfiltration, got %v", report.Verdict)
	}
}
