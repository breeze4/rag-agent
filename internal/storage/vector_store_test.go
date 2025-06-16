package storage

import (
	"os"
	"testing"
	"time"
)

func TestVectorStoreOperations(t *testing.T) {
	// Skip test if no Chroma URL is provided
	chromaURL := os.Getenv("CHROMA_URL")
	if chromaURL == "" {
		chromaURL = "http://localhost:8000"
	}

	// Create vector store
	vs, err := NewVectorStore(chromaURL)
	if err != nil {
		t.Skipf("Skipping test: failed to connect to Chroma at %s: %v", chromaURL, err)
	}

	// Test collection creation
	testCollectionName := "test_collection_" + time.Now().Format("20060102_150405")
	err = vs.EnsureCollection(testCollectionName)
	if err != nil {
		t.Fatalf("Failed to create collection: %v", err)
	}

	// Clean up collection after test
	defer func() {
		vs.DeleteCollection(testCollectionName)
	}()

	// Test data
	chunks := []DocumentChunk{
		{
			ID:         "test_1",
			Content:    "This is a test document about artificial intelligence.",
			DocumentID: 1,
			ChunkIndex: 0,
			Metadata:   map[string]string{"topic": "ai"},
		},
		{
			ID:         "test_2",
			Content:    "Machine learning is a subset of artificial intelligence.",
			DocumentID: 1,
			ChunkIndex: 1,
			Metadata:   map[string]string{"topic": "ml"},
		},
	}

	// Mock embeddings (in real usage, these would come from OpenAI API)
	embeddings := [][]float32{
		{0.1, 0.2, 0.3, 0.4, 0.5},
		{0.2, 0.3, 0.4, 0.5, 0.6},
	}

	// Test adding chunks
	err = vs.AddChunks(chunks, embeddings)
	if err != nil {
		t.Fatalf("Failed to add chunks: %v", err)
	}

	// Test search
	queryEmbedding := []float32{0.15, 0.25, 0.35, 0.45, 0.55}
	results, err := vs.SearchSimilar(queryEmbedding, 2)
	if err != nil {
		t.Fatalf("Failed to search: %v", err)
	}

	if len(results) == 0 {
		t.Fatal("Expected search results, got none")
	}

	t.Logf("Found %d results", len(results))
	for i, result := range results {
		t.Logf("Result %d: ID=%s, Score=%.3f, Content=%s", i, result.ID, result.Score, result.Content[:50]+"...")
	}

	// Test collection info
	info, err := vs.GetCollectionInfo()
	if err != nil {
		t.Fatalf("Failed to get collection info: %v", err)
	}

	count, ok := info["count"].(int32)
	if !ok || count != 2 {
		t.Fatalf("Expected count=2, got %v", info["count"])
	}

	// Test deletion
	err = vs.DeleteChunks([]string{"test_1"})
	if err != nil {
		t.Fatalf("Failed to delete chunk: %v", err)
	}

	// Verify deletion
	info, err = vs.GetCollectionInfo()
	if err != nil {
		t.Fatalf("Failed to get collection info after deletion: %v", err)
	}

	count, ok = info["count"].(int32)
	if !ok || count != 1 {
		t.Fatalf("Expected count=1 after deletion, got %v", info["count"])
	}

	t.Log("All vector store operations completed successfully")
}

func TestVectorService(t *testing.T) {
	chromaURL := os.Getenv("CHROMA_URL")
	if chromaURL == "" {
		chromaURL = "http://localhost:8000"
	}

	vs, err := NewVectorService(chromaURL)
	if err != nil {
		t.Skipf("Skipping test: failed to connect to Chroma at %s: %v", chromaURL, err)
	}

	// Test storing document chunks
	documentID := 999 // Use a test document ID
	chunks := []string{
		"This is the first chunk of text.",
		"This is the second chunk of text.",
	}
	embeddings := [][]float32{
		{0.1, 0.2, 0.3},
		{0.4, 0.5, 0.6},
	}
	metadata := map[string]string{
		"filename": "test.pdf",
		"author":   "test_author",
	}

	err = vs.StoreDocumentChunks(documentID, chunks, embeddings, metadata)
	if err != nil {
		t.Fatalf("Failed to store document chunks: %v", err)
	}

	// Clean up after test
	defer func() {
		vs.DeleteDocumentChunks(documentID)
	}()

	// Test searching
	queryEmbedding := []float32{0.2, 0.3, 0.4}
	results, err := vs.SearchRelevantChunks(queryEmbedding, 5)
	if err != nil {
		t.Fatalf("Failed to search chunks: %v", err)
	}

	if len(results) == 0 {
		t.Fatal("Expected search results, got none")
	}

	// Verify results contain our test data
	found := false
	for _, result := range results {
		if result.DocumentID == documentID {
			found = true
			break
		}
	}

	if !found {
		t.Fatal("Expected to find test document in search results")
	}

	// Test stats
	stats, err := vs.GetStats()
	if err != nil {
		t.Fatalf("Failed to get stats: %v", err)
	}

	t.Logf("Collection stats: %+v", stats)

	t.Log("Vector service test completed successfully")
}