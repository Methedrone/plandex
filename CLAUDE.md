# Plandex - Claude Code Development Guide

## 🎯 PROJECT OVERVIEW

**Plandex** is a terminal-based AI development tool designed for large coding tasks that span multiple files and complex implementations. It's built for real-world projects with intelligent context management, multi-model AI integration, and automated debugging capabilities.

### Core Capabilities
- **Large Task Management**: Handle up to 2M tokens of context directly
- **Multi-File Operations**: Plan and execute changes across dozens of files
- **AI Model Integration**: Combine models from Anthropic, OpenAI, Google, and open source
- **Context Intelligence**: Tree-sitter project maps supporting 30+ languages
- **Automated Debugging**: Built-in debugging for commands, builds, tests, and deployments

### Architecture Overview
- **Language**: Go 1.23+ (multi-module workspace)
- **Database**: PostgreSQL with optimized connection pooling
- **Architecture**: Client-server model (CLI + Server components)
- **Deployment**: Docker-based with multi-stage optimization
- **Context Management**: Intelligent file loading and token optimization

---

## 🚨 CRITICAL IMPLEMENTATION PHASES

## Phase 1: Security & Dependencies (IMMEDIATE PRIORITY)

### 🔴 CRITICAL SECURITY ISSUES
**MUST BE COMPLETED BEFORE ANY OTHER WORK**

#### Go Runtime Security Upgrade
- **Current**: Go 1.23.3 (CRITICAL - 7 security patches missing)  
- **Target**: Go 1.23.10+
- **Impact**: Multiple CVEs affecting all modules
- **Files**: `/app/cli/go.mod`, `/app/server/go.mod`, `/app/shared/go.mod`, Docker files

#### PostgreSQL Security Fix
- **Current**: `postgres:latest` (HIGH RISK - CVE-2025-1094)
- **Target**: Pin to `postgres:17.5` with security configuration
- **Impact**: SQL injection vulnerability
- **Files**: `docker-compose.yml`, database connection strings

#### Security Headers Implementation
- **Missing**: CORS, CSP, HSTS, X-Frame-Options
- **Required**: Comprehensive security middleware
- **Files**: Create `/app/server/middleware/security.go`

#### Input Validation System
- **Current**: Basic XML escaping only
- **Required**: Comprehensive validation framework
- **Files**: Create `/app/server/validation/` package

#### API Key Protection
- **Risk**: Potential credential exposure in logs
- **Required**: Secure logging with sensitive data redaction
- **Files**: Create `/app/server/logging/` package

### Phase 1 Success Criteria
- ✅ All critical security vulnerabilities resolved (0 high-risk findings)
- ✅ Security scanning passes with comprehensive coverage
- ✅ All modules build and tests pass after upgrades
- ✅ Performance regression <5%
- ✅ Rollback procedures tested and documented

---

## Phase 2: Performance & Infrastructure Optimization

### 🚀 PERFORMANCE TARGETS
Optimized specifically for MacBook 2012 constraints (4-8GB RAM, dual/quad-core CPU)

#### Database Layer Optimization
- **Target**: 50-70% improvement in query response times
- **Implementation**: Optimized connection pooling, strategic indexing, query pattern optimization
- **MacBook 2012**: Conservative pool settings (8 max connections vs 50 production)

#### Memory Management Optimization
- **Target**: 30-40% reduction in memory usage
- **Implementation**: Object pooling, streaming optimization, GC tuning
- **MacBook 2012**: Aggressive GC (GOGC=50), memory limits (512MB)

#### Build & Docker Optimization
- **Target**: 40% faster build times
- **Implementation**: Multi-stage Docker builds, UPX compression, build caching
- **Results**: 60%+ smaller images, improved layer caching

#### Caching Implementation
- **Target**: 40-60% reduction in repeated operations
- **Implementation**: Multi-level caching (L1 memory + L2 Redis), context caching
- **Hit Rate**: >70% for frequent operations

### MacBook 2012 Specific Optimizations
```bash
# Environment variables for optimal performance
export GOGC=50                    # More aggressive GC
export GOMEMLIMIT=512MiB         # Memory limit
export GOMAXPROCS=4              # CPU cores
export DOCKER_DEFAULT_PLATFORM=linux/amd64
```

---

## Phase 3-5: Advanced Feature Development

### Phase 3: Testing & DevOps (Week 3-4)
- **Target**: 80%+ test coverage with automated quality gates
- **CI/CD**: Complete GitHub Actions pipeline with security scanning
- **Monitoring**: Structured logging, metrics collection, health checks

