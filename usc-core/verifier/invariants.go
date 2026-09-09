package verifier

// InvariantID identifies each of the 9 Zero-Trust Invariants.
type InvariantID string

const (
	INV1 InvariantID = "INV-1: Origin Neutrality"
	INV2 InvariantID = "INV-2: Zero Grant Baseline"
	INV3 InvariantID = "INV-3: Authority Severance"
	INV4 InvariantID = "INV-4: Candidate Distrust"
	INV5 InvariantID = "INV-5: Secrets as Capabilities"
	INV6 InvariantID = "INV-6: Strict Provenance"
	INV7 InvariantID = "INV-7: Monotonic Privilege Reduction"
	INV8 InvariantID = "INV-8: Deterministic Confinement"
	INV9 InvariantID = "INV-9: Physical Clean-Room Barrier"
)

// InvariantStatus reflects the verification status of an invariant.
type InvariantStatus struct {
	ID        InvariantID `json:"invariant_id"`
	Passed    bool        `json:"passed"`
	Violation string      `json:"violation,omitempty"`
}

// InvariantSuite aggregates verification of all 9 invariants.
type InvariantSuite struct {
	Results []InvariantStatus `json:"results"`
}

// AllPassed evaluates if all invariants in the suite hold strictly.
func (s *InvariantSuite) AllPassed() bool {
	for _, r := range s.Results {
		if !r.Passed {
			return false
		}
	}
	return true
}
