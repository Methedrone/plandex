# Phase 2: Performance & Infrastructure Optimization Guide
## Plandex App Upgrade - Comprehensive Performance Enhancement

---

## 🎯 EXECUTIVE SUMMARY

This guide provides detailed implementation of performance optimizations and infrastructure improvements for the Plandex application. Building on the secure foundation from Phase 1, this phase focuses on maximizing performance, especially for resource-constrained environments like MacBook 2012.

### Current Performance Analysis
- **Database**: PostgreSQL with basic connection pooling (needs optimization)
- **Memory Usage**: High goroutine count, unbounded channels (needs management)
- **Build Process**: Standard Go builds (missing optimization flags)
- **File I/O**: Synchronous operations in hot paths (needs async optimization)
- **Caching**: Minimal caching layers (needs comprehensive strategy)

### Performance Targets
- **Database Queries**: 50-70% faster response times
- **Memory Usage**: 30-40% reduction in peak memory
- **Build Time**: 40% faster Docker builds  
- **API Response**: 25-50% improvement in latency
- **MacBook 2012**: Optimized for 4-8GB RAM constraints

---

## 🔧 CLAUDE CODE WORKFLOW INTEGRATION

### MCP Servers Utilization
```bash
# Research performance best practices
use context7

# For complex optimization analysis
use SequentialThinking

# Break down performance tasks
use Task-Master with performance-focused PRD
```

### Performance Profiling Tools
```bash
# Enable pprof endpoints for continuous profiling
go tool pprof http://localhost:8080/debug/pprof/profile

# Memory profiling
go tool pprof http://localhost:8080/debug/pprof/heap

# Goroutine analysis
go tool pprof http://localhost:8080/debug/pprof/goroutine
```

---

## 📋 DETAILED IMPLEMENTATION PLAN

## Phase 2A: Database Layer Optimization
### 🎯 HIGH IMPACT OPTIMIZATION
**Target**: 50-70% improvement in database query performance

### Implementation Steps

#### Step 2A.1: Connection Pool Optimization
**File: `/app/server/db/pool.go`** (create new file)
```go
package db

import (
    "context"
    "database/sql"
    "fmt"
    "log"
    "time"
    
    "github.com/jackc/pgx/v5/pgxpool"
    _ "github.com/lib/pq"
)

type DatabasePool struct {
    Pool   *pgxpool.Pool
    Config *pgxpool.Config
}

// NewOptimizedPool creates an optimized connection pool
func NewOptimizedPool(connString string, isProduction bool) (*DatabasePool, error) {
    config, err := pgxpool.ParseConfig(connString)
    if err != nil {
        return nil, fmt.Errorf("failed to parse connection string: %w", err)
    }
    
    // Optimize pool settings based on environment
    if isProduction {
        config.MaxConns = 50                    // Production: higher concurrency  
        config.MinConns = 10                    // Keep connections warm
        config.MaxConnLifetime = time.Hour      // Rotate connections
        config.MaxConnIdleTime = 30 * time.Minute
        config.HealthCheckPeriod = 5 * time.Minute
    } else {
        // MacBook 2012 optimized settings
        config.MaxConns = 8                     // Limited by older hardware
        config.MinConns = 2                     // Minimal warm connections
        config.MaxConnLifetime = 30 * time.Minute
        config.MaxConnIdleTime = 10 * time.Minute
        config.HealthCheckPeriod = 2 * time.Minute
    }
    
    // Connection timeout settings
    config.ConnConfig.ConnectTimeout = 10 * time.Second
    config.ConnConfig.CommandTimeout = 30 * time.Second
    
    pool, err := pgxpool.NewWithConfig(context.Background(), config)
    if err != nil {
        return nil, fmt.Errorf("failed to create connection pool: %w", err)
    }
    
    return &DatabasePool{
        Pool:   pool,
        Config: config,
    }, nil
}

// GetPoolStats returns connection pool statistics
func (dp *DatabasePool) GetPoolStats() map[string]int32 {
    stats := dp.Pool.Stat()
    return map[string]int32{
        "total_conns":        stats.TotalConns(),
        "acquired_conns":     stats.AcquiredConns(),
        "idle_conns":        stats.IdleConns(),
        "max_conns":         stats.MaxConns(),
        "acquire_count":     stats.AcquireCount(),
        "acquire_duration":  int32(stats.AcquireDuration().Milliseconds()),
    }
}

// HealthCheck verifies database connectivity
func (dp *DatabasePool) HealthCheck(ctx context.Context) error {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()
    
    return dp.Pool.Ping(ctx)
}
```

**TodoWrite Task**: `Implement optimized database connection pooling`

#### Step 2A.2: Database Index Analysis and Optimization
**File: `/app/server/db/indexes.sql`** (create new file)
```sql
-- Analysis queries to identify missing indexes
-- Run these to analyze current query performance

-- Find slow queries (if pg_stat_statements is enabled)
SELECT 
    query,
    calls,
    total_time,
    mean_time,
    rows
FROM pg_stat_statements 
ORDER BY mean_time DESC 
LIMIT 20;

-- Find missing indexes on foreign keys
SELECT 
    c.conname AS constraint_name,
    t.relname AS table_name,
    ARRAY_AGG(a.attname ORDER BY a.attnum) AS columns
FROM pg_constraint c
JOIN pg_class t ON c.conrelid = t.oid
JOIN pg_attribute a ON a.attrelid = t.oid AND a.attnum = ANY(c.conkey)
WHERE c.contype = 'f'
AND NOT EXISTS (
    SELECT 1 FROM pg_index i 
    WHERE i.indrelid = c.conrelid 
    AND i.indkey::int[] @> c.conkey::int[]
)
GROUP BY c.conname, t.relname;

-- Suggested index optimizations based on Plandex schema analysis
-- Add these indexes based on common query patterns

-- User authentication queries
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_email 
ON users(email) WHERE deleted_at IS NULL;

-- Plan and project queries  
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_plans_user_id_status 
ON plans(user_id, status) WHERE archived_at IS NULL;

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_plans_org_id_created 
ON plans(org_id, created_at DESC) WHERE archived_at IS NULL;

-- Context and file operations
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_contexts_plan_id_active 
ON contexts(plan_id) WHERE active = true;

-- Message and conversation queries
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_conversations_plan_id_created 
ON conversations(plan_id, created_at DESC);

-- Model configuration lookups
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_models_org_id_active 
ON models(org_id) WHERE active = true;

-- Invite and permission queries
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_invites_email_status 
ON invites(email, status) WHERE expires_at > NOW();

-- Composite indexes for complex queries
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_plan_files_plan_path 
ON plan_files(plan_id, path) WHERE deleted_at IS NULL;

-- Partial indexes for soft deletes (common pattern in Plandex)
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_active_plans 
ON plans(created_at DESC) WHERE archived_at IS NULL AND deleted_at IS NULL;

-- Analyze tables after index creation
ANALYZE users;
ANALYZE plans; 
ANALYZE contexts;
ANALYZE conversations;
ANALYZE models;
ANALYZE invites;
ANALYZE plan_files;
```

**TodoWrite Task**: `Create and apply database performance indexes`

