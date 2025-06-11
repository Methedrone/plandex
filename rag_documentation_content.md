# Retrieval Augmented Generation (RAG) in Plandex

Retrieval Augmented Generation (RAG) enhances Plandex's ability to understand and discuss your project by automatically fetching relevant snippets from your codebase to augment the context provided to the AI model.

## Concept

When RAG is enabled, Plandex creates a searchable index of your project's files. When you interact with the AI (e.g., using the `plandex tell` command), Plandex:
1. Embeds your query to understand its meaning.
2. Searches the project index for text chunks that are semantically similar to your query.
3. Includes these retrieved chunks as additional context for the AI model.

This allows the AI to have more relevant, up-to-date information from your files, even if those files weren't explicitly added to the context manually or by the auto-context feature. It's particularly useful for large projects where manually managing all relevant context can be challenging.

## Setup & Indexing

1.  **Enable RAG**:
    To use RAG, you first need to enable it in your plan's configuration:
    ```bash
    plandex set-config ragenabled true
    ```
    This setting can be applied to the current plan or set globally using `plandex set-config default ragenabled true`.

2.  **OpenAI API Key**:
    RAG uses OpenAI embeddings to create the index and to embed your queries. Ensure your `OPENAI_API_KEY` environment variable is correctly set up, similar to how it's needed for general Plandex AI interactions.

3.  **Create/Update Project Index**:
    Once enabled, you need to index your project. Run the following command from your project's root directory:
    ```bash
    plandex index
    ```
    This command will:
    *   Scan your project files (respecting `.gitignore` and `.plandexignore`).
    *   Split files into manageable text chunks.
    *   Generate vector embeddings for each chunk using OpenAI.
    *   Store these embeddings in a local database within your project's `.plandex` directory (specifically, `.plandex/rag.db`).

    The first indexing can take some time for large projects. Subsequent runs will update the index based on changed or new files.

## Configuration

You can customize RAG behavior through plan settings:

*   **`rag.topn`**: Determines how many top relevant chunks are retrieved and added to the context.
    *   Example: `plandex set-config rag.topn 5`
    *   Default: If not set or set to 0, a system default (e.g., 3) is used.
    *   A higher number provides more context but consumes more tokens.

*   **`rag.chunksizetokens`**: Defines the approximate target size (in tokens) for text chunks during indexing.
    *   Example: `plandex set-config rag.chunksizetokens 256`
    *   Default: If not set or set to 0, a system default (e.g., 512 tokens, which translates to an approximate character limit) is used.
    *   Smaller chunks can provide more granular context but might miss broader relationships. Larger chunks provide more context per hit but might be less focused. The current implementation uses an approximate character count based on this token value.

You can view your current RAG settings using `plandex config`.

## How It Works (Briefly)

1.  **Indexing**: The `plandex index` command processes your code. Files are broken down into smaller text "chunks." Each chunk is converted into a numerical representation (a vector embedding) that captures its semantic meaning. These embeddings are stored locally.
2.  **Retrieval**: When you make a query (e.g., "implement a function to handle user uploads"), Plandex converts your query into an embedding.
3.  **Search**: This query embedding is compared against the embeddings of all indexed chunks to find the ones most semantically similar to your query. The `rag.topn` setting controls how many of these are selected.
4.  **Augmentation**: The text content of these selected chunks is then added to the context window of the AI model, just before your actual query and other context files. This gives the model more specific, relevant information from your codebase to draw upon when generating its response.

## Troubleshooting & Notes

