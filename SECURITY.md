# 🛡️ PDw-X Security Guide

## 🔒 Security Overview

PDw-X implements enterprise-grade security measures designed to protect sensitive development workflows, AI model interactions, and user data. This comprehensive guide covers security architecture, best practices, and compliance standards.

## 🎯 Security Principles

### Defense in Depth
- **Multi-layered Protection**: Multiple security controls at different levels
- **Principle of Least Privilege**: Minimal access rights for users and processes
- **Zero Trust Architecture**: Verify all requests regardless of source
- **Secure by Default**: Security-first configuration and deployment

### Threat Model
- **Data Protection**: Source code, AI interactions, and user credentials
- **Infrastructure Security**: Server hardening and network protection
- **Supply Chain Security**: Dependency and build pipeline protection
- **Privacy Protection**: User data and development pattern privacy

## 🏗️ Security Architecture

### System Components
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   CLI Client    │    │   Web Dashboard │    │   Mobile App    │
│   (Local Auth)  │◄──►│   (OAuth 2.0)   │◄──►│   (JWT Tokens)  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────────────────────────────────────────────────┐
│                    API Gateway + Load Balancer              │
│            ✓ Rate Limiting  ✓ DDoS Protection               │
│            ✓ SSL/TLS        ✓ Request Validation             │
└─────────────────────────────────────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────────────────────────────┐
│                      PDw-X Server                           │
│    ✓ Input Validation    ✓ Secure Headers                   │
│    ✓ RBAC               ✓ Audit Logging                     │
│    ✓ Encryption         ✓ Session Management               │
└─────────────────────────────────────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────────────────────────────┐
│                PostgreSQL 17.5 Database                     │
│    ✓ Encryption at Rest  ✓ Row-Level Security               │
│    ✓ Connection Encryption ✓ Backup Encryption             │
└─────────────────────────────────────────────────────────────┘
```

## 🔐 Authentication & Authorization

### Authentication Methods

#### 1. JWT Token Authentication
```go
// JWT configuration with security hardening
type JWTConfig struct {
    Secret          []byte        `json:"-"`          // Never log
    Algorithm       string        `json:"algorithm"`   // RS256 recommended
    ExpirationTime  time.Duration `json:"expires_in"`  // 15 minutes default
    RefreshTime     time.Duration `json:"refresh_in"`  // 7 days default
    Issuer          string        `json:"issuer"`
    Audience        []string      `json:"audience"`
}
```

**Security Features:**
- Short-lived access tokens (15 minutes)
- Secure refresh token rotation
- RS256 asymmetric signing
- Token revocation support
- Device fingerprinting

#### 2. OAuth 2.0 / OIDC Integration
```yaml
# OAuth configuration
oauth:
  providers:
    github:
      client_id: "${GITHUB_CLIENT_ID}"
      client_secret: "${GITHUB_CLIENT_SECRET}"
      scopes: ["user:email", "read:org"]
    google:
      client_id: "${GOOGLE_CLIENT_ID}"
      client_secret: "${GOOGLE_CLIENT_SECRET}"
      scopes: ["openid", "profile", "email"]
```

#### 3. API Key Management
```go
// Secure API key structure
type APIKey struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    HashedKey   string    `json:"-"`           // Never expose
    Scopes      []string  `json:"scopes"`
    LastUsed    time.Time `json:"last_used"`
    ExpiresAt   time.Time `json:"expires_at"`
    IPWhitelist []string  `json:"ip_whitelist,omitempty"`
}
```

### Authorization (RBAC)

#### Role Definitions
```yaml
roles:
  admin:
    description: "Full system access"
    permissions: ["*"]
    
  developer:
    description: "Development and deployment access"
    permissions:
      - "plans:read"
      - "plans:write"
      - "models:use"
      - "context:load"
      
  viewer:
    description: "Read-only access"
    permissions:
      - "plans:read"
      - "context:read"
      
  agent:
    description: "AI agent execution"
    permissions:
      - "models:execute"
      - "context:analyze"
