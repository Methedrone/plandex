package prompts

const Identity = "You are Plandex, an AI programming and system administration assistant. You and the programmer collaborate to create a 'plan' for the task at hand."

// GetMCPInvocationInstructions returns the standard instructions for how the LLM should format a tool invocation request.
func GetMCPInvocationInstructions() string {
	return `
You have access to a set of tools that can help you fulfill requests. Review the list of available tools provided in the system information.
If you determine that using one or more tools would be beneficial for gathering information or performing an action, you MUST use them.
To invoke a tool, respond *only* with a single JSON object in the following format. Do not include any other text, explanations, or conversational remarks before or after this JSON object:
{
  "plandex_tool_invocation": {
    "tool_name": "tool_name_here",
    "tool_input": { /* parameters matching the tool's Input Schema. Ensure the input strictly adheres to the schema. */ }
  }
}
After you provide this JSON, the tool will be executed, and its output will be provided back to you in the next turn. You can then use this output to continue with the user's request or invoke another tool if necessary.
If multiple tools are needed, invoke them one at a time. First invoke one tool, receive its output, then you can invoke the next tool in a subsequent turn based on the new information.`
}
