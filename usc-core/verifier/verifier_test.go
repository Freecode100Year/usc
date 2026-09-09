package verifier

import (
	"testing"

	"github.com/Freecode100Year/usc/usc-core/capability"
)

func TestVerifyINV7(t *testing.T) {
	bp := capability.NewSet(capability.Capability{
		CapabilityID: "cap1", Kind: capability.KindNetHTTP, Actions: []string{"GET"},
		Resource: capability.Resource{Host: capability.HostSpec{Type: "EXACT", Value: "api.github.com"}},
	})
	genSafe := capability.NewSet(capability.Capability{
		CapabilityID: "cap1", Kind: capability.KindNetHTTP, Actions: []string{"GET"},
		Resource: capability.Resource{Host: capability.HostSpec{Type: "EXACT", Value: "api.github.com"}},
	})
	genWild := capability.NewSet(capability.Capability{
		CapabilityID: "cap_wild", Kind: capability.KindNetHTTP, Actions: []string{"GET", "POST"},
		Resource: capability.Resource{Host: capability.HostSpec{Type: "WILDCARD", Value: "*"}},
	})
	if r1 := VerifyINV7(bp, genSafe); !r1.Holds {
		t.Fatalf("expected safe generation to pass INV-7")
	}
	if r2 := VerifyINV7(bp, genWild); r2.Holds {
		t.Fatalf("expected wildcard generation to trigger HARD_BLOCK on INV-7")
	}
}

func TestCleanRoomAndProofObligations(t *testing.T) {
	proof := MeasureCleanRoom(0, 0, 0, false)
	if !proof.Verified {
		t.Fatalf("clean room measurement failed")
	}
	po := ProofObligations{
		PO1_GeneratedWithinBlueprint: true,
		PO2_CleanRoomZeroSourceBytes: true,
		PO3_ZeroSecretMaterialInSkill: true,
		PO4_RuntimeEffectsConfinement: true,
		PO5_HashChainIntegrity:       true,
	}
	if !po.EvaluateProofObligations() {
		t.Fatalf("expected all proof obligations to be satisfied")
	}
}
