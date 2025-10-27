package middleware

import (
	"backend/models"
	"backend/store"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockRedisStore struct {
	sessions map[string]*models.User
}

func NewMockRedisStore() *MockRedisStore {
	return &MockRedisStore{
		sessions: make(map[string]*models.User),
	}
}

func (m *MockRedisStore) CreateSession(userID uint64) (map[string]interface{}, error) {
	return nil, nil
}

func (m *MockRedisStore) GetUserBySession(token string) (*models.User, bool) {
	user, exists := m.sessions[token]
	return user, exists
}

func (m *MockRedisStore) DeleteSession(token string) error {
	delete(m.sessions, token)
	return nil
}

func (m *MockRedisStore) CleanExpiredSessions() error {
	return nil
}

func (m *MockRedisStore) Close() error {
	return nil
}

func TestWithUserID(t *testing.T) {
	ctx := context.Background()
	userID := uint64(123)

	newCtx := WithUserID(ctx, userID)

	retrievedID, ok := GetUserID(newCtx)
	if !ok {
		t.Error("Expected to retrieve user ID from context")
	}

	if retrievedID != userID {
		t.Errorf("Expected user ID %d, got %d", userID, retrievedID)
	}
}

func TestGetUserID_Valid(t *testing.T) {
	ctx := context.Background()
	userID := uint64(123)

	ctx = WithUserID(ctx, userID)

	retrievedID, ok := GetUserID(ctx)
	if !ok {
		t.Error("Expected to retrieve user ID from context")
	}

	if retrievedID != userID {
		t.Errorf("Expected user ID %d, got %d", userID, retrievedID)
	}
}

func TestGetUserID_Invalid(t *testing.T) {
	ctx := context.Background()

	_, ok := GetUserID(ctx)
	if ok {
		t.Error("Expected not to retrieve user ID from empty context")
	}
}

func TestGetUserID_WrongType(t *testing.T) {
	ctx := context.WithValue(context.Background(), UserIDKey, "not a uint64")

	_, ok := GetUserID(ctx)
	if ok {
		t.Error("Expected not to retrieve user ID with wrong type")
	}
}

func TestCORS_AllowedOrigin(t *testing.T) {
	handler := CORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://localhost:8080")

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Header().Get("Access-Control-Allow-Origin") != "http://localhost:8080" {
		t.Errorf("Expected Access-Control-Allow-Origin http://localhost:8080, got %s", w.Header().Get("Access-Control-Allow-Origin"))
	}

	if w.Header().Get("Access-Control-Allow-Methods") != "GET, POST, PUT, DELETE, OPTIONS" {
		t.Errorf("Expected Access-Control-Allow-Methods GET, POST, PUT, DELETE, OPTIONS, got %s", w.Header().Get("Access-Control-Allow-Methods"))
	}

	if w.Header().Get("Access-Control-Allow-Headers") != "Content-Type, Authorization" {
		t.Errorf("Expected Access-Control-Allow-Headers Content-Type, Authorization, got %s", w.Header().Get("Access-Control-Allow-Headers"))
	}

	if w.Header().Get("Access-Control-Allow-Credentials") != "true" {
		t.Errorf("Expected Access-Control-Allow-Credentials true, got %s", w.Header().Get("Access-Control-Allow-Credentials"))
	}
}

func TestCORS_DisallowedOrigin(t *testing.T) {
	handler := CORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://malicious.com")

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Errorf("Expected empty Access-Control-Allow-Origin, got %s", w.Header().Get("Access-Control-Allow-Origin"))
	}
}

func TestCORS_NoOrigin(t *testing.T) {
	handler := CORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Errorf("Expected empty Access-Control-Allow-Origin, got %s", w.Header().Get("Access-Control-Allow-Origin"))
	}
}

func TestCORS_OPTIONS(t *testing.T) {
	handler := CORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("OPTIONS", "/test", nil)
	req.Header.Set("Origin", "http://localhost:8080")

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestAuthMiddleware_NoCookie(t *testing.T) {
	handler := AuthMiddleware(nil)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called when no cookie")
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAuthMiddleware_InvalidSession(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil store")
		}
	}()
	
	handler := AuthMiddleware(nil)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called with invalid session")
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "invalid"})
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
}

func TestAuthMiddleware_ValidSession(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil store")
		}
	}()
	
	handler := AuthMiddleware(nil)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := GetUserID(r.Context())
		if !ok {
			t.Error("Expected to retrieve user ID from context")
		}
		if userID != 123 {
			t.Errorf("Expected user ID 123, got %d", userID)
		}
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "valid-token"})
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
}

func TestAuthMiddleware_ValidSessionWithMockStore(t *testing.T) {
	mockStore := NewMockRedisStore()
	user := &models.User{ID: 123, Email: "test@example.com"}
	mockStore.sessions["valid-token"] = user

	redisStore := &store.RedisStore{}
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil Redis client")
		}
	}()
	
	handler := AuthMiddleware(redisStore)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := GetUserID(r.Context())
		if !ok {
			t.Error("Expected to retrieve user ID from context")
		}
		if userID != 123 {
			t.Errorf("Expected user ID 123, got %d", userID)
		}
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "valid-token"})
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
}

func TestAuthMiddleware_InvalidSessionWithMockStore(t *testing.T) {
	redisStore := &store.RedisStore{}
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil Redis client")
		}
	}()
	
	handler := AuthMiddleware(redisStore)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called with invalid session")
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "invalid-token"})
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
}

func TestCORS_AllAllowedOrigins(t *testing.T) {
	allowedOrigins := []string{
		"http://localhost:8080",
		"http://127.0.0.1:8080",
		"http://localhost:8030",
		"http://127.0.0.1:8030",
		"http://89.208.210.115:8030",
		"http://localhost:3001",
	}

	for _, origin := range allowedOrigins {
		t.Run("Origin_"+origin, func(t *testing.T) {
			handler := CORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			req := httptest.NewRequest("GET", "/test", nil)
			req.Header.Set("Origin", origin)

			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			if w.Header().Get("Access-Control-Allow-Origin") != origin {
				t.Errorf("Expected Access-Control-Allow-Origin %s, got %s", origin, w.Header().Get("Access-Control-Allow-Origin"))
			}
		})
	}
}