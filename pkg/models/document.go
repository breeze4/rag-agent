package models

import "time"

type Document struct {
	ID          int       `json:"id" db:"id"`
	FileName    string    `json:"file_name" db:"file_name"`
	FilePath    string    `json:"file_path" db:"file_path"`
	FileSize    int64     `json:"file_size" db:"file_size"`
	ContentHash string    `json:"content_hash" db:"content_hash"`
	UploadedAt  time.Time `json:"uploaded_at" db:"uploaded_at"`
	ProcessedAt *time.Time `json:"processed_at,omitempty" db:"processed_at"`
	Status      string    `json:"status" db:"status"` // "pending", "processing", "completed", "failed"
}

const (
	DocumentStatusPending    = "pending"
	DocumentStatusProcessing = "processing"
	DocumentStatusCompleted  = "completed"
	DocumentStatusFailed     = "failed"
)