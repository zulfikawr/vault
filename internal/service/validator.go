package service

import (
	"net/http"

	"github.com/zulfikawr/vault/internal/errors"
	"github.com/zulfikawr/vault/internal/models"
)

func ValidateRecord(col *models.Collection, data map[string]any) error {
	details := make(map[string]any)

	for _, f := range col.Fields {
		val, ok := data[f.Name]

		if f.Required && (!ok || val == nil || val == "") {
			details[f.Name] = "this field is required"
			continue
		}

		if ok && val != nil {
			// Type validation (simplified)
			switch f.Type {
			case models.FieldTypeNumber:
				if _, ok := val.(float64); !ok {
					if _, ok := val.(int); !ok {
						details[f.Name] = "must be a number"
					}
				}
			case models.FieldTypeBool:
				if _, ok := val.(bool); !ok {
					details[f.Name] = "must be a boolean"
				}
			}
		}
	}

	if len(details) > 0 {
		return errors.NewError(http.StatusBadRequest, "VALIDATION_FAILED", "Data validation failed").WithDetails(details)
	}

	return nil
}
