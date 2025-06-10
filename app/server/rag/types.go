package rag

import "time"

// IndexerConfig holds configuration for the Indexer service.
type IndexerConfig struct {
	AllowedFileExtensions []string `json:"allowedFileExtensions"` // e.g., [".go", ".md", ".txt"]
	MaxChunkSizeTokens    int      `json:"maxChunkSizeTokens"`    // Max tokens per chunk for LLM processing
	EmbeddingModelName    string   `json:"embeddingModelName"`    // Name of the embedding model used
}

// IndexedDocument represents a document that has been processed and stored.
type IndexedDocument struct {
	ID          string    `json:"id"`          // Unique identifier for the document or chunk
	FilePath    string    `json:"filePath"`    // Path to the source file
	ContentHash string    `json:"contentHash"` // Hash of the TextChunk to detect changes
	TextChunk   string    `json:"textChunk"`   // The actual text content chunk
	Embedding   []float32 `json:"embedding"`   // The vector embedding of the TextChunk
	IndexedAt   time.Time `json:"indexedAt"`   // Timestamp of when the document was indexed
	Metadata    map[string]interface{} `json:"metadata"`    // Any other relevant info (e.g., line numbers, symbol type)
}
