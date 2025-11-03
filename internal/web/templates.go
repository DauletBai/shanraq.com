package web

import (
	"html/template"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"time"

	agencyservice "shanraq.com/internal/services/agency"
	listingservice "shanraq.com/internal/services/listing"
	transportservice "shanraq.com/internal/services/transport"
)

// Renderer wraps parsed HTML templates and exposes render helpers.
type Renderer struct {
	mu   sync.RWMutex
	base *template.Template
	fsys fs.FS
}

// BasePageData carries common metadata consumed by the shared layout.
type BasePageData struct {
	Theme       string
	PageTitle   string
	BrandName   string
	Description string
	PageID      string
	CurrentYear int
}

// HomePageData captures the dynamic properties injected into the landing page.
type HomePageData struct {
	BasePageData
	FeaturedListings  []ListingCard
	FeaturedAgencies  []AgencyCard
	FeaturedRealtors  []RealtorCard
	FeaturedTransport []TransportCard
}

// NewRenderer parses templates from the web directory.

func NewRenderer() (*Renderer, error) {
	webRoot := locateWebDir()
	fsys := os.DirFS(webRoot)

	tmpl := template.New("layout.html").Funcs(template.FuncMap{
		"statusColor": statusColor,
	})

	parsed, err := tmpl.ParseFS(fsys,
		"layout.html",
		"partials/*.html",
	)
	if err != nil {
		return nil, err
	}

	return &Renderer{
		base: parsed,
		fsys: fsys,
	}, nil
}

// RenderHome renders the primary landing page.
func (r *Renderer) RenderHome(w io.Writer, data *HomePageData) error {
	if data == nil {
		data = &HomePageData{}
	}

	if data.BrandName == "" {
		data.BrandName = "Shanraq"
	}
	if data.PageTitle == "" {
		data.PageTitle = "Discover Global Properties · "
	}
	if data.Description == "" {
		data.Description = "Search, compare, and manage international real estate listings from a single platform."
	}
	if data.PageID == "" {
		data.PageID = "home"
	}
	if data.CurrentYear == 0 {
		data.CurrentYear = time.Now().Year()
	}
	if data.Theme == "" {
		data.Theme = "auto"
	}

	r.mu.RLock()
	clone, err := r.base.Clone()
	r.mu.RUnlock()
	if err != nil {
		return err
	}

	if _, err := clone.ParseFS(r.fsys, "pages/home.html"); err != nil {
		return err
	}

	return clone.ExecuteTemplate(w, "layout.html", data)
}

// Unwrap exposes the underlying template for advanced use cases.
func (r *Renderer) Unwrap() *template.Template {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.base
}

func locateWebDir() string {
	wd, err := os.Getwd()
	if err != nil {
		return "web"
	}

	for i := 0; i < 10; i++ {
		candidate := filepath.Join(wd, "web")
		if _, statErr := os.Stat(filepath.Join(candidate, "layout.html")); statErr == nil {
			return candidate
		}
		parent := filepath.Dir(wd)
		if parent == wd {
			break
		}
		wd = parent
	}

	return "web"
}

func statusColor(status string) string {
	switch status {
	case "pending", "queued":
		return "warning"
	case "running", "in_progress":
		return "info"
	case "retry":
		return "secondary"
	case "failed", "error":
		return "danger"
	case "done", "completed", "success":
		return "success"
	default:
		return "secondary"
	}
}

// ListingCard represents a curated property listing for the landing page.
type ListingCard struct {
	ID          string
	Title       string
	Location    string
	Summary     string
	Price       string
	Thumbnail   string
	PropertyURL string
}

// AgencyCard represents an agency highlight.
type AgencyCard struct {
	ID      string
	Name    string
	Country string
	Website string
	Tagline string
	LogoURL string
}

// RealtorCard represents a featured realtor profile.
type RealtorCard struct {
	ID        string
	Name      string
	Agency    string
	Languages []string
	Region    string
	Email     string
}

// TransportCard represents a moving/logistics provider.
type TransportCard struct {
	ID          string
	Name        string
	CountryCode string
	Services    []string
	Coverage    []string
}

// MapListings converts listing service models into template cards.
func MapListings(listings []listingservice.Listing) []ListingCard {
	result := make([]ListingCard, 0, len(listings))
	for _, l := range listings {
		result = append(result, ListingCard{
			ID:          l.ID.String(),
			Title:       l.Title,
			Location:    l.LocationString(),
			Summary:     l.Summary,
			Price:       l.DisplayPrice(),
			Thumbnail:   l.ImageURL,
			PropertyURL: l.DetailsURL,
		})
	}
	return result
}

// MapAgencies converts agency service models into template cards.
func MapAgencies(agencies []agencyservice.Agency) []AgencyCard {
	result := make([]AgencyCard, 0, len(agencies))
	for _, a := range agencies {
		result = append(result, AgencyCard{
			ID:      a.ID.String(),
			Name:    a.Name,
			Country: a.Country,
			Website: a.Website,
			Tagline: a.Tagline,
			LogoURL: a.LogoURL,
		})
	}
	return result
}

// MapRealtors converts realtor models into template cards.
func MapRealtors(realtors []agencyservice.Realtor) []RealtorCard {
	result := make([]RealtorCard, 0, len(realtors))
	for _, r := range realtors {
		result = append(result, RealtorCard{
			ID:        r.ID.String(),
			Name:      r.FullName,
			Agency:    r.AgencyName,
			Languages: append([]string(nil), r.Languages...),
			Region:    r.Region,
			Email:     r.Email,
		})
	}
	return result
}

// MapTransportCompanies converts transport company models into template cards.
func MapTransportCompanies(companies []transportservice.Company) []TransportCard {
	result := make([]TransportCard, 0, len(companies))
	for _, c := range companies {
		result = append(result, TransportCard{
			ID:          c.ID.String(),
			Name:        c.Name,
			CountryCode: c.CountryCode,
			Services:    append([]string(nil), c.ServicesOffered...),
			Coverage:    append([]string(nil), c.CoverageRegions...),
		})
	}
	return result
}
