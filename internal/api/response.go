package api

import (
	"encoding/json"
	"math"
	"net/http"
)

type Response struct {
	Data any `json:"data"`
	Meta any `json:"meta,omitempty"`
}

type PaginatedResponse struct {
	Page       int `json:"page"`
	PerPage    int `json:"perPage"`
	TotalItems int `json:"totalItems"`
	TotalPages int `json:"totalPages"`
	Items      any `json:"items"`
}

func SendJSON(w http.ResponseWriter, status int, data any, meta any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// If meta contains pagination info, we might want to use PaginatedResponse
	if pagination, ok := meta.(map[string]any); ok {
		page, hasPage := pagination["page"].(int)
		perPage, hasPerPage := pagination["perPage"].(int)
		totalItems, hasTotalItems := pagination["totalItems"].(int)

		if hasPage && hasPerPage && hasTotalItems {
			totalPages := 0
			if perPage > 0 {
				totalPages = int(math.Ceil(float64(totalItems) / float64(perPage)))
			}

			_ = json.NewEncoder(w).Encode(PaginatedResponse{
				Page:       page,
				PerPage:    perPage,
				TotalItems: totalItems,
				TotalPages: totalPages,
				Items:      data,
			})
			return
		}
	}

	resp := Response{
		Data: data,
		Meta: meta,
	}

	_ = json.NewEncoder(w).Encode(resp)
}
