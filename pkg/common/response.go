// pkg/common/response.go
package common

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	response := Response{
		Status: "success",
		Data:   data,
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