*   **Re-indexing**: If you make significant changes to your project files, add new files, or remove old ones, you should run `plandex index` again to keep the RAG index up-to-date and ensure the AI gets the most relevant information.
*   **Index Location**: The RAG index is stored in `.plandex/rag.db` within your project. This file can be safely deleted if you want to force a complete re-index from scratch (you'll need to run `plandex index` afterwards). It's generally recommended to add `.plandex/` to your project's `.gitignore` file.
*   **Cost**: Generating embeddings uses the OpenAI API and will incur costs, similar to other AI interactions. Indexing a very large project for the first time might result in noticeable API usage.
*   **Relevance**: The effectiveness of RAG depends on the quality of the embeddings and the nature of your queries. Sometimes, the retrieved chunks might not be perfectly relevant, but often they provide useful context the AI wouldn't have otherwise.
*   **Supported Files**: The indexer processes text-based files and currently uses a list of common source code file extensions. Binary files are ignored.
*   **Interaction with other context**: RAG-retrieved context is added alongside context files you've manually loaded or that Plandex's auto-context feature has loaded. It's one of several sources of information for the AI.

---
MCP Documentation Content (for mcp.md)
---

# Model Context Protocol (MCP) in Plandex

The Model Context Protocol (MCP) allows Plandex to extend the capabilities of the AI model by enabling it to invoke predefined tools. These tools can range from server-side functions to external HTTP APIs, allowing the AI to gather information, perform actions, and interact with other systems.

## Concept

MCP empowers the AI model to go beyond text generation and code modification. When the AI determines that it needs to perform an action or fetch information not readily available in its current context, it can request the invocation of a specific tool. Plandex then executes this tool and feeds the output back to the AI, allowing it to continue its task with new information or the result of an action.

This creates a powerful feedback loop where the AI can:
1.  Analyze a request.
2.  Identify the need for a tool.
3.  Request tool invocation with specific inputs.
4.  Receive the tool's output.
5.  Use this output to generate a more informed response, make further plans, or even invoke other tools.

## Defining Tools

MCP Tools are defined in JSON format. Each tool definition specifies its name, description, input/output schemas, and how it should be executed.

**Tool Definition Structure (`shared.MCPToolDefinition`)**:

```json
{
  "toolName": "unique_tool_name",
  "description": "Detailed description for the LLM.",
  "inputSchema": "{...JSON Schema...}",
  "outputSchema": "{...JSON Schema...}",
  "executionType": "http | predefined_server_function",
  "executionDetails": {
    // Structure depends on executionType
  }
}
```

**Field Explanations**:

*   `toolName` (string, required): Unique name for the tool (e.g., `getWeather`, `searchProjectIssues`).
*   `description` (string, required): Detailed description for the LLM to understand the tool's purpose, capabilities, and when to use it. This is critical for the LLM to make good decisions.
*   `inputSchema` (string, required): A JSON Schema string defining the expected input parameters for the tool. The LLM will attempt to generate inputs matching this schema. Example:
    ```json
    // For a tool expecting a "city" string and "days" integer:
    {
      "type": "object",
      "properties": {
        "city": { "type": "string", "description": "The city name." },
        "days": { "type": "integer", "description": "Number of forecast days." }
      },
      "required": ["city"]
    }
    ```
*   `outputSchema` (string, optional): A JSON Schema string defining the expected output structure from the tool. This is mainly for documentation, validation on the server side, or potentially for the LLM if it needs to understand the output structure in advance.
*   `executionType` (string, required): Specifies how the tool is executed. Current valid values:
    *   `"http"`: The tool makes an HTTP request to an external service.
    *   `"predefined_server_function"`: The tool calls a function that is already defined within the Plandex server.
*   `executionDetails` (object, required): A map containing specific details for the execution, the structure of which depends on `executionType`.
    *   **For `http`**:
        *   `"url"` (string, required): The endpoint URL. Can include placeholders like `{param}` that will be replaced by values from the `tool_input` (see HTTP Tools section).
        *   `"method"` (string, required): The HTTP method (e.g., "GET", "POST", "PUT").
        *   (Optional) Other fields like `"headers"` (map[string]string) or specific parameter mapping configurations could be added here in future extensions.
    *   **For `predefined_server_function`**:
        *   `"functionName"` (string, required): The registered name of the server-side function to call.

**Example Tool Definition File (`get_weather_tool.json`)**:
```json
{
  "toolName": "getWeatherForecast",
  "description": "Fetches the weather forecast for a specified city for a number of days. Use this if the user asks about weather conditions.",
  "inputSchema": "{\n  \"type\": \"object\",\n  \"properties\": {\n    \"city\": { \"type\": \"string\", \"description\": \"The city for which to get the weather forecast.\" },\n    \"num_days\": { \"type\": \"integer\", \"description\": \"The number of days for the forecast, e.g., 1 for today, 7 for a week.\", \"default\": 1 }\n  },\n  \"required\": [\"city\"]\n}",
  "outputSchema": "{\n  \"type\": \"object\",\n  \"properties\": {\n    \"forecast\": { \"type\": \"string\", \"description\": \"A summary of the weather forecast.\" },\n    \"details\": { \"type\": \"array\", \"items\": { \"type\": \"object\", \"properties\": { \"day\": {\"type\":\"string\"}, \"condition\":{\"type\":\"string\"}, \"temp_high\":{\"type\":\"string\"}, \"temp_low\":{\"type\":\"string\"}}}} \n  }\n}",
  "executionType": "http",
  "executionDetails": {
    "url": "https://api.weather fictional.com/forecast?city={city}&days={num_days}",
    "method": "GET"
  }
}
```

## Managing Tools

Tool definitions are associated with a specific plan and are stored within the plan's configuration.

1.  **Enable MCP**:
    First, enable MCP for the current plan (or default):
    ```bash
    plandex set-config mcpenabled true
    ```

2.  **Add a Tool**:
    Use the `plandex mcp add-tool` command with a JSON file containing the tool definition:
    ```bash
    plandex mcp add-tool --file path/to/your_tool_definition.json
    ```
    Plandex will validate the tool definition (e.g., check for `toolName` uniqueness within the plan).

3.  **List Tools**:
    To see all tools configured for the current plan:
    ```bash
    plandex mcp list-tools
    ```
    This will display a table of tool names, descriptions, and execution types.

4.  **Remove a Tool**:
    To remove a tool by its name:
    ```bash
    plandex mcp remove-tool <toolName>
    ```
    You will be asked for confirmation before the tool is removed.

## How It Works

1.  **LLM Awareness**: When MCP is enabled and tools are defined, Plandex includes the list of available tools (name, description, input schema) in the system prompt sent to the LLM.
2.  **Invocation Request**: If the LLM determines a tool is needed, it is instructed to respond *only* with a specific JSON object:
    ```json
    {
      "plandex_tool_invocation": {
        "tool_name": "the_tool_name_to_call",
        "tool_input": { /* parameters matching the tool's Input Schema */ }
      }
    }
    ```
3.  **Server Parsing & Validation**: The Plandex server parses this JSON.
    *   It checks if the `tool_name` is valid and defined for the plan.
    *   It validates the provided `tool_input` against the tool's `inputSchema` using a JSON Schema validator.
4.  **Server Execution**: If the request is valid, the server executes the tool:
    *   For HTTP tools, it makes the configured HTTP request.
    *   For predefined tools, it calls the corresponding server-side Go function.
5.  **Feedback Loop**: The output from the tool (or an error message if execution failed) is then formatted into a message and sent back to the LLM. This message typically includes the `tool_name` and its `tool_output`.
6.  **LLM Continues**: The LLM receives the tool's output and can then use this information to formulate its next response to the user, potentially make further plans, or even decide to call another tool.

## Predefined Server Tools

Plandex can include built-in server-side functions that tools can invoke. Here are some examples:

*   **`echoTool`**:
    *   **Description**: A simple debugging tool that returns its input as its output.
    *   **Input Schema**: Can be any valid JSON (e.g., `{"type": "object"}` or more specific).
    *   **ExecutionDetails**: `{"functionName": "echoTool"}`
    *   **Output**: The JSON string representation of the `tool_input` it received.

*   **`simpleCalculator`**:
    *   **Description**: Performs basic arithmetic (add, subtract).
    *   **Input Schema**:
        ```json
        {
          "type": "object",
          "properties": {
            "a": { "type": "number", "description": "First operand." },
            "b": { "type": "number", "description": "Second operand." },
            "operation": { "type": "string", "enum": ["add", "subtract"], "description": "The operation to perform." }
          },
          "required": ["a", "b", "operation"]
        }
        ```
    *   **ExecutionDetails**: `{"functionName": "simpleCalculator"}`
    *   **Output**: A string representing the numerical result of the calculation (e.g., `"25"`). Returns an error message if inputs are invalid or the operation is unsupported.

*(More predefined tools can be added to the Plandex server over time.)*

## HTTP Tools

HTTP tools allow the LLM to interact with external APIs.

*   **Configuration**:
    *   `url`: The full URL for the API endpoint.
    *   `method`: Standard HTTP methods like "GET", "POST", "PUT", "DELETE", "PATCH".
*   **Parameter Passing**:
    *   **URL Path Templating**: You can include placeholders in the URL like `https://api.example.com/users/{userID}`. If `userID` is a key in the `tool_input` provided by the LLM, its value will be substituted into the URL.
    *   **Query Parameters (for GET/DELETE)**: For GET or DELETE requests, any `tool_input` fields not used in URL path templating will be automatically converted to URL query parameters.
    *   **JSON Request Body (for POST/PUT/PATCH)**: For these methods, `tool_input` fields not used in URL path templating are typically sent as the JSON request body. The `Content-Type` header is automatically set to `application/json`.
    *   *(Future: More advanced parameter mapping from `tool_input` to request body, headers, or query params might be supported via an explicit `parameterMapping` field in `executionDetails`.)*
*   **Security Notes**:
    *   **Authentication**: Currently, direct support for complex authentication schemes (like OAuth2) within the tool definition is limited. If an API key is needed, it might need to be part of the `url` (if the API supports it as a query param) or embedded in `headers` if that feature is added to `executionDetails`. For more complex auth, a custom predefined server tool acting as a proxy might be necessary.
    *   **Server-Side Requests**: All HTTP requests are made from the Plandex server, not the user's client.
    *   **Data Exposure**: Be mindful of what data the LLM might send to external APIs via `tool_input`.
    *   **Response Size**: The Plandex server limits the size of the response body read from the HTTP tool (e.g., to 1MB) to prevent abuse.

## Writing Effective Tool Descriptions

The `description` field in your `MCPToolDefinition` is crucial. The LLM relies heavily on this description to:
*   Understand what the tool does.
*   Know what kind of inputs it expects (beyond the schema, in natural language).
*   Decide when it's appropriate to use the tool versus trying to answer directly or asking the user for information.

**Tips for good descriptions**:
*   Be specific about the tool's capabilities and limitations.
*   Mention key input parameters and their purpose.
*   Provide clear indicators of when the tool should be used (e.g., "Use this tool if the user asks for X," or "Use to perform action Y").
*   If the tool has side effects (e.g., creates data), mention this.

## Troubleshooting & Notes

*   **Tool Not Found**: If the LLM tries to invoke a tool not defined in the plan's `MCPSettings`, Plandex will return an error to the LLM.
*   **Input Schema Validation**: If the LLM provides `tool_input` that doesn't match the tool's `inputSchema`, Plandex will return a validation error to the LLM. The LLM might then attempt to correct its input.
*   **Execution Errors**: If a tool fails during execution (e.g., HTTP request times out, predefined function panics), an error message is returned to the LLM as the `tool_output`.
*   **Debugging**: Check server logs for details on tool parsing, validation, and execution steps. The LLM's raw JSON request for tool invocation will also be logged.
*   **Iterative Process**: Tool usage is often iterative. The LLM might call a tool, get a result, and then decide to call another tool or the same tool with different parameters based on the new information.
*   **Token Limits**: Tool definitions, invocation JSON, and tool outputs all consume tokens in the conversation with the LLM. Be mindful of this, especially with verbose tool outputs or many tools.