#### Step 2A.3: Query Optimization Patterns
**File: `/app/server/db/optimized_queries.go`** (create new file)
```go
package db

import (
    "context"
    "fmt"
    "strings"
    "time"
)

// OptimizedQueries contains performance-optimized database queries
type OptimizedQueries struct {
    pool *DatabasePool
}

func NewOptimizedQueries(pool *DatabasePool) *OptimizedQueries {
    return &OptimizedQueries{pool: pool}
}

// GetUserPlansOptimized - Optimized version avoiding N+1 queries
func (q *OptimizedQueries) GetUserPlansOptimized(ctx context.Context, userID string, limit int) ([]Plan, error) {
    query := `
        SELECT 
            p.id, p.name, p.status, p.created_at, p.updated_at,
            COUNT(c.id) as context_count,
            COUNT(pf.id) as file_count,
            COALESCE(latest_msg.created_at, p.created_at) as last_activity
        FROM plans p
        LEFT JOIN contexts c ON c.plan_id = p.id AND c.active = true
        LEFT JOIN plan_files pf ON pf.plan_id = p.id AND pf.deleted_at IS NULL
        LEFT JOIN LATERAL (
            SELECT created_at 
            FROM conversations 
            WHERE plan_id = p.id 
            ORDER BY created_at DESC 
            LIMIT 1
        ) latest_msg ON true
        WHERE p.user_id = $1 
        AND p.archived_at IS NULL 
        AND p.deleted_at IS NULL
        GROUP BY p.id, p.name, p.status, p.created_at, p.updated_at, latest_msg.created_at
        ORDER BY last_activity DESC
        LIMIT $2`
    
    rows, err := q.pool.Pool.Query(ctx, query, userID, limit)
    if err != nil {
        return nil, fmt.Errorf("failed to get user plans: %w", err)
    }
    defer rows.Close()
    
    var plans []Plan
    for rows.Next() {
        var plan Plan
        var contextCount, fileCount int
        var lastActivity time.Time
        
        err := rows.Scan(
            &plan.ID, &plan.Name, &plan.Status, &plan.CreatedAt, &plan.UpdatedAt,
            &contextCount, &fileCount, &lastActivity,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to scan plan: %w", err)
        }
        
        plan.ContextCount = contextCount
        plan.FileCount = fileCount
        plan.LastActivity = lastActivity
        plans = append(plans, plan)
    }
    
    return plans, nil
}

// BatchGetPlanContexts - Batch load contexts to avoid N+1
func (q *OptimizedQueries) BatchGetPlanContexts(ctx context.Context, planIDs []string) (map[string][]Context, error) {
    if len(planIDs) == 0 {
        return make(map[string][]Context), nil
    }
    
    // Create parameter placeholders for IN clause
    placeholders := make([]string, len(planIDs))
    params := make([]interface{}, len(planIDs))
    for i, id := range planIDs {
        placeholders[i] = fmt.Sprintf("$%d", i+1)
        params[i] = id
    }
    
    query := fmt.Sprintf(`
        SELECT plan_id, id, name, content, created_at
        FROM contexts 
        WHERE plan_id IN (%s) AND active = true
        ORDER BY plan_id, created_at`, 
        strings.Join(placeholders, ","))
    
    rows, err := q.pool.Pool.Query(ctx, query, params...)
    if err != nil {
        return nil, fmt.Errorf("failed to batch get contexts: %w", err)
    }
    defer rows.Close()
    
    result := make(map[string][]Context)
    for rows.Next() {
        var context Context
        var planID string
        
        err := rows.Scan(&planID, &context.ID, &context.Name, &context.Content, &context.CreatedAt)
        if err != nil {
            return nil, fmt.Errorf("failed to scan context: %w", err)
        }
        
        result[planID] = append(result[planID], context)
    }
    
    return result, nil
}

// GetPlanWithStatsOptimized - Single query for plan with statistics
func (q *OptimizedQueries) GetPlanWithStatsOptimized(ctx context.Context, planID string) (*PlanWithStats, error) {
    query := `
        SELECT 
            p.id, p.name, p.status, p.created_at, p.updated_at,
            p.user_id, p.org_id,
            COUNT(DISTINCT c.id) as context_count,
            COUNT(DISTINCT pf.id) as file_count,
            COUNT(DISTINCT conv.id) as conversation_count,
            SUM(CASE WHEN c.active = true THEN 1 ELSE 0 END) as active_context_count,
            AVG(LENGTH(c.content)) as avg_context_size
        FROM plans p
        LEFT JOIN contexts c ON c.plan_id = p.id
        LEFT JOIN plan_files pf ON pf.plan_id = p.id AND pf.deleted_at IS NULL
        LEFT JOIN conversations conv ON conv.plan_id = p.id
        WHERE p.id = $1 AND p.deleted_at IS NULL
        GROUP BY p.id, p.name, p.status, p.created_at, p.updated_at, p.user_id, p.org_id`
    
    var plan PlanWithStats
    var avgContextSize sql.NullFloat64
    
    err := q.pool.Pool.QueryRow(ctx, query, planID).Scan(
        &plan.ID, &plan.Name, &plan.Status, &plan.CreatedAt, &plan.UpdatedAt,
        &plan.UserID, &plan.OrgID,
        &plan.ContextCount, &plan.FileCount, &plan.ConversationCount,
        &plan.ActiveContextCount, &avgContextSize,
    )
    
    if err != nil {
        return nil, fmt.Errorf("failed to get plan with stats: %w", err)
    }
    
    if avgContextSize.Valid {
        plan.AvgContextSize = int(avgContextSize.Float64)
    }
    
    return &plan, nil
}
```

**TodoWrite Task**: `Implement optimized database query patterns`

#### Step 2A.4: Database Performance Monitoring
**File: `/app/server/monitoring/db_monitor.go`** (create new file)
```go
package monitoring

import (
    "context"
    "fmt"
    "log"
    "time"
    "yourapp/db"
)

type DBMonitor struct {
    pool     *db.DatabasePool
    metrics  *MetricsCollector
    slowQueryThreshold time.Duration
}

func NewDBMonitor(pool *db.DatabasePool, metrics *MetricsCollector) *DBMonitor {
    return &DBMonitor{
        pool:               pool,
        metrics:           metrics,
        slowQueryThreshold: 100 * time.Millisecond, // Log queries > 100ms
    }
}

// MonitorQueries wraps database queries with performance monitoring
func (m *DBMonitor) MonitorQuery(ctx context.Context, queryName string, fn func() error) error {
    start := time.Now()
    err := fn()
    duration := time.Since(start)
    
    // Record metrics
    m.metrics.RecordDBQuery(queryName, duration, err == nil)
    
    // Log slow queries
    if duration > m.slowQueryThreshold {
        log.Printf("SLOW QUERY [%s]: %v", queryName, duration)
    }
    
    return err
}

// CollectPoolMetrics gathers connection pool statistics
func (m *DBMonitor) CollectPoolMetrics() {
    stats := m.pool.GetPoolStats()
    
    m.metrics.SetGauge("db_pool_total_conns", float64(stats["total_conns"]))
    m.metrics.SetGauge("db_pool_acquired_conns", float64(stats["acquired_conns"]))
    m.metrics.SetGauge("db_pool_idle_conns", float64(stats["idle_conns"]))
    m.metrics.SetGauge("db_pool_max_conns", float64(stats["max_conns"]))
    
    // Connection pool utilization percentage
    utilization := float64(stats["acquired_conns"]) / float64(stats["max_conns"]) * 100
    m.metrics.SetGauge("db_pool_utilization_percent", utilization)
}

// StartPeriodicCollection starts background metrics collection
func (m *DBMonitor) StartPeriodicCollection() {
    ticker := time.NewTicker(30 * time.Second)
    go func() {
        for range ticker.C {
            m.CollectPoolMetrics()
        }
    }()
}
```

