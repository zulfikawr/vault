package models

type Record struct {
	ID         string         `json:"id"`
	Collection string         `json:"collection"`
	Data       map[string]any `json:"data"`
	Expand     map[string]any `json:"expand,omitempty"`
	Created    string         `json:"created"`
	Updated    string         `json:"updated"`
}

func (r *Record) HideField(name string) {
	delete(r.Data, name)
}

func NewRecord(collection string) *Record {
	return &Record{
		Collection: collection,
		Data:       make(map[string]any),
		Expand:     make(map[string]any),
	}
}

func (r *Record) GetString(key string) string {
	if val, ok := r.Data[key].(string); ok {
		return val
	}
	return ""
}

func (r *Record) GetInt(key string) int {
	if val, ok := r.Data[key].(float64); ok { // JSON unmarshals numbers as float64
		return int(val)
	}
	if val, ok := r.Data[key].(int); ok {
		return val
	}
	return 0
}
