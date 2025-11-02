package transport

import (
	"time"

	"github.com/google/uuid"

	transportservice "shanraq.com/internal/services/transport"
)

type companyResponse struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	Slug            string    `json:"slug"`
	CountryCode     string    `json:"country_code"`
	CoverageRegions []string  `json:"coverage_regions"`
	ServicesOffered []string  `json:"services_offered"`
	ContactEmail    string    `json:"contact_email,omitempty"`
	ContactPhone    string    `json:"contact_phone,omitempty"`
	Website         string    `json:"website,omitempty"`
	Description     string    `json:"description,omitempty"`
	Active          bool      `json:"active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type listResponse struct {
	Data []companyResponse `json:"data"`
	Meta listMeta          `json:"meta"`
}

type listMeta struct {
	Total   int    `json:"total"`
	Limit   int    `json:"limit"`
	Offset  int    `json:"offset"`
	Country string `json:"country,omitempty"`
}

type createRequest struct {
	Name            string   `json:"name"`
	CountryCode     string   `json:"country_code"`
	CoverageRegions []string `json:"coverage_regions"`
	ServicesOffered []string `json:"services_offered"`
	ContactEmail    string   `json:"contact_email"`
	ContactPhone    string   `json:"contact_phone"`
	Website         string   `json:"website"`
	Description     string   `json:"description"`
	Active          *bool    `json:"active"`
}

type updateRequest struct {
	Name            *string  `json:"name"`
	CountryCode     *string  `json:"country_code"`
	CoverageRegions []string `json:"coverage_regions"`
	ServicesOffered []string `json:"services_offered"`
	ContactEmail    *string  `json:"contact_email"`
	ContactPhone    *string  `json:"contact_phone"`
	Website         *string  `json:"website"`
	Description     *string  `json:"description"`
	Active          *bool    `json:"active"`
}

func mapToResponse(company transportservice.Company) companyResponse {
	return companyResponse{
		ID:              company.ID,
		Name:            company.Name,
		Slug:            company.Slug,
		CountryCode:     company.CountryCode,
		CoverageRegions: company.CoverageRegions,
		ServicesOffered: company.ServicesOffered,
		ContactEmail:    company.ContactEmail,
		ContactPhone:    company.ContactPhone,
		Website:         company.Website,
		Description:     company.Description,
		Active:          company.Active,
		CreatedAt:       company.CreatedAt,
		UpdatedAt:       company.UpdatedAt,
	}
}
