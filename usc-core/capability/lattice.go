package capability

import "reflect"

// Subsumes checks if parent subsumes child (child ⊑ parent).
func Subsumes(parent, child Capability) bool {
	if parent.Kind != child.Kind {
		return false
	}
	if !ActionsSubsumes(parent.Actions, child.Actions) {
		return false
	}
	if !ResourceSubsumes(parent.Resource, child.Resource) {
		return false
	}
	return checkConstraints(parent.Constraints, child.Constraints)
}

func checkConstraints(parent, child Constraints) bool {
	if !parent.AllowRedirect && child.AllowRedirect {
		return false
	}
	if parent.MaxResponseBytes > 0 && child.MaxResponseBytes > parent.MaxResponseBytes {
		return false
	}
	if parent.RateLimitPerMinute > 0 && child.RateLimitPerMinute > parent.RateLimitPerMinute {
		return false
	}
	return true
}

// VerifyReflexive checks reflexivity: a ⊑ a.
func VerifyReflexive(a Capability) bool {
	return Subsumes(a, a)
}

// VerifyAntisymmetric checks antisymmetry: a ⊑ b ∧ b ⊑ a => a == b.
func VerifyAntisymmetric(a, b Capability) bool {
	if Subsumes(a, b) && Subsumes(b, a) {
		return reflect.DeepEqual(a.Kind, b.Kind) &&
			reflect.DeepEqual(a.Resource, b.Resource)
	}
	return true
}

// VerifyTransitive checks transitivity: a ⊑ b ∧ b ⊑ c => a ⊑ c.
func VerifyTransitive(a, b, c Capability) bool {
	if Subsumes(b, a) && Subsumes(c, b) {
		return Subsumes(c, a)
	}
	return true
}
