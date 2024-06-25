package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"

	contextUtils "github.com/AsakoKabe/go-yandex-shortener/internal/app/context"
	"github.com/AsakoKabe/go-yandex-shortener/internal/app/server/errs"
	middlewareUtils "github.com/AsakoKabe/go-yandex-shortener/internal/app/server/middleware"
	"github.com/AsakoKabe/go-yandex-shortener/internal/app/shortener"
	"github.com/AsakoKabe/go-yandex-shortener/internal/logger"
)

type Handler struct {
	urlShortener shortener.URLShortener
	prefixURL    string
	deleteJobs   chan deleteJob
}

const numDeleteJobs = 5
const numWorkers = 5

type deleteJob struct {
	shortURL []string
	userID   string
}

func NewHandler(
	urlShortener shortener.URLShortener,
	prefixURL string,
) *Handler {
	jobs := make(chan deleteJob, numDeleteJobs)

	for w := 1; w <= numWorkers; w++ {
		go deleteWorker(urlShortener, jobs)
	}

	return &Handler{
		urlShortener: urlShortener,
		prefixURL:    prefixURL + "/",
		deleteJobs:   jobs,
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
	if isURLEmpty(url) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userID := contextUtils.GetUserID(r.Context())
	shortURL, err := h.urlShortener.Add(r.Context(), url, userID)
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

func (h *Handler) getURL(w http.ResponseWriter, r *http.Request) {
	shortURL := chi.URLParam(r, "id")

	if isURLEmpty(shortURL) {
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
	if isURLEmpty(url.OriginalURL) {
		logger.Log.Error("URL not found")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if url.DeletedFlag {
		w.WriteHeader(http.StatusGone)
		return
	}
	w.Header().Set("Location", url.OriginalURL)
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

	if isURLEmpty(sr.URL) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	userID := contextUtils.GetUserID(r.Context())
	shortURL, err := h.urlShortener.Add(r.Context(), sr.URL, userID)
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

	var originalURLs []string
	for _, originalURL := range urlBatch {
		originalURLs = append(originalURLs, originalURL.OriginalURL)
	}
	userID := contextUtils.GetUserID(r.Context())
	shortURLs, err := h.urlShortener.AddBatch(r.Context(), originalURLs, userID)
	if err != nil {
		logger.Log.Error("error to create short url", zap.String("err", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var shortURLBatch []ShortenResponseBatch
	for i, shortURL := range *shortURLs {
		shortURLBatch = append(shortURLBatch, ShortenResponseBatch{
			ShortURL:      h.prefixURL + shortURL,
			CorrelationID: urlBatch[i].CorrelationID,
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

func (h *Handler) getURLsByUser(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie(middlewareUtils.CookieName)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	userID := contextUtils.GetUserID(r.Context())
	urls, err := h.urlShortener.GetByUserID(r.Context(), userID)
	if err != nil {
		logger.Log.Error("error to get url")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var shortURLBatch []ShortenUserResponseBatch
	for _, url := range *urls {
		shortURLBatch = append(shortURLBatch, ShortenUserResponseBatch{
			ShortURL:    h.prefixURL + url.ShortURL,
			OriginalURL: url.OriginalURL,
		})
	}

	if len(shortURLBatch) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(shortURLBatch)
	if err != nil {
		logger.Log.Error("error to create response", zap.String("err", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}

func (h *Handler) deleteShorURLs(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var shortURLs []string
	err := json.NewDecoder(r.Body).Decode(&shortURLs)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID := contextUtils.GetUserID(r.Context())
	h.deleteJobs <- deleteJob{
		shortURL: shortURLs,
		userID:   userID,
	}

	w.WriteHeader(http.StatusAccepted)

}

func isURLEmpty(url string) bool {
	return url == ""
}

func deleteWorker(urlShortener shortener.URLShortener, jobs <-chan deleteJob) {
	for j := range jobs {
		err := urlShortener.DeleteShortURLs(context.Background(), j.shortURL, j.userID)
		if err != nil {
			logger.Log.Error("error to delete url", zap.String("err", err.Error()))
		}
	}
}
