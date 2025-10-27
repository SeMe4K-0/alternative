package repository

import (
	"backend/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRepository_InvalidConfig(t *testing.T) {
	cfg := &config.DatabaseConfig{
		Host:     "invalid-host",
		Port:     5432,
		User:     "invalid-user",
		Password: "invalid-password",
		DBName:   "invalid-db",
		SSLMode:  "disable",
	}

	repo, err := NewRepository(cfg)

	assert.Error(t, err)
	assert.Nil(t, repo)
}

func TestNewRepository_EmptyConfig(t *testing.T) {
	cfg := &config.DatabaseConfig{}

	repo, err := NewRepository(cfg)

	assert.Error(t, err)
	assert.Nil(t, repo)
}

func TestNewRepository_NilConfig(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil config")
		}
	}()
	
	NewRepository(nil)
}

func TestRepository_CreateComet_NoConnection(t *testing.T) {
	repo := &Repository{}

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil database")
		}
	}()
	
	repo.CreateComet(nil)
}

func TestRepository_GetCometByID_NoConnection(t *testing.T) {
	repo := &Repository{}

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil database")
		}
	}()
	
	repo.GetCometByID(1)
}

func TestRepository_UpdateComet_NoConnection(t *testing.T) {
	repo := &Repository{}

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil database")
		}
	}()
	
	repo.UpdateComet(nil)
}

func TestRepository_DeleteComet_NoConnection(t *testing.T) {
	repo := &Repository{}

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil database")
		}
	}()
	
	repo.DeleteComet(1)
}

func TestRepository_ListCometsByUserID_NoConnection(t *testing.T) {
	repo := &Repository{}

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil database")
		}
	}()
	
	repo.ListCometsByUserID(1)
}

func TestRepository_CreateUser_NoConnection(t *testing.T) {
	repo := &Repository{}

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil database")
		}
	}()
	
	repo.CreateUser(nil)
}

func TestRepository_GetUserByEmail_NoConnection(t *testing.T) {
	repo := &Repository{}

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil database")
		}
	}()
	
	repo.GetUserByEmail("test@example.com")
}

func TestRepository_GetUserByUsername_NoConnection(t *testing.T) {
	repo := &Repository{}

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil database")
		}
	}()
	
	repo.GetUserByUsername("testuser")
}

func TestRepository_GetUserByID_NoConnection(t *testing.T) {
	repo := &Repository{}

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil database")
		}
	}()
	
	repo.GetUserByID(1)
}

func TestRepository_UpdateUser_NoConnection(t *testing.T) {
	repo := &Repository{}

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil database")
		}
	}()
	
	repo.UpdateUser(nil)
}

func TestRepository_CreateObservation_NoConnection(t *testing.T) {
	repo := &Repository{}

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil database")
		}
	}()
	
	repo.CreateObservation(nil)
}

func TestRepository_GetObservationByID_NoConnection(t *testing.T) {
	repo := &Repository{}

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil database")
		}
	}()
	
	repo.GetObservationByID(1)
}

func TestRepository_UpdateObservation_NoConnection(t *testing.T) {
	repo := &Repository{}

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil database")
		}
	}()
	
	repo.UpdateObservation(nil)
}

func TestRepository_DeleteObservation_NoConnection(t *testing.T) {
	repo := &Repository{}

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil database")
		}
	}()
	
	repo.DeleteObservation(1)
}

func TestRepository_ListObservationsByCometID_NoConnection(t *testing.T) {
	repo := &Repository{}

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil database")
		}
	}()
	
	repo.ListObservationsByCometID(1)
}

func TestRepository_CreateOrbitalCalculation_NoConnection(t *testing.T) {
	repo := &Repository{}

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil database")
		}
	}()
	
	repo.CreateOrbitalCalculation(nil)
}

func TestRepository_GetOrbitalCalculationByID_NoConnection(t *testing.T) {
	repo := &Repository{}

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil database")
		}
	}()
	
	repo.GetOrbitalCalculationByID(1)
}

func TestRepository_UpdateOrbitalCalculation_NoConnection(t *testing.T) {
	repo := &Repository{}

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil database")
		}
	}()
	
	repo.UpdateOrbitalCalculation(nil)
}

func TestRepository_DeleteOrbitalCalculation_NoConnection(t *testing.T) {
	repo := &Repository{}

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil database")
		}
	}()
	
	repo.DeleteOrbitalCalculation(1)
}

func TestRepository_ListCalculationsByCometID_NoConnection(t *testing.T) {
	repo := &Repository{}

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil database")
		}
	}()
	
	repo.ListCalculationsByCometID(1)
}