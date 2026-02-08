package api

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Data any `json:"data"`
	Meta any `json:"meta,omitempty"`
}

func SendJSON(w http.ResponseWriter, status int, data any, meta any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	
	resp := Response{
		Data: data,
		Meta: meta,
	}
	
	json.NewEncoder(w).Encode(resp)
}
