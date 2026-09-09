package capability

import "strings"

// Join computes the least upper bound (a ⊔ b) for normalization/analysis.
func Join(a, b Capability) Capability {
	res := Capability{
		CapabilityID:  a.CapabilityID + "_join_" + b.CapabilityID,
		Kind:          a.Kind,
		Actions:       unionActions(a.Actions, b.Actions),
		Resource:      joinResource(a.Resource, b.Resource),
		Constraints:   joinConstraints(a.Constraints, b.Constraints),
		Justification: a.Justification,
	}
	return res
}

func unionActions(a, b []string) []string {
	seen := make(map[string]struct{})
	var res []string
	for _, act := range append(a, b...) {
		norm := strings.ToUpper(strings.TrimSpace(act))
		if _, ok := seen[norm]; !ok {
			seen[norm] = struct{}{}
			res = append(res, norm)
		}
	}
	return res
}

func joinResource(a, b Resource) Resource {
	r := a
	if a.Host.Value != b.Host.Value {
		r.Host = HostSpec{Type: "WILDCARD", Value: "*"}
	}
	if a.Path != b.Path {
		r.Path = "/*"
	}
	return r
}

func joinConstraints(a, b Constraints) Constraints {
	c := Constraints{
		AllowRedirect:      a.AllowRedirect || b.AllowRedirect,
		MaxResponseBytes:   a.MaxResponseBytes,
		RateLimitPerMinute: a.RateLimitPerMinute,
	}
	if b.MaxResponseBytes > c.MaxResponseBytes {
		c.MaxResponseBytes = b.MaxResponseBytes
	}
	if b.RateLimitPerMinute > c.RateLimitPerMinute {
		c.RateLimitPerMinute = b.RateLimitPerMinute
	}
	return c
}
