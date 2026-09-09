package intent

import "github.com/Freecode100Year/usc/usc-core/capability"

// NecessaryBehavior defines the minimum behavior set B_nec required to fulfill intent.
type NecessaryBehavior struct {
	IntentID     string         `json:"intent_id"`
	Capabilities capability.Set `json:"capabilities"`
}

// DeriveNecessaryCapabilities infers necessary capabilities from declared intent steps.
func DeriveNecessaryCapabilities(intent DeclaredIntent) *NecessaryBehavior {
	caps := make(capability.Set)
	for _, step := range intent.Steps {
		if capObj, ok := mapStepToCapability(step); ok {
			caps[capObj.CapabilityID] = capObj
		}
	}
	return &NecessaryBehavior{
		IntentID:     intent.SkillName,
		Capabilities: caps,
	}
}

func mapStepToCapability(step DeclaredStep) (capability.Capability, bool) {
	switch step.Action {
	case "http.get", "GET":
		return capability.Capability{
			CapabilityID: "nec_http_" + step.StepID,
			Kind:         capability.KindNetHTTP,
			Actions:      []string{"GET"},
			Resource: capability.Resource{
				Scheme: "https",
				Host:   capability.HostSpec{Type: "EXACT", Value: step.Target},
				Path:   "/*",
			},
			Justification: capability.Justification{
				IntentStep: step.StepID,
				RequiredBy: step.Description,
			},
		}, true
	case "fs.read", "READ":
		return capability.Capability{
			CapabilityID: "nec_fs_" + step.StepID,
			Kind:         capability.KindFSRead,
			Actions:      []string{"READ"},
			Resource: capability.Resource{
				Path: step.Target,
			},
			Justification: capability.Justification{
				IntentStep: step.StepID,
				RequiredBy: step.Description,
			},
		}, true
	default:
		return capability.Capability{}, false
	}
}
