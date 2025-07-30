package cache

import (
	"context"
	"fmt"
	"regexp"
	"sync"
	"time"
)

// L1MemoryCache implements a thread-safe in-memory cache with TTL support
type L1MemoryCache struct {
	items         map[string]*cacheItem
	mutex         sync.RWMutex
	maxItems      int
	defaultTTL    time.Duration
	cleanupTicker *time.Ticker
	stats         *CacheStats
	stopCleanup   chan struct{}
}

// cacheItem represents an item stored in the cache
type cacheItem struct {
	value     []byte
	expiresAt time.Time
	createdAt time.Time
	hits      int64
}

// NewL1MemoryCache creates a new L1 memory cache with specified configuration
func NewL1MemoryCache(config L1Config) *L1MemoryCache {
	cache := &L1MemoryCache{
		items:       make(map[string]*cacheItem),
		maxItems:    config.MaxItems,
		defaultTTL:  config.DefaultTTL,
		stopCleanup: make(chan struct{}),
		stats: &CacheStats{
			DetailsByLayer: make(map[string]*LayerStats),
		},
	}
	
	// Initialize L1 layer stats
	cache.stats.DetailsByLayer["L1"] = &LayerStats{}
	
	// Start background cleanup if cleanup interval is specified
	if config.CleanupInterval > 0 {
		cache.cleanupTicker = time.NewTicker(config.CleanupInterval)
		go cache.runCleanup()
	}
	
	return cache
}

// Get retrieves a value from the cache
func (c *L1MemoryCache) Get(ctx context.Context, key string) ([]byte, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	
	item, exists := c.items[key]
	if !exists {
		c.updateStats(false, "L1")
		return nil, NewCacheError("get", key, fmt.Errorf("key not found"))
	}
	
	// Check if item has expired
	if time.Now().After(item.expiresAt) {
		c.mutex.RUnlock()
		c.mutex.Lock()
		delete(c.items, key)
		c.mutex.Unlock()
		c.mutex.RLock()
		
		c.updateStats(false, "L1")
		return nil, NewCacheError("get", key, fmt.Errorf("key expired"))
	}
	
	// Update hit count and stats
	item.hits++
	c.updateStats(true, "L1")
	
	// Return a copy to prevent external modifications
	result := make([]byte, len(item.value))
	copy(result, item.value)
	return result, nil
}

// Set stores a value in the cache with specified TTL
func (c *L1MemoryCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	if ttl <= 0 {
		ttl = c.defaultTTL
	}
	
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	// Check if we need to evict items to make room
	if len(c.items) >= c.maxItems {
		c.evictOldest()
	}
	
	// Store a copy to prevent external modifications
	valueCopy := make([]byte, len(value))
	copy(valueCopy, value)
	
	now := time.Now()
	c.items[key] = &cacheItem{
		value:     valueCopy,
		expiresAt: now.Add(ttl),
		createdAt: now,
		hits:      0,
	}
	
	c.updateItemCount()
	return nil
}

// Delete removes a key from the cache
func (c *L1MemoryCache) Delete(ctx context.Context, key string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	delete(c.items, key)
	c.updateItemCount()
	return nil
}

// Exists checks if a key exists in the cache and is not expired
func (c *L1MemoryCache) Exists(ctx context.Context, key string) (bool, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	
	item, exists := c.items[key]
	if !exists {
		return false, nil
	}
	
	// Check expiration
	if time.Now().After(item.expiresAt) {
		return false, nil
	}
	
	return true, nil
}

// GetMultiple retrieves multiple values from the cache
func (c *L1MemoryCache) GetMultiple(ctx context.Context, keys []string) (map[string][]byte, error) {
	result := make(map[string][]byte)
	
	for _, key := range keys {
		if value, err := c.Get(ctx, key); err == nil {
			result[key] = value
		}
	}
	
	return result, nil
}

