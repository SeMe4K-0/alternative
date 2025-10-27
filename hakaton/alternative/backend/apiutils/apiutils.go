package apiutils

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)
type ErrorResponse struct {
	Error string `json:"error" example:"сообщение об ошибке"`
}

func WriteJSON(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Error().Err(err).Msg("json encode error")
	}
}

func WriteError(w http.ResponseWriter, code int, message string) {
	response := ErrorResponse{Error: message}
	WriteJSON(w, code, response)
}
