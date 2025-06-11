package rag

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sashabaranov/go-openai"
)

// Indexer handles the process of reading files, chunking them, generating embeddings,
// and storing them in a VectorStore.
type Indexer struct {
	store        VectorStore
	config       IndexerConfig
	openAIClient *openai.Client // Client for OpenAI API, must be pre-configured with API key.
}

// NewIndexer creates a new Indexer.
// The openAIClient must be pre-configured with an API key.
func NewIndexer(store VectorStore, config IndexerConfig, openAIClient *openai.Client) *Indexer {
	return &Indexer{
		store:        store,
		config:       config,
		openAIClient: openAIClient,
	}
}

// generateEmbedding uses the OpenAI client to create an embedding for the given text chunk.
func (i *Indexer) generateEmbedding(ctx context.Context, textChunk string) ([]float32, error) {
	if i.openAIClient == nil {
		return nil, fmt.Errorf("OpenAI client is not initialized")
	}

	model := openai.EmbeddingModel(i.config.EmbeddingModelName)
	if model == "" {
		// Default to text-embedding-ada-002 if not specified
		model = openai.AdaEmbeddingV2
	}

	req := openai.EmbeddingRequest{
		Input: []string{textChunk}, // OpenAI API expects a slice of strings
		Model: model,
	}

	resp, err := i.openAIClient.CreateEmbeddings(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error creating embeddings: %w", err)
	}

	if len(resp.Data) == 0 || len(resp.Data[0].Embedding) == 0 {
		return nil, fmt.Errorf("received empty embedding data from API")
	}

	// Assuming we only sent one input string, so we take the first result.
	return resp.Data[0].Embedding, nil
}

