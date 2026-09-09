package verifier

// ProofObligations encapsulates the 5 mandatory formal proof obligations.
type ProofObligations struct {
	PO1_GeneratedWithinBlueprint bool `json:"po1_generated_within_blueprint"`
	PO2_CleanRoomZeroSourceBytes bool `json:"po2_clean_room_zero_source_bytes"`
	PO3_ZeroSecretMaterialInSkill bool `json:"po3_zero_secret_material_in_skill"`
	PO4_RuntimeEffectsConfinement bool `json:"po4_runtime_effects_confinement"`
	PO5_HashChainIntegrity       bool `json:"po5_hash_chain_integrity"`
}

// EvaluateProofObligations checks whether all proof obligations are satisfied.
func (po *ProofObligations) EvaluateProofObligations() bool {
	return po.PO1_GeneratedWithinBlueprint &&
		po.PO2_CleanRoomZeroSourceBytes &&
		po.PO3_ZeroSecretMaterialInSkill &&
		po.PO4_RuntimeEffectsConfinement &&
		po.PO5_HashChainIntegrity
}
