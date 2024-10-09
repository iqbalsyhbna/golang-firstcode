package common

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	response := Response{
		Status:  "success",
		Data:    data,
		Message: message,
	}

	json.NewEncoder(w).Encode(response)
}

func WriteError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	response := Response{
		Status:  "error",
		Message: message,
	}
	json.NewEncoder(w).Encode(response)
}
