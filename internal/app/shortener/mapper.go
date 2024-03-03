package shortener

import "github.com/AsakoKabe/go-yandex-shortener/internal/app/utils"

type UrlMapper struct {
	mapping        map[string]string
	maxLenShortUrl int
}

func NewUrlMapper(maxLenShortUrl int) UrlShortener {
	return &UrlMapper{
		mapping:        make(map[string]string),
		maxLenShortUrl: maxLenShortUrl,
	}
}

func (m *UrlMapper) Add(url string) (string, error) {
	shortUrl := "/" + utils.RandStringRunes(m.maxLenShortUrl)
	m.mapping[shortUrl] = url
	return shortUrl, nil
}

func (m *UrlMapper) Get(shortUrl string) (string, error) {
	url, ok := m.mapping[shortUrl]
	if ok {
		return url, nil
	}
	return "", nil
}