// SetMultiple stores multiple values in the cache
func (c *L1MemoryCache) SetMultiple(ctx context.Context, items map[string]CacheItem) error {
	for key, item := range items {
		if err := c.Set(ctx, key, item.Value, item.TTL); err != nil {
			return err
		}
	}
	return nil
}

// DeletePattern removes all keys matching a pattern
func (c *L1MemoryCache) DeletePattern(ctx context.Context, pattern string) error {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return NewCacheError("delete_pattern", pattern, err)
	}
	
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	var keysToDelete []string
	for key := range c.items {
		if regex.MatchString(key) {
			keysToDelete = append(keysToDelete, key)
		}
	}
	
	for _, key := range keysToDelete {
		delete(c.items, key)
	}
	
	c.updateItemCount()
	return nil
}

// Clear removes all items from the cache
func (c *L1MemoryCache) Clear(ctx context.Context) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	c.items = make(map[string]*cacheItem)
	c.updateItemCount()
	return nil
}

// Stats returns current cache statistics
func (c *L1MemoryCache) Stats(ctx context.Context) (*CacheStats, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	
	// Update current stats
	c.stats.ItemCount = int64(len(c.items))
	c.stats.CalculateHitRate()
	c.stats.LastAccess = time.Now()
	
	// Update L1 layer stats
	if l1Stats := c.stats.DetailsByLayer["L1"]; l1Stats != nil {
		l1Stats.ItemCount = c.stats.ItemCount
		l1Stats.LastAccess = c.stats.LastAccess
	}
	
	// Calculate memory usage estimate
	var memoryUsage int64
	for _, item := range c.items {
		memoryUsage += int64(len(item.value))
	}
	c.stats.MemoryUsage = memoryUsage
	
	// Return a copy
	statsCopy := *c.stats
	return &statsCopy, nil
}

// Close stops the cleanup routine and cleans up resources
func (c *L1MemoryCache) Close() error {
	if c.cleanupTicker != nil {
		c.cleanupTicker.Stop()
		close(c.stopCleanup)
	}
	return nil
}

// updateStats updates cache hit/miss statistics
func (c *L1MemoryCache) updateStats(hit bool, layer string) {
	if hit {
		c.stats.HitCount++
		if layerStats := c.stats.DetailsByLayer[layer]; layerStats != nil {
			layerStats.HitCount++
		}
	} else {
		c.stats.MissCount++
		if layerStats := c.stats.DetailsByLayer[layer]; layerStats != nil {
			layerStats.MissCount++
		}
	}
}

// updateItemCount updates the item count in stats
func (c *L1MemoryCache) updateItemCount() {
	c.stats.ItemCount = int64(len(c.items))
	if layerStats := c.stats.DetailsByLayer["L1"]; layerStats != nil {
		layerStats.ItemCount = c.stats.ItemCount
	}
}

// evictOldest removes the oldest item to make room for new items
func (c *L1MemoryCache) evictOldest() {
	var oldestKey string
	var oldestTime time.Time
	first := true
	
	for key, item := range c.items {
		if first || item.createdAt.Before(oldestTime) {
			oldestKey = key
			oldestTime = item.createdAt
			first = false
		}
	}
	
	if oldestKey != "" {
		delete(c.items, oldestKey)
		c.stats.EvictionCount++
		if layerStats := c.stats.DetailsByLayer["L1"]; layerStats != nil {
			layerStats.EvictionCount++
		}
	}
}

// runCleanup runs the background cleanup routine
func (c *L1MemoryCache) runCleanup() {
	for {
		select {
		case <-c.cleanupTicker.C:
			c.cleanupExpired()
		case <-c.stopCleanup:
			return
		}
	}
}

// cleanupExpired removes expired items from the cache
func (c *L1MemoryCache) cleanupExpired() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	now := time.Now()
	var expiredKeys []string
	
	for key, item := range c.items {
		if now.After(item.expiresAt) {
			expiredKeys = append(expiredKeys, key)
		}
	}
	
	for _, key := range expiredKeys {
		delete(c.items, key)
	}
	
	if len(expiredKeys) > 0 {
		c.updateItemCount()
	}
}