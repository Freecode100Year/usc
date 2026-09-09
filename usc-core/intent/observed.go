package intent

import "github.com/Freecode100Year/usc/usc-core/capability"

// ObservedFact represents a single factual behavior detected from static/dynamic analysis.
type ObservedFact struct {
	FactID     string                `json:"fact_id"`
	SourceFile string                `json:"source_file"`
	LineNumber int                   `json:"line_number"`
	Capability capability.Capability `json:"capability"`
	RawSnippet string                `json:"raw_snippet,omitempty"`
}

// ObservedBehavior encapsulates the observed behavior set B_obs.
type ObservedBehavior struct {
	SkillName    string         `json:"skill_name"`
	Facts        []ObservedFact `json:"facts"`
	Capabilities capability.Set `json:"capabilities"`
}

// NewObservedBehavior initializes an ObservedBehavior instance.
func NewObservedBehavior(skillName string, facts []ObservedFact) *ObservedBehavior {
	caps := make(capability.Set, len(facts))
	for _, f := range facts {
		caps[f.Capability.CapabilityID] = f.Capability
	}
	return &ObservedBehavior{
		SkillName:    skillName,
		Facts:        facts,
		Capabilities: caps,
	}
}
