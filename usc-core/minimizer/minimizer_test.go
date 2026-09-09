package minimizer

import (
	"testing"

	"github.com/Freecode100Year/usc/usc-core/capability"
	"github.com/Freecode100Year/usc/usc-core/intent"
)

func TestGreedySolver(t *testing.T) {
	solver := NewGreedySolver()
	in := intent.DeclaredIntent{
		SkillName: "weather-cli",
		Steps: []intent.DeclaredStep{
			{StepID: "fetch_weather", Action: "http.get", Target: "api.weather.gov"},
		},
	}
	wf := Workflow{WorkflowID: "wf1", StepOrder: []string{"fetch_weather"}}
	target := RuntimeSpec{
		TargetPlatform: "hermes",
		SupportedKinds: []string{"net.http", "fs.read"},
	}
	policy := Policy{
		PolicyID: "enterprise_standard",
		AllowedCapabilities: capability.NewSet(
			capability.Capability{
				CapabilityID: "allow_all_net",
				Kind:         capability.KindNetHTTP,
				Actions:      []string{"GET", "POST"},
				Resource: capability.Resource{
					Host: capability.HostSpec{Type: "WILDCARD", Value: "*"},
					Path: "/*",
				},
			},
		),
	}
	res, proof, err := solver.Solve(in, wf, target, policy)
	if err != nil || !proof.Satisfied || len(res) == 0 {
		t.Fatalf("solver failed: %v", err)
	}
}
