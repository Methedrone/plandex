# Phase 5: Next-Gen Capabilities Implementation Guide
## Plandex App Upgrade - Cutting-Edge AI & Innovation Features

---

## 🎯 EXECUTIVE SUMMARY

This guide provides comprehensive implementation of next-generation capabilities for the Plandex application. Building on the secure, performant, well-tested, and feature-rich foundation from Phases 1-4, this phase introduces cutting-edge AI capabilities, enterprise-grade features, and innovative technologies that position Plandex as a leading-edge development platform.

### Current Innovation State
- **AI Integration**: Basic model API calls with text-only processing
- **Context Understanding**: File-based context with limited intelligence
- **Mobile Experience**: None (desktop/web only)
- **Enterprise Features**: Basic organization support
- **Analytics**: Limited usage tracking
- **Multi-Modal AI**: Not implemented

### Next-Generation Targets
- **Multi-Modal AI Integration**: Vision models for diagrams, UI screenshots, video analysis
- **Intelligent Context Optimization**: ML-powered context selection and pattern learning
- **Progressive Web App**: Mobile-first, offline-capable interface
- **Enterprise Security & Compliance**: Advanced RBAC, audit logging, compliance frameworks
- **Advanced Analytics**: Predictive insights, project health monitoring, ML-driven recommendations
- **Distributed AI Processing**: Horizontal scaling, cost optimization, edge computing

---

## 🔧 CLAUDE CODE WORKFLOW INTEGRATION

### MCP Servers Utilization
```bash
# Research cutting-edge AI technologies
use context7

# Complex AI architecture planning
use SequentialThinking

# Innovation-focused development tasks
use Task-Master with next-gen capabilities PRD
```

### Innovation Development TodoWrite Strategy
This guide provides detailed TodoWrite checkpoints for each next-generation capability, ensuring systematic implementation of advanced features with future-proof architecture.

---

## 📋 DETAILED IMPLEMENTATION PLAN

## Phase 5A: Multi-Modal AI Integration
### 🔮 ADVANCED AI CAPABILITIES TARGET
**Goal**: Implement vision models, multi-modal understanding, and advanced AI processing

### Implementation Steps

#### Step 5A.1: Vision AI Integration
**File: `/app/server/ai/vision_processor.go`** (create vision AI processor)
```go
package ai

import (
    "bytes"
    "context"
    "encoding/base64"
    "encoding/json"
    "fmt"
    "image"
    "image/jpeg"
    "image/png"
    "io"
    "mime/multipart"
    "net/http"
    "strings"
    "time"
    
    "github.com/disintegration/imaging"
)

// VisionProcessor handles multi-modal AI vision capabilities
type VisionProcessor struct {
    aiClient      *MultiModalAIClient
    imageCache    *ImageCache
    optimizer     *ImageOptimizer
    analyzer      *VisualAnalyzer
    metrics       *VisionMetrics
}

// VisionAnalysisRequest represents a vision analysis request
type VisionAnalysisRequest struct {
    Images      []ImageInput     `json:"images"`
    Context     string          `json:"context,omitempty"`
    TaskType    VisionTaskType  `json:"task_type"`
    Language    string          `json:"language,omitempty"`
    Options     VisionOptions   `json:"options,omitempty"`
}

// ImageInput represents an input image for analysis
type ImageInput struct {
    Data        []byte           `json:"data,omitempty"`
    URL         string           `json:"url,omitempty"`
    Format      string           `json:"format"`
    Description string           `json:"description,omitempty"`
    Metadata    ImageMetadata    `json:"metadata,omitempty"`
}

// VisionAnalysisResult represents the result of vision analysis
type VisionAnalysisResult struct {
    TaskType        VisionTaskType        `json:"task_type"`
    Results         []VisionResult        `json:"results"`
    Confidence      float64               `json:"confidence"`
    ProcessingTime  time.Duration         `json:"processing_time"`
    TokensUsed      int                   `json:"tokens_used"`
    CostEstimate    float64               `json:"cost_estimate"`
    Suggestions     []ActionSuggestion    `json:"suggestions"`
    Metadata        AnalysisMetadata      `json:"metadata"`
}

// VisionResult represents a specific vision analysis result
type VisionResult struct {
    ImageIndex      int                   `json:"image_index"`
    Type           ResultType            `json:"type"`
    Content        string                `json:"content"`
    Coordinates    []BoundingBox         `json:"coordinates,omitempty"`
    Confidence     float64               `json:"confidence"`
    Details        map[string]interface{} `json:"details,omitempty"`
}

// Vision task types
type VisionTaskType string
const (
    VisionTaskOCR                VisionTaskType = "ocr"
    VisionTaskUIAnalysis         VisionTaskType = "ui_analysis"
    VisionTaskDiagramAnalysis    VisionTaskType = "diagram_analysis"
    VisionTaskCodeScreenshot     VisionTaskType = "code_screenshot"
    VisionTaskArchitectureDiagram VisionTaskType = "architecture_diagram"
    VisionTaskWireframeAnalysis  VisionTaskType = "wireframe_analysis"
    VisionTaskErrorScreenshot    VisionTaskType = "error_screenshot"
    VisionTaskDesignReview       VisionTaskType = "design_review"
)

// NewVisionProcessor creates a new vision processor
func NewVisionProcessor(aiClient *MultiModalAIClient) *VisionProcessor {
    return &VisionProcessor{
        aiClient:   aiClient,
        imageCache: NewImageCache(),
        optimizer:  NewImageOptimizer(),
        analyzer:   NewVisualAnalyzer(),
        metrics:    NewVisionMetrics(),
    }
}

// AnalyzeImages performs multi-modal vision analysis
func (vp *VisionProcessor) AnalyzeImages(ctx context.Context, request VisionAnalysisRequest) (*VisionAnalysisResult, error) {
    start := time.Now()
    defer func() {
        vp.metrics.RecordAnalysis(request.TaskType, time.Since(start))
    }()
    
    // Validate and prepare images
    processedImages, err := vp.prepareImages(request.Images)
    if err != nil {
        return nil, fmt.Errorf("failed to prepare images: %w", err)
    }
    
    // Build vision prompt based on task type
    prompt := vp.buildVisionPrompt(request.TaskType, request.Context, processedImages)
    
    // Perform vision analysis
    aiResponse, err := vp.aiClient.AnalyzeVision(ctx, VisionRequest{
        Images:      processedImages,
        Prompt:      prompt,
        MaxTokens:   request.Options.MaxTokens,
        Temperature: request.Options.Temperature,
    })
    
    if err != nil {
        vp.metrics.RecordError(request.TaskType, err)
        return nil, fmt.Errorf("vision analysis failed: %w", err)
    }
    
    // Parse and structure results
    results, err := vp.parseVisionResults(aiResponse, request.TaskType)
    if err != nil {
        return nil, fmt.Errorf("failed to parse vision results: %w", err)
    }
    
    // Generate actionable suggestions
    suggestions := vp.generateActionSuggestions(results, request.TaskType)
    
    return &VisionAnalysisResult{
        TaskType:       request.TaskType,
        Results:        results,
        Confidence:     vp.calculateConfidence(results),
        ProcessingTime: time.Since(start),
        TokensUsed:     aiResponse.TokensUsed,
        CostEstimate:   aiResponse.CostEstimate,
        Suggestions:    suggestions,
        Metadata: AnalysisMetadata{
            Model:     aiResponse.Model,
            Timestamp: time.Now(),
            Version:   "5.0.0",
        },
    }, nil
}

// AnalyzeUIScreenshot analyzes UI screenshots for usability and design
func (vp *VisionProcessor) AnalyzeUIScreenshot(ctx context.Context, imageData []byte, context string) (*UIAnalysisResult, error) {
    request := VisionAnalysisRequest{
        Images: []ImageInput{{
            Data:   imageData,
            Format: "image/png",
        }},
        Context:  context,
        TaskType: VisionTaskUIAnalysis,
    }
    
    result, err := vp.AnalyzeImages(ctx, request)
    if err != nil {
        return nil, err
    }
    
    return vp.convertToUIAnalysis(result), nil
}

// AnalyzeArchitectureDiagram analyzes system architecture diagrams
func (vp *VisionProcessor) AnalyzeArchitectureDiagram(ctx context.Context, imageData []byte) (*ArchitectureAnalysisResult, error) {
    request := VisionAnalysisRequest{
        Images: []ImageInput{{
            Data:   imageData,
            Format: "image/png",
        }},
        TaskType: VisionTaskArchitectureDiagram,
    }
    
    result, err := vp.AnalyzeImages(ctx, request)
    if err != nil {
        return nil, err
    }
    
    return vp.convertToArchitectureAnalysis(result), nil
}

// ExtractCodeFromScreenshot extracts code from screenshots using OCR + AI
func (vp *VisionProcessor) ExtractCodeFromScreenshot(ctx context.Context, imageData []byte, language string) (*CodeExtractionResult, error) {
    request := VisionAnalysisRequest{
        Images: []ImageInput{{
            Data:   imageData,
            Format: "image/png",
        }},
        Language: language,
        TaskType: VisionTaskCodeScreenshot,
    }
    
    result, err := vp.AnalyzeImages(ctx, request)
    if err != nil {
        return nil, err
    }
    
    return vp.convertToCodeExtraction(result), nil
}

// prepareImages optimizes and validates images for analysis
func (vp *VisionProcessor) prepareImages(images []ImageInput) ([]ProcessedImage, error) {
    var processed []ProcessedImage
    
    for i, img := range images {
        // Load image data
        var imageData []byte
        var err error
        
        if img.URL != "" {
            imageData, err = vp.downloadImage(img.URL)
            if err != nil {
                return nil, fmt.Errorf("failed to download image %d: %w", i, err)
            }
        } else {
            imageData = img.Data
        }
        
        // Validate image format
        format, err := vp.detectImageFormat(imageData)
        if err != nil {
            return nil, fmt.Errorf("invalid image format for image %d: %w", i, err)
        }
        
        // Optimize image for AI processing
        optimized, err := vp.optimizer.OptimizeForAI(imageData, format)
        if err != nil {
            return nil, fmt.Errorf("failed to optimize image %d: %w", i, err)
        }
        
        processed = append(processed, ProcessedImage{
            Index:       i,
            Data:        optimized.Data,
            Format:      optimized.Format,
            Width:       optimized.Width,
            Height:      optimized.Height,
            Size:        len(optimized.Data),
            Description: img.Description,
        })
    }
    
    return processed, nil
}

// buildVisionPrompt creates task-specific prompts for vision analysis
func (vp *VisionProcessor) buildVisionPrompt(taskType VisionTaskType, context string, images []ProcessedImage) string {
    basePrompt := fmt.Sprintf("You are an expert AI assistant analyzing %d image(s) for %s.", len(images), taskType)
    
    if context != "" {
        basePrompt += fmt.Sprintf(" Context: %s", context)
    }
    
    switch taskType {
    case VisionTaskUIAnalysis:
        return basePrompt + `
        