### Phase 4: Feature Enhancement (Week 4-8)
- **Web Dashboard**: Modern React interface with real-time collaboration
- **AI Assistant**: Code quality analysis with intelligent suggestions
- **IDE Integration**: VS Code extension, JetBrains plugins
- **Git Integration**: Automated PR workflows, semantic versioning

### Phase 5: Next-Gen Capabilities (Week 8-12)
- **Multi-Modal AI**: Vision processing for screenshots/diagrams
- **Progressive Web App**: Mobile-responsive with offline capabilities
- **Enterprise Security**: Advanced RBAC, audit logging, compliance
- **Analytics**: Predictive insights, ML-driven recommendations

---

## 🔧 CLAUDE CODE WORKFLOW INTEGRATION

### MCP Servers Strategy
```bash
# ALWAYS start with current documentation
use context7

# For complex multi-phase analysis
use SequentialThinking

# For comprehensive task management
use Task-Master with PRD-driven development
```

### TodoWrite Master Framework
This project requires systematic TodoWrite management across all 5 phases:

#### Phase 1 TodoWrite Example
```bash
# Phase 1 Critical Security Path
- [ ] Backup current codebase and create rollback plan
- [ ] Upgrade Go runtime to 1.23.10+ across all modules  
- [ ] Pin PostgreSQL to secure version 17.5
- [ ] Implement comprehensive security headers
- [ ] Create input validation framework
- [ ] Audit and secure API key handling
- [ ] Run security scans and validate fixes
- [ ] Update documentation and security procedures
```

#### Cross-Phase Dependencies
- **Phase 1 → Phase 2**: Security foundation required before performance optimization
- **Phase 2 → Phase 3**: Performance baseline needed for testing infrastructure
- **Phase 3 → Phase 4**: Quality gates required before feature development
- **Phase 4 → Phase 5**: Feature foundation needed for next-gen capabilities

---

## 🏗️ DEVELOPMENT ARCHITECTURE

### Module Structure
```
/app/
├── cli/              # Command-line interface
│   ├── go.mod       # CLI module dependencies
│   └── cmd/         # CLI commands
├── server/          # Core server application
│   ├── go.mod       # Server module dependencies
│   ├── main.go      # Server entry point
│   ├── db/          # Database layer
│   ├── handlers/    # HTTP handlers
│   └── middleware/  # HTTP middleware
├── shared/          # Shared utilities
│   ├── go.mod       # Shared module dependencies
│   └── types/       # Common types
└── docs/            # Docusaurus documentation site
```

### Key Design Patterns
- **Context Management**: Intelligent loading with tree-sitter parsing
- **Streaming Architecture**: Memory-efficient processing for large operations
- **Multi-Model Integration**: Flexible AI provider abstraction
- **Error Handling**: Comprehensive error recovery with rollback capabilities

---

## 🛡️ SECURITY & COMPLIANCE

### Security Framework
- **Authentication**: JWT-based with secure token handling
- **Authorization**: Role-based access control (RBAC)
- **Data Protection**: Encryption at rest and in transit
- **Audit Logging**: Comprehensive security event logging
- **Compliance**: SOC 2, GDPR, HIPAA ready (Phase 5)

### Security Headers Implementation
```go
// Essential security headers for all responses
w.Header().Set("X-Content-Type-Options", "nosniff")
w.Header().Set("X-Frame-Options", "DENY")  
w.Header().Set("X-XSS-Protection", "1; mode=block")
w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'...")
```

---

## ⚡ PERFORMANCE OPTIMIZATION

### Database Performance
- **Connection Pooling**: Optimized for MacBook 2012 (8 connections max)
- **Query Optimization**: Eliminate N+1 patterns, strategic indexing
- **Monitoring**: Real-time pool statistics and slow query detection

### Memory Management
- **Object Pooling**: Reuse buffers, JSON encoders, string builders
- **Streaming**: Chunked processing for large files and AI responses
- **GC Tuning**: Aggressive collection for memory-constrained environments

### Caching Strategy
```go
// Multi-level caching for optimal performance
L1 Cache: In-memory (1000 items, 10min TTL)
L2 Cache: Redis (optional, longer TTL)
Context Cache: Project contexts, file content, model responses
Hit Rate Target: >70% for frequent operations
```

---

## 📋 QUALITY ASSURANCE

### Testing Strategy
- **Unit Tests**: 80%+ coverage across all modules
- **Integration Tests**: Database, API, file system operations
- **E2E Tests**: Complete CLI workflows with real project scenarios
- **Performance Tests**: Benchmarking suite with regression detection
- **Security Tests**: Vulnerability scanning, penetration testing

