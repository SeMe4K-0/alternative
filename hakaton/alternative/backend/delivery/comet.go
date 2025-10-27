package delivery

import (
	"backend/apiutils"
	"backend/middleware"
	"backend/named_errors"
	"backend/usecase"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type CometHandler struct {
	cometUsecase *usecase.CometUsecase
}

func NewCometHandler(cometUsecase *usecase.CometUsecase) *CometHandler {
	return &CometHandler{cometUsecase: cometUsecase}
}

type CreateCometRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateCometRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// @Summary Create a new comet
// @Description Creates a new comet entry associated with the authenticated user.
// @Tags comets
// @Accept  json
// @Produce json
// @Param   comet body CreateCometRequest true "Comet information"
// @Success 201 {object} models.Comet
// @Failure 400 {object} apiutils.ErrorResponse "Invalid request body or name is required"
// @Failure 401 {object} apiutils.ErrorResponse "User not authenticated"
// @Failure 500 {object} apiutils.ErrorResponse "Internal server error"
// @Security ApiKeyAuth
// @Router /comets [post]
func (h *CometHandler) CreateComet(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		apiutils.WriteError(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	var req CreateCometRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apiutils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	comet, err := h.cometUsecase.CreateComet(userID, req.Name, req.Description)
	if err != nil {
		if errors.Is(err, named_errors.ErrInvalidInput) {
			apiutils.WriteError(w, http.StatusBadRequest, "invalid input: name is required")
			return
		}
		log.Error().Err(err).Msg("failed to create comet")
		apiutils.WriteError(w, http.StatusInternalServerError, "failed to create comet")
		return
	}

	apiutils.WriteJSON(w, http.StatusCreated, comet)
}

// @Summary Get a comet by ID
// @Description Retrieves a specific comet by its ID. User must be the owner.
// @Tags comets
// @Produce json
// @Param   id path int true "Comet ID"
// @Success 200 {object} models.Comet
// @Failure 400 {object} apiutils.ErrorResponse "Invalid comet ID"
// @Failure 401 {object} apiutils.ErrorResponse "User not authenticated"
// @Failure 403 {object} apiutils.ErrorResponse "Access denied"
// @Failure 404 {object} apiutils.ErrorResponse "Comet not found"
// @Security ApiKeyAuth
// @Router /comets/{id} [get]
func (h *CometHandler) GetComet(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		apiutils.WriteError(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	vars := mux.Vars(r)
	cometID, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		apiutils.WriteError(w, http.StatusBadRequest, "invalid comet ID")
		return
	}

	comet, err := h.cometUsecase.GetCometByID(userID, cometID)
	if err != nil {
		switch {
		case errors.Is(err, named_errors.ErrNotFound):
			apiutils.WriteError(w, http.StatusNotFound, "comet not found")
		case errors.Is(err, named_errors.ErrAccessDenied):
			apiutils.WriteError(w, http.StatusForbidden, "access denied")
		default:
			log.Error().Err(err).Msg("failed to get comet")
			apiutils.WriteError(w, http.StatusInternalServerError, "could not retrieve comet")
		}
		return
	}

	apiutils.WriteJSON(w, http.StatusOK, comet)
}

// @Summary List all comets for the current user
// @Description Retrieves a list of all comets owned by the authenticated user.
// @Tags comets
// @Produce json
// @Success 200 {array} models.Comet
// @Failure 401 {object} apiutils.ErrorResponse "User not authenticated"
// @Failure 500 {object} apiutils.ErrorResponse "Internal server error"
// @Security ApiKeyAuth
// @Router /comets [get]
func (h *CometHandler) ListUserComets(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		apiutils.WriteError(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	comets, err := h.cometUsecase.ListCometsByUserID(userID)
	if err != nil {
		log.Error().Err(err).Msg("failed to list comets")
		apiutils.WriteError(w, http.StatusInternalServerError, "could not retrieve comets")
		return
	}

	apiutils.WriteJSON(w, http.StatusOK, comets)
}

// @Summary Update a comet
// @Description Updates a comet's name and/or description. User must be the owner.
// @Tags comets
// @Accept  json
// @Produce json
// @Param   id    path int                true "Comet ID"
// @Param   comet body UpdateCometRequest true "Comet data to update"
// @Success 200 {object} models.Comet
// @Failure 400 {object} apiutils.ErrorResponse "Invalid comet ID or request body"
// @Failure 401 {object} apiutils.ErrorResponse "User not authenticated"
// @Failure 403 {object} apiutils.ErrorResponse "Access denied"
// @Failure 404 {object} apiutils.ErrorResponse "Comet not found"
// @Security ApiKeyAuth
// @Router /comets/{id} [put]
func (h *CometHandler) UpdateComet(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		apiutils.WriteError(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	vars := mux.Vars(r)
	cometID, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		apiutils.WriteError(w, http.StatusBadRequest, "invalid comet ID")
		return
	}

	var req UpdateCometRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apiutils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	comet, err := h.cometUsecase.UpdateComet(userID, cometID, req.Name, req.Description)
	if err != nil {
		switch {
		case errors.Is(err, named_errors.ErrNotFound):
			apiutils.WriteError(w, http.StatusNotFound, "comet not found")
		case errors.Is(err, named_errors.ErrAccessDenied):
			apiutils.WriteError(w, http.StatusForbidden, "access denied")
		default:
			log.Error().Err(err).Msg("failed to update comet")
			apiutils.WriteError(w, http.StatusInternalServerError, "could not update comet")
		}
		return
	}

	apiutils.WriteJSON(w, http.StatusOK, comet)
}

// @Summary Delete a comet
// @Description Deletes a comet by its ID. User must be the owner.
// @Tags comets
// @Produce json
// @Param   id path int true "Comet ID"
// @Success 200 {object} map[string]string "Success message"
// @Failure 400 {object} apiutils.ErrorResponse "Invalid comet ID"
// @Failure 401 {object} apiutils.ErrorResponse "User not authenticated"
// @Failure 403 {object} apiutils.ErrorResponse "Access denied"
// @Failure 404 {object} apiutils.ErrorResponse "Comet not found"
// @Security ApiKeyAuth
// @Router /comets/{id} [delete]
func (h *CometHandler) DeleteComet(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		apiutils.WriteError(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	vars := mux.Vars(r)
	cometID, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		apiutils.WriteError(w, http.StatusBadRequest, "invalid comet ID")
		return
	}

	if err := h.cometUsecase.DeleteComet(userID, cometID); err != nil {
		switch {
		case errors.Is(err, named_errors.ErrNotFound):
			apiutils.WriteError(w, http.StatusNotFound, "comet not found")
		case errors.Is(err, named_errors.ErrAccessDenied):
			apiutils.WriteError(w, http.StatusForbidden, "access denied")
		default:
			log.Error().Err(err).Msg("failed to delete comet")
			apiutils.WriteError(w, http.StatusInternalServerError, "could not delete comet")
		}
		return
	}

	apiutils.WriteJSON(w, http.StatusOK, map[string]string{"message": "comet deleted successfully"})
}