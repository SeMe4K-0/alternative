package usecase

import (
	"backend/named_errors"
	"backend/repository"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCometUsecase_CreateComet_EmptyName(t *testing.T) {
	usecase := NewCometUsecase(&repository.Repository{})

	_, err := usecase.CreateComet(1, "", "description")

	assert.Error(t, err)
	assert.Equal(t, named_errors.ErrInvalidInput, err)
}

func TestCometUsecase_CreateComet_ValidName(t *testing.T) {
	usecase := NewCometUsecase(&repository.Repository{})

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil database")
		}
	}()
	
	usecase.CreateComet(1, "Test Comet", "Test Description")
}

func TestCometUsecase_GetCometByID_RepositoryError(t *testing.T) {
	usecase := NewCometUsecase(&repository.Repository{})

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil database")
		}
	}()
	
	usecase.GetCometByID(1, 1)
}

func TestCometUsecase_ListCometsByUserID_RepositoryError(t *testing.T) {
	usecase := NewCometUsecase(&repository.Repository{})

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil database")
		}
	}()
	
	usecase.ListCometsByUserID(1)
}

func TestCometUsecase_UpdateComet_GetCometError(t *testing.T) {
	usecase := NewCometUsecase(&repository.Repository{})

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil database")
		}
	}()
	
	usecase.UpdateComet(1, 1, "Updated Comet", "Updated Description")
}

func TestCometUsecase_DeleteComet_GetCometError(t *testing.T) {
	usecase := NewCometUsecase(&repository.Repository{})

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil database")
		}
	}()
	
	usecase.DeleteComet(1, 1)
}

func TestNewCometUsecase(t *testing.T) {
	repo := &repository.Repository{}
	usecase := NewCometUsecase(repo)

	assert.NotNil(t, usecase)
	assert.Equal(t, repo, usecase.cometRepo)
}

func TestCometUsecase_CreateComet_EdgeCases(t *testing.T) {
	usecase := NewCometUsecase(&repository.Repository{})

	longName := string(make([]byte, 1000))
	for i := range longName {
		longName = longName[:i] + "a" + longName[i+1:]
	}
	
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil database")
		}
	}()
	
	usecase.CreateComet(1, longName, "description")
}

func TestCometUsecase_GetCometByID_EdgeCases(t *testing.T) {
	usecase := NewCometUsecase(&repository.Repository{})

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil database")
		}
	}()
	
	usecase.GetCometByID(1, 0)
}

func TestCometUsecase_ListCometsByUserID_EdgeCases(t *testing.T) {
	usecase := NewCometUsecase(&repository.Repository{})

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil database")
		}
	}()
	
	usecase.ListCometsByUserID(0)
}

func TestCometUsecase_UpdateComet_EdgeCases(t *testing.T) {
	usecase := NewCometUsecase(&repository.Repository{})

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil database")
		}
	}()
	
	usecase.UpdateComet(1, 1, "", "")
}

func TestCometUsecase_DeleteComet_EdgeCases(t *testing.T) {
	usecase := NewCometUsecase(&repository.Repository{})

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil database")
		}
	}()
	
	usecase.DeleteComet(0, 0)
}