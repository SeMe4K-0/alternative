package delivery

import (
	"backend/apiutils"
	"backend/middleware"
	"backend/models"
	"backend/named_errors"
	"backend/usecase"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type ObservationHandler struct {
	observationUsecase *usecase.ObservationUsecase
}

func NewObservationHandler(uc *usecase.ObservationUsecase) *ObservationHandler {
	return &ObservationHandler{observationUsecase: uc}
}

type CreateObservationRequest struct {
	ObservedAt string  `json:"observed_at"`
	RA         float64 `json:"right_ascension"`
	Dec        float64 `json:"declination"`
	Notes      string  `json:"notes,omitempty"`
}

type UpdateObservationRequest struct {
	ObservedAt string  `json:"observed_at"`
	RA         float64 `json:"right_ascension"`
	Dec        float64 `json:"declination"`
	Notes      string  `json:"notes,omitempty"`
	ImageURL   string  `json:"image_url,omitempty"`
}

// parseCometID извлекает ID кометы из URL.
func parseID(r *http.Request, key string) (uint64, error) {
	vars := mux.Vars(r)
	idStr, ok := vars[key]
	if !ok {
		return 0, errors.New("missing id in URL path: " + key)
	}
	return strconv.ParseUint(idStr, 10, 64)
}

// @Summary Create an observation for a comet
// @Description Adds a new observation record to a specific comet. Can accept JSON or multipart/form-data for image uploads.
// @Tags observations
// @Accept  json
// @Accept  mpfd
// @Produce json
// @Param   comet_id        path     int    true  "ID of the comet to add observation to"
// @Param   observation     body     CreateObservationRequest false "Observation data (if using JSON)"
// @Param   observed_at     formData string true  "Observation timestamp in RFC3339 format (if using multipart)"
// @Param   right_ascension formData number true  "Right Ascension in degrees (if using multipart)"
// @Param   declination     formData number true  "Declination in degrees (if using multipart)"
// @Param   notes           formData string false "Notes about the observation (if using multipart)"
// @Param   image           formData file   false "Image of the observation (if using multipart)"
// @Success 201 {object} models.Observation
// @Failure 400 {object} apiutils.ErrorResponse "Invalid request data or comet ID"
// @Failure 401 {object} apiutils.ErrorResponse "User not authenticated"
// @Failure 403 {object} apiutils.ErrorResponse "Access denied to parent comet"
// @Failure 404 {object} apiutils.ErrorResponse "Parent comet not found"
// @Security ApiKeyAuth
// @Router /comets/{comet_id}/observations [post]
func (h *ObservationHandler) CreateObservationForComet(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		apiutils.WriteError(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	cometID, err := parseID(r, "comet_id")
	if err != nil {
		apiutils.WriteError(w, http.StatusBadRequest, "invalid comet ID in URL")
		return
	}

	contentType := r.Header.Get("Content-Type")

	var obs *models.Observation

	if strings.Contains(contentType, "multipart/form-data") {
		err = r.ParseMultipartForm(10 << 20)
		if err != nil {
			log.Error().Err(err).Msg("failed to parse multipart form")
			apiutils.WriteError(w, http.StatusBadRequest, "failed to parse multipart form")
			return
		}

		observedAtStr := r.FormValue("observed_at")
		raStr := r.FormValue("right_ascension")
		decStr := r.FormValue("declination")
		notes := r.FormValue("notes")

		if observedAtStr == "" {
			apiutils.WriteError(w, http.StatusBadRequest, "observed_at is required")
			return
		}
		if raStr == "" {
			apiutils.WriteError(w, http.StatusBadRequest, "right_ascension is required")
			return
		}
		if decStr == "" {
			apiutils.WriteError(w, http.StatusBadRequest, "declination is required")
			return
		}

		observedAt, err := time.Parse(time.RFC3339, observedAtStr)
		if err != nil {
			log.Error().Err(err).Str("observed_at", observedAtStr).Msg("invalid date format")
			apiutils.WriteError(w, http.StatusBadRequest, "invalid date format for observedAt, use RFC3339")
			return
		}

		ra, err := strconv.ParseFloat(raStr, 64)
		if err != nil {
			log.Error().Err(err).Str("right_ascension", raStr).Msg("invalid right_ascension value")
			apiutils.WriteError(w, http.StatusBadRequest, "invalid right_ascension value")
			return
		}

		dec, err := strconv.ParseFloat(decStr, 64)
		if err != nil {
			log.Error().Err(err).Str("declination", decStr).Msg("invalid declination value")
			apiutils.WriteError(w, http.StatusBadRequest, "invalid declination value")
			return
		}

		obs, err = h.observationUsecase.CreateObservation(userID, cometID, observedAt, ra, dec, notes)
		if err != nil {
			switch {
			case errors.Is(err, named_errors.ErrNotFound):
				apiutils.WriteError(w, http.StatusNotFound, "parent comet not found")
			case errors.Is(err, named_errors.ErrAccessDenied):
				apiutils.WriteError(w, http.StatusForbidden, "access denied to parent comet")
			default:
				log.Error().Err(err).Msg("failed to create observation")
				apiutils.WriteError(w, http.StatusInternalServerError, "failed to create observation")
			}
			return
		}

		file, header, err := r.FormFile("image")
		if err == nil {
			defer file.Close()

			if header.Size > 10<<20 {
				apiutils.WriteError(w, http.StatusBadRequest, "file too large, maximum size is 10MB")
				return
			}

			fileContentType := header.Header.Get("Content-Type")
			if fileContentType == "" {
				fileContentType = "application/octet-stream"
			}

			obs, err = h.observationUsecase.UploadObservationImage(userID, obs.ID, file, header.Filename, fileContentType)
			if err != nil {
				switch {
				case errors.Is(err, named_errors.ErrInvalidInput):
					apiutils.WriteError(w, http.StatusBadRequest, "invalid file type, only images are allowed")
				default:
					log.Error().Err(err).Msg("failed to upload observation image")
					apiutils.WriteError(w, http.StatusInternalServerError, "failed to upload image")
				}
				return
			}
		}
	} else {
		var req CreateObservationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			apiutils.WriteError(w, http.StatusBadRequest, "invalid request body")
			log.Error().Err(err).Msg("failed to decode create observation request")
			return
		}

		observedAt, err := time.Parse(time.RFC3339, req.ObservedAt)
		if err != nil {
			apiutils.WriteError(w, http.StatusBadRequest, "invalid date format for observed_at, use RFC3339")
			return
		}

		obs, err = h.observationUsecase.CreateObservation(userID, cometID, observedAt, req.RA, req.Dec, req.Notes)
		if err != nil {
			switch {
			case errors.Is(err, named_errors.ErrNotFound):
				apiutils.WriteError(w, http.StatusNotFound, "parent comet not found")
			case errors.Is(err, named_errors.ErrAccessDenied):
				apiutils.WriteError(w, http.StatusForbidden, "access denied to parent comet")
			default:
				log.Error().Err(err).Msg("failed to create observation")
				apiutils.WriteError(w, http.StatusInternalServerError, "failed to create observation")
			}
			return
		}
	}

	apiutils.WriteJSON(w, http.StatusCreated, obs)
}

// @Summary List all observations for a comet
// @Description Retrieves all observation records for a specific comet.
// @Tags observations
// @Produce json
// @Param   comet_id path int true "ID of the comet"
// @Success 200 {array} models.Observation
// @Failure 400 {object} apiutils.ErrorResponse "Invalid comet ID"
// @Failure 401 {object} apiutils.ErrorResponse "User not authenticated"
// @Failure 403 {object} apiutils.ErrorResponse "Access denied to comet"
// @Failure 404 {object} apiutils.ErrorResponse "Comet not found"
// @Security ApiKeyAuth
// @Router /comets/{comet_id}/observations [get]
func (h *ObservationHandler) ListObservationsForComet(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		apiutils.WriteError(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	cometID, err := parseID(r, "comet_id")
	if err != nil {
		apiutils.WriteError(w, http.StatusBadRequest, "invalid comet ID in URL")
		return
	}

	observations, err := h.observationUsecase.ListObservationsByComet(userID, cometID)
	if err != nil {
		switch {
		case errors.Is(err, named_errors.ErrNotFound):
			apiutils.WriteError(w, http.StatusNotFound, "comet not found")
		case errors.Is(err, named_errors.ErrAccessDenied):
			apiutils.WriteError(w, http.StatusForbidden, "access denied to comet")
		default:
			log.Error().Err(err).Msg("failed to list observations")
			apiutils.WriteError(w, http.StatusInternalServerError, "could not retrieve observations")
		}
		return
	}

	apiutils.WriteJSON(w, http.StatusOK, observations)
}

// @Summary Get a single observation
// @Description Retrieves a single observation by its ID.
// @Tags observations
// @Produce json
// @Param   observation_id path int true "ID of the observation"
// @Success 200 {object} models.Observation
// @Failure 400 {object} apiutils.ErrorResponse "Invalid observation ID"
// @Failure 401 {object} apiutils.ErrorResponse "User not authenticated"
// @Failure 403 {object} apiutils.ErrorResponse "Access denied"
// @Failure 404 {object} apiutils.ErrorResponse "Observation not found"
// @Security ApiKeyAuth
// @Router /observations/{observation_id} [get]
func (h *ObservationHandler) GetObservation(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		apiutils.WriteError(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	observationID, err := parseID(r, "observation_id")
	if err != nil {
		apiutils.WriteError(w, http.StatusBadRequest, "invalid observation ID in URL")
		return
	}

	obs, err := h.observationUsecase.GetObservationByID(userID, observationID)
	if err != nil {
		switch {
		case errors.Is(err, named_errors.ErrNotFound):
			apiutils.WriteError(w, http.StatusNotFound, "observation not found")
		case errors.Is(err, named_errors.ErrAccessDenied):
			apiutils.WriteError(w, http.StatusForbidden, "access denied")
		default:
			log.Error().Err(err).Msg("failed to get observation")
			apiutils.WriteError(w, http.StatusInternalServerError, "could not retrieve observation")
		}
		return
	}

	apiutils.WriteJSON(w, http.StatusOK, obs)
}

// @Summary Update an observation
// @Description Updates an existing observation record.
// @Tags observations
// @Accept  json
// @Accept  mpfd
// @Produce json
// @Param   observation_id  path     int    true  "ID of the observation to update"
// @Param   observation     body     UpdateObservationRequest false "Updated observation data (if using JSON)"
// @Param   observed_at     formData string true  "Updated timestamp in RFC3339 format (if using multipart)"
// @Param   right_ascension formData number true  "Updated Right Ascension in degrees (if using multipart)"
// @Param   declination     formData number true  "Updated Declination in degrees (if using multipart)"
// @Param   notes           formData string false "Updated notes (if using multipart)"
// @Param   image           formData file   false "New image for the observation (if using multipart)"
// @Success 200 {object} models.Observation
// @Failure 400 {object} apiutils.ErrorResponse "Invalid request data or observation ID"
// @Failure 401 {object} apiutils.ErrorResponse "User not authenticated"
// @Failure 403 {object} apiutils.ErrorResponse "Access denied"
// @Failure 404 {object} apiutils.ErrorResponse "Observation not found"
// @Security ApiKeyAuth
// @Router /observations/{observation_id} [put]
func (h *ObservationHandler) UpdateObservation(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		apiutils.WriteError(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	observationID, err := parseID(r, "observation_id")
	if err != nil {
		apiutils.WriteError(w, http.StatusBadRequest, "invalid observation ID in URL")
		return
	}

	contentType := r.Header.Get("Content-Type")

	var obs *models.Observation

	if strings.Contains(contentType, "multipart/form-data") {
		err = r.ParseMultipartForm(10 << 20)
		if err != nil {
			apiutils.WriteError(w, http.StatusBadRequest, "failed to parse multipart form")
			return
		}

		observedAtStr := r.FormValue("observed_at")
		raStr := r.FormValue("right_ascension")
		decStr := r.FormValue("declination")
		notes := r.FormValue("notes")

		if observedAtStr == "" {
			apiutils.WriteError(w, http.StatusBadRequest, "observed_at is required")
			return
		}
		if raStr == "" {
			apiutils.WriteError(w, http.StatusBadRequest, "right_ascension is required")
			return
		}
		if decStr == "" {
			apiutils.WriteError(w, http.StatusBadRequest, "declination is required")
			return
		}

		observedAt, err := time.Parse(time.RFC3339, observedAtStr)
		if err != nil {
			apiutils.WriteError(w, http.StatusBadRequest, "invalid date format for observed_at, use RFC3339")
			return
		}

		ra, err := strconv.ParseFloat(raStr, 64)
		if err != nil {
			apiutils.WriteError(w, http.StatusBadRequest, "invalid right_ascension value")
			return
		}

		dec, err := strconv.ParseFloat(decStr, 64)
		if err != nil {
			apiutils.WriteError(w, http.StatusBadRequest, "invalid declination value")
			return
		}

		updateData := models.Observation{
			ObservedAt: observedAt,
			RA:         ra,
			Dec:        dec,
			Notes:      notes,
		}

		obs, err = h.observationUsecase.UpdateObservation(userID, observationID, updateData)
		if err != nil {
			switch {
			case errors.Is(err, named_errors.ErrNotFound):
				apiutils.WriteError(w, http.StatusNotFound, "observation not found")
			case errors.Is(err, named_errors.ErrAccessDenied):
				apiutils.WriteError(w, http.StatusForbidden, "access denied")
			default:
				log.Error().Err(err).Msg("failed to update observation")
				apiutils.WriteError(w, http.StatusInternalServerError, "could not update observation")
			}
			return
		}

		file, header, err := r.FormFile("image")
		if err == nil {
			defer file.Close()

			if header.Size > 10<<20 {
				apiutils.WriteError(w, http.StatusBadRequest, "file too large, maximum size is 10MB")
				return
			}

			fileContentType := header.Header.Get("Content-Type")
			if fileContentType == "" {
				fileContentType = "application/octet-stream"
			}

			obs, err = h.observationUsecase.UploadObservationImage(userID, observationID, file, header.Filename, fileContentType)
			if err != nil {
				switch {
				case errors.Is(err, named_errors.ErrInvalidInput):
					apiutils.WriteError(w, http.StatusBadRequest, "invalid file type, only images are allowed")
				default:
					log.Error().Err(err).Msg("failed to upload observation image")
					apiutils.WriteError(w, http.StatusInternalServerError, "failed to upload image")
				}
				return
			}
		}
	} else {
		var req UpdateObservationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			apiutils.WriteError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		observedAt, err := time.Parse(time.RFC3339, req.ObservedAt)
		if err != nil {
			apiutils.WriteError(w, http.StatusBadRequest, "invalid date format for observed_at, use RFC3339")
			return
		}

		updateData := models.Observation{
			ObservedAt: observedAt,
			RA:         req.RA,
			Dec:        req.Dec,
			Notes:      req.Notes,
			ImageURL:   req.ImageURL,
		}

		obs, err = h.observationUsecase.UpdateObservation(userID, observationID, updateData)
		if err != nil {
			switch {
			case errors.Is(err, named_errors.ErrNotFound):
				apiutils.WriteError(w, http.StatusNotFound, "observation not found")
			case errors.Is(err, named_errors.ErrAccessDenied):
				apiutils.WriteError(w, http.StatusForbidden, "access denied")
			default:
				log.Error().Err(err).Msg("failed to update observation")
				apiutils.WriteError(w, http.StatusInternalServerError, "could not update observation")
			}
			return
		}
	}

	apiutils.WriteJSON(w, http.StatusOK, obs)
}

// @Summary Delete an observation
// @Description Deletes an observation by its ID.
// @Tags observations
// @Produce json
// @Param   observation_id path int true "ID of the observation to delete"
// @Success 200 {object} map[string]string "Success message"
// @Failure 400 {object} apiutils.ErrorResponse "Invalid observation ID"
// @Failure 401 {object} apiutils.ErrorResponse "User not authenticated"
// @Failure 403 {object} apiutils.ErrorResponse "Access denied"
// @Failure 404 {object} apiutils.ErrorResponse "Observation not found"
// @Security ApiKeyAuth
// @Router /observations/{observation_id} [delete]
func (h *ObservationHandler) DeleteObservation(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		apiutils.WriteError(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	observationID, err := parseID(r, "observation_id")
	if err != nil {
		apiutils.WriteError(w, http.StatusBadRequest, "invalid observation ID in URL")
		return
	}

	if err := h.observationUsecase.DeleteObservation(userID, observationID); err != nil {
		switch {
		case errors.Is(err, named_errors.ErrNotFound):
			apiutils.WriteError(w, http.StatusNotFound, "observation not found")
		case errors.Is(err, named_errors.ErrAccessDenied):
			apiutils.WriteError(w, http.StatusForbidden, "access denied")
		default:
			log.Error().Err(err).Msg("failed to delete observation")
			apiutils.WriteError(w, http.StatusInternalServerError, "could not delete observation")
		}
		return
	}

	apiutils.WriteJSON(w, http.StatusOK, map[string]string{"message": "observation deleted successfully"})
}
