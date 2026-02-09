package realtime

import "github.com/zulfikawr/vault/internal/models"

type Message struct {
	Action     string         `json:"action"`
	Collection string         `json:"collection"`
	Record     *models.Record `json:"record"`
}
