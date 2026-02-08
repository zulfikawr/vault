package models

type FieldType string

const (
	FieldTypeText     FieldType = "text"
	FieldTypeNumber   FieldType = "number"
	FieldTypeBool     FieldType = "bool"
	FieldTypeJSON     FieldType = "json"
	FieldTypeDate     FieldType = "date"
	FieldTypeRelation FieldType = "relation"
)

type Field struct {
	Name     string    `json:"name"`
	Type     FieldType `json:"type"`
	Required bool      `json:"required"`
	Unique   bool      `json:"unique"`
	Options  any       `json:"options,omitempty"`
}
