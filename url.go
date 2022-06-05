package url

type Item struct {
	ShortURL    string `json:"short_url" db:"short_url"`
	OriginalURL string `json:"original_url" db:"base_url"`
}

type ShortenBatchInput struct {
	UserID        string `db:"user_id"`
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url" db:"base_url"`
	ShortenURL    string `db:"short_url"`
}

type ShortenBatchResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}
