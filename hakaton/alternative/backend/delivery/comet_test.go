package delivery

import (
	"backend/middleware"
	"backend/usecase"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestCometHandler_CreateComet_NoUserID(t *testing.T) {
	handler := NewCometHandler(&usecase.CometUsecase{})

	reqData := map[string]string{
		"name":        "Test Comet",
		"description": "Test Description",
	}
	reqBody, _ := json.Marshal(reqData)

	req := httptest.NewRequest("POST", "/comets", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.Background()) // No user ID in context
	w := httptest.NewRecorder()

	handler.CreateComet(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestCometHandler_CreateComet_InvalidJSON(t *testing.T) {
	handler := NewCometHandler(&usecase.CometUsecase{})

	req := httptest.NewRequest("POST", "/comets", strings.NewReader("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(middleware.WithUserID(context.Background(), uint64(1)))
	w := httptest.NewRecorder()

	handler.CreateComet(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCometHandler_CreateComet_EmptyName(t *testing.T) {
	handler := NewCometHandler(&usecase.CometUsecase{})

	reqData := map[string]string{
		"name":        "",
		"description": "Test Description",
	}
	reqBody, _ := json.Marshal(reqData)

	req := httptest.NewRequest("POST", "/comets", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(middleware.WithUserID(context.Background(), uint64(1)))
	w := httptest.NewRecorder()

	handler.CreateComet(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCometHandler_CreateComet_ValidData(t *testing.T) {
	handler := NewCometHandler(&usecase.CometUsecase{})

	reqData := map[string]string{
		"name":        "Test Comet",
		"description": "Test Description",
	}
	reqBody, _ := json.Marshal(reqData)

	req := httptest.NewRequest("POST", "/comets", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(middleware.WithUserID(context.Background(), uint64(1)))
	w := httptest.NewRecorder()

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil repository")
		}
	}()
	
	handler.CreateComet(w, req)
}

func TestCometHandler_GetComet_NoUserID(t *testing.T) {
	handler := NewCometHandler(&usecase.CometUsecase{})

	req := httptest.NewRequest("GET", "/comets/1", nil)
	req = req.WithContext(context.Background()) // No user ID in context
	w := httptest.NewRecorder()

	handler.GetComet(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestCometHandler_GetComet_InvalidID(t *testing.T) {
	handler := NewCometHandler(&usecase.CometUsecase{})

	req := httptest.NewRequest("GET", "/comets/invalid", nil)
	req = req.WithContext(middleware.WithUserID(context.Background(), uint64(1)))
	w := httptest.NewRecorder()

	handler.GetComet(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCometHandler_GetComet_ValidID(t *testing.T) {
	handler := NewCometHandler(&usecase.CometUsecase{})

	req := httptest.NewRequest("GET", "/comets/1", nil)
	req = req.WithContext(middleware.WithUserID(context.Background(), uint64(1)))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil repository")
		}
	}()
	
	handler.GetComet(w, req)
}

func TestCometHandler_ListUserComets_NoUserID(t *testing.T) {
	handler := NewCometHandler(&usecase.CometUsecase{})

	req := httptest.NewRequest("GET", "/comets", nil)
	req = req.WithContext(context.Background()) // No user ID in context
	w := httptest.NewRecorder()

	handler.ListUserComets(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestCometHandler_ListUserComets_ValidUserID(t *testing.T) {
	handler := NewCometHandler(&usecase.CometUsecase{})

	req := httptest.NewRequest("GET", "/comets", nil)
	req = req.WithContext(middleware.WithUserID(context.Background(), uint64(1)))
	w := httptest.NewRecorder()

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil repository")
		}
	}()
	
	handler.ListUserComets(w, req)
}

func TestCometHandler_UpdateComet_NoUserID(t *testing.T) {
	handler := NewCometHandler(&usecase.CometUsecase{})

	reqData := map[string]string{
		"name":        "Updated Comet",
		"description": "Updated Description",
	}
	reqBody, _ := json.Marshal(reqData)

	req := httptest.NewRequest("PUT", "/comets/1", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.Background()) // No user ID in context
	w := httptest.NewRecorder()

	handler.UpdateComet(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestCometHandler_UpdateComet_InvalidID(t *testing.T) {
	handler := NewCometHandler(&usecase.CometUsecase{})

	reqData := map[string]string{
		"name":        "Updated Comet",
		"description": "Updated Description",
	}
	reqBody, _ := json.Marshal(reqData)

	req := httptest.NewRequest("PUT", "/comets/invalid", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(middleware.WithUserID(context.Background(), uint64(1)))
	w := httptest.NewRecorder()

	handler.UpdateComet(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCometHandler_UpdateComet_InvalidJSON(t *testing.T) {
	handler := NewCometHandler(&usecase.CometUsecase{})

	req := httptest.NewRequest("PUT", "/comets/1", strings.NewReader("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(middleware.WithUserID(context.Background(), uint64(1)))
	w := httptest.NewRecorder()

	handler.UpdateComet(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCometHandler_UpdateComet_ValidData(t *testing.T) {
	handler := NewCometHandler(&usecase.CometUsecase{})

	reqData := map[string]string{
		"name":        "Updated Comet",
		"description": "Updated Description",
	}
	reqBody, _ := json.Marshal(reqData)

	req := httptest.NewRequest("PUT", "/comets/1", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(middleware.WithUserID(context.Background(), uint64(1)))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil repository")
		}
	}()
	
	handler.UpdateComet(w, req)
}

func TestCometHandler_DeleteComet_NoUserID(t *testing.T) {
	handler := NewCometHandler(&usecase.CometUsecase{})

	req := httptest.NewRequest("DELETE", "/comets/1", nil)
	req = req.WithContext(context.Background()) // No user ID in context
	w := httptest.NewRecorder()

	handler.DeleteComet(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestCometHandler_DeleteComet_InvalidID(t *testing.T) {
	handler := NewCometHandler(&usecase.CometUsecase{})

	req := httptest.NewRequest("DELETE", "/comets/invalid", nil)
	req = req.WithContext(middleware.WithUserID(context.Background(), uint64(1)))
	w := httptest.NewRecorder()

	handler.DeleteComet(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCometHandler_DeleteComet_ValidID(t *testing.T) {
	handler := NewCometHandler(&usecase.CometUsecase{})

	req := httptest.NewRequest("DELETE", "/comets/1", nil)
	req = req.WithContext(middleware.WithUserID(context.Background(), uint64(1)))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil repository")
		}
	}()
	
	handler.DeleteComet(w, req)
}

func TestNewCometHandler(t *testing.T) {
	cometUsecase := &usecase.CometUsecase{}
	
	handler := NewCometHandler(cometUsecase)

	assert.NotNil(t, handler)
	assert.Equal(t, cometUsecase, handler.cometUsecase)
}

func TestCometHandler_CreateComet_EdgeCases(t *testing.T) {
	handler := NewCometHandler(&usecase.CometUsecase{})

	longName := string(make([]byte, 1000))
	for i := range longName {
		longName = longName[:i] + "a" + longName[i+1:]
	}
	
	reqData := map[string]string{
		"name":        longName,
		"description": "Test Description",
	}
	reqBody, _ := json.Marshal(reqData)

	req := httptest.NewRequest("POST", "/comets", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(middleware.WithUserID(context.Background(), uint64(1)))
	w := httptest.NewRecorder()

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil repository")
		}
	}()
	
	handler.CreateComet(w, req)
}

func TestCometHandler_GetComet_EdgeCases(t *testing.T) {
	handler := NewCometHandler(&usecase.CometUsecase{})

	req := httptest.NewRequest("GET", "/comets/0", nil)
	req = req.WithContext(middleware.WithUserID(context.Background(), uint64(1)))
	req = mux.SetURLVars(req, map[string]string{"id": "0"})
	w := httptest.NewRecorder()

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil repository")
		}
	}()
	
	handler.GetComet(w, req)
}

func TestCometHandler_UpdateComet_EdgeCases(t *testing.T) {
	handler := NewCometHandler(&usecase.CometUsecase{})

	reqData := map[string]string{
		"name":        "",
		"description": "",
	}
	reqBody, _ := json.Marshal(reqData)

	req := httptest.NewRequest("PUT", "/comets/1", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(middleware.WithUserID(context.Background(), uint64(1)))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil repository")
		}
	}()
	
	handler.UpdateComet(w, req)
}

func TestCometHandler_DeleteComet_EdgeCases(t *testing.T) {
	handler := NewCometHandler(&usecase.CometUsecase{})

	req := httptest.NewRequest("DELETE", "/comets/0", nil)
	req = req.WithContext(middleware.WithUserID(context.Background(), uint64(1)))
	req = mux.SetURLVars(req, map[string]string{"id": "0"})
	w := httptest.NewRecorder()

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil repository")
		}
	}()
	
	handler.DeleteComet(w, req)
}
