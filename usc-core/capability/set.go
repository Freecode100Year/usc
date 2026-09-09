package capability

// Set represents an immutable or working set of capabilities.
type Set map[string]Capability

// NewSet creates an initialized Capability Set.
func NewSet(caps ...Capability) Set {
	s := make(Set, len(caps))
	for _, c := range caps {
		s[c.CapabilityID] = c
	}
	return s
}

// SubsumesSet checks if parent set subsumes child set (childSet ⊑ parentSet).
func SubsumesSet(parentSet, childSet Set) bool {
	for _, child := range childSet {
		found := false
		for _, parent := range parentSet {
			if Subsumes(parent, child) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// IntersectSets computes the meet-based intersection of two capability sets.
func IntersectSets(s1, s2 Set) Set {
	res := make(Set)
	for _, c1 := range s1 {
		for _, c2 := range s2 {
			if m, ok := Meet(c1, c2); ok {
				res[m.CapabilityID] = *m
			}
		}
	}
	return res
}

// DifferenceSets computes set difference (s1 \ s2) based on subsumption.
func DifferenceSets(s1, s2 Set) Set {
	res := make(Set)
	for id, c1 := range s1 {
		subsumed := false
		for _, c2 := range s2 {
			if Subsumes(c2, c1) {
				subsumed = true
				break
			}
		}
		if !subsumed {
			res[id] = c1
		}
	}
	return res
}
