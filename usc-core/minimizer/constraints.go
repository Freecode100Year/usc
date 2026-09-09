package minimizer

import (
	"errors"

	"github.com/Freecode100Year/usc/usc-core/capability"
)

var (
	ErrPolicyViolation    = errors.New("USC-E7002: capability exceeds policy allowance")
	ErrUnsupportedRuntime = errors.New("USC-E7003: capability kind unsupported by target runtime")
)

// ValidateAgainstPolicy verifies if candidate set is subsumed by allowed policy.
func ValidateAgainstPolicy(candidate capability.Set, policy Policy) error {
	if len(policy.AllowedCapabilities) == 0 {
		return nil
	}
	if !capability.SubsumesSet(policy.AllowedCapabilities, candidate) {
		return ErrPolicyViolation
	}
	return nil
}

// ValidateAgainstRuntime verifies if candidate kinds are supported by runtime.
func ValidateAgainstRuntime(candidate capability.Set, target RuntimeSpec) error {
	if len(target.SupportedKinds) == 0 {
		return nil
	}
	supp := make(map[string]struct{}, len(target.SupportedKinds))
	for _, k := range target.SupportedKinds {
		supp[k] = struct{}{}
	}
	for _, c := range candidate {
		if _, ok := supp[string(c.Kind)]; !ok {
			return ErrUnsupportedRuntime
		}
	}
	return nil
}
