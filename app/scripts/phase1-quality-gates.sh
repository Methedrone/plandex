#!/bin/bash

# Phase 1 Quality Gates - Security Foundation Validation
# This script validates all Phase 1 security objectives are met

set -e

echo "🔒 Phase 1 Quality Gates - Security Foundation Validation"
echo "========================================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Results tracking
TOTAL_CHECKS=0
PASSED_CHECKS=0
FAILED_CHECKS=0

# Function to run a check
run_check() {
    local check_name="$1"
    local command="$2"
    
    TOTAL_CHECKS=$((TOTAL_CHECKS + 1))
    echo -n "Checking $check_name... "
    
    if eval "$command" > /dev/null 2>&1; then
        echo -e "${GREEN}✅ PASS${NC}"
        PASSED_CHECKS=$((PASSED_CHECKS + 1))
        return 0
    else
        echo -e "${RED}❌ FAIL${NC}"
        FAILED_CHECKS=$((FAILED_CHECKS + 1))
        return 1
    fi
}

# Function to run a security pattern check
check_security_pattern() {
    local description="$1"
    local pattern="$2"
    local files="$3"
    local should_exist="$4"
    
    TOTAL_CHECKS=$((TOTAL_CHECKS + 1))
    echo -n "Security Check: $description... "
    
    local found=$(grep -r "$pattern" $files 2>/dev/null | wc -l)
    
    if [ "$should_exist" = "true" ]; then
        if [ "$found" -gt 0 ]; then
            echo -e "${GREEN}✅ PASS${NC}"
            PASSED_CHECKS=$((PASSED_CHECKS + 1))
            return 0
        else
            echo -e "${RED}❌ FAIL${NC}"
            FAILED_CHECKS=$((FAILED_CHECKS + 1))
            return 1
        fi
    else
        if [ "$found" -eq 0 ]; then
            echo -e "${GREEN}✅ PASS${NC}"
            PASSED_CHECKS=$((PASSED_CHECKS + 1))
            return 0
        else
            echo -e "${RED}❌ FAIL (Found: $found)${NC}"
            FAILED_CHECKS=$((FAILED_CHECKS + 1))
            return 1
        fi
    fi
}

echo "📋 1. Go Runtime and Build Validation"
echo "-------------------------------------"

# Check Go version
run_check "Go 1.23.10 in all modules" "grep -r 'go 1.23.10' */go.mod"

# Check if modules can be downloaded (simulated)
run_check "CLI module structure" "test -f cli/go.mod && test -f cli/main.go"
run_check "Server module structure" "test -f server/go.mod && test -f server/main.go"
run_check "Shared module structure" "test -f shared/go.mod"

echo ""
echo "🐘 2. PostgreSQL Security Validation"
echo "------------------------------------"

# Check PostgreSQL pinning
run_check "PostgreSQL 17.5 pinned" "grep 'postgres:17.5-alpine' docker-compose.yml"

# Check SCRAM-SHA-256 configuration
run_check "SCRAM-SHA-256 auth configured" "grep 'scram-sha-256' docker-compose.yml"

# Check security configuration files
run_check "PostgreSQL security config exists" "test -f postgresql/postgresql.conf"
run_check "PostgreSQL auth config exists" "test -f postgresql/pg_hba.conf"

# Check SSL preference in connection string
run_check "SSL preference configured" "grep 'sslmode=prefer' docker-compose.yml"

echo ""
echo "🛡️ 3. Security Headers and Middleware Validation"
echo "-----------------------------------------------"

# Check security middleware exists
run_check "Security middleware implemented" "test -f server/security/middleware.go"

# Check security features
check_security_pattern "Rate limiting implemented" "rate.NewLimiter" "server/security/" true
check_security_pattern "CORS implementation" "Access-Control-Allow-Origin" "server/security/" true
check_security_pattern "Security headers set" "X-Frame-Options" "server/security/" true
check_security_pattern "CSP implementation" "Content-Security-Policy" "server/security/" true
check_security_pattern "HSTS implementation" "Strict-Transport-Security" "server/security/" true

echo ""
echo "🔐 4. API Key Security Validation"
echo "--------------------------------"

# Check API key security
run_check ".env.example template exists" "test -f .env.example"
run_check ".env properly ignored" "grep '^.env$' .gitignore"
run_check "Security alert documentation" "test -f SECURITY_ALERT.md"

