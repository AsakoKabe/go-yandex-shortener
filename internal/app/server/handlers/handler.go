package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/AsakoKabe/go-yandex-shortener/internal/app/server/errs"
	"github.com/AsakoKabe/go-yandex-shortener/internal/logger"
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
	shortURL, err := h.urlShortener.Add(r.Context(), url)
	if errors.Is(err, errs.ErrConflictOriginalURL) {
		logger.Log.Info("original url already exist", zap.String("err", err.Error()))
		w.WriteHeader(http.StatusConflict)
	} else if err != nil {
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

	url, ok := h.urlShortener.Get(r.Context(), shortURL)
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
	w.Header().Set("Content-Type", "application/json")

	shortURL, err := h.urlShortener.Add(r.Context(), sr.URL)
	if errors.Is(err, errs.ErrConflictOriginalURL) {
		logger.Log.Info("original url already exist", zap.String("err", err.Error()))
		w.WriteHeader(http.StatusConflict)
	} else if err != nil {
		logger.Log.Error("error to create short url", zap.String("err", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		w.WriteHeader(http.StatusCreated)
	}

	response := ShortenerResponse{Result: h.prefixURL + shortURL}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		logger.Log.Error("error to create response", zap.String("err", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (h *Handler) createFromBatch(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var urlBatch []ShortenRequestBatch
	err := json.NewDecoder(r.Body).Decode(&urlBatch)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var shortURLBatch []ShortenResponseBatch
	for _, url := range urlBatch {
		shortURL, err := h.urlShortener.Add(r.Context(), url.OriginalURL)
		if err != nil {
			logger.Log.Error("error to create short url", zap.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		shortURLBatch = append(shortURLBatch, ShortenResponseBatch{
			ShortURL:      h.prefixURL + shortURL,
			CorrelationID: url.CorrelationID,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(shortURLBatch)
	if err != nil {
		logger.Log.Error("error to create response", zap.String("err", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func urlNotEmpty(url string) bool {
	return url == ""
}