**TodoWrite Task**: `Implement database performance monitoring`

### KPIs for Phase 2A
- ✅ Database query response time improved by 50-70%
- ✅ Connection pool utilization optimized (80-90% efficiency)
- ✅ Slow query count reduced by 80%+
- ✅ Database connection errors eliminated
- ✅ N+1 query patterns eliminated

---

## Phase 2B: Memory Management Optimization
### 🧠 MEMORY EFFICIENCY TARGET
**Target**: 30-40% reduction in memory usage, optimized for MacBook 2012

### Implementation Steps

#### Step 2B.1: Object Pool Implementation
**File: `/app/server/pool/object_pools.go`** (create new file)
```go
package pool

import (
    "bytes"
    "encoding/json"
    "sync"
)

// ObjectPools contains all reusable object pools
type ObjectPools struct {
    BufferPool     *BufferPool
    JSONEncoderPool *JSONEncoderPool
    StringBuilderPool *StringBuilderPool
}

// NewObjectPools creates optimized object pools
func NewObjectPools() *ObjectPools {
    return &ObjectPools{
        BufferPool:       NewBufferPool(),
        JSONEncoderPool:  NewJSONEncoderPool(),
        StringBuilderPool: NewStringBuilderPool(),
    }
}

// BufferPool manages reusable byte buffers for I/O operations
type BufferPool struct {
    pool sync.Pool
}

func NewBufferPool() *BufferPool {
    return &BufferPool{
        pool: sync.Pool{
            New: func() interface{} {
                // Start with 4KB buffers, will grow as needed
                return bytes.NewBuffer(make([]byte, 0, 4096))
            },
        },
    }
}

func (bp *BufferPool) Get() *bytes.Buffer {
    buf := bp.pool.Get().(*bytes.Buffer)
    buf.Reset()
    return buf
}

func (bp *BufferPool) Put(buf *bytes.Buffer) {
    // Prevent memory leaks from oversized buffers
    if buf.Cap() > 64*1024 { // 64KB limit
        return
    }
    bp.pool.Put(buf)
}

// JSONEncoderPool manages reusable JSON encoders
type JSONEncoderPool struct {
    pool sync.Pool
}

func NewJSONEncoderPool() *JSONEncoderPool {
    return &JSONEncoderPool{
        pool: sync.Pool{
            New: func() interface{} {
                return json.NewEncoder(nil)
            },
        },
    }
}

func (jp *JSONEncoderPool) Get() *json.Encoder {
    return jp.pool.Get().(*json.Encoder)
}

func (jp *JSONEncoderPool) Put(enc *json.Encoder) {
    jp.pool.Put(enc)
}

// StringBuilderPool manages reusable string builders
type StringBuilderPool struct {
    pool sync.Pool
}

func NewStringBuilderPool() *StringBuilderPool {
    return &StringBuilderPool{
        pool: sync.Pool{
            New: func() interface{} {
                var sb strings.Builder
                sb.Grow(1024) // Pre-allocate 1KB
                return &sb
            },
        },
    }
}

func (sp *StringBuilderPool) Get() *strings.Builder {
    sb := sp.pool.Get().(*strings.Builder)
    sb.Reset()
    return sb
}

func (sp *StringBuilderPool) Put(sb *strings.Builder) {
    // Prevent memory retention from large builders
    if sb.Cap() > 32*1024 { // 32KB limit
        return
    }
    sp.pool.Put(sb)
}

// ContextPool manages reusable context objects for AI operations
type ContextPool struct {
    pool sync.Pool
}

type PooledContext struct {
    Files   []string
    Content []byte
    Tokens  int
}

func NewContextPool() *ContextPool {
    return &ContextPool{
        pool: sync.Pool{
            New: func() interface{} {
                return &PooledContext{
                    Files:   make([]string, 0, 10),
                    Content: make([]byte, 0, 1024*1024), // 1MB initial capacity
                }
            },
        },
    }
}

func (cp *ContextPool) Get() *PooledContext {
    ctx := cp.pool.Get().(*PooledContext)
    ctx.Files = ctx.Files[:0]
    ctx.Content = ctx.Content[:0]
    ctx.Tokens = 0
    return ctx
}

func (cp *ContextPool) Put(ctx *PooledContext) {
    // Prevent memory leaks from oversized contexts
    if cap(ctx.Content) > 10*1024*1024 { // 10MB limit
        return
    }
    cp.pool.Put(ctx)
}
```

**TodoWrite Task**: `Implement object pooling for memory efficiency`

