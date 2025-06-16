package storage

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type FileStorage struct {
	dataDir string
}

func NewFileStorage(dataDir string) (*FileStorage, error) {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	documentsDir := filepath.Join(dataDir, "documents")
	if err := os.MkdirAll(documentsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create documents directory: %w", err)
	}

	return &FileStorage{
		dataDir: dataDir,
	}, nil
}

func (fs *FileStorage) SaveDocument(fileName string, content io.Reader) (string, string, int64, error) {
	hash := sha256.New()
	
	tempFile, err := os.CreateTemp(fs.dataDir, "upload_*.tmp")
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tempFile.Name())

	teeReader := io.TeeReader(content, hash)
	size, err := io.Copy(tempFile, teeReader)
	if err != nil {
		tempFile.Close()
		return "", "", 0, fmt.Errorf("failed to copy content: %w", err)
	}
	tempFile.Close()

	contentHash := fmt.Sprintf("%x", hash.Sum(nil))
	
	timestamp := time.Now().Format("20060102_150405")
	safeFileName := fmt.Sprintf("%s_%s", timestamp, fileName)
	finalPath := filepath.Join(fs.dataDir, "documents", safeFileName)

	if err := os.Rename(tempFile.Name(), finalPath); err != nil {
		return "", "", 0, fmt.Errorf("failed to move file to final location: %w", err)
	}

	return finalPath, contentHash, size, nil
}

func (fs *FileStorage) DeleteDocument(filePath string) error {
	return os.Remove(filePath)
}

func (fs *FileStorage) DocumentExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}