# Check for API key exposure in version-controlled code files only
TOTAL_CHECKS=$((TOTAL_CHECKS + 1))
echo -n "Security Check: No API keys in version-controlled files... "
API_KEYS_IN_CODE=$(find . -name "*.go" -o -name "*.js" -o -name "*.ts" | grep -v "/node_modules/" | xargs grep -l "sk-or-v1-\|AIzaSy\|pplx-" 2>/dev/null | wc -l)
if [ "$API_KEYS_IN_CODE" -eq 0 ]; then
    echo -e "${GREEN}✅ PASS${NC}"
    PASSED_CHECKS=$((PASSED_CHECKS + 1))
else
    echo -e "${RED}❌ FAIL (Found in $API_KEYS_IN_CODE files)${NC}"
    FAILED_CHECKS=$((FAILED_CHECKS + 1))
fi

echo ""
echo "📝 5. Logging Security Validation"
echo "--------------------------------"

# Check secure logging implementation
run_check "Secure logging framework exists" "test -f server/security/logging.go"

# Check for sensitive data exposure in logs (should not exist)
check_security_pattern "No PIN logging in verification" "pin is.*for email" "server/" false
check_security_pattern "No cookie content logging" "setting auth cookie.*cookie" "server/" false
check_security_pattern "No API key env var exposure" "api key env var:" "server/" false

# Check for secure logging patterns (should exist)
check_security_pattern "Secure authentication logging" "Authentication successful.*User:" "server/" true
check_security_pattern "User data masking implemented" "MaskUserID\|MaskEmail" "server/security/" true

echo ""
echo "🚀 6. Docker Security Validation"
echo "-------------------------------"

# Check Docker security improvements
run_check "Multi-stage Dockerfile" "grep 'AS builder' server/Dockerfile"
run_check "Non-root user in Dockerfile" "grep 'USER appuser' server/Dockerfile"
run_check "Security options in compose" "grep 'no-new-privileges' docker-compose.yml"
run_check "Resource limits configured" "grep 'memory: 1G' docker-compose.yml"
run_check "Health checks configured" "grep 'healthcheck:' docker-compose.yml"

echo ""
echo "⚡ 7. Performance Optimization Validation"
echo "----------------------------------------"

# Check MacBook 2012 optimizations
check_security_pattern "Go GC optimization" "GOGC=50" "server/Dockerfile" true
check_security_pattern "Memory limit set" "GOMEMLIMIT=512MiB" "server/Dockerfile" true
check_security_pattern "CPU optimization" "GOMAXPROCS=4" "server/Dockerfile" true

# Check PostgreSQL optimizations
check_security_pattern "Memory optimization" "shared_buffers = 256MB" "postgresql/" true
check_security_pattern "Connection limits" "max_connections = 50" "postgresql/" true

echo ""
echo "📚 8. Documentation Validation"
echo "-----------------------------"

# Check documentation exists
run_check "Security middleware documentation" "test -f server/security/README.md"
run_check "CLAUDE.md exists" "test -f CLAUDE.md"
run_check "Security alert documentation" "test -f SECURITY_ALERT.md"

echo ""
echo "🔍 9. Input Validation Security"
echo "------------------------------"

# Check input validation features
check_security_pattern "Request size validation" "MaxRequestSize" "server/security/" true
check_security_pattern "Content type validation" "Content-Type" "server/security/" true
check_security_pattern "Host header validation" "Host.*header" "server/security/" true
check_security_pattern "Path traversal protection" "\.\..*//.*Invalid path" "server/security/" true

echo ""
echo "🏗️ 10. Integration Validation"
echo "----------------------------"

# Check middleware integration
check_security_pattern "Security middleware integrated" "securityMiddleware.Apply" "server/setup/" true
check_security_pattern "Secure logging integrated" "security.NewSecureLogger" "server/setup/" true

echo ""
echo "📊 Phase 1 Quality Gates Results"
echo "================================"

echo "Total Checks: $TOTAL_CHECKS"
echo -e "Passed: ${GREEN}$PASSED_CHECKS${NC}"
echo -e "Failed: ${RED}$FAILED_CHECKS${NC}"

PASS_RATE=$((PASSED_CHECKS * 100 / TOTAL_CHECKS))
echo "Pass Rate: $PASS_RATE%"

echo ""
if [ $FAILED_CHECKS -eq 0 ]; then
    echo -e "${GREEN}🎉 All Phase 1 Quality Gates PASSED!${NC}"
    echo -e "${GREEN}✅ Security foundation is ready for Phase 2${NC}"
    exit 0
else
    echo -e "${RED}❌ Phase 1 Quality Gates FAILED${NC}"
    echo -e "${YELLOW}⚠️  Please fix failing checks before proceeding to Phase 2${NC}"
    exit 1
fi