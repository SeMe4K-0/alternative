package delivery

import (
	"backend/apiutils"
	"backend/middleware"
	"backend/named_errors"
	"backend/usecase"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/rs/zerolog/log"
)

type CalculationHandler struct {
	calculationUsecase *usecase.CalculationUsecase
}

func NewCalculationHandler(uc *usecase.CalculationUsecase) *CalculationHandler {
	return &CalculationHandler{calculationUsecase: uc}
}

// --- Хендлер для РУЧКИ 1 ---
type CreateOrbitRequest struct {
	ObservationIDs []uint64 `json:"observation_ids"`
}

// @Summary (Step 1) Create a new orbit calculation
// @Description ...
// @Tags calculations
// @Accept  json
// @Produce json
// @Param   comet_id path int true "ID of the comet"
// @Param   observation_ids body CreateOrbitRequest true "List of observation IDs"
// @Success 201 {object} models.OrbitalCalculation "Returns calculation with 6 orbit params, but empty approach data"
// @Router /comets/{comet_id}/calculations [post]
func (h *CalculationHandler) CreateOrbitCalculation(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		apiutils.WriteError(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	cometID, err := parseID(r, "comet_id")
	if err != nil {
		apiutils.WriteError(w, http.StatusBadRequest, "invalid comet ID")
		return
	}

	var req CreateOrbitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apiutils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	calc, err := h.calculationUsecase.CreateOrbitCalculation(userID, cometID, req.ObservationIDs)
	if err != nil {
		apiutils.WriteError(w, http.StatusInternalServerError, "failed to trigger calculation")
		return
	}
	apiutils.WriteJSON(w, http.StatusCreated, calc)
}

// --- Хендлер для РУЧКИ 2 ---

// @Summary (Step 2) Calculate and add close approach data to an existing orbit
// @Description Takes 6 orbital parameters, calculates close approach and UPDATES an existing calculation record with the results.
// @Tags calculations
// @Accept  json
// @Produce json
// @Param   calculation_id path int true "ID of the existing Orbit Calculation to update"
// @Param   orbit_parameters body clients.ApproachCalculationRequest true "The 6 orbital parameters"
// @Success 200 {object} models.OrbitalCalculation "Returns the FULL calculation object, now updated with approach data"
// @Router /calculations/{calculation_id}/approach [post]
func (h *CalculationHandler) AddCloseApproachData(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		apiutils.WriteError(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	calcID, err := parseID(r, "calculation_id")
	if err != nil {
		apiutils.WriteError(w, http.StatusBadRequest, "invalid calculation ID")
		return
	}

	updatedCalc, err := h.calculationUsecase.UpdateWithCloseApproach(userID, calcID)
	if err != nil {
		apiutils.WriteError(w, http.StatusInternalServerError, "failed to calculate and update close approach")
		return
	}
	apiutils.WriteJSON(w, http.StatusOK, updatedCalc)
}

// @Summary List calculation history for a comet
// @Description Retrieves the history of all orbital calculations performed for a specific comet.
// @Tags calculations
// @Produce json
// @Param   comet_id path int true "ID of the comet"
// @Success 200 {array} models.OrbitalCalculation
// @Failure 400 {object} apiutils.ErrorResponse "Invalid comet ID"
// @Failure 401 {object} apiutils.ErrorResponse "User not authenticated"
// @Failure 403 {object} apiutils.ErrorResponse "Access denied to comet"
// @Failure 404 {object} apiutils.ErrorResponse "Comet not found"
// @Security ApiKeyAuth
// @Router /comets/{comet_id}/calculations [get]
func (h *CalculationHandler) ListCalculationsForComet(w http.ResponseWriter, r *http.Request) {
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

	calcs, err := h.calculationUsecase.ListCalculationsByComet(userID, cometID)
	if err != nil {
		switch {
		case errors.Is(err, named_errors.ErrNotFound):
			apiutils.WriteError(w, http.StatusNotFound, "comet not found")
		case errors.Is(err, named_errors.ErrAccessDenied):
			apiutils.WriteError(w, http.StatusForbidden, "access denied to comet")
		default:
			log.Error().Err(err).Msg("failed to list calculations")
			apiutils.WriteError(w, http.StatusInternalServerError, "could not retrieve calculations")
		}
		return
	}

	apiutils.WriteJSON(w, http.StatusOK, calcs)
}

// @Summary Get a specific calculation result
// @Description Retrieves a single orbital calculation result by its ID.
// @Tags calculations
// @Produce json
// @Param   calculation_id path int true "ID of the calculation"
// @Success 200 {object} models.OrbitalCalculation
// @Failure 400 {object} apiutils.ErrorResponse "Invalid calculation ID"
// @Failure 401 {object} apiutils.ErrorResponse "User not authenticated"
// @Failure 403 {object} apiutils.ErrorResponse "Access denied"
// @Failure 404 {object} apiutils.ErrorResponse "Calculation not found"
// @Security ApiKeyAuth
// @Router /calculations/{calculation_id} [get]
func (h *CalculationHandler) GetCalculation(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		apiutils.WriteError(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	calcID, err := parseID(r, "calculation_id")
	if err != nil {
		apiutils.WriteError(w, http.StatusBadRequest, "invalid calculation ID in URL")
		return
	}

	calc, err := h.calculationUsecase.GetCalculationByID(userID, calcID)
	if err != nil {
		switch {
		case errors.Is(err, named_errors.ErrNotFound):
			apiutils.WriteError(w, http.StatusNotFound, "calculation not found")
		case errors.Is(err, named_errors.ErrAccessDenied):
			apiutils.WriteError(w, http.StatusForbidden, "access denied")
		default:
			log.Error().Err(err).Msg("failed to get calculation")
			apiutils.WriteError(w, http.StatusInternalServerError, "could not retrieve calculation")
		}
		return
	}

	apiutils.WriteJSON(w, http.StatusOK, calc)
}

// @Summary Delete a calculation result
// @Description Deletes a calculation result by its ID.
// @Tags calculations
// @Produce json
// @Param   calculation_id path int true "ID of the calculation to delete"
// @Success 200 {object} map[string]string "Success message"
// @Failure 400 {object} apiutils.ErrorResponse "Invalid calculation ID"
// @Failure 401 {object} apiutils.ErrorResponse "User not authenticated"
// @Failure 403 {object} apiutils.ErrorResponse "Access denied"
// @Failure 404 {object} apiutils.ErrorResponse "Calculation not found"
// @Security ApiKeyAuth
// @Router /calculations/{calculation_id} [delete]
func (h *CalculationHandler) DeleteCalculation(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		apiutils.WriteError(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	calcID, err := parseID(r, "calculation_id")
	if err != nil {
		apiutils.WriteError(w, http.StatusBadRequest, "invalid calculation ID in URL")
		return
	}

	if err := h.calculationUsecase.DeleteCalculation(userID, calcID); err != nil {
		switch {
		case errors.Is(err, named_errors.ErrNotFound):
			apiutils.WriteError(w, http.StatusNotFound, "calculation not found")
		case errors.Is(err, named_errors.ErrAccessDenied):
			apiutils.WriteError(w, http.StatusForbidden, "access denied")
		default:
			log.Error().Err(err).Msg("failed to delete calculation")
			apiutils.WriteError(w, http.StatusInternalServerError, "could not delete calculation")
		}
		return
	}

	apiutils.WriteJSON(w, http.StatusOK, map[string]string{"message": "calculation deleted successfully"})
}
