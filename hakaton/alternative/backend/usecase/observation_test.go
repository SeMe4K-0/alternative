package usecase

import (
	"backend/models"
	"backend/repository"
	"backend/store"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestObservationUsecase_CreateObservation_RepositoryError(t *testing.T) {
	usecase := NewObservationUsecase(&repository.Repository{}, &store.MinIOStore{})

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil database")
		}
	}()
	
	usecase.CreateObservation(1, 1, time.Now(), 1.0, 2.0, "test notes")
}

func TestObservationUsecase_GetObservationByID_RepositoryError(t *testing.T) {
	usecase := NewObservationUsecase(&repository.Repository{}, &store.MinIOStore{})

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil database")
		}
	}()
	
	usecase.GetObservationByID(1, 1)
}

func TestObservationUsecase_ListObservationsByComet_RepositoryError(t *testing.T) {
	usecase := NewObservationUsecase(&repository.Repository{}, &store.MinIOStore{})

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil database")
		}
	}()
	
	usecase.ListObservationsByComet(1, 1)
}

func TestObservationUsecase_UpdateObservation_RepositoryError(t *testing.T) {
	usecase := NewObservationUsecase(&repository.Repository{}, &store.MinIOStore{})

	observation := models.Observation{
		CometID:    1,
		ObservedAt: time.Now(),
		RA:         1.0,
		Dec:        2.0,
		Notes:      "updated notes",
	}

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil database")
		}
	}()
	
	usecase.UpdateObservation(1, 1, observation)
}

func TestObservationUsecase_DeleteObservation_RepositoryError(t *testing.T) {
	usecase := NewObservationUsecase(&repository.Repository{}, &store.MinIOStore{})

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil database")
		}
	}()
	
	usecase.DeleteObservation(1, 1)
}

func TestObservationUsecase_UploadObservationImage_RepositoryError(t *testing.T) {
	usecase := NewObservationUsecase(&repository.Repository{}, &store.MinIOStore{})

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil database")
		}
	}()
	
	usecase.UploadObservationImage(1, 1, nil, "test.jpg", "image/jpeg")
}

func TestObservationUsecase_DeleteObservationImage_RepositoryError(t *testing.T) {
	usecase := NewObservationUsecase(&repository.Repository{}, &store.MinIOStore{})

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil database")
		}
	}()
	
	usecase.DeleteObservationImage(1, 1)
}

// Test the NewObservationUsecase constructor
func TestNewObservationUsecase(t *testing.T) {
	repo := &repository.Repository{}
	minioStore := &store.MinIOStore{}
	
	usecase := NewObservationUsecase(repo, minioStore)

	assert.NotNil(t, usecase)
	assert.Equal(t, repo, usecase.repo)
	assert.Equal(t, minioStore, usecase.minioStore)
}

// Test the checkCometOwnership method
func TestObservationUsecase_checkCometOwnership_RepositoryError(t *testing.T) {
	usecase := NewObservationUsecase(&repository.Repository{}, &store.MinIOStore{})

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil database")
		}
	}()
	
	usecase.checkCometOwnership(1, 1)
}
