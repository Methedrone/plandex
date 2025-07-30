package performance

import (
	"log"
	"os"
	"runtime"
	"runtime/debug"
	"strconv"
	"time"
)

// GCConfig holds garbage collection optimization settings
type GCConfig struct {
	GOGC        int    // GC target percentage
	MemLimit    string // Memory limit (e.g., "512MiB")
	MaxProcs    int    // Maximum CPU cores to use
	Development bool   // Whether running in development mode
}

// MacBook2012Config returns optimized GC settings for MacBook 2012 constraints
func MacBook2012Config() GCConfig {
	isDev := os.Getenv("GOENV") == "development"
	
	return GCConfig{
		GOGC:        50,         // More aggressive GC (default is 100)
		MemLimit:    "512MiB",   // Conservative memory limit
		MaxProcs:    4,          // MacBook 2012 typically has 2-4 cores
		Development: isDev,
	}
}

// ProductionConfig returns optimized GC settings for production environments
func ProductionConfig() GCConfig {
	return GCConfig{
		GOGC:        100,        // Standard GC target
		MemLimit:    "2GiB",     // More generous memory limit
		MaxProcs:    0,          // Use all available cores
		Development: false,
	}
}

// InitializeGCOptimization configures the Go runtime for optimal performance
func InitializeGCOptimization() {
	var config GCConfig
	
	// Determine configuration based on environment
	if isConstrainedEnvironment() {
		config = MacBook2012Config()
		log.Println("Applying MacBook 2012 optimized GC settings")
	} else {
		config = ProductionConfig()
		log.Println("Applying production GC settings")
	}
	
	// Apply the configuration
	applyGCConfig(config)
	
	// Log the applied settings
	logGCSettings()
	
	// Set up periodic GC stats logging if in development
	if config.Development {
		go periodicGCStatsLogging()
	}
}

// applyGCConfig applies the garbage collection configuration
func applyGCConfig(config GCConfig) {
	// Set GOGC if not already set by environment variable
	if os.Getenv("GOGC") == "" {
		debug.SetGCPercent(config.GOGC)
		log.Printf("Set GOGC to %d (more aggressive GC)", config.GOGC)
	} else {
		log.Printf("Using GOGC from environment: %s", os.Getenv("GOGC"))
	}
	
	// Set memory limit if not already set by environment variable
	if os.Getenv("GOMEMLIMIT") == "" {
		debug.SetMemoryLimit(parseMemoryLimit(config.MemLimit))
		log.Printf("Set GOMEMLIMIT to %s", config.MemLimit)
	} else {
		log.Printf("Using GOMEMLIMIT from environment: %s", os.Getenv("GOMEMLIMIT"))
	}
	
	// Set GOMAXPROCS if not already set by environment variable
	if os.Getenv("GOMAXPROCS") == "" && config.MaxProcs > 0 {
		runtime.GOMAXPROCS(config.MaxProcs)
		log.Printf("Set GOMAXPROCS to %d", config.MaxProcs)
	} else {
		log.Printf("Using GOMAXPROCS: %d", runtime.GOMAXPROCS(0))
	}
}

// isConstrainedEnvironment detects if we're running in a resource-constrained environment
func isConstrainedEnvironment() bool {
	// Check for explicit environment variable
	if os.Getenv("CONSTRAINED_ENV") == "true" {
		return true
	}
	
	// Check for MacBook 2012 indicators
	if os.Getenv("GOENV") == "development" {
		return true
	}
	
	// Check system memory (if we can detect it)
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	
	// Heuristic: if available memory seems limited, apply constraints
	// This is a rough estimate and may not be completely accurate
	return runtime.NumCPU() <= 4
}

// parseMemoryLimit converts human-readable memory limits to bytes
func parseMemoryLimit(limit string) int64 {
	if limit == "" {
		return -1 // No limit
	}
	
	// Extract number and unit
	var value float64
	var unit string
	
	if n, err := strconv.ParseFloat(limit[:len(limit)-3], 64); err == nil {
		value = n
		unit = limit[len(limit)-3:]
	} else if n, err := strconv.ParseFloat(limit[:len(limit)-2], 64); err == nil {
		value = n
		unit = limit[len(limit)-2:]
	} else {
		log.Printf("Warning: Could not parse memory limit %s, using no limit", limit)
		return -1
	}
	
	// Convert to bytes
	switch unit {
	case "KiB", "KB":
		return int64(value * 1024)
	case "MiB", "MB":
		return int64(value * 1024 * 1024)
	case "GiB", "GB":
		return int64(value * 1024 * 1024 * 1024)
	default:
		log.Printf("Warning: Unknown memory unit %s, using no limit", unit)
		return -1
	}
}

