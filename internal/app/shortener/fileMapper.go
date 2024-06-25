package shortener

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"strings"
	"sync"

	"go.uber.org/zap"

	"github.com/AsakoKabe/go-yandex-shortener/internal/app/shortener/models"
	"github.com/AsakoKabe/go-yandex-shortener/internal/app/utils"
	"github.com/AsakoKabe/go-yandex-shortener/internal/logger"
)

type FileURLMapper struct {
	mappingByShortURL sync.Map
	maxLenShortURL    int
	fileStoragePath   string
	fileMutex         sync.Mutex
}

func NewFileURLMapper(maxLenShortURL int, fileStoragePath string) *FileURLMapper {
	mapper := &FileURLMapper{
		maxLenShortURL:  maxLenShortURL,
		fileStoragePath: fileStoragePath,
	}
	err := mapper.loadFromFile()
	if err != nil {
		panic(err)
	}
	return mapper
}

func (m *FileURLMapper) Add(_ context.Context, url string, userID string) (string, error) {
	shortURL := utils.RandStringRunes(m.maxLenShortURL)
	su := models.URL{
		ShortURL:    shortURL,
		OriginalURL: url,
		UserID:      userID,
	}
	m.mappingByShortURL.Store(shortURL, su)
	err := m.saveToFile(su)
	if err != nil {
		return "", err
	}
	return shortURL, nil
}

func (m *FileURLMapper) AddBatch(_ context.Context, originalURLs []string, userID string) (*[]string, error) {
	var shortURLs []string
	for _, originalURL := range originalURLs {
		shortURL := utils.RandStringRunes(m.maxLenShortURL)
		su := models.URL{
			ShortURL:    shortURL,
			OriginalURL: originalURL,
			UserID:      userID,
		}
		m.mappingByShortURL.Store(shortURL, su)
		shortURLs = append(shortURLs, shortURL)
		err := m.saveToFile(su)
		if err != nil {
			return nil, err
		}
	}

	return &shortURLs, nil
}

func (m *FileURLMapper) Get(_ context.Context, shortURL string) (*models.URL, bool) {
	su, ok := m.mappingByShortURL.Load(shortURL)

	if ok {
		val, ok := su.(models.URL)
		if !ok {
			return nil, false
		}
		return &val, true
	}
	return nil, false
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
		m.mappingByShortURL.Store(su.ShortURL, su)
	}

	return nil
}

func (m *FileURLMapper) saveToFile(su models.URL) error {
	m.fileMutex.Lock()
	defer m.fileMutex.Unlock()

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

func (m *FileURLMapper) GetByUserID(_ context.Context, userID string) (*[]models.URL, error) {
	var urls []models.URL

	m.mappingByShortURL.Range(func(key, value interface{}) bool {
		url := value.(models.URL)
		if url.UserID == userID {
			urls = append(urls, url)
		}
		return true
	})

	return &urls, nil
}

func (m *FileURLMapper) DeleteShortURLs(_ context.Context, shortURLs []string, userID string) error {
	for _, shortURL := range shortURLs {
		value, ok := m.mappingByShortURL.Load(shortURL)
		if !ok {
			continue
		}
		url := value.(models.URL)
		if url.UserID == userID {
			url.DeletedFlag = true
		}
		m.mappingByShortURL.Store(shortURL, url)
	}

	return nil
}
