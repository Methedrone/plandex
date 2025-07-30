package security

import (
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// SecureLogger provides secure logging utilities with automatic data sanitization
type SecureLogger struct {
	isDevelopment bool
}

// NewSecureLogger creates a new secure logger instance
func NewSecureLogger(isDevelopment bool) *SecureLogger {
	return &SecureLogger{
		isDevelopment: isDevelopment,
	}
}

// SensitiveDataPatterns contains regex patterns for identifying sensitive data
var SensitiveDataPatterns = []*regexp.Regexp{
	// API Keys
	regexp.MustCompile(`(?i)(api[_-]?key|token|secret|password)\s*[:=]\s*['"]*([a-zA-Z0-9_\-\.]{20,})['"]*`),
	
	// Email addresses
	regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`),
	
	// Base64 tokens (likely JWT or auth tokens)
	regexp.MustCompile(`[A-Za-z0-9+/]{50,}={0,2}`),
	
	// UUIDs (might be sensitive user/org IDs)
	regexp.MustCompile(`[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`),
	
	// Credit card numbers
	regexp.MustCompile(`\b\d{4}[\s-]?\d{4}[\s-]?\d{4}[\s-]?\d{4}\b`),
	
	// Social Security Numbers
	regexp.MustCompile(`\b\d{3}-\d{2}-\d{4}\b`),
	
	// IP addresses (might be internal)
	regexp.MustCompile(`\b(?:[0-9]{1,3}\.){3}[0-9]{1,3}\b`),
}

// SanitizeString removes or masks sensitive data from a string
func (sl *SecureLogger) SanitizeString(input string) string {
	sanitized := input
	
	for _, pattern := range SensitiveDataPatterns {
		sanitized = pattern.ReplaceAllStringFunc(sanitized, func(match string) string {
			// Different sanitization based on pattern type
			if strings.Contains(strings.ToLower(match), "@") {
				// Email address - mask middle part
				parts := strings.Split(match, "@")
				if len(parts) == 2 {
					local := parts[0]
					domain := parts[1]
					if len(local) > 2 {
						local = local[:1] + "*****" + local[len(local)-1:]
					}
					return local + "@" + domain
				}
			}
			
			// For other sensitive data, show only first and last 4 characters
			if len(match) > 8 {
				return match[:4] + "*****" + match[len(match)-4:]
			} else if len(match) > 4 {
				return match[:2] + "***" + match[len(match)-2:]
			}
			
			return "***REDACTED***"
		})
	}
	
	return sanitized
}

// MaskUserID masks user IDs for logging while preserving some identifying info
func (sl *SecureLogger) MaskUserID(userID string) string {
	if userID == "" {
		return "unknown"
	}
	if len(userID) <= 8 {
		return "user-***"
	}
	return "user-" + userID[:4] + "***" + userID[len(userID)-4:]
}

// MaskEmail masks email addresses for logging
func (sl *SecureLogger) MaskEmail(email string) string {
	if email == "" {
		return "unknown@domain"
	}
	
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "invalid@email"
	}
	
	local := parts[0]
	domain := parts[1]
	
	if len(local) <= 2 {
		return "*@" + domain
	}
	
	return local[:1] + "***@" + domain
}

// LogAuthSuccess logs successful authentication without exposing sensitive data
func (sl *SecureLogger) LogAuthSuccess(userID, email, orgID string) {
	log.Printf("Authentication successful - User: %s, Org: %s", 
		sl.MaskUserID(userID), 
		sl.MaskUserID(orgID))
}

// LogAuthFailure logs authentication failures securely
func (sl *SecureLogger) LogAuthFailure(email, reason string) {
	log.Printf("Authentication failed - Email: %s, Reason: %s", 
		sl.MaskEmail(email), 
		reason)
}

// LogVerificationPIN logs PIN-related events securely
func (sl *SecureLogger) LogVerificationPIN(email string, action string) {
	if sl.isDevelopment {
		log.Printf("Development mode: %s for email %s", action, sl.MaskEmail(email))
	} else {
		log.Printf("Email verification %s initiated", action)
	}
}

// LogAuthCookie logs cookie-related events without exposing cookie content
func (sl *SecureLogger) LogAuthCookie(action string) {
	log.Printf("Auth cookie %s", action)
}

// LogRequest logs HTTP requests securely
func (sl *SecureLogger) LogRequest(r *http.Request, duration time.Duration) {
	// Skip logging for monitoring endpoints
	if r.URL.Path == "/health" || r.URL.Path == "/version" {
		return
	}
	
	// Sanitize URL path to remove potential sensitive data
	sanitizedPath := sl.SanitizeURLPath(r.URL.Path)
	
	// Log basic request info without sensitive details
	log.Printf("Request: %s %s (duration: %v)", r.Method, sanitizedPath, duration)
}

// SanitizeURLPath removes sensitive data from URL paths
func (sl *SecureLogger) SanitizeURLPath(path string) string {
	// Replace UUIDs in paths with placeholder
	uuidPattern := regexp.MustCompile(`/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`)
	sanitized := uuidPattern.ReplaceAllString(path, "/[ID]")
	
	// Replace other potential IDs (long alphanumeric strings)
	idPattern := regexp.MustCompile(`/[a-zA-Z0-9]{20,}`)
	sanitized = idPattern.ReplaceAllString(sanitized, "/[ID]")
	
	return sanitized
}

// LogError logs errors securely without exposing sensitive information
func (sl *SecureLogger) LogError(err error, context string) {
	if err == nil {
		return
	}
	
	errorMessage := err.Error()
	sanitizedError := sl.SanitizeString(errorMessage)
	
	log.Printf("Error in %s: %s", context, sanitizedError)
}

// LogSecurityEvent logs security-related events
func (sl *SecureLogger) LogSecurityEvent(event, details string) {
	sanitizedDetails := sl.SanitizeString(details)
	log.Printf("SECURITY EVENT: %s - %s", event, sanitizedDetails)
}

// LogAPIKeyEvent logs API key related events without exposing the key
func (sl *SecureLogger) LogAPIKeyEvent(event, keyName string) {
	log.Printf("API Key Event: %s for key configuration: %s", event, keyName)
}

// LogDatabaseEvent logs database events securely
func (sl *SecureLogger) LogDatabaseEvent(event, details string) {
	// Remove any potential connection strings or sensitive data
	sanitizedDetails := sl.SanitizeString(details)
	log.Printf("Database: %s - %s", event, sanitizedDetails)
}

// SetupSecureLogging initializes secure logging for the application
func SetupSecureLogging(isDevelopment bool) *SecureLogger {
	// Configure the default logger to include milliseconds and short file paths
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	
	// Add security notice to log output
	if isDevelopment {
		log.Println("🔒 Secure logging initialized (Development Mode)")
	} else {
		log.Println("🔒 Secure logging initialized (Production Mode)")
	}
	
	return NewSecureLogger(isDevelopment)
}

// StructuredLogEntry represents a structured log entry
type StructuredLogEntry struct {
	Timestamp string                 `json:"timestamp"`
	Level     string                 `json:"level"`
	Event     string                 `json:"event"`
	Details   map[string]interface{} `json:"details"`
}

// LogStructured logs structured data while sanitizing sensitive fields
func (sl *SecureLogger) LogStructured(level, event string, details map[string]interface{}) {
	// Sanitize sensitive fields in the details map
	sanitizedDetails := make(map[string]interface{})
	for key, value := range details {
		if sl.isSensitiveField(key) {
			if strValue, ok := value.(string); ok {
				sanitizedDetails[key] = sl.SanitizeString(strValue)
			} else {
				sanitizedDetails[key] = "***REDACTED***"
			}
		} else {
			sanitizedDetails[key] = value
		}
	}
	
	entry := StructuredLogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Level:     level,
		Event:     event,
		Details:   sanitizedDetails,
	}
	
	log.Printf("[%s] %s: %+v", entry.Level, entry.Event, entry.Details)
}

// isSensitiveField checks if a field name indicates sensitive data
func (sl *SecureLogger) isSensitiveField(fieldName string) bool {
	sensitiveFields := []string{
		"password", "token", "secret", "key", "pin", "cookie", 
		"email", "ssn", "credit_card", "api_key", "auth",
		"authorization", "bearer", "jwt", "session",
	}
	
	fieldLower := strings.ToLower(fieldName)
	for _, sensitive := range sensitiveFields {
		if strings.Contains(fieldLower, sensitive) {
			return true
		}
	}
	
	return false
}