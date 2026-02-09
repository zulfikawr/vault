package core

import (
	"encoding/json"
	"net/http"
)

type VaultError struct {
	Status  int            `json:"-"`
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Details map[string]any `json:"details,omitempty"`
}

func (e *VaultError) Error() string {
	return e.Message
}

func NewError(status int, code string, message string) *VaultError {
	return &VaultError{
		Status:  status,
		Code:    code,
		Message: message,
	}
}

func (e *VaultError) WithDetails(details map[string]any) *VaultError {
	e.Details = details
	return e
}

func SendError(w http.ResponseWriter, err error) {
	ve, ok := err.(*VaultError)
	if !ok {
		ve = NewError(http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(ve.Status)
	_ = json.NewEncoder(w).Encode(map[string]*VaultError{"error": ve})
}
