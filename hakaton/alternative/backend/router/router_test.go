package router

import (
	"backend/clients"
	"backend/config"
	"backend/repository"
	"backend/store"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRouter(t *testing.T) {
	redisStore := &store.RedisStore{}
	minioStore := &store.MinIOStore{}
	repo := &repository.Repository{}
	pythonClient := &clients.PythonServiceClient{}
	cfg := &config.Config{}

	router := NewRouter(redisStore, minioStore, repo, pythonClient, cfg)

	assert.NotNil(t, router)
}

func TestRouter_HealthEndpoint(t *testing.T) {
	redisStore := &store.RedisStore{}
	minioStore := &store.MinIOStore{}
	repo := &repository.Repository{}
	pythonClient := &clients.PythonServiceClient{}
	cfg := &config.Config{}

	router := NewRouter(redisStore, minioStore, repo, pythonClient, cfg)

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRouter_SwaggerEndpoint(t *testing.T) {
	redisStore := &store.RedisStore{}
	minioStore := &store.MinIOStore{}
	repo := &repository.Repository{}
	pythonClient := &clients.PythonServiceClient{}
	cfg := &config.Config{}

	router := NewRouter(redisStore, minioStore, repo, pythonClient, cfg)

	req := httptest.NewRequest("GET", "/swagger/", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Swagger endpoint should be accessible (may redirect)
	assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusMovedPermanently)
}

func TestRouter_RegisterEndpoint(t *testing.T) {
	redisStore := &store.RedisStore{}
	minioStore := &store.MinIOStore{}
	repo := &repository.Repository{}
	pythonClient := &clients.PythonServiceClient{}
	cfg := &config.Config{}

	router := NewRouter(redisStore, minioStore, repo, pythonClient, cfg)

	req := httptest.NewRequest("POST", "/api/auth/register", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// This will fail because we don't have a real database connection
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRouter_LoginEndpoint(t *testing.T) {
	redisStore := &store.RedisStore{}
	minioStore := &store.MinIOStore{}
	repo := &repository.Repository{}
	pythonClient := &clients.PythonServiceClient{}
	cfg := &config.Config{}

	router := NewRouter(redisStore, minioStore, repo, pythonClient, cfg)

	req := httptest.NewRequest("POST", "/api/auth/login", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// This will fail because we don't have a real database connection
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRouter_LogoutEndpoint(t *testing.T) {
	redisStore := &store.RedisStore{}
	minioStore := &store.MinIOStore{}
	repo := &repository.Repository{}
	pythonClient := &clients.PythonServiceClient{}
	cfg := &config.Config{}

	router := NewRouter(redisStore, minioStore, repo, pythonClient, cfg)

	req := httptest.NewRequest("POST", "/api/auth/logout", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// This will fail because we don't have a real database connection
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRouter_CometsEndpoint(t *testing.T) {
	redisStore := &store.RedisStore{}
	minioStore := &store.MinIOStore{}
	repo := &repository.Repository{}
	pythonClient := &clients.PythonServiceClient{}
	cfg := &config.Config{}

	router := NewRouter(redisStore, minioStore, repo, pythonClient, cfg)

	req := httptest.NewRequest("GET", "/api/comets", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// This will fail because we don't have a real database connection
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
