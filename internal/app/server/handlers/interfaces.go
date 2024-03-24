package handlers

type URLShortener interface {
	Add(url string) (string, error)
	Get(shortURL string) (string, bool)
}