#### Step 2B.2: Streaming Response Optimization
**File: `/app/server/streaming/optimized_stream.go`** (create new file)
```go
package streaming

import (
    "bufio"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"
    "yourapp/pool"
)

type OptimizedStreamer struct {
    pools   *pool.ObjectPools
    bufSize int
}

func NewOptimizedStreamer(pools *pool.ObjectPools) *OptimizedStreamer {
    return &OptimizedStreamer{
        pools:   pools,
        bufSize: 4096, // 4KB buffer size optimized for MacBook 2012
    }
}

// StreamAIResponse streams AI model responses with memory optimization
func (s *OptimizedStreamer) StreamAIResponse(ctx context.Context, w http.ResponseWriter, aiResponse io.Reader) error {
    // Set up streaming headers
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.Header().Set("Cache-Control", "no-cache")
    w.Header().Set("Connection", "keep-alive")
    
    // Use flusher for immediate response streaming
    flusher, ok := w.(http.Flusher)
    if !ok {
        return fmt.Errorf("streaming not supported")
    }
    
    // Get buffer from pool
    buffer := s.pools.BufferPool.Get()
    defer s.pools.BufferPool.Put(buffer)
    
    // Use buffered reader for efficient I/O
    reader := bufio.NewReaderSize(aiResponse, s.bufSize)
    
    // Stream in chunks to control memory usage
    chunk := make([]byte, 1024) // 1KB chunks for responsive streaming
    
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
            n, err := reader.Read(chunk)
            if n > 0 {
                // Write chunk to response
                if _, writeErr := w.Write(chunk[:n]); writeErr != nil {
                    return fmt.Errorf("write error: %w", writeErr)
                }
                flusher.Flush()
            }
            
            if err == io.EOF {
                return nil
            }
            if err != nil {
                return fmt.Errorf("read error: %w", err)
            }
            
            // Small delay to prevent overwhelming slower hardware
            time.Sleep(1 * time.Millisecond)
        }
    }
}

// StreamJSONResponse streams JSON responses with memory optimization
func (s *OptimizedStreamer) StreamJSONResponse(ctx context.Context, w http.ResponseWriter, data interface{}) error {
    w.Header().Set("Content-Type", "application/json")
    
    // Get JSON encoder from pool
    encoder := s.pools.JSONEncoderPool.Get()
    defer s.pools.JSONEncoderPool.Put(encoder)
    
    // Get buffer from pool for encoding
    buffer := s.pools.BufferPool.Get()
    defer s.pools.BufferPool.Put(buffer)
    
    // Encode to buffer first to handle potential errors
    encoder = json.NewEncoder(buffer)
    if err := encoder.Encode(data); err != nil {
        return fmt.Errorf("JSON encoding error: %w", err)
    }
    
    // Stream from buffer to response
    _, err := io.Copy(w, buffer)
    return err
}

// BatchStreamProcessor handles batch processing with memory limits
type BatchStreamProcessor struct {
    maxBatchSize int
    maxMemory    int64
    pools        *pool.ObjectPools
}

func NewBatchStreamProcessor(pools *pool.ObjectPools, maxMemoryMB int) *BatchStreamProcessor {
    return &BatchStreamProcessor{
        maxBatchSize: 100,                    // Max items per batch
        maxMemory:    int64(maxMemoryMB) * 1024 * 1024, // Convert MB to bytes
        pools:        pools,
    }
}

func (bp *BatchStreamProcessor) ProcessInBatches(ctx context.Context, items []interface{}, processor func(batch []interface{}) error) error {
    var currentMemory int64
    var batch []interface{}
    
    for _, item := range items {
        // Estimate memory usage (rough approximation)
        itemSize := int64(estimateObjectSize(item))
        
        // Check if adding this item would exceed memory limit
        if currentMemory+itemSize > bp.maxMemory || len(batch) >= bp.maxBatchSize {
            // Process current batch
            if len(batch) > 0 {
                if err := processor(batch); err != nil {
                    return fmt.Errorf("batch processing error: %w", err)
                }
                
                // Reset batch
                batch = batch[:0]
                currentMemory = 0
                
                // Check for context cancellation
                select {
                case <-ctx.Done():
                    return ctx.Err()
                default:
                }
            }
        }
        
        batch = append(batch, item)
        currentMemory += itemSize
    }
    
    // Process final batch
    if len(batch) > 0 {
        return processor(batch)
    }
    
    return nil
}

// estimateObjectSize provides rough memory usage estimation
func estimateObjectSize(obj interface{}) int {
    // This is a simplified estimation - in production, consider using
    // a more sophisticated memory measurement approach
    switch v := obj.(type) {
    case string:
        return len(v) + 16 // String overhead
    case []byte:
        return len(v) + 24 // Slice overhead
    case map[string]interface{}:
        size := 48 // Map overhead
        for k, val := range v {
            size += len(k) + estimateObjectSize(val)
        }
        return size
    default:
        return 64 // Default object overhead
    }
}
```

**TodoWrite Task**: `Implement streaming optimization for memory efficiency`

#### Step 2B.3: Garbage Collection Tuning for MacBook 2012
**File: `/app/server/runtime/gc_tuning.go`** (create new file)
```go
package runtime

import (
    "fmt"
    "log"
    "os"
    "runtime"
    "runtime/debug"
    "strconv"
    "time"
)

// GCTuner optimizes garbage collection for resource-constrained environments
type GCTuner struct {
    targetMemoryMB int
    isLowMemory    bool
}

// NewGCTuner creates a GC tuner optimized for MacBook 2012
func NewGCTuner() *GCTuner {
    tuner := &GCTuner{
        targetMemoryMB: 512, // Target 512MB for main process
    }
    
    // Detect if running on resource-constrained hardware
    tuner.detectHardwareConstraints()
    tuner.applyOptimalSettings()
    
    return tuner
}

func (gc *GCTuner) detectHardwareConstraints() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    // Check available system memory (approximation)
    totalMemoryMB := int(m.Sys / 1024 / 1024)
    
    // If total allocated memory is low, assume constrained environment
    if totalMemoryMB < 1024 { // Less than 1GB allocated suggests limited system
        gc.isLowMemory = true
        gc.targetMemoryMB = 256 // More aggressive memory target
        log.Println("Detected resource-constrained environment, applying conservative memory settings")
    }
}

func (gc *GCTuner) applyOptimalSettings() {
    if gc.isLowMemory {
        // Conservative settings for MacBook 2012
        debug.SetGCPercent(50)  // More frequent GC (default is 100)
        debug.SetMemoryLimit(int64(gc.targetMemoryMB) * 1024 * 1024) // Set memory limit
        
        // Reduce max processors if system has limited cores
        if runtime.NumCPU() <= 4 {
            runtime.GOMAXPROCS(runtime.NumCPU()) // Use all available cores
        } else {
            runtime.GOMAXPROCS(4) // Limit to 4 cores for older hardware
        }
        
        log.Printf("Applied low-memory GC settings: GCPercent=50, MemoryLimit=%dMB, GOMAXPROCS=%d", 
                   gc.targetMemoryMB, runtime.GOMAXPROCS(-1))
    } else {
        // Standard performance settings
        debug.SetGCPercent(100)
        debug.SetMemoryLimit(int64(gc.targetMemoryMB * 2) * 1024 * 1024) // More generous limit
        
        log.Printf("Applied standard GC settings: GCPercent=100, MemoryLimit=%dMB", gc.targetMemoryMB*2)
    }
    
    // Override from environment if specified
    if gcPercent := os.Getenv("GOGC"); gcPercent != "" {
        if val, err := strconv.Atoi(gcPercent); err == nil {
            debug.SetGCPercent(val)
            log.Printf("GC percent overridden by environment: %d", val)
        }
    }
}

// MonitorMemoryUsage provides periodic memory usage monitoring
func (gc *GCTuner) MonitorMemoryUsage() {
    ticker := time.NewTicker(60 * time.Second) // Check every minute
    go func() {
        for range ticker.C {
            var m runtime.MemStats
            runtime.ReadMemStats(&m)
            
            currentMB := int(m.Alloc / 1024 / 1024)
            
            // Log memory usage
            log.Printf("Memory usage: %dMB allocated, %dMB system, %d GC cycles", 
                      currentMB, int(m.Sys/1024/1024), m.NumGC)
            
            // Trigger manual GC if memory usage is high
            if currentMB > gc.targetMemoryMB {
                log.Printf("Memory usage (%dMB) exceeds target (%dMB), forcing GC", 
                          currentMB, gc.targetMemoryMB)
                runtime.GC()
            }
        }
    }()
}

// GetMemoryStats returns current memory statistics
func (gc *GCTuner) GetMemoryStats() map[string]interface{} {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    return map[string]interface{}{
        "alloc_mb":      int(m.Alloc / 1024 / 1024),
        "total_alloc_mb": int(m.TotalAlloc / 1024 / 1024),
        "sys_mb":        int(m.Sys / 1024 / 1024),
        "num_gc":        m.NumGC,
        "gc_percent":    debug.SetGCPercent(-1), // Get current setting
        "num_goroutine": runtime.NumGoroutine(),
        "gomaxprocs":    runtime.GOMAXPROCS(-1),
    }
}
```

**TodoWrite Task**: `Implement GC tuning for MacBook 2012 optimization`

### KPIs for Phase 2B
- ✅ Memory usage reduced by 30-40%
- ✅ GC pressure reduced (fewer, more efficient collections)
- ✅ Object allocation rate decreased by 50%+
- ✅ Buffer reuse efficiency >90%
- ✅ Streaming memory overhead <10MB

