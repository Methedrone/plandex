package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// CacheInterface defines the contract for all cache implementations
type CacheInterface interface {
	// Basic operations
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	
	// Advanced operations
	GetMultiple(ctx context.Context, keys []string) (map[string][]byte, error)
	SetMultiple(ctx context.Context, items map[string]CacheItem) error
	DeletePattern(ctx context.Context, pattern string) error
	
	// Cache management
	Clear(ctx context.Context) error
	Stats(ctx context.Context) (*CacheStats, error)
	Close() error
}

// CacheItem represents a cache entry with TTL
type CacheItem struct {
	Key       string        `json:"key"`
	Value     []byte        `json:"value"`
	TTL       time.Duration `json:"ttl"`
	CreatedAt time.Time     `json:"created_at"`
	ExpiresAt time.Time     `json:"expires_at"`
}

// CacheStats provides cache performance metrics
type CacheStats struct {
	HitCount       int64                    `json:"hit_count"`
	MissCount      int64                    `json:"miss_count"`
	HitRate        float64                  `json:"hit_rate"`
	ItemCount      int64                    `json:"item_count"`
	MemoryUsage    int64                    `json:"memory_usage"`
	EvictionCount  int64                    `json:"eviction_count"`
	LastAccess     time.Time                `json:"last_access"`
	DetailsByLayer map[string]*LayerStats   `json:"details_by_layer,omitempty"`
}

// LayerStats provides per-layer cache statistics
type LayerStats struct {
	HitCount      int64     `json:"hit_count"`
	MissCount     int64     `json:"miss_count"`
	ItemCount     int64     `json:"item_count"`
	MemoryUsage   int64     `json:"memory_usage"`
	EvictionCount int64     `json:"eviction_count"`
	LastAccess    time.Time `json:"last_access"`
}

// CacheStrategy defines the multi-layer caching strategy
type CacheStrategy struct {
	L1Config L1Config `json:"l1_config"`
	L2Config L2Config `json:"l2_config"`
}

// L1Config configures the in-memory cache layer
type L1Config struct {
	MaxItems     int           `json:"max_items"`
	DefaultTTL   time.Duration `json:"default_ttl"`
	CleanupInterval time.Duration `json:"cleanup_interval"`
}

// L2Config configures the Redis cache layer
type L2Config struct {
	Enabled        bool          `json:"enabled"`
	Host           string        `json:"host"`
	Port           int           `json:"port"`
	Password       string        `json:"password"`
	Database       int           `json:"database"`
	DefaultTTL     time.Duration `json:"default_ttl"`
	ConnTimeout    time.Duration `json:"conn_timeout"`
	ReadTimeout    time.Duration `json:"read_timeout"`
	WriteTimeout   time.Duration `json:"write_timeout"`
}

// CacheError represents cache-specific errors
type CacheError struct {
	Operation string
	Key       string
	Err       error
}

func (e *CacheError) Error() string {
	return fmt.Sprintf("cache %s error for key '%s': %v", e.Operation, e.Key, e.Err)
}

// Common cache keys and prefixes
const (
	// Context caching
	ContextKey         = "ctx:"
	ContextContentKey  = "ctx:content:"
	ContextTokensKey   = "ctx:tokens:"
	
	// AI Model responses
	ModelResponseKey   = "model:response:"
	ModelTokenUsageKey = "model:tokens:"
	
	// File mapping and syntax
	FileMapKey         = "file:map:"
	SyntaxTreeKey      = "syntax:tree:"
	
	// Database queries
	DBQueryKey         = "db:query:"
	DBResultKey        = "db:result:"
	
	// API responses
	APIResponseKey     = "api:response:"
	APIHealthKey       = "api:health:"
	
	// Default TTLs
	DefaultL1TTL       = 10 * time.Minute
	DefaultL2TTL       = 1 * time.Hour
	ContextTTL         = 30 * time.Minute
	ModelResponseTTL   = 2 * time.Hour
	FileMapTTL         = 1 * time.Hour
	DBQueryTTL         = 5 * time.Minute
	APIResponseTTL     = 1 * time.Minute
)

// Helper functions for common operations
func GenerateKey(prefix, identifier string) string {
	return fmt.Sprintf("%s%s", prefix, identifier)
}

func MarshalValue(value interface{}) ([]byte, error) {
	return json.Marshal(value)
}

func UnmarshalValue(data []byte, target interface{}) error {
	return json.Unmarshal(data, target)
}

// NewCacheError creates a new cache error
func NewCacheError(operation, key string, err error) *CacheError {
	return &CacheError{
		Operation: operation,
		Key:       key,
		Err:       err,
	}
}

// IsExpired checks if a cache item has expired
func (item *CacheItem) IsExpired() bool {
	return time.Now().After(item.ExpiresAt)
}

// TimeToLive returns the remaining TTL for a cache item
func (item *CacheItem) TimeToLive() time.Duration {
	if item.IsExpired() {
		return 0
	}
	return time.Until(item.ExpiresAt)
}

// CalculateHitRate calculates cache hit rate as a percentage
func (stats *CacheStats) CalculateHitRate() {
	total := stats.HitCount + stats.MissCount
	if total > 0 {
		stats.HitRate = float64(stats.HitCount) / float64(total) * 100
	} else {
		stats.HitRate = 0
	}
}