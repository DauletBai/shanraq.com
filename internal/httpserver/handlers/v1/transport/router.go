package transport

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"shanraq.com/internal/config"
	transportservice "shanraq.com/internal/services/transport"
)

// Router configures routes for transportation logistics partners.
func Router(cfg config.Config, logger zerolog.Logger, service transportservice.Service) chi.Router {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		if !cfg.Features.EnableTransportCompanies {
			respondError(w, http.StatusNotFound, "feature_disabled")
			return
		}
		filter := transportservice.ListFilter{
			CountryCode: r.URL.Query().Get("country"),
			ActiveOnly:  r.URL.Query().Get("active") == "true",
		}
		if limit := r.URL.Query().Get("limit"); limit != "" {
			if v, err := strconv.Atoi(limit); err == nil && v >= 0 {
				filter.Limit = v
			}
		}
		if offset := r.URL.Query().Get("offset"); offset != "" {
			if v, err := strconv.Atoi(offset); err == nil && v >= 0 {
				filter.Offset = v
			}
		}

		items, total, err := service.List(r.Context(), filter)
		if err != nil {
			logger.Error().Err(err).Msg("list_transport_companies")
			respondError(w, http.StatusInternalServerError, "list_failed")
			return
		}

		resp := listResponse{
			Data: make([]companyResponse, 0, len(items)),
			Meta: listMeta{
				Total:   total,
				Limit:   filter.Limit,
				Offset:  filter.Offset,
				Country: filter.CountryCode,
			},
		}
		for _, c := range items {
			resp.Data = append(resp.Data, mapToResponse(c))
		}
		writeJSON(w, http.StatusOK, resp)
	})

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		if !cfg.Features.EnableTransportCompanies {
			respondError(w, http.StatusNotFound, "feature_disabled")
			return
		}
		var payload createRequest
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			respondError(w, http.StatusBadRequest, "invalid_payload")
			return
		}
		defer r.Body.Close()

		input := transportservice.CreateInput{
			Name:            payload.Name,
			CountryCode:     payload.CountryCode,
			CoverageRegions: payload.CoverageRegions,
			ServicesOffered: payload.ServicesOffered,
			ContactEmail:    payload.ContactEmail,
			ContactPhone:    payload.ContactPhone,
			Website:         payload.Website,
			Description:     payload.Description,
		}
		if payload.Active != nil {
			input.Active = *payload.Active
		} else {
			input.Active = true
		}

		company, err := service.Create(r.Context(), input)
		if err != nil {
			logger.Warn().Err(err).Msg("create_transport_company")
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}

		writeJSON(w, http.StatusCreated, mapToResponse(company))
	})

	r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		if !cfg.Features.EnableTransportCompanies {
			respondError(w, http.StatusNotFound, "feature_disabled")
			return
		}
		id, err := parseUUID(chi.URLParam(r, "id"))
		if err != nil {
			respondError(w, http.StatusBadRequest, "invalid_id")
			return
		}

		company, err := service.Get(r.Context(), id)
		if err != nil {
			if errors.Is(err, transportservice.ErrNotFound) {
				respondError(w, http.StatusNotFound, "not_found")
				return
			}
			logger.Error().Err(err).Str("id", id.String()).Msg("get_transport_company")
			respondError(w, http.StatusInternalServerError, "get_failed")
			return
		}
		writeJSON(w, http.StatusOK, mapToResponse(company))
	})

	r.Put("/{id}", func(w http.ResponseWriter, r *http.Request) {
		if !cfg.Features.EnableTransportCompanies {
			respondError(w, http.StatusNotFound, "feature_disabled")
			return
		}
		id, err := parseUUID(chi.URLParam(r, "id"))
		if err != nil {
			respondError(w, http.StatusBadRequest, "invalid_id")
			return
		}

		var payload updateRequest
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			respondError(w, http.StatusBadRequest, "invalid_payload")
			return
		}
		defer r.Body.Close()

		input := transportservice.UpdateInput{
			Name:         payload.Name,
			CountryCode:  payload.CountryCode,
			ContactEmail: payload.ContactEmail,
			ContactPhone: payload.ContactPhone,
			Website:      payload.Website,
			Description:  payload.Description,
			Active:       payload.Active,
		}
		if payload.CoverageRegions != nil {
			input.CoverageRegions = &payload.CoverageRegions
		}
		if payload.ServicesOffered != nil {
			input.ServicesOffered = &payload.ServicesOffered
		}

		company, err := service.Update(r.Context(), id, input)
		if err != nil {
			if errors.Is(err, transportservice.ErrNotFound) {
				respondError(w, http.StatusNotFound, "not_found")
				return
			}
			logger.Warn().Err(err).Str("id", id.String()).Msg("update_transport_company")
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}

		writeJSON(w, http.StatusOK, mapToResponse(company))
	})

	r.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {
		if !cfg.Features.EnableTransportCompanies {
			respondError(w, http.StatusNotFound, "feature_disabled")
			return
		}
		id, err := parseUUID(chi.URLParam(r, "id"))
		if err != nil {
			respondError(w, http.StatusBadRequest, "invalid_id")
			return
		}

		if err := service.Delete(r.Context(), id); err != nil {
			if errors.Is(err, transportservice.ErrNotFound) {
				respondError(w, http.StatusNotFound, "not_found")
				return
			}
			logger.Error().Err(err).Str("id", id.String()).Msg("delete_transport_company")
			respondError(w, http.StatusInternalServerError, "delete_failed")
			return
		}

		writeJSON(w, http.StatusNoContent, nil)
	})

	return r
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload == nil {
		return
	}
	_ = json.NewEncoder(w).Encode(payload)
}

func respondError(w http.ResponseWriter, status int, code string) {
	writeJSON(w, status, map[string]string{
		"error": code,
	})
}

func parseUUID(value string) (uuid.UUID, error) {
	return uuid.Parse(value)
}
