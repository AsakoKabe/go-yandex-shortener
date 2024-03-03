package shortener

type UrlShortener interface {
	Add(url string) (string, error)
	Get(shortUrl string) (string, error)
}
