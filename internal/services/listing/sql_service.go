package listing

import (
    "context"
    "database/sql"
    "encoding/json"
    "errors"
    "fmt"
    "strings"

    "github.com/google/uuid"
)

type sqlService struct {
	db *sql.DB
}

// NewSQLService builds a listing service backed by PostgreSQL.
func NewSQLService(db *sql.DB) (Service, error) {
	return &sqlService{db: db}, nil
}

func (s *sqlService) List(ctx context.Context) ([]Listing, error) {
	rows, err := s.db.QueryContext(ctx, `
        SELECT l.id, l.title, l.listing_type, l.country_code, l.city, l.region, l.neighborhood,
               l.summary, l.price, l.currency, l.bedrooms, l.bathrooms, l.area_sqm,
               l.hero_image_url, l.details_url, COALESCE(array_to_json(l.tags)::text, '[]'),
               l.agency_id, COALESCE(a.name, '')
        FROM property_listings l
        LEFT JOIN real_estate_agencies a ON a.id = l.agency_id
        ORDER BY l.created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	listings := make([]Listing, 0)
	for rows.Next() {
		var record Listing
		var agencyName string
		var price, bathrooms, area sql.NullFloat64
		var heroURL, detailsURL sql.NullString
		var tagsJSON string
		if err := rows.Scan(
			&record.ID,
			&record.Title,
			&record.Type,
			&record.Country,
			&record.City,
			&record.Region,
			&record.Neighborhood,
			&record.Summary,
			&price,
			&record.Currency,
			&record.Bedrooms,
			&bathrooms,
			&area,
			&heroURL,
			&detailsURL,
			&tagsJSON,
			&record.AgencyID,
			&agencyName,
		); err != nil {
			return nil, err
		}
		record.Price = price.Float64
		record.Bathrooms = bathrooms.Float64
		record.AreaSqM = area.Float64
		record.ImageURL = strings.TrimSpace(heroURL.String)
		record.DetailsURL = strings.TrimSpace(detailsURL.String)
		record.AgencyName = agencyName
		if err := json.Unmarshal([]byte(tagsJSON), &record.Tags); err != nil {
			record.Tags = nil
		}
		listings = append(listings, record)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return listings, nil
}

func (s *sqlService) Featured(ctx context.Context, limit int) ([]Listing, error) {
	if limit <= 0 {
		limit = 6
	}
	rows, err := s.db.QueryContext(ctx, `
        SELECT l.id, l.title, l.listing_type, l.country_code, l.city, l.region, l.neighborhood,
               l.summary, l.price, l.currency, l.bedrooms, l.bathrooms, l.area_sqm,
               l.hero_image_url, l.details_url, COALESCE(array_to_json(l.tags)::text, '[]'),
               l.agency_id, COALESCE(a.name, '')
        FROM property_listings l
        LEFT JOIN real_estate_agencies a ON a.id = l.agency_id
        ORDER BY l.created_at DESC
        LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	listings := make([]Listing, 0, limit)
	for rows.Next() {
		var record Listing
		var agencyName string
		var price, bathrooms, area sql.NullFloat64
		var heroURL, detailsURL sql.NullString
		var tagsJSON string
		if err := rows.Scan(
			&record.ID,
			&record.Title,
			&record.Type,
			&record.Country,
			&record.City,
			&record.Region,
			&record.Neighborhood,
			&record.Summary,
			&price,
			&record.Currency,
			&record.Bedrooms,
			&bathrooms,
			&area,
			&heroURL,
			&detailsURL,
			&tagsJSON,
			&record.AgencyID,
			&agencyName,
		); err != nil {
			return nil, err
		}
		record.Price = price.Float64
		record.Bathrooms = bathrooms.Float64
		record.AreaSqM = area.Float64
		record.ImageURL = strings.TrimSpace(heroURL.String)
		record.DetailsURL = strings.TrimSpace(detailsURL.String)
		record.AgencyName = agencyName
		if err := json.Unmarshal([]byte(tagsJSON), &record.Tags); err != nil {
			record.Tags = nil
		}
		listings = append(listings, record)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return listings, nil
}

func (s *sqlService) Get(ctx context.Context, id uuid.UUID) (Listing, error) {
	var record Listing
	var agencyName string
	var price, bathrooms, area sql.NullFloat64
	var heroURL, detailsURL sql.NullString
	var tagsJSON string
	err := s.db.QueryRowContext(ctx, `
        SELECT l.id, l.title, l.listing_type, l.country_code, l.city, l.region, l.neighborhood,
               l.summary, l.price, l.currency, l.bedrooms, l.bathrooms, l.area_sqm,
               l.hero_image_url, l.details_url, COALESCE(array_to_json(l.tags)::text, '[]'),
               l.agency_id, COALESCE(a.name, '')
        FROM property_listings l
        LEFT JOIN real_estate_agencies a ON a.id = l.agency_id
        WHERE l.id = $1`, id).
		Scan(
			&record.ID,
			&record.Title,
			&record.Type,
			&record.Country,
			&record.City,
			&record.Region,
			&record.Neighborhood,
			&record.Summary,
			&price,
			&record.Currency,
			&record.Bedrooms,
			&bathrooms,
			&area,
			&heroURL,
			&detailsURL,
			&tagsJSON,
			&record.AgencyID,
			&agencyName,
		)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return Listing{}, fmt.Errorf("listing %s not found", id)
        }
        return Listing{}, err
    }
	record.Price = price.Float64
	record.Bathrooms = bathrooms.Float64
	record.AreaSqM = area.Float64
	record.ImageURL = strings.TrimSpace(heroURL.String)
	record.DetailsURL = strings.TrimSpace(detailsURL.String)
	record.AgencyName = agencyName
	if err := json.Unmarshal([]byte(tagsJSON), &record.Tags); err != nil {
		record.Tags = nil
	}
	return record, nil
}

var _ Service = (*sqlService)(nil)
