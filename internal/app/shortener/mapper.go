package shortener

import (
	"encoding/json"
	"github.com/AsakoKabe/go-yandex-shortener/internal/app/utils"
	"github.com/AsakoKabe/go-yandex-shortener/internal/logger"
	"go.uber.org/zap"
	"io"
	"os"
	"strings"
)

type URLMapper struct {
	mapping         map[string]ShortURL
	maxLenShortURL  int
	fileStoragePath string
}

func NewURLMapper(maxLenShortURL int, fileStoragePath string) *URLMapper {
	mapper := &URLMapper{
		mapping:         make(map[string]ShortURL),
		maxLenShortURL:  maxLenShortURL,
		fileStoragePath: fileStoragePath,
	}
	err := mapper.loadFromFile()
	if err != nil {
		panic(err)
	}
	return mapper
}

func (m *URLMapper) Add(url string) (string, error) {
	shortURL := "/" + utils.RandStringRunes(m.maxLenShortURL)
	su := ShortURL{
		ShortURL:    shortURL,
		OriginalURL: url,
	}
	m.mapping[shortURL] = su
	err := m.saveToFile(su)
	if err != nil {
		return "", err
	}
	return shortURL, nil
}

func (m *URLMapper) Get(shortURL string) (string, bool) {
	su, ok := m.mapping[shortURL]
	if ok {
		return su.OriginalURL, true
	}
	return "", false
}

func (m *URLMapper) loadFromFile() error {
	data, err := os.ReadFile(m.fileStoragePath)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		logger.Log.Error(
			"error to read file",
			zap.String("file path", m.fileStoragePath),
			zap.String("err", err.Error()),
		)
		return err
	}

	dec := json.NewDecoder(strings.NewReader(string(data)))
	for {
		var su ShortURL

		err = dec.Decode(&su)
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Log.Error("error to parse json", zap.String("err", err.Error()))
			return err
		}
		m.mapping[su.ShortURL] = su
	}

	return nil
}

func (m *URLMapper) saveToFile(su ShortURL) error {
	content, _ := json.Marshal(su)

	f, err := os.OpenFile(m.fileStoragePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		logger.Log.Error(
			"error to open file",
			zap.String("file path", m.fileStoragePath),
			zap.String("err", err.Error()),
		)
	}
	defer f.Close()

	_, err = f.Write(content)
	if err != nil {
		logger.Log.Error(
			"error to save data to file",
			zap.String("file path", m.fileStoragePath),
			zap.String("err", err.Error()),
		)
		return err
	}
	return nil
}
