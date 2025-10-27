package clients

import (
	"backend/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewPythonServiceClient(t *testing.T) {
	client := NewPythonServiceClient("http://localhost:8000")

	assert.NotNil(t, client)
	assert.Equal(t, "http://localhost:8000", client.baseURL)
	assert.NotNil(t, client.client)
}

func TestPythonServiceClient_CalculateOrbit_InvalidURL(t *testing.T) {
	client := NewPythonServiceClient("invalid-url")

	observations := []models.Observation{
		{CometID: 1, ObservedAt: time.Now(), RA: 1.0, Dec: 2.0},
		{CometID: 1, ObservedAt: time.Now(), RA: 1.1, Dec: 2.1},
		{CometID: 1, ObservedAt: time.Now(), RA: 1.2, Dec: 2.2},
		{CometID: 1, ObservedAt: time.Now(), RA: 1.3, Dec: 2.3},
		{CometID: 1, ObservedAt: time.Now(), RA: 1.4, Dec: 2.4},
	}

	_, err := client.CalculateOrbit(observations)

	assert.Error(t, err)
}

func TestPythonServiceClient_CalculateCloseApproach_InvalidURL(t *testing.T) {
	client := NewPythonServiceClient("invalid-url")

	params := ApproachCalculationRequest{
		SemiMajorAxis:    1.5,
		Eccentricity:     0.1,
		Inclination:      0.2,
		LonAscendingNode: 0.3,
		ArgPeriapsis:     0.4,
		TimePerihelion:   time.Now(),
	}

	_, err := client.CalculateCloseApproach(params)

	assert.Error(t, err)
}

func TestPythonServiceClient_CalculateOrbit_EmptyObservations(t *testing.T) {
	client := NewPythonServiceClient("http://localhost:8000")

	observations := []models.Observation{}

	_, err := client.CalculateOrbit(observations)

	assert.Error(t, err)
}

func TestPythonServiceClient_CalculateOrbit_ValidObservations(t *testing.T) {
	client := NewPythonServiceClient("http://localhost:8000")

	observations := []models.Observation{
		{CometID: 1, ObservedAt: time.Now(), RA: 1.0, Dec: 2.0},
		{CometID: 1, ObservedAt: time.Now(), RA: 1.1, Dec: 2.1},
		{CometID: 1, ObservedAt: time.Now(), RA: 1.2, Dec: 2.2},
		{CometID: 1, ObservedAt: time.Now(), RA: 1.3, Dec: 2.3},
		{CometID: 1, ObservedAt: time.Now(), RA: 1.4, Dec: 2.4},
	}

	_, err := client.CalculateOrbit(observations)

	assert.Error(t, err)
}

func TestPythonServiceClient_CalculateCloseApproach_ValidParams(t *testing.T) {
	client := NewPythonServiceClient("http://localhost:8000")

	params := ApproachCalculationRequest{
		SemiMajorAxis:    1.5,
		Eccentricity:     0.1,
		Inclination:      0.2,
		LonAscendingNode: 0.3,
		ArgPeriapsis:     0.4,
		TimePerihelion:   time.Now(),
	}

	_, err := client.CalculateCloseApproach(params)

	assert.Error(t, err)
}

func TestGenerateMockOrbitResponse(t *testing.T) {
	response := generateMockOrbitResponse()
	
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.SemiMajorAxis)
	assert.NotEmpty(t, response.Eccentricity)
	assert.NotEmpty(t, response.Inclination)
	assert.NotEmpty(t, response.LonAscendingNode)
	assert.NotEmpty(t, response.ArgPeriapsis)
	assert.NotZero(t, response.TimePerihelion)
}

func TestGenerateMockApproachResponse(t *testing.T) {
	response := generateMockApproachResponse()
	
	assert.NotNil(t, response)
	assert.NotZero(t, response.ApproachDate)
	assert.NotZero(t, response.DistanceAU)
	assert.NotZero(t, response.DistanceKM)
}

func TestPythonServiceClient_CalculateOrbit_MockResponse(t *testing.T) {
	client := NewPythonServiceClient("http://localhost:8000")

	observations := []models.Observation{
		{CometID: 1, ObservedAt: time.Now(), RA: 1.0, Dec: 2.0},
		{CometID: 1, ObservedAt: time.Now(), RA: 1.1, Dec: 2.1},
		{CometID: 1, ObservedAt: time.Now(), RA: 1.2, Dec: 2.2},
		{CometID: 1, ObservedAt: time.Now(), RA: 1.3, Dec: 2.3},
		{CometID: 1, ObservedAt: time.Now(), RA: 1.4, Dec: 2.4},
	}

	mockResponse := generateMockOrbitResponse()
	assert.NotNil(t, mockResponse)
	
	_, err := client.CalculateOrbit(observations)
	assert.Error(t, err)
}

func TestPythonServiceClient_CalculateCloseApproach_MockResponse(t *testing.T) {
	client := NewPythonServiceClient("http://localhost:8000")

	params := ApproachCalculationRequest{
		SemiMajorAxis:    1.5,
		Eccentricity:     0.1,
		Inclination:      0.2,
		LonAscendingNode: 0.3,
		ArgPeriapsis:     0.4,
		TimePerihelion:   time.Now(),
	}

	mockResponse := generateMockApproachResponse()
	assert.NotNil(t, mockResponse)
	
	_, err := client.CalculateCloseApproach(params)
	assert.Error(t, err)
}
