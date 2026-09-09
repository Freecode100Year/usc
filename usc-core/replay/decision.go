package replay

import (
	"fmt"

	"github.com/Freecode100Year/usc/usc-core/capability"
)

// ReplayStepResult captures the outcome of an individual decision replay step.
type ReplayStepResult struct {
	Index           int    `json:"index"`
	ExpectedVerdict string `json:"expected_verdict"`
	ReplayedVerdict string `json:"replayed_verdict"`
	Matches         bool   `json:"matches"`
}

// ReplayDecisions re-evaluates all recorded events against frozen capability policy.
func ReplayDecisions(tf *TraceFile, frozenPolicy capability.Set) ([]ReplayStepResult, bool) {
	results := make([]ReplayStepResult, 0, len(tf.Events))
	allMatch := true
	for _, evt := range tf.Events {
		res := replaySingleEvent(evt, frozenPolicy)
		if !res.Matches {
			allMatch = false
		}
		results = append(results, res)
	}
	return results, allMatch
}

func replaySingleEvent(evt TraceEvent, policy capability.Set) ReplayStepResult {
	req := capability.Capability{
		Kind:    capability.KindNetHTTP,
		Actions: []string{evt.Action},
		Resource: capability.Resource{
			Host: capability.HostSpec{Type: "EXACT", Value: evt.Target},
		},
	}
	replayedVerdict := "DENY"
	for _, granted := range policy {
		if capability.Subsumes(granted, req) {
			replayedVerdict = "ALLOW"
			break
		}
	}
	matches := evt.Decision == replayedVerdict
	return ReplayStepResult{
		Index:           evt.Index,
		ExpectedVerdict: evt.Decision,
		ReplayedVerdict: replayedVerdict,
		Matches:         matches,
	}
}

// FormatReplaySummary produces human-readable diagnostic output for replay steps.
func FormatReplaySummary(results []ReplayStepResult) string {
	var summary string
	for _, r := range results {
		status := "MATCH"
		if !r.Matches {
			status = "MISMATCH"
		}
		summary += fmt.Sprintf("Step %02d: Expected=%s Replayed=%s [%s]\n", r.Index, r.ExpectedVerdict, r.ReplayedVerdict, status)
	}
	return summary
}
