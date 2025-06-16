package storage

import (
	"database/sql"
	"fmt"
	"time"

	"rag-therapist/pkg/models"
)

type DocumentRepository struct {
	db *Database
}

func NewDocumentRepository(db *Database) *DocumentRepository {
	return &DocumentRepository{db: db}
}

func (r *DocumentRepository) Insert(doc *models.Document) error {
	query := `
		INSERT INTO documents (file_name, file_path, file_size, content_hash, uploaded_at, status)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	
	result, err := r.db.db.Exec(query, doc.FileName, doc.FilePath, doc.FileSize, doc.ContentHash, doc.UploadedAt, doc.Status)
	if err != nil {
		return fmt.Errorf("failed to insert document: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	doc.ID = int(id)
	return nil
}

func (r *DocumentRepository) GetByID(id int) (*models.Document, error) {
	query := `
		SELECT id, file_name, file_path, file_size, content_hash, uploaded_at, processed_at, status
		FROM documents WHERE id = ?
	`
	
	row := r.db.db.QueryRow(query, id)
	
	var doc models.Document
	var processedAt sql.NullTime
	
	err := row.Scan(&doc.ID, &doc.FileName, &doc.FilePath, &doc.FileSize, &doc.ContentHash, &doc.UploadedAt, &processedAt, &doc.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("document not found")
		}
		return nil, fmt.Errorf("failed to scan document: %w", err)
	}

	if processedAt.Valid {
		doc.ProcessedAt = &processedAt.Time
	}

	return &doc, nil
}

func (r *DocumentRepository) GetByContentHash(hash string) (*models.Document, error) {
	query := `
		SELECT id, file_name, file_path, file_size, content_hash, uploaded_at, processed_at, status
		FROM documents WHERE content_hash = ?
	`
	
	row := r.db.db.QueryRow(query, hash)
	
	var doc models.Document
	var processedAt sql.NullTime
	
	err := row.Scan(&doc.ID, &doc.FileName, &doc.FilePath, &doc.FileSize, &doc.ContentHash, &doc.UploadedAt, &processedAt, &doc.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("document not found")
		}
		return nil, fmt.Errorf("failed to scan document: %w", err)
	}

	if processedAt.Valid {
		doc.ProcessedAt = &processedAt.Time
	}

	return &doc, nil
}

func (r *DocumentRepository) List(limit, offset int) ([]*models.Document, error) {
	query := `
		SELECT id, file_name, file_path, file_size, content_hash, uploaded_at, processed_at, status
		FROM documents ORDER BY uploaded_at DESC LIMIT ? OFFSET ?
	`
	
	rows, err := r.db.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query documents: %w", err)
	}
	defer rows.Close()

	var documents []*models.Document
	for rows.Next() {
		var doc models.Document
		var processedAt sql.NullTime
		
		err := rows.Scan(&doc.ID, &doc.FileName, &doc.FilePath, &doc.FileSize, &doc.ContentHash, &doc.UploadedAt, &processedAt, &doc.Status)
		if err != nil {
			return nil, fmt.Errorf("failed to scan document: %w", err)
		}

		if processedAt.Valid {
			doc.ProcessedAt = &processedAt.Time
		}

		documents = append(documents, &doc)
	}

	return documents, nil
}

func (r *DocumentRepository) UpdateStatus(id int, status string, processedAt *time.Time) error {
	query := `UPDATE documents SET status = ?, processed_at = ? WHERE id = ?`
	
	_, err := r.db.db.Exec(query, status, processedAt, id)
	if err != nil {
		return fmt.Errorf("failed to update document status: %w", err)
	}

	return nil
}

func (r *DocumentRepository) Delete(id int) error {
	query := `DELETE FROM documents WHERE id = ?`
	
	_, err := r.db.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}

	return nil
}

func (r *DocumentRepository) GetByStatus(status string) ([]*models.Document, error) {
	query := `
		SELECT id, file_name, file_path, file_size, content_hash, uploaded_at, processed_at, status
		FROM documents WHERE status = ? ORDER BY uploaded_at ASC
	`
	
	rows, err := r.db.db.Query(query, status)
	if err != nil {
		return nil, fmt.Errorf("failed to query documents by status: %w", err)
	}
	defer rows.Close()

	var documents []*models.Document
	for rows.Next() {
		var doc models.Document
		var processedAt sql.NullTime
		
		err := rows.Scan(&doc.ID, &doc.FileName, &doc.FilePath, &doc.FileSize, &doc.ContentHash, &doc.UploadedAt, &processedAt, &doc.Status)
		if err != nil {
			return nil, fmt.Errorf("failed to scan document: %w", err)
		}

		if processedAt.Valid {
			doc.ProcessedAt = &processedAt.Time
		}

		documents = append(documents, &doc)
	}

	return documents, nil
}