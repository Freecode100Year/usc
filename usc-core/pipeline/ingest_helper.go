package pipeline

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Freecode100Year/usc/usc-core/capability"
	"github.com/Freecode100Year/usc/usc-core/intent"
)

type ingestedMeta struct {
	Name      string   `json:"name"`
	Endpoints []string `json:"endpoints"`
	EnvVars   []string `json:"env_vars"`
}

func extractIngestMetadata(skillName, sourcePath string) ([]intent.DeclaredStep, []intent.ObservedFact) {
	metaFile := filepath.Join(sourcePath, "metadata.json")
	data, err := os.ReadFile(metaFile)
	if err != nil {
		return defaultIngestFacts(skillName)
	}
	var meta ingestedMeta
	if err := json.Unmarshal(data, &meta); err != nil || len(meta.Endpoints) == 0 {
		return defaultIngestFacts(skillName)
	}
	return buildFactsFromMeta(skillName, &meta)
}

func buildFactsFromMeta(skillName string, meta *ingestedMeta) ([]intent.DeclaredStep, []intent.ObservedFact) {
	var steps []intent.DeclaredStep
	var facts []intent.ObservedFact
	for i, ep := range meta.Endpoints {
		stepID := fmt.Sprintf("step_%d", i+1)
		steps = append(steps, intent.DeclaredStep{
			StepID: stepID,
			Action: "http.get",
			Target: ep,
		})
		facts = append(facts, intent.ObservedFact{
			FactID: fmt.Sprintf("obs_%d", i+1),
			Capability: capability.Capability{
				CapabilityID: fmt.Sprintf("cap_obs_%d", i+1),
				Kind:         capability.KindNetHTTP,
				Actions:      []string{"GET"},
				Resource: capability.Resource{
					Scheme: "https",
					Host:   capability.HostSpec{Type: "EXACT", Value: ep},
					Path:   "/*",
				},
			},
		})
	}
	return steps, facts
}

func defaultIngestFacts(skillName string) ([]intent.DeclaredStep, []intent.ObservedFact) {
	steps := []intent.DeclaredStep{
		{StepID: "main_step", Action: "http.get", Target: "api.github.com"},
	}
	facts := []intent.ObservedFact{
		{
			FactID: "obs_1",
			Capability: capability.Capability{
				CapabilityID: "cap_obs_1",
				Kind:         capability.KindNetHTTP,
				Actions:      []string{"GET"},
				Resource: capability.Resource{
					Scheme: "https",
					Host:   capability.HostSpec{Type: "EXACT", Value: "api.github.com"},
					Path:   "/repos/*/pulls/*",
				},
			},
		},
	}
	return steps, facts
}
