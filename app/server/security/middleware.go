package security

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/time/rate"
)

// SecurityConfig holds configuration for security middleware
type SecurityConfig struct {
	TrustedOrigins   []string
	IsDevelopment    bool
	RateLimit        rate.Limit
	RateBurst        int
	MaxRequestSize   int64
	EnableCSP        bool
	EnableHSTS       bool
	EnableCORS       bool
}

// DefaultSecurityConfig returns a secure default configuration
func DefaultSecurityConfig() *SecurityConfig {
	isDev := os.Getenv("GOENV") == "development"
	
	config := &SecurityConfig{
		TrustedOrigins: []string{"http://localhost:3000", "http://localhost:8080"},
		IsDevelopment:  isDev,
		RateLimit:      rate.Limit(100), // 100 requests per second
		RateBurst:      200,             // Burst of 200 requests
		MaxRequestSize: 10 * 1024 * 1024, // 10MB (reduced from 1GB for security)
		EnableCSP:      true,
		EnableHSTS:     !isDev, // Only enable HSTS in production
		EnableCORS:     true,
	}
	
	// Production-specific trusted origins
	if !isDev {
		config.TrustedOrigins = []string{
			"https://yourdomain.com",
			"https://app.yourdomain.com",
		}
	}
	
	return config
}

// SecurityMiddleware provides comprehensive web security
type SecurityMiddleware struct {
	config      *SecurityConfig
	rateLimiter *rate.Limiter
}

// NewSecurityMiddleware creates a new security middleware instance
func NewSecurityMiddleware(config *SecurityConfig) *SecurityMiddleware {
	if config == nil {
		config = DefaultSecurityConfig()
	}
	
	return &SecurityMiddleware{
		config:      config,
		rateLimiter: rate.NewLimiter(config.RateLimit, config.RateBurst),
	}
}

// Apply applies all security middleware to the provided handler
func (sm *SecurityMiddleware) Apply(next http.Handler) http.Handler {
	// Apply middleware in the correct order (outer to inner)
	handler := next
	
	// 1. Security headers (innermost)
	handler = sm.securityHeadersMiddleware(handler)
	
	// 2. Content Security Policy
	if sm.config.EnableCSP {
		handler = sm.cspMiddleware(handler)
	}
	
	// 3. HSTS (only in production)
	if sm.config.EnableHSTS {
		handler = sm.hstsMiddleware(handler)
	}
	
	// 4. CORS
	if sm.config.EnableCORS {
		handler = sm.corsMiddleware(handler)
	}
	
	// 5. Input validation
	handler = sm.inputValidationMiddleware(handler)
	
	// 6. Rate limiting (outermost)
	handler = sm.rateLimitMiddleware(handler)
	
	return handler
}

// rateLimitMiddleware implements rate limiting
func (sm *SecurityMiddleware) rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip rate limiting for health checks
		if r.URL.Path == "/health" || r.URL.Path == "/version" {
			next.ServeHTTP(w, r)
			return
		}
		
		if !sm.rateLimiter.Allow() {
			w.Header().Set("Retry-After", "60")
			w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%.0f", float64(sm.config.RateLimit)))
			w.Header().Set("X-RateLimit-Remaining", "0")
			w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(time.Minute).Unix()))
			
			http.Error(w, "Rate limit exceeded. Please try again later.", http.StatusTooManyRequests)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

