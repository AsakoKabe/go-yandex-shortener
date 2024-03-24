package handlers

import (
	"log"
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
		log.Printf("error to read body, err: %+v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if h.emptyURL(url) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	shortURL, err := h.urlShortener.Add(url)
	if err != nil {
		log.Printf("error to create short url, err: %+v", err)
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
		log.Printf("shortURL not found")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	url, ok := h.urlShortener.Get(shortURL)
	if !ok {
		log.Printf("error to get url")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if urlNotEmpty(url) {
		log.Printf("URL not found")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func urlNotEmpty(url string) bool {
	return url == ""
}
