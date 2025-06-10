package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"plandex-server/db" // For database interactions like GetPlanSettings, UpdatePlanSettings
	"plandex-server/server/client_helper" // For GetPlanFromRequest
	"plandex-server/types"                // For types.ServerAuth if needed by GetPlanFromRequest

	shared "plandex-shared" // For MCPToolDefinition, MCPConfig

	"github.com/gorilla/mux" // For mux.Vars to get path parameters
)

// HandleAddMCPTool adds a new MCP tool definition to a plan's configuration.
func HandleAddMCPTool(w http.ResponseWriter, r *http.Request) {
	plan, auth, apiErr := client_helper.GetPlanFromRequest(r, db.DB)
	if apiErr != nil {
		http.Error(w, apiErr.Msg, apiErr.Status)
		return
	}

	var newTool shared.MCPToolDefinition
	if err := json.NewDecoder(r.Body).Decode(&newTool); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if newTool.ToolName == "" {
		http.Error(w, "ToolName is required", http.StatusBadRequest)
		return
	}
	// Basic validation for other fields could be added here if necessary

	planConfig, err := db.GetPlanSettings(plan, true) // true for include_internal
	if err != nil {
		http.Error(w, "Failed to get plan configuration: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if planConfig == nil { // Should not happen if GetPlanSettings returns no error for existing plan
		planConfig = &shared.PlanConfig{}
	}

	if planConfig.MCPSettings == nil {
		planConfig.MCPSettings = &shared.MCPConfig{}
	}
	if planConfig.MCPSettings.Tools == nil {
		planConfig.MCPSettings.Tools = []shared.MCPToolDefinition{}
	}

	// Validate: Check if a tool with the same ToolName already exists
	for _, existingTool := range planConfig.MCPSettings.Tools {
		if existingTool.ToolName == newTool.ToolName {
			http.Error(w, fmt.Sprintf("Tool with name '%s' already exists", newTool.ToolName), http.StatusConflict)
			return
		}
	}

	planConfig.MCPSettings.Tools = append(planConfig.MCPSettings.Tools, newTool)

	err = db.UpdatePlanSettings(plan.Id, planConfig, auth.User.Id)
	if err != nil {
		http.Error(w, "Failed to update plan configuration: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTool)
}

// HandleListMCPTools lists all MCP tool definitions for a plan.
func HandleListMCPTools(w http.ResponseWriter, r *http.Request) {
	plan, _, apiErr := client_helper.GetPlanFromRequest(r, db.DB)
	if apiErr != nil {
		http.Error(w, apiErr.Msg, apiErr.Status)
		return
	}

	planConfig, err := db.GetPlanSettings(plan, true) // true for include_internal
	if err != nil {
		http.Error(w, "Failed to get plan configuration: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if planConfig == nil || planConfig.MCPSettings == nil || planConfig.MCPSettings.Tools == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]shared.MCPToolDefinition{}) // Return empty list
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(planConfig.MCPSettings.Tools)
}

// HandleRemoveMCPTool removes an MCP tool definition from a plan's configuration.
func HandleRemoveMCPTool(w http.ResponseWriter, r *http.Request) {
	plan, auth, apiErr := client_helper.GetPlanFromRequest(r, db.DB)
	if apiErr != nil {
		http.Error(w, apiErr.Msg, apiErr.Status)
		return
	}

	vars := mux.Vars(r)
	toolName, ok := vars["toolName"]
	if !ok || toolName == "" {
		http.Error(w, "Tool name is missing in URL", http.StatusBadRequest)
		return
	}

	planConfig, err := db.GetPlanSettings(plan, true) // true for include_internal
	if err != nil {
		http.Error(w, "Failed to get plan configuration: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if planConfig == nil || planConfig.MCPSettings == nil || planConfig.MCPSettings.Tools == nil {
		http.Error(w, fmt.Sprintf("Tool with name '%s' not found", toolName), http.StatusNotFound)
		return
	}

	found := false
	var updatedTools []shared.MCPToolDefinition
	for _, existingTool := range planConfig.MCPSettings.Tools {
		if existingTool.ToolName == toolName {
			found = true
		} else {
			updatedTools = append(updatedTools, existingTool)
		}
	}

	if !found {
		http.Error(w, fmt.Sprintf("Tool with name '%s' not found", toolName), http.StatusNotFound)
		return
	}

	planConfig.MCPSettings.Tools = updatedTools

	err = db.UpdatePlanSettings(plan.Id, planConfig, auth.User.Id)
	if err != nil {
		http.Error(w, "Failed to update plan configuration: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // Or http.StatusOK with a success message
}
