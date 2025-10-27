package apiutils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteJSON(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]string{"message": "test"}

	WriteJSON(w, http.StatusOK, data)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", w.Header().Get("Content-Type"))
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["message"] != "test" {
		t.Errorf("Expected message 'test', got %s", response["message"])
	}
}

func TestWriteJSON_WithError(t *testing.T) {
	w := httptest.NewRecorder()
	// Create a channel to test error handling
	data := make(chan int)

	WriteJSON(w, http.StatusInternalServerError, data)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestWriteError(t *testing.T) {
	w := httptest.NewRecorder()
	message := "test error"

	WriteError(w, http.StatusBadRequest, message)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", w.Header().Get("Content-Type"))
	}

	var response ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Error != message {
		t.Errorf("Expected error message '%s', got '%s'", message, response.Error)
	}
}

func TestWriteError_EmptyMessage(t *testing.T) {
	w := httptest.NewRecorder()
	message := ""

	WriteError(w, http.StatusInternalServerError, message)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}

	var response ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Error != "" {
		t.Errorf("Expected empty error message, got '%s'", response.Error)
	}
}

func TestErrorResponse_Structure(t *testing.T) {
	response := ErrorResponse{Error: "test error"}
	
	if response.Error != "test error" {
		t.Errorf("Expected error 'test error', got '%s'", response.Error)
	}
}
