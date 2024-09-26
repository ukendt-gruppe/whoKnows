package utils

import (
	"encoding/json"
	"net/http"
)

// JSONRepsonse writes a response with the provided status code and data.
func JSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

type StandardResponse struct {
	Data interface{} `json:"data"`
}
