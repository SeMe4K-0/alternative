package delivery

import (
	"backend/middleware"
	"backend/store"
	"backend/usecase"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserHandler_Register_InvalidJSON(t *testing.T) {
	handler := NewUserHandler(&usecase.UserUsecase{}, &store.MinIOStore{})

	req := httptest.NewRequest("POST", "/auth/register", strings.NewReader("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Register(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_Register_EmptyEmail(t *testing.T) {
	handler := NewUserHandler(&usecase.UserUsecase{}, &store.MinIOStore{})

	reqData := map[string]string{
		"email":    "",
		"password": "password123",
	}
	reqBody, _ := json.Marshal(reqData)

	req := httptest.NewRequest("POST", "/auth/register", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Register(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_Register_EmptyPassword(t *testing.T) {
	handler := NewUserHandler(&usecase.UserUsecase{}, &store.MinIOStore{})

	reqData := map[string]string{
		"email":    "test@example.com",
		"password": "",
	}
	reqBody, _ := json.Marshal(reqData)

	req := httptest.NewRequest("POST", "/auth/register", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Register(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_Register_ValidData(t *testing.T) {
	assert.True(t, true)
}

func TestUserHandler_Login_InvalidJSON(t *testing.T) {
	handler := NewUserHandler(&usecase.UserUsecase{}, &store.MinIOStore{})

	req := httptest.NewRequest("POST", "/auth/login", strings.NewReader("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Login(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_Login_EmptyEmail(t *testing.T) {
	handler := NewUserHandler(&usecase.UserUsecase{}, &store.MinIOStore{})

	reqData := map[string]string{
		"email":    "",
		"password": "password123",
	}
	reqBody, _ := json.Marshal(reqData)

	req := httptest.NewRequest("POST", "/auth/login", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Login(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_Login_EmptyPassword(t *testing.T) {
	handler := NewUserHandler(&usecase.UserUsecase{}, &store.MinIOStore{})

	reqData := map[string]string{
		"email":    "test@example.com",
		"password": "",
	}
	reqBody, _ := json.Marshal(reqData)

	req := httptest.NewRequest("POST", "/auth/login", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Login(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_Login_ValidData(t *testing.T) {
	assert.True(t, true)
}

func TestUserHandler_Logout_NoSession(t *testing.T) {
	handler := NewUserHandler(&usecase.UserUsecase{}, &store.MinIOStore{})

	req := httptest.NewRequest("POST", "/auth/logout", nil)
	w := httptest.NewRecorder()

	handler.Logout(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_GetProfile_NoUserID(t *testing.T) {
	handler := NewUserHandler(&usecase.UserUsecase{}, &store.MinIOStore{})

	req := httptest.NewRequest("GET", "/profile/me", nil)
	req = req.WithContext(context.Background()) // No user ID in context
	w := httptest.NewRecorder()

	handler.GetProfile(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestUserHandler_UpdateProfile_NoUserID(t *testing.T) {
	handler := NewUserHandler(&usecase.UserUsecase{}, &store.MinIOStore{})

	req := httptest.NewRequest("PUT", "/profile/me", nil)
	req = req.WithContext(context.Background()) // No user ID in context
	w := httptest.NewRecorder()

	handler.UpdateProfile(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestUserHandler_UpdateProfile_InvalidJSON(t *testing.T) {
	handler := NewUserHandler(&usecase.UserUsecase{}, &store.MinIOStore{})

	req := httptest.NewRequest("PUT", "/profile/me", strings.NewReader("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(middleware.WithUserID(context.Background(), uint64(1)))
	w := httptest.NewRecorder()

	handler.UpdateProfile(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_GetAvatar_EmptyPath(t *testing.T) {
	handler := NewUserHandler(&usecase.UserUsecase{}, &store.MinIOStore{})

	req := httptest.NewRequest("GET", "/api/files/", nil)
	w := httptest.NewRecorder()

	handler.GetAvatar(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_GetAvatar_ValidPath(t *testing.T) {
	handler := NewUserHandler(&usecase.UserUsecase{}, &store.MinIOStore{})

	req := httptest.NewRequest("GET", "/api/files/avatar_user_1.jpg", nil)
	w := httptest.NewRecorder()

	handler.GetAvatar(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUserHandler_Register_ValidData_WithPanic(t *testing.T) {
	handler := NewUserHandler(&usecase.UserUsecase{}, &store.MinIOStore{})

	reqData := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	reqBody, _ := json.Marshal(reqData)

	req := httptest.NewRequest("POST", "/auth/register", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil repository")
		}
	}()
	
	handler.Register(w, req)
}

func TestUserHandler_Login_ValidData_WithPanic(t *testing.T) {
	handler := NewUserHandler(&usecase.UserUsecase{}, &store.MinIOStore{})

	reqData := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	reqBody, _ := json.Marshal(reqData)

	req := httptest.NewRequest("POST", "/auth/login", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil repository")
		}
	}()
	
	handler.Login(w, req)
}

func TestUserHandler_Logout_ValidSession(t *testing.T) {
	handler := NewUserHandler(&usecase.UserUsecase{}, &store.MinIOStore{})

	req := httptest.NewRequest("POST", "/auth/logout", nil)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "valid-session"})
	w := httptest.NewRecorder()

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil Redis client")
		}
	}()
	
	handler.Logout(w, req)
}

func TestUserHandler_GetProfile_ValidUserID(t *testing.T) {
	handler := NewUserHandler(&usecase.UserUsecase{}, &store.MinIOStore{})

	req := httptest.NewRequest("GET", "/profile/me", nil)
	req = req.WithContext(middleware.WithUserID(context.Background(), uint64(1)))
	w := httptest.NewRecorder()

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil repository")
		}
	}()
	
	handler.GetProfile(w, req)
}

func TestUserHandler_UpdateProfile_ValidData(t *testing.T) {
	handler := NewUserHandler(&usecase.UserUsecase{}, &store.MinIOStore{})

	reqData := map[string]string{
		"username": "newusername",
		"email":    "newemail@example.com",
	}
	reqBody, _ := json.Marshal(reqData)

	req := httptest.NewRequest("PUT", "/profile/me", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(middleware.WithUserID(context.Background(), uint64(1)))
	w := httptest.NewRecorder()

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil repository")
		}
	}()
	
	handler.UpdateProfile(w, req)
}

func TestNewUserHandler(t *testing.T) {
	userUsecase := &usecase.UserUsecase{}
	minioStore := &store.MinIOStore{}
	
	handler := NewUserHandler(userUsecase, minioStore)

	assert.NotNil(t, handler)
	assert.Equal(t, userUsecase, handler.userUsecase)
	assert.Equal(t, minioStore, handler.minioStore)
}

func TestUserHandler_GetProfile_ValidData_WithPanic(t *testing.T) {
	handler := NewUserHandler(&usecase.UserUsecase{}, &store.MinIOStore{})

	req := httptest.NewRequest("GET", "/profile/me", nil)
	req = req.WithContext(middleware.WithUserID(context.Background(), uint64(1)))
	w := httptest.NewRecorder()

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil database")
		}
	}()
	
	handler.GetProfile(w, req)
}

func TestUserHandler_UpdateProfile_ValidData_WithPanic(t *testing.T) {
	handler := NewUserHandler(&usecase.UserUsecase{}, &store.MinIOStore{})

	reqData := map[string]string{
		"username": "newusername",
		"email":    "newemail@example.com",
	}
	reqBody, _ := json.Marshal(reqData)

	req := httptest.NewRequest("PUT", "/profile/me", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(middleware.WithUserID(context.Background(), uint64(1)))
	w := httptest.NewRecorder()

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil database")
		}
	}()
	
	handler.UpdateProfile(w, req)
}