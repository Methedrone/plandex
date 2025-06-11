package plan

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	shared "plandex-shared"
)

// PredefinedFunctionMap maps function names to their implementations.
var PredefinedFunctionMap = map[string]func(inputs map[string]interface{}) (string, error){
	"echoTool":         echoTool,
	"simpleCalculator": simpleCalculator,
	// Add other predefined tools here
}

// ExecutePendingToolCall executes the tool call stored in state.pendingToolCall.
func ExecutePendingToolCall(ctx context.Context, state *activeTellStreamState) (toolOutput string, err error) {
	if state.pendingToolCall == nil || state.pendingToolDefinition == nil {
		return "", nil // No pending tool call
	}

	toolName := state.pendingToolCall.ToolName
	toolDef := state.pendingToolDefinition
	toolInput := state.pendingToolCall.ToolInput

	log.Printf("MCP Executor: Executing tool '%s' with type '%s'", toolName, toolDef.ExecutionType)
	log.Printf("MCP Executor: Tool input: %+v", toolInput)

	switch toolDef.ExecutionType {
	case shared.MCPExecutionTypeHTTP:
		return executeHTTPTool(ctx, toolDef, toolInput)
	case shared.MCPExecutionTypePredefined:
		return executePredefinedTool(ctx, toolDef, toolInput)
	default:
		return "", fmt.Errorf("unknown MCP execution type: %s for tool %s", toolDef.ExecutionType, toolName)
	}
}

func executeHTTPTool(ctx context.Context, toolDef *shared.MCPToolDefinition, inputs map[string]interface{}) (string, error) {
	details := toolDef.ExecutionDetails
	rawURL, okURL := details["url"].(string)
	method, okMethod := details["method"].(string)

	if !okURL || rawURL == "" {
		return "", fmt.Errorf("HTTP tool '%s' missing or invalid 'url' in ExecutionDetails", toolDef.ToolName)
	}
	if !okMethod || method == "" {
		return "", fmt.Errorf("HTTP tool '%s' missing or invalid 'method' in ExecutionDetails", toolDef.ToolName)
	}
	method = strings.ToUpper(method)

	// Basic URL templating: replace {param} with value from inputs
	finalURL, err := performURLTemplating(rawURL, inputs)
	if err != nil {
		return "", fmt.Errorf("error processing URL template for tool '%s': %w", toolDef.ToolName, err)
	}

	var reqBody io.Reader
	if method == "POST" || method == "PUT" || method == "PATCH" {
		if inputs != nil {
			// Send all inputs as JSON body, after filtering out those used in URL templating.
			// A more sophisticated approach might use parameterMapping from ExecutionDetails.
			bodyInputs := make(map[string]interface{})
			for k, v := range inputs {
				if !strings.Contains(rawURL, "{"+k+"}") { // Only include if not used in URL path
					bodyInputs[k] = v
				}
			}
			if len(bodyInputs) > 0 {
				jsonBody, err := json.Marshal(bodyInputs)
				if err != nil {
					return "", fmt.Errorf("error marshalling JSON body for tool '%s': %w", toolDef.ToolName, err)
				}
				reqBody = bytes.NewBuffer(jsonBody)
				log.Printf("MCP Executor: HTTP %s to %s with body: %s", method, finalURL, string(jsonBody))
			} else {
				log.Printf("MCP Executor: HTTP %s to %s with no body (inputs used in URL or empty).", method, finalURL)
			}
		}
	} else if (method == "GET" || method == "DELETE") && len(inputs) > 0 {
		// For GET/DELETE, if there are inputs not used in URL templating, consider adding as query params.
		// This is a simplified version; a full implementation would use parameterMapping.
		queryParams := url.Values{}
		anyParamsAdded := false
		for k, v := range inputs {
			if !strings.Contains(rawURL, "{"+k+"}") { // Only include if not used in URL path
				// Convert v to string. This is a basic conversion.
				queryParams.Add(k, fmt.Sprint(v))
				anyParamsAdded = true
			}
		}
		if anyParamsAdded {
			if strings.Contains(finalURL, "?") {
				finalURL += "&" + queryParams.Encode()
			} else {
				finalURL += "?" + queryParams.Encode()
			}
			log.Printf("MCP Executor: HTTP %s to %s with query parameters from input.", method, finalURL)
		} else {
			log.Printf("MCP Executor: HTTP %s to %s. Inputs used in URL or no remaining inputs for query.", method, finalURL)
		}
	} else {
		log.Printf("MCP Executor: HTTP %s to %s", method, finalURL)
	}


	req, err := http.NewRequestWithContext(ctx, method, finalURL, reqBody)
	if err != nil {
		return "", fmt.Errorf("error creating HTTP request for tool '%s': %w", toolDef.ToolName, err)
	}
	if reqBody != nil && (method == "POST" || method == "PUT" || method == "PATCH") {
		req.Header.Set("Content-Type", "application/json")
	}
	// TODO: Add other headers from ExecutionDetails if specified (e.g., Authorization)

	client := &http.Client{Timeout: 15 * time.Second} // Configurable timeout
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error executing HTTP request for tool '%s': %w", toolDef.ToolName, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Try to read body for error details, but limit size
		errorBodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 1*1024*1024)) // Limit to 1MB
		return "", fmt.Errorf("HTTP tool '%s' request failed with status %s: %s", toolDef.ToolName, resp.Status, string(errorBodyBytes))
	}

	// Limit response body size
	responseBodyBytes, err := io.ReadAll(io.LimitReader(resp.Body, 1*1024*1024)) // Limit to 1MB
	if err != nil {
		return "", fmt.Errorf("error reading HTTP response body for tool '%s': %w", toolDef.ToolName, err)
	}

	log.Printf("MCP Executor: HTTP tool '%s' executed successfully. Response status: %s", toolDef.ToolName, resp.Status)
	return string(responseBodyBytes), nil
}

