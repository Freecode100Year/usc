package attestation

// ArtifactInfo identifies the compiled target artifact.
type ArtifactInfo struct {
	Name           string `json:"name"`
	SHA256         string `json:"sha256"`
	TargetPlatform string `json:"target_platform"`
}

// Claims encapsulates the zero-trust claims verified during compilation.
type Claims struct {
	SourceTaintQuarantined             bool    `json:"source_taint_quarantined"`
	CapabilityCountReduction           float64 `json:"capability_count_reduction"`
	AttackSurfaceReduction             float64 `json:"attack_surface_reduction"`
	CleanroomIsolated                  bool    `json:"cleanroom_isolated"`
	GeneratedAuthorityWithinBlueprint  bool    `json:"generated_authority_within_blueprint"`
	DryRunValidation                   string  `json:"dry_run_validation"`
}

// ProofReferences links all cryptographic digests forming the proof bundle.
type ProofReferences struct {
	AuditChainRoot         string `json:"audit_chain_root"`
	CanonicalIntentDigest  string `json:"canonical_intent_digest"`
	BlueprintDigest        string `json:"blueprint_digest"`
	PolicyDigest           string `json:"policy_digest"`
	CleanroomProofDigest   string `json:"cleanroom_proof_digest"`
}

// Signature records cryptographic verification from an authority key.
type Signature struct {
	KeyID     string `json:"key_id"`
	Algorithm string `json:"algorithm"` // ed25519
	Sig       string `json:"sig"`
}

// Attestation represents the machine-verifiable proof document.
type Attestation struct {
	Schema          string          `json:"$schema"`
	Version         string          `json:"version"`
	Artifact        ArtifactInfo    `json:"artifact"`
	Claims          Claims          `json:"claims"`
	ProofReferences ProofReferences `json:"proof_references"`
	Signature       Signature       `json:"signature"`
}
