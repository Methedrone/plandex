package cache

import (
	"context"
	"log"
	"sync"
	"time"
)

// CacheManager manages the caching system and provides easy access to caches
type CacheManager struct {
	l1Cache CacheInterface
	mu      sync.RWMutex
	config  CacheStrategy
}

// GlobalCacheManager is the global cache manager instance
var (
	GlobalCacheManager *CacheManager
	cacheOnce         sync.Once
)

// InitializeCache initializes the global cache manager with optimized settings for MacBook 2012
func InitializeCache() {
	cacheOnce.Do(func() {
		// MacBook 2012 optimized cache configuration
		config := CacheStrategy{
			L1Config: L1Config{
				MaxItems:        1000,                // Conservative limit for memory constraints
				DefaultTTL:      DefaultL1TTL,        // 10 minutes
				CleanupInterval: 5 * time.Minute,     // Regular cleanup
			},
			L2Config: L2Config{
				Enabled: false, // Start with L1 only, L2 can be added later
			},
		}
		
		GlobalCacheManager = NewCacheManager(config)
		log.Println("Cache manager initialized with L1 memory cache")
	})
}

// NewCacheManager creates a new cache manager with the specified configuration
func NewCacheManager(config CacheStrategy) *CacheManager {
	manager := &CacheManager{
		config: config,
	}
	
	// Initialize L1 cache
	manager.l1Cache = NewL1MemoryCache(config.L1Config)
	
	return manager
}

// GetL1Cache returns the L1 cache instance
func (cm *CacheManager) GetL1Cache() CacheInterface {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.l1Cache
}

// CacheContext caches context data with appropriate TTL
func (cm *CacheManager) CacheContext(ctx context.Context, key string, data []byte) error {
	cache := cm.GetL1Cache()
	return cache.Set(ctx, GenerateKey(ContextKey, key), data, ContextTTL)
}

// GetContext retrieves cached context data
func (cm *CacheManager) GetContext(ctx context.Context, key string) ([]byte, error) {
	cache := cm.GetL1Cache()
	return cache.Get(ctx, GenerateKey(ContextKey, key))
}

// CacheContextContent caches context content with extended TTL
func (cm *CacheManager) CacheContextContent(ctx context.Context, key string, data []byte) error {
	cache := cm.GetL1Cache()
	return cache.Set(ctx, GenerateKey(ContextContentKey, key), data, ContextTTL)
}

// GetContextContent retrieves cached context content
func (cm *CacheManager) GetContextContent(ctx context.Context, key string) ([]byte, error) {
	cache := cm.GetL1Cache()
	return cache.Get(ctx, GenerateKey(ContextContentKey, key))
}

// CacheModelResponse caches AI model responses
func (cm *CacheManager) CacheModelResponse(ctx context.Context, key string, data []byte) error {
	cache := cm.GetL1Cache()
	return cache.Set(ctx, GenerateKey(ModelResponseKey, key), data, ModelResponseTTL)
}

// GetModelResponse retrieves cached AI model responses
func (cm *CacheManager) GetModelResponse(ctx context.Context, key string) ([]byte, error) {
	cache := cm.GetL1Cache()
	return cache.Get(ctx, GenerateKey(ModelResponseKey, key))
}

// CacheFileMap caches file mapping data
func (cm *CacheManager) CacheFileMap(ctx context.Context, key string, data []byte) error {
	cache := cm.GetL1Cache()
	return cache.Set(ctx, GenerateKey(FileMapKey, key), data, FileMapTTL)
}

// GetFileMap retrieves cached file mapping data
func (cm *CacheManager) GetFileMap(ctx context.Context, key string) ([]byte, error) {
	cache := cm.GetL1Cache()
	return cache.Get(ctx, GenerateKey(FileMapKey, key))
}

// CacheSyntaxTree caches syntax tree data
func (cm *CacheManager) CacheSyntaxTree(ctx context.Context, key string, data []byte) error {
	cache := cm.GetL1Cache()
	return cache.Set(ctx, GenerateKey(SyntaxTreeKey, key), data, FileMapTTL)
}

// GetSyntaxTree retrieves cached syntax tree data
func (cm *CacheManager) GetSyntaxTree(ctx context.Context, key string) ([]byte, error) {
	cache := cm.GetL1Cache()
	return cache.Get(ctx, GenerateKey(SyntaxTreeKey, key))
}

// CacheDBQuery caches database query results
func (cm *CacheManager) CacheDBQuery(ctx context.Context, key string, data []byte) error {
	cache := cm.GetL1Cache()
	return cache.Set(ctx, GenerateKey(DBQueryKey, key), data, DBQueryTTL)
}

// GetDBQuery retrieves cached database query results
func (cm *CacheManager) GetDBQuery(ctx context.Context, key string) ([]byte, error) {
	cache := cm.GetL1Cache()
	return cache.Get(ctx, GenerateKey(DBQueryKey, key))
}

// CacheAPIResponse caches API responses
func (cm *CacheManager) CacheAPIResponse(ctx context.Context, key string, data []byte) error {
	cache := cm.GetL1Cache()
	return cache.Set(ctx, GenerateKey(APIResponseKey, key), data, APIResponseTTL)
}

// GetAPIResponse retrieves cached API responses
func (cm *CacheManager) GetAPIResponse(ctx context.Context, key string) ([]byte, error) {
	cache := cm.GetL1Cache()
	return cache.Get(ctx, GenerateKey(APIResponseKey, key))
}

// InvalidatePattern invalidates all cache entries matching a pattern
func (cm *CacheManager) InvalidatePattern(ctx context.Context, pattern string) error {
	cache := cm.GetL1Cache()
	return cache.DeletePattern(ctx, pattern)
}

// InvalidateContext invalidates context-related cache entries
func (cm *CacheManager) InvalidateContext(ctx context.Context, key string) error {
	cache := cm.GetL1Cache()
	// Invalidate both context and context content
	cache.Delete(ctx, GenerateKey(ContextKey, key))
	return cache.Delete(ctx, GenerateKey(ContextContentKey, key))
}

// GetStats returns cache statistics
func (cm *CacheManager) GetStats(ctx context.Context) (*CacheStats, error) {
	cache := cm.GetL1Cache()
	return cache.Stats(ctx)
}

// WarmupCommonData pre-loads frequently accessed data into cache
func (cm *CacheManager) WarmupCommonData(ctx context.Context) {
	// This can be implemented to pre-load common data
	log.Println("Cache warmup initiated")
}

// Cleanup performs cache maintenance and cleanup
func (cm *CacheManager) Cleanup() {
	cm.mu.RLock()
	cache := cm.l1Cache
	cm.mu.RUnlock()
	
	if cache != nil {
		cache.Close()
	}
}

// Helper functions for common cache operations

// CacheWithFallback attempts to get from cache, calls fallback if miss, then caches result
func CacheWithFallback[T any](
	ctx context.Context,
	key string,
	fallback func() (T, error),
	marshal func(T) ([]byte, error),
	unmarshal func([]byte) (T, error),
	cacheFunc func(context.Context, string, []byte) error,
	getFunc func(context.Context, string) ([]byte, error),
) (T, error) {
	var zero T
	
	// Try to get from cache first
	if data, err := getFunc(ctx, key); err == nil {
		if result, err := unmarshal(data); err == nil {
			return result, nil
		}
	}
	
	// Cache miss or unmarshal error, call fallback
	result, err := fallback()
	if err != nil {
		return zero, err
	}
	
	// Cache the result for next time
	if data, err := marshal(result); err == nil {
		cacheFunc(ctx, key, data) // Ignore cache errors
	}
	
	return result, nil
}