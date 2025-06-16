package storage

import (
	"fmt"
)

const DefaultCollectionName = "document_chunks"

type VectorService struct {
	store *VectorStore
}

func NewVectorService(chromaURL string) (*VectorService, error) {
	store, err := NewVectorStore(chromaURL)
	if err != nil {
		return nil, err
	}

	// Ensure the default collection exists
	if err := store.EnsureCollection(DefaultCollectionName); err != nil {
		return nil, fmt.Errorf("failed to ensure collection: %w", err)
	}

	return &VectorService{
		store: store,
	}, nil
}

func (vs *VectorService) StoreDocumentChunks(documentID int, chunks []string, embeddings [][]float32, metadata map[string]string) error {
	if len(chunks) != len(embeddings) {
		return fmt.Errorf("chunks and embeddings length mismatch")
	}

	var documentChunks []DocumentChunk
	for i, chunk := range chunks {
		chunkID := fmt.Sprintf("doc_%d_chunk_%d", documentID, i)
		
		// Create chunk metadata
		chunkMetadata := make(map[string]string)
		for k, v := range metadata {
			chunkMetadata[k] = v
		}
		
		documentChunk := DocumentChunk{
			ID:         chunkID,
			Content:    chunk,
			DocumentID: documentID,
			ChunkIndex: i,
			Metadata:   chunkMetadata,
		}
		
		documentChunks = append(documentChunks, documentChunk)
	}

	return vs.store.AddChunks(documentChunks, embeddings)
}

func (vs *VectorService) SearchRelevantChunks(queryEmbedding []float32, limit int) ([]SearchResult, error) {
	return vs.store.SearchSimilar(queryEmbedding, limit)
}

func (vs *VectorService) DeleteDocumentChunks(documentID int) error {
	return vs.store.DeleteByDocumentID(documentID)
}

func (vs *VectorService) GetStats() (map[string]interface{}, error) {
	return vs.store.GetCollectionInfo()
}

// Helper function to generate chunk ID
func GenerateChunkID(documentID, chunkIndex int) string {
	return fmt.Sprintf("doc_%d_chunk_%d", documentID, chunkIndex)
}

// Helper function to parse chunk ID
func ParseChunkID(chunkID string) (documentID, chunkIndex int, err error) {
	var docID, chunkIdx int
	_, err = fmt.Sscanf(chunkID, "doc_%d_chunk_%d", &docID, &chunkIdx)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse chunk ID: %w", err)
	}
	return docID, chunkIdx, nil
}