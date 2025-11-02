package listing

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/google/uuid"
)

// ListingType enumerates property categories.
type ListingType string

const (
	ListingTypeResidential ListingType = "residential"
	ListingTypeCommercial  ListingType = "commercial"
	ListingTypeLand        ListingType = "land"
)

// Listing represents an individual property.
type Listing struct {
	ID           uuid.UUID   `json:"id"`
	Title        string      `json:"title"`
	Type         ListingType `json:"type"`
	Country      string      `json:"country"`
	City         string      `json:"city"`
	Region       string      `json:"region"`
	Neighborhood string      `json:"neighborhood"`
	Summary      string      `json:"summary"`
	Price        float64     `json:"price"`
	Currency     string      `json:"currency"`
	Bedrooms     int         `json:"bedrooms"`
	Bathrooms    float64     `json:"bathrooms"`
	AreaSqM      float64     `json:"area_sqm"`
	ImageURL     string      `json:"image_url"`
	DetailsURL   string      `json:"details_url"`
	AgencyID     uuid.UUID   `json:"agency_id"`
	AgencyName   string      `json:"agency_name"`
	Tags         []string    `json:"tags"`
}

// Service exposes read access to listings.
type Service interface {
	List(ctx context.Context) ([]Listing, error)
	Featured(ctx context.Context, limit int) ([]Listing, error)
	Get(ctx context.Context, id uuid.UUID) (Listing, error)
}

// InMemoryService is a seeded implementation.
type InMemoryService struct {
	mu       sync.RWMutex
	listings []Listing
}

// NewInMemoryService seeds demo listings.
func NewInMemoryService() *InMemoryService {
	svc := &InMemoryService{}
	svc.seed()
	return svc
}

func (s *InMemoryService) List(_ context.Context) ([]Listing, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]Listing, len(s.listings))
	copy(out, s.listings)
	sort.Slice(out, func(i, j int) bool {
		return out[i].Title < out[j].Title
	})
	return out, nil
}

func (s *InMemoryService) Featured(ctx context.Context, limit int) ([]Listing, error) {
	listings, err := s.List(ctx)
	if err != nil {
		return nil, err
	}
	if limit <= 0 || limit > len(listings) {
		limit = len(listings)
	}
	return listings[:limit], nil
}

func (s *InMemoryService) Get(_ context.Context, id uuid.UUID) (Listing, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, l := range s.listings {
		if l.ID == id {
			return l, nil
		}
	}
	return Listing{}, fmt.Errorf("listing %s not found", id)
}

// LocationString returns a formatted location string.
func (l Listing) LocationString() string {
	parts := []string{}
	if l.Neighborhood != "" {
		parts = append(parts, l.Neighborhood)
	}
	if l.City != "" {
		parts = append(parts, l.City)
	}
	if l.Region != "" {
		parts = append(parts, l.Region)
	}
	if l.Country != "" {
		parts = append(parts, l.Country)
	}
	return strings.Join(parts, ", ")
}

// DisplayPrice returns a formatted price string.
func (l Listing) DisplayPrice() string {
	if l.Price == 0 {
		return "Price on request"
	}
	return fmt.Sprintf("%s %.0f", strings.ToUpper(l.Currency), l.Price)
}

