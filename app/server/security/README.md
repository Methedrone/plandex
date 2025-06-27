# Plandex Security Middleware

This package provides comprehensive web security middleware for the Plandex server, implementing modern security best practices and standards.

## Features

### 🛡️ **Security Headers**
- **X-Frame-Options**: Prevents clickjacking attacks
- **X-Content-Type-Options**: Prevents MIME sniffing
- **X-XSS-Protection**: Disables legacy XSS filter (prevents bypasses)
- **Referrer-Policy**: Controls referrer information sent with requests
- **Cross-Origin-***: Modern cross-origin isolation policies
- **Permissions-Policy**: Restricts dangerous browser features

### 🚫 **Rate Limiting**
- Configurable requests per second (default: 100/sec)
- Burst capacity (default: 200 requests)
- Automatic retry headers
- Health check exemption

### 🌐 **CORS (Cross-Origin Resource Sharing)**
- Trusted origin validation
- Secure credential handling
- Preflight request support
- Development and production configurations

### 🔒 **Content Security Policy (CSP)**
- Nonce-based script security
- Strict policies for production
- Relaxed policies for development
- Context-aware CSP nonce injection

### 🔐 **HTTP Strict Transport Security (HSTS)**
- Production-only activation
- Subdomains inclusion
- Preload directive support
- TLS detection

### ✅ **Input Validation**
- Request size limiting (10MB default)
- Content-type validation
- Host header injection prevention
- Path traversal protection

## Usage

### Basic Setup

```go
import "plandex-server/security"

// Use default configuration
securityConfig := security.DefaultSecurityConfig()
securityMiddleware := security.NewSecurityMiddleware(securityConfig)
handler = securityMiddleware.Apply(handler)
```

### Custom Configuration

```go
config := &security.SecurityConfig{
    TrustedOrigins: []string{"https://yourdomain.com"},
    IsDevelopment:  false,
    RateLimit:      rate.Limit(50),  // 50 req/sec
    RateBurst:      100,             // 100 req burst
    MaxRequestSize: 5 * 1024 * 1024, // 5MB
    EnableCSP:      true,
    EnableHSTS:     true,
    EnableCORS:     true,
}

securityMiddleware := security.NewSecurityMiddleware(config)
```

## Security Configuration

### Development Mode
- **CORS**: Allows localhost origins
- **CSP**: Relaxed policy with unsafe-inline and unsafe-eval
- **HSTS**: Disabled
- **Rate Limiting**: Standard limits

### Production Mode
- **CORS**: Strict origin validation
- **CSP**: Strict policy with nonce-based execution
- **HSTS**: Enabled with preload
- **Rate Limiting**: Enforced limits

## Middleware Order

The security middleware applies protections in the optimal order:

1. **Rate Limiting** (outermost)
2. **Input Validation**
3. **CORS**
4. **HSTS**
5. **CSP**
6. **Security Headers** (innermost)

## Environment Variables

- `GOENV=development`: Enables development mode with relaxed security

## Security Headers Applied

| Header | Value | Purpose |
|--------|-------|---------|
| X-Frame-Options | SAMEORIGIN | Prevent clickjacking |
| X-Content-Type-Options | nosniff | Prevent MIME sniffing |
| X-XSS-Protection | 0 | Disable legacy XSS filter |
| Referrer-Policy | strict-origin-when-cross-origin | Control referrer info |
| Cross-Origin-Opener-Policy | same-origin | Cross-origin isolation |
| Cross-Origin-Embedder-Policy | require-corp | Resource isolation |
| Cross-Origin-Resource-Policy | cross-origin | Resource sharing |
| Permissions-Policy | camera=(), microphone=()... | Restrict features |
| Cache-Control | no-cache, no-store | Secure API endpoints |

## CSP Nonce Usage

For templates that need to use the CSP nonce:

```go
nonce := security.GetCSPNonce(r)
// Use nonce in your templates: <script nonce="{{.Nonce}}">
```

## Rate Limiting Headers

When rate limiting is triggered:

- `Retry-After`: Seconds until retry allowed
- `X-RateLimit-Limit`: Configured rate limit
- `X-RateLimit-Remaining`: Remaining requests
- `X-RateLimit-Reset`: Reset timestamp

## Performance Optimizations

- **Minimal Memory Allocation**: Efficient header setting
- **Skip Health Checks**: No rate limiting for monitoring
- **Context Reuse**: Efficient nonce context handling
- **Configurable Limits**: Adjustable for hardware constraints

## Security Best Practices

1. **Keep Dependencies Updated**: Regularly update security dependencies
2. **Monitor Rate Limits**: Adjust based on usage patterns
3. **Review Trusted Origins**: Regularly audit CORS origins
4. **Test CSP Policies**: Validate CSP doesn't break functionality
5. **Enable HSTS Preload**: Submit domain to HSTS preload list

## Threat Mitigation

This middleware protects against:

- **DoS/DDoS**: Rate limiting and request size limits
- **Clickjacking**: X-Frame-Options header
- **XSS**: Content Security Policy
- **CSRF**: CORS and SameSite protections
- **Information Disclosure**: Secure headers and error handling
- **Path Traversal**: Input validation
- **Host Header Injection**: Host validation
- **MIME Sniffing**: X-Content-Type-Options
- **Protocol Downgrade**: HSTS

## Monitoring and Logging

The middleware integrates with the existing logging system:

- Rate limit violations are logged
- Security header violations can be monitored
- CSP violations can be reported (when configured)

## Future Enhancements

Planned security features:

- **IP-based Rate Limiting**: Per-IP rate limits
- **Geolocation Blocking**: Country-based restrictions
- **API Key Validation**: Enhanced authentication
- **Request Signing**: Message integrity validation
- **Security Monitoring**: Real-time threat detection