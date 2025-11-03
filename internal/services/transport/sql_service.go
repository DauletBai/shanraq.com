package transport

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type sqlService struct {
	db *sql.DB
}

// NewSQLService builds a transport company service backed by PostgreSQL.
func NewSQLService(db *sql.DB) (Service, error) {
	return &sqlService{db: db}, nil
}

func (s *sqlService) List(ctx context.Context, filter ListFilter) ([]Company, int, error) {
	clauses := make([]string, 0)
	args := make([]interface{}, 0)

	if filter.CountryCode != "" {
		args = append(args, strings.ToUpper(filter.CountryCode))
		clauses = append(clauses, fmt.Sprintf("country_code = $%d", len(args)))
	}
	if filter.ActiveOnly {
		args = append(args, true)
		clauses = append(clauses, fmt.Sprintf("active = $%d", len(args)))
	}

	where := ""
	if len(clauses) > 0 {
		where = "WHERE " + strings.Join(clauses, " AND ")
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM transport_companies %s", where)
	var total int
	if err := s.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	queryArgs := append([]interface{}{}, args...)
	limit := filter.Limit
	offset := filter.Offset
	if limit <= 0 {
		limit = 50
	}
	queryArgs = append(queryArgs, limit, offset)

	listQuery := fmt.Sprintf(`
        SELECT id, name, slug, country_code,
               COALESCE(array_to_json(coverage_regions)::text, '[]') AS coverage,
               COALESCE(array_to_json(services_offered)::text, '[]') AS services,
               contact_email, contact_phone, website, description, active,
               created_at, updated_at
        FROM transport_companies
        %s
        ORDER BY name
        LIMIT $%d OFFSET $%d`, where, len(queryArgs)-1, len(queryArgs))

	rows, err := s.db.QueryContext(ctx, listQuery, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	companies := make([]Company, 0)
	for rows.Next() {
		company, err := scanCompany(rows)
		if err != nil {
			return nil, 0, err
		}
		companies = append(companies, company)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return companies, total, nil
}

func (s *sqlService) Create(ctx context.Context, input CreateInput) (Company, error) {
	if strings.TrimSpace(input.Name) == "" {
		return Company{}, errors.New("name is required")
	}
	if len(input.CountryCode) != 2 {
		return Company{}, errors.New("country code must be ISO 3166-1 alpha-2")
	}

	slugBase := slugify(input.Name)
	slug, err := s.generateUniqueSlug(ctx, slugBase)
	if err != nil {
		return Company{}, err
	}

	coverage := dedupeStrings(input.CoverageRegions)
	services := dedupeStrings(input.ServicesOffered)

	now := time.Now().UTC()

	var company Company
	err = s.db.QueryRowContext(ctx, `
        INSERT INTO transport_companies
            (name, slug, country_code, coverage_regions, services_offered,
             contact_email, contact_phone, website, description, active, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $11)
        RETURNING id, created_at, updated_at`,
		strings.TrimSpace(input.Name),
		slug,
		strings.ToUpper(input.CountryCode),
		coverage,
		services,
		strings.TrimSpace(input.ContactEmail),
		strings.TrimSpace(input.ContactPhone),
		strings.TrimSpace(input.Website),
		strings.TrimSpace(input.Description),
		input.Active,
		now,
	).Scan(&company.ID, &company.CreatedAt, &company.UpdatedAt)
	if err != nil {
		return Company{}, err
	}

	company.Name = strings.TrimSpace(input.Name)
	company.Slug = slug
	company.CountryCode = strings.ToUpper(input.CountryCode)
	company.CoverageRegions = dedupeStrings(input.CoverageRegions)
	company.ServicesOffered = dedupeStrings(input.ServicesOffered)
	company.ContactEmail = strings.TrimSpace(input.ContactEmail)
	company.ContactPhone = strings.TrimSpace(input.ContactPhone)
	company.Website = strings.TrimSpace(input.Website)
	company.Description = strings.TrimSpace(input.Description)
	company.Active = input.Active
	return company, nil
}

func (s *sqlService) Get(ctx context.Context, id uuid.UUID) (Company, error) {
    row := s.db.QueryRowContext(ctx, `
        SELECT id, name, slug, country_code,
               COALESCE(array_to_json(coverage_regions)::text, '[]'),
               COALESCE(array_to_json(services_offered)::text, '[]'),
               contact_email, contact_phone, website, description, active,
               created_at, updated_at
        FROM transport_companies
        WHERE id = $1`, id)
    company, err := scanCompany(row)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return Company{}, ErrNotFound
        }
        return Company{}, err
    }
    return company, nil
}

func (s *sqlService) Update(ctx context.Context, id uuid.UUID, input UpdateInput) (Company, error) {
	existing, err := s.Get(ctx, id)
	if err != nil {
		return Company{}, err
	}

	if input.Name != nil {
		trimmed := strings.TrimSpace(*input.Name)
		if trimmed == "" {
			return Company{}, errors.New("name cannot be empty")
		}
		if !strings.EqualFold(trimmed, existing.Name) {
			slug, err := s.generateUniqueSlug(ctx, slugify(trimmed))
			if err != nil {
				return Company{}, err
			}
			existing.Slug = slug
		}
		existing.Name = trimmed
	}
	if input.CountryCode != nil {
		if len(*input.CountryCode) != 2 {
			return Company{}, errors.New("country code must be ISO 3166-1 alpha-2")
		}
		existing.CountryCode = strings.ToUpper(strings.TrimSpace(*input.CountryCode))
	}
	if input.CoverageRegions != nil {
		existing.CoverageRegions = dedupeStrings(*input.CoverageRegions)
	}
	if input.ServicesOffered != nil {
		existing.ServicesOffered = dedupeStrings(*input.ServicesOffered)
	}
	if input.ContactEmail != nil {
		existing.ContactEmail = strings.TrimSpace(*input.ContactEmail)
	}
	if input.ContactPhone != nil {
		existing.ContactPhone = strings.TrimSpace(*input.ContactPhone)
	}
	if input.Website != nil {
		existing.Website = strings.TrimSpace(*input.Website)
	}
	if input.Description != nil {
		existing.Description = strings.TrimSpace(*input.Description)
	}
	if input.Active != nil {
		existing.Active = *input.Active
	}

	coverage := dedupeStrings(existing.CoverageRegions)
	services := dedupeStrings(existing.ServicesOffered)

	err = s.db.QueryRowContext(ctx, `
        UPDATE transport_companies
        SET name = $1,
            slug = $2,
            country_code = $3,
            coverage_regions = $4,
            services_offered = $5,
            contact_email = $6,
            contact_phone = $7,
            website = $8,
            description = $9,
            active = $10,
            updated_at = NOW()
        WHERE id = $11
        RETURNING updated_at`,
		existing.Name,
		existing.Slug,
		existing.CountryCode,
		coverage,
		services,
		existing.ContactEmail,
		existing.ContactPhone,
		existing.Website,
		existing.Description,
		existing.Active,
		id,
	).Scan(&existing.UpdatedAt)
	if err != nil {
		return Company{}, err
	}

	return existing, nil
}

func (s *sqlService) Delete(ctx context.Context, id uuid.UUID) error {
    result, err := s.db.ExecContext(ctx, `DELETE FROM transport_companies WHERE id = $1`, id)
    if err != nil {
        return err
    }
    if rows, _ := result.RowsAffected(); rows == 0 {
        return ErrNotFound
    }
    return nil
}

func (s *sqlService) generateUniqueSlug(ctx context.Context, base string) (string, error) {
	if base == "" {
		base = "transport-company"
	}
	slug := base
	counter := 1
	for {
		var exists bool
		err := s.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM transport_companies WHERE slug = $1)`, slug).Scan(&exists)
		if err != nil {
			return "", err
		}
		if !exists {
			return slug, nil
		}
		counter++
		slug = fmt.Sprintf("%s-%d", base, counter)
	}
}

func scanCompany(scanner interface {
	Scan(dest ...any) error
}) (Company, error) {
	var company Company
	var coverageJSON, servicesJSON string
	var email, phone, website, description sql.NullString
	if err := scanner.Scan(
		&company.ID,
		&company.Name,
		&company.Slug,
		&company.CountryCode,
		&coverageJSON,
		&servicesJSON,
		&email,
		&phone,
		&website,
		&description,
		&company.Active,
		&company.CreatedAt,
		&company.UpdatedAt,
	); err != nil {
		return Company{}, err
	}
	if err := json.Unmarshal([]byte(coverageJSON), &company.CoverageRegions); err != nil {
		company.CoverageRegions = nil
	}
	if err := json.Unmarshal([]byte(servicesJSON), &company.ServicesOffered); err != nil {
		company.ServicesOffered = nil
	}
	company.ContactEmail = strings.TrimSpace(email.String)
	company.ContactPhone = strings.TrimSpace(phone.String)
	company.Website = strings.TrimSpace(website.String)
	company.Description = strings.TrimSpace(description.String)
	return company, nil
}

var _ Service = (*sqlService)(nil)
