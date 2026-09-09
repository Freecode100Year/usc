package capability

import "strings"

// Meet computes the greatest lower bound (a ⊓ b) of two capabilities.
func Meet(a, b Capability) (*Capability, bool) {
	if a.Kind != b.Kind {
		return nil, false
	}
	commonActions := intersectActions(a.Actions, b.Actions)
	if len(commonActions) == 0 {
		return nil, false
	}
	res, ok := intersectResource(a.Resource, b.Resource)
	if !ok {
		return nil, false
	}
	return &Capability{
		CapabilityID:  a.CapabilityID + "_meet_" + b.CapabilityID,
		Kind:          a.Kind,
		Actions:       commonActions,
		Resource:      res,
		Constraints:   meetConstraints(a.Constraints, b.Constraints),
		Justification: a.Justification,
	}, true
}

func intersectActions(a, b []string) []string {
	bm := make(map[string]struct{}, len(b))
	for _, act := range b {
		bm[strings.ToUpper(strings.TrimSpace(act))] = struct{}{}
	}
	var res []string
	for _, act := range a {
		norm := strings.ToUpper(strings.TrimSpace(act))
		if _, exists := bm[norm]; exists {
			res = append(res, norm)
		}
	}
	return res
}

func intersectResource(a, b Resource) (Resource, bool) {
	if a.Scheme != b.Scheme && a.Scheme != "" && b.Scheme != "" {
		return Resource{}, false
	}
	scheme := a.Scheme
	if scheme == "" {
		scheme = b.Scheme
	}
	host := a.Host
	if a.Host.Value != b.Host.Value {
		if a.Host.Type == "WILDCARD" && HostSubsumes(a.Host, b.Host) {
			host = b.Host
		} else if b.Host.Type == "WILDCARD" && HostSubsumes(b.Host, a.Host) {
			host = a.Host
		} else {
			return Resource{}, false
		}
	}
	return Resource{Scheme: scheme, Host: host, Path: a.Path}, true
}

func meetConstraints(a, b Constraints) Constraints {
	c := Constraints{
		AllowRedirect:      a.AllowRedirect && b.AllowRedirect,
		MaxResponseBytes:   a.MaxResponseBytes,
		RateLimitPerMinute: a.RateLimitPerMinute,
	}
	if b.MaxResponseBytes > 0 && (c.MaxResponseBytes == 0 || b.MaxResponseBytes < c.MaxResponseBytes) {
		c.MaxResponseBytes = b.MaxResponseBytes
	}
	if b.RateLimitPerMinute > 0 && (c.RateLimitPerMinute == 0 || b.RateLimitPerMinute < c.RateLimitPerMinute) {
		c.RateLimitPerMinute = b.RateLimitPerMinute
	}
	return c
}
