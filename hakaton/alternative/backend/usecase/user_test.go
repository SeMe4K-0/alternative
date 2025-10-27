package usecase

import (
	"backend/models"
	"backend/repository"
	"backend/services"
	"backend/store"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserUsecase_RegisterUser_RepositoryError(t *testing.T) {
	usecase := NewUserUsecase(&repository.Repository{}, &store.RedisStore{}, &store.MinIOStore{}, &services.EmailService{})

	// This will panic because we don't have a real database connection
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil database")
		}
	}()
	
	usecase.RegisterUser("test@example.com", nil, "password123")
}

func TestUserUsecase_LoginUser_RepositoryError(t *testing.T) {
	usecase := NewUserUsecase(&repository.Repository{}, &store.RedisStore{}, &store.MinIOStore{}, &services.EmailService{})

	// This will panic because we don't have a real database connection
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil database")
		}
	}()
	
	usecase.LoginUser("test@example.com", "password123")
}

func TestUserUsecase_GetUserProfile_RepositoryError(t *testing.T) {
	usecase := NewUserUsecase(&repository.Repository{}, &store.RedisStore{}, &store.MinIOStore{}, &services.EmailService{})

	// This will panic because we don't have a real database connection
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil database")
		}
	}()
	
	usecase.GetUserProfile(1)
}

func TestUserUsecase_UpdateUser_RepositoryError(t *testing.T) {
	usecase := NewUserUsecase(&repository.Repository{}, &store.RedisStore{}, &store.MinIOStore{}, &services.EmailService{})

	user := &models.User{
		ID:       1,
		Email:    "test@example.com",
		Username: "testuser",
	}

	// This will panic because we don't have a real database connection
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil database")
		}
	}()
	
	usecase.UpdateUser(user)
}

func TestUserUsecase_UpdateUserAvatar_RepositoryError(t *testing.T) {
	usecase := NewUserUsecase(&repository.Repository{}, &store.RedisStore{}, &store.MinIOStore{}, &services.EmailService{})

	// This will panic because we don't have a real database connection
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil database")
		}
	}()
	
	usecase.UpdateUserAvatar(1, "https://example.com/avatar.jpg")
}

func TestUserUsecase_LogoutUser_StoreError(t *testing.T) {
	usecase := NewUserUsecase(&repository.Repository{}, &store.RedisStore{}, &store.MinIOStore{}, &services.EmailService{})

	// This will panic because we don't have a real Redis connection
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil Redis client")
		}
	}()
	
	usecase.LogoutUser("token123")
}

func TestUserUsecase_GetUserBySession_StoreError(t *testing.T) {
	usecase := NewUserUsecase(&repository.Repository{}, &store.RedisStore{}, &store.MinIOStore{}, &services.EmailService{})

	// This will panic because we don't have a real Redis connection
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil Redis client")
		}
	}()
	
	usecase.GetUserBySession("token123")
}

// Test the NewUserUsecase constructor
func TestNewUserUsecase(t *testing.T) {
	repo := &repository.Repository{}
	redisStore := &store.RedisStore{}
	minioStore := &store.MinIOStore{}
	emailService := &services.EmailService{}
	
	usecase := NewUserUsecase(repo, redisStore, minioStore, emailService)

	assert.NotNil(t, usecase)
	assert.Equal(t, repo, usecase.repo)
	assert.Equal(t, redisStore, usecase.redisStore)
	assert.Equal(t, minioStore, usecase.minioStore)
	assert.Equal(t, emailService, usecase.emailService)
}