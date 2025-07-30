package handlers

import (
	"encoding/json"
	"net/http"
	"plandex-server/cache"
	"plandex-server/performance"
	"runtime"
	"time"
)

// PerformanceStats represents comprehensive performance metrics
type PerformanceStats struct {
	Timestamp    time.Time                  `json:"timestamp"`
	Memory       MemoryStats               `json:"memory"`
	GC           GCStats                   `json:"gc"`
	Cache        *cache.CacheStats         `json:"cache,omitempty"`
	ObjectPools  ObjectPoolStats           `json:"object_pools"`
	WorkerPool   *performance.PoolMetrics  `json:"worker_pool,omitempty"`
	System       SystemStats               `json:"system"`
}

// MemoryStats represents memory usage statistics
type MemoryStats struct {
	HeapAlloc      uint64  `json:"heap_alloc_mb"`
	HeapSys        uint64  `json:"heap_sys_mb"`
	HeapObjects    uint64  `json:"heap_objects"`
	StackInUse     uint64  `json:"stack_inuse_mb"`
	TotalAlloc     uint64  `json:"total_alloc_mb"`
	MemoryPressure int     `json:"memory_pressure_percent"`
}

// GCStats represents garbage collection statistics
type GCStats struct {
	NextGC       uint64        `json:"next_gc_mb"`
	LastGC       time.Time     `json:"last_gc"`
	NumGC        uint32        `json:"num_gc"`
	PauseTotal   time.Duration `json:"pause_total_ns"`
	AvgPause     time.Duration `json:"avg_pause_ns"`
	GCPercent    int           `json:"gc_percent"`
}

// ObjectPoolStats represents object pool statistics
type ObjectPoolStats struct {
	Initialized bool   `json:"initialized"`
	Available   bool   `json:"available"`
	Pools       []Pool `json:"pools"`
}

// Pool represents individual pool statistics
type Pool struct {
	Name   string `json:"name"`
	Active bool   `json:"active"`
}

// SystemStats represents system-level statistics
type SystemStats struct {
	NumCPU       int           `json:"num_cpu"`
	GOMAXPROCS   int           `json:"gomaxprocs"`
	NumGoroutine int           `json:"num_goroutine"`
	Uptime       time.Duration `json:"uptime_seconds"`
}

var startTime = time.Now()

// GetPerformanceStats returns comprehensive performance statistics
func GetPerformanceStats(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	stats := collectPerformanceStats()
	
	w.Header().Set("Content-Type", "application/json")
	
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		http.Error(w, "Error encoding performance stats", http.StatusInternalServerError)
		return
	}
}

// collectPerformanceStats gathers all performance metrics
func collectPerformanceStats() PerformanceStats {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	
	stats := PerformanceStats{
		Timestamp: time.Now(),
		Memory: MemoryStats{
			HeapAlloc:      mem.HeapAlloc / 1024 / 1024,      // Convert to MB
			HeapSys:        mem.HeapSys / 1024 / 1024,        // Convert to MB
			HeapObjects:    mem.HeapObjects,
			StackInUse:     mem.StackInuse / 1024 / 1024,     // Convert to MB
			TotalAlloc:     mem.TotalAlloc / 1024 / 1024,     // Convert to MB
			MemoryPressure: performance.GetMemoryPressure(),
		},
		GC: GCStats{
			NextGC:    mem.NextGC / 1024 / 1024,             // Convert to MB
			LastGC:    time.Unix(0, int64(mem.LastGC)),
			NumGC:     mem.NumGC,
			PauseTotal: time.Duration(mem.PauseTotalNs),
			GCPercent: runtime.GC,
		},
		ObjectPools: ObjectPoolStats{
			Initialized: performance.GlobalObjectPools != nil,
			Available:   performance.GlobalObjectPools != nil,
			Pools: []Pool{
				{Name: "JSONEncoder", Active: performance.GlobalObjectPools != nil},
				{Name: "JSONDecoder", Active: performance.GlobalObjectPools != nil},
				{Name: "ByteBuffer", Active: performance.GlobalObjectPools != nil},
				{Name: "StringBuilder", Active: performance.GlobalObjectPools != nil},
				{Name: "RequestMap", Active: performance.GlobalObjectPools != nil},
				{Name: "ResponseMap", Active: performance.GlobalObjectPools != nil},
			},
		},
		System: SystemStats{
			NumCPU:       runtime.NumCPU(),
			GOMAXPROCS:   runtime.GOMAXPROCS(0),
			NumGoroutine: runtime.NumGoroutine(),
			Uptime:       time.Since(startTime),
		},
	}
	
	// Calculate average GC pause
	if mem.NumGC > 0 {
		stats.GC.AvgPause = time.Duration(mem.PauseTotalNs / uint64(mem.NumGC))
	}
	
	// Get cache statistics if available
	if cache.GlobalCacheManager != nil {
		if cacheStats, err := cache.GlobalCacheManager.GetStats(r.Context()); err == nil {
			stats.Cache = cacheStats
		}
	}
	
	// Get worker pool statistics if available
	if performance.GlobalWorkerPool != nil {
		workerStats := performance.GlobalWorkerPool.GetMetrics()
		stats.WorkerPool = &workerStats
	}
	
	return stats
}

// GetPerformanceHealth returns a simple health check with performance indicators
func GetPerformanceHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	pressure := performance.GetMemoryPressure()
	var status string
	var statusCode int
	
	switch {
	case pressure < 50:
		status = "healthy"
		statusCode = http.StatusOK
	case pressure < 80:
		status = "warning"
		statusCode = http.StatusOK
	default:
		status = "critical"
		statusCode = http.StatusServiceUnavailable
	}
	
	health := map[string]interface{}{
		"status":           status,
		"memory_pressure":  pressure,
		"timestamp":        time.Now(),
		"pools_active":     performance.GlobalObjectPools != nil,
		"cache_active":     cache.GlobalCacheManager != nil,
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(health)
}