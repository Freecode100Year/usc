package intent

import (
	"github.com/Freecode100Year/usc/usc-core/capability"
)

// Verdict defines the security verdict for intent decontamination.
type Verdict string

const (
	VerdictPass      Verdict = "PASS"
	VerdictHardBlock Verdict = "HARD_BLOCK"
	VerdictStrip     Verdict = "STRIP"
)

// DecontaminationReport summarizes the decontamination pass.
type DecontaminationReport struct {
	ExtraBehavior   capability.Set `json:"extra_behavior"`   // B_extra
	MissingBehavior capability.Set `json:"missing_behavior"` // B_missing
	DivergenceScore float64        `json:"divergence_score"` // D_N
	Verdict         Verdict        `json:"verdict"`
	Findings        []string       `json:"findings"`
}

// Decontaminate performs three-way intent decontamination analysis.
func Decontaminate(obs *ObservedBehavior, nec *NecessaryBehavior, criticalTau float64) *DecontaminationReport {
	extra := capability.DifferenceSets(obs.Capabilities, nec.Capabilities)
	missing := capability.DifferenceSets(nec.Capabilities, obs.Capabilities)
	dn, findings := evaluateDivergence(extra)

	verdict := VerdictPass
	if dn > criticalTau {
		verdict = VerdictHardBlock
	} else if len(extra) > 0 {
		verdict = VerdictStrip
	}

	return &DecontaminationReport{
		ExtraBehavior:   extra,
		MissingBehavior: missing,
		DivergenceScore: dn,
		Verdict:         verdict,
		Findings:        findings,
	}
}

func evaluateDivergence(extra capability.Set) (float64, []string) {
	dn := 0.0
	var findings []string
	weights := capability.DefaultWeights()
	for _, c := range extra {
		cv := capability.ComputeCostVector(c)
		weight := cv.NetworkBreadth*weights.Wn + cv.FilesystemBreadth*weights.Wf +
			cv.SecretExposure*weights.Ws + cv.ExecAuthority*weights.We
		dn += weight
		findings = append(findings, "USC-INTENT-DECONTAM: extra unverified capability "+string(c.Kind))
	}
	return dn, findings
}
