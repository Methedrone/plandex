package cmd

import (
	"fmt"
	"log"

	"github.com/PlandexAI/plandex/app/cli/pkg/api"
	"github.com/PlandexAI/plandex/app/cli/pkg/auth"
	"github.com/PlandexAI/plandex/app/cli/pkg/lib"
	"github.com/spf13/cobra"
)

// indexCmd represents the index command
var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "Manages project indexes for Retrieval Augmented Generation (RAG).",
	Long: `The index command is used to create, update, or manage search indexes
for your project. These indexes are utilized by the RAG system to provide
contextually relevant information to the AI models.

Currently, this command initiates the indexing process for the active project.`,
	Run: func(cmd *cobra.Command, args []string) {
		auth.MustResolveAuthWithOrg()
		lib.MustResolveProject()

		projectID := lib.CurrentProject.Id
		if projectID == "" {
			log.Fatal("No active project selected or project ID is missing. Use 'plandex project switch' or ensure your project is configured.")
		}

		fmt.Printf("Initiating indexing for project ID: %s (%s)...\n", lib.CurrentProject.Name, projectID)

		endpoint := fmt.Sprintf("/api/v1/project/%s/index", projectID)
		var responseBody interface{} // Or a more specific struct if the response is structured

		spinner := lib.DisplaySpinner("Requesting indexing from server...")
		err := api.Client.Post(endpoint, nil, &responseBody)
		spinner.Stop()
		if err != nil {
			lib.DisplayError(fmt.Errorf("error initiating indexing for project %s: %v", lib.CurrentProject.Name, err))
			if apiErr, ok := err.(*api.APIError); ok {
				lib.DisplayError(fmt.Errorf("server responded with status %d: %s", apiErr.StatusCode, apiErr.Body))
			}
			return
		}

		// Assuming the responseBody might contain a message, e.g., {"message": "Indexing started"}
		// For now, we'll print a generic success message if no error occurred.
		// Type assertion or more robust parsing can be done if responseBody structure is known.
		if respMap, ok := responseBody.(map[string]interface{}); ok {
			if msg, ok := respMap["message"].(string); ok {
				lib.DisplaySuccess(fmt.Sprintf("Server response: %s", msg))
			} else {
				lib.DisplaySuccess(fmt.Sprintf("Indexing request for project %s accepted by the server.", lib.CurrentProject.Name))
			}
		} else {
			lib.DisplaySuccess(fmt.Sprintf("Indexing request for project %s accepted by the server.", lib.CurrentProject.Name))
		}
		fmt.Println("You can check the server logs for detailed indexing progress.")
	},
}

func init() {
	// This is where rootCmd is typically available if this file is in the same package as root.go
	// This is where rootCmd is typically available if this file is in the same package as root.go
	// If rootCmd is defined in main.go or another package, you might need a different way to add this command,
	// e.g., by having an exported function in this package that root.go calls.
	// rootCmd is accessible here as root.go is in the same directory.
	rootCmd.AddCommand(indexCmd)
}
