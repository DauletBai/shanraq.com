package web

import (
	"bytes"
	"strings"
	"testing"
)

func TestRenderHome(t *testing.T) {
	renderer, err := NewRenderer()
	if err != nil {
		t.Fatalf("NewRenderer() error = %v", err)
	}

	data := &HomePageData{}
	data.BrandName = "Shanraq"
	data.PageTitle = "Test Title · "
	data.Description = "Templating smoke test"
	data.PageID = "home"
	data.FeaturedListings = []ListingCard{
		{
			ID:        "listing-1",
			Title:     "Palm Jumeirah Sky Villa",
			Location:  "Dubai, UAE",
			Summary:   "Infinity pool with skyline views.",
			Price:     "USD 6800000",
			Thumbnail: "https://example.com/hero.jpg",
		},
	}
	data.FeaturedAgencies = []AgencyCard{{
		ID:      "agency-1",
		Name:    "Shanraq Global Realty",
		Country: "AE",
		Website: "https://shanraq.com/agency/global",
		Tagline: "Luxury Estates Across Continents",
	}}
	data.FeaturedRealtors = []RealtorCard{{
		ID:        "realtor-1",
		Name:      "Layla Al-Mansouri",
		Agency:    "Shanraq Global Realty",
		Languages: []string{"Arabic", "English"},
		Region:    "MENA",
		Email:     "layla@example.com",
	}}
	data.FeaturedTransport = []TransportCard{{
		ID:          "transport-1",
		Name:        "Atlas Relocation Partners",
		CountryCode: "AE",
		Services:    []string{"packing", "customs"},
		Coverage:    []string{"Middle East"},
	}}

	var buf bytes.Buffer
	if err := renderer.RenderHome(&buf, data); err != nil {
		t.Fatalf("RenderHome() error = %v", err)
	}

	html := buf.String()
	mustContain := []string{
		"Test Title",
		"Discover global real estate",
		"Palm Jumeirah Sky Villa",
		"Shanraq Global Realty",
		"Atlas Relocation Partners",
	}

	for _, token := range mustContain {
		if !strings.Contains(html, token) {
			t.Fatalf("rendered home page missing %q", token)
		}
	}

	// ensure the content block is wrapped by the layout's main container
	if !strings.Contains(html, "<main class=\"container py-5\">") {
		t.Fatalf("layout main container not found in rendered output")
	}
}
