package delivery

import (
	"backend/middleware"
	"backend/usecase"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestNewObservationHandler(t *testing.T) {
	observationUsecase := &usecase.ObservationUsecase{}
	
	handler := NewObservationHandler(observationUsecase)

	assert.NotNil(t, handler)
	assert.Equal(t, observationUsecase, handler.observationUsecase)
}

func TestObservationHandler_CreateObservationForComet_NoUserID(t *testing.T) {
	handler := NewObservationHandler(&usecase.ObservationUsecase{})

	reqData := map[string]interface{}{
		"observed_at": time.Now().Format(time.RFC3339),
		"ra":          1.0,
		"dec":         2.0,
		"notes":       "Test observation",
	}
	reqBody, _ := json.Marshal(reqData)

	req := httptest.NewRequest("POST", "/comets/1/observations", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.Background()) // No user ID in context
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.CreateObservationForComet(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestObservationHandler_CreateObservationForComet_InvalidJSON(t *testing.T) {
	handler := NewObservationHandler(&usecase.ObservationUsecase{})

	req := httptest.NewRequest("POST", "/comets/1/observations", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(middleware.WithUserID(context.Background(), uint64(1)))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.CreateObservationForComet(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestObservationHandler_CreateObservationForComet_ValidData(t *testing.T) {
	handler := NewObservationHandler(&usecase.ObservationUsecase{})

	reqData := map[string]interface{}{
		"observed_at": time.Now().Format(time.RFC3339),
		"ra":          1.0,
		"dec":         2.0,
		"notes":       "Test observation",
	}
	reqBody, _ := json.Marshal(reqData)

	req := httptest.NewRequest("POST", "/comets/1/observations", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(middleware.WithUserID(context.Background(), uint64(1)))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.CreateObservationForComet(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestObservationHandler_ListObservationsForComet_NoUserID(t *testing.T) {
	handler := NewObservationHandler(&usecase.ObservationUsecase{})

	req := httptest.NewRequest("GET", "/comets/1/observations", nil)
	req = req.WithContext(context.Background()) // No user ID in context
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.ListObservationsForComet(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestObservationHandler_ListObservationsForComet_ValidUserID(t *testing.T) {
	handler := NewObservationHandler(&usecase.ObservationUsecase{})

	req := httptest.NewRequest("GET", "/comets/1/observations", nil)
	req = req.WithContext(middleware.WithUserID(context.Background(), uint64(1)))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.ListObservationsForComet(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestObservationHandler_GetObservation_NoUserID(t *testing.T) {
	handler := NewObservationHandler(&usecase.ObservationUsecase{})

	req := httptest.NewRequest("GET", "/comets/1/observations/1", nil)
	req = req.WithContext(context.Background()) // No user ID in context
	req = mux.SetURLVars(req, map[string]string{"id": "1", "observation_id": "1"})
	w := httptest.NewRecorder()

	handler.GetObservation(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestObservationHandler_GetObservation_ValidUserID(t *testing.T) {
	handler := NewObservationHandler(&usecase.ObservationUsecase{})

	req := httptest.NewRequest("GET", "/comets/1/observations/1", nil)
	req = req.WithContext(middleware.WithUserID(context.Background(), uint64(1)))
	req = mux.SetURLVars(req, map[string]string{"id": "1", "observation_id": "1"})
	w := httptest.NewRecorder()

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil repository")
		}
	}()
	
	handler.GetObservation(w, req)
}

func TestObservationHandler_UpdateObservation_NoUserID(t *testing.T) {
	handler := NewObservationHandler(&usecase.ObservationUsecase{})

	reqData := map[string]interface{}{
		"observed_at": time.Now().Format(time.RFC3339),
		"ra":          1.5,
		"dec":         2.5,
		"notes":       "Updated observation",
	}
	reqBody, _ := json.Marshal(reqData)

	req := httptest.NewRequest("PUT", "/comets/1/observations/1", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.Background()) // No user ID in context
	req = mux.SetURLVars(req, map[string]string{"id": "1", "observation_id": "1"})
	w := httptest.NewRecorder()

	handler.UpdateObservation(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestObservationHandler_UpdateObservation_InvalidJSON(t *testing.T) {
	handler := NewObservationHandler(&usecase.ObservationUsecase{})

	req := httptest.NewRequest("PUT", "/comets/1/observations/1", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(middleware.WithUserID(context.Background(), uint64(1)))
	req = mux.SetURLVars(req, map[string]string{"id": "1", "observation_id": "1"})
	w := httptest.NewRecorder()

	handler.UpdateObservation(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestObservationHandler_UpdateObservation_ValidData(t *testing.T) {
	handler := NewObservationHandler(&usecase.ObservationUsecase{})

	reqData := map[string]interface{}{
		"observed_at": time.Now().Format(time.RFC3339),
		"ra":          1.5,
		"dec":         2.5,
		"notes":       "Updated observation",
	}
	reqBody, _ := json.Marshal(reqData)

	req := httptest.NewRequest("PUT", "/comets/1/observations/1", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(middleware.WithUserID(context.Background(), uint64(1)))
	req = mux.SetURLVars(req, map[string]string{"id": "1", "observation_id": "1"})
	w := httptest.NewRecorder()

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil repository")
		}
	}()
	
	handler.UpdateObservation(w, req)
}

func TestObservationHandler_DeleteObservation_NoUserID(t *testing.T) {
	handler := NewObservationHandler(&usecase.ObservationUsecase{})

	req := httptest.NewRequest("DELETE", "/comets/1/observations/1", nil)
	req = req.WithContext(context.Background()) // No user ID in context
	req = mux.SetURLVars(req, map[string]string{"id": "1", "observation_id": "1"})
	w := httptest.NewRecorder()

	handler.DeleteObservation(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestObservationHandler_DeleteObservation_ValidUserID(t *testing.T) {
	handler := NewObservationHandler(&usecase.ObservationUsecase{})

	req := httptest.NewRequest("DELETE", "/comets/1/observations/1", nil)
	req = req.WithContext(middleware.WithUserID(context.Background(), uint64(1)))
	req = mux.SetURLVars(req, map[string]string{"id": "1", "observation_id": "1"})
	w := httptest.NewRecorder()

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil repository")
		}
	}()
	
	handler.DeleteObservation(w, req)
}

func TestParseID(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "123"})
	
	id, err := parseID(req, "id")
	assert.NoError(t, err)
	assert.Equal(t, uint64(123), id)
}

func TestParseID_InvalidID(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "invalid"})
	
	_, err := parseID(req, "id")
	assert.Error(t, err)
}

func TestParseID_MissingID(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	req = mux.SetURLVars(req, map[string]string{})
	
	_, err := parseID(req, "id")
	assert.Error(t, err)
}

func TestObservationHandler_CreateObservationForComet_ValidData_WithPanic(t *testing.T) {
	handler := NewObservationHandler(&usecase.ObservationUsecase{})

	reqData := map[string]interface{}{
		"observed_at": time.Now().Format(time.RFC3339),
		"ra":          1.0,
		"dec":         2.0,
		"notes":       "Test observation",
	}
	reqBody, _ := json.Marshal(reqData)

	req := httptest.NewRequest("POST", "/comets/1/observations", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(middleware.WithUserID(context.Background(), uint64(1)))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.CreateObservationForComet(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestObservationHandler_ListObservationsForComet_ValidData_WithPanic(t *testing.T) {
	handler := NewObservationHandler(&usecase.ObservationUsecase{})

	req := httptest.NewRequest("GET", "/comets/1/observations", nil)
	req = req.WithContext(middleware.WithUserID(context.Background(), uint64(1)))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.ListObservationsForComet(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestObservationHandler_UpdateObservation_ValidData_WithPanic(t *testing.T) {
	handler := NewObservationHandler(&usecase.ObservationUsecase{})

	reqData := map[string]interface{}{
		"observed_at": time.Now().Format(time.RFC3339),
		"ra":          1.5,
		"dec":         2.5,
		"notes":       "Updated observation",
	}
	reqBody, _ := json.Marshal(reqData)

	req := httptest.NewRequest("PUT", "/comets/1/observations/1", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(middleware.WithUserID(context.Background(), uint64(1)))
	req = mux.SetURLVars(req, map[string]string{"id": "1", "observation_id": "1"})
	w := httptest.NewRecorder()

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil repository")
		}
	}()
	
	handler.UpdateObservation(w, req)
}

func TestObservationHandler_DeleteObservation_ValidData_WithPanic(t *testing.T) {
	handler := NewObservationHandler(&usecase.ObservationUsecase{})

	req := httptest.NewRequest("DELETE", "/comets/1/observations/1", nil)
	req = req.WithContext(middleware.WithUserID(context.Background(), uint64(1)))
	req = mux.SetURLVars(req, map[string]string{"id": "1", "observation_id": "1"})
	w := httptest.NewRecorder()

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil repository")
		}
	}()
	
	handler.DeleteObservation(w, req)
}
