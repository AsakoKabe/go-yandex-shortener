package shortener

import (
	"github.com/AsakoKabe/go-yandex-shortener/internal/app/utils"
)

type URLMapper struct {
	mapping        map[string]string
	maxLenShortURL int
}

func NewURLMapper(maxLenShortURL int) *URLMapper {
	return &URLMapper{
		mapping:        make(map[string]string),
		maxLenShortURL: maxLenShortURL,
	}
}

func (m *URLMapper) Add(url string) (string, error) {
	shortURL := "/" + utils.RandStringRunes(m.maxLenShortURL)
	m.mapping[shortURL] = url
	return shortURL, nil
}

func (m *URLMapper) Get(shortURL string) (string, bool) {
	url, ok := m.mapping[shortURL]
	if ok {
		return url, true
	}
	return "", false
}
