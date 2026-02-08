package models

type CollectionType string

const (
	CollectionTypeBase   CollectionType = "base"
	CollectionTypeAuth   CollectionType = "auth"
	CollectionTypeSystem CollectionType = "system"
)

type Collection struct {
	ID      string         `json:"id"`
	Name    string         `json:"name"`
	Type    CollectionType `json:"type"`
	Fields  []Field        `json:"fields"`
	Indexes []string       `json:"indexes"`
	
	// API Rules (simple string filters for now)
	ListRule   *string `json:"list_rule"`
	ViewRule   *string `json:"view_rule"`
	CreateRule *string `json:"create_rule"`
	UpdateRule *string `json:"update_rule"`
	DeleteRule *string `json:"delete_rule"`

	Created string `json:"created"`
	Updated string `json:"updated"`
}
