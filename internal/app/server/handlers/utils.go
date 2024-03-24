package handlers

import (
	"io"
	"log"
)

func readBody(reqBody io.ReadCloser) (string, error) {
	body, err := io.ReadAll(reqBody)
	if err != nil {
		log.Fatalln(err)
		return "", err
	}

	return string(body), nil
}
