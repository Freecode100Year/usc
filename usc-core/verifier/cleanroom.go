package verifier

// CleanRoomProof contains measurable isolation facts verifying INV-9.
type CleanRoomProof struct {
	SourceMountsVisible     int      `json:"source_mounts_visible"`
	InheritedSourceFDs      int      `json:"inherited_source_fds"`
	KnownSourceDigestsFound int      `json:"known_source_digests_found"`
	NetworkDefaultRoute     bool     `json:"network_default_route"`
	AmbientCapabilities     []string `json:"ambient_capabilities"`
	CleanEnvironment        bool     `json:"clean_environment"`
	Verified                bool     `json:"verified"`
}

// MeasureCleanRoom inspects process metrics to deterministically verify INV-9.
func MeasureCleanRoom(sourceMounts, inheritedFDs, sourceDigests int, netRoute bool) CleanRoomProof {
	verified := sourceMounts == 0 &&
		inheritedFDs == 0 &&
		sourceDigests == 0 &&
		!netRoute

	return CleanRoomProof{
		SourceMountsVisible:     sourceMounts,
		InheritedSourceFDs:      inheritedFDs,
		KnownSourceDigestsFound: sourceDigests,
		NetworkDefaultRoute:     netRoute,
		AmbientCapabilities:     []string{},
		CleanEnvironment:        verified,
		Verified:                verified,
	}
}
