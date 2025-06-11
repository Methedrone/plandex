package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/PlandexAI/plandex/app/server/rag" // Assuming this is the correct path to your rag package
	"github.com/gorilla/mux"
	"github.com/sashabaranov/go-openai"
)

// TODO: Implement proper project path resolution from DB or config
// For now, using a relative path for local testing.
// This assumes the server is run from a directory where ./projects/{projectID} is valid.
func getProjectRootPath(projectID string) (string, error) {
	// For local dev, Plandex projects are typically cloned into a 'projects' subdirectory
	// relative to where the server is running.
	// Example: if server runs in /plandex-app/, project "xyz" is in /plandex-app/projects/xyz/
	// This needs to be robust for production.
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error getting current working directory: %w", err)
	}
	log.Printf("Assuming project root relative to current working directory: %s", wd)
	return filepath.Join(wd, "projects", projectID), nil
}

// HandleIndexProject handles the HTTP request to index a project.
func HandleIndexProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID, ok := vars["projectId"]
	if !ok {
		http.Error(w, "Project ID is missing in URL", http.StatusBadRequest)
		return
	}

	log.Printf("Received indexing request for project ID: %s", projectID)

	projectRootPath, err := getProjectRootPath(projectID)
	if err != nil {
		log.Printf("Error getting project root path for %s: %v", projectID, err)
		http.Error(w, fmt.Sprintf("Error resolving project path: %v", err), http.StatusInternalServerError)
		return
	}
	log.Printf("Project root path resolved to: %s", projectRootPath)

	// Ensure project root exists
	if _, err := os.Stat(projectRootPath); os.IsNotExist(err) {
		log.Printf("Project root path %s does not exist.", projectRootPath)
		http.Error(w, fmt.Sprintf("Project directory not found at: %s", projectRootPath), http.StatusNotFound)
		return
	}

	// Database path
	plandexDir := filepath.Join(projectRootPath, ".plandex")
	if err := os.MkdirAll(plandexDir, 0755); err != nil {
		log.Printf("Error creating .plandex directory for %s: %v", projectID, err)
		http.Error(w, fmt.Sprintf("Error creating .plandex directory: %v", err), http.StatusInternalServerError)
		return
	}
	dbPath := filepath.Join(plandexDir, "rag.db")
	log.Printf("RAG database path for project %s: %s", projectID, dbPath)

	// Initialize VectorStore
	store, err := rag.NewSQLiteVectorStore(dbPath)
	if err != nil {
		log.Printf("Error initializing SQLite vector store for %s (path: %s): %v", projectID, dbPath, err)
		http.Error(w, fmt.Sprintf("Error initializing vector store: %v", err), http.StatusInternalServerError)
		return
	}
	defer store.Close()
	log.Printf("SQLite vector store initialized for project %s", projectID)

	// OpenAI Client
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Println("OPENAI_API_KEY environment variable not set.")
		http.Error(w, "OpenAI API key not configured on server", http.StatusInternalServerError)
		return
	}
	openAIClient := openai.NewClient(apiKey)
	log.Printf("OpenAI client initialized for project %s", projectID)

	// Indexer Config
	// Indexer Config
	// Fetch RAG settings from plan configuration (or use defaults)
	planConfig, err := db.GetPlanSettings(plan, true) // plan is from client_helper.GetPlanFromRequest
	if err != nil {
		log.Printf("Warning: Could not fetch plan settings for project %s to configure RAG indexer: %v. Using default indexer settings.", projectID, err)
		planConfig = &shared.PlanConfig{} // Use empty config, defaults will apply
	}

	var ragChunkSizeTokens int = 512 // Default chunk size
	if planConfig.RAGSettings != nil && planConfig.RAGSettings.ChunkSizeTokens > 0 {
		ragChunkSizeTokens = planConfig.RAGSettings.ChunkSizeTokens
		log.Printf("RAG: Using ChunkSizeTokens from plan config: %d", ragChunkSizeTokens)
	} else {
		log.Printf("RAG: Using default ChunkSizeTokens: %d", ragChunkSizeTokens)
	}

	// TODO: Make AllowedFileExtensions configurable as well if needed
	config := rag.IndexerConfig{
		AllowedFileExtensions: []string{".go", ".md", ".txt", ".py", ".js", ".ts", ".java", ".c", ".cpp", ".h", ".hpp", ".cs", ".rb", ".php", ".html", ".css", ".json", ".yaml", ".yml"},
		MaxChunkSizeTokens:    ragChunkSizeTokens,
		EmbeddingModelName:    string(openai.AdaEmbeddingV2), // TODO: Make EmbeddingModelName configurable from RAGSettings
	}
	log.Printf("Indexer configuration set for project %s: %+v", projectID, config)

	// Create Indexer
	// The NewIndexer function signature was: NewIndexer(store VectorStore, config IndexerConfig, openAIClient *openai.Client)
	indexer := rag.NewIndexer(store, config, openAIClient)
	log.Printf("Indexer created for project %s", projectID)

	// Start indexing in a goroutine so the HTTP request returns immediately.
	// The client will be notified that indexing has started.
	// Actual status updates would require a more complex system (e.g., websockets, polling).
	go func() {
		log.Printf("Starting background indexing for project %s at path %s", projectID, projectRootPath)
		if err := indexer.IndexProjectFiles(projectRootPath); err != nil {
			// Log critical error during indexing.
			// This error won't be directly sent to the client as the HTTP response is already sent.
			// Requires proper operational monitoring.
			log.Printf("CRITICAL: Error during background indexing of project %s: %v", projectID, err)
		} else {
			log.Printf("Background indexing completed successfully for project %s", projectID)
		}
	}()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted) // 202 Accepted: request received, processing started
	json.NewEncoder(w).Encode(map[string]string{"message": "Project indexing initiated. Check server logs for progress."})
}
