package shared

// MCPExecutionType defines the type of execution for an MCP tool.
type MCPExecutionType string

const (
	// MCPExecutionTypeHTTP indicates that the tool is executed via an HTTP request.
	MCPExecutionTypeHTTP MCPExecutionType = "http"
	// MCPExecutionTypePredefined indicates that the tool is a predefined function on the server.
	MCPExecutionTypePredefined MCPExecutionType = "predefined_server_function"
)

// MCPToolDefinition defines the structure for a tool that can be invoked by the Model-Calling-Plandex (MCP) system.
type MCPToolDefinition struct {
	// ToolName is a unique name for the tool.
	ToolName string `json:"toolName"`
	// Description is a detailed description for the LLM to understand the tool's purpose and when to use it.
	Description string `json:"description"`
	// InputSchema is a JSON Schema string defining the expected input parameters for the tool.
	InputSchema string `json:"inputSchema"`
	// OutputSchema is a JSON Schema string defining the expected output structure from the tool. (For validation/documentation)
	OutputSchema string `json:"outputSchema"`
	// ExecutionType specifies how the tool is executed (e.g., 'http', 'predefined_server_function').
	ExecutionType MCPExecutionType `json:"executionType"`
	// ExecutionDetails is a map containing specific details for the execution, structure depends on ExecutionType.
	// For MCPExecutionTypeHTTP, details might include:
	//   - "url": string (The endpoint URL)
	//   - "method": string (e.g., "GET", "POST")
	//   - "parameterMapping": map[string]string (e.g., maps inputSchema properties to URL query params or request body)
	// For MCPExecutionTypePredefined, details might include:
	//   - "functionName": string (The name of the server-side function to call)
	ExecutionDetails map[string]interface{} `json:"executionDetails"`
}

// MCPConfig holds the configuration for Model-Calling-Plandex, primarily a list of available tools.
type MCPConfig struct {
	Tools []MCPToolDefinition `json:"tools,omitempty"`
}
