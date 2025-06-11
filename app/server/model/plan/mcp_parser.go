package plan

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	shared "plandex-shared"

	"github.com/santhosh-tekuri/jsonschema/v5"
	_ "github.com/santhosh-tekuri/jsonschema/v5/httploader" // Required for loading remote schemas, if any
)

// TryParseToolInvocation attempts to parse the given responseText as a Plandex Tool Invocation.
// It returns the parsed request, a boolean indicating if it's a tool call, and any error during parsing/validation.
func TryParseToolInvocation(responseText string, availableTools []shared.MCPToolDefinition) (
	toolRequest *shared.PlandexToolInvocationRequest,
	isToolCall bool,
	validationError error, // Error specifically from validation if it is a tool call but input is wrong
	toolDefinition *shared.MCPToolDefinition, // The matched tool definition if valid
) {
	// First, try to unmarshal into the wrapper structure.
	// The LLM is instructed to respond *only* with the JSON object.
	// We trim whitespace in case of minor deviations.
	trimmedResponse := strings.TrimSpace(responseText)

	var wrapper shared.PlandexToolInvocationWrapper
	err := json.Unmarshal([]byte(trimmedResponse), &wrapper)
	if err != nil {
		// Not a JSON object or not the correct wrapper structure.
		// This is not necessarily an error, just means it's not a tool call.
		return nil, false, nil, nil
	}

	// Check if the expected plandex_tool_invocation key is present and ToolName is non-empty
	if wrapper.Invocation.ToolName == "" {
		// It was JSON, but not a valid PlandexToolInvocationWrapper structure as we define it.
		// This could happen if the LLM produces a different JSON object.
		return nil, false, nil, nil
	}

	// At this point, we consider it an attempted tool call.
	isToolCall = true
	parsedRequest := wrapper.Invocation

	log.Printf("MCP Parser: Detected potential tool invocation for tool: %s", parsedRequest.ToolName)

	// Validate: Find the tool definition
	var matchedToolDef shared.MCPToolDefinition
	foundTool := false
	for _, t := range availableTools {
		if t.ToolName == parsedRequest.ToolName {
			matchedToolDef = t
			foundTool = true
			break
		}
	}

	if !foundTool {
		log.Printf("MCP Parser: LLM requested unknown tool '%s'", parsedRequest.ToolName)
		return &parsedRequest, true, fmt.Errorf("tool '%s' not found in defined tools", parsedRequest.ToolName), nil
	}
	log.Printf("MCP Parser: Found definition for tool '%s'", matchedToolDef.ToolName)
	toolDefinition = &matchedToolDef

	// Validate InputSchema if provided
	if matchedToolDef.InputSchema == "" {
		log.Printf("MCP Parser: Tool '%s' has no InputSchema defined. Assuming any input (or no input) is valid.", matchedToolDef.ToolName)
		// If schema is empty, but input is provided, it might be an error based on stricter interpretation.
		// For now, allow it, as an empty schema might mean "accepts anything" or "accepts nothing".
		// If ToolInput is also empty/nil, it's fine. If ToolInput is provided, it passes.
		return &parsedRequest, true, nil, toolDefinition
	}

	// Compile the schema
	compiler := jsonschema.NewCompiler()
	compiler.Draft = jsonschema.Draft2020_12 // Or your preferred draft
	if err := compiler.AddResource("schema.json", strings.NewReader(matchedToolDef.InputSchema)); err != nil {
		log.Printf("MCP Parser: Error compiling InputSchema for tool '%s': %v", matchedToolDef.ToolName, err)
		// This is a server-side schema error, not an LLM input error.
		return &parsedRequest, true, fmt.Errorf("internal error: could not compile schema for tool '%s': %w", matchedToolDef.ToolName, err), toolDefinition
	}
	schema, err := compiler.Compile("schema.json")
	if err != nil {
		log.Printf("MCP Parser: Error compiling InputSchema for tool '%s': %v", matchedToolDef.ToolName, err)
		return &parsedRequest, true, fmt.Errorf("internal error: could not compile schema for tool '%s': %w", matchedToolDef.ToolName, err), toolDefinition
	}

	// Convert toolRequest.ToolInput (map[string]interface{}) back to JSON bytes for validation
	var inputJSONBytes []byte
	if parsedRequest.ToolInput != nil {
		inputJSONBytes, err = json.Marshal(parsedRequest.ToolInput)
		if err != nil {
			log.Printf("MCP Parser: Error marshalling ToolInput for validation for tool '%s': %v", matchedToolDef.ToolName, err)
			return &parsedRequest, true, fmt.Errorf("error preparing tool input for validation: %w", err), toolDefinition
		}
	} else {
		// If ToolInput is nil, use an empty JSON object for validation against the schema
		inputJSONBytes = []byte("{}")
	}

	log.Printf("MCP Parser: Validating input for tool '%s' against schema. Input: %s", matchedToolDef.ToolName, string(inputJSONBytes))

	// Validate the input
	if err := schema.Validate(bytes.NewReader(inputJSONBytes)); err != nil {
		log.Printf("MCP Parser: LLM input for tool '%s' failed schema validation: %v", matchedToolDef.ToolName, err)
		return &parsedRequest, true, fmt.Errorf("input for tool '%s' is invalid: %w", matchedToolDef.ToolName, err), toolDefinition
	}

	log.Printf("MCP Parser: Input for tool '%s' validated successfully.", matchedToolDef.ToolName)
	return &parsedRequest, true, nil, toolDefinition
}
