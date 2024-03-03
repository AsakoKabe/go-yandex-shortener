package endpoints

import (
	"github.com/AsakoKabe/go-yandex-shortener/internal/app/shortener"
	"log"
	"net/http"
)

type Handler struct {
	urlShortener shortener.UrlShortener
}

func NewHandler(urlShortener shortener.UrlShortener) *Handler {
	return &Handler{urlShortener: urlShortener}
}

func (h *Handler) root(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.createShortUrl(w, r)
	case http.MethodGet:
		h.getUrl(w, r)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}

}

func (h *Handler) createShortUrl(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	url, err := readBody(r.Body)
	if err != nil {
		log.Printf("error to read body, err: %+v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if h.emptyUrl(url) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	shortUlr, err := h.urlShortener.Add(url)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(r.Host + shortUlr))
}

func (h *Handler) emptyUrl(url string) bool {
	return url == ""
}

func (h *Handler) getUrl(w http.ResponseWriter, r *http.Request) {
	url, err := h.urlShortener.Get(r.URL.Path)
	if err != nil {
		log.Printf("error to get url, err: %+v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if h.shortUrlNotExist(url) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *Handler) shortUrlNotExist(url string) bool {
	return url == ""
}
