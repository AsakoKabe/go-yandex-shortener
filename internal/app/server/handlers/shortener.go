package handlers

import "context"

type URLShortener interface {
	Add(ctx context.Context, url string) (string, error)
	Get(ctx context.Context, shortURL string) (string, bool)
}
