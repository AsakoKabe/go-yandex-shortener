package handlers

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenerResponse struct {
	Result string `json:"result"`
}

type ShortenRequestBatch struct {
	OriginalURL   string `json:"original_url"`
	CorrelationID string `json:"correlation_id"`
}

type ShortenResponseBatch struct {
	ShortURL      string `json:"short_url"`
	CorrelationID string `json:"correlation_id"`
}
