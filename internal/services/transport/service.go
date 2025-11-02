package transport

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Company models a logistics partner supporting property relocations.
type Company struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	Slug            string    `json:"slug"`
	CountryCode     string    `json:"country_code"`
	CoverageRegions []string  `json:"coverage_regions"`
	ServicesOffered []string  `json:"services_offered"`
	ContactEmail    string    `json:"contact_email"`
	ContactPhone    string    `json:"contact_phone"`
	Website         string    `json:"website"`
	Description     string    `json:"description"`
	Active          bool      `json:"active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// ListFilter captures query parameters for listing companies.
type ListFilter struct {
	CountryCode string
	Limit       int
	Offset      int
	ActiveOnly  bool
}

// CreateInput defines attributes required to create a transport company.
type CreateInput struct {
	Name            string
	CountryCode     string
	CoverageRegions []string
	ServicesOffered []string
	ContactEmail    string
	ContactPhone    string
	Website         string
	Description     string
	Active          bool
}

// UpdateInput defines mutable fields for a transport company.
type UpdateInput struct {
	Name            *string
	CountryCode     *string
	CoverageRegions *[]string
	ServicesOffered *[]string
	ContactEmail    *string
	ContactPhone    *string
	Website         *string
	Description     *string
	Active          *bool
}

// Service exposes the CRUD capabilities for transport companies.
type Service interface {
	List(ctx context.Context, filter ListFilter) ([]Company, int, error)
	Create(ctx context.Context, input CreateInput) (Company, error)
	Get(ctx context.Context, id uuid.UUID) (Company, error)
	Update(ctx context.Context, id uuid.UUID, input UpdateInput) (Company, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// ErrNotFound is returned when a company cannot be located.
var ErrNotFound = errors.New("transport company not found")

// InMemoryService is a lightweight implementation used until a database-backed repo is plugged in.
type InMemoryService struct {
	mu        sync.RWMutex
	companies map[uuid.UUID]Company
}

// NewInMemoryService provides a concurrency-safe in-memory repository.
func NewInMemoryService() *InMemoryService {
	svc := &InMemoryService{
		companies: make(map[uuid.UUID]Company),
	}
	svc.seed()
	return svc
}

// List returns transport companies applying filter criteria.
func (s *InMemoryService) List(_ context.Context, filter ListFilter) ([]Company, int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var results []Company
	for _, c := range s.companies {
		if filter.CountryCode != "" && !strings.EqualFold(filter.CountryCode, c.CountryCode) {
			continue
		}
		if filter.ActiveOnly && !c.Active {
			continue
		}
		results = append(results, c)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Name < results[j].Name
	})

	start := filter.Offset
	if start > len(results) {
		return []Company{}, len(results), nil
	}

	end := len(results)
	if filter.Limit > 0 && start+filter.Limit < end {
		end = start + filter.Limit
	}

	return results[start:end], len(results), nil
}

// Create registers a new transport company and generates a unique slug.
func (s *InMemoryService) Create(_ context.Context, input CreateInput) (Company, error) {
	if strings.TrimSpace(input.Name) == "" {
		return Company{}, errors.New("name is required")
	}
	if len(input.CountryCode) != 2 {
		return Company{}, errors.New("country code must be ISO 3166-1 alpha-2")
	}

	now := time.Now().UTC()
	if !strings.Contains(input.ContactEmail, "@") && input.ContactEmail != "" {
		return Company{}, errors.New("invalid contact email")
	}

	company := Company{
		ID:              uuid.New(),
		Name:            strings.TrimSpace(input.Name),
		Slug:            "",
		CountryCode:     strings.ToUpper(input.CountryCode),
		CoverageRegions: dedupeStrings(input.CoverageRegions),
		ServicesOffered: dedupeStrings(input.ServicesOffered),
		ContactEmail:    strings.TrimSpace(input.ContactEmail),
		ContactPhone:    strings.TrimSpace(input.ContactPhone),
		Website:         strings.TrimSpace(input.Website),
		Description:     strings.TrimSpace(input.Description),
		Active:          input.Active,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if company.Active == false {
		company.Active = true
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	company.Slug = s.generateUniqueSlug(company.Name)
	s.companies[company.ID] = company
	return company, nil
}

// Get retrieves a transport company by ID.
func (s *InMemoryService) Get(_ context.Context, id uuid.UUID) (Company, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	company, ok := s.companies[id]
	if !ok {
		return Company{}, ErrNotFound
	}
	return company, nil
}

// Update mutates fields for a transport company.
func (s *InMemoryService) Update(_ context.Context, id uuid.UUID, input UpdateInput) (Company, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	company, ok := s.companies[id]
	if !ok {
		return Company{}, ErrNotFound
	}

	if input.Name != nil {
		if strings.TrimSpace(*input.Name) == "" {
			return Company{}, errors.New("name cannot be empty")
		}
		company.Name = strings.TrimSpace(*input.Name)
		company.Slug = s.generateUniqueSlug(company.Name)
	}
	if input.CountryCode != nil {
		if len(*input.CountryCode) != 2 {
			return Company{}, errors.New("country code must be ISO 3166-1 alpha-2")
		}
		company.CountryCode = strings.ToUpper(strings.TrimSpace(*input.CountryCode))
	}
	if input.CoverageRegions != nil {
		company.CoverageRegions = dedupeStrings(*input.CoverageRegions)
	}
	if input.ServicesOffered != nil {
		company.ServicesOffered = dedupeStrings(*input.ServicesOffered)
	}
	if input.ContactEmail != nil {
		if email := strings.TrimSpace(*input.ContactEmail); email != "" && !strings.Contains(email, "@") {
			return Company{}, errors.New("invalid contact email")
		} else {
			company.ContactEmail = email
		}
	}
	if input.ContactPhone != nil {
		company.ContactPhone = strings.TrimSpace(*input.ContactPhone)
	}
	if input.Website != nil {
		company.Website = strings.TrimSpace(*input.Website)
	}
	if input.Description != nil {
		company.Description = strings.TrimSpace(*input.Description)
	}
	if input.Active != nil {
		company.Active = *input.Active
	}

	company.UpdatedAt = time.Now().UTC()
	s.companies[id] = company
	return company, nil
}

// Delete removes a transport company permanently.
func (s *InMemoryService) Delete(_ context.Context, id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.companies[id]; !ok {
		return ErrNotFound
	}
	delete(s.companies, id)
	return nil
}

func (s *InMemoryService) generateUniqueSlug(name string) string {
	base := slugify(name)
	if base == "" {
		base = "transport-company"
	}
	slug := base
	counter := 1
	for s.slugExists(slug) {
		counter++
		slug = fmt.Sprintf("%s-%d", base, counter)
	}
	return slug
}

func (s *InMemoryService) slugExists(slug string) bool {
	for _, c := range s.companies {
		if c.Slug == slug {
			return true
		}
	}
	return false
}

func slugify(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	value = strings.ReplaceAll(value, "&", "and")
	value = strings.ReplaceAll(value, ".", " ")
	value = strings.ReplaceAll(value, "_", " ")
	value = strings.ReplaceAll(value, "/", " ")
	value = strings.ReplaceAll(value, "'", "")
	fields := strings.Fields(value)
	return strings.Join(fields, "-")
}

func dedupeStrings(values []string) []string {
	set := make(map[string]struct{})
	var result []string
	for _, v := range values {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		if _, exists := set[v]; exists {
			continue
		}
		set[v] = struct{}{}
		result = append(result, v)
	}
	sort.Strings(result)
	return result
}

func (s *InMemoryService) seed() {
	s.mu.Lock()
	if len(s.companies) > 0 {
		s.mu.Unlock()
		return
	}
	s.mu.Unlock()

	seedData := []CreateInput{
		{
			Name:            "Atlas Relocation Partners",
			CountryCode:     "AE",
			CoverageRegions: []string{"Middle East", "North Africa", "South Asia"},
			ServicesOffered: []string{"packing", "customs-clearance", "storage", "visa-support"},
			ContactEmail:    "hello@atlasrelocation.global",
			ContactPhone:    "+971-4-555-1200",
			Website:         "https://atlasrelocation.global",
			Description:     "Premium relocations for luxury and corporate clients with AI-driven logistics.",
			Active:          true,
		},
		{
			Name:            "Nordic Move & Care",
			CountryCode:     "SE",
			CoverageRegions: []string{"Scandinavia", "Baltics", "Northern Europe"},
			ServicesOffered: []string{"packing", "temperature-controlled", "pet-relocation"},
			ContactEmail:    "support@nordicmove.se",
			ContactPhone:    "+46-8-777-9911",
			Website:         "https://nordicmove.se",
			Description:     "Sustainable moving services with carbon accounting and smart storage.",
			Active:          true,
		},
		{
			Name:            "Pacifica Trans-Pacific Logistics",
			CountryCode:     "US",
			CoverageRegions: []string{"North America", "East Asia", "Oceania"},
			ServicesOffered: []string{"packing", "freight-forwarding", "corporate-relocation"},
			ContactEmail:    "contact@pacifica-transpac.com",
			ContactPhone:    "+1-415-555-4550",
			Website:         "https://pacifica-transpac.com",
			Description:     "AI-optimized freight routes and compliance for cross-border moves.",
			Active:          true,
		},
		{
			Name:            "Heritage Art Movers",
			CountryCode:     "IT",
			CoverageRegions: []string{"Southern Europe", "Middle East"},
			ServicesOffered: []string{"art-handling", "climate-storage", "white-glove"},
			ContactEmail:    "info@heritageartmovers.it",
			ContactPhone:    "+39-055-555-889",
			Website:         "https://heritageartmovers.it",
			Description:     "Boutique relocations for art collections and heritage properties.",
			Active:          true,
		},
		{
			Name:            "Southern Cross Mobility",
			CountryCode:     "AU",
			CoverageRegions: []string{"Australia", "New Zealand", "Pacific Islands"},
			ServicesOffered: []string{"packing", "vehicle-shipping", "remote-installations"},
			ContactEmail:    "team@southerncrossmobility.au",
			ContactPhone:    "+61-2-5555-4422",
			Website:         "https://southerncrossmobility.au",
			Description:     "Remote area specialists with autonomous convoy support.",
			Active:          true,
		},
	}

	for _, input := range seedData {
		_, _ = s.Create(context.Background(), input)
	}
}

// Ensure InMemoryService fulfils Service interface.
var _ Service = (*InMemoryService)(nil)
