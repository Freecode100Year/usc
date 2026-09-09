package audit

import "testing"

func TestAuditLedgerAndHashChain(t *testing.T) {
	genesis := ComputeGenesisDigest("src_hash", "bp_hash", "policy_hash", "usc_v0.1")
	ledger := NewLedger(genesis)

	_, _, err := ledger.Append(StageIngest, "RULE_SOURCE_TAINT", DecisionPass, "ingested 3 files")
	if err != nil {
		t.Fatalf("append failed: %v", err)
	}
	_, _, err = ledger.Append(StageDecontam, "RULE_EXTRA_BEHAVIOR", DecisionDrop, "stripped unauthorized telemetry")
	if err != nil {
		t.Fatalf("append failed: %v", err)
	}

	valid, err := VerifyChain(genesis, ledger.Events(), ledger.LatestDigest())
	if err != nil || !valid {
		t.Fatalf("audit chain verification failed: %v", err)
	}
}

func TestCanonicalize(t *testing.T) {
	raw := []byte(`{"z":1,"a":2,"m":{"b":3,"a":4}}`)
	canon, err := Canonicalize(raw)
	if err != nil {
		t.Fatalf("canonicalize failed: %v", err)
	}
	expected := `{"a":2,"m":{"a":4,"b":3},"z":1}`
	if string(canon) != expected {
		t.Fatalf("canonicalization mismatch, expected %s, got %s", expected, string(canon))
	}
}
