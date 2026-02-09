package models

type FileMetadata struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
	Mime string `json:"mime"`
	URL  string `json:"url"`
}
