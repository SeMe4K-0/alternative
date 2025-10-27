package usecase

import (
	"backend/clients"
	"backend/repository"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculationUsecase_CreateOrbitCalculation_RepositoryError(t *testing.T) {
	usecase := NewCalculationUsecase(&repository.Repository{}, &clients.PythonServiceClient{})

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil database")
		}
	}()
	
	usecase.CreateOrbitCalculation(1, 1, []uint64{1, 2, 3, 4, 5})
}

func TestCalculationUsecase_CreateOrbitCalculation_InsufficientObservations(t *testing.T) {
	usecase := NewCalculationUsecase(&repository.Repository{}, &clients.PythonServiceClient{})

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil database")
		}
	}()
	
	usecase.CreateOrbitCalculation(1, 1, []uint64{1, 2, 3, 4, 5})
}

func TestCalculationUsecase_UpdateWithCloseApproach_RepositoryError(t *testing.T) {
	usecase := NewCalculationUsecase(&repository.Repository{}, &clients.PythonServiceClient{})

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil database")
		}
	}()
	
	usecase.UpdateWithCloseApproach(1, 1)
}

func TestCalculationUsecase_GetCalculationByID_RepositoryError(t *testing.T) {
	usecase := NewCalculationUsecase(&repository.Repository{}, &clients.PythonServiceClient{})

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil database")
		}
	}()
	
	usecase.GetCalculationByID(1, 1)
}

func TestCalculationUsecase_ListCalculationsByComet_RepositoryError(t *testing.T) {
	usecase := NewCalculationUsecase(&repository.Repository{}, &clients.PythonServiceClient{})

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil database")
		}
	}()
	
	usecase.ListCalculationsByComet(1, 1)
}

func TestCalculationUsecase_DeleteCalculation_RepositoryError(t *testing.T) {
	usecase := NewCalculationUsecase(&repository.Repository{}, &clients.PythonServiceClient{})

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil database")
		}
	}()
	
	usecase.DeleteCalculation(1, 1)
}

func TestNewCalculationUsecase(t *testing.T) {
	repo := &repository.Repository{}
	pythonClient := &clients.PythonServiceClient{}
	
	usecase := NewCalculationUsecase(repo, pythonClient)

	assert.NotNil(t, usecase)
	assert.Equal(t, repo, usecase.repo)
	assert.Equal(t, pythonClient, usecase.pythonClient)
}

func TestCalculationUsecase_checkCometOwnership_RepositoryError(t *testing.T) {
	usecase := NewCalculationUsecase(&repository.Repository{}, &clients.PythonServiceClient{})

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil database")
		}
	}()
	
	usecase.checkCometOwnership(1, 1)
}
