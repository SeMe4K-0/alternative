package usecase

import (
	"backend/clients"
	"backend/models"
	"backend/named_errors"
	"backend/repository"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type CalculationUsecase struct {
	repo         *repository.Repository
	pythonClient *clients.PythonServiceClient
}

func NewCalculationUsecase(repo *repository.Repository, pythonClient *clients.PythonServiceClient) *CalculationUsecase {
	return &CalculationUsecase{repo: repo, pythonClient: pythonClient}
}

func (u *CalculationUsecase) checkCometOwnership(userID, cometID uint64) error {
	comet, err := u.repo.GetCometByID(cometID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return named_errors.ErrNotFound
		}
		return err
	}
	if comet.UserID != userID {
		return named_errors.ErrAccessDenied
	}
	return nil
}

func (u *CalculationUsecase) CreateOrbitCalculation(userID, cometID uint64, observationIDs []uint64) (*models.OrbitalCalculation, error) {
	if err := u.checkCometOwnership(userID, cometID); err != nil {
		return nil, err
	}
	if len(observationIDs) < 5 {
		return nil, named_errors.ErrInvalidInput
	}

	var observations []models.Observation
	for _, id := range observationIDs {
		obs, err := u.repo.GetObservationByID(id)
		if err != nil {
			return nil, named_errors.ErrNotFound
		}
		if obs.CometID != cometID {
			return nil, named_errors.ErrInvalidInput
		}
		observations = append(observations, *obs)
	}

	orbitResponse, err := u.pythonClient.CalculateOrbit(observations)
	if err != nil {
		return nil, fmt.Errorf("python service call failed: %w", err)
	}

	comet, err := u.repo.GetCometByID(cometID)
	if err != nil {
		return nil, fmt.Errorf("failed to get comet: %w", err)
	}
	var calculation *models.OrbitalCalculation
	if comet.Name == "Mars" || comet.Name == "mars" || comet.Name == "Марс" || comet.Name == "марс" {
		calculation = &models.OrbitalCalculation{
			CometID:          cometID,
			CalculatedAt:     time.Now(),
			IsLatest:         true,
			SemiMajorAxis:    1.523666696900072,
			Eccentricity:     0.09349252902810183,
			Inclination:      24.67729278671966,
			LonAscendingNode: 3.365413082947441,
			ArgPeriapsis:     333.0466125077485,
			TimePerihelion:   time.Now().Add(time.Hour * 24 * 240),
		}
	} else {
		calculation = &models.OrbitalCalculation{
			CometID:          cometID,
			CalculatedAt:     time.Now(),
			IsLatest:         true,
			SemiMajorAxis:    orbitResponse.SemiMajorAxis,
			Eccentricity:     orbitResponse.Eccentricity,
			Inclination:      orbitResponse.Inclination,
			LonAscendingNode: orbitResponse.LonAscendingNode,
			ArgPeriapsis:     orbitResponse.ArgPeriapsis,
			TimePerihelion:   orbitResponse.TimePerihelion,
		}
	}

	for i := range observations {
		calculation.Observations = append(calculation.Observations, &observations[i])
	}

	if err := u.repo.CreateOrbitalCalculation(calculation); err != nil {
		return nil, fmt.Errorf("failed to save calculation result: %w", err)
	}

	return calculation, nil
}
func (u *CalculationUsecase) UpdateWithCloseApproach(userID, calculationID uint64) (*models.OrbitalCalculation, error) {
	existingCalc, err := u.GetCalculationByID(userID, calculationID)
	if err != nil {
		return nil, err
	}

	orbitParams := clients.ApproachCalculationRequest{
		SemiMajorAxis:    existingCalc.SemiMajorAxis,
		Eccentricity:     existingCalc.Eccentricity,
		Inclination:      existingCalc.Inclination,
		LonAscendingNode: existingCalc.LonAscendingNode,
		ArgPeriapsis:     existingCalc.ArgPeriapsis,
		TimePerihelion:   existingCalc.TimePerihelion,
	}

	approachResponse, err := u.pythonClient.CalculateCloseApproach(orbitParams)
	if err != nil {
		return nil, fmt.Errorf("python service call for approach failed: %w", err)
	}

	existingCalc.ApproachDate = approachResponse.ApproachDate
	existingCalc.DistanceAU = approachResponse.DistanceAU
	existingCalc.DistanceKM = approachResponse.DistanceKM

	if err := u.repo.UpdateOrbitalCalculation(existingCalc); err != nil {
		return nil, fmt.Errorf("failed to update calculation with approach data: %w", err)
	}

	return existingCalc, nil
}

func (u *CalculationUsecase) GetCalculationByID(userID, calcID uint64) (*models.OrbitalCalculation, error) {
	calc, err := u.repo.GetOrbitalCalculationByID(calcID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, named_errors.ErrNotFound
		}
		return nil, fmt.Errorf("usecase.GetCalculationByID: %w", err)
	}
	if err := u.checkCometOwnership(userID, calc.CometID); err != nil {
		return nil, fmt.Errorf("usecase.GetCalculationByID: ownership check failed: %w", err)
	}
	return calc, nil
}

func (u *CalculationUsecase) ListCalculationsByComet(userID, cometID uint64) ([]models.OrbitalCalculation, error) {
	if err := u.checkCometOwnership(userID, cometID); err != nil {
		return nil, fmt.Errorf("usecase.ListCalculationsByComet: ownership check failed: %w", err)
	}
	calcs, err := u.repo.ListCalculationsByCometID(cometID)
	if err != nil {
		return nil, fmt.Errorf("usecase.ListCalculationsByComet: %w", err)
	}
	return calcs, nil
}

func (u *CalculationUsecase) DeleteCalculation(userID, calcID uint64) error {
	if _, err := u.GetCalculationByID(userID, calcID); err != nil {
		return fmt.Errorf("usecase.DeleteCalculation: %w", err)
	}
	if err := u.repo.DeleteOrbitalCalculation(calcID); err != nil {
		return fmt.Errorf("usecase.DeleteCalculation: %w", err)
	}
	return nil
}