---

## Phase 2C: Build & Docker Optimization
### 🚀 BUILD PERFORMANCE TARGET  
**Target**: 40% faster builds, optimized Docker images

### Implementation Steps

#### Step 2C.1: Multi-Stage Docker Optimization
**File: `/app/server/Dockerfile.optimized`** (create optimized version)
```dockerfile
# Multi-stage build optimized for performance and size
ARG GO_VERSION=1.23.10
ARG ALPINE_VERSION=3.19

# Build stage - optimized for compilation speed
FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS builder

# Install build dependencies
RUN apk add --no-cache \
    git \
    ca-certificates \
    tzdata \
    upx

# Set build environment for performance
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://proxy.golang.org,direct \
    GOSUMDB=sum.golang.org

# Create non-root user for security
RUN adduser -D -s /bin/sh -u 1001 appuser

# Set working directory
WORKDIR /build

# Copy go mod files first for better layer caching
COPY go.mod go.sum ./
COPY shared/go.mod shared/go.sum ./shared/
COPY server/go.mod server/go.sum ./server/

# Download dependencies with optimizations
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Copy source code
COPY . .

# Build with optimization flags
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o plandex-server \
    ./server/main.go

# Compress binary with UPX (optional, can reduce size by 60%+)
RUN upx --best --lzma plandex-server

# Runtime stage - minimal image
FROM scratch AS runtime

# Copy CA certificates for HTTPS requests
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy user information
COPY --from=builder /etc/passwd /etc/passwd

# Copy the binary
COPY --from=builder /build/plandex-server /usr/local/bin/plandex-server

# Use non-root user
USER appuser

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/usr/local/bin/plandex-server", "--health-check"]

# Expose port
EXPOSE 8080

# Run the binary
ENTRYPOINT ["/usr/local/bin/plandex-server"]
```

**File: `/app/cli/Dockerfile.optimized`** (create optimized CLI version)
```dockerfile
# Multi-stage build for CLI
ARG GO_VERSION=1.23.10
ARG ALPINE_VERSION=3.19

FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS builder

RUN apk add --no-cache git ca-certificates upx

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

# Optimize layer caching
COPY go.mod go.sum ./
COPY shared/go.mod shared/go.sum ./shared/
COPY cli/go.mod cli/go.sum ./cli/

RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

# Build multiple architectures if needed
ARG TARGETARCH=amd64
ARG TARGETOS=linux

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    GOARCH=${TARGETARCH} GOOS=${TARGETOS} \
    go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o plandex-cli \
    ./cli/main.go

RUN upx --best --lzma plandex-cli

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/plandex-cli /usr/local/bin/plandex
ENTRYPOINT ["/usr/local/bin/plandex"]
```

**TodoWrite Task**: `Create optimized multi-stage Docker builds`

#### Step 2C.2: Go Build Optimization
**File: `/app/scripts/build-optimized.sh`** (create new script)
```bash
#!/bin/bash

# Optimized build script for Plandex
set -euo pipefail

# Configuration
BUILD_DIR="./dist"
VERSION="${VERSION:-$(git describe --tags --always)}"
BUILD_TIME="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
GIT_COMMIT="$(git rev-parse HEAD)"

# Build flags for optimization
LDFLAGS="-w -s"
LDFLAGS="$LDFLAGS -X main.Version=$VERSION"
LDFLAGS="$LDFLAGS -X main.BuildTime=$BUILD_TIME"
LDFLAGS="$LDFLAGS -X main.GitCommit=$GIT_COMMIT"

# Create build directory
mkdir -p "$BUILD_DIR"

# Set build environment for performance
export CGO_ENABLED=0
export GOPROXY=https://proxy.golang.org,direct
export GOSUMDB=sum.golang.org

# Function to build for target
build_target() {
    local os=$1
    local arch=$2
    local output_name=$3
    
    echo "Building $output_name for $os/$arch..."
    
    GOOS=$os GOARCH=$arch go build \
        -ldflags="$LDFLAGS" \
        -a -installsuffix cgo \
        -trimpath \
        -o "$BUILD_DIR/$output_name" \
        "$4"
    
    # Compress with UPX if available
    if command -v upx >/dev/null 2>&1; then
        echo "Compressing $output_name with UPX..."
        upx --best --lzma "$BUILD_DIR/$output_name" || echo "UPX compression failed, continuing..."
    fi
    
    echo "Built $output_name ($(du -h "$BUILD_DIR/$output_name" | cut -f1))"
}

echo "Starting optimized build process..."

# Build server
build_target linux amd64 plandex-server-linux-amd64 ./server/main.go
build_target darwin amd64 plandex-server-darwin-amd64 ./server/main.go
build_target darwin arm64 plandex-server-darwin-arm64 ./server/main.go

# Build CLI
build_target linux amd64 plandex-cli-linux-amd64 ./cli/main.go
build_target darwin amd64 plandex-cli-darwin-amd64 ./cli/main.go
build_target darwin arm64 plandex-cli-darwin-arm64 ./cli/main.go
build_target windows amd64 plandex-cli-windows-amd64.exe ./cli/main.go

echo "Build complete! Artifacts in $BUILD_DIR/"
ls -lh "$BUILD_DIR/"

# Generate checksums
echo "Generating checksums..."
cd "$BUILD_DIR"
sha256sum * > checksums.txt
cd ..

echo "Build optimization complete!"
```

**TodoWrite Task**: `Create optimized build scripts with compression`

#### Step 2C.3: Docker Compose Performance Optimization
**File: `/docker-compose.performance.yml`** (create performance-optimized version)
```yaml
version: '3.8'

services:
  postgres:
    image: postgres:17.5-alpine
    restart: unless-stopped
    environment:
      POSTGRES_DB: plandex
      POSTGRES_USER: plandex
      POSTGRES_PASSWORD: plandex
      POSTGRES_INITDB_ARGS: "--auth-host=scram-sha-256 --auth-local=scram-sha-256"
    command: >
      postgres
      -c shared_buffers=256MB
      -c effective_cache_size=1GB
      -c maintenance_work_mem=64MB
      -c work_mem=4MB
      -c random_page_cost=1.1  
      -c effective_io_concurrency=200
      -c max_connections=100
      -c max_worker_processes=4
      -c max_parallel_workers_per_gather=2
      -c max_parallel_workers=4
      -c checkpoint_completion_target=0.9
      -c wal_buffers=16MB
      -c default_statistics_target=100
      -c log_statement=none
      -c log_min_duration_statement=1000
      -c log_connections=off
      -c log_disconnections=off
      -c log_lock_waits=on
      -c deadlock_timeout=1s
      -c shared_preload_libraries=pg_stat_statements
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./postgresql.conf:/etc/postgresql/postgresql.conf:ro
    networks:
      - plandex-network
    deploy:
      resources:
        limits:
          memory: 512M  # Optimized for MacBook 2012
          cpus: '2.0'
        reservations:
          memory: 256M
          cpus: '1.0'

  plandex-server:
    build:
      context: .
      dockerfile: ./server/Dockerfile.optimized
      cache_from:
        - plandex-server:latest
      args:
        - BUILDKIT_INLINE_CACHE=1
    restart: unless-stopped
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgres://plandex:plandex@postgres:5432/plandex?sslmode=disable
      - PORT=8080
      - ENVIRONMENT=production
      - LOG_LEVEL=info
      - GOGC=50  # More aggressive GC for limited memory
      - GOMEMLIMIT=512MiB  # Set memory limit
    depends_on:
      - postgres
    networks:
      - plandex-network
    healthcheck:
      test: ["CMD", "/usr/local/bin/plandex-server", "--health-check"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s
    deploy:
      resources:
        limits:
          memory: 512M  # Optimized for MacBook 2012
          cpus: '2.0'
        reservations:
          memory: 256M
          cpus: '1.0'

  # Redis for caching (optional performance enhancement)
  redis:
    image: redis:7.2-alpine
    restart: unless-stopped
    command: >
      redis-server
      --maxmemory 128mb
      --maxmemory-policy allkeys-lru
      --save ""
      --appendonly no
    networks:
      - plandex-network
    deploy:
      resources:
        limits:
          memory: 128M
          cpus: '0.5'

volumes:
  postgres_data:
    driver: local

networks:
  plandex-network:
    driver: bridge
    ipam:
      config:
        - subnet: 172.20.0.0/16
```

