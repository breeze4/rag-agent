package storage

import (
	"context"
	"fmt"
	"strconv"

	chroma "github.com/amikos-tech/chroma-go"
	"github.com/amikos-tech/chroma-go/types"
)

type VectorStore struct {
	client     *chroma.Client
	collection *chroma.Collection
}

type DocumentChunk struct {
	ID         string            `json:"id"`
	Content    string            `json:"content"`
	DocumentID int               `json:"document_id"`
	ChunkIndex int               `json:"chunk_index"`
	Metadata   map[string]string `json:"metadata"`
}

type SearchResult struct {
	ID         string            `json:"id"`
	Content    string            `json:"content"`
	DocumentID int               `json:"document_id"`
	ChunkIndex int               `json:"chunk_index"`
	Score      float32           `json:"score"`
	Metadata   map[string]string `json:"metadata"`
}

func NewVectorStore(chromaURL string) (*VectorStore, error) {
	client, err := chroma.NewClient(chroma.WithBasePath(chromaURL))
	if err != nil {
		return nil, fmt.Errorf("failed to create chroma client: %w", err)
	}

	ctx := context.Background()
	
	// Test connection
	_, err = client.Heartbeat(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to chroma: %w", err)
	}

	return &VectorStore{
		client: client,
	}, nil
}

func (vs *VectorStore) EnsureCollection(name string) error {
	ctx := context.Background()
	
	// Check if collection exists
	collection, err := vs.client.GetCollection(ctx, name, nil)
	if err != nil {
		// Collection doesn't exist, create it
		collection, err = vs.client.CreateCollection(ctx, name, nil, true, nil, types.L2)
		if err != nil {
			return fmt.Errorf("failed to create collection: %w", err)
		}
	}
	
	vs.collection = collection
	return nil
}

func (vs *VectorStore) AddChunks(chunks []DocumentChunk, embeddings [][]float32) error {
	if len(chunks) != len(embeddings) {
		return fmt.Errorf("chunks and embeddings length mismatch: %d vs %d", len(chunks), len(embeddings))
	}

	ctx := context.Background()
	
	var ids []string
	var documents []string
	var metadatas []map[string]interface{}
	var embeddingsList [][]float32

	for i, chunk := range chunks {
		ids = append(ids, chunk.ID)
		documents = append(documents, chunk.Content)
		
		// Convert metadata to interface{} map
		metadata := make(map[string]interface{})
		metadata["document_id"] = strconv.Itoa(chunk.DocumentID)
		metadata["chunk_index"] = strconv.Itoa(chunk.ChunkIndex)
		
		// Add custom metadata
		for k, v := range chunk.Metadata {
			metadata[k] = v
		}
		
		metadatas = append(metadatas, metadata)
		embeddingsList = append(embeddingsList, embeddings[i])
	}

	// Convert embeddings to proper type
	chromaEmbeddings := types.NewEmbeddingsFromFloat32(embeddingsList)
	
	_, err := vs.collection.Add(ctx, chromaEmbeddings, metadatas, documents, ids)
	if err != nil {
		return fmt.Errorf("failed to add chunks to collection: %w", err)
	}

	return nil
}

func (vs *VectorStore) SearchSimilar(queryEmbedding []float32, limit int) ([]SearchResult, error) {
	ctx := context.Background()
	
	// Convert embedding to proper type
	embedding := types.NewEmbeddingFromFloat32(queryEmbedding)
	
	results, err := vs.collection.QueryWithOptions(ctx,
		types.WithQueryEmbedding(embedding),
		types.WithNResults(int32(limit)),
		types.WithInclude(types.IDocuments, types.IMetadatas, types.IDistances),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query collection: %w", err)
	}

	var searchResults []SearchResult
	
	if len(results.Documents) == 0 || len(results.Documents[0]) == 0 {
		return searchResults, nil
	}

	for i := 0; i < len(results.Documents[0]); i++ {
		var documentID int
		var chunkIndex int
		var metadata map[string]string
		
		if len(results.Metadatas) > 0 && len(results.Metadatas[0]) > i {
			meta := results.Metadatas[0][i]
			metadata = make(map[string]string)
			
			if docIDStr, ok := meta["document_id"].(string); ok {
				if docID, err := strconv.Atoi(docIDStr); err == nil {
					documentID = docID
				}
			}
			
			if chunkIndexStr, ok := meta["chunk_index"].(string); ok {
				if chunkIdx, err := strconv.Atoi(chunkIndexStr); err == nil {
					chunkIndex = chunkIdx
				}
			}
			
			// Extract other metadata
			for k, v := range meta {
				if k != "document_id" && k != "chunk_index" {
					if strVal, ok := v.(string); ok {
						metadata[k] = strVal
					}
				}
			}
		}

		var score float32
		if len(results.Distances) > 0 && len(results.Distances[0]) > i {
			// Convert distance to similarity score (lower distance = higher similarity)
			distance := results.Distances[0][i]
			score = 1.0 - distance
		}

		searchResult := SearchResult{
			ID:         results.Ids[0][i],
			Content:    results.Documents[0][i],
			DocumentID: documentID,
			ChunkIndex: chunkIndex,
			Score:      score,
			Metadata:   metadata,
		}
		
		searchResults = append(searchResults, searchResult)
	}

	return searchResults, nil
}

func (vs *VectorStore) DeleteChunks(chunkIDs []string) error {
	ctx := context.Background()
	
	_, err := vs.collection.Delete(ctx, chunkIDs, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to delete chunks: %w", err)
	}

	return nil
}

func (vs *VectorStore) DeleteByDocumentID(documentID int) error {
	ctx := context.Background()
	
	where := map[string]interface{}{
		"document_id": strconv.Itoa(documentID),
	}
	
	_, err := vs.collection.Delete(ctx, nil, where, nil)
	if err != nil {
		return fmt.Errorf("failed to delete chunks by document ID: %w", err)
	}

	return nil
}

func (vs *VectorStore) GetCollectionInfo() (map[string]interface{}, error) {
	ctx := context.Background()
	
	count, err := vs.collection.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get collection count: %w", err)
	}

	info := map[string]interface{}{
		"name":  vs.collection.Name,
		"count": count,
	}

	return info, nil
}

func (vs *VectorStore) DeleteCollection(collectionName string) error {
	ctx := context.Background()
	
	_, err := vs.client.DeleteCollection(ctx, collectionName)
	if err != nil {
		return fmt.Errorf("failed to delete collection: %w", err)
	}
	
	return nil
}