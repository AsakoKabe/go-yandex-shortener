package endpoints

import (
	"github.com/AsakoKabe/go-yandex-shortener/internal/app/shortener"
	"log"
	"net/http"
)

type Handler struct {
	urlShortener shortener.URLShortener
}

func NewHandler(urlShortener shortener.URLShortener) *Handler {
	return &Handler{urlShortener: urlShortener}
}

func (h *Handler) root(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.createShortURL(w, r)
	case http.MethodGet:
		h.getURL(w, r)
	default:
		w.WriteHeader(http.StatusBadRequest)
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
	w.Write([]byte("http://" + r.Host + shortURL))
}

func (h *Handler) emptyURL(url string) bool {
	return url == ""
}

func (h *Handler) getURL(w http.ResponseWriter, r *http.Request) {
	url, err := h.urlShortener.Get(r.URL.Path)
	if err != nil {
		log.Printf("error to get url, err: %+v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if h.shortURLNotExist(url) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *Handler) shortURLNotExist(url string) bool {
	return url == ""
}
