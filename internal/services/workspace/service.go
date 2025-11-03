package workspace

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"

	"shanraq.com/internal/auth"
)

// Workspace captures a user's business plans and activity summary.
type Workspace struct {
	ID            uuid.UUID      `json:"id"`
	OwnerID       string         `json:"owner_id"`
	OwnerName     string         `json:"owner_name"`
	OwnerEmail    string         `json:"owner_email"`
	BusinessPlans []BusinessPlan `json:"business_plans"`
	Metrics       Metrics        `json:"metrics"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

// BusinessPlan is a simple descriptor for demo purposes.
type BusinessPlan struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Summary   string    `json:"summary"`
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Metrics contains lightweight analytics for the user's workspace.
type Metrics struct {
	ActivePlans    int `json:"active_plans"`
	CompletedPlans int `json:"completed_plans"`
	Watchers       int `json:"watchers"`
}

// Service defines behaviour for managing user workspaces.
type Service interface {
	GetOrCreate(ctx context.Context, identity auth.Identity) (Workspace, error)
	AddPlan(ctx context.Context, identity auth.Identity, plan BusinessPlan) (Workspace, error)
}

// InMemoryService is a demo implementation backed by an in-memory map.
type InMemoryService struct {
	mu         sync.RWMutex
	workspaces map[string]Workspace
}

// NewInMemoryService instantiates the demo workspace service.
func NewInMemoryService() *InMemoryService {
	return &InMemoryService{workspaces: make(map[string]Workspace)}
}

func (s *InMemoryService) GetOrCreate(_ context.Context, identity auth.Identity) (Workspace, error) {
	key := identity.Subject
	if key == "" {
		key = identity.Email
	}
	if key == "" {
		key = identity.Provider
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	workspace, ok := s.workspaces[key]
	if !ok {
		workspace = Workspace{
			ID:         uuid.New(),
			OwnerID:    key,
			OwnerName:  identity.FullName,
			OwnerEmail: identity.Email,
			BusinessPlans: []BusinessPlan{
				{
					ID:        uuid.New(),
					Title:     "Global Expansion Roadmap",
					Summary:   "Steps for launching operations in three new regions over the next six months.",
					Status:    "active",
					UpdatedAt: time.Now().UTC(),
				},
			},
			Metrics: Metrics{
				ActivePlans:    1,
				CompletedPlans: 0,
				Watchers:       3,
			},
			UpdatedAt: time.Now().UTC(),
		}
		s.workspaces[key] = workspace
	}
	return workspace, nil
}

func (s *InMemoryService) AddPlan(_ context.Context, identity auth.Identity, plan BusinessPlan) (Workspace, error) {
	key := identity.Subject
	if key == "" {
		key = identity.Email
	}
	if key == "" {
		key = identity.Provider
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	workspace, ok := s.workspaces[key]
	if !ok {
		workspace = Workspace{
			ID:         uuid.New(),
			OwnerID:    key,
			OwnerName:  identity.FullName,
			OwnerEmail: identity.Email,
		}
	}

	if plan.ID == uuid.Nil {
		plan.ID = uuid.New()
	}
	if plan.UpdatedAt.IsZero() {
		plan.UpdatedAt = time.Now().UTC()
	}
	if plan.Status == "" {
		plan.Status = "draft"
	}

	workspace.BusinessPlans = append(workspace.BusinessPlans, plan)
	workspace.Metrics.ActivePlans++
	workspace.UpdatedAt = time.Now().UTC()

	s.workspaces[key] = workspace
	return workspace, nil
}

var _ Service = (*InMemoryService)(nil)