Please analyze the UI screenshot(s) and provide:
1. Overall design assessment (layout, typography, color scheme)
2. Usability issues and improvements
3. Accessibility concerns
4. Mobile responsiveness observations
5. Best practice recommendations
6. Component identification and structure
7. User experience flow analysis
        
Provide specific, actionable feedback with coordinates when possible.`

    case VisionTaskArchitectureDiagram:
        return basePrompt + `
        
Please analyze the architecture diagram(s) and provide:
1. System components identification
2. Data flow analysis
3. Integration points and dependencies
4. Scalability considerations
5. Security implications
6. Best practice recommendations
7. Potential bottlenecks or issues
8. Technology stack analysis
        
Structure your response with clear component relationships and recommendations.`

    case VisionTaskCodeScreenshot:
        return basePrompt + `
        
Please extract and analyze the code from the screenshot(s):
1. Extract all visible code with proper formatting
2. Identify the programming language
3. Analyze code quality and structure
4. Identify potential bugs or improvements
5. Suggest optimizations
6. Check for security issues
7. Verify syntax correctness
        
Return the extracted code in proper markdown code blocks with language specification.`

    case VisionTaskErrorScreenshot:
        return basePrompt + `
        
Please analyze the error screenshot(s) and provide:
1. Error message extraction and interpretation
2. Root cause analysis
3. Step-by-step troubleshooting guide
4. Prevention strategies
5. Related documentation links
6. Code fixes if applicable
        
Focus on providing actionable solutions to resolve the error.`

    case VisionTaskDiagramAnalysis:
        return basePrompt + `
        
Please analyze the diagram(s) and provide:
1. Diagram type identification (flowchart, UML, network, etc.)
2. Element and relationship analysis
3. Logic flow verification
4. Completeness assessment
5. Improvement suggestions
6. Standards compliance check
        
Provide detailed analysis of the diagram's effectiveness and accuracy.`

    default:
        return basePrompt + "\n\nPlease provide a detailed analysis of the image(s) with specific observations and recommendations."
    }
}

// parseVisionResults parses AI vision responses into structured results
func (vp *VisionProcessor) parseVisionResults(response *VisionResponse, taskType VisionTaskType) ([]VisionResult, error) {
    var results []VisionResult
    
    // Parse JSON response if structured
    if response.IsStructured {
        var structuredResults []VisionResult
        if err := json.Unmarshal([]byte(response.Content), &structuredResults); err == nil {
            return structuredResults, nil
        }
    }
    
    // Parse text response based on task type
    content := response.Content
    
    switch taskType {
    case VisionTaskCodeScreenshot:
        code := vp.extractCodeBlocks(content)
        for i, block := range code {
            results = append(results, VisionResult{
                ImageIndex: 0,
                Type:       ResultTypeCode,
                Content:    block.Code,
                Confidence: block.Confidence,
                Details: map[string]interface{}{
                    "language": block.Language,
                    "line_count": strings.Count(block.Code, "\n") + 1,
                },
            })
        }
        
    case VisionTaskUIAnalysis:
        analysis := vp.parseUIAnalysis(content)
        results = append(results, VisionResult{
            ImageIndex: 0,
            Type:       ResultTypeUIAnalysis,
            Content:    analysis.Summary,
            Confidence: analysis.Confidence,
            Details: map[string]interface{}{
                "issues":          analysis.Issues,
                "recommendations": analysis.Recommendations,
                "components":      analysis.Components,
            },
        })
        
    case VisionTaskArchitectureDiagram:
        architecture := vp.parseArchitectureAnalysis(content)
        results = append(results, VisionResult{
            ImageIndex: 0,
            Type:       ResultTypeArchitecture,
            Content:    architecture.Summary,
            Confidence: architecture.Confidence,
            Details: map[string]interface{}{
                "components":    architecture.Components,
                "connections":   architecture.Connections,
                "technologies":  architecture.Technologies,
                "recommendations": architecture.Recommendations,
            },
        })
        
    default:
        results = append(results, VisionResult{
            ImageIndex: 0,
            Type:       ResultTypeGeneral,
            Content:    content,
            Confidence: 0.8, // Default confidence
        })
    }
    
    return results, nil
}

// Image optimization for AI processing
type ImageOptimizer struct {
    maxWidth    int
    maxHeight   int
    quality     int
    maxFileSize int
}

func NewImageOptimizer() *ImageOptimizer {
    return &ImageOptimizer{
        maxWidth:    1920,
        maxHeight:   1080,
        quality:     85,
        maxFileSize: 2 * 1024 * 1024, // 2MB
    }
}

func (io *ImageOptimizer) OptimizeForAI(data []byte, format string) (*OptimizedImage, error) {
    // Decode image
    img, err := io.decodeImage(data, format)
    if err != nil {
        return nil, err
    }
    
    // Resize if too large
    bounds := img.Bounds()
    width, height := bounds.Max.X, bounds.Max.Y
    
    if width > io.maxWidth || height > io.maxHeight {
        img = imaging.Fit(img, io.maxWidth, io.maxHeight, imaging.Lanczos)
        bounds = img.Bounds()
        width, height = bounds.Max.X, bounds.Max.Y
    }
    
    // Encode optimized image
    var buf bytes.Buffer
    switch format {
    case "image/jpeg", "image/jpg":
        err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: io.quality})
    case "image/png":
        err = png.Encode(&buf, img)
    default:
        err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: io.quality})
        format = "image/jpeg"
    }
    
    if err != nil {
        return nil, err
    }
    
    optimizedData := buf.Bytes()
    
    // Check file size and reduce quality if necessary
    if len(optimizedData) > io.maxFileSize && format == "image/jpeg" {
        for quality := io.quality - 10; quality >= 50 && len(optimizedData) > io.maxFileSize; quality -= 10 {
            buf.Reset()
            jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality})
            optimizedData = buf.Bytes()
        }
    }
    
    return &OptimizedImage{
        Data:   optimizedData,
        Format: format,
        Width:  width,
        Height: height,
    }, nil
}

// Multi-modal AI client for vision tasks
type MultiModalAIClient struct {
    openAIClient     *OpenAIVisionClient
    anthropicClient  *AnthropicVisionClient
    googleClient     *GoogleVisionClient
    defaultProvider  string
    fallbackChain    []string
}

type VisionRequest struct {
    Images      []ProcessedImage `json:"images"`
    Prompt      string          `json:"prompt"`
    MaxTokens   int             `json:"max_tokens"`
    Temperature float64         `json:"temperature"`
    Model       string          `json:"model,omitempty"`
}

type VisionResponse struct {
    Content      string          `json:"content"`
    TokensUsed   int             `json:"tokens_used"`
    CostEstimate float64         `json:"cost_estimate"`
    Model        string          `json:"model"`
    IsStructured bool            `json:"is_structured"`
    Confidence   float64         `json:"confidence"`
}

func (mmac *MultiModalAIClient) AnalyzeVision(ctx context.Context, request VisionRequest) (*VisionResponse, error) {
    providers := append([]string{mmac.defaultProvider}, mmac.fallbackChain...)
    
    for _, provider := range providers {
        response, err := mmac.callProvider(ctx, provider, request)
        if err == nil {
            return response, nil
        }
        
        // Log error and try next provider
        fmt.Printf("Provider %s failed: %v\n", provider, err)
    }
    
    return nil, fmt.Errorf("all vision providers failed")
}

func (mmac *MultiModalAIClient) callProvider(ctx context.Context, provider string, request VisionRequest) (*VisionResponse, error) {
    switch provider {
    case "openai":
        return mmac.openAIClient.AnalyzeVision(ctx, request)
    case "anthropic":
        return mmac.anthropicClient.AnalyzeVision(ctx, request)
    case "google":
        return mmac.googleClient.AnalyzeVision(ctx, request)
    default:
        return nil, fmt.Errorf("unknown provider: %s", provider)
    }
}

// Specialized analysis results
type UIAnalysisResult struct {
    Summary         string              `json:"summary"`
    OverallScore    float64             `json:"overall_score"`
    Issues          []UIIssue           `json:"issues"`
    Recommendations []UIRecommendation  `json:"recommendations"`
    Components      []UIComponent       `json:"components"`
    Accessibility   AccessibilityReport `json:"accessibility"`
    Performance     UIPerformanceReport `json:"performance"`
}

type ArchitectureAnalysisResult struct {
    Summary         string                     `json:"summary"`
    Components      []ArchitectureComponent    `json:"components"`
    Connections     []ComponentConnection      `json:"connections"`
    Technologies    []TechnologyStack         `json:"technologies"`
    SecurityAnalysis SecurityArchitectureReport `json:"security"`
    Scalability     ScalabilityReport         `json:"scalability"`
    Recommendations []ArchitectureRecommendation `json:"recommendations"`
}

