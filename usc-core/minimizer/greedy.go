package minimizer

import (
	"github.com/Freecode100Year/usc/usc-core/capability"
	"github.com/Freecode100Year/usc/usc-core/intent"
)

// GreedySolver implements greedy minimal capability reduction.
type GreedySolver struct{}

// NewGreedySolver creates a GreedySolver instance.
func NewGreedySolver() *GreedySolver {
	return &GreedySolver{}
}

// Solve performs greedy minimization on necessary intent capabilities.
func (g *GreedySolver) Solve(
	in intent.DeclaredIntent,
	wf Workflow,
	target RuntimeSpec,
	policy Policy,
) (capability.Set, Proof, error) {
	nec := intent.DeriveNecessaryCapabilities(in)
	minimized := make(capability.Set)

	for id, capItem := range nec.Capabilities {
		tightCap := tightenCapability(capItem)
		minimized[id] = tightCap
	}

	if err := ValidateAgainstPolicy(minimized, policy); err != nil {
		return nil, Proof{Satisfied: false, Reason: err.Error()}, err
	}
	if err := ValidateAgainstRuntime(minimized, target); err != nil {
		return nil, Proof{Satisfied: false, Reason: err.Error()}, err
	}

	cost := capability.TotalCost(minimized, capability.DefaultWeights())
	return minimized, Proof{Satisfied: true, FinalCost: cost, Reason: "Greedy minimal bound verified"}, nil
}

func tightenCapability(c capability.Capability) capability.Capability {
	tight := c
	if tight.Constraints.MaxResponseBytes == 0 {
		tight.Constraints.MaxResponseBytes = 1048576 // default 1MB cap
	}
	if tight.Constraints.RateLimitPerMinute == 0 {
		tight.Constraints.RateLimitPerMinute = 60
	}
	tight.Constraints.AllowRedirect = false
	return tight
}
