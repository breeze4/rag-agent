package storage

import (
	"io"
	"time"

	"rag-therapist/pkg/models"
)

type StorageService struct {
	fileStorage *FileStorage
	docRepo     *DocumentRepository
}

func NewStorageService(dataDir string) (*StorageService, error) {
	database, err := NewDatabase(dataDir)
	if err != nil {
		return nil, err
	}

	fileStorage, err := NewFileStorage(dataDir)
	if err != nil {
		return nil, err
	}

	docRepo := NewDocumentRepository(database)

	return &StorageService{
		fileStorage: fileStorage,
		docRepo:     docRepo,
	}, nil
}

func (s *StorageService) StoreDocument(fileName string, content io.Reader) (*models.Document, error) {
	filePath, contentHash, fileSize, err := s.fileStorage.SaveDocument(fileName, content)
	if err != nil {
		return nil, err
	}

	existing, err := s.docRepo.GetByContentHash(contentHash)
	if err == nil {
		s.fileStorage.DeleteDocument(filePath)
		return existing, nil
	}

	doc := &models.Document{
		FileName:    fileName,
		FilePath:    filePath,
		FileSize:    fileSize,
		ContentHash: contentHash,
		UploadedAt:  time.Now(),
		Status:      models.DocumentStatusPending,
	}

	if err := s.docRepo.Insert(doc); err != nil {
		s.fileStorage.DeleteDocument(filePath)
		return nil, err
	}

	return doc, nil
}

func (s *StorageService) GetDocument(id int) (*models.Document, error) {
	return s.docRepo.GetByID(id)
}

func (s *StorageService) ListDocuments(limit, offset int) ([]*models.Document, error) {
	return s.docRepo.List(limit, offset)
}

func (s *StorageService) UpdateDocumentStatus(id int, status string) error {
	now := time.Now()
	return s.docRepo.UpdateStatus(id, status, &now)
}

func (s *StorageService) GetPendingDocuments() ([]*models.Document, error) {
	return s.docRepo.GetByStatus(models.DocumentStatusPending)
}

func (s *StorageService) DeleteDocument(id int) error {
	doc, err := s.docRepo.GetByID(id)
	if err != nil {
		return err
	}

	if err := s.fileStorage.DeleteDocument(doc.FilePath); err != nil {
		return err
	}

	return s.docRepo.Delete(id)
}