**TodoWrite Task**: `Create performance-optimized Docker Compose configuration`

### KPIs for Phase 2C
- ✅ Docker build time reduced by 40%+
- ✅ Final image size reduced by 60%+ (with UPX compression)
- ✅ Build cache hit rate >80%
- ✅ Multi-architecture build support
- ✅ Resource usage optimized for MacBook 2012

---

## Phase 2D: Caching Implementation
### ⚡ STRATEGIC CACHING TARGET
**Target**: 40-60% reduction in repeated operations

### Implementation Steps

#### Step 2D.1: Multi-Level Caching Strategy
**File: `/app/server/cache/cache_manager.go`** (create new file)
```go
package cache

import (
    "context"
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "sync"
    "time"
)

// CacheManager provides multi-level caching
type CacheManager struct {
    l1Cache    *MemoryCache    // In-memory L1 cache
    l2Cache    *RedisCache     // Redis L2 cache (optional)
    stats      *CacheStats
    keyPrefix  string
}

// CacheStats tracks cache performance
type CacheStats struct {
    mu          sync.RWMutex
    Hits        int64
    Misses      int64
    Sets        int64
    Errors      int64
    Evictions   int64
}

func (s *CacheStats) RecordHit() {
    s.mu.Lock()
    s.Hits++
    s.mu.Unlock()
}

func (s *CacheStats) RecordMiss() {
    s.mu.Lock()
    s.Misses++
    s.mu.Unlock()
}

func (s *CacheStats) GetStats() map[string]int64 {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    total := s.Hits + s.Misses
    hitRate := int64(0)
    if total > 0 {
        hitRate = (s.Hits * 100) / total
    }
    
    return map[string]int64{
        "hits":       s.Hits,
        "misses":     s.Misses,
        "sets":       s.Sets,
        "errors":     s.Errors,
        "evictions":  s.Evictions,
        "hit_rate":   hitRate,
        "total":      total,
    }
}

// NewCacheManager creates an optimized cache manager
func NewCacheManager(redisURL string) *CacheManager {
    cm := &CacheManager{
        l1Cache:   NewMemoryCache(1000, 10*time.Minute), // 1000 items, 10min TTL
        stats:     &CacheStats{},
        keyPrefix: "plandex:",
    }
    
    // Initialize Redis L2 cache if URL provided
    if redisURL != "" {
        if redis, err := NewRedisCache(redisURL); err == nil {
            cm.l2Cache = redis
        }
    }
    
    return cm
}

// Get retrieves value from cache (L1 first, then L2)
func (cm *CacheManager) Get(ctx context.Context, key string) ([]byte, bool) {
    fullKey := cm.keyPrefix + key
    
    // Try L1 cache first
    if value, found := cm.l1Cache.Get(fullKey); found {
        cm.stats.RecordHit()
        return value, true
    }
    
    // Try L2 cache if available
    if cm.l2Cache != nil {
        if value, err := cm.l2Cache.Get(ctx, fullKey); err == nil && value != nil {
            // Store in L1 for future access
            cm.l1Cache.Set(fullKey, value, 5*time.Minute)
            cm.stats.RecordHit()
            return value, true
        }
    }
    
    cm.stats.RecordMiss()
    return nil, false
}

// Set stores value in both cache levels
func (cm *CacheManager) Set(ctx context.Context, key string, value []byte, ttl time.Duration) {
    fullKey := cm.keyPrefix + key
    
    // Store in L1 cache
    cm.l1Cache.Set(fullKey, value, ttl)
    
    // Store in L2 cache if available
    if cm.l2Cache != nil {
        cm.l2Cache.Set(ctx, fullKey, value, ttl)
    }
    
    cm.stats.mu.Lock()
    cm.stats.Sets++
    cm.stats.mu.Unlock()
}

// GetOrSet implements cache-aside pattern
func (cm *CacheManager) GetOrSet(ctx context.Context, key string, ttl time.Duration, loader func() ([]byte, error)) ([]byte, error) {
    // Try to get from cache first
    if value, found := cm.Get(ctx, key); found {
        return value, nil
    }
    
    // Load value using provided function
    value, err := loader()
    if err != nil {
        return nil, err
    }
    
    // Store in cache for future use
    cm.Set(ctx, key, value, ttl)
    
    return value, nil
}

// GenerateKey creates a cache key from multiple components
func (cm *CacheManager) GenerateKey(components ...string) string {
    hasher := sha256.New()
    for _, comp := range components {
        hasher.Write([]byte(comp))
    }
    return hex.EncodeToString(hasher.Sum(nil))[:16] // Use first 16 chars
}
```

**TodoWrite Task**: `Implement multi-level cache manager`