### CI/CD Pipeline
```yaml
# GitHub Actions workflow stages
1. Code Quality: linting, formatting, security scanning
2. Testing: unit, integration, e2e test suites
3. Performance: benchmark comparison and regression detection
4. Security: vulnerability scanning, dependency audit
5. Build: optimized Docker images with multi-architecture support
6. Deploy: automated deployment with health checks
```

---

## 📊 MONITORING & METRICS

### Performance KPIs
- **Database**: 50-70% query performance improvement
- **Memory**: 30-40% usage reduction
- **Build**: 40% faster Docker builds
- **API**: 25-50% response time improvement
- **Cache**: >70% hit rate efficiency

### Monitoring Implementation
```go
// Comprehensive performance metrics
type PerformanceMetrics struct {
    Memory      MemoryStats      `json:"memory"`
    Database    DatabaseStats    `json:"database"`
    Cache       CacheStats       `json:"cache"`
    Runtime     RuntimeStats     `json:"runtime"`
    Uptime      time.Duration    `json:"uptime"`
}
```

---

## 🚀 DEPLOYMENT & OPERATIONS

### Docker Optimization
```dockerfile
# Multi-stage build optimized for performance
FROM golang:1.23.10-alpine AS builder
# Optimized build flags: -ldflags='-w -s' -a -installsuffix cgo
# UPX compression for 60%+ size reduction
FROM scratch AS runtime
# Minimal runtime with security hardening
```

### Development Environment
```bash
# MacBook 2012 optimized setup
cd /Users/Mahavir/code/plandex
docker-compose -f docker-compose.performance.yml up -d

# Development tools
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install golang.org/x/vuln/cmd/govulncheck@latest
```

---

## 🔄 DEVELOPMENT WORKFLOW

### Pre-Development Checklist
- [ ] `use context7` - Research latest best practices
- [ ] Run `SequentialThinking` for complex tasks
- [ ] Check `Task-Master` priority queue
- [ ] Verify current security status (Phase 1 complete)
- [ ] Confirm performance baseline (Phase 2 complete)

### Implementation Pattern
1. **Research Phase**: Use Context7 for current documentation and best practices
2. **Planning Phase**: Apply SequentialThinking for complex problem analysis
3. **Execution Phase**: TodoWrite with systematic progress tracking
4. **Validation Phase**: Testing, security scanning, performance benchmarking

### Code Quality Standards
- **Go Standards**: Follow effective Go practices, comprehensive error handling
- **Security**: Input validation, secure logging, credential protection
- **Performance**: Memory efficiency, database optimization, caching strategies
- **Testing**: Comprehensive coverage with automated quality gates

---

## 📚 DOCUMENTATION STRUCTURE

### Core Documentation
- **README.md**: Project overview, installation, quick start
- **API Documentation**: Generated from code with examples
- **Development Guide**: Setup, architecture, contributing guidelines
- **Deployment Guide**: Production deployment, monitoring, maintenance

### User Documentation (Docusaurus)
```bash
cd /Users/Mahavir/code/plandex/docs
yarn install      # Install dependencies
yarn start        # Local development server
yarn build        # Production build
```

---

## 🎯 SUCCESS CRITERIA

### Technical Excellence
- ✅ Zero critical security vulnerabilities
- ✅ 30-70% performance improvements across all metrics
- ✅ 80%+ test coverage with automated quality gates
- ✅ MacBook 2012 optimization validated
- ✅ Enterprise-grade security and compliance ready

### Innovation Leadership
- ✅ Multi-modal AI processing capabilities
- ✅ Real-time collaboration features
- ✅ Advanced analytics with ML-driven insights
- ✅ Progressive Web App with offline capabilities
- ✅ Cutting-edge developer experience

---

## 🚨 CRITICAL REMINDERS

### Security First
**NEVER proceed to Phase 2+ until Phase 1 security issues are completely resolved**
- Go 1.23.10+ upgrade is CRITICAL
- PostgreSQL CVE-2025-1094 is HIGH RISK
- Security headers are MANDATORY

### MacBook 2012 Optimization
**All performance work must consider resource constraints:**
- Memory: 4-8GB RAM limitation
- CPU: 2.5-2.9GHz dual/quad-core
- Storage: HDD/SSD hybrid considerations

### Quality Gates
**Every phase must pass quality validation before proceeding:**
- Security scanning clean
- Performance benchmarks met
- Test coverage targets achieved
- Documentation updated

---

*This guide provides the comprehensive framework for Claude Code to work effectively with the Plandex codebase, emphasizing security priorities, performance optimization for resource-constrained environments, and systematic development practices across all implementation phases.*