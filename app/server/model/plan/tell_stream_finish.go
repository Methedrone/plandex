package plan

import (
	"fmt"
	"log"
	"net/http"
	"plandex-server/db"
	"plandex-server/notify"
	"plandex-server/types"
	"runtime/debug"
	"time"

	shared "plandex-shared"

	"github.com/davecgh/go-spew/spew"
	"strings" // Added for TryParseToolInvocation logic, though mcp_parser uses it.
)

const MaxAutoContinueIterations = 200

type handleStreamFinishedResult struct {
	shouldContinueMainLoop bool
	shouldReturn           bool
}

func (state *activeTellStreamState) handleStreamFinished() handleStreamFinishedResult {
	planId := state.plan.Id
	branch := state.branch
	auth := state.auth
	plan := state.plan
	req := state.req
	clients := state.clients
	settings := state.settings
	currentOrgId := state.currentOrgId
	summaries := state.summaries
	convo := state.convo
	iteration := state.iteration
	replyOperations := state.chunkProcessor.replyOperations

	err := state.setActivePlan()
	if err != nil {
		state.onActivePlanMissingError()
		return handleStreamFinishedResult{
			shouldContinueMainLoop: true,
			shouldReturn:           false,
		}
	}

	active := state.activePlan

	// --- MCP Tool Invocation Parsing ---
	var isToolCall bool
	var toolRequest *shared.PlandexToolInvocationRequest
	var matchedToolDef *shared.MCPToolDefinition
	// Ensure settings and MCP settings are available before trying to parse
	if state.settings != nil && state.settings.MCPSettings != nil && state.settings.MCPSettings.Enabled && len(state.settings.MCPSettings.Tools) > 0 {
		var validationErr error
		log.Printf("MCP: Checking for tool invocation in response: %s", active.CurrentReplyContent)
		toolRequest, isToolCall, validationErr, matchedToolDef = TryParseToolInvocation(active.CurrentReplyContent, state.settings.MCPSettings.Tools)

		if isToolCall {
			if validationErr != nil {
				log.Printf("MCP: Invalid tool invocation attempt: %v. Treating as normal message.", validationErr)
				isToolCall = false // Invalidate the tool call due to error
				toolRequest = nil
				matchedToolDef = nil
			} else {
				log.Printf("MCP: Valid tool invocation request for tool '%s' detected and validated.", toolRequest.ToolName)
				// TODO: Store toolRequest and matchedToolDef in state for execution in the next cycle.
				// For now, we will prevent this message from being further processed as a regular reply.
				// The actual tool execution step will handle sending a result back to the LLM.

				// Clear the current reply content as it's a tool call, not a message to the user.
				// This might need adjustment based on how `storeOnFinished` and `summarizeConvo` use CurrentReplyContent.
				// For now, let's assume this is the right approach to prevent user display.
				// active.CurrentReplyContent = fmt.Sprintf("[Tool call to '%s' processing...]", toolRequest.ToolName) // Placeholder content

				// Signal that a tool call is pending. This needs a new field in activeTellStreamState or activePlan.
				// For now, we'll just log and then potentially alter flow later in this function.
				// This is where the flow would diverge significantly.
				// For this subtask, we'll focus on detection and validation.
				// The next subtask will handle the execution flow.

				// For now, effectively stop normal processing of this "message"
				// We might need to send a specific signal back to mainLoop or tell_exec
				// or set a state that the next action is tool execution.

				// Let's assume we set something on the state:
				state.pendingToolCall = toolRequest
				state.pendingToolDefinition = matchedToolDef

				// The rest of handleStreamFinished might need to be conditional based on pendingToolCall.
				// For now, let's just log and proceed to see where it breaks or needs adjustment.
			}
		}
	}
	// --- End MCP Tool Invocation Parsing ---

	// If it's a tool call, we might want to bypass much of the standard message processing.
	if state.pendingToolCall != nil {
		log.Printf("MCP: Tool call to '%s' is pending. Skipping standard stream finish processing.", state.pendingToolCall.ToolName)
		// TODO: Implement the actual tool execution flow.
		// For now, we will let it go through summarizeConvo and storeOnFinished,
		// but these might need to be aware of the tool call.
		// A simple solution for now is to just mark the reply as "processed" so it doesn't show to user,
		// and then the next iteration of the tell loop would pick up the pendingToolCall.
		// This means the 'summarizeConvo' and 'storeOnFinished' might operate on an empty/modified reply.

		// Send a message to active.CurrentReplyDoneCh to unblock the main loop
		// This is important because the main loop waits on this channel.
		log.Println("MCP: Signaling CurrentReplyDoneCh for pending tool call.")
		active.CurrentReplyDoneCh <- true
		log.Println("MCP: Resetting active.CurrentReplyDoneCh for pending tool call.")
		UpdateActivePlan(planId, branch, func(ap *types.ActivePlan) {
			ap.CurrentStreamingReplyId = ""
			ap.CurrentReplyDoneCh = nil
			// Potentially set a status on ActivePlan like "awaiting_tool_execution"
		})

		// TODO: Trigger the next step in execTellPlan for tool execution
		// This might involve returning a specific result or setting a flag that execTellPlan checks.
		// For now, returning true to stop further processing in this function.
		return handleStreamFinishedResult{shouldReturn: true} // Stop further processing in this function
	}


	time.Sleep(30 * time.Millisecond)
	active.FlushStreamBuffer()
	time.Sleep(100 * time.Millisecond)

	active.Stream(shared.StreamMessage{
		Type: shared.StreamMessageDescribing,
	})
	active.FlushStreamBuffer()

	err = db.SetPlanStatus(planId, branch, shared.PlanStatusDescribing, "")
	if err != nil {
		res := state.onError(onErrorParams{
			streamErr: fmt.Errorf("failed to set plan status to describing: %v", err),
			storeDesc: true,
		})

		return handleStreamFinishedResult{
			shouldContinueMainLoop: res.shouldContinueMainLoop,
			shouldReturn:           res.shouldReturn,
		}
	}

	autoLoadContextResult := state.checkAutoLoadContext()
	addedSubtasks := state.checkNewSubtasks()
	removedSubtasks := state.checkRemoveSubtasks()
	hasNewSubtasks := len(addedSubtasks) > 0

	log.Println("removedSubtasks:\n", spew.Sdump(removedSubtasks))
	log.Println("addedSubtasks:\n", spew.Sdump(addedSubtasks))
	log.Println("hasNewSubtasks:\n", hasNewSubtasks)

	handleDescAndExecStatusRes := state.handleDescAndExecStatus()
	if handleDescAndExecStatusRes.shouldContinueMainLoop || handleDescAndExecStatusRes.shouldReturn {
		return handleDescAndExecStatusRes.handleStreamFinishedResult
	}
	generatedDescription := handleDescAndExecStatusRes.generatedDescription
	subtaskFinished := handleDescAndExecStatusRes.subtaskFinished

	log.Printf("subtaskFinished: %v\n", subtaskFinished)

	storeOnFinishedResult := state.storeOnFinished(storeOnFinishedParams{
		replyOperations:       replyOperations,
		generatedDescription:  generatedDescription,
		subtaskFinished:       subtaskFinished,
		hasNewSubtasks:        hasNewSubtasks,
		autoLoadContextResult: autoLoadContextResult,
		addedSubtasks:         addedSubtasks,
		removedSubtasks:       removedSubtasks,
	})
	if storeOnFinishedResult.shouldContinueMainLoop || storeOnFinishedResult.shouldReturn {
		return storeOnFinishedResult.handleStreamFinishedResult
	}
	allSubtasksFinished := storeOnFinishedResult.allSubtasksFinished

	log.Println("allSubtasksFinished:\n", spew.Sdump(allSubtasksFinished))


	// summarize convo needs to come *after* the reply is stored in order to correctly summarize the latest message
	// If it was a tool call, CurrentReplyContent might be empty or a placeholder.
	// Summarization logic should ideally handle this (e.g., not summarize tool calls or summarize the intent).
	log.Println("summarizing convo in background")
	// summarize in the background
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("panic in summarizeConvo: %v\n%s", r, debug.Stack())
				active.StreamDoneCh <- &shared.ApiError{
					Type:   shared.ApiErrorTypeOther,
					Status: http.StatusInternalServerError,
					Msg:    fmt.Sprintf("Error summarizing convo: %v", r),
				}
			}
		}()

		err := summarizeConvo(clients, settings.ModelPack.PlanSummary, summarizeConvoParams{
			auth:                  auth,
			plan:                  plan,
			branch:                branch,
			convo:                 convo,
			summaries:             summaries,
			userPrompt:            state.userPrompt,
			currentOrgId:          currentOrgId,
			currentReply:          active.CurrentReplyContent,
			currentReplyNumTokens: active.NumTokens,
			modelPackName:         settings.ModelPack.Name,
		}, active.SummaryCtx)

		if err != nil {
			log.Printf("Error summarizing convo: %v\n", err)
			active.StreamDoneCh <- err
		}
	}()

	log.Println("Sending active.CurrentReplyDoneCh <- true")

	active.CurrentReplyDoneCh <- true

	log.Println("Resetting active.CurrentReplyDoneCh")

	UpdateActivePlan(planId, branch, func(ap *types.ActivePlan) {
		ap.CurrentStreamingReplyId = ""
		ap.CurrentReplyDoneCh = nil
	})

	autoLoadPaths := autoLoadContextResult.autoLoadPaths
	log.Printf("len(autoLoadPaths): %d\n", len(autoLoadPaths))
	if len(autoLoadPaths) > 0 {
		log.Println("Sending stream message to load context files")

		go func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("panic streaming auto-load context: %v\n%s", r, debug.Stack())
					go notify.NotifyErr(notify.SeverityError, fmt.Errorf("panic streaming auto-load context: %v\n%s", r, debug.Stack()))
				}
			}()

			active.Stream(shared.StreamMessage{
				Type:             shared.StreamMessageLoadContext,
				LoadContextFiles: autoLoadPaths,
			})
			active.FlushStreamBuffer()
		}()

		log.Println("Waiting for client to auto load context (30s timeout)")

		select {
		case <-active.Ctx.Done():
			log.Println("Context cancelled while waiting for auto load context")
			state.execHookOnStop(false)
			return handleStreamFinishedResult{
				shouldContinueMainLoop: false,
				shouldReturn:           true,
			}
		case <-time.After(30 * time.Second):
			log.Println("Timeout waiting for auto load context")
			res := state.onError(onErrorParams{
				streamErr: fmt.Errorf("timeout waiting for auto load context response"),
				storeDesc: true,
			})
			return handleStreamFinishedResult{
				shouldContinueMainLoop: res.shouldContinueMainLoop,
				shouldReturn:           res.shouldReturn,
			}
		case <-active.AutoLoadContextCh:
		}
	}

	willContinue := state.willContinuePlan(willContinuePlanParams{
		hasNewSubtasks:      hasNewSubtasks,
		allSubtasksFinished: allSubtasksFinished,
		activatePaths:       autoLoadContextResult.activatePaths,
		removedSubtasks:     len(removedSubtasks) > 0,
		hasExplicitPaths:    autoLoadContextResult.hasExplicitPaths,
	})

	if willContinue {
		log.Println("Auto continue plan")
		// continue plan
		execTellPlan(execTellPlanParams{
			clients:   clients,
			plan:      plan,
			branch:    branch,
			auth:      auth,
			req:       req,
			iteration: iteration + 1,
		})
	} else {
		var buildFinished bool
		UpdateActivePlan(planId, branch, func(ap *types.ActivePlan) {
			buildFinished = ap.BuildFinished()
			ap.RepliesFinished = true
		})

		log.Printf("Won't continue plan. Build finished: %v\n", buildFinished)

		time.Sleep(50 * time.Millisecond)

		if buildFinished {
			log.Println("Reply is finished and build is finished, calling active.Finish()")
			active := GetActivePlan(planId, branch)

			if active == nil {
				state.onActivePlanMissingError()
				return handleStreamFinishedResult{
					shouldContinueMainLoop: true,
					shouldReturn:           false,
				}
			}

			active.Finish()
		} else {
			log.Println("Plan is still building")
			log.Println("Updating status to building")
			err := db.SetPlanStatus(planId, branch, shared.PlanStatusBuilding, "")
			if err != nil {
				log.Printf("Error setting plan status to building: %v\n", err)
				go notify.NotifyErr(notify.SeverityError, fmt.Errorf("error setting plan status to building: %v", err))

				active.StreamDoneCh <- &shared.ApiError{
					Type:   shared.ApiErrorTypeOther,
					Status: http.StatusInternalServerError,
					Msg:    "Error setting plan status to building",
				}

				return handleStreamFinishedResult{
					shouldContinueMainLoop: true,
					shouldReturn:           false,
				}
			}

			log.Println("Sending RepliesFinished stream message")
			active.Stream(shared.StreamMessage{
				Type: shared.StreamMessageRepliesFinished,
			})

		}
	}

	return handleStreamFinishedResult{}
}