// IndexProjectFiles traverses the projectRoot directory, processes allowed files,
// and updates the vector store.
func (i *Indexer) IndexProjectFiles(projectRoot string) error {
	log.Printf("Starting indexing process for project root: %s", projectRoot)

	var newDocuments []IndexedDocument
	processedFilePaths := make(map[string]bool)

	err := filepath.WalkDir(projectRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Printf("Error accessing path %q: %v\n", path, err)
			return err
		}

		// Skip directories, hidden files/dirs, and non-regular files
		if d.IsDir() || !d.Type().IsRegular() {
			if strings.HasPrefix(d.Name(), ".") && d.Name() != "." && d.Name() != ".." {
				log.Printf("Skipping hidden path: %s", path)
				if d.IsDir() {
					return filepath.SkipDir // Skip walking the contents of hidden directories
				}
				return nil
			}
			if d.IsDir() { // Regular directory, continue walking
				return nil
			}
			log.Printf("Skipping non-regular file: %s", path)
			return nil
		}

		// Check file extension
		allowed := false
		for _, ext := range i.config.AllowedFileExtensions {
			if strings.HasSuffix(strings.ToLower(path), strings.ToLower(ext)) {
				allowed = true
				break
			}
		}
		if !allowed {
			log.Printf("Skipping file with unallowed extension: %s", path)
			return nil
		}

		// Mark this path as processed in this run
		processedFilePaths[path] = true

		// Read file content
		content, err := os.ReadFile(path)
		if err != nil {
			log.Printf("Error reading file %q: %v\n", path, err)
			return nil // Continue with other files
		}

		// Calculate content hash
		hashBytes := sha256.Sum256(content)
		contentHash := hex.EncodeToString(hashBytes[:])

		// File Change Check
		// We use GetDocumentByPath to fetch the *first* chunk of an existing document.
		// This is a simplification; a more robust check might involve a dedicated metadata store
		// or querying for all chunks and verifying hashes.
		existingDoc, err := i.store.GetDocumentByPath(path)
		if err == nil && existingDoc != nil { // Document exists
			if existingDoc.ContentHash == contentHash {
				log.Printf("File %s unchanged (hash: %s), skipping.", path, contentHash)
				return nil
			}
			log.Printf("File %s changed (old hash: %s, new hash: %s). Re-indexing.", path, existingDoc.ContentHash, contentHash)
			if err := i.store.RemoveDocumentsByPath(path); err != nil {
				log.Printf("Error removing old version of %s from store: %v", path, err)
				// Continue to try and index the new version
			}
		} else {
			// Could be an error other than "not found", or simply not indexed yet.
			// If err is not nil, it might be an actual DB error. For now, we proceed.
			log.Printf("File %s is new or not found in store. Indexing.", path)
		}

		// Content Chunking
		var chunks []string
		charLimit := 0
		if i.config.MaxChunkSizeTokens > 0 {
			// Approximate character limit: average 3.5 chars/token. This is a rough heuristic.
			// Using a proper tokenizer (e.g., tiktoken) would be more accurate but adds complexity/dependencies.
			// TODO: Replace with tokenizer-based chunking for better accuracy.
			charLimit = int(float64(i.config.MaxChunkSizeTokens) * 3.5)
			log.Printf("RAG Indexer: Using character-based chunking for %s. Target chars/chunk: %d (from MaxChunkSizeTokens: %d)", path, charLimit, i.config.MaxChunkSizeTokens)

			lines := strings.Split(string(content), "\n")
			var currentChunk strings.Builder
			for _, line := range lines {
				// If a single line exceeds the charLimit, it will be a chunk on its own (or could be hard-split).
				// For now, a long line becomes its own chunk.
				if currentChunk.Len() > 0 && currentChunk.Len()+len(line)+1 > charLimit && currentChunk.Len() > charLimit/2 { // +1 for newline, ensure chunk is reasonably large
					chunks = append(chunks, currentChunk.String())
					currentChunk.Reset()
				}
				if currentChunk.Len() > 0 {
					currentChunk.WriteString("\n")
				}
				currentChunk.WriteString(line)

				// If a line itself is very long, make it its own chunk (or split it further - current: own chunk)
				if currentChunk.Len() > charLimit {
					chunks = append(chunks, currentChunk.String())
					currentChunk.Reset()
				}
			}
			if currentChunk.Len() > 0 {
				chunks = append(chunks, currentChunk.String())
			}
		} else {
			log.Printf("RAG Indexer: MaxChunkSizeTokens not configured or is 0 for %s. Falling back to paragraph splitting.", path)
			// Fallback to basic newline splitting (paragraphs)
			chunks = strings.Split(string(content), "\n\n")
			if len(chunks) == 0 && len(content) > 0 {
				chunks = []string{string(content)}
			}
		}

		if len(chunks) == 0 && len(content) > 0 { // Ensure at least one chunk if content exists
			log.Printf("RAG Indexer: Warning - No chunks created for %s despite content. Adding entire content as one chunk.", path)
			chunks = []string{string(content)}
		}

		log.Printf("RAG Indexer: File %s split into %d chunks.", path, len(chunks))

		for chunkIndex, chunkText := range chunks {
			trimmedChunkText := strings.TrimSpace(chunkText)
			if trimmedChunkText == "" {
				continue
			}

			// Embedding Generation
			embedding, err := i.generateEmbedding(context.Background(), trimmedChunkText)
			if err != nil {
				log.Printf("Error generating embedding for chunk %d of %s: %v. Skipping chunk.", chunkIndex+1, path, err)
				continue
			}
			if len(embedding) != EmbeddingDimension {
				log.Printf("Warning: Embedding dimension mismatch for chunk %d of %s. Expected %d, got %d. Skipping chunk.", chunkIndex+1, path, EmbeddingDimension, len(embedding))
				continue
			}

			docID := uuid.NewString()
			indexedDoc := IndexedDocument{
				ID:          docID,
				FilePath:    path,
				ContentHash: contentHash, // File-level hash for all chunks of this file
				TextChunk:   trimmedChunkText, // Use trimmed chunk
				Embedding:   embedding,
				IndexedAt:   time.Now().UTC(),
				Metadata: map[string]interface{}{
					"chunkNumber": chunkIndex + 1,
					"source":      "file_system_indexer",
				},
			}
			newDocuments = append(newDocuments, indexedDoc)
		}
		log.Printf("Processed and chunked file: %s. %d chunks created.", path, len(chunks))
		return nil
	})

	if err != nil {
		log.Printf("Error during file walk: %v", err)
		// Decide if we should proceed with adding documents found so far or return
		// For now, we proceed to add what we have and then handle orphans.
	}

	// Add all new/updated documents to the store
	if len(newDocuments) > 0 {
		log.Printf("Adding %d new/updated document chunks to the vector store...", len(newDocuments))
		if err := i.store.AddDocuments(newDocuments); err != nil {
			log.Printf("Error adding documents to store: %v", err)
			// Potentially return this error, but we'll try to do orphan removal first
		} else {
			log.Printf("%d document chunks added successfully.", len(newDocuments))
		}
	} else {
		log.Println("No new or updated documents to add to the store.")
	}

	// Orphaned Document Removal
	log.Println("Starting orphaned document removal process...")
	indexedFilePaths, err := i.store.GetIndexedFilePaths()
	if err != nil {
		log.Printf("Error getting all indexed file paths from store: %v. Skipping orphan removal.", err)
		return fmt.Errorf("error during indexing (getting indexed paths): %w", err)
	}

	orphansRemoved := 0
	for _, indexedPath := range indexedFilePaths {
		if !processedFilePaths[indexedPath] {
			log.Printf("File %s found in store but not in current project scan. Removing as orphan.", indexedPath)
			if err := i.store.RemoveDocumentsByPath(indexedPath); err != nil {
				log.Printf("Error removing orphaned path %s from store: %v", indexedPath, err)
				// Continue trying to remove other orphans
			} else {
				orphansRemoved++
			}
		}
	}
	log.Printf("Orphaned document removal complete. %d orphaned file paths removed.", orphansRemoved)

	log.Printf("Indexing process completed for project root: %s", projectRoot)
	if err != nil { // Return the original walk error if it occurred
		return fmt.Errorf("error during file walk: %w", err)
	}
	return nil // If walk was fine, but other errors might have been logged
}
