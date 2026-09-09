package audit

// Stage defines the 7 observable compiler stages.
type Stage string

const (
	StageIngest       Stage = "INGEST"
	StageDecontam     Stage = "DECONTAMINATE"
	StageMinimize     Stage = "MINIMIZE"
	StageCleanRebuild Stage = "CLEAN_REBUILD"
	StageReAudit      Stage = "RE_AUDIT"
	StageSandbox      Stage = "SANDBOX"
	StageAttest       Stage = "ATTEST"
)

// Decision defines the audit verdict for an event.
type Decision string

const (
	DecisionPass      Decision = "PASS"
	DecisionHardBlock Decision = "HARD_BLOCK"
	DecisionWarn      Decision = "WARN"
	DecisionDrop      Decision = "DROP"
)

// Event represents a lossless audit record in the hash chain.
type Event struct {
	Seq            uint64   `json:"seq"`
	PreviousDigest string   `json:"previous_digest"`
	Stage          Stage    `json:"stage"`
	Rule           string   `json:"rule"`
	Decision       Decision `json:"decision"`
	Payload        string   `json:"payload"`
	Timestamp      int64    `json:"timestamp"`
}
