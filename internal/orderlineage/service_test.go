package orderlineage

import (
	"path/filepath"
	"testing"
)

func TestRecordAndResolveFollowRolloverChain(t *testing.T) {
	t.Parallel()

	service := NewService(filepath.Join(t.TempDir(), "lineage.json"))
	if err := service.Record("2026-03-13/1", Entry{CurrentOrderID: "2026-03-13/2", Kind: "amend"}); err != nil {
		t.Fatalf("Record returned error: %v", err)
	}
	if err := service.Record("2026-03-13/2", Entry{CurrentOrderID: "2026-03-13/3", Kind: "cancel"}); err != nil {
		t.Fatalf("Record returned error: %v", err)
	}

	resolved, ok, err := service.Resolve("2026-03-13/1")
	if err != nil {
		t.Fatalf("Resolve returned error: %v", err)
	}
	if !ok {
		t.Fatal("expected lineage resolution")
	}
	if resolved != "2026-03-13/3" {
		t.Fatalf("expected final order id 2026-03-13/3, got %q", resolved)
	}
}

func TestResolveReturnsNoAliasWhenUnchanged(t *testing.T) {
	t.Parallel()

	service := NewService(filepath.Join(t.TempDir(), "lineage.json"))
	resolved, ok, err := service.Resolve("2026-03-13/1")
	if err != nil {
		t.Fatalf("Resolve returned error: %v", err)
	}
	if ok {
		t.Fatalf("expected no alias, got %q", resolved)
	}
}

func TestRecordKeepsUnresolvedRecoveryHint(t *testing.T) {
	t.Parallel()

	service := NewService(filepath.Join(t.TempDir(), "lineage.json"))
	if err := service.Record("2026-03-13/3", Entry{
		Kind:      "cancel",
		Symbol:    "TSLL",
		Market:    "us",
		Quantity:  1,
		Price:     500,
		OrderDate: "2026-03-13",
	}); err != nil {
		t.Fatalf("Record returned error: %v", err)
	}

	entry, ok, err := service.Lookup("2026-03-13/3")
	if err != nil {
		t.Fatalf("Lookup returned error: %v", err)
	}
	if !ok {
		t.Fatal("expected lookup hit")
	}
	if entry.CurrentOrderID != "" {
		t.Fatalf("expected unresolved current id, got %q", entry.CurrentOrderID)
	}
	if entry.Kind != "cancel" || entry.Symbol != "TSLL" || entry.Market != "us" {
		t.Fatalf("unexpected entry metadata: %#v", entry)
	}
	if entry.OrderDate != "2026-03-13" || entry.Quantity != 1 || entry.Price != 500 {
		t.Fatalf("unexpected entry match hint: %#v", entry)
	}
	if entry.UpdatedAt.IsZero() {
		t.Fatal("expected updated_at to be recorded")
	}
}
