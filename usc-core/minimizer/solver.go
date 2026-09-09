package minimizer

import (
	"github.com/Freecode100Year/usc/usc-core/capability"
	"github.com/Freecode100Year/usc/usc-core/intent"
)

// Workflow describes the execution steps and ordering.
type Workflow struct {
	WorkflowID string   `json:"workflow_id"`
	StepOrder  []string `json:"step_order"`
}

// RuntimeSpec describes the target runtime constraints and features.
type RuntimeSpec struct {
	TargetPlatform string   `json:"target_platform"` // hermes, openclaw, langgraph
	SupportedKinds []string `json:"supported_kinds"`
}

// Policy defines enterprise security bounds and allowed capabilities.
type Policy struct {
	PolicyID            string         `json:"policy_id"`
	AllowedCapabilities capability.Set `json:"allowed_capabilities"`
	MaxDivergence       float64        `json:"max_divergence"`
}

// Proof records the mathematical and logical justification for the minimal set.
type Proof struct {
	Satisfied bool    `json:"satisfied"`
	FinalCost float64 `json:"final_cost"`
	Reason    string  `json:"reason"`
}

// Solver abstracts capability minimization.
type Solver interface {
	Solve(
		intent intent.DeclaredIntent,
		wf Workflow,
		target RuntimeSpec,
		policy Policy,
	) (capability.Set, Proof, error)
}