// performURLTemplating replaces placeholders like {key} in the rawURL with values from inputs.
// It also removes these keys from the inputs map so they are not reused in body/query params.
func performURLTemplating(rawURL string, inputs map[string]interface{}) (string, error) {
    finalURL := rawURL
    // Iterate over a copy of keys if modifying the map during iteration, or collect keys first.
    // For simplicity, we assume that if a key is used in the path, it's okay if it's also in the inputs map
    // that might be used for a JSON body (e.g. POST). A stricter version would remove used keys.
    for key, val := range inputs {
        placeholder := "{" + key + "}"
        strVal, ok := val.(string) // Basic templating assumes string values for path params
        if !ok {
            // If not a string, could attempt to convert or skip. For path params, strings are typical.
            // For this basic version, we'll skip non-string replaceable params to avoid format issues.
            if strings.Contains(finalURL, placeholder) {
                 log.Printf("MCP Executor: URL template placeholder %s found, but input value for key '%s' is not a string. Skipping replacement for this placeholder.", placeholder, key)
            }
            continue
        }
        if strings.Contains(finalURL, placeholder) {
            finalURL = strings.ReplaceAll(finalURL, placeholder, strVal)
        }
    }
    // Check if any placeholders remain
    if strings.Contains(finalURL, "{") && strings.Contains(finalURL, "}") {
        // This is a basic check. A more robust check would use regex to find unreplaced {key} patterns.
        log.Printf("MCP Executor: Warning - URL '%s' may still contain unresolved template parameters after substitution.", finalURL)
    }
    return finalURL, nil
}


func executePredefinedTool(ctx context.Context, toolDef *shared.MCPToolDefinition, inputs map[string]interface{}) (string, error) {
	functionName, ok := toolDef.ExecutionDetails["functionName"].(string)
	if !ok || functionName == "" {
		return "", fmt.Errorf("predefined tool '%s' missing or invalid 'functionName' in ExecutionDetails", toolDef.ToolName)
	}

	fn, exists := PredefinedFunctionMap[functionName]
	if !exists {
		return "", fmt.Errorf("predefined function '%s' for tool '%s' not found in server map", functionName, toolDef.ToolName)
	}

	log.Printf("MCP Executor: Executing predefined function '%s' for tool '%s'", functionName, toolDef.ToolName)
	output, err := fn(inputs)
	if err != nil {
		return "", fmt.Errorf("error executing predefined function '%s' for tool '%s': %w", functionName, toolDef.ToolName, err)
	}
	log.Printf("MCP Executor: Predefined function '%s' executed successfully for tool '%s'", functionName, toolDef.ToolName)
	return output, nil
}

// --- Predefined Tool Implementations ---

func echoTool(inputs map[string]interface{}) (string, error) {
	log.Printf("MCP Predefined: echoTool called with inputs: %+v", inputs)
	outputBytes, err := json.Marshal(inputs)
	if err != nil {
		return "", fmt.Errorf("echoTool: error marshalling inputs to JSON: %w", err)
	}
	return string(outputBytes), nil
}

func simpleCalculator(inputs map[string]interface{}) (string, error) {
	log.Printf("MCP Predefined: simpleCalculator called with inputs: %+v", inputs)
	rawA, okA := inputs["a"]
	rawB, okB := inputs["b"]
	operation, okOp := inputs["operation"].(string)

	if !okA || !okB || !okOp {
		return "", fmt.Errorf("simpleCalculator: missing required inputs 'a', 'b', or 'operation'")
	}

	var a, b float64
	var err error

	// Convert a and b to float64, accommodating string or number types from JSON
	switch valA := rawA.(type) {
	case float64:
		a = valA
	case string:
		a, err = strconv.ParseFloat(valA, 64)
		if err != nil {
			return "", fmt.Errorf("simpleCalculator: invalid number format for 'a': %v (input: %s)", err, valA)
		}
	default:
		return "", fmt.Errorf("simpleCalculator: 'a' must be a number or a string representing a number, got %T", rawA)
	}

	switch valB := rawB.(type) {
	case float64:
		b = valB
	case string:
		b, err = strconv.ParseFloat(valB, 64)
		if err != nil {
			return "", fmt.Errorf("simpleCalculator: invalid number format for 'b': %v (input: %s)", err, valB)
		}
	default:
		return "", fmt.Errorf("simpleCalculator: 'b' must be a number or a string representing a number, got %T", rawB)
	}

	var result float64
	switch strings.ToLower(operation) {
	case "add":
		result = a + b
	case "subtract":
		result = a - b
	// TODO: Add multiply, divide with check for division by zero
	default:
		return "", fmt.Errorf("simpleCalculator: unknown operation '%s'. Supported operations: 'add', 'subtract'", operation)
	}

	return fmt.Sprintf("%f", result), nil
}
