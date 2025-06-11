package plan

import (
	"plandex-server/db"
	"plandex-server/model"
	"plandex-server/types"
	"time"

	shared "plandex-shared"

	"github.com/PlandexAI/plandex/app/server/rag"
	"github.com/sashabaranov/go-openai"
)

type activeTellStreamState struct {
	activePlan            *types.ActivePlan
	modelStreamId         string
	clients               map[string]model.ClientInfo
	req                   *shared.TellPlanRequest
	auth                  *types.ServerAuth
	currentOrgId          string
	currentUserId         string
	plan                  *db.Plan
	branch                string
	iteration             int
	replyId               string
	modelContext          []*db.Context
	hasContextMap         bool
	contextMapEmpty       bool
	convo                 []*db.ConvoMessage
	promptConvoMessage    *db.ConvoMessage
	currentPlanState      *shared.CurrentPlanState
	missingFileResponse   shared.RespondMissingFileChoice
	summaries             []*db.ConvoSummary
	summarizedToMessageId string
	latestSummaryTokens   int
	userPrompt            string
	promptMessage         *openai.ChatCompletionMessage
	replyParser           *types.ReplyParser
	replyNumTokens        int
	messages              []types.ExtendedChatMessage
	tokensBeforeConvo     int
	totalRequestTokens    int
	settings              *shared.PlanSettings
	subtasks              []*db.Subtask
	currentSubtask        *db.Subtask
	hasAssistantReply     bool
	currentStage          shared.CurrentStage
	chunkProcessor        *chunkProcessor
	generationId          string

	requestStartedAt time.Time
	firstTokenAt     time.Time
	originalReq      *types.ExtendedChatCompletionRequest
	modelConfig      *shared.ModelRoleConfig
	fallbackRes      shared.FallbackResult

	skipConvoMessages map[string]bool

	manualStop []string

	numErrorRetry     int
	numFallbackRetry  int
	modelErr          *shared.ModelError
	noCacheSupportErr bool
	ragVectorStore    *rag.SQLiteVectorStore // For RAG context retrieval

	// Fields for MCP Tool Call
	pendingToolCall       *shared.PlandexToolInvocationRequest
	pendingToolDefinition *shared.MCPToolDefinition
	// allConvoMessages stores the accumulated conversation history for the current "tell" session,
	// including original user prompt, assistant replies, tool calls, and tool results.
	// It is used to build the messages sent to the LLM for each turn.
	allConvoMessages []types.ExtendedChatMessage
}

// ClearPendingToolCall resets the pending tool call and definition on the state.
func (state *activeTellStreamState) ClearPendingToolCall() {
	state.pendingToolCall = nil
	state.pendingToolDefinition = nil
}

type chunkProcessor struct {
	replyOperations                 []*shared.Operation
	chunksReceived                  int
	maybeRedundantOpeningTagContent string
	fileOpen                        bool
	contentBuffer                   string
	awaitingBlockOpeningTag         bool
	awaitingBlockClosingTag         bool
	awaitingOpClosingTag            bool
	awaitingBackticks               bool
}
