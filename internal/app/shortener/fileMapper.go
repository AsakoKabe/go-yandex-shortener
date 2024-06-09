package shortener

import (
	"context"
	"encoding/json"
	"github.com/AsakoKabe/go-yandex-shortener/internal/app/shortener/models"
	"github.com/AsakoKabe/go-yandex-shortener/internal/app/utils"
	"github.com/AsakoKabe/go-yandex-shortener/internal/logger"
	"go.uber.org/zap"
	"io"
	"os"
	"strings"
)

type FileURLMapper struct {
	mapping         map[string]models.URL
	maxLenShortURL  int
	fileStoragePath string
}

func NewFileURLMapper(maxLenShortURL int, fileStoragePath string) *FileURLMapper {
	mapper := &FileURLMapper{
		mapping:         make(map[string]models.URL),
		maxLenShortURL:  maxLenShortURL,
		fileStoragePath: fileStoragePath,
	}
	err := mapper.loadFromFile()
	if err != nil {
		panic(err)
	}
	return mapper
}

func (m *FileURLMapper) Add(_ context.Context, url string) (string, error) {
	shortURL := "/" + utils.RandStringRunes(m.maxLenShortURL)
	su := models.URL{
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

func (m *FileURLMapper) Get(_ context.Context, shortURL string) (string, bool) {
	su, ok := m.mapping[shortURL]
	if ok {
		return su.OriginalURL, true
	}
	return "", false
}

func (m *FileURLMapper) loadFromFile() error {
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
		var su models.URL

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

func (m *FileURLMapper) saveToFile(su models.URL) error {
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
