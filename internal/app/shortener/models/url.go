package models

type URL struct {
	ID          int
	UserID      string
	ShortURL    string `json:"short_url,omitempty"`
	OriginalURL string `json:"original_url,omitempty"`
}
