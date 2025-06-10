package rag

import (
	"database/sql"
	"fmt"
	// "github.com/mattn/go-sqlite3" // We'll add this when we actually implement
)

const (
	// EmbeddingDimension is the size of the embeddings.
	// For OpenAI text-embedding-ada-002, this is 1536.
	EmbeddingDimension = 1536
)

// VectorStore defines the interface for interacting with a vector database.
type VectorStore interface {
	AddDocuments(docs []IndexedDocument) error
	SearchSimilar(queryEmbedding []float32, topN int, filePathFilter string) ([]IndexedDocument, error)
	GetDocumentByPath(filePath string) (*IndexedDocument, error) // Could return []IndexedDocument if a file is split into many
	RemoveDocumentsByPath(filePath string) error
	IsPathIndexed(filePath string) (bool, error)
	GetIndexedFilePaths() ([]string, error)
}

// SQLiteVectorStore is an implementation of VectorStore using SQLite with a vector extension.
type SQLiteVectorStore struct {
	db *sql.DB
}

// NewSQLiteVectorStore creates a new SQLiteVectorStore.
// dbPath is the path to the SQLite database file.
func NewSQLiteVectorStore(dbPath string) (*SQLiteVectorStore, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	store := &SQLiteVectorStore{db: db}
	if err := store.initialize(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}
	return store, nil
}

// initialize sets up the necessary tables in the SQLite database.
func (s *SQLiteVectorStore) initialize() error {
	// Table for storing document metadata and actual text content
	// We store TextChunk here to be able to retrieve it directly.
	// ContentHash helps in quickly checking if a file needs re-indexing.
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS documents (
			id TEXT PRIMARY KEY,
			file_path TEXT NOT NULL,
			content_hash TEXT NOT NULL,
			text_chunk TEXT NOT NULL,
			indexed_at DATETIME NOT NULL,
			metadata TEXT -- JSON stored as TEXT
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create documents table: %w", err)
	}

	// Create an index on file_path for faster lookups
	_, err = s.db.Exec(`CREATE INDEX IF NOT EXISTS idx_documents_file_path ON documents(file_path);`)
	if err != nil {
		return fmt.Errorf("failed to create index on file_path: %w", err)
	}

	// Virtual table for vector search using an extension like sqlite-vss
	// The embedding dimension (e.g., 1536) should match the model used.
	// The `embedding` column in `vss_docs` will store the vectors.
	// We link it to the `documents` table via rowid.
	// Note: Actual integration of sqlite-vss or similar is pending.
	// This DDL assumes such an extension is available.
	_, err = s.db.Exec(fmt.Sprintf(`
		CREATE VIRTUAL TABLE IF NOT EXISTS vss_docs USING vss0(
			embedding(%d)
		);`, EmbeddingDimension))
	if err != nil {
		// This might fail if the vss0 extension is not loaded.
		// For now, we'll log this but not necessarily fail initialization,
		// as the core Plandex app might function without RAG initially.
		fmt.Printf("Warning: Failed to create vss_docs virtual table (this is expected if VSS extension is not available): %v\n", err)
		// return fmt.Errorf("failed to create vss_docs virtual table: %w", err)
	}

	return nil
}

// AddDocuments adds a batch of documents to the vector store.
// This will involve inserting metadata into 'documents' and vectors into 'vss_docs'.
func (s *SQLiteVectorStore) AddDocuments(docs []IndexedDocument) error {
	// TODO: Implement batch insertion
	// For each doc:
	// 1. Insert into 'documents' table (handle metadata serialization, e.g., to JSON)
	// 2. Insert embedding into 'vss_docs' table, getting the rowid from 'documents' insert.
	return fmt.Errorf("AddDocuments not implemented")
}

// SearchSimilar finds documents with embeddings similar to the queryEmbedding.
// topN specifies the number of similar documents to return.
// filePathFilter can be used to restrict search to a specific file (optional).
func (s *SQLiteVectorStore) SearchSimilar(queryEmbedding []float32, topN int, filePathFilter string) ([]IndexedDocument, error) {
	// TODO: Implement similarity search
	// 1. Construct a query for 'vss_docs' to find nearest neighbors.
	//    - Use `vss_search` or similar function provided by the VSS extension.
	//    - `SELECT rowid, distance FROM vss_docs WHERE vss_search(embedding, ?)`
	// 2. Join with 'documents' table to retrieve metadata and text.
	//    - `SELECT d.id, d.file_path, d.content_hash, d.text_chunk, d.indexed_at, d.metadata FROM documents d JOIN vss_docs v ON d.rowid = v.rowid WHERE vss_search...`
	// 3. If filePathFilter is provided, add a WHERE clause for `d.file_path`.
	// 4. Limit results to topN.
	// 5. Deserialize metadata.
	return nil, fmt.Errorf("SearchSimilar not implemented")
}

// GetDocumentByPath retrieves all document chunks associated with a given file path.
// This is useful for checking if a file is indexed or what parts of it are.
func (s *SQLiteVectorStore) GetDocumentByPath(filePath string) (*IndexedDocument, error) {
	// TODO: Implement retrieval by file path
	// This might return []IndexedDocument if a file can have multiple chunks.
	// For now, let's assume one doc per path for simplicity or the main/first one.
	// `SELECT id, file_path, content_hash, text_chunk, indexed_at, metadata FROM documents WHERE file_path = ?`
	return nil, fmt.Errorf("GetDocumentByPath not implemented")
}

// RemoveDocumentsByPath removes all document chunks associated with a given file path.
func (s *SQLiteVectorStore) RemoveDocumentsByPath(filePath string) error {
	// TODO: Implement removal by file path
	// 1. Find all document rowids associated with the filePath from 'documents' table.
	// 2. Delete entries from 'vss_docs' using these rowids.
	// 3. Delete entries from 'documents' table using filePath.
	// Needs to be transactional.
	return fmt.Errorf("RemoveDocumentsByPath not implemented")
}

// IsPathIndexed checks if a given file path has any indexed documents.
func (s *SQLiteVectorStore) IsPathIndexed(filePath string) (bool, error) {
	// TODO: Implement check for indexed path
	// `SELECT EXISTS(SELECT 1 FROM documents WHERE file_path = ?)`
	return false, fmt.Errorf("IsPathIndexed not implemented")
}

// GetIndexedFilePaths returns a list of all unique file paths that have been indexed.
func (s *SQLiteVectorStore) GetIndexedFilePaths() ([]string, error) {
	// TODO: Implement retrieval of all indexed file paths
	// `SELECT DISTINCT file_path FROM documents`
	return nil, fmt.Errorf("GetIndexedFilePaths not implemented")
}

// Close closes the database connection.
func (s *SQLiteVectorStore) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}
