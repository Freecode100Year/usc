package capability

// CostVector defines the multi-dimensional attack surface cost vector.
type CostVector struct {
	NetworkBreadth    float64 `json:"network_breadth"`    // Rn
	FilesystemBreadth float64 `json:"filesystem_breadth"` // Rf
	SecretExposure    float64 `json:"secret_exposure"`    // Rs
	WriteAuthority    float64 `json:"write_authority"`    // Rw
	ExecAuthority     float64 `json:"exec_authority"`     // Re
	PrivilegeLevel    float64 `json:"privilege_level"`    // Rp
	Persistence       float64 `json:"persistence"`        // Rpersist
}

// Weights defines the policy weighting vector.
type Weights struct {
	Wn float64
	Wf float64
	Ws float64
	Ww float64
	We float64
	Wp float64
	Wpersist float64
}

// DefaultWeights provides standard Enterprise Strict policy weights.
func DefaultWeights() Weights {
	return Weights{
		Ws: 100.0, Wp: 100.0, We: 50.0,
		Ww: 30.0, Wn: 20.0, Wf: 10.0, Wpersist: 20.0,
	}
}

// ComputeCostVector calculates the multi-dimensional cost vector for a capability.
func ComputeCostVector(c Capability) CostVector {
	v := CostVector{}
	switch c.Kind {
	case KindNetHTTP, KindNetEgress:
		v.NetworkBreadth = evalNetBreadth(c)
	case KindFSRead:
		v.FilesystemBreadth = 10.0
	case KindFSWrite:
		v.FilesystemBreadth = 10.0
		v.WriteAuthority = 25.0
	case KindSysExec:
		v.ExecAuthority = 50.0
		v.PrivilegeLevel = 30.0
	case KindSecretAccess:
		v.SecretExposure = 100.0
	}
	return v
}

func evalNetBreadth(c Capability) float64 {
	if c.Resource.Host.Type == "WILDCARD" && c.Resource.Host.Value == "*" {
		return 50.0
	}
	return 5.0
}

// TotalCost computes weighted scalar cost for a set of capabilities.
func TotalCost(s Set, w Weights) float64 {
	total := 0.0
	for _, c := range s {
		cv := ComputeCostVector(c)
		total += cv.NetworkBreadth*w.Wn + cv.FilesystemBreadth*w.Wf +
			cv.SecretExposure*w.Ws + cv.WriteAuthority*w.Ww +
			cv.ExecAuthority*w.We + cv.PrivilegeLevel*w.Wp +
			cv.Persistence*w.Wpersist
	}
	return total
}

// ComputeReduction computes CCR and ASR metrics.
func ComputeReduction(observed, granted Set, w Weights) (ccr, asr float64) {
	obsCount := float64(len(observed))
	grantCount := float64(len(granted))
	if obsCount > 0 {
		ccr = 1.0 - (grantCount / obsCount)
	}
	obsCost := TotalCost(observed, w)
	grantCost := TotalCost(granted, w)
	if obsCost > 0 {
		asr = 1.0 - (grantCost / obsCost)
	}
	return ccr, asr
}