// logGCSettings logs the current garbage collection settings
func logGCSettings() {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	
	log.Printf("GC Settings Applied:")
	// Get GOGC from environment since debug.GetGCPercent() is only available in Go 1.19+
	gogcEnv := os.Getenv("GOGC")
	if gogcEnv == "" {
		gogcEnv = "100" // default value
	}
	log.Printf("  GOGC: %s", gogcEnv)
	log.Printf("  GOMAXPROCS: %d", runtime.GOMAXPROCS(0))
	log.Printf("  NumCPU: %d", runtime.NumCPU())
	log.Printf("  Current Heap: %d KB", mem.HeapAlloc/1024)
	log.Printf("  System Memory: %d KB", mem.Sys/1024)
}

// periodicGCStatsLogging logs GC statistics periodically (development only)
func periodicGCStatsLogging() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	
	for range ticker.C {
		logGCStats()
	}
}

// logGCStats logs current garbage collection statistics
func logGCStats() {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	
	log.Printf("GC Stats:")
	log.Printf("  Heap Alloc: %d KB", mem.HeapAlloc/1024)
	log.Printf("  Heap Sys: %d KB", mem.HeapSys/1024)
	log.Printf("  Heap Objects: %d", mem.HeapObjects)
	log.Printf("  GC Cycles: %d", mem.NumGC)
	log.Printf("  Last GC: %v ago", time.Since(time.Unix(0, int64(mem.LastGC))))
	log.Printf("  Pause Total: %v", time.Duration(mem.PauseTotalNs))
	
	if mem.NumGC > 0 {
		log.Printf("  Avg Pause: %v", time.Duration(mem.PauseTotalNs/uint64(mem.NumGC)))
	}
}

// ForceGC triggers a garbage collection cycle (use sparingly)
func ForceGC() {
	before := getCurrentMemUsage()
	runtime.GC()
	after := getCurrentMemUsage()
	
	log.Printf("Forced GC: %d KB -> %d KB (freed %d KB)", 
		before/1024, after/1024, (before-after)/1024)
}

// getCurrentMemUsage returns current memory usage in bytes
func getCurrentMemUsage() uint64 {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	return mem.HeapAlloc
}

// GetMemoryPressure returns a score indicating memory pressure (0-100)
func GetMemoryPressure() int {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	
	// Simple heuristic based on heap usage vs system memory
	if mem.Sys == 0 {
		return 0
	}
	
	pressure := int((mem.HeapAlloc * 100) / mem.Sys)
	if pressure > 100 {
		pressure = 100
	}
	
	return pressure
}

// ShouldTriggerCleanup determines if we should trigger cleanup based on memory pressure
func ShouldTriggerCleanup() bool {
	pressure := GetMemoryPressure()
	return pressure > 75 // Trigger cleanup if memory pressure > 75%
}

// OptimizeForLowMemory applies additional optimizations when memory is constrained
func OptimizeForLowMemory() {
	log.Println("Applying low memory optimizations...")
	
	// More aggressive GC
	debug.SetGCPercent(25)
	
	// Force immediate garbage collection
	runtime.GC()
	
	// Force freeing of OS memory
	debug.FreeOSMemory()
	
	log.Println("Low memory optimizations applied")
}

// MemoryWatchdog monitors memory usage and applies optimizations when needed
func StartMemoryWatchdog() {
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		
		for range ticker.C {
			if ShouldTriggerCleanup() {
				log.Printf("High memory pressure detected (%d%%), triggering cleanup", GetMemoryPressure())
				
				// Trigger object pool cleanup if available
				if GlobalObjectPools != nil {
					GlobalObjectPools.Cleanup()
				}
				
				// Force GC if pressure is very high
				if GetMemoryPressure() > 90 {
					ForceGC()
				}
			}
		}
	}()
	
	log.Println("Memory watchdog started")
}