#### Step 2D.2: Context Caching System
**File: `/app/server/cache/context_cache.go`** (create new file)
```go
package cache

import (
    "context"
    "encoding/json"
    "fmt"
    "time"
    "yourapp/types"
)

// ContextCache manages caching for AI context operations
type ContextCache struct {
    cache   *CacheManager
    metrics *CacheMetrics
}

type CachedContext struct {
    Content   string    `json:"content"`
    Files     []string  `json:"files"`
    TokenCount int      `json:"token_count"`
    CreatedAt time.Time `json:"created_at"`
    Hash      string    `json:"hash"`
}

func NewContextCache(cache *CacheManager) *ContextCache {
    return &ContextCache{
        cache:   cache,
        metrics: NewCacheMetrics("context"),
    }
}

// GetProjectContext retrieves cached project context
func (cc *ContextCache) GetProjectContext(ctx context.Context, projectPath string, fileList []string) (*CachedContext, bool) {
    // Generate cache key based on project path and file list
    key := cc.cache.GenerateKey("project_context", projectPath, fmt.Sprintf("%v", fileList))
    
    data, found := cc.cache.Get(ctx, key)
    if !found {
        cc.metrics.RecordMiss("project_context")
        return nil, false
    }
    
    var cachedCtx CachedContext
    if err := json.Unmarshal(data, &cachedCtx); err != nil {
        cc.metrics.RecordError("project_context")
        return nil, false
    }
    
    cc.metrics.RecordHit("project_context")
    return &cachedCtx, true
}

// SetProjectContext stores project context in cache
func (cc *ContextCache) SetProjectContext(ctx context.Context, projectPath string, fileList []string, context *CachedContext) {
    key := cc.cache.GenerateKey("project_context", projectPath, fmt.Sprintf("%v", fileList))
    
    data, err := json.Marshal(context)
    if err != nil {
        cc.metrics.RecordError("project_context")
        return
    }
    
    // Cache for 30 minutes (contexts change frequently during development)
    cc.cache.Set(ctx, key, data, 30*time.Minute)
    cc.metrics.RecordSet("project_context")
}

// GetFileContent retrieves cached file content
func (cc *ContextCache) GetFileContent(ctx context.Context, filePath string, lastModified time.Time) ([]byte, bool) {
    key := cc.cache.GenerateKey("file_content", filePath, lastModified.Format(time.RFC3339))
    
    data, found := cc.cache.Get(ctx, key)
    if !found {
        cc.metrics.RecordMiss("file_content")
        return nil, false
    }
    
    cc.metrics.RecordHit("file_content")
    return data, true
}

// SetFileContent stores file content in cache
func (cc *ContextCache) SetFileContent(ctx context.Context, filePath string, lastModified time.Time, content []byte) {
    key := cc.cache.GenerateKey("file_content", filePath, lastModified.Format(time.RFC3339))
    
    // Cache file content for 1 hour
    cc.cache.Set(ctx, key, content, time.Hour)
    cc.metrics.RecordSet("file_content")
}

// GetModelResponse retrieves cached AI model response
func (cc *ContextCache) GetModelResponse(ctx context.Context, modelName, prompt string, temperature float64) ([]byte, bool) {
    key := cc.cache.GenerateKey("model_response", modelName, prompt, fmt.Sprintf("%.2f", temperature))
    
    data, found := cc.cache.Get(ctx, key)
    if !found {
        cc.metrics.RecordMiss("model_response")
        return nil, false
    }
    
    cc.metrics.RecordHit("model_response")
    return data, true
}

// SetModelResponse stores AI model response in cache
func (cc *ContextCache) SetModelResponse(ctx context.Context, modelName, prompt string, temperature float64, response []byte) {
    key := cc.cache.GenerateKey("model_response", modelName, prompt, fmt.Sprintf("%.2f", temperature))
    
    // Cache model responses for 2 hours (balance between performance and freshness)
    cc.cache.Set(ctx, key, response, 2*time.Hour)
    cc.metrics.RecordSet("model_response")
}

// InvalidateProjectCache removes cached data for a project
func (cc *ContextCache) InvalidateProjectCache(ctx context.Context, projectPath string) {
    // This would require a more sophisticated implementation
    // For now, we rely on TTL expiration
    cc.metrics.RecordEviction("project_context")
}

// CacheMetrics tracks cache performance by operation type
type CacheMetrics struct {
    name    string
    ops     map[string]*OperationStats
    mutex   sync.RWMutex
}

type OperationStats struct {
    Hits   int64
    Misses int64
    Sets   int64
    Errors int64
}

func NewCacheMetrics(name string) *CacheMetrics {
    return &CacheMetrics{
        name: name,
        ops:  make(map[string]*OperationStats),
    }
}

func (cm *CacheMetrics) getOrCreateStats(operation string) *OperationStats {
    cm.mutex.Lock()
    defer cm.mutex.Unlock()
    
    if stats, exists := cm.ops[operation]; exists {
        return stats
    }
    
    stats := &OperationStats{}
    cm.ops[operation] = stats
    return stats
}

func (cm *CacheMetrics) RecordHit(operation string) {
    stats := cm.getOrCreateStats(operation)
    stats.Hits++
}

func (cm *CacheMetrics) RecordMiss(operation string) {
    stats := cm.getOrCreateStats(operation)
    stats.Misses++
}

func (cm *CacheMetrics) RecordSet(operation string) {
    stats := cm.getOrCreateStats(operation)
    stats.Sets++
}

func (cm *CacheMetrics) RecordError(operation string) {
    stats := cm.getOrCreateStats(operation)
    stats.Errors++
}

func (cm *CacheMetrics) RecordEviction(operation string) {
    // Implementation for eviction tracking
}

func (cm *CacheMetrics) GetMetrics() map[string]map[string]int64 {
    cm.mutex.RLock()
    defer cm.mutex.RUnlock()
    
    result := make(map[string]map[string]int64)
    
    for op, stats := range cm.ops {
        total := stats.Hits + stats.Misses
        hitRate := int64(0)
        if total > 0 {
            hitRate = (stats.Hits * 100) / total
        }
        
        result[op] = map[string]int64{
            "hits":     stats.Hits,
            "misses":   stats.Misses,
            "sets":     stats.Sets,
            "errors":   stats.Errors,
            "hit_rate": hitRate,
            "total":    total,
        }
    }
    
    return result
}
```

**TodoWrite Task**: `Implement context-specific caching system`

### KPIs for Phase 2D
- ✅ Cache hit rate >70% for frequent operations
- ✅ Context loading time reduced by 50%+
- ✅ File content retrieval 80% faster
- ✅ Model response caching reduces API calls by 60%+
- ✅ Memory usage for caching <100MB

---

## 🧪 COMPREHENSIVE TESTING & BENCHMARKING

### Performance Benchmarking Suite
**File: `/app/benchmarks/performance_test.go`** (create new file)
```go
package benchmarks

import (
    "context"
    "testing"
    "time"
    "yourapp/db"
    "yourapp/cache"
)

func BenchmarkDatabaseQueries(b *testing.B) {
    // Setup optimized database pool
    pool, err := db.NewOptimizedPool("postgres://test:test@localhost/test", false)
    if err != nil {
        b.Fatal(err)
    }
    defer pool.Pool.Close()
    
    ctx := context.Background()
    
    b.Run("GetUserPlansOptimized", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            queries := db.NewOptimizedQueries(pool)
            _, err := queries.GetUserPlansOptimized(ctx, "test-user-id", 10)
            if err != nil {
                b.Error(err)
            }
        }
    })
    
    b.Run("BatchGetPlanContexts", func(b *testing.B) {
        planIDs := []string{"plan1", "plan2", "plan3", "plan4", "plan5"}
        for i := 0; i < b.N; i++ {
            queries := db.NewOptimizedQueries(pool)
            _, err := queries.BatchGetPlanContexts(ctx, planIDs)
            if err != nil {
                b.Error(err)
            }
        }
    })
}

func BenchmarkCachePerformance(b *testing.B) {
    cache := cache.NewCacheManager("")
    ctx := context.Background()
    
    testData := make([]byte, 1024) // 1KB test data
    for i := range testData {
        testData[i] = byte(i % 256)
    }
    
    b.Run("CacheSet", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            key := fmt.Sprintf("test-key-%d", i)
            cache.Set(ctx, key, testData, time.Minute)
        }
    })
    
    b.Run("CacheGet", func(b *testing.B) {
        // Pre-populate cache
        for i := 0; i < 1000; i++ {
            key := fmt.Sprintf("test-key-%d", i)
            cache.Set(ctx, key, testData, time.Minute)
        }
        
        b.ResetTimer()
        for i := 0; i < b.N; i++ {
            key := fmt.Sprintf("test-key-%d", i%1000)
            _, _ = cache.Get(ctx, key)
        }
    })
}

func BenchmarkMemoryPools(b *testing.B) {
    pools := pool.NewObjectPools()
    
    b.Run("BufferPoolUsage", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            buf := pools.BufferPool.Get()
            buf.WriteString("test data for buffer pool performance testing")
            pools.BufferPool.Put(buf)
        }
    })
    
    b.Run("StringBuilderPool", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            sb := pools.StringBuilderPool.Get()
            sb.WriteString("test data for string builder pool performance testing")
            pools.StringBuilderPool.Put(sb)
        }
    })
}
```

