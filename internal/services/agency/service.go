package agency

import (
	"context"
	"sort"
	"strings"
	"sync"

	"github.com/google/uuid"
)

// Agency represents a real estate agency or brokerage.
type Agency struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Tagline    string    `json:"tagline"`
	Country    string    `json:"country"`
	Website    string    `json:"website"`
	LogoURL    string    `json:"logo_url"`
	HeadOffice string    `json:"head_office"`
}

// Realtor represents an individual agent affiliated with an agency.
type Realtor struct {
	ID         uuid.UUID `json:"id"`
	AgencyID   uuid.UUID `json:"agency_id"`
	AgencyName string    `json:"agency_name"`
	FullName   string    `json:"full_name"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	Languages  []string  `json:"languages"`
	Region     string    `json:"region"`
	PhotoURL   string    `json:"photo_url"`
}

// Service exposes agency and realtor data.
type Service interface {
	ListAgencies(ctx context.Context) ([]Agency, error)
	Featured(ctx context.Context, limit int) ([]Agency, error)
	ListRealtors(ctx context.Context) ([]Realtor, error)
	FeaturedRealtors(ctx context.Context, limit int) ([]Realtor, error)
}

// InMemoryService provides seeded demo data.
type InMemoryService struct {
	mu       sync.RWMutex
	agencies []Agency
	realtors []Realtor
}

// NewInMemoryService seeds demo agencies and realtors.
func NewInMemoryService() *InMemoryService {
	service := &InMemoryService{}
	service.seed()
	return service
}

func (s *InMemoryService) ListAgencies(_ context.Context) ([]Agency, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	agencies := make([]Agency, len(s.agencies))
	copy(agencies, s.agencies)
	sort.Slice(agencies, func(i, j int) bool {
		return agencies[i].Name < agencies[j].Name
	})
	return agencies, nil
}

func (s *InMemoryService) Featured(ctx context.Context, limit int) ([]Agency, error) {
	agencies, err := s.ListAgencies(ctx)
	if err != nil {
		return nil, err
	}
	if limit <= 0 || limit > len(agencies) {
		limit = len(agencies)
	}
	return agencies[:limit], nil
}

func (s *InMemoryService) ListRealtors(_ context.Context) ([]Realtor, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	realtors := make([]Realtor, len(s.realtors))
	copy(realtors, s.realtors)
	sort.Slice(realtors, func(i, j int) bool {
		return realtors[i].FullName < realtors[j].FullName
	})
	return realtors, nil
}

func (s *InMemoryService) FeaturedRealtors(ctx context.Context, limit int) ([]Realtor, error) {
	realtors, err := s.ListRealtors(ctx)
	if err != nil {
		return nil, err
	}
	if limit <= 0 || limit > len(realtors) {
		limit = len(realtors)
	}
	return realtors[:limit], nil
}

func (s *InMemoryService) seed() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.agencies) > 0 {
		return
	}

	s.agencies = []Agency{
		{
			ID:         uuid.New(),
			Name:       "Shanraq Global Realty",
			Tagline:    "Luxury Estates Across Continents",
			Country:    "AE",
			Website:    "https://shanraq.com/agency/global",
			LogoURL:    "/static/brand/logo_light.svg",
			HeadOffice: "Dubai, UAE",
		},
		{
			ID:         uuid.New(),
			Name:       "Nordic Skyline Partners",
			Tagline:    "Scandinavian waterfront & alpine living",
			Country:    "SE",
			Website:    "https://shanraq.com/agency/nordic-skyline",
			LogoURL:    "",
			HeadOffice: "Stockholm, Sweden",
		},
		{
			ID:         uuid.New(),
			Name:       "Pacifica Urban Advisors",
			Tagline:    "Smart investments across the Pacific Rim",
			Country:    "US",
			Website:    "https://shanraq.com/agency/pacifica-urban",
			LogoURL:    "",
			HeadOffice: "San Francisco, USA",
		},
		{
			ID:         uuid.New(),
			Name:       "Atlas Heritage Homes",
			Tagline:    "Historic residences and cultural landmarks",
			Country:    "IT",
			Website:    "https://shanraq.com/agency/atlas-heritage",
			LogoURL:    "",
			HeadOffice: "Florence, Italy",
		},
	}

	s.realtors = []Realtor{
		{
			ID:         uuid.New(),
			AgencyID:   s.agencies[0].ID,
			AgencyName: s.agencies[0].Name,
			FullName:   "Layla Al-Mansouri",
			Email:      "layla@shanraq.com",
			Phone:      "+971-4-555-0147",
			Languages:  []string{"Arabic", "English", "Hindi"},
			Region:     "Middle East & North Africa",
			PhotoURL:   "",
		},
		{
			ID:         uuid.New(),
			AgencyID:   s.agencies[1].ID,
			AgencyName: s.agencies[1].Name,
			FullName:   "Karl Johansson",
			Email:      "karl@nordicskyline.com",
			Phone:      "+46-8-555-0199",
			Languages:  []string{"Swedish", "Norwegian", "English"},
			Region:     "Nordics & Northern Europe",
			PhotoURL:   "",
		},
		{
			ID:         uuid.New(),
			AgencyID:   s.agencies[2].ID,
			AgencyName: s.agencies[2].Name,
			FullName:   "Maya Chen",
			Email:      "maya@pacificaurban.com",
			Phone:      "+1-415-555-0901",
			Languages:  []string{"English", "Mandarin"},
			Region:     "Pacific Rim & Silicon Valley",
			PhotoURL:   "",
		},
		{
			ID:         uuid.New(),
			AgencyID:   s.agencies[3].ID,
			AgencyName: s.agencies[3].Name,
			FullName:   "Giulia Romano",
			Email:      "giulia@atlasheritage.it",
			Phone:      "+39-055-555-221",
			Languages:  []string{"Italian", "English", "French"},
			Region:     "Southern Europe & Mediterranean",
			PhotoURL:   "",
		},
		{
			ID:         uuid.New(),
			AgencyID:   s.agencies[2].ID,
			AgencyName: s.agencies[2].Name,
			FullName:   "Diego Alvarez",
			Email:      "diego@pacificaurban.com",
			Phone:      "+561-555-1758",
			Languages:  []string{"Spanish", "English"},
			Region:     "Latin America & US Sunbelt",
			PhotoURL:   "",
		},
	}

	for idx := range s.agencies {
		s.agencies[idx].Website = strings.TrimSpace(s.agencies[idx].Website)
	}
}

var _ Service = (*InMemoryService)(nil)