// inputValidationMiddleware validates incoming requests
func (sm *SecurityMiddleware) inputValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Content length validation
		if r.ContentLength > sm.config.MaxRequestSize {
			http.Error(w, "Request entity too large", http.StatusRequestEntityTooLarge)
			return
		}
		
		// Content type validation for POST/PUT/PATCH
		if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" {
			contentType := r.Header.Get("Content-Type")
			if contentType != "" &&
				!strings.HasPrefix(contentType, "application/json") &&
				!strings.HasPrefix(contentType, "application/x-www-form-urlencoded") &&
				!strings.HasPrefix(contentType, "multipart/form-data") {
				http.Error(w, "Unsupported Media Type", http.StatusUnsupportedMediaType)
				return
			}
		}
		
		// Host header validation (prevent Host header injection)
		host := r.Host
		if host == "" {
			http.Error(w, "Missing Host header", http.StatusBadRequest)
			return
		}
		
		// Basic path traversal protection
		if strings.Contains(r.URL.Path, "..") || strings.Contains(r.URL.Path, "//") {
			http.Error(w, "Invalid path", http.StatusBadRequest)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

// corsMiddleware implements Cross-Origin Resource Sharing
func (sm *SecurityMiddleware) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		
		// Check if origin is in trusted list
		if origin != "" {
			isAllowed := false
			for _, trusted := range sm.config.TrustedOrigins {
				if origin == trusted {
					isAllowed = true
					break
				}
			}
			
			if isAllowed {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}
		}
		
		// Set CORS headers for preflight requests
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, X-CSRF-Token")
			w.Header().Set("Access-Control-Expose-Headers", "X-CSRF-Token")
			w.Header().Set("Access-Control-Max-Age", "86400") // 24 hours
			
			if origin != "" {
				for _, trusted := range sm.config.TrustedOrigins {
					if origin == trusted {
						w.Header().Set("Access-Control-Allow-Credentials", "true")
						break
					}
				}
			}
			
			w.WriteHeader(http.StatusNoContent)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

// hstsMiddleware implements HTTP Strict Transport Security
func (sm *SecurityMiddleware) hstsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only apply HSTS to HTTPS requests
		if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
			w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		}
		next.ServeHTTP(w, r)
	})
}

// cspMiddleware implements Content Security Policy
func (sm *SecurityMiddleware) cspMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Generate nonce for inline scripts/styles
		nonce := generateSecureNonce()
		
		// API-focused CSP policy
		var csp string
		if sm.config.IsDevelopment {
			// More relaxed CSP for development
			csp = fmt.Sprintf(`default-src 'self' 'unsafe-inline' 'unsafe-eval'; connect-src 'self' ws: wss:; img-src 'self' data:; nonce-%s`, nonce)
		} else {
			// Strict CSP for production
			csp = fmt.Sprintf(`default-src 'none'; connect-src 'self'; base-uri 'none'; frame-ancestors 'none'; upgrade-insecure-requests; nonce-%s`, nonce)
		}
		
		w.Header().Set("Content-Security-Policy", csp)
		
		// Store nonce in context for use in templates if needed
		ctx := context.WithValue(r.Context(), "csp-nonce", nonce)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// securityHeadersMiddleware sets comprehensive security headers
func (sm *SecurityMiddleware) securityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headers := w.Header()
		
		// Prevent clickjacking
		headers.Set("X-Frame-Options", "SAMEORIGIN")
		
		// Prevent MIME sniffing
		headers.Set("X-Content-Type-Options", "nosniff")
		
		// XSS Protection (disable legacy filter that can be bypassed)
		headers.Set("X-XSS-Protection", "0")
		
		// Referrer Policy
		headers.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		
		// Cross-Origin policies
		headers.Set("Cross-Origin-Opener-Policy", "same-origin")
		headers.Set("Cross-Origin-Embedder-Policy", "require-corp")
		headers.Set("Cross-Origin-Resource-Policy", "cross-origin")
		
		// Permissions Policy (restrict dangerous features)
		headers.Set("Permissions-Policy", "camera=(), microphone=(), geolocation=(), payment=(), usb=()")
		
		// Remove server information
		headers.Set("Server", "")
		headers.Del("X-Powered-By")
		
		// Cache control for security-sensitive endpoints
		if strings.HasPrefix(r.URL.Path, "/api/") {
			headers.Set("Cache-Control", "no-cache, no-store, must-revalidate")
			headers.Set("Pragma", "no-cache")
			headers.Set("Expires", "0")
		}
		
		// Set secure cookie defaults (for future cookie usage)
		if !sm.config.IsDevelopment {
			// These headers help ensure secure cookie defaults
			headers.Add("Set-Cookie", "SameSite=Strict; Secure; HttpOnly")
		}
		
		next.ServeHTTP(w, r)
	})
}

// generateSecureNonce generates a cryptographically secure nonce for CSP
func generateSecureNonce() string {
	bytes := make([]byte, 16) // 128-bit nonce
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to less secure but still functional nonce
		return fmt.Sprintf("fallback-%d", time.Now().UnixNano())
	}
	return base64.StdEncoding.EncodeToString(bytes)
}

// GetCSPNonce extracts the CSP nonce from the request context
func GetCSPNonce(r *http.Request) string {
	if nonce, ok := r.Context().Value("csp-nonce").(string); ok {
		return nonce
	}
	return ""
}