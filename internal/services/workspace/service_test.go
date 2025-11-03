package workspace

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"shanraq.com/internal/auth"
)

func TestInMemoryServiceGetOrCreateSeedsWorkspace(t *testing.T) {
	service := NewInMemoryService()
	identity := auth.Identity{
		Subject:  "user-42",
		FullName: "Future Founder",
		Email:    "founder@example.com",
	}

	ws, err := service.GetOrCreate(context.Background(), identity)
	if err != nil {
		t.Fatalf("GetOrCreate() error = %v", err)
	}

	if ws.OwnerID == "" {
		t.Error("OwnerID is empty")
	}
	if ws.OwnerName != identity.FullName {
		t.Errorf("OwnerName = %q, want %q", ws.OwnerName, identity.FullName)
	}
	if len(ws.BusinessPlans) == 0 {
		t.Fatal("expected seeded BusinessPlans, got none")
	}
	if ws.Metrics.ActivePlans != len(ws.BusinessPlans) {
		t.Errorf("ActivePlans = %d, want %d", ws.Metrics.ActivePlans, len(ws.BusinessPlans))
	}
	if ws.UpdatedAt.IsZero() {
		t.Error("UpdatedAt is zero")
	}
}

func TestInMemoryServiceGetOrCreateReturnsExistingWorkspace(t *testing.T) {
	service := NewInMemoryService()
	identity := auth.Identity{Subject: "repeat-user"}

	first, err := service.GetOrCreate(context.Background(), identity)
	if err != nil {
		t.Fatalf("GetOrCreate() first call error = %v", err)
	}

	time.Sleep(1 * time.Millisecond)

	second, err := service.GetOrCreate(context.Background(), identity)
	if err != nil {
		t.Fatalf("GetOrCreate() second call error = %v", err)
	}

	if first.ID != second.ID {
		t.Fatalf("workspace ID changed between calls: %s vs %s", first.ID, second.ID)
	}
	if len(second.BusinessPlans) != len(first.BusinessPlans) {
		t.Fatalf("BusinessPlans length changed: %d vs %d", len(second.BusinessPlans), len(first.BusinessPlans))
	}
	if !second.UpdatedAt.After(first.UpdatedAt) && !second.UpdatedAt.Equal(first.UpdatedAt) {
		t.Fatalf("UpdatedAt unexpectedly older; first=%v second=%v", first.UpdatedAt, second.UpdatedAt)
	}
}

func TestInMemoryServiceAddPlanAppendsAndUpdatesMetrics(t *testing.T) {
	service := NewInMemoryService()
	identity := auth.Identity{Subject: "planner"}

	ws, err := service.GetOrCreate(context.Background(), identity)
	if err != nil {
		t.Fatalf("GetOrCreate() error = %v", err)
	}
	initialPlans := len(ws.BusinessPlans)

	input := BusinessPlan{
		Title:   "MENA Expansion",
		Summary: "Launch Riyadh brokerage pod and logistics lane.",
		Status:  "active",
	}
	updated, err := service.AddPlan(context.Background(), identity, input)
	if err != nil {
		t.Fatalf("AddPlan() error = %v", err)
	}

	if len(updated.BusinessPlans) != initialPlans+1 {
		t.Fatalf("len(BusinessPlans) = %d, want %d", len(updated.BusinessPlans), initialPlans+1)
	}

	lastPlan := updated.BusinessPlans[len(updated.BusinessPlans)-1]
	if lastPlan.Title != input.Title {
		t.Errorf("lastPlan.Title = %q, want %q", lastPlan.Title, input.Title)
	}
	if lastPlan.ID == uuid.Nil {
		t.Error("expected plan ID to be generated")
	}
	if updated.Metrics.ActivePlans != ws.Metrics.ActivePlans+1 {
		t.Fatalf("ActivePlans = %d, want %d", updated.Metrics.ActivePlans, ws.Metrics.ActivePlans+1)
	}
	if updated.UpdatedAt.Before(ws.UpdatedAt) {
		t.Fatalf("UpdatedAt not refreshed; before=%v after=%v", ws.UpdatedAt, updated.UpdatedAt)
	}
}
