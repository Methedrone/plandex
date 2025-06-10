package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"plandex-cli/api"
	"plandex-cli/auth"
	"plandex-cli/lib"
	"plandex-cli/term"

	shared "plandex-shared"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Manage Model Context Protocol (MCP) tools for the current plan",
	Long: `The mcp command group allows you to list, add, or remove
Model Context Protocol (MCP) tools associated with the current Plandex plan.
MCP tools extend the capabilities of the AI model by allowing it to interact
with external services or predefined server functions.`,
}

var listToolsCmd = &cobra.Command{
	Use:   "list-tools",
	Short: "List MCP tools for the current plan",
	Long:  "Retrieves and displays all MCP tools currently registered with the active Plandex plan.",
	Run:   listMCPTools,
}

var addToolCmd = &cobra.Command{
	Use:   "add-tool",
	Short: "Add a new MCP tool to the current plan from a JSON file",
	Long: `Adds a new Model Context Protocol (MCP) tool to the current Plandex plan.
The tool definition must be provided as a JSON file.`,
	Run: addMCPTool,
}

var removeToolCmd = &cobra.Command{
	Use:   "remove-tool [toolName]",
	Short: "Remove an MCP tool from the current plan",
	Long:  "Removes a specific Model Context Protocol (MCP) tool, identified by its name, from the current Plandex plan.",
	Args:  cobra.ExactArgs(1), // Requires exactly one argument: toolName
	Run:   removeMCPTool,
}

func listMCPTools(cmd *cobra.Command, args []string) {
	auth.MustResolveAuthWithOrg()
	lib.MustResolveProject()

	if lib.CurrentPlanId == "" {
		term.OutputErrorAndExit("No active plan selected. Use 'plandex checkout' to select a plan.")
		return
	}

	term.StartSpinner(fmt.Sprintf("Fetching MCP tools for plan %s...", lib.CurrentPlanId))
	path := fmt.Sprintf("/api/v1/plan/%s/mcp/tools", lib.CurrentPlanId)
	var toolsList []shared.MCPToolDefinition
	apiErr := api.Client.Get(path, &toolsList)
	term.StopSpinner()

	if apiErr != nil {
		term.OutputErrorAndExit("Error fetching MCP tools: %v", apiErr)
		return
	}

	if len(toolsList) == 0 {
		fmt.Println("No MCP tools are configured for the current plan.")
		return
	}

	fmt.Printf("MCP Tools for plan %s (%s):\n", lib.CurrentPlanName, lib.CurrentPlanId)
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Tool Name", "Description", "Execution Type"})
	table.SetAutoWrapText(true)
	table.SetColWidth(60) // For description

	for _, tool := range toolsList {
		table.Append([]string{tool.ToolName, tool.Description, string(tool.ExecutionType)})
	}
	table.Render()
}

// removeToolCmd will be defined here

func addMCPTool(cmd *cobra.Command, args []string) {
	auth.MustResolveAuthWithOrg()
	lib.MustResolveProject()

	if lib.CurrentPlanId == "" {
		term.OutputErrorAndExit("No active plan selected. Use 'plandex checkout' to select a plan.")
		return
	}

	filePath, _ := cmd.Flags().GetString("file")
	if filePath == "" { // Should be caught by MarkFlagRequired, but good practice to check
		term.OutputErrorAndExit("File path for tool definition is required. Use --file <path>.")
		return
	}

	// Read the file content
	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		term.OutputErrorAndExit("Error resolving file path '%s': %v", filePath, err)
		return
	}

	term.StartSpinner(fmt.Sprintf("Reading tool definition from %s...", absFilePath))
	fileContent, err := os.ReadFile(absFilePath)
	if err != nil {
		term.StopSpinner()
		term.OutputErrorAndExit("Error reading tool definition file '%s': %v", absFilePath, err)
		return
	}
	term.StopSpinner()

	// Unmarshal the JSON content
	var toolDefinition shared.MCPToolDefinition
	if err := json.Unmarshal(fileContent, &toolDefinition); err != nil {
		term.OutputErrorAndExit("Error parsing JSON tool definition from '%s': %v\n\nFile content should be a valid JSON object representing a single MCPToolDefinition.", absFilePath, err)
		return
	}

	if toolDefinition.ToolName == "" {
		term.OutputErrorAndExit("ToolName is a required field in the JSON tool definition.")
		return
	}
	// Add more validation as necessary based on MCPToolDefinition fields

	term.StartSpinner(fmt.Sprintf("Adding MCP tool '%s' to plan %s...", toolDefinition.ToolName, lib.CurrentPlanId))
	apiPath := fmt.Sprintf("/api/v1/plan/%s/mcp/tools", lib.CurrentPlanId)
	var response shared.MCPToolDefinition // Expecting the added tool back
	apiErr := api.Client.Post(apiPath, toolDefinition, &response)
	term.StopSpinner()

	if apiErr != nil {
		term.OutputErrorAndExit("Error adding MCP tool '%s': %v", toolDefinition.ToolName, apiErr)
		return
	}

	fmt.Printf("✅ MCP tool '%s' added successfully to plan %s.\n", response.ToolName, lib.CurrentPlanId)
	// Optionally display details of the added tool here
}

func removeMCPTool(cmd *cobra.Command, args []string) {
	auth.MustResolveAuthWithOrg()
	lib.MustResolveProject()

	if lib.CurrentPlanId == "" {
		term.OutputErrorAndExit("No active plan selected. Use 'plandex checkout' to select a plan.")
		return
	}

	toolName := args[0]
	// Basic validation for toolName, e.g., non-empty, though cobra.ExactArgs(1) ensures it's present.
	if toolName == "" {
		term.OutputErrorAndExit("Tool name cannot be empty.") // Should not be reached due to Args check
		return
	}

	// Confirmation prompt
	confirm, err := term.GetConfirmation(fmt.Sprintf("Are you sure you want to remove MCP tool '%s' from plan %s?", toolName, lib.CurrentPlanName), false)
	if err != nil {
		term.OutputErrorAndExit("Error getting confirmation: %v", err)
		return
	}
	if !confirm {
		fmt.Println("Tool removal cancelled.")
		return
	}

	term.StartSpinner(fmt.Sprintf("Removing MCP tool '%s' from plan %s...", toolName, lib.CurrentPlanId))
	// URL encoding for toolName might be needed if names can contain special characters.
	// For now, assuming simple names or that the API client/server handles it.
	apiPath := fmt.Sprintf("/api/v1/plan/%s/mcp/tools/%s", lib.CurrentPlanId, toolName)
	var response interface{} // DELETE typically doesn't return a body, or just a status message
	apiErr := api.Client.Delete(apiPath, nil, &response)
	term.StopSpinner()

	if apiErr != nil {
		term.OutputErrorAndExit("Error removing MCP tool '%s': %v", toolName, apiErr)
		return
	}

	fmt.Printf("✅ MCP tool '%s' removed successfully from plan %s.\n", toolName, lib.CurrentPlanId)
}

func init() {
	RootCmd.AddCommand(mcpCmd)
	mcpCmd.AddCommand(listToolsCmd)
	mcpCmd.AddCommand(addToolCmd)
	addToolCmd.Flags().StringP("file", "f", "", "Path to the JSON file containing the MCP tool definition")
	addToolCmd.MarkFlagRequired("file")
	mcpCmd.AddCommand(removeToolCmd)
}
