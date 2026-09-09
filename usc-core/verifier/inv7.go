package verifier

import (
	"fmt"

	"github.com/Freecode100Year/usc/usc-core/capability"
)

// MonotonicPrivilegeResult represents the result of INV-7 verification.
type MonotonicPrivilegeResult struct {
	Holds       bool                   `json:"holds"`
	Diagnostic  string                 `json:"diagnostic,omitempty"`
	Violations  []capability.Capability `json:"violations,omitempty"`
	Contraction []capability.Capability `json:"contraction,omitempty"`
}

// VerifyINV7 asserts that C_out ⊑ C_in (monotonic privilege reduction).
func VerifyINV7(cIn, cOut capability.Set) MonotonicPrivilegeResult {
	var violations []capability.Capability
	for _, outCap := range cOut {
		subsumed := false
		for _, inCap := range cIn {
			if capability.Subsumes(inCap, outCap) {
				subsumed = true
				break
			}
		}
		if !subsumed {
			violations = append(violations, outCap)
		}
	}
	if len(violations) > 0 {
		diag := fmt.Sprintf("USC-E7007: INV-7 Monotonic Privilege Violation: %d capabilities expanded beyond blueprint", len(violations))
		return MonotonicPrivilegeResult{Holds: false, Diagnostic: diag, Violations: violations}
	}
	contraction := capability.DifferenceSets(cIn, cOut)
	var contractedList []capability.Capability
	for _, c := range contraction {
		contractedList = append(contractedList, c)
	}
	return MonotonicPrivilegeResult{Holds: true, Diagnostic: "INV-7 PASSED", Contraction: contractedList}
}