```

## 🔒 Data Protection

### Encryption Standards

#### Data at Rest
- **Database Encryption**: AES-256 encryption for PostgreSQL
- **File System Encryption**: Full disk encryption required
- **Backup Encryption**: AES-256 encrypted backups with key rotation
- **Secret Management**: HashiCorp Vault or AWS Secrets Manager

#### Data in Transit
```go
// TLS configuration
tlsConfig := &tls.Config{
    MinVersion:               tls.VersionTLS13,
    CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
    PreferServerCipherSuites: true,
    CipherSuites: []uint16{
        tls.TLS_AES_256_GCM_SHA384,
        tls.TLS_CHACHA20_POLY1305_SHA256,
        tls.TLS_AES_128_GCM_SHA256,
    },
}
```

### Data Classification
- **Public**: Documentation, open source code
- **Internal**: Configuration, non-sensitive logs
- **Confidential**: User code, AI model responses
- **Restricted**: Authentication credentials, API keys

## 🛡️ Input Validation & Sanitization

### Comprehensive Validation Framework
```go
// Multi-layer validation
type ValidationFramework struct {
    Syntax   SyntaxValidator   // Format validation
    Semantic SemanticValidator // Business logic validation
    Security SecurityValidator // Security threat detection
    Size     SizeValidator     // Resource limit validation
}

// Example validation implementation
func (v *ValidationFramework) ValidateRequest(req *Request) error {
    if err := v.Syntax.Validate(req); err != nil {
        return fmt.Errorf("syntax validation failed: %w", err)
    }
    
    if err := v.Security.ScanForThreats(req); err != nil {
        return fmt.Errorf("security validation failed: %w", err)
    }
    
    if err := v.Size.CheckLimits(req); err != nil {
        return fmt.Errorf("size validation failed: %w", err)
    }
    
    return v.Semantic.ValidateBusinessRules(req)
}
```

### SQL Injection Prevention
```go
// Parameterized queries only
func GetUserByID(db *sql.DB, userID string) (*User, error) {
    query := `SELECT id, name, email FROM users WHERE id = $1`
    
    var user User
    err := db.QueryRow(query, userID).Scan(&user.ID, &user.Name, &user.Email)
    if err != nil {
        return nil, fmt.Errorf("failed to get user: %w", err)
    }
    
    return &user, nil
}
```

### XSS Protection
```go
// HTML sanitization
import "github.com/microcosm-cc/bluemonday"

func SanitizeUserInput(input string) string {
    p := bluemonday.UGCPolicy()
    return p.Sanitize(input)
}
```

## 🔍 Security Headers

### HTTP Security Headers
```go
// Complete security header implementation
func SecurityHeadersMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Prevent MIME type sniffing
        w.Header().Set("X-Content-Type-Options", "nosniff")
        
        // Prevent clickjacking
        w.Header().Set("X-Frame-Options", "DENY")
        
        // XSS protection
        w.Header().Set("X-XSS-Protection", "1; mode=block")
        
        // HSTS (HTTPS enforcement)
        w.Header().Set("Strict-Transport-Security", 
            "max-age=31536000; includeSubDomains; preload")
        
        // CSP (Content Security Policy)
        w.Header().Set("Content-Security-Policy", 
            "default-src 'self'; "+
            "script-src 'self' 'unsafe-inline'; "+
            "style-src 'self' 'unsafe-inline'; "+
            "img-src 'self' data: https:; "+
            "connect-src 'self' wss:; "+
            "frame-ancestors 'none'")
        
        // Referrer policy
        w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
        
        // Permissions policy
        w.Header().Set("Permissions-Policy", 
            "camera=(), microphone=(), geolocation=()")
        
        next.ServeHTTP(w, r)
    })
}
```

## 📊 Secure Logging & Monitoring

### Logging Framework
```go
// Secure logging with sensitive data redaction
type SecureLogger struct {
    logger     *logrus.Logger
    redactor   *SensitiveDataRedactor
    encryptor  *LogEncryptor
}

type SensitiveDataRedactor struct {
    patterns map[string]*regexp.Regexp
}

