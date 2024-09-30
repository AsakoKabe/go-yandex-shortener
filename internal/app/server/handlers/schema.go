package handlers

// ShortenRequest Структура для запроса создания сжатого URL
type ShortenRequest struct {
	URL string `json:"url"`
}

// ShortenerResponse Структура для ответа сжатия URL
type ShortenerResponse struct {
	Result string `json:"result"`
}

// ShortenRequestBatch Структура для запроса создания батча сжатых URL
type ShortenRequestBatch struct {
	OriginalURL   string `json:"original_url"`
	CorrelationID string `json:"correlation_id"`
}

// ShortenResponseBatch Структура для ответа батча сжатых URL
type ShortenResponseBatch struct {
	ShortURL      string `json:"short_url"`
	CorrelationID string `json:"correlation_id"`
}

// ShortenUserResponseBatch Структура для ответа сжатых URL по пользователю
type ShortenUserResponseBatch struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// InternalStats Структура с внутренней статистикой
type InternalStats struct {
	Urls  int `json:"urls"`
	Users int `json:"users"`
}
