package agency

import (
	"context"
	"database/sql"
	"encoding/json"
	"sort"
	"strings"
)

type sqlRepository struct {
	db *sql.DB
}

type sqlService struct {
	repo *sqlRepository
}

// NewSQLService returns a Service backed by PostgreSQL data.
func NewSQLService(db *sql.DB) (Service, error) {
	return &sqlService{repo: &sqlRepository{db: db}}, nil
}

func (s *sqlService) ListAgencies(ctx context.Context) ([]Agency, error) {
	return s.repo.listAgencies(ctx)
}

func (s *sqlService) Featured(ctx context.Context, limit int) ([]Agency, error) {
	agencies, err := s.repo.featuredAgencies(ctx, limit)
	if err != nil {
		return nil, err
	}
	return agencies, nil
}

func (s *sqlService) ListRealtors(ctx context.Context) ([]Realtor, error) {
	return s.repo.listRealtors(ctx)
}

func (s *sqlService) FeaturedRealtors(ctx context.Context, limit int) ([]Realtor, error) {
	return s.repo.featuredRealtors(ctx, limit)
}

func (r *sqlRepository) listAgencies(ctx context.Context) ([]Agency, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, tagline, country_code, website, logo_url, head_office FROM real_estate_agencies ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	agencies := make([]Agency, 0)
	for rows.Next() {
		var agency Agency
		var website, logo, headOffice sql.NullString
		if err := rows.Scan(&agency.ID, &agency.Name, &agency.Tagline, &agency.Country, &website, &logo, &headOffice); err != nil {
			return nil, err
		}
		agency.Website = strings.TrimSpace(website.String)
		agency.LogoURL = strings.TrimSpace(logo.String)
		agency.HeadOffice = headOffice.String
		agencies = append(agencies, agency)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return agencies, nil
}

func (r *sqlRepository) featuredAgencies(ctx context.Context, limit int) ([]Agency, error) {
	if limit <= 0 {
		limit = 3
	}
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, tagline, country_code, website, logo_url, head_office FROM real_estate_agencies ORDER BY name LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	agencies := make([]Agency, 0)
	for rows.Next() {
		var agency Agency
		var website, logo, headOffice sql.NullString
		if err := rows.Scan(&agency.ID, &agency.Name, &agency.Tagline, &agency.Country, &website, &logo, &headOffice); err != nil {
			return nil, err
		}
		agency.Website = strings.TrimSpace(website.String)
		agency.LogoURL = strings.TrimSpace(logo.String)
		agency.HeadOffice = headOffice.String
		agencies = append(agencies, agency)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return agencies, nil
}

func (r *sqlRepository) listRealtors(ctx context.Context) ([]Realtor, error) {
	rows, err := r.db.QueryContext(ctx, `
        SELECT r.id, r.agency_id, COALESCE(a.name, ''), r.full_name, r.email, r.phone,
               COALESCE(array_to_json(r.languages)::text, '[]'), r.region, r.photo_url
        FROM realtors r
        LEFT JOIN real_estate_agencies a ON a.id = r.agency_id
        ORDER BY r.full_name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	realtors := make([]Realtor, 0)
	for rows.Next() {
		var realtor Realtor
		var agencyName string
		var email, phone, region, photo sql.NullString
		var langsJSON string
		if err := rows.Scan(&realtor.ID, &realtor.AgencyID, &agencyName, &realtor.FullName, &email, &phone, &langsJSON, &region, &photo); err != nil {
			return nil, err
		}
		realtor.AgencyName = agencyName
		realtor.Email = strings.TrimSpace(email.String)
		realtor.Phone = strings.TrimSpace(phone.String)
		realtor.Region = region.String
		realtor.PhotoURL = photo.String
		if err := json.Unmarshal([]byte(langsJSON), &realtor.Languages); err != nil {
			realtor.Languages = nil
		}
		realtors = append(realtors, realtor)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	sort.Slice(realtors, func(i, j int) bool {
		return realtors[i].FullName < realtors[j].FullName
	})
	return realtors, nil
}

func (r *sqlRepository) featuredRealtors(ctx context.Context, limit int) ([]Realtor, error) {
	if limit <= 0 {
		limit = 4
	}
	rows, err := r.db.QueryContext(ctx, `
        SELECT r.id, r.agency_id, COALESCE(a.name, ''), r.full_name, r.email, r.phone,
               COALESCE(array_to_json(r.languages)::text, '[]'), r.region, r.photo_url
        FROM realtors r
        LEFT JOIN real_estate_agencies a ON a.id = r.agency_id
        ORDER BY r.full_name
        LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]Realtor, 0)
	for rows.Next() {
		var realtor Realtor
		var agencyName string
		var email, phone, region, photo sql.NullString
		var langsJSON string
		if err := rows.Scan(&realtor.ID, &realtor.AgencyID, &agencyName, &realtor.FullName, &email, &phone, &langsJSON, &region, &photo); err != nil {
			return nil, err
		}
		realtor.AgencyName = agencyName
		realtor.Email = strings.TrimSpace(email.String)
		realtor.Phone = strings.TrimSpace(phone.String)
		realtor.Region = region.String
		realtor.PhotoURL = photo.String
		if err := json.Unmarshal([]byte(langsJSON), &realtor.Languages); err != nil {
			realtor.Languages = nil
		}
		list = append(list, realtor)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

var _ Service = (*sqlService)(nil)
