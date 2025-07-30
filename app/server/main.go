package main

import (
	"log"
	"os"
	"plandex-server/cache"
	"plandex-server/performance"
	"plandex-server/routes"
	"plandex-server/setup"

	"github.com/gorilla/mux"
)

func main() {
	// Configure the default logger to include milliseconds in timestamps
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	
	// Initialize GC optimization early for maximum benefit
	performance.InitializeGCOptimization()
	
	// Start memory watchdog for resource-constrained environments
	performance.StartMemoryWatchdog()

	routes.RegisterHandlePlandex(func(router *mux.Router, path string, isStreaming bool, handler routes.PlandexHandler) *mux.Route {
		return router.HandleFunc(path, handler)
	})

	r := mux.NewRouter()
	routes.AddHealthRoutes(r)
	routes.AddApiRoutes(r)
	routes.AddProxyableApiRoutes(r)
	setup.MustLoadIp()
	setup.MustInitDb()
	
	// Initialize performance pools for memory efficiency
	performance.InitializePools()
	log.Println("Performance pools initialized")
	
	// Initialize L1 memory cache system
	cache.InitializeCache()
	log.Println("Cache system initialized")
	
	// Register cache cleanup on shutdown
	setup.RegisterShutdownHook(func() {
		log.Println("Cleaning up cache system...")
		if cache.GlobalCacheManager != nil {
			cache.GlobalCacheManager.Cleanup()
		}
	})
	
	setup.StartServer(r, nil, nil)
	os.Exit(0)
}
