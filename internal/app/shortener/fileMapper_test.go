package shortener

import (
	"context"
	"log/slog"
	"math/rand"
	"os"
	"testing"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

const triesN = 1000

const maxLenShortURL = 6
const fileStoragePath = "/tmp/mapper.json"

func deleteFileStorageMapping() {
	err := os.Remove(fileStoragePath)
	if os.IsNotExist(err) {
		return
	}
	if err != nil {
		slog.Error("error to delete file storage mapping for benchmark")
	}
}

func BenchmarkFileURLMapper_Add(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		fm := NewFileURLMapper(maxLenShortURL, fileStoragePath)
		url := randSeq(10)
		b.StartTimer()
		_, _ = fm.Add(context.Background(), url, "")
		b.StopTimer()
		deleteFileStorageMapping()
	}
}

func BenchmarkFileURLMapper_AddBatch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		fm := NewFileURLMapper(maxLenShortURL, fileStoragePath)
		urls := []string{randSeq(10), randSeq(10), randSeq(10), randSeq(10)}
		b.StartTimer()
		_, _ = fm.AddBatch(context.Background(), urls, "")
		b.StopTimer()
		deleteFileStorageMapping()
	}
}

func BenchmarkFileURLMapper_Get(b *testing.B) {
	b.Run(
		"get not existed", func(b *testing.B) {
			for i := 0; i < triesN; i++ {
				b.StopTimer()
				fm := NewFileURLMapper(maxLenShortURL, fileStoragePath)
				url := randSeq(10)
				b.StartTimer()
				_, _ = fm.Get(context.Background(), url)
				b.StopTimer()
				deleteFileStorageMapping()
			}
		},
	)

	b.Run(
		"get existed", func(b *testing.B) {
			for i := 0; i < triesN; i++ {
				b.StopTimer()
				fm := NewFileURLMapper(maxLenShortURL, fileStoragePath)
				url := randSeq(10)
				shortURL, _ := fm.Add(context.Background(), url, "")
				b.StartTimer()
				_, _ = fm.Get(context.Background(), shortURL)
				b.StopTimer()
				deleteFileStorageMapping()
			}
		},
	)
}

func BenchmarkFileURLMapper_DeleteShortURLs(b *testing.B) {
	b.Run(
		"delete not existed", func(b *testing.B) {
			for i := 0; i < triesN; i++ {
				b.StopTimer()
				fm := NewFileURLMapper(maxLenShortURL, fileStoragePath)
				urls := []string{randSeq(10), randSeq(10), randSeq(10), randSeq(10)}
				b.StartTimer()
				_ = fm.DeleteShortURLs(context.Background(), urls, "")
				b.StopTimer()
				deleteFileStorageMapping()
			}
		},
	)

	b.Run(
		"delete existed", func(b *testing.B) {
			for i := 0; i < triesN; i++ {
				b.StopTimer()
				fm := NewFileURLMapper(maxLenShortURL, fileStoragePath)
				urls := []string{randSeq(10), randSeq(10), randSeq(10), randSeq(10)}
				shortURLs, _ := fm.AddBatch(context.Background(), urls, "")
				b.StartTimer()
				_ = fm.DeleteShortURLs(context.Background(), *shortURLs, "")
				b.StopTimer()
				deleteFileStorageMapping()
			}
		},
	)
}

func BenchmarkFileURLMapper_GetByUserID(b *testing.B) {
	b.Run(
		"get by userID not existed", func(b *testing.B) {
			for i := 0; i < triesN; i++ {
				b.StopTimer()
				fm := NewFileURLMapper(maxLenShortURL, fileStoragePath)
				b.StartTimer()
				_, _ = fm.GetByUserID(context.Background(), "")
				b.StopTimer()
				deleteFileStorageMapping()
			}
		},
	)

	b.Run(
		"get by userID existed", func(b *testing.B) {
			for i := 0; i < triesN; i++ {
				b.StopTimer()
				fm := NewFileURLMapper(maxLenShortURL, fileStoragePath)
				_, _ = fm.Add(context.Background(), randSeq(10), "")
				b.StartTimer()
				_, _ = fm.GetByUserID(context.Background(), "")
				b.StopTimer()
				deleteFileStorageMapping()
			}
		},
	)
}
