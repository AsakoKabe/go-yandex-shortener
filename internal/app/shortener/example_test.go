package shortener

import (
	"context"
	"fmt"
)

func Example() {
	const maxLenShortURL = 6
	const fileStoragePath = "/tmp/mapper.json"

	// Создадим URLShortener с хранение в файле
	fm := NewFileURLMapper(maxLenShortURL, fileStoragePath)

	// Добавление и сжатие одного URL
	url := "ya.ru"
	userID := "123"
	shortURL, err := fm.Add(context.Background(), url, userID)
	if err != nil {
		return
	}
	fmt.Println("Short URL: " + shortURL + " Original URL: " + url)

	// Добавить батч из URL
	urls := []string{"ya.com", "google.com"}
	shortURLs, err := fm.AddBatch(context.Background(), urls, "")
	if err != nil {
		return
	}
	fmt.Println("Short URLs: ", shortURLs, "Original URLs: ", urls)

	// Получить оригинальный URL по сжатому
	originalURL, ok := fm.Get(context.Background(), url)
	if ok {
		fmt.Println("Short URL: ", url, "Original URL: ", originalURL)
	} else {
		fmt.Println("URL not found")
	}

	// Удалить сжатые URL для пользователя
	err = fm.DeleteShortURLs(context.Background(), *shortURLs, userID)
	if err != nil {
		return
	}
	fmt.Println("delete success")

	// Получить все сжатые URL для пользователя
	existedShortURLs, err := fm.GetByUserID(context.Background(), userID)
	if err != nil {
		return
	}
	fmt.Println("Short URLs: ", existedShortURLs, " userID", userID)
}