func NewSensitiveDataRedactor() *SensitiveDataRedactor {
    return &SensitiveDataRedactor{
        patterns: map[string]*regexp.Regexp{
            "api_key":      regexp.MustCompile(`(api[_-]?key[_-]?=\s*)([a-zA-Z0-9]{20,})`),
            "password":     regexp.MustCompile(`(password[_-]?=\s*)([^\s]+)`),
            "jwt_token":    regexp.MustCompile(`(Bearer\s+)([a-zA-Z0-9\-_]+\.[a-zA-Z0-9\-_]+\.[a-zA-Z0-9\-_]+)`),
            "email":        regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`),
            "credit_card":  regexp.MustCompile(`\b\d{4}[\s-]?\d{4}[\s-]?\d{4}[\s-]?\d{4}\b`),
        },
    }
}

func (r *SensitiveDataRedactor) Redact(message string) string {
    for name, pattern := range r.patterns {
        message = pattern.ReplaceAllString(message, fmt.Sprintf("${1}[REDACTED_%s]", strings.ToUpper(name)))
    }
    return message
}
```

### Audit Logging
```go
// Security event audit logging
type AuditEvent struct {
    Timestamp    time.Time            `json:"timestamp"`
    EventType    string              `json:"event_type"`
    UserID       string              `json:"user_id,omitempty"`
    SessionID    string              `json:"session_id,omitempty"`
    IPAddress    string              `json:"ip_address"`
    UserAgent    string              `json:"user_agent"`
    Resource     string              `json:"resource"`
    Action       string              `json:"action"`
    Result       string              `json:"result"` // SUCCESS, FAILURE, ERROR
    Details      map[string]interface{} `json:"details,omitempty"`
    RiskScore    int                 `json:"risk_score"` // 0-100
}

// Audit event types
const (
    EventTypeLogin           = "user.login"
    EventTypeLogout          = "user.logout"
    EventTypeAPIKeyCreated   = "api_key.created"
    EventTypeAPIKeyRevoked   = "api_key.revoked"
    EventTypePermissionDenied = "access.denied"
    EventTypeDataAccess      = "data.access"
    EventTypeConfigChange    = "config.change"
    EventTypeSecurityAlert   = "security.alert"
)
```

## 🚨 Incident Response

### Security Incident Classification
- **P0 (Critical)**: Active security breach, data exposure
- **P1 (High)**: Potential security vulnerability, service compromise
- **P2 (Medium)**: Security policy violation, suspicious activity
- **P3 (Low)**: Security configuration issue, minor policy violation

### Incident Response Workflow
1. **Detection**: Automated monitoring and manual reporting
2. **Assessment**: Risk evaluation and impact analysis
3. **Containment**: Immediate threat mitigation
4. **Investigation**: Root cause analysis and evidence collection
5. **Recovery**: System restoration and security improvements
6. **Lessons Learned**: Post-incident review and process updates

### Automated Response Actions
```go
// Automated security responses
type SecurityAutomation struct {
    RateLimiter    *rate.Limiter
    BanList        sync.Map
    AlertManager   *AlertManager
}

func (s *SecurityAutomation) HandleSecurityEvent(event AuditEvent) {
    switch event.EventType {
    case EventTypePermissionDenied:
        if event.RiskScore > 70 {
            s.TemporaryBan(event.IPAddress, 30*time.Minute)
            s.AlertManager.SendAlert("High risk access denied", event)
        }
        
    case EventTypeSecurityAlert:
        if event.RiskScore > 90 {
            s.EmergencyLockdown(event.UserID)
            s.AlertManager.SendCriticalAlert("Security breach detected", event)
        }
    }
}
```

## 🔧 Security Configuration

### PostgreSQL Security Hardening
```sql
-- Database security configuration
-- Enable SSL/TLS
ssl = on
ssl_cert_file = '/etc/ssl/certs/server.crt'
ssl_key_file = '/etc/ssl/private/server.key'
ssl_ca_file = '/etc/ssl/certs/ca.crt'

-- Authentication
authentication_timeout = 10s
password_encryption = scram-sha-256

-- Connection limits
max_connections = 100
superuser_reserved_connections = 3

-- Logging
log_connections = on
log_disconnections = on
log_statement = 'all'
log_line_prefix = '%t [%p]: [%l-1] user=%u,db=%d,app=%a,client=%h '

-- Row-level security
row_security = on
```

### Environment Configuration
```bash
# Security environment variables
export PLANDEX_SECURITY_MODE=strict
export PLANDEX_AUDIT_LOGGING=enabled
export PLANDEX_ENCRYPTION_AT_REST=enabled
export PLANDEX_SESSION_TIMEOUT=900  # 15 minutes
export PLANDEX_MAX_LOGIN_ATTEMPTS=3
export PLANDEX_RATE_LIMIT_REQUESTS=100
export PLANDEX_RATE_LIMIT_WINDOW=60  # per minute
```

## 🛠️ Security Tools & Scripts

### Security Scanning
```bash
#!/bin/bash
# security-scan.sh - Comprehensive security scanning

echo "Running PDw-X Security Scan..."

# Vulnerability scanning
echo "1. Scanning for vulnerabilities..."
govulncheck ./...

# Dependency security audit
echo "2. Auditing dependencies..."
go list -json -deps ./... | nancy sleuth

# Static security analysis
echo "3. Static security analysis..."
gosec ./...

# Secret scanning
echo "4. Scanning for secrets..."
truffleHog --regex --entropy=False .

# Container security scanning
echo "5. Container security scan..."
docker run --rm -v "$PWD":/app clair-scanner:latest

# TLS configuration test
echo "6. TLS configuration test..."
testssl.sh --quiet --color 0 localhost:8099

echo "Security scan completed."
```

### Security Validation
```bash
#!/bin/bash
# validate-security.sh - Security configuration validation

# Check file permissions
find . -type f -perm -o+w | while read file; do
    echo "WARNING: World-writable file: $file"
done

# Check for hardcoded secrets
grep -r "password\|secret\|key" --include="*.go" . | grep -v test

# Validate TLS configuration
openssl s_client -connect localhost:8099 -tls1_3 -quiet

# Check database permissions
psql -c "SELECT rolname, rolsuper, rolcreaterole, rolcreatedb FROM pg_roles;"
```

## 📋 Security Checklist

### Development Security
- [ ] Input validation implemented for all endpoints
- [ ] SQL injection prevention (parameterized queries)
- [ ] XSS protection and output encoding
- [ ] CSRF protection for state-changing operations
- [ ] Secure authentication and session management
- [ ] Proper error handling (no information leakage)
- [ ] Security headers implemented
- [ ] Rate limiting and DDoS protection
- [ ] Audit logging for security events
- [ ] Secrets management (no hardcoded credentials)

### Infrastructure Security
- [ ] TLS 1.3 encryption for all communications
- [ ] Database encryption at rest
- [ ] Regular security updates applied
- [ ] Network segmentation and firewall rules
- [ ] Intrusion detection and monitoring
- [ ] Backup encryption and security
- [ ] Access controls and privilege management
- [ ] Security monitoring and alerting
- [ ] Incident response procedures
- [ ] Regular security assessments

### Compliance Requirements
- [ ] GDPR compliance (data protection and privacy)
- [ ] SOC 2 Type II controls implementation
- [ ] HIPAA compliance for healthcare data
- [ ] PCI DSS compliance for payment data
- [ ] ISO 27001 security management
- [ ] Data retention and deletion policies
- [ ] Privacy policy and user consent
- [ ] Security training and awareness
- [ ] Third-party security assessments
- [ ] Regular compliance audits

## 🆘 Security Support

### Vulnerability Reporting
**Email**: [security@pdw-x.dev](mailto:security@pdw-x.dev)

**PGP Key**: [Download Public Key](https://pdw-x.dev/.well-known/security.txt)

### Security Advisory Process
1. **Report**: Send detailed vulnerability report
2. **Acknowledgment**: Response within 24 hours
3. **Assessment**: Risk evaluation and timeline
4. **Fix**: Security patch development
5. **Disclosure**: Coordinated vulnerability disclosure
6. **Recognition**: Security researcher acknowledgment

### Emergency Contacts
- **Security Team**: security@pdw-x.dev
- **Incident Response**: incident@pdw-x.dev
- **Emergency Hotline**: +1-555-SECURITY

---

## 📚 Additional Resources

- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [NIST Cybersecurity Framework](https://www.nist.gov/cyberframework)
- [CIS Security Controls](https://www.cisecurity.org/controls/)
- [Go Security Guide](https://github.com/Checkmarx/Go-SCP)
- [PostgreSQL Security](https://www.postgresql.org/docs/current/security.html)

---

*This security guide is regularly updated to reflect the latest security practices and threat landscape. Last updated: $(date)*