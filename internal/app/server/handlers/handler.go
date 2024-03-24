package handlers

import (
	"encoding/json"
	"github.com/AsakoKabe/go-yandex-shortener/internal/logger"
	"go.uber.org/zap"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	urlShortener URLShortener
	prefixURL    string
}

func NewHandler(
	urlShortener URLShortener,
	prefixURL string,
) *Handler {
	return &Handler{
		urlShortener: urlShortener,
		prefixURL:    prefixURL,
	}
}

func (h *Handler) createShortURL(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	url, err := readBody(r.Body)
	if err != nil {
		logger.Log.Error("error to read body", zap.String("err", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if h.emptyURL(url) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	shortURL, err := h.urlShortener.Add(url)
	if err != nil {
		logger.Log.Error("error to create short url", zap.String("err", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(h.prefixURL + shortURL))
}

func (h *Handler) emptyURL(url string) bool {
	return url == ""
}

func (h *Handler) getURL(w http.ResponseWriter, r *http.Request) {
	shortURL := "/" + chi.URLParam(r, "id")

	if urlNotEmpty(shortURL) {
		logger.Log.Error("shortURL not found")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	url, ok := h.urlShortener.Get(shortURL)
	if !ok {
		logger.Log.Error("error to get url")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if urlNotEmpty(url) {
		logger.Log.Error("URL not found")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *Handler) createShortURLJson(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var sr ShortenRequest
	err := json.NewDecoder(r.Body).Decode(&sr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if h.emptyURL(sr.URL) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	shortURL, err := h.urlShortener.Add(sr.URL)
	if err != nil {
		logger.Log.Error("error to create short url", zap.String("err", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response := ShortenResponse{Result: h.prefixURL + shortURL}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		logger.Log.Error("error to create response", zap.String("err", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func urlNotEmpty(url string) bool {
	return url == ""
}