type CodeExtractionResult struct {
    ExtractedCode   string              `json:"extracted_code"`
    Language        string              `json:"language"`
    Confidence      float64             `json:"confidence"`
    QualityAnalysis CodeQualityReport   `json:"quality_analysis"`
    Suggestions     []CodeSuggestion    `json:"suggestions"`
    Errors          []ExtractedError    `json:"errors"`
}
```

**TodoWrite Task**: `Implement multi-modal AI vision processing capabilities`

#### Step 5A.2: Advanced Context Intelligence
**File: `/app/server/ai/context_intelligence.go`** (create intelligent context system)
```go
package ai

import (
    "context"
    "encoding/json"
    "fmt"
    "math"
    "sort"
    "strings"
    "time"
    
    "gonum.org/v1/gonum/mat"
    "github.com/sajari/fuzzy"
)

// ContextIntelligence provides ML-powered context optimization
type ContextIntelligence struct {
    vectorizer       *ContextVectorizer
    similarityEngine *SimilarityEngine
    patternLearner   *PatternLearner
    optimizer        *ContextOptimizer
    cache           *IntelligentCache
    metrics         *ContextMetrics
}

// IntelligentContextRequest represents a request for intelligent context selection
type IntelligentContextRequest struct {
    ProjectID       string                    `json:"project_id"`
    UserQuery       string                    `json:"user_query"`
    CurrentContext  []string                  `json:"current_context"`
    AvailableFiles  []FileMetadata           `json:"available_files"`
    UserHistory     []InteractionHistory     `json:"user_history"`
    TokenLimit      int                      `json:"token_limit"`
    Preferences     ContextPreferences       `json:"preferences"`
    TaskType        TaskType                 `json:"task_type"`
}

// IntelligentContextResponse represents the optimized context selection
type IntelligentContextResponse struct {
    SelectedFiles     []SelectedFile           `json:"selected_files"`
    RelevanceScores   map[string]float64       `json:"relevance_scores"`
    TokenEstimate     int                      `json:"token_estimate"`
    Confidence        float64                  `json:"confidence"`
    Reasoning         string                   `json:"reasoning"`
    Suggestions       []ContextSuggestion      `json:"suggestions"`
    LearningInsights  []LearningInsight        `json:"learning_insights"`
    OptimizationStats OptimizationStats        `json:"optimization_stats"`
}

// SelectedFile represents a file selected for context
type SelectedFile struct {
    Path            string                   `json:"path"`
    RelevanceScore  float64                  `json:"relevance_score"`
    TokenWeight     int                      `json:"token_weight"`
    InclusionReason string                   `json:"inclusion_reason"`
    ContentSummary  string                   `json:"content_summary"`
    Dependencies    []string                 `json:"dependencies"`
    Priority        Priority                 `json:"priority"`
}

// Context pattern learning
type PatternLearner struct {
    userPatterns    map[string]*UserPattern
    projectPatterns map[string]*ProjectPattern
    globalPatterns  *GlobalPatterns
    mlModel        *ContextMLModel
}

type UserPattern struct {
    UserID              string                      `json:"user_id"`
    PreferredFileTypes  map[string]float64          `json:"preferred_file_types"`
    CommonContextSizes  []int                       `json:"common_context_sizes"`
    FrequentPairs       map[string]map[string]float64 `json:"frequent_pairs"`
    TaskPreferences     map[TaskType]ContextProfile `json:"task_preferences"`
    SuccessfulContexts  []ContextSuccess            `json:"successful_contexts"`
    LastUpdated         time.Time                   `json:"last_updated"`
}

type ProjectPattern struct {
    ProjectID           string                      `json:"project_id"`
    ArchitectureType    string                      `json:"architecture_type"`
    CoreFiles          []string                    `json:"core_files"`
    ModuleDependencies  DependencyGraph             `json:"module_dependencies"`
    FileRelationships   RelationshipGraph           `json:"file_relationships"`
    OptimalContexts     []OptimalContext            `json:"optimal_contexts"`
    SemanticClusters    []SemanticCluster           `json:"semantic_clusters"`
}

// NewContextIntelligence creates a new intelligent context system
func NewContextIntelligence() *ContextIntelligence {
    return &ContextIntelligence{
        vectorizer:       NewContextVectorizer(),
        similarityEngine: NewSimilarityEngine(),
        patternLearner:   NewPatternLearner(),
        optimizer:        NewContextOptimizer(),
        cache:           NewIntelligentCache(),
        metrics:         NewContextMetrics(),
    }
}

// OptimizeContext performs intelligent context selection
func (ci *ContextIntelligence) OptimizeContext(ctx context.Context, request IntelligentContextRequest) (*IntelligentContextResponse, error) {
    start := time.Now()
    defer func() {
        ci.metrics.RecordOptimization(request.TaskType, time.Since(start))
    }()
    
    // Check cache for similar queries
    if cached := ci.cache.GetSimilar(request); cached != nil {
        ci.metrics.RecordCacheHit()
        return cached, nil
    }
    
    // Load user and project patterns
    userPattern := ci.patternLearner.GetUserPattern(request.ProjectID)
    projectPattern := ci.patternLearner.GetProjectPattern(request.ProjectID)
    
    // Vectorize query and files
    queryVector := ci.vectorizer.VectorizeQuery(request.UserQuery, request.TaskType)
    fileVectors := ci.vectorizer.VectorizeFiles(request.AvailableFiles)
    
    // Calculate similarity scores
    similarities := ci.similarityEngine.CalculateSimilarities(queryVector, fileVectors)
    
    // Apply pattern-based scoring
    patternScores := ci.applyPatternScoring(request, userPattern, projectPattern)
    
    // Combine scores with intelligent weighting
    finalScores := ci.combineScores(similarities, patternScores, request)
    
    // Select optimal context within token limit
    selection := ci.optimizer.SelectOptimalContext(finalScores, request.TokenLimit, request.Preferences)
    
    // Generate reasoning and insights
    reasoning := ci.generateReasoning(selection, request)
    insights := ci.generateLearningInsights(selection, userPattern, projectPattern)
    
    response := &IntelligentContextResponse{
        SelectedFiles:     selection.Files,
        RelevanceScores:   finalScores,
        TokenEstimate:     selection.TokenCount,
        Confidence:        selection.Confidence,
        Reasoning:         reasoning,
        Suggestions:       ci.generateSuggestions(selection, request),
        LearningInsights:  insights,
        OptimizationStats: selection.Stats,
    }
    
    // Cache result
    ci.cache.Store(request, response)
    
    // Learn from this interaction
    ci.patternLearner.LearnFromSelection(request, response)
    
    return response, nil
}

// LearnFromFeedback learns from user feedback on context quality
func (ci *ContextIntelligence) LearnFromFeedback(ctx context.Context, feedback ContextFeedback) error {
    // Update patterns based on user feedback
    if feedback.Helpful {
        ci.patternLearner.ReinforcePatternsSuccess(feedback.ContextID, feedback.UserID)
    } else {
        ci.patternLearner.ReducePatternWeights(feedback.ContextID, feedback.UserID)
    }
    
    // Update ML model with feedback
    return ci.patternLearner.mlModel.UpdateWithFeedback(feedback)
}

// Context vectorization for semantic similarity
type ContextVectorizer struct {
    tokenizer    *SemanticTokenizer
    embeddings   *EmbeddingModel
    vocabulary   *Vocabulary
    transformer  *ContextTransformer
}

func NewContextVectorizer() *ContextVectorizer {
    return &ContextVectorizer{
        tokenizer:   NewSemanticTokenizer(),
        embeddings:  NewEmbeddingModel(),
        vocabulary:  NewVocabulary(),
        transformer: NewContextTransformer(),
    }
}

func (cv *ContextVectorizer) VectorizeQuery(query string, taskType TaskType) *QueryVector {
    // Tokenize and extract features
    tokens := cv.tokenizer.Tokenize(query)
    
    // Extract task-specific features
    taskFeatures := cv.extractTaskFeatures(query, taskType)
    
    // Generate embeddings
    embeddings := cv.embeddings.Embed(tokens)
    
    // Combine with task context
    vector := cv.transformer.CombineFeatures(embeddings, taskFeatures)
    
    return &QueryVector{
        Vector:       vector,
        Tokens:       tokens,
        TaskType:     taskType,
        Features:     taskFeatures,
        Confidence:   cv.calculateConfidence(vector),
    }
}

func (cv *ContextVectorizer) VectorizeFiles(files []FileMetadata) map[string]*FileVector {
    vectors := make(map[string]*FileVector)
    
    for _, file := range files {
        // Extract file content features
        contentFeatures := cv.extractContentFeatures(file)
        
        // Extract structural features
        structuralFeatures := cv.extractStructuralFeatures(file)
        
        // Extract semantic features
        semanticFeatures := cv.extractSemanticFeatures(file)
        
        // Combine all features
        vector := cv.transformer.CombineFileFeatures(
            contentFeatures,
            structuralFeatures,
            semanticFeatures,
        )
        
        vectors[file.Path] = &FileVector{
            Vector:             vector,
            ContentFeatures:    contentFeatures,
            StructuralFeatures: structuralFeatures,
            SemanticFeatures:   semanticFeatures,
            FileType:          file.Language,
            Size:              file.Size,
            LastModified:      file.LastModified,
        }
    }
    
    return vectors
}

// Advanced similarity calculation
type SimilarityEngine struct {
    algorithms map[string]SimilarityAlgorithm
    weights    map[string]float64
    cache      *SimilarityCache
}

type SimilarityAlgorithm interface {
    Calculate(query *QueryVector, file *FileVector) float64
    Name() string
}

// Cosine similarity algorithm
type CosineSimilarity struct{}

func (cs *CosineSimilarity) Calculate(query *QueryVector, file *FileVector) float64 {
    return cosineSimilarity(query.Vector, file.Vector)
}

func (cs *CosineSimilarity) Name() string {
    return "cosine"
}

// Jaccard similarity for token overlap
type JaccardSimilarity struct{}

func (js *JaccardSimilarity) Calculate(query *QueryVector, file *FileVector) float64 {
    return jaccardSimilarity(query.Tokens, file.ContentFeatures.Tokens)
}

func (js *JaccardSimilarity) Name() string {
    return "jaccard"
}

// Semantic similarity using embeddings
type SemanticSimilarity struct {
    embeddings *EmbeddingModel
}

func (ss *SemanticSimilarity) Calculate(query *QueryVector, file *FileVector) float64 {
    queryEmbedding := ss.embeddings.EmbedText(strings.Join(query.Tokens, " "))
    fileEmbedding := ss.embeddings.EmbedText(file.ContentFeatures.Summary)
    
    return cosineSimilarity(queryEmbedding, fileEmbedding)
}

func (ss *SemanticSimilarity) Name() string {
    return "semantic"
}

func NewSimilarityEngine() *SimilarityEngine {
    return &SimilarityEngine{
        algorithms: map[string]SimilarityAlgorithm{
            "cosine":   &CosineSimilarity{},
            "jaccard":  &JaccardSimilarity{},
            "semantic": &SemanticSimilarity{embeddings: NewEmbeddingModel()},
        },
        weights: map[string]float64{
            "cosine":   0.4,
            "jaccard":  0.3,
            "semantic": 0.3,
        },
        cache: NewSimilarityCache(),
    }
}

func (se *SimilarityEngine) CalculateSimilarities(query *QueryVector, files map[string]*FileVector) map[string]float64 {
    similarities := make(map[string]float64)
    
    for path, fileVector := range files {
        // Check cache first
        if cached := se.cache.Get(query, fileVector); cached != nil {
            similarities[path] = *cached
            continue
        }
        
        var weightedSum float64
        var totalWeight float64
        
        // Calculate similarity using multiple algorithms
        for name, algorithm := range se.algorithms {
            weight := se.weights[name]
            similarity := algorithm.Calculate(query, fileVector)
            
            weightedSum += similarity * weight
            totalWeight += weight
        }
        
        finalSimilarity := weightedSum / totalWeight
        similarities[path] = finalSimilarity
        
        // Cache result
        se.cache.Store(query, fileVector, finalSimilarity)
    }
    
    return similarities
}

// Context optimization with constraints
type ContextOptimizer struct {
    strategies map[string]OptimizationStrategy
    selector   *FileSelector
    estimator  *TokenEstimator
}

type OptimizationStrategy interface {
    Optimize(scores map[string]float64, limit int, preferences ContextPreferences) *OptimizationResult
    Name() string
}

// Greedy optimization strategy
type GreedyOptimization struct{}

func (go *GreedyOptimization) Optimize(scores map[string]float64, limit int, preferences ContextPreferences) *OptimizationResult {
    // Sort files by score
    type fileScore struct {
        path  string
        score float64
    }
    
    var sorted []fileScore
    for path, score := range scores {
        sorted = append(sorted, fileScore{path, score})
    }
    
    sort.Slice(sorted, func(i, j int) bool {
        return sorted[i].score > sorted[j].score
    })
    
    // Select files within token limit
    var selected []SelectedFile
    totalTokens := 0
    
    for _, fs := range sorted {
        // Estimate tokens for this file
        tokens := estimateFileTokens(fs.path)
        
        if totalTokens+tokens <= limit {
            selected = append(selected, SelectedFile{
                Path:           fs.path,
                RelevanceScore: fs.score,
                TokenWeight:    tokens,
                InclusionReason: "High relevance score",
                Priority:       getPriority(fs.score),
            })
            totalTokens += tokens
        }
        
        if len(selected) >= preferences.MaxFiles {
            break
        }
    }
    
    return &OptimizationResult{
        Files:      selected,
        TokenCount: totalTokens,
        Confidence: calculateSelectionConfidence(selected),
        Strategy:   "greedy",
        Stats: OptimizationStats{
            FilesConsidered: len(sorted),
            FilesSelected:   len(selected),
            TokenUtilization: float64(totalTokens) / float64(limit),
        },
    }
}

// Machine learning for context optimization
type ContextMLModel struct {
    model          *NeuralNetwork
    featureExtractor *FeatureExtractor
    trainer        *ModelTrainer
    evaluator      *ModelEvaluator
}

func NewContextMLModel() *ContextMLModel {
    return &ContextMLModel{
        model:           NewNeuralNetwork(),
        featureExtractor: NewFeatureExtractor(),
        trainer:         NewModelTrainer(),
        evaluator:       NewModelEvaluator(),
    }
}

func (ml *ContextMLModel) PredictRelevance(query *QueryVector, file *FileVector, context ContextFeatures) float64 {
    // Extract features for ML model
    features := ml.featureExtractor.ExtractFeatures(query, file, context)
    
    // Predict relevance score
    prediction := ml.model.Predict(features)
    
    return prediction
}

func (ml *ContextMLModel) Train(trainingData []TrainingExample) error {
    // Prepare training data
    features, labels := ml.prepareTrainingData(trainingData)
    
    // Train the model
    return ml.trainer.Train(ml.model, features, labels)
}

func (ml *ContextMLModel) UpdateWithFeedback(feedback ContextFeedback) error {
    // Convert feedback to training example
    example := ml.feedbackToTrainingExample(feedback)
    
    // Perform online learning update
    return ml.trainer.OnlineUpdate(ml.model, example)
}

// Advanced caching with semantic similarity
type IntelligentCache struct {
    entries    map[string]*CacheEntry
    index      *SemanticIndex
    maxSize    int
    ttl        time.Duration
}

type CacheEntry struct {
    Request    IntelligentContextRequest     `json:"request"`
    Response   *IntelligentContextResponse   `json:"response"`
    AccessTime time.Time                     `json:"access_time"`
    HitCount   int                          `json:"hit_count"`
    Embedding  []float64                    `json:"embedding"`
}

func (ic *IntelligentCache) GetSimilar(request IntelligentContextRequest) *IntelligentContextResponse {
    // Create embedding for the request
    embedding := ic.createRequestEmbedding(request)
    
    // Find most similar cached request
    similar := ic.index.FindMostSimilar(embedding, 0.9) // 90% similarity threshold
    
    if similar != nil {
        entry := ic.entries[similar.ID]
        entry.AccessTime = time.Now()
        entry.HitCount++
        
        return entry.Response
    }
    
    return nil
}

// Utility functions
func cosineSimilarity(a, b []float64) float64 {
    if len(a) != len(b) {
        return 0
    }
    
    var dotProduct, normA, normB float64
    for i := range a {
        dotProduct += a[i] * b[i]
        normA += a[i] * a[i]
        normB += b[i] * b[i]
    }
    
    if normA == 0 || normB == 0 {
        return 0
    }
    
    return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}

func jaccardSimilarity(a, b []string) float64 {
    setA := make(map[string]bool)
    setB := make(map[string]bool)
    
    for _, item := range a {
        setA[item] = true
    }
    for _, item := range b {
        setB[item] = true
    }
    
    intersection := 0
    union := len(setA)
    
    for item := range setB {
        if setA[item] {
            intersection++
        } else {
            union++
        }
    }
    
    if union == 0 {
        return 0
    }
    
    return float64(intersection) / float64(union)
}
```

**TodoWrite Task**: `Implement ML-powered intelligent context optimization`

### KPIs for Phase 5A
- ✅ Multi-modal AI processing with 95%+ accuracy
- ✅ Context optimization reducing token usage by 40%+
- ✅ Real-time vision analysis <5 seconds
- ✅ Intelligent pattern learning with user feedback
- ✅ Advanced semantic understanding and reasoning
- ✅ Seamless integration with existing AI workflows

---

## Phase 5B: Progressive Web App Development
### 📱 MOBILE-FIRST EXPERIENCE TARGET
**Goal**: Create offline-capable, mobile-responsive PWA with native app performance

### Implementation Steps

#### Step 5B.1: PWA Architecture Setup
**File: `/app/web-dashboard/public/manifest.json`** (create PWA manifest)
```json
{
  "name": "Plandex - AI Development Assistant",
  "short_name": "Plandex",
  "description": "AI-powered development assistant for large-scale coding projects",
  "start_url": "/",
  "display": "standalone",
  "orientation": "portrait-primary",
  "theme_color": "#4F46E5",
  "background_color": "#FFFFFF",
  "lang": "en",
  "scope": "/",
  "categories": ["productivity", "developer", "business"],
  "icons": [
    {
      "src": "/icons/icon-72x72.png",
      "sizes": "72x72",
      "type": "image/png",
      "purpose": "maskable any"
    },
    {
      "src": "/icons/icon-96x96.png",
      "sizes": "96x96",
      "type": "image/png",
      "purpose": "maskable any"
    },
    {
      "src": "/icons/icon-128x128.png",
      "sizes": "128x128",
      "type": "image/png",
      "purpose": "maskable any"
    },
    {
      "src": "/icons/icon-144x144.png",
      "sizes": "144x144",
      "type": "image/png",
      "purpose": "maskable any"
    },
    {
      "src": "/icons/icon-152x152.png",
      "sizes": "152x152",
      "type": "image/png",
      "purpose": "maskable any"
    },
    {
      "src": "/icons/icon-192x192.png",
      "sizes": "192x192",
      "type": "image/png",
      "purpose": "maskable any"
    },
    {
      "src": "/icons/icon-384x384.png",
      "sizes": "384x384",
      "type": "image/png",
      "purpose": "maskable any"
    },
    {
      "src": "/icons/icon-512x512.png",
      "sizes": "512x512",
      "type": "image/png",
      "purpose": "maskable any"
    }
  ],
  "screenshots": [
    {
      "src": "/screenshots/desktop-1.png",
      "sizes": "1920x1080",
      "type": "image/png",
      "form_factor": "wide",
      "label": "Plans Dashboard"
    },
    {
      "src": "/screenshots/mobile-1.png",
      "sizes": "390x844",
      "type": "image/png",
      "form_factor": "narrow",
      "label": "Mobile Plans View"
    }
  ],
  "shortcuts": [
    {
      "name": "New Plan",
      "short_name": "New Plan",
      "description": "Create a new AI development plan",
      "url": "/plans/new",
      "icons": [
        {
          "src": "/icons/shortcut-new.png",
          "sizes": "96x96",
          "type": "image/png"
        }
      ]
    },
    {
      "name": "AI Chat",
      "short_name": "Chat",
      "description": "Start AI conversation",
      "url": "/chat",
      "icons": [
        {
          "src": "/icons/shortcut-chat.png",
          "sizes": "96x96",
          "type": "image/png"
        }
      ]
    }
  ],
  "prefer_related_applications": false,
  "related_applications": [],
  "file_handlers": [
    {
      "action": "/handle-file",
      "accept": {
        "text/plain": [".txt", ".md"],
        "text/javascript": [".js", ".mjs"],
        "text/typescript": [".ts"],
        "text/python": [".py"],
        "text/go": [".go"]
      }
    }
  ],
  "protocol_handlers": [
    {
      "protocol": "plandex",
      "url": "/handle-protocol?url=%s"
    }
  ],
  "launch_handler": {
    "client_mode": ["navigate-existing", "navigate-new"]
  },
  "edge_side_panel": {
    "preferred_width": 400
  }
}
```

**File: `/app/web-dashboard/src/hooks/usePWA.ts`** (create PWA hooks)
```typescript
import { useState, useEffect, useCallback } from 'react';

// PWA Installation Hook
export const usePWAInstall = () => {
  const [deferredPrompt, setDeferredPrompt] = useState<any>(null);
  const [isInstallable, setIsInstallable] = useState(false);
  const [isInstalled, setIsInstalled] = useState(false);

  useEffect(() => {
    // Check if already installed
    const checkInstalled = () => {
      const isStandalone = window.matchMedia('(display-mode: standalone)').matches;
      const isInApp = (navigator as any).standalone === true;
      setIsInstalled(isStandalone || isInApp);
    };

    checkInstalled();

    // Listen for install prompt
    const handleBeforeInstallPrompt = (e: Event) => {
      e.preventDefault();
      setDeferredPrompt(e);
      setIsInstallable(true);
    };

    // Listen for app installed
    const handleAppInstalled = () => {
      setIsInstalled(true);
      setIsInstallable(false);
      setDeferredPrompt(null);
    };

    window.addEventListener('beforeinstallprompt', handleBeforeInstallPrompt);
    window.addEventListener('appinstalled', handleAppInstalled);

    return () => {
      window.removeEventListener('beforeinstallprompt', handleBeforeInstallPrompt);
      window.removeEventListener('appinstalled', handleAppInstalled);
    };
  }, []);

  const install = useCallback(async () => {
    if (!deferredPrompt) return false;

    try {
      deferredPrompt.prompt();
      const choiceResult = await deferredPrompt.userChoice;
      
      if (choiceResult.outcome === 'accepted') {
        setIsInstallable(false);
        setDeferredPrompt(null);
        return true;
      }
      return false;
    } catch (error) {
      console.error('PWA install failed:', error);
      return false;
    }
  }, [deferredPrompt]);

  return {
    isInstallable,
    isInstalled,
    install
  };
};

// Service Worker Hook
export const useServiceWorker = () => {
  const [isOnline, setIsOnline] = useState(navigator.onLine);
  const [updateAvailable, setUpdateAvailable] = useState(false);
  const [registration, setRegistration] = useState<ServiceWorkerRegistration | null>(null);

  useEffect(() => {
    // Register service worker
    if ('serviceWorker' in navigator) {
      navigator.serviceWorker.register('/sw.js')
        .then((reg) => {
          setRegistration(reg);
          
          // Check for updates
          reg.addEventListener('updatefound', () => {
            const newWorker = reg.installing;
            if (newWorker) {
              newWorker.addEventListener('statechange', () => {
                if (newWorker.state === 'installed' && navigator.serviceWorker.controller) {
                  setUpdateAvailable(true);
                }
              });
            }
          });
        })
        .catch((error) => {
          console.error('Service Worker registration failed:', error);
        });
    }

    // Listen for online/offline changes
    const handleOnline = () => setIsOnline(true);
    const handleOffline = () => setIsOnline(false);

    window.addEventListener('online', handleOnline);
    window.addEventListener('offline', handleOffline);

    return () => {
      window.removeEventListener('online', handleOnline);
      window.removeEventListener('offline', handleOffline);
    };
  }, []);

  const updateApp = useCallback(() => {
    if (registration && registration.waiting) {
      registration.waiting.postMessage({ type: 'SKIP_WAITING' });
      window.location.reload();
    }
  }, [registration]);

  return {
    isOnline,
    updateAvailable,
    updateApp
  };
};

// Offline Storage Hook
export const useOfflineStorage = <T>(key: string, defaultValue: T) => {
  const [value, setValue] = useState<T>(() => {
    try {
      const item = localStorage.getItem(key);
      return item ? JSON.parse(item) : defaultValue;
    } catch {
      return defaultValue;
    }
  });

  const setStoredValue = useCallback((newValue: T | ((prev: T) => T)) => {
    try {
      const valueToStore = typeof newValue === 'function' 
        ? (newValue as (prev: T) => T)(value)
        : newValue;
      
      setValue(valueToStore);
      localStorage.setItem(key, JSON.stringify(valueToStore));
    } catch (error) {
      console.error('Failed to store value:', error);
    }
  }, [key, value]);

  return [value, setStoredValue] as const;
};

// Background Sync Hook
export const useBackgroundSync = () => {
  const [pendingOperations, setPendingOperations] = useState<any[]>([]);

  const addPendingOperation = useCallback((operation: any) => {
    setPendingOperations(prev => [...prev, operation]);
    
    // Store in IndexedDB for persistence
    if ('serviceWorker' in navigator && 'sync' in window.ServiceWorkerRegistration.prototype) {
      navigator.serviceWorker.ready.then((registration) => {
        return registration.sync.register('background-sync');
      });
    }
  }, []);

  const clearPendingOperations = useCallback(() => {
    setPendingOperations([]);
  }, []);

  return {
    pendingOperations,
    addPendingOperation,
    clearPendingOperations
  };
};

// Push Notifications Hook
export const usePushNotifications = () => {
  const [permission, setPermission] = useState<NotificationPermission>('default');
  const [subscription, setSubscription] = useState<PushSubscription | null>(null);

  useEffect(() => {
    if ('Notification' in window) {
      setPermission(Notification.permission);
    }
  }, []);

  const requestPermission = useCallback(async () => {
    if (!('Notification' in window)) {
      throw new Error('This browser does not support notifications');
    }

    const permission = await Notification.requestPermission();
    setPermission(permission);
    return permission === 'granted';
  }, []);

  const subscribe = useCallback(async () => {
    if (!('serviceWorker' in navigator) || !('PushManager' in window)) {
      throw new Error('Push notifications not supported');
    }

    const registration = await navigator.serviceWorker.ready;
    const subscription = await registration.pushManager.subscribe({
      userVisibleOnly: true,
      applicationServerKey: process.env.REACT_APP_VAPID_PUBLIC_KEY
    });

    setSubscription(subscription);
    
    // Send subscription to server
    await fetch('/api/push/subscribe', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(subscription),
    });

    return subscription;
  }, []);

  const unsubscribe = useCallback(async () => {
    if (subscription) {
      await subscription.unsubscribe();
      setSubscription(null);
      
      // Remove subscription from server
      await fetch('/api/push/unsubscribe', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(subscription),
      });
    }
  }, [subscription]);

  return {
    permission,
    subscription,
    requestPermission,
    subscribe,
    unsubscribe
  };
};

// Network Status Hook
export const useNetworkStatus = () => {
  const [isOnline, setIsOnline] = useState(navigator.onLine);
  const [connectionType, setConnectionType] = useState<string>('unknown');
  const [effectiveType, setEffectiveType] = useState<string>('unknown');

  useEffect(() => {
    const updateOnlineStatus = () => setIsOnline(navigator.onLine);
    
    const updateConnectionInfo = () => {
      const connection = (navigator as any).connection || 
                        (navigator as any).mozConnection || 
                        (navigator as any).webkitConnection;
      
      if (connection) {
        setConnectionType(connection.type || 'unknown');
        setEffectiveType(connection.effectiveType || 'unknown');
      }
    };

    updateConnectionInfo();

    window.addEventListener('online', updateOnlineStatus);
    window.addEventListener('offline', updateOnlineStatus);

    const connection = (navigator as any).connection;
    if (connection) {
      connection.addEventListener('change', updateConnectionInfo);
    }

    return () => {
      window.removeEventListener('online', updateOnlineStatus);
      window.removeEventListener('offline', updateOnlineStatus);
      
      if (connection) {
        connection.removeEventListener('change', updateConnectionInfo);
      }
    };
  }, []);

  return {
    isOnline,
    connectionType,
    effectiveType,
    isSlowConnection: effectiveType === 'slow-2g' || effectiveType === '2g'
  };
};
```

**File: `/app/web-dashboard/public/sw.js`** (create service worker)
```javascript
// Plandex PWA Service Worker
const CACHE_NAME = 'plandex-v1.0.0';
const API_CACHE_NAME = 'plandex-api-v1.0.0';
const RUNTIME_CACHE_NAME = 'plandex-runtime-v1.0.0';

// Files to cache immediately
const PRECACHE_URLS = [
  '/',
  '/static/js/bundle.js',
  '/static/css/main.css',
  '/manifest.json',
  '/icons/icon-192x192.png',
  '/icons/icon-512x512.png'
];

// API endpoints to cache
const API_CACHE_PATTERNS = [
  /^https:\/\/api\.plandex\.com\/api\/plans/,
  /^https:\/\/api\.plandex\.com\/api\/user/,
  /^https:\/\/api\.plandex\.com\/api\/organizations/
];

// Runtime caching patterns
const RUNTIME_CACHE_PATTERNS = [
  {
    pattern: /^https:\/\/fonts\.googleapis\.com/,
    strategy: 'StaleWhileRevalidate',
    cacheName: 'google-fonts-stylesheets'
  },
  {
    pattern: /^https:\/\/fonts\.gstatic\.com/,
    strategy: 'CacheFirst',
    cacheName: 'google-fonts-webfonts'
  }
];

// Install event - cache core files
self.addEventListener('install', (event) => {
  event.waitUntil(
    caches.open(CACHE_NAME)
      .then((cache) => {
        console.log('Precaching core files');
        return cache.addAll(PRECACHE_URLS);
      })
      .then(() => {
        console.log('Service worker installed');
        return self.skipWaiting();
      })
  );
});

// Activate event - clean up old caches
self.addEventListener('activate', (event) => {
  event.waitUntil(
    caches.keys()
      .then((cacheNames) => {
        return Promise.all(
          cacheNames.map((cacheName) => {
            if (cacheName !== CACHE_NAME && 
                cacheName !== API_CACHE_NAME && 
                cacheName !== RUNTIME_CACHE_NAME) {
              console.log('Deleting old cache:', cacheName);
              return caches.delete(cacheName);
            }
          })
        );
      })
      .then(() => {
        console.log('Service worker activated');
        return self.clients.claim();
      })
  );
});

// Fetch event - network-first with fallback strategies
self.addEventListener('fetch', (event) => {
  const { request } = event;
  const url = new URL(request.url);

  // Handle different types of requests
  if (request.method === 'GET') {
    // API requests - network first with cache fallback
    if (isAPIRequest(request)) {
      event.respondWith(networkFirstStrategy(request, API_CACHE_NAME));
    }
    // Navigation requests - cache first with network fallback
    else if (request.mode === 'navigate') {
      event.respondWith(navigationStrategy(request));
    }
    // Static assets - cache first
    else if (isStaticAsset(request)) {
      event.respondWith(cacheFirstStrategy(request, CACHE_NAME));
    }
    // Runtime caching for external resources
    else {
      const cachePattern = RUNTIME_CACHE_PATTERNS.find(
        pattern => pattern.pattern.test(request.url)
      );
      
      if (cachePattern) {
        if (cachePattern.strategy === 'CacheFirst') {
          event.respondWith(cacheFirstStrategy(request, cachePattern.cacheName));
        } else if (cachePattern.strategy === 'StaleWhileRevalidate') {
          event.respondWith(staleWhileRevalidateStrategy(request, cachePattern.cacheName));
        }
      }
    }
  }
  // Handle POST requests for offline functionality
  else if (request.method === 'POST') {
    if (isAPIRequest(request)) {
      event.respondWith(handleOfflinePost(request));
    }
  }
});

// Background sync for offline operations
self.addEventListener('sync', (event) => {
  if (event.tag === 'background-sync') {
    event.waitUntil(
      processOfflineOperations()
    );
  }
});

// Push notification handling
self.addEventListener('push', (event) => {
  const options = {
    body: event.data ? event.data.text() : 'New notification from Plandex',
    icon: '/icons/icon-192x192.png',
    badge: '/icons/badge-72x72.png',
    vibrate: [100, 50, 100],
    data: {
      dateOfArrival: Date.now(),
      primaryKey: 1
    },
    actions: [
      {
        action: 'explore',
        title: 'View',
        icon: '/icons/view.png'
      },
      {
        action: 'close',
        title: 'Close',
        icon: '/icons/close.png'
      }
    ]
  };

  event.waitUntil(
    self.registration.showNotification('Plandex', options)
  );
});

// Notification click handling
self.addEventListener('notificationclick', (event) => {
  event.notification.close();

  if (event.action === 'explore') {
    event.waitUntil(
      clients.openWindow('/')
    );
  }
});

// Strategy implementations
async function networkFirstStrategy(request, cacheName) {
  try {
    // Try network first
    const networkResponse = await fetch(request);
    
    // Cache successful responses
    if (networkResponse.ok) {
      const cache = await caches.open(cacheName);
      cache.put(request.clone(), networkResponse.clone());
    }
    
    return networkResponse;
  } catch (error) {
    // Fallback to cache
    console.log('Network failed, falling back to cache');
    const cachedResponse = await caches.match(request);
    
    if (cachedResponse) {
      return cachedResponse;
    }
    
    // Return offline page for navigation requests
    if (request.mode === 'navigate') {
      return caches.match('/offline.html');
    }
    
    throw error;
  }
}

async function cacheFirstStrategy(request, cacheName) {
  const cachedResponse = await caches.match(request);
  
  if (cachedResponse) {
    return cachedResponse;
  }
  
  try {
    const networkResponse = await fetch(request);
    const cache = await caches.open(cacheName);
    cache.put(request.clone(), networkResponse.clone());
    return networkResponse;
  } catch (error) {
    console.log('Cache and network failed for:', request.url);
    throw error;
  }
}

async function staleWhileRevalidateStrategy(request, cacheName) {
  const cache = await caches.open(cacheName);
  const cachedResponse = await cache.match(request);
  
  const fetchPromise = fetch(request).then((networkResponse) => {
    cache.put(request.clone(), networkResponse.clone());
    return networkResponse;
  });
  
  return cachedResponse || fetchPromise;
}

async function navigationStrategy(request) {
  try {
    // Try network first for navigation
    const networkResponse = await fetch(request);
    
    // Cache the response
    const cache = await caches.open(CACHE_NAME);
    cache.put(request.clone(), networkResponse.clone());
    
    return networkResponse;
  } catch (error) {
    // Fallback to cached page or offline page
    const cachedResponse = await caches.match(request);
    return cachedResponse || caches.match('/offline.html');
  }
}

async function handleOfflinePost(request) {
  try {
    return await fetch(request);
  } catch (error) {
    // Store for background sync
    const data = await request.clone().json();
    await storeOfflineOperation({
      url: request.url,
      method: request.method,
      data: data,
      timestamp: Date.now()
    });
    
    return new Response(
      JSON.stringify({ message: 'Operation queued for when online' }),
      {
        status: 202,
        headers: { 'Content-Type': 'application/json' }
      }
    );
  }
}

async function processOfflineOperations() {
  const operations = await getOfflineOperations();
  
  for (const operation of operations) {
    try {
      const response = await fetch(operation.url, {
        method: operation.method,
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(operation.data)
      });
      
      if (response.ok) {
        await removeOfflineOperation(operation.id);
      }
    } catch (error) {
      console.log('Failed to sync operation:', operation.id, error);
    }
  }
}

// Helper functions
function isAPIRequest(request) {
  return API_CACHE_PATTERNS.some(pattern => pattern.test(request.url));
}

function isStaticAsset(request) {
  const url = new URL(request.url);
  return url.pathname.startsWith('/static/') ||
         url.pathname.endsWith('.js') ||
         url.pathname.endsWith('.css') ||
         url.pathname.endsWith('.png') ||
         url.pathname.endsWith('.jpg') ||
         url.pathname.endsWith('.svg');
}

// IndexedDB operations for offline storage
async function storeOfflineOperation(operation) {
  return new Promise((resolve, reject) => {
    const request = indexedDB.open('PlandexOffline', 1);
    
    request.onerror = () => reject(request.error);
    request.onsuccess = () => {
      const db = request.result;
      const transaction = db.transaction(['operations'], 'readwrite');
      const store = transaction.objectStore('operations');
      
      operation.id = Date.now() + Math.random();
      store.add(operation);
      
      transaction.oncomplete = () => resolve();
      transaction.onerror = () => reject(transaction.error);
    };
    
    request.onupgradeneeded = () => {
      const db = request.result;
      if (!db.objectStoreNames.contains('operations')) {
        db.createObjectStore('operations', { keyPath: 'id' });
      }
    };
  });
}

async function getOfflineOperations() {
  return new Promise((resolve, reject) => {
    const request = indexedDB.open('PlandexOffline', 1);
    
    request.onerror = () => reject(request.error);
    request.onsuccess = () => {
      const db = request.result;
      const transaction = db.transaction(['operations'], 'readonly');
      const store = transaction.objectStore('operations');
      const getAllRequest = store.getAll();
      
      getAllRequest.onsuccess = () => resolve(getAllRequest.result);
      getAllRequest.onerror = () => reject(getAllRequest.error);
    };
  });
}

async function removeOfflineOperation(id) {
  return new Promise((resolve, reject) => {
    const request = indexedDB.open('PlandexOffline', 1);
    
    request.onerror = () => reject(request.error);
    request.onsuccess = () => {
      const db = request.result;
      const transaction = db.transaction(['operations'], 'readwrite');
      const store = transaction.objectStore('operations');
      
      store.delete(id);
      
      transaction.oncomplete = () => resolve();
      transaction.onerror = () => reject(transaction.error);
    };
  });
}

// Handle messages from main thread
self.addEventListener('message', (event) => {
  if (event.data && event.data.type === 'SKIP_WAITING') {
    self.skipWaiting();
  }
});
```

**TodoWrite Task**: `Implement comprehensive PWA architecture with offline capabilities`

### KPIs for Phase 5B
- ✅ PWA with offline functionality and app-like experience
- ✅ Mobile-responsive design with touch optimization
- ✅ <3 second load time on slow connections
- ✅ 90%+ functionality available offline
- ✅ Push notifications for real-time updates
- ✅ Native app installation and shortcuts

---

## Phase 5C: Enterprise Security & Compliance
### 🛡️ ENTERPRISE-GRADE SECURITY TARGET
**Goal**: Advanced RBAC, audit logging, compliance frameworks, and enterprise security features

### Implementation Steps

#### Step 5C.1: Advanced RBAC System
**File: `/app/server/auth/rbac_system.go`** (create advanced RBAC)
```go
package auth

import (
    "context"
    "encoding/json"
    "fmt"
    "strings"
    "time"
    
    "github.com/casbin/casbin/v2"
    "github.com/casbin/casbin/v2/model"
    "github.com/casbin/casbin/v2/persist"
)

// AdvancedRBACSystem provides enterprise-grade role-based access control
type AdvancedRBACSystem struct {
    enforcer        *casbin.Enforcer
    policyManager   *PolicyManager
    roleManager     *RoleManager
    permissionCache *PermissionCache
    auditLogger     *AuditLogger
    validator       *PolicyValidator
}

// Permission represents a fine-grained permission
type Permission struct {
    ID          string            `json:"id"`
    Name        string            `json:"name"`
    Resource    string            `json:"resource"`
    Action      string            `json:"action"`
    Conditions  []Condition       `json:"conditions,omitempty"`
    Metadata    map[string]string `json:"metadata,omitempty"`
    CreatedAt   time.Time         `json:"created_at"`
    UpdatedAt   time.Time         `json:"updated_at"`
}

// Role represents a collection of permissions
type Role struct {
    ID              string            `json:"id"`
    Name            string            `json:"name"`
    Description     string            `json:"description"`
    OrganizationID  string            `json:"organization_id"`
    Permissions     []Permission      `json:"permissions"`
    Inherits        []string          `json:"inherits,omitempty"`
    IsSystem        bool              `json:"is_system"`
    IsActive        bool              `json:"is_active"`
    Metadata        map[string]string `json:"metadata,omitempty"`
    CreatedAt       time.Time         `json:"created_at"`
    UpdatedAt       time.Time         `json:"updated_at"`
    CreatedBy       string            `json:"created_by"`
}

// Condition represents a permission condition
type Condition struct {
    Type     ConditionType     `json:"type"`
    Field    string            `json:"field"`
    Operator ConditionOperator `json:"operator"`
    Value    interface{}       `json:"value"`
    Logic    LogicOperator     `json:"logic,omitempty"`
}

// Resource hierarchy and permissions
type Resource struct {
    Type        ResourceType      `json:"type"`
    ID          string            `json:"id"`
    Parent      *Resource         `json:"parent,omitempty"`
    Children    []Resource        `json:"children,omitempty"`
    Attributes  map[string]string `json:"attributes,omitempty"`
    Owners      []string          `json:"owners,omitempty"`
}

// Enums
type ResourceType string
const (
    ResourceTypeOrganization ResourceType = "organization"
    ResourceTypeProject      ResourceType = "project"
    ResourceTypePlan         ResourceType = "plan"
    ResourceTypeFile         ResourceType = "file"
    ResourceTypeContext      ResourceType = "context"
    ResourceTypeConversation ResourceType = "conversation"
    ResourceTypeModel        ResourceType = "model"
    ResourceTypeUser         ResourceType = "user"
    ResourceTypeRole         ResourceType = "role"
    ResourceTypeAuditLog     ResourceType = "audit_log"
)

type ActionType string
const (
    ActionCreate ActionType = "create"
    ActionRead   ActionType = "read"
    ActionUpdate ActionType = "update"
    ActionDelete ActionType = "delete"
    ActionExecute ActionType = "execute"
    ActionShare  ActionType = "share"
    ActionManage ActionType = "manage"
    ActionAudit  ActionType = "audit"
)

type ConditionType string
const (
    ConditionTypeTime        ConditionType = "time"
    ConditionTypeLocation    ConditionType = "location"
    ConditionTypeDevice      ConditionType = "device"
    ConditionTypeAttribute   ConditionType = "attribute"
    ConditionTypeOwnership   ConditionType = "ownership"
    ConditionTypeHierarchy   ConditionType = "hierarchy"
)

type ConditionOperator string
const (
    OperatorEquals        ConditionOperator = "eq"
    OperatorNotEquals     ConditionOperator = "ne"
    OperatorGreaterThan   ConditionOperator = "gt"
    OperatorLessThan      ConditionOperator = "lt"
    OperatorIn            ConditionOperator = "in"
    OperatorNotIn         ConditionOperator = "not_in"
    OperatorContains      ConditionOperator = "contains"
    OperatorStartsWith    ConditionOperator = "starts_with"
    OperatorEndsWith      ConditionOperator = "ends_with"
    OperatorMatches       ConditionOperator = "matches"
)

type LogicOperator string
const (
    LogicAnd LogicOperator = "and"
    LogicOr  LogicOperator = "or"
    LogicNot LogicOperator = "not"
)

// NewAdvancedRBACSystem creates a new RBAC system
func NewAdvancedRBACSystem(adapter persist.Adapter) (*AdvancedRBACSystem, error) {
    // Load Casbin model
    modelText := `
[request_definition]
r = sub, obj, act, ctx

[policy_definition]
p = sub, obj, act, eft, conditions

[role_definition]
g = _, _
g2 = _, _

[policy_effect]
e = some(where (p.eft == allow)) && !some(where (p.eft == deny))

[matchers]
m = g(r.sub, p.sub) && keyMatch3(r.obj, p.obj) && regexMatch(r.act, p.act) && eval(p.conditions)
`
    
    model, err := model.NewModelFromString(modelText)
    if err != nil {
        return nil, fmt.Errorf("failed to create RBAC model: %w", err)
    }
    
    enforcer, err := casbin.NewEnforcer(model, adapter)
    if err != nil {
        return nil, fmt.Errorf("failed to create RBAC enforcer: %w", err)
    }
    
    return &AdvancedRBACSystem{
        enforcer:        enforcer,
        policyManager:   NewPolicyManager(enforcer),
        roleManager:     NewRoleManager(),
        permissionCache: NewPermissionCache(),
        auditLogger:     NewAuditLogger(),
        validator:       NewPolicyValidator(),
    }, nil
}

// CheckPermission validates if a user has permission for an action on a resource
func (rbac *AdvancedRBACSystem) CheckPermission(ctx context.Context, userID, resource, action string, context map[string]interface{}) (bool, error) {
    start := time.Now()
    defer func() {
        rbac.auditLogger.LogPermissionCheck(userID, resource, action, time.Since(start))
    }()
    
    // Check cache first
    cacheKey := fmt.Sprintf("%s:%s:%s", userID, resource, action)
    if cached := rbac.permissionCache.Get(cacheKey); cached != nil {
        return *cached, nil
    }
    
    // Prepare context for evaluation
    contextStr := rbac.serializeContext(context)
    
    // Check permission using Casbin
    allowed, err := rbac.enforcer.Enforce(userID, resource, action, contextStr)
    if err != nil {
        return false, fmt.Errorf("permission check failed: %w", err)
    }
    
    // Cache result
    rbac.permissionCache.Set(cacheKey, allowed, 5*time.Minute)
    
    // Log the permission check
    rbac.auditLogger.LogPermissionCheck(userID, resource, action, time.Since(start))
    
    return allowed, nil
}

// GrantPermission grants a permission to a user
func (rbac *AdvancedRBACSystem) GrantPermission(ctx context.Context, userID, resource, action string, conditions []Condition, grantedBy string) error {
    // Validate the permission grant
    if err := rbac.validator.ValidatePermissionGrant(userID, resource, action, conditions); err != nil {
        return fmt.Errorf("invalid permission grant: %w", err)
    }
    
    // Serialize conditions
    conditionsStr := rbac.serializeConditions(conditions)
    
    // Add policy
    _, err := rbac.enforcer.AddPolicy(userID, resource, action, "allow", conditionsStr)
    if err != nil {
        return fmt.Errorf("failed to grant permission: %w", err)
    }
    
    // Invalidate cache
    rbac.permissionCache.Invalidate(userID)
    
    // Log the grant
    rbac.auditLogger.LogPermissionGrant(userID, resource, action, grantedBy)
    
    return nil
}

// RevokePermission revokes a permission from a user
func (rbac *AdvancedRBACSystem) RevokePermission(ctx context.Context, userID, resource, action string, revokedBy string) error {
    // Remove policy
    _, err := rbac.enforcer.RemovePolicy(userID, resource, action)
    if err != nil {
        return fmt.Errorf("failed to revoke permission: %w", err)
    }
    
    // Invalidate cache
    rbac.permissionCache.Invalidate(userID)
    
    // Log the revocation
    rbac.auditLogger.LogPermissionRevoke(userID, resource, action, revokedBy)
    
    return nil
}

// AssignRole assigns a role to a user
func (rbac *AdvancedRBACSystem) AssignRole(ctx context.Context, userID, roleID string, assignedBy string) error {
    // Validate role exists
    role, err := rbac.roleManager.GetRole(roleID)
    if err != nil {
        return fmt.Errorf("role not found: %w", err)
    }
    
    // Check if assigner has permission to assign this role
    canAssign, err := rbac.CheckPermission(ctx, assignedBy, fmt.Sprintf("role:%s", roleID), "assign", nil)
    if err != nil {
        return fmt.Errorf("permission check failed: %w", err)
    }
    
    if !canAssign {
        return fmt.Errorf("insufficient permissions to assign role")
    }
    
    // Add role assignment
    _, err = rbac.enforcer.AddGroupingPolicy(userID, roleID)
    if err != nil {
        return fmt.Errorf("failed to assign role: %w", err)
    }
    
    // Invalidate cache
    rbac.permissionCache.Invalidate(userID)
    
    // Log the assignment
    rbac.auditLogger.LogRoleAssignment(userID, roleID, assignedBy)
    
    return nil
}

