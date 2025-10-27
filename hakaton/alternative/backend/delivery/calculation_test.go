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

func TestNewCalculationHandler(t *testing.T) {
	calculationUsecase := &usecase.CalculationUsecase{}
	
	handler := NewCalculationHandler(calculationUsecase)

	assert.NotNil(t, handler)
	assert.Equal(t, calculationUsecase, handler.calculationUsecase)
}

func TestCalculationHandler_CreateOrbitCalculation_NoUserID(t *testing.T) {
	handler := NewCalculationHandler(&usecase.CalculationUsecase{})

	reqData := map[string]interface{}{
		"observations": []map[string]interface{}{
			{
				"observed_at": time.Now().Format(time.RFC3339),
				"ra":          1.0,
				"dec":         2.0,
			},
		},
	}
	reqBody, _ := json.Marshal(reqData)

	req := httptest.NewRequest("POST", "/comets/1/calculations/orbit", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.Background()) // No user ID in context
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.CreateOrbitCalculation(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestCalculationHandler_CreateOrbitCalculation_InvalidJSON(t *testing.T) {
	handler := NewCalculationHandler(&usecase.CalculationUsecase{})

	req := httptest.NewRequest("POST", "/comets/1/calculations/orbit", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(middleware.WithUserID(context.Background(), uint64(1)))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.CreateOrbitCalculation(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCalculationHandler_CreateOrbitCalculation_ValidData(t *testing.T) {
	handler := NewCalculationHandler(&usecase.CalculationUsecase{})

	reqData := map[string]interface{}{
		"observations": []map[string]interface{}{
			{
				"observed_at": time.Now().Format(time.RFC3339),
				"ra":          1.0,
				"dec":         2.0,
			},
		},
	}
	reqBody, _ := json.Marshal(reqData)

	req := httptest.NewRequest("POST", "/comets/1/calculations/orbit", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(middleware.WithUserID(context.Background(), uint64(1)))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.CreateOrbitCalculation(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCalculationHandler_AddCloseApproachData_NoUserID(t *testing.T) {
	handler := NewCalculationHandler(&usecase.CalculationUsecase{})

	reqData := map[string]interface{}{
		"semi_major_axis":     1.5,
		"eccentricity":        0.1,
		"inclination":         0.2,
		"lon_ascending_node":  0.3,
		"arg_periapsis":       0.4,
		"time_perihelion":     time.Now().Format(time.RFC3339),
	}
	reqBody, _ := json.Marshal(reqData)

	req := httptest.NewRequest("POST", "/comets/1/calculations/1/close-approach", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.Background()) // No user ID in context
	req = mux.SetURLVars(req, map[string]string{"id": "1", "calculation_id": "1"})
	w := httptest.NewRecorder()

	handler.AddCloseApproachData(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestCalculationHandler_AddCloseApproachData_InvalidJSON(t *testing.T) {
	handler := NewCalculationHandler(&usecase.CalculationUsecase{})

	req := httptest.NewRequest("POST", "/comets/1/calculations/1/close-approach", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(middleware.WithUserID(context.Background(), uint64(1)))
	req = mux.SetURLVars(req, map[string]string{"id": "1", "calculation_id": "1"})
	w := httptest.NewRecorder()

	handler.AddCloseApproachData(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCalculationHandler_AddCloseApproachData_ValidData(t *testing.T) {
	handler := NewCalculationHandler(&usecase.CalculationUsecase{})

	reqData := map[string]interface{}{
		"semi_major_axis":     1.5,
		"eccentricity":        0.1,
		"inclination":         0.2,
		"lon_ascending_node":  0.3,
		"arg_periapsis":       0.4,
		"time_perihelion":     time.Now().Format(time.RFC3339),
	}
	reqBody, _ := json.Marshal(reqData)

	req := httptest.NewRequest("POST", "/comets/1/calculations/1/close-approach", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(middleware.WithUserID(context.Background(), uint64(1)))
	req = mux.SetURLVars(req, map[string]string{"id": "1", "calculation_id": "1"})
	w := httptest.NewRecorder()

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil repository")
		}
	}()
	
	handler.AddCloseApproachData(w, req)
}

func TestCalculationHandler_ListCalculationsForComet_NoUserID(t *testing.T) {
	handler := NewCalculationHandler(&usecase.CalculationUsecase{})

	req := httptest.NewRequest("GET", "/comets/1/calculations", nil)
	req = req.WithContext(context.Background()) // No user ID in context
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.ListCalculationsForComet(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestCalculationHandler_ListCalculationsForComet_ValidUserID(t *testing.T) {
	handler := NewCalculationHandler(&usecase.CalculationUsecase{})

	req := httptest.NewRequest("GET", "/comets/1/calculations", nil)
	req = req.WithContext(middleware.WithUserID(context.Background(), uint64(1)))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.ListCalculationsForComet(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCalculationHandler_GetCalculation_NoUserID(t *testing.T) {
	handler := NewCalculationHandler(&usecase.CalculationUsecase{})

	req := httptest.NewRequest("GET", "/comets/1/calculations/1", nil)
	req = req.WithContext(context.Background()) // No user ID in context
	req = mux.SetURLVars(req, map[string]string{"id": "1", "calculation_id": "1"})
	w := httptest.NewRecorder()

	handler.GetCalculation(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestCalculationHandler_GetCalculation_ValidUserID(t *testing.T) {
	handler := NewCalculationHandler(&usecase.CalculationUsecase{})

	req := httptest.NewRequest("GET", "/comets/1/calculations/1", nil)
	req = req.WithContext(middleware.WithUserID(context.Background(), uint64(1)))
	req = mux.SetURLVars(req, map[string]string{"id": "1", "calculation_id": "1"})
	w := httptest.NewRecorder()

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil repository")
		}
	}()
	
	handler.GetCalculation(w, req)
}

func TestCalculationHandler_DeleteCalculation_NoUserID(t *testing.T) {
	handler := NewCalculationHandler(&usecase.CalculationUsecase{})

	req := httptest.NewRequest("DELETE", "/comets/1/calculations/1", nil)
	req = req.WithContext(context.Background()) // No user ID in context
	req = mux.SetURLVars(req, map[string]string{"id": "1", "calculation_id": "1"})
	w := httptest.NewRecorder()

	handler.DeleteCalculation(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestCalculationHandler_DeleteCalculation_ValidUserID(t *testing.T) {
	handler := NewCalculationHandler(&usecase.CalculationUsecase{})

	req := httptest.NewRequest("DELETE", "/comets/1/calculations/1", nil)
	req = req.WithContext(middleware.WithUserID(context.Background(), uint64(1)))
	req = mux.SetURLVars(req, map[string]string{"id": "1", "calculation_id": "1"})
	w := httptest.NewRecorder()

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil repository")
		}
	}()
	
	handler.DeleteCalculation(w, req)
}