func (s *InMemoryService) seed() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.listings) > 0 {
		return
	}

	s.listings = []Listing{
		{
			ID:           uuid.New(),
			Title:        "Palm Jumeirah Sky Villa",
			Type:         ListingTypeResidential,
			Country:      "AE",
			City:         "Dubai",
			Region:       "Dubai",
			Neighborhood: "Palm Jumeirah",
			Summary:      "Four-bedroom duplex with private infinity pool and uninterrupted Gulf views.",
			Price:        6800000,
			Currency:     "USD",
			Bedrooms:     4,
			Bathrooms:    4.5,
			AreaSqM:      620,
			ImageURL:     "https://images.shanraq.com/demo/palm-jumeirah.jpg",
			DetailsURL:   "/listings/palm-jumeirah-sky-villa",
			AgencyName:   "Shanraq Global Realty",
			Tags:         []string{"waterfront", "smart-home", "penthouse"},
		},
		{
			ID:           uuid.New(),
			Title:        "Östermalm Art Nouveau Residence",
			Type:         ListingTypeResidential,
			Country:      "SE",
			City:         "Stockholm",
			Region:       "Stockholm County",
			Neighborhood: "Östermalm",
			Summary:      "Restored 1903 apartment with modern energy systems and winter garden.",
			Price:        14500000,
			Currency:     "SEK",
			Bedrooms:     3,
			Bathrooms:    2,
			AreaSqM:      210,
			ImageURL:     "https://images.shanraq.com/demo/ostermalm.jpg",
			DetailsURL:   "/listings/ostermalm-art-nouveau",
			AgencyName:   "Nordic Skyline Partners",
			Tags:         []string{"heritage", "city-center"},
		},
		{
			ID:           uuid.New(),
			Title:        "Kyoto Machiya Boutique Hotel",
			Type:         ListingTypeCommercial,
			Country:      "JP",
			City:         "Kyoto",
			Region:       "Kansai",
			Neighborhood: "Gion",
			Summary:      "Six-key licensed machiya hotel blending traditional architecture with modern amenities.",
			Price:        215000000,
			Currency:     "JPY",
			Bedrooms:     6,
			Bathrooms:    6.5,
			AreaSqM:      380,
			ImageURL:     "https://images.shanraq.com/demo/kyoto-machiya.jpg",
			DetailsURL:   "/listings/kyoto-machiya-hotel",
			AgencyName:   "Pacifica Urban Advisors",
			Tags:         []string{"hospitality", "licensed", "turnkey"},
		},
		{
			ID:           uuid.New(),
			Title:        "Lisbon Digital District Loft",
			Type:         ListingTypeResidential,
			Country:      "PT",
			City:         "Lisbon",
			Region:       "Lisbon",
			Neighborhood: "Parque das Nações",
			Summary:      "Smart loft with Tagus river views, co-working mezzanine, and EV-ready parking.",
			Price:        890000,
			Currency:     "EUR",
			Bedrooms:     2,
			Bathrooms:    2,
			AreaSqM:      165,
			ImageURL:     "https://images.shanraq.com/demo/lisbon-loft.jpg",
			DetailsURL:   "/listings/lisbon-digital-loft",
			AgencyName:   "Pacifica Urban Advisors",
			Tags:         []string{"smart-home", "waterfront", "digital-nomad"},
		},
		{
			ID:           uuid.New(),
			Title:        "Tuscany Heritage Vineyard Estate",
			Type:         ListingTypeCommercial,
			Country:      "IT",
			City:         "Siena",
			Region:       "Tuscany",
			Neighborhood: "Chianti Classico",
			Summary:      "Organic DOCG vineyard with 18th-century villa, guest suites, and agritourism permit.",
			Price:        6200000,
			Currency:     "EUR",
			Bedrooms:     8,
			Bathrooms:    7,
			AreaSqM:      950,
			ImageURL:     "https://images.shanraq.com/demo/tuscany-vineyard.jpg",
			DetailsURL:   "/listings/tuscany-heritage-vineyard",
			AgencyName:   "Atlas Heritage Homes",
			Tags:         []string{"vineyard", "heritage", "agritourism"},
		},
		{
			ID:           uuid.New(),
			Title:        "Singapore Sky Garden Duplex",
			Type:         ListingTypeResidential,
			Country:      "SG",
			City:         "Singapore",
			Region:       "Central Region",
			Neighborhood: "Marina Bay",
			Summary:      "Biophilic duplex with hydroponic greenhouse, concierge robotics, and Marina skyline views.",
			Price:        12800000,
			Currency:     "SGD",
			Bedrooms:     4,
			Bathrooms:    4,
			AreaSqM:      420,
			ImageURL:     "https://images.shanraq.com/demo/singapore-sky-garden.jpg",
			DetailsURL:   "/listings/singapore-sky-garden",
			AgencyName:   "Shanraq Global Realty",
			Tags:         []string{"biophilic", "city-center"},
		},
		{
			ID:           uuid.New(),
			Title:        "Reykjavík Geothermal Retreat",
			Type:         ListingTypeResidential,
			Country:      "IS",
			City:         "Reykjavík",
			Region:       "Capital Region",
			Neighborhood: "Mosfellsbær",
			Summary:      "Net-zero villa with geothermal spa wing, aurora lounge, and drone landing pad.",
			Price:        325000000,
			Currency:     "ISK",
			Bedrooms:     5,
			Bathrooms:    4,
			AreaSqM:      480,
			ImageURL:     "https://images.shanraq.com/demo/reykjavik-retreat.jpg",
			DetailsURL:   "/listings/reykjavik-geothermal-retreat",
			AgencyName:   "Nordic Skyline Partners",
			Tags:         []string{"net-zero", "luxury", "spa"},
		},
		{
			ID:           uuid.New(),
			Title:        "Cape Town Atlantic Seaboard Villa",
			Type:         ListingTypeResidential,
			Country:      "ZA",
			City:         "Cape Town",
			Region:       "Western Cape",
			Neighborhood: "Bantry Bay",
			Summary:      "Secure cliffside villa with desalination system, solar microgrid, and cinematic pavilion.",
			Price:        39500000,
			Currency:     "ZAR",
			Bedrooms:     6,
			Bathrooms:    6.5,
			AreaSqM:      720,
			ImageURL:     "https://images.shanraq.com/demo/capetown-villa.jpg",
			DetailsURL:   "/listings/cape-town-atlantic-villa",
			AgencyName:   "Shanraq Global Realty",
			Tags:         []string{"coastal", "security", "solar"},
		},
		{
			ID:           uuid.New(),
			Title:        "São Paulo Innovation Hub Loft",
			Type:         ListingTypeCommercial,
			Country:      "BR",
			City:         "São Paulo",
			Region:       "São Paulo",
			Neighborhood: "Vila Olímpia",
			Summary:      "Adaptive reuse warehouse with 5G infrastructure, studios, and data lounge.",
			Price:        11800000,
			Currency:     "BRL",
			Bedrooms:     0,
			Bathrooms:    4,
			AreaSqM:      980,
			ImageURL:     "https://images.shanraq.com/demo/sao-paulo-hub.jpg",
			DetailsURL:   "/listings/sao-paulo-innovation-hub",
			AgencyName:   "Pacifica Urban Advisors",
			Tags:         []string{"innovation", "mixed-use"},
		},
		{
			ID:           uuid.New(),
			Title:        "British Columbia Wilderness Lodge",
			Type:         ListingTypeCommercial,
			Country:      "CA",
			City:         "Whistler",
			Region:       "British Columbia",
			Neighborhood: "Callaghan Valley",
			Summary:      "Heli-access eco lodge with 12 guest suites, carbon-negative design, and heli pads.",
			Price:        8600000,
			Currency:     "CAD",
			Bedrooms:     12,
			Bathrooms:    12,
			AreaSqM:      1250,
			ImageURL:     "https://images.shanraq.com/demo/bc-wilderness.jpg",
			DetailsURL:   "/listings/bc-wilderness-lodge",
			AgencyName:   "Nordic Skyline Partners",
			Tags:         []string{"eco", "adventure", "hospitality"},
		},
	}
}

var _ Service = (*InMemoryService)(nil)