// CreateRole creates a new role with permissions
func (rbac *AdvancedRBACSystem) CreateRole(ctx context.Context, role Role, createdBy string) error {
    // Validate role
    if err := rbac.validator.ValidateRole(role); err != nil {
        return fmt.Errorf("invalid role: %w", err)
    }
    
    // Check permission to create roles
    canCreate, err := rbac.CheckPermission(ctx, createdBy, "role", "create", nil)
    if err != nil {
        return fmt.Errorf("permission check failed: %w", err)
    }
    
    if !canCreate {
        return fmt.Errorf("insufficient permissions to create role")
    }
    
    // Store role
    role.CreatedBy = createdBy
    role.CreatedAt = time.Now()
    role.UpdatedAt = time.Now()
    
    if err := rbac.roleManager.CreateRole(role); err != nil {
        return fmt.Errorf("failed to create role: %w", err)
    }
    
    // Add role permissions to Casbin
    for _, permission := range role.Permissions {
        conditionsStr := rbac.serializeConditions(permission.Conditions)
        _, err := rbac.enforcer.AddPolicy(role.ID, permission.Resource, permission.Action, "allow", conditionsStr)
        if err != nil {
            return fmt.Errorf("failed to add role permission: %w", err)
        }
    }
    
    // Add role inheritance
    for _, parentRole := range role.Inherits {
        _, err := rbac.enforcer.AddGroupingPolicy(role.ID, parentRole)
        if err != nil {
            return fmt.Errorf("failed to add role inheritance: %w", err)
        }
    }
    
    // Log role creation
    rbac.auditLogger.LogRoleCreation(role.ID, createdBy)
    
    return nil
}

// GetUserPermissions returns all effective permissions for a user
func (rbac *AdvancedRBACSystem) GetUserPermissions(ctx context.Context, userID string) ([]Permission, error) {
    // Get direct permissions
    directPermissions := rbac.enforcer.GetPermissionsForUser(userID)
    
    // Get role-based permissions
    roles := rbac.enforcer.GetRolesForUser(userID)
    var rolePermissions [][]string
    
    for _, role := range roles {
        permissions := rbac.enforcer.GetPermissionsForUser(role)
        rolePermissions = append(rolePermissions, permissions...)
    }
    
    // Combine and convert to Permission objects
    allPermissions := append(directPermissions, rolePermissions...)
    var permissions []Permission
    
    for _, perm := range allPermissions {
        if len(perm) >= 3 {
            permission := Permission{
                Resource: perm[1],
                Action:   perm[2],
            }
            
            if len(perm) >= 5 {
                permission.Conditions = rbac.deserializeConditions(perm[4])
            }
            
            permissions = append(permissions, permission)
        }
    }
    
    return permissions, nil
}

// Policy validation
type PolicyValidator struct{}

func NewPolicyValidator() *PolicyValidator {
    return &PolicyValidator{}
}

func (pv *PolicyValidator) ValidateRole(role Role) error {
    if role.Name == "" {
        return fmt.Errorf("role name is required")
    }
    
    if role.OrganizationID == "" && !role.IsSystem {
        return fmt.Errorf("organization ID is required for non-system roles")
    }
    
    // Validate permissions
    for _, permission := range role.Permissions {
        if err := pv.ValidatePermission(permission); err != nil {
            return fmt.Errorf("invalid permission: %w", err)
        }
    }
    
    return nil
}

func (pv *PolicyValidator) ValidatePermission(permission Permission) error {
    if permission.Resource == "" {
        return fmt.Errorf("permission resource is required")
    }
    
    if permission.Action == "" {
        return fmt.Errorf("permission action is required")
    }
    
    // Validate conditions
    for _, condition := range permission.Conditions {
        if err := pv.ValidateCondition(condition); err != nil {
            return fmt.Errorf("invalid condition: %w", err)
        }
    }
    
    return nil
}

func (pv *PolicyValidator) ValidateCondition(condition Condition) error {
    if condition.Field == "" {
        return fmt.Errorf("condition field is required")
    }
    
    if condition.Value == nil {
        return fmt.Errorf("condition value is required")
    }
    
    return nil
}

// Helper methods
func (rbac *AdvancedRBACSystem) serializeContext(context map[string]interface{}) string {
    if context == nil {
        return "{}"
    }
    
    data, _ := json.Marshal(context)
    return string(data)
}

func (rbac *AdvancedRBACSystem) serializeConditions(conditions []Condition) string {
    if conditions == nil {
        return "true"
    }
    
    // Convert conditions to evaluable expression
    var expressions []string
    
    for _, condition := range conditions {
        expr := rbac.conditionToExpression(condition)
        expressions = append(expressions, expr)
    }
    
    return strings.Join(expressions, " && ")
}

func (rbac *AdvancedRBACSystem) conditionToExpression(condition Condition) string {
    switch condition.Type {
    case ConditionTypeTime:
        return fmt.Sprintf("timeInRange(ctx.time, '%v')", condition.Value)
    case ConditionTypeLocation:
        return fmt.Sprintf("locationMatches(ctx.location, '%v')", condition.Value)
    case ConditionTypeAttribute:
        return fmt.Sprintf("ctx.%s %s '%v'", condition.Field, condition.Operator, condition.Value)
    case ConditionTypeOwnership:
        return fmt.Sprintf("isOwner(ctx.user, ctx.resource)")
    default:
        return "true"
    }
}

func (rbac *AdvancedRBACSystem) deserializeConditions(conditionsStr string) []Condition {
    // This would implement the reverse of serializeConditions
    // For brevity, returning empty slice
    return []Condition{}
}

// Permission caching
type PermissionCache struct {
    cache map[string]*CacheEntry
    mutex sync.RWMutex
}

type CacheEntry struct {
    Value     bool
    ExpiresAt time.Time
}

func NewPermissionCache() *PermissionCache {
    return &PermissionCache{
        cache: make(map[string]*CacheEntry),
    }
}

func (pc *PermissionCache) Get(key string) *bool {
    pc.mutex.RLock()
    defer pc.mutex.RUnlock()
    
    entry, exists := pc.cache[key]
    if !exists || time.Now().After(entry.ExpiresAt) {
        return nil
    }
    
    return &entry.Value
}

func (pc *PermissionCache) Set(key string, value bool, ttl time.Duration) {
    pc.mutex.Lock()
    defer pc.mutex.Unlock()
    
    pc.cache[key] = &CacheEntry{
        Value:     value,
        ExpiresAt: time.Now().Add(ttl),
    }
}

func (pc *PermissionCache) Invalidate(userID string) {
    pc.mutex.Lock()
    defer pc.mutex.Unlock()
    
    // Remove all entries for the user
    for key := range pc.cache {
        if strings.HasPrefix(key, userID+":") {
            delete(pc.cache, key)
        }
    }
}
```

**TodoWrite Task**: `Implement enterprise-grade RBAC system with fine-grained permissions`

### KPIs for Phase 5C
- ✅ Enterprise-grade RBAC with fine-grained permissions
- ✅ Comprehensive audit logging for compliance
- ✅ SOC 2, GDPR, and HIPAA compliance features
- ✅ Advanced security policies and enforcement
- ✅ Multi-factor authentication and SSO integration
- ✅ Real-time security monitoring and alerting

---

## 🎯 PHASE 5 SUCCESS METRICS & VALIDATION

### Next-Generation Capabilities Metrics
- **Multi-Modal AI**: 95%+ accuracy in vision analysis tasks
- **Context Intelligence**: 40%+ reduction in token usage through optimization
- **PWA Experience**: 90%+ functionality available offline, <3s load times
- **Enterprise Security**: SOC 2 Type II compliance, comprehensive audit trails
- **Advanced Analytics**: Predictive insights with 85%+ accuracy
- **Innovation Index**: Leading-edge features adopted by 70%+ of users

### Technology Innovation Validation
```bash
# Test multi-modal AI capabilities
npm run test:vision-analysis
npm run test:context-intelligence

# Validate PWA functionality
npm run test:pwa-offline
npm run test:service-worker

# Enterprise security testing
npm run test:rbac-system
npm run test:audit-logging

# Performance benchmarking
npm run benchmark:next-gen-features
```

### Enterprise Readiness Checklist
- [ ] Multi-modal AI processing operational
- [ ] Intelligent context optimization deployed
- [ ] PWA with offline capabilities available
- [ ] Enterprise security and compliance verified
- [ ] Advanced analytics providing insights
- [ ] All innovation features tested and documented

---

## 🚀 FINAL SYSTEM TRANSFORMATION

With Phase 5 complete, Plandex has evolved into:

### Next-Generation AI Platform
- **Multi-Modal Intelligence**: Vision processing, advanced context understanding
- **Intelligent Optimization**: ML-powered context selection and pattern learning
- **Mobile-First Experience**: PWA with native app capabilities
- **Enterprise Ready**: Advanced security, compliance, and audit features
- **Predictive Analytics**: Data-driven insights and recommendations

### Innovation Leadership
The next-generation capabilities in Phase 5 position Plandex as:
- **Technology Pioneer**: Leading-edge AI and ML integration
- **Enterprise Standard**: Security and compliance for large organizations
- **Mobile Innovation**: Progressive web app with offline-first design
- **Intelligence Platform**: Context-aware, learning, and adaptive system

### Future-Proof Architecture
- **Scalable Infrastructure**: Distributed processing and edge computing ready
- **Extensible Platform**: Plugin architecture for continued innovation
- **AI Integration**: Multi-provider, multi-modal AI capabilities
- **Enterprise Platform**: Complete solution for organizations of any size

---

*This comprehensive next-generation capabilities guide transforms Plandex into a cutting-edge, enterprise-ready, AI-powered development platform that sets new standards for innovation, security, and user experience in the developer tools space.*