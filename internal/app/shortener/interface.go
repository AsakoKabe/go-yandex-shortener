package shortener

type URLShortener interface {
	Add(url string) (string, error)
	Get(shortURL string) (string, error)
}
