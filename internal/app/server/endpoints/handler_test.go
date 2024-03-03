package endpoints

import (
	"github.com/AsakoKabe/go-yandex-shortener/internal/app/shortener"
	"github.com/AsakoKabe/go-yandex-shortener/internal/app/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_createShortURL(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name      string
		want      want
		body      io.Reader
		shortener shortener.URLShortener
	}{
		{
			name: "simple positive",
			want: want{
				code:        http.StatusCreated,
				contentType: "text/plain",
			},
			body:      strings.NewReader("https://ya.ru"),
			shortener: shortener.NewURLMapper(5),
		},
		{
			name: "return status 400 for empty url",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "",
			},
			body:      strings.NewReader(""),
			shortener: shortener.NewURLMapper(5),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/", test.body)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			h := NewHandler(test.shortener)
			h.createShortURL(w, request)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			if res.StatusCode != http.StatusBadRequest {
				assert.NotEmpty(t, string(resBody))
			}
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}

func Test_getURL(t *testing.T) {
	urlMap, h := setUpSimple()

	for url, shortURL := range urlMap {
		t.Run("positive, url: "+url, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, shortURL, nil)
			w := httptest.NewRecorder()
			h.getURL(w, request)
			res := w.Result()

			assert.Equal(t, http.StatusTemporaryRedirect, res.StatusCode)
			assert.Equal(t, res.Header.Get("Location"), url)
		})
	}

	shortURL := "/" + utils.RandStringRunes(5)
	t.Run("negative, random shortURL: "+shortURL, func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, shortURL, nil)
		w := httptest.NewRecorder()
		h.getURL(w, request)
		res := w.Result()

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})
}

func setUpSimple() (map[string]string, *Handler) {
	urlMap := make(map[string]string)
	urls := []string{
		"https://ya.ru",
		"https://example.ru",
	}
	h := NewHandler(shortener.NewURLMapper(5))
	for _, url := range urls {
		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(url))
		w := httptest.NewRecorder()
		h.createShortURL(w, request)
		res := w.Result()
		resBody, _ := io.ReadAll(res.Body)
		short := strings.Split(string(resBody), "/")
		urlMap[url] = "/" + short[(len(short)-1)]
	}
	return urlMap, h
}
