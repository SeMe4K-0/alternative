// ./clients/python_client.go
package clients

import (
	"backend/models"
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog/log"
)

// --- Структуры для ПЕРВОЙ ручки (расчет орбиты) ---
type OrbitCalculationRequest struct {
	Observations []models.Observation `json:"observations"`
}

type OrbitCalculationResponse struct {
	SemiMajorAxis    float64   `json:"semi_major_axis"`
	Eccentricity     float64   `json:"eccentricity"`
	Inclination      float64   `json:"inclination"`
	LonAscendingNode float64   `json:"lon_ascending_node"`
	ArgPeriapsis     float64   `json:"arg_periapsis"`
	TimePerihelion   time.Time `json:"time_perihelion"`
}

// --- Структуры для ВТОРОЙ ручки (расчет сближения) ---
type ApproachCalculationRequest struct {
	SemiMajorAxis    float64   `json:"semi_major_axis"`
	Eccentricity     float64   `json:"eccentricity"`
	Inclination      float64   `json:"inclination"`
	LonAscendingNode float64   `json:"lon_ascending_node"`
	ArgPeriapsis     float64   `json:"arg_periapsis"`
	TimePerihelion   time.Time `json:"time_perihelion"`
}

type ApproachCalculationResponse struct {
	ApproachDate time.Time `json:"approach_date"`
	DistanceAU   float64   `json:"distance_au"`
	DistanceKM   float64   `json:"distance_km"`
}

// --- Клиент ---
type PythonServiceClient struct {
	client  *http.Client
	baseURL string
}

func NewPythonServiceClient(baseURL string) *PythonServiceClient {
	return &PythonServiceClient{
		client:  &http.Client{Timeout: 30 * time.Second},
		baseURL: baseURL,
	}
}

// РУЧКА 1: Отправляет наблюдения, получает 6 параметров орбиты
func (c *PythonServiceClient) CalculateOrbit(observations []models.Observation) (*OrbitCalculationResponse, error) {
	if os.Getenv("USE_PYTHON_MOCK") == "true" {
		log.Warn().Msg("--- MOCK: CalculateOrbit ---")
		return generateMockOrbitResponse(), nil
	}

	log.Info().Msg("--- CALLING: /calculate-orbit ---")
	requestBody, err := json.Marshal(OrbitCalculationRequest{Observations: observations})
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Post(c.baseURL+"/calculate-orbit", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("python service returned non-200 status: %d", resp.StatusCode)
	}

	var orbitResponse OrbitCalculationResponse
	if err := json.NewDecoder(resp.Body).Decode(&orbitResponse); err != nil {
		return nil, err
	}
	return &orbitResponse, nil
}

// РУЧКА 2: Отправляет 6 параметров орбиты, получает 3 параметра сближения
func (c *PythonServiceClient) CalculateCloseApproach(orbitParams ApproachCalculationRequest) (*ApproachCalculationResponse, error) {
	if os.Getenv("USE_PYTHON_MOCK") == "true" {
		log.Warn().Msg("--- MOCK: CalculateCloseApproach ---")
		return generateMockApproachResponse(), nil
	}

	log.Info().Msg("--- CALLING: /closest-approach ---")
	requestBody, err := json.Marshal(orbitParams)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Post(c.baseURL+"/closest-approach", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("python service returned non-200 status: %d", resp.StatusCode)
	}

	var approachResponse ApproachCalculationResponse
	if err := json.NewDecoder(resp.Body).Decode(&approachResponse); err != nil {
		return nil, err
	}
	return &approachResponse, nil
}

// --- Заглушки ---
func generateMockOrbitResponse() *OrbitCalculationResponse {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return &OrbitCalculationResponse{
		SemiMajorAxis:    3.1 + r.Float64(),
		Eccentricity:     0.75 + r.Float64()*0.1,
		Inclination:      22.0 + r.Float64()*5,
		LonAscendingNode: 150.0 + r.Float64()*10,
		ArgPeriapsis:     110.0 + r.Float64()*10,
		TimePerihelion:   time.Now().AddDate(0, 8, 0),
	}
}

func generateMockApproachResponse() *ApproachCalculationResponse {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return &ApproachCalculationResponse{
		ApproachDate: time.Now().AddDate(1, 2, 0),
		DistanceAU:   0.05 + r.Float64()*0.1,
		DistanceKM:   (0.05 + r.Float64()*0.1) * 149597870.7,
	}
}