**TodoWrite Task**: `Create comprehensive performance benchmarking suite`

### Load Testing Scripts
**File: `/app/scripts/load-test.sh`** (create new script)
```bash
#!/bin/bash

# Load testing script for performance validation
set -euo pipefail

# Configuration
BASE_URL="${BASE_URL:-http://localhost:8080}"
CONCURRENT_USERS="${CONCURRENT_USERS:-10}"
TEST_DURATION="${TEST_DURATION:-60s}"
RAMP_UP_TIME="${RAMP_UP_TIME:-10s}"

echo "Starting load tests against $BASE_URL"
echo "Concurrent users: $CONCURRENT_USERS"
echo "Test duration: $TEST_DURATION"

# Test health endpoint
echo "Testing health endpoint..."
hey -n 1000 -c $CONCURRENT_USERS -t 30 "$BASE_URL/api/health"

# Test authentication endpoint
echo "Testing authentication endpoint..."
hey -n 500 -c 5 -t 30 -m POST -H "Content-Type: application/json" \
    -d '{"email":"test@example.com","password":"testpassword"}' \
    "$BASE_URL/api/auth/login"

# Test plan listing (requires auth)
echo "Testing plan listing endpoint..."
if [ -f "auth_token.txt" ]; then
    TOKEN=$(cat auth_token.txt)
    hey -n 1000 -c $CONCURRENT_USERS -t 30 \
        -H "Authorization: Bearer $TOKEN" \
        "$BASE_URL/api/plans"
else
    echo "Skipping authenticated endpoints (no auth token found)"
fi

echo "Load testing complete!"
```

**TodoWrite Task**: `Create load testing scripts for performance validation`

---

## 📊 SUCCESS METRICS & VALIDATION

### Performance KPIs Dashboard
- **Database Performance**: 50-70% improvement in query response times
- **Memory Usage**: 30-40% reduction in peak memory consumption  
- **Build Performance**: 40% faster Docker builds and deployments
- **API Response Times**: 25-50% improvement in average latency
- **Cache Efficiency**: >70% hit rate for frequent operations
- **Resource Utilization**: Optimized for MacBook 2012 constraints

### Monitoring Implementation
**File: `/app/server/monitoring/performance_monitor.go`** (create new file)
```go
package monitoring

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
    "runtime"
    "time"
    "yourapp/cache"
    "yourapp/db"
)

type PerformanceMonitor struct {
    dbPool    *db.DatabasePool
    cache     *cache.CacheManager
    startTime time.Time
}

func NewPerformanceMonitor(dbPool *db.DatabasePool, cache *cache.CacheManager) *PerformanceMonitor {
    return &PerformanceMonitor{
        dbPool:    dbPool,
        cache:     cache,
        startTime: time.Now(),
    }
}

// GetPerformanceMetrics returns comprehensive performance metrics
func (pm *PerformanceMonitor) GetPerformanceMetrics() map[string]interface{} {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    metrics := map[string]interface{}{
        "timestamp": time.Now().UTC(),
        "uptime_seconds": time.Since(pm.startTime).Seconds(),
        
        // Memory metrics
        "memory": map[string]interface{}{
            "alloc_mb":      int(m.Alloc / 1024 / 1024),
            "total_alloc_mb": int(m.TotalAlloc / 1024 / 1024),
            "sys_mb":        int(m.Sys / 1024 / 1024),
            "num_gc":        m.NumGC,
            "gc_cpu_fraction": m.GCCPUFraction,
        },
        
        // Runtime metrics  
        "runtime": map[string]interface{}{
            "num_goroutine": runtime.NumGoroutine(),
            "num_cpu":       runtime.NumCPU(),
            "gomaxprocs":    runtime.GOMAXPROCS(-1),
            "go_version":    runtime.Version(),
        },
    }
    
    // Database metrics
    if pm.dbPool != nil {
        metrics["database"] = pm.dbPool.GetPoolStats()
    }
    
    // Cache metrics  
    if pm.cache != nil {
        metrics["cache"] = pm.cache.stats.GetStats()
    }
    
    return metrics
}

// MetricsHandler provides HTTP endpoint for metrics
func (pm *PerformanceMonitor) MetricsHandler(w http.ResponseWriter, r *http.Request) {
    metrics := pm.GetPerformanceMetrics()
    
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(metrics); err != nil {
        http.Error(w, "Failed to encode metrics", http.StatusInternalServerError)
        return
    }
}

// StartPeriodicLogging logs performance metrics periodically
func (pm *PerformanceMonitor) StartPeriodicLogging() {
    ticker := time.NewTicker(5 * time.Minute)
    go func() {
        for range ticker.C {
            metrics := pm.GetPerformanceMetrics()
            log.Printf("PERFORMANCE_METRICS: %s", mustMarshalJSON(metrics))
        }
    }()
}

func mustMarshalJSON(v interface{}) string {
    data, _ := json.Marshal(v)
    return string(data)
}
```

**TodoWrite Task**: `Implement performance monitoring and metrics collection`

---

## 🔄 ROLLBACK & RECOVERY PROCEDURES

### Performance Rollback Strategy
```bash
# Quick rollback for performance issues
git checkout performance-optimizations-backup

# Selective rollback options
git checkout HEAD~1 -- app/server/db/pool.go          # Database optimizations
git checkout HEAD~1 -- app/server/cache/              # Caching layer
git checkout HEAD~1 -- app/server/pool/               # Object pooling
git checkout HEAD~1 -- docker-compose.performance.yml # Docker optimizations

# Emergency fallback
docker-compose -f docker-compose.yml up -d  # Revert to original config
```

### Performance Regression Detection
```bash
# Run performance regression tests
./scripts/performance-regression-test.sh

# Compare benchmarks
go test -bench=. -benchmem > new_benchmarks.txt
benchcmp old_benchmarks.txt new_benchmarks.txt
```

---

## 🎯 HANDOFF TO PHASE 3

Once Phase 2 performance optimizations are complete and validated:

### Pre-Phase 3 Checklist
- [ ] All performance KPIs met and validated
- [ ] MacBook 2012 optimization confirmed  
- [ ] Performance monitoring active
- [ ] Benchmark suite established
- [ ] Load testing passing
- [ ] Memory usage within targets
- [ ] Database performance improved
- [ ] Caching layer operational

### Performance Baseline for Phase 3
The optimized performance foundation from Phase 2 will support the comprehensive testing and CI/CD infrastructure in Phase 3, ensuring that quality assurance processes don't compromise the performance gains achieved.

---

*This guide provides comprehensive performance optimization specifically tailored for Claude Code execution, with detailed TodoWrite integration, MacBook 2012-specific optimizations, and measurable performance improvements across database, memory, build, and caching layers.*