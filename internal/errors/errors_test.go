package errors

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
)

func TestVaultError(t *testing.T) {
	err := NewError(400, "VALIDATION_ERROR", "invalid field")
	if err.Error() != "invalid field" {
		t.Errorf("expected 'invalid field', got %s", err.Error())
	}

	w := httptest.NewRecorder()
	SendError(w, err)

	if w.Code != 400 {
		t.Errorf("expected status 400, got %d", w.Code)
	}

	var resp map[string]VaultError
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}

	if resp["error"].Code != "VALIDATION_ERROR" {
		t.Errorf("expected code VALIDATION_ERROR, got %s", resp["error"].Code)
	}
}
