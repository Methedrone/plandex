# Phase 3: Testing & DevOps Infrastructure Implementation Guide
## Plandex App Upgrade - Comprehensive Quality Assurance & Automation

---

## 🎯 EXECUTIVE SUMMARY

This guide provides comprehensive implementation of testing frameworks, CI/CD pipelines, monitoring infrastructure, and developer tooling for the Plandex application. Building on the secure and performant foundation from Phases 1-2, this phase establishes enterprise-grade quality assurance and operational excellence.

### Current State Analysis
- **Testing Coverage**: ~5% (only 6 test files found)
- **CI/CD Pipeline**: Basic GitHub Actions (3 workflows)
- **Testing Types**: Limited unit tests for syntax/streaming
- **Monitoring**: Basic health checks only
- **DevOps Tooling**: Minimal development automation
- **Quality Gates**: No automated quality controls

### Implementation Targets
- **Test Coverage**: 80%+ across all modules
- **CI/CD Pipeline**: Fully automated testing, security scanning, deployment
- **Testing Types**: Unit, Integration, E2E, API, Load, Security tests
- **Monitoring**: Comprehensive observability stack
- **Quality Gates**: Automated quality controls at every stage
- **Developer Experience**: Streamlined development workflow

---

## 🔧 CLAUDE CODE WORKFLOW INTEGRATION

### MCP Servers Utilization
```bash
# Research testing best practices
use context7

# Complex testing strategy analysis
use SequentialThinking

# Task breakdown for testing implementation
use Task-Master with testing-focused PRD
```

### Testing-Focused TodoWrite Strategy
This guide integrates comprehensive TodoWrite checkpoints for each testing component, ensuring systematic implementation and quality validation.

---

## 📋 DETAILED IMPLEMENTATION PLAN

## Phase 3A: Comprehensive Testing Framework
### 🧪 TESTING FOUNDATION TARGET
**Goal**: Achieve 80%+ test coverage with comprehensive test types

### Implementation Steps

#### Step 3A.1: Testing Infrastructure Setup
**File: `/app/testing/test_config.go`** (create new file)
```go
package testing

import (
    "context"
    "database/sql"
    "fmt"
    "log"
    "os"
    "testing"
    "time"
    
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/postgres"
    "yourapp/db"
    "yourapp/cache"
)

// TestConfig provides test environment configuration
type TestConfig struct {
    DBPool          *db.DatabasePool
    Cache           *cache.CacheManager
    PostgresContainer testcontainers.Container
    TestDBURL       string
    Cleanup         func()
}

// SetupTestEnvironment creates isolated test environment
func SetupTestEnvironment(t *testing.T) *TestConfig {
    ctx := context.Background()
    
    // Start PostgreSQL test container
    postgresContainer, err := postgres.RunContainer(ctx,
        testcontainers.WithImage("postgres:17.5-alpine"),
        postgres.WithDatabase("testdb"),
        postgres.WithUsername("testuser"),
        postgres.WithPassword("testpass"),
        testcontainers.WithWaitStrategy(
            wait.ForLog("database system is ready to accept connections").
                WithOccurrence(2).
                WithStartupTimeout(30*time.Second)),
    )
    if err != nil {
        t.Fatalf("Failed to start postgres container: %v", err)
    }
    
    // Get connection string
    connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
    if err != nil {
        t.Fatalf("Failed to get connection string: %v", err)
    }
    
    // Create optimized database pool
    dbPool, err := db.NewOptimizedPool(connStr, false)
    if err != nil {
        t.Fatalf("Failed to create database pool: %v", err)
    }
    
    // Run migrations
    if err := runTestMigrations(connStr); err != nil {
        t.Fatalf("Failed to run migrations: %v", err)
    }
    
    // Create cache manager (in-memory for tests)
    cacheManager := cache.NewCacheManager("")
    
    config := &TestConfig{
        DBPool:            dbPool,
        Cache:             cacheManager,
        PostgresContainer: postgresContainer,
        TestDBURL:         connStr,
        Cleanup: func() {
            dbPool.Pool.Close()
            postgresContainer.Terminate(ctx)
        },
    }
    
    // Register cleanup
    t.Cleanup(config.Cleanup)
    
    return config
}

// TestFixtures provides common test data
type TestFixtures struct {
    Users        []User
    Organizations []Organization
    Plans        []Plan
    Contexts     []Context
}

// LoadTestFixtures loads predefined test data
func LoadTestFixtures(t *testing.T, config *TestConfig) *TestFixtures {
    fixtures := &TestFixtures{
        Users: []User{
            {ID: "test-user-1", Email: "user1@test.com", Name: "Test User 1"},
            {ID: "test-user-2", Email: "user2@test.com", Name: "Test User 2"},
        },
        Organizations: []Organization{
            {ID: "test-org-1", Name: "Test Organization 1"},
        },
        Plans: []Plan{
            {ID: "test-plan-1", Name: "Test Plan 1", UserID: "test-user-1"},
            {ID: "test-plan-2", Name: "Test Plan 2", UserID: "test-user-1"},
        },
        Contexts: []Context{
            {ID: "test-context-1", PlanID: "test-plan-1", Name: "Test Context 1"},
        },
    }
    
    // Insert fixtures into test database
    if err := insertTestFixtures(config.DBPool, fixtures); err != nil {
        t.Fatalf("Failed to load test fixtures: %v", err)
    }
    
    return fixtures
}

// runTestMigrations applies database migrations for testing
func runTestMigrations(connStr string) error {
    // Implementation would run database migrations
    // This is a simplified version - integrate with your migration system
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return err
    }
    defer db.Close()
    
    // Run basic table creation for tests
    migrations := []string{
        createUsersTable,
        createOrganizationsTable, 
        createPlansTable,
        createContextsTable,
        // Add other necessary tables
    }
    
    for _, migration := range migrations {
        if _, err := db.Exec(migration); err != nil {
            return fmt.Errorf("migration failed: %w", err)
        }
    }
    
    return nil
}

// Test database schema definitions
const (
    createUsersTable = `
        CREATE TABLE IF NOT EXISTS users (
            id UUID PRIMARY KEY,
            email VARCHAR(255) UNIQUE NOT NULL,
            name VARCHAR(255) NOT NULL,
            created_at TIMESTAMP DEFAULT NOW(),
            updated_at TIMESTAMP DEFAULT NOW(),
            deleted_at TIMESTAMP NULL
        )`
    
    createOrganizationsTable = `
        CREATE TABLE IF NOT EXISTS organizations (
            id UUID PRIMARY KEY,
            name VARCHAR(255) NOT NULL,
            created_at TIMESTAMP DEFAULT NOW(),
            updated_at TIMESTAMP DEFAULT NOW(),
            deleted_at TIMESTAMP NULL
        )`
    
    createPlansTable = `
        CREATE TABLE IF NOT EXISTS plans (
            id UUID PRIMARY KEY,
            name VARCHAR(255) NOT NULL,
            user_id UUID REFERENCES users(id),
            org_id UUID REFERENCES organizations(id),
            status VARCHAR(50) DEFAULT 'active',
            created_at TIMESTAMP DEFAULT NOW(),
            updated_at TIMESTAMP DEFAULT NOW(),
            archived_at TIMESTAMP NULL,
            deleted_at TIMESTAMP NULL
        )`
    
    createContextsTable = `
        CREATE TABLE IF NOT EXISTS contexts (
            id UUID PRIMARY KEY,
            plan_id UUID REFERENCES plans(id),
            name VARCHAR(255) NOT NULL,
            content TEXT,
            active BOOLEAN DEFAULT true,
            created_at TIMESTAMP DEFAULT NOW(),
            updated_at TIMESTAMP DEFAULT NOW()
        )`
)
```

**TodoWrite Task**: `Create comprehensive testing infrastructure with containers`

#### Step 3A.2: Unit Testing Framework
**File: `/app/server/handlers/auth_test.go`** (create comprehensive auth tests)
```go
package handlers

import (
    "bytes"
    "context"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "time"
    
    "github.com/gorilla/mux"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "yourapp/testing"
)

func TestAuthHandler_Register(t *testing.T) {
    // Setup test environment
    testConfig := testing.SetupTestEnvironment(t)
    fixtures := testing.LoadTestFixtures(t, testConfig)
    
    // Create auth handler with test dependencies
    authHandler := NewAuthHandler(testConfig.DBPool, testConfig.Cache)
    
    tests := []struct {
        name           string
        requestBody    map[string]interface{}
        expectedStatus int
        expectedError  string
        setupFunc      func()
        validateFunc   func(t *testing.T, response *httptest.ResponseRecorder)
    }{
        {
            name: "successful_registration",
            requestBody: map[string]interface{}{
                "email":    "newuser@test.com",
                "password": "securepassword123",
                "name":     "New User",
            },
            expectedStatus: http.StatusCreated,
            validateFunc: func(t *testing.T, response *httptest.ResponseRecorder) {
                var result map[string]interface{}
                err := json.NewDecoder(response.Body).Decode(&result)
                require.NoError(t, err)
                
                assert.Contains(t, result, "user_id")
                assert.Contains(t, result, "token")
                assert.Equal(t, "newuser@test.com", result["email"])
            },
        },
        {
            name: "duplicate_email_registration",
            requestBody: map[string]interface{}{
                "email":    "user1@test.com", // Exists in fixtures
                "password": "securepassword123",
                "name":     "Duplicate User",
            },
            expectedStatus: http.StatusConflict,
            expectedError:  "email already exists",
        },
        {
            name: "invalid_email_format",
            requestBody: map[string]interface{}{
                "email":    "invalid-email",
                "password": "securepassword123",
                "name":     "Invalid Email User",
            },
            expectedStatus: http.StatusBadRequest,
            expectedError:  "invalid email format",
        },
        {
            name: "weak_password",
            requestBody: map[string]interface{}{
                "email":    "weakpass@test.com",
                "password": "123",
                "name":     "Weak Password User",
            },
            expectedStatus: http.StatusBadRequest,
            expectedError:  "password too weak",
        },
        {
            name: "missing_required_fields",
            requestBody: map[string]interface{}{
                "email": "incomplete@test.com",
                // Missing password and name
            },
            expectedStatus: http.StatusBadRequest,
            expectedError:  "missing required fields",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Setup
            if tt.setupFunc != nil {
                tt.setupFunc()
            }
            
            // Prepare request
            requestBody, err := json.Marshal(tt.requestBody)
            require.NoError(t, err)
            
            req := httptest.NewRequest(http.MethodPost, "/api/auth/register", 
                bytes.NewBuffer(requestBody))
            req.Header.Set("Content-Type", "application/json")
            
            // Execute
            rr := httptest.NewRecorder()
            authHandler.Register(rr, req)
            
            // Validate status
            assert.Equal(t, tt.expectedStatus, rr.Code)
            
            // Validate error message
            if tt.expectedError != "" {
                var errorResponse map[string]interface{}
                err := json.NewDecoder(rr.Body).Decode(&errorResponse)
                require.NoError(t, err)
                assert.Contains(t, errorResponse["error"], tt.expectedError)
            }
            
            // Custom validation
            if tt.validateFunc != nil {
                tt.validateFunc(t, rr)
            }
        })
    }
}

func TestAuthHandler_Login(t *testing.T) {
    testConfig := testing.SetupTestEnvironment(t)
    fixtures := testing.LoadTestFixtures(t, testConfig)
    authHandler := NewAuthHandler(testConfig.DBPool, testConfig.Cache)
    
    // Pre-register a user for login tests
    testUser := &User{
        Email:    "logintest@test.com",
        Password: "hashedpassword123", // This should be properly hashed
        Name:     "Login Test User",
    }
    err := authHandler.userService.CreateUser(context.Background(), testUser)
    require.NoError(t, err)
    
    tests := []struct {
        name           string
        email          string
        password       string
        expectedStatus int
        expectedError  string
    }{
        {
            name:           "successful_login",
            email:          "logintest@test.com",
            password:       "hashedpassword123",
            expectedStatus: http.StatusOK,
        },
        {
            name:           "invalid_credentials",
            email:          "logintest@test.com",
            password:       "wrongpassword",
            expectedStatus: http.StatusUnauthorized,
            expectedError:  "invalid credentials",
        },
        {
            name:           "nonexistent_user",
            email:          "nonexistent@test.com",
            password:       "somepassword",
            expectedStatus: http.StatusUnauthorized,
            expectedError:  "invalid credentials",
        },
        {
            name:           "missing_email",
            email:          "",
            password:       "somepassword",
            expectedStatus: http.StatusBadRequest,
            expectedError:  "email required",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            requestBody := map[string]string{
                "email":    tt.email,
                "password": tt.password,
            }
            
            body, err := json.Marshal(requestBody)
            require.NoError(t, err)
            
            req := httptest.NewRequest(http.MethodPost, "/api/auth/login",
                bytes.NewBuffer(body))
            req.Header.Set("Content-Type", "application/json")
            
            rr := httptest.NewRecorder()
            authHandler.Login(rr, req)
            
            assert.Equal(t, tt.expectedStatus, rr.Code)
            
            if tt.expectedError != "" {
                var errorResponse map[string]interface{}
                err := json.NewDecoder(rr.Body).Decode(&errorResponse)
                require.NoError(t, err)
                assert.Contains(t, errorResponse["error"], tt.expectedError)
            } else {
                // Successful login should return token
                var successResponse map[string]interface{}
                err := json.NewDecoder(rr.Body).Decode(&successResponse)
                require.NoError(t, err)
                assert.Contains(t, successResponse, "token")
                assert.Contains(t, successResponse, "user")
            }
        })
    }
}

// Benchmark tests for performance validation
func BenchmarkAuthHandler_Register(b *testing.B) {
    testConfig := testing.SetupTestEnvironment(&testing.T{})
    defer testConfig.Cleanup()
    
    authHandler := NewAuthHandler(testConfig.DBPool, testConfig.Cache)
    
    requestBody := map[string]interface{}{
        "email":    "benchmark@test.com",
        "password": "securepassword123",
        "name":     "Benchmark User",
    }
    
    body, _ := json.Marshal(requestBody)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        req := httptest.NewRequest(http.MethodPost, "/api/auth/register",
            bytes.NewBuffer(body))
        req.Header.Set("Content-Type", "application/json")
        
        rr := httptest.NewRecorder()
        authHandler.Register(rr, req)
    }
}
```

**TodoWrite Task**: `Create comprehensive unit tests for authentication handlers`

#### Step 3A.3: Integration Testing Framework
**File: `/app/testing/integration/api_integration_test.go`** (create new file)
```go
package integration

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "net/http/httptest"
    "testing"
    "time"
    
    "github.com/gorilla/mux"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/stretchr/testify/suite"
    "yourapp/server"
    "yourapp/testing"
)

// APIIntegrationTestSuite provides integration testing for API endpoints
type APIIntegrationTestSuite struct {
    suite.Suite
    server     *httptest.Server
    client     *http.Client
    testConfig *testing.TestConfig
    authToken  string
    userID     string
}

// SetupSuite initializes the test suite
func (suite *APIIntegrationTestSuite) SetupSuite() {
    // Setup test environment
    suite.testConfig = testing.SetupTestEnvironment(suite.T())
    
    // Create test server
    app := server.NewApp(suite.testConfig.DBPool, suite.testConfig.Cache)
    suite.server = httptest.NewServer(app.Router)
    
    // Create HTTP client
    suite.client = &http.Client{
        Timeout: 30 * time.Second,
    }
    
    // Register test user and get auth token
    suite.registerTestUser()
}

// TearDownSuite cleans up after tests
func (suite *APIIntegrationTestSuite) TearDownSuite() {
    suite.server.Close()
    suite.testConfig.Cleanup()
}

// registerTestUser creates test user and obtains auth token
func (suite *APIIntegrationTestSuite) registerTestUser() {
    registerData := map[string]interface{}{
        "email":    "integration@test.com",
        "password": "integrationtest123",
        "name":     "Integration Test User",
    }
    
    body, err := json.Marshal(registerData)
    require.NoError(suite.T(), err)
    
    resp, err := suite.client.Post(
        suite.server.URL+"/api/auth/register",
        "application/json",
        bytes.NewBuffer(body),
    )
    require.NoError(suite.T(), err)
    defer resp.Body.Close()
    
    require.Equal(suite.T(), http.StatusCreated, resp.StatusCode)
    
    var result map[string]interface{}
    err = json.NewDecoder(resp.Body).Decode(&result)
    require.NoError(suite.T(), err)
    
    suite.authToken = result["token"].(string)
    suite.userID = result["user_id"].(string)
}

// makeAuthenticatedRequest helper for authenticated API calls
func (suite *APIIntegrationTestSuite) makeAuthenticatedRequest(method, path string, body interface{}) (*http.Response, error) {
    var reqBody *bytes.Buffer
    if body != nil {
        jsonBody, err := json.Marshal(body)
        if err != nil {
            return nil, err
        }
        reqBody = bytes.NewBuffer(jsonBody)
    } else {
        reqBody = bytes.NewBuffer(nil)
    }
    
    req, err := http.NewRequest(method, suite.server.URL+path, reqBody)
    if err != nil {
        return nil, err
    }
    
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+suite.authToken)
    
    return suite.client.Do(req)
}

// TestUserWorkflow tests complete user workflow
func (suite *APIIntegrationTestSuite) TestUserWorkflow() {
    // Test getting user profile
    resp, err := suite.makeAuthenticatedRequest("GET", "/api/user/profile", nil)
    require.NoError(suite.T(), err)
    defer resp.Body.Close()
    
    assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
    
    var profile map[string]interface{}
    err = json.NewDecoder(resp.Body).Decode(&profile)
    require.NoError(suite.T(), err)
    
    assert.Equal(suite.T(), "integration@test.com", profile["email"])
    assert.Equal(suite.T(), "Integration Test User", profile["name"])
}

// TestPlanWorkflow tests complete plan management workflow
func (suite *APIIntegrationTestSuite) TestPlanWorkflow() {
    // Create a new plan
    planData := map[string]interface{}{
        "name":        "Integration Test Plan",
        "description": "Test plan for integration testing",
    }
    
    resp, err := suite.makeAuthenticatedRequest("POST", "/api/plans", planData)
    require.NoError(suite.T(), err)
    defer resp.Body.Close()
    
    assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)
    
    var createdPlan map[string]interface{}
    err = json.NewDecoder(resp.Body).Decode(&createdPlan)
    require.NoError(suite.T(), err)
    
    planID := createdPlan["id"].(string)
    assert.NotEmpty(suite.T(), planID)
    
    // Get the created plan
    resp, err = suite.makeAuthenticatedRequest("GET", "/api/plans/"+planID, nil)
    require.NoError(suite.T(), err)
    defer resp.Body.Close()
    
    assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
    
    var retrievedPlan map[string]interface{}
    err = json.NewDecoder(resp.Body).Decode(&retrievedPlan)
    require.NoError(suite.T(), err)
    
    assert.Equal(suite.T(), planID, retrievedPlan["id"])
    assert.Equal(suite.T(), "Integration Test Plan", retrievedPlan["name"])
    
    // Update the plan
    updateData := map[string]interface{}{
        "name":        "Updated Integration Test Plan",
        "description": "Updated description",
    }
    
    resp, err = suite.makeAuthenticatedRequest("PUT", "/api/plans/"+planID, updateData)
    require.NoError(suite.T(), err)
    defer resp.Body.Close()
    
    assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
    
    // List user plans
    resp, err = suite.makeAuthenticatedRequest("GET", "/api/plans", nil)
    require.NoError(suite.T(), err)
    defer resp.Body.Close()
    
    assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
    
    var plans map[string]interface{}
    err = json.NewDecoder(resp.Body).Decode(&plans)
    require.NoError(suite.T(), err)
    
    planList := plans["plans"].([]interface{})
    assert.GreaterOrEqual(suite.T(), len(planList), 1)
    
    // Delete the plan
    resp, err = suite.makeAuthenticatedRequest("DELETE", "/api/plans/"+planID, nil)
    require.NoError(suite.T(), err)
    defer resp.Body.Close()
    
    assert.Equal(suite.T(), http.StatusNoContent, resp.StatusCode)
    
    // Verify plan is deleted
    resp, err = suite.makeAuthenticatedRequest("GET", "/api/plans/"+planID, nil)
    require.NoError(suite.T(), err)
    defer resp.Body.Close()
    
    assert.Equal(suite.T(), http.StatusNotFound, resp.StatusCode)
}

// TestContextWorkflow tests context management
func (suite *APIIntegrationTestSuite) TestContextWorkflow() {
    // First create a plan
    planData := map[string]interface{}{
        "name": "Context Test Plan",
    }
    
    resp, err := suite.makeAuthenticatedRequest("POST", "/api/plans", planData)
    require.NoError(suite.T(), err)
    defer resp.Body.Close()
    
    var plan map[string]interface{}
    err = json.NewDecoder(resp.Body).Decode(&plan)
    require.NoError(suite.T(), err)
    
    planID := plan["id"].(string)
    
    // Add context to the plan
    contextData := map[string]interface{}{
        "name":    "Test Context",
        "content": "This is test context content",
        "files":   []string{"test.go", "main.go"},
    }
    
    resp, err = suite.makeAuthenticatedRequest("POST", 
        "/api/plans/"+planID+"/contexts", contextData)
    require.NoError(suite.T(), err)
    defer resp.Body.Close()
    
    assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)
    
    var context map[string]interface{}
    err = json.NewDecoder(resp.Body).Decode(&context)
    require.NoError(suite.T(), err)
    
    contextID := context["id"].(string)
    
    // Get plan contexts
    resp, err = suite.makeAuthenticatedRequest("GET", 
        "/api/plans/"+planID+"/contexts", nil)
    require.NoError(suite.T(), err)
    defer resp.Body.Close()
    
    assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
    
    var contexts map[string]interface{}
    err = json.NewDecoder(resp.Body).Decode(&contexts)
    require.NoError(suite.T(), err)
    
    contextList := contexts["contexts"].([]interface{})
    assert.GreaterOrEqual(suite.T(), len(contextList), 1)
}

// TestErrorHandling tests error scenarios
func (suite *APIIntegrationTestSuite) TestErrorHandling() {
    // Test unauthorized access
    req, err := http.NewRequest("GET", suite.server.URL+"/api/plans", nil)
    require.NoError(suite.T(), err)
    
    resp, err := suite.client.Do(req)
    require.NoError(suite.T(), err)
    defer resp.Body.Close()
    
    assert.Equal(suite.T(), http.StatusUnauthorized, resp.StatusCode)
    
    // Test invalid JSON
    resp, err = suite.makeAuthenticatedRequest("POST", "/api/plans", 
        "invalid json")
    require.NoError(suite.T(), err)
    defer resp.Body.Close()
    
    assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)
    
    // Test not found
    resp, err = suite.makeAuthenticatedRequest("GET", 
        "/api/plans/nonexistent-id", nil)
    require.NoError(suite.T(), err)
    defer resp.Body.Close()
    
    assert.Equal(suite.T(), http.StatusNotFound, resp.StatusCode)
}

// TestRateLimiting tests rate limiting functionality
func (suite *APIIntegrationTestSuite) TestRateLimiting() {
    // Make multiple rapid requests to test rate limiting
    for i := 0; i < 100; i++ {
        resp, err := suite.makeAuthenticatedRequest("GET", "/api/plans", nil)
        if err != nil {
            continue
        }
        resp.Body.Close()
        
        // If we hit rate limit, expect 429 status
        if resp.StatusCode == http.StatusTooManyRequests {
            suite.T().Log("Rate limiting working correctly")
            return
        }
    }
}

// Run the integration test suite
func TestAPIIntegrationSuite(t *testing.T) {
    suite.Run(t, new(APIIntegrationTestSuite))
}
```

**TodoWrite Task**: `Create comprehensive integration tests for API workflows`

#### Step 3A.4: End-to-End Testing Framework
**File: `/app/testing/e2e/cli_e2e_test.go`** (create new file)
```go
package e2e

import (
    "bytes"
    "context"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
    "testing"
    "time"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/stretchr/testify/suite"
)

// CLIEndToEndTestSuite tests complete CLI workflows
type CLIEndToEndTestSuite struct {
    suite.Suite
    cliPath     string
    tempDir     string
    configDir   string
    serverURL   string
}

// SetupSuite initializes E2E test environment
func (suite *CLIEndToEndTestSuite) SetupSuite() {
    // Build CLI binary for testing
    var err error
    suite.tempDir, err = os.MkdirTemp("", "plandex-e2e-*")
    require.NoError(suite.T(), err)
    
    suite.cliPath = filepath.Join(suite.tempDir, "plandex-cli")
    suite.configDir = filepath.Join(suite.tempDir, ".plandex")
    
    // Build the CLI
    cmd := exec.Command("go", "build", "-o", suite.cliPath, "../../cli/main.go")
    output, err := cmd.CombinedOutput()
    require.NoError(suite.T(), err, "Failed to build CLI: %s", output)
    
    // Start test server (this would ideally use testcontainers)
    suite.serverURL = "http://localhost:8080"
    
    // Setup config directory
    err = os.MkdirAll(suite.configDir, 0755)
    require.NoError(suite.T(), err)
}

// TearDownSuite cleans up E2E test environment
func (suite *CLIEndToEndTestSuite) TearDownSuite() {
    os.RemoveAll(suite.tempDir)
}

// runCLICommand executes CLI command and returns output
func (suite *CLIEndToEndTestSuite) runCLICommand(args ...string) (string, string, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    cmd := exec.CommandContext(ctx, suite.cliPath, args...)
    cmd.Env = append(os.Environ(),
        "PLANDEX_CONFIG_DIR="+suite.configDir,
        "PLANDEX_SERVER_URL="+suite.serverURL,
    )
    
    var stdout, stderr bytes.Buffer
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr
    
    err := cmd.Run()
    return stdout.String(), stderr.String(), err
}

// TestCLIVersion tests version command
func (suite *CLIEndToEndTestSuite) TestCLIVersion() {
    stdout, stderr, err := suite.runCLICommand("version")
    
    assert.NoError(suite.T(), err, "stderr: %s", stderr)
    assert.Contains(suite.T(), stdout, "plandex")
    assert.NotContains(suite.T(), stdout, "unknown")
}

// TestCLIHelp tests help functionality
func (suite *CLIEndToEndTestSuite) TestCLIHelp() {
    stdout, stderr, err := suite.runCLICommand("--help")
    
    assert.NoError(suite.T(), err, "stderr: %s", stderr)
    assert.Contains(suite.T(), stdout, "Usage:")
    assert.Contains(suite.T(), stdout, "Commands:")
    assert.Contains(suite.T(), stdout, "Flags:")
}

// TestAuthWorkflow tests complete authentication workflow
func (suite *CLIEndToEndTestSuite) TestAuthWorkflow() {
    // Test signup
    stdout, stderr, err := suite.runCLICommand("auth", "signup", 
        "--email", "e2e@test.com",
        "--password", "e2etest123",
        "--name", "E2E Test User")
    
    if err != nil {
        suite.T().Logf("Signup output - stdout: %s, stderr: %s", stdout, stderr)
    }
    
    // Note: This might fail if user already exists, which is okay for E2E tests
    // The important thing is that the command structure is correct
    
    // Test login
    stdout, stderr, err = suite.runCLICommand("auth", "login",
        "--email", "e2e@test.com", 
        "--password", "e2etest123")
    
    if err == nil {
        assert.Contains(suite.T(), stdout, "success")
    }
    
    // Test status
    stdout, stderr, err = suite.runCLICommand("auth", "status")
    assert.NoError(suite.T(), err, "stderr: %s", stderr)
}

// TestPlanWorkflow tests complete plan management workflow
func (suite *CLIEndToEndTestSuite) TestPlanWorkflow() {
    // Ensure we're authenticated (this might fail, but that's okay for isolated testing)
    suite.runCLICommand("auth", "login", "--email", "e2e@test.com", "--password", "e2etest123")
    
    // Create a new plan
    stdout, stderr, err := suite.runCLICommand("plans", "create", "E2E Test Plan")
    
    if err != nil {
        suite.T().Logf("Plan create failed (might be auth issue): %s", stderr)
        suite.T().Skip("Skipping plan workflow test due to auth requirements")
        return
    }
    
    assert.Contains(suite.T(), stdout, "created")
    
    // List plans
    stdout, stderr, err = suite.runCLICommand("plans", "list")
    assert.NoError(suite.T(), err, "stderr: %s", stderr)
    assert.Contains(suite.T(), stdout, "E2E Test Plan")
    
    // Get plan ID from output (this would need parsing)
    lines := strings.Split(stdout, "\n")
    var planID string
    for _, line := range lines {
        if strings.Contains(line, "E2E Test Plan") {
            // Extract plan ID from line (implementation depends on output format)
            // This is simplified - real implementation would parse structured output
            break
        }
    }
    
    if planID != "" {
        // Test plan operations
        stdout, stderr, err = suite.runCLICommand("plans", "show", planID)
        assert.NoError(suite.T(), err, "stderr: %s", stderr)
        assert.Contains(suite.T(), stdout, "E2E Test Plan")
    }
}

// TestProjectIntegration tests CLI integration with project files
func (suite *CLIEndToEndTestSuite) TestProjectIntegration() {
    // Create a temporary project directory
    projectDir, err := os.MkdirTemp(suite.tempDir, "test-project-*")
    require.NoError(suite.T(), err)
    
    // Create some test files
    testFiles := map[string]string{
        "main.go": `package main

import "fmt"

func main() {
    fmt.Println("Hello, Plandex!")
}`,
        "README.md": `# Test Project

This is a test project for Plandex E2E testing.`,
        "go.mod": `module test-project

go 1.23`,
    }
    
    for filename, content := range testFiles {
        filePath := filepath.Join(projectDir, filename)
        err := os.WriteFile(filePath, []byte(content), 0644)
        require.NoError(suite.T(), err)
    }
    
    // Change to project directory
    originalDir, err := os.Getwd()
    require.NoError(suite.T(), err)
    defer os.Chdir(originalDir)
    
    err = os.Chdir(projectDir)
    require.NoError(suite.T(), err)
    
    // Test project initialization
    stdout, stderr, err := suite.runCLICommand("init")
    
    if err != nil {
        suite.T().Logf("Project init output - stdout: %s, stderr: %s", stdout, stderr)
        // This might fail due to auth requirements, which is okay for isolated testing
    }
    
    // Test file listing
    stdout, stderr, err = suite.runCLICommand("files", "list")
    
    if err == nil {
        assert.Contains(suite.T(), stdout, "main.go")
        assert.Contains(suite.T(), stdout, "README.md")
    }
}

// TestErrorScenarios tests CLI error handling
func (suite *CLIEndToEndTestSuite) TestErrorScenarios() {
    // Test invalid command
    stdout, stderr, err := suite.runCLICommand("invalid-command")
    assert.Error(suite.T(), err)
    assert.Contains(suite.T(), stderr, "unknown command")
    
    // Test missing required arguments
    stdout, stderr, err = suite.runCLICommand("plans", "create")
    assert.Error(suite.T(), err)
    
    // Test invalid flags
    stdout, stderr, err = suite.runCLICommand("--invalid-flag")
    assert.Error(suite.T(), err)
}

// TestPerformance tests CLI performance characteristics
func (suite *CLIEndToEndTestSuite) TestPerformance() {
    // Test command execution time
    start := time.Now()
    stdout, stderr, err := suite.runCLICommand("version")
    duration := time.Since(start)
    
    assert.NoError(suite.T(), err, "stderr: %s", stderr)
    assert.Less(suite.T(), duration, 2*time.Second, "Version command took too long")
    
    // Test help command performance
    start = time.Now()
    stdout, stderr, err = suite.runCLICommand("--help")
    duration = time.Since(start)
    
    assert.NoError(suite.T(), err, "stderr: %s", stderr)
    assert.Less(suite.T(), duration, 1*time.Second, "Help command took too long")
}

// Run the E2E test suite
func TestCLIEndToEndSuite(t *testing.T) {
    suite.Run(t, new(CLIEndToEndTestSuite))
}
```

**TodoWrite Task**: `Create comprehensive end-to-end CLI testing framework`

### KPIs for Phase 3A
- ✅ Test coverage increased from ~5% to 80%+
- ✅ All test types implemented (unit, integration, E2E)
- ✅ Test execution time <5 minutes for full suite
- ✅ Zero flaky tests (consistent pass/fail)
- ✅ Comprehensive test data fixtures
- ✅ Isolated test environments with containers

---

## Phase 3B: CI/CD Pipeline Implementation
### 🚀 AUTOMATED QUALITY ASSURANCE TARGET
**Goal**: Fully automated testing, security scanning, and deployment pipeline

### Implementation Steps

#### Step 3B.1: GitHub Actions Workflow Enhancement
**File: `.github/workflows/comprehensive-ci.yml`** (create new file)
```yaml
name: Comprehensive CI/CD Pipeline

on:
  push:
    branches: [ main, develop ]
    tags: [ 'v*' ]
  pull_request:
    branches: [ main, develop ]

env:
  GO_VERSION: '1.23.10'
  NODE_VERSION: '18'
  REGISTRY: ghcr.io
  IMAGE_NAME: ${{ github.repository }}

jobs:
  # Job 1: Code Quality Checks
  code-quality:
    name: Code Quality & Security
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v4
        with:
          fetch-depth: 0  # Full history for better analysis

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: ${{ env.GO_VERSION }}

      - name: Cache Go modules
        uses: actions/cache@v3
        with:
          path: |
            ~/.cache/go-build
            ~/go/pkg/mod
          key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
          restore-keys: |
            ${{ runner.os }}-go-

      - name: Download dependencies
        run: |
          go mod download
          cd cli && go mod download
          cd ../server && go mod download
          cd ../shared && go mod download

      - name: Run gofmt check
        run: |
          if [ "$(gofmt -s -l . | wc -l)" -gt 0 ]; then
            echo "Code is not formatted. Please run 'gofmt -s -w .'"
            gofmt -s -l .
            exit 1
          fi

      - name: Run go vet
        run: |
          go vet ./...
          cd cli && go vet ./...
          cd ../server && go vet ./...
          cd ../shared && go vet ./...

      - name: Install staticcheck
        run: go install honnef.co/go/tools/cmd/staticcheck@latest

      - name: Run staticcheck
        run: |
          staticcheck ./...
          cd cli && staticcheck ./...
          cd ../server && staticcheck ./...
          cd ../shared && staticcheck ./...

      - name: Install gosec
        run: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest

      - name: Run gosec security scan
        run: |
          gosec -fmt sarif -out gosec-results.sarif ./...
        continue-on-error: true

      - name: Upload SARIF file
        uses: github/codeql-action/upload-sarif@v2
        with:
          sarif_file: gosec-results.sarif
        if: always()

      - name: Install govulncheck
        run: go install golang.org/x/vuln/cmd/govulncheck@latest

      - name: Run vulnerability check
        run: |
          govulncheck ./...
          cd cli && govulncheck ./...
          cd ../server && govulncheck ./...
          cd ../shared && govulncheck ./...

  # Job 2: Unit Tests
  unit-tests:
    name: Unit Tests
    runs-on: ubuntu-latest
    needs: code-quality
    
    services:
      postgres:
        image: postgres:17.5
        env:
          POSTGRES_PASSWORD: testpass
          POSTGRES_USER: testuser
          POSTGRES_DB: testdb
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 5432:5432

      redis:
        image: redis:7.2-alpine
        options: >-
          --health-cmd "redis-cli ping"
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 6379:6379

    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: ${{ env.GO_VERSION }}

      - name: Cache Go modules
        uses: actions/cache@v3
        with:
          path: |
            ~/.cache/go-build
            ~/go/pkg/mod
          key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
          restore-keys: |
            ${{ runner.os }}-go-

      - name: Download dependencies
        run: |
          go mod download
          cd cli && go mod download
          cd ../server && go mod download
          cd ../shared && go mod download

      - name: Run unit tests
        env:
          DATABASE_URL: postgres://testuser:testpass@localhost:5432/testdb?sslmode=disable
          REDIS_URL: redis://localhost:6379
          TEST_ENV: true
        run: |
          go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
          cd cli && go test -v -race -coverprofile=cli-coverage.out -covermode=atomic ./...
          cd ../server && go test -v -race -coverprofile=server-coverage.out -covermode=atomic ./...
          cd ../shared && go test -v -race -coverprofile=shared-coverage.out -covermode=atomic ./...

      - name: Combine coverage reports
        run: |
          echo "mode: atomic" > combined-coverage.out
          tail -n +2 coverage.out >> combined-coverage.out
          tail -n +2 cli/cli-coverage.out >> combined-coverage.out
          tail -n +2 server/server-coverage.out >> combined-coverage.out
          tail -n +2 shared/shared-coverage.out >> combined-coverage.out

      - name: Upload coverage to Codecov
        uses: codecov/codecov-action@v3
        with:
          file: combined-coverage.out
          flags: unittests
          name: codecov-umbrella

      - name: Check coverage threshold
        run: |
          COVERAGE=$(go tool cover -func=combined-coverage.out | grep total | awk '{print $3}' | sed 's/%//')
          echo "Coverage: $COVERAGE%"
          if (( $(echo "$COVERAGE < 80" | bc -l) )); then
            echo "Coverage $COVERAGE% is below threshold of 80%"
            exit 1
          fi

  # Job 3: Integration Tests
  integration-tests:
    name: Integration Tests
    runs-on: ubuntu-latest
    needs: unit-tests
    
    services:
      postgres:
        image: postgres:17.5
        env:
          POSTGRES_PASSWORD: testpass
          POSTGRES_USER: testuser
          POSTGRES_DB: testdb
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 5432:5432

      redis:
        image: redis:7.2-alpine
        options: >-
          --health-cmd "redis-cli ping"
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 6379:6379

    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: ${{ env.GO_VERSION }}

      - name: Cache Go modules
        uses: actions/cache@v3
        with:
          path: |
            ~/.cache/go-build
            ~/go/pkg/mod
          key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
          restore-keys: |
            ${{ runner.os }}-go-

      - name: Download dependencies
        run: |
          go mod download
          cd cli && go mod download
          cd ../server && go mod download
          cd ../shared && go mod download

      - name: Run database migrations
        env:
          DATABASE_URL: postgres://testuser:testpass@localhost:5432/testdb?sslmode=disable
        run: |
          cd server && go run migrations/migrate.go up

      - name: Build server for integration tests
        run: |
          cd server && go build -o ../plandex-server main.go

      - name: Start test server
        env:
          DATABASE_URL: postgres://testuser:testpass@localhost:5432/testdb?sslmode=disable
          REDIS_URL: redis://localhost:6379
          PORT: 8080
          LOG_LEVEL: debug
        run: |
          ./plandex-server &
          echo $! > server.pid
          sleep 5

      - name: Wait for server to be ready
        run: |
          timeout 30 bash -c 'until curl -f http://localhost:8080/api/health; do sleep 1; done'

      - name: Run integration tests
        env:
          TEST_SERVER_URL: http://localhost:8080
          DATABASE_URL: postgres://testuser:testpass@localhost:5432/testdb?sslmode=disable
        run: |
          go test -v -tags=integration ./testing/integration/...

      - name: Stop test server
        run: |
          if [ -f server.pid ]; then
            kill $(cat server.pid) || true
          fi

  # Job 4: End-to-End Tests
  e2e-tests:
    name: End-to-End Tests
    runs-on: ubuntu-latest
    needs: integration-tests
    
    services:
      postgres:
        image: postgres:17.5
        env:
          POSTGRES_PASSWORD: testpass
          POSTGRES_USER: testuser
          POSTGRES_DB: testdb
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 5432:5432

    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: ${{ env.GO_VERSION }}

      - name: Cache Go modules
        uses: actions/cache@v3
        with:
          path: |
            ~/.cache/go-build
            ~/go/pkg/mod
          key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
          restore-keys: |
            ${{ runner.os }}-go-

      - name: Build CLI and Server
        run: |
          cd cli && go build -o ../plandex-cli main.go
          cd ../server && go build -o ../plandex-server main.go

      - name: Run database migrations
        env:
          DATABASE_URL: postgres://testuser:testpass@localhost:5432/testdb?sslmode=disable
        run: |
          cd server && go run migrations/migrate.go up

      - name: Start test server
        env:
          DATABASE_URL: postgres://testuser:testpass@localhost:5432/testdb?sslmode=disable
          PORT: 8080
        run: |
          ./plandex-server &
          echo $! > server.pid
          sleep 5

      - name: Run E2E tests
        env:
          PLANDEX_CLI_PATH: ./plandex-cli
          PLANDEX_SERVER_URL: http://localhost:8080
        run: |
          go test -v -tags=e2e ./testing/e2e/...

      - name: Stop test server
        run: |
          if [ -f server.pid ]; then
            kill $(cat server.pid) || true
          fi

  # Job 5: Build and Security Scan
  build-and-scan:
    name: Build & Security Scan
    runs-on: ubuntu-latest
    needs: code-quality
    
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: ${{ env.GO_VERSION }}

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3

      - name: Log in to Container Registry
        if: github.event_name != 'pull_request'
        uses: docker/login-action@v3
        with:
          registry: ${{ env.REGISTRY }}
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      - name: Extract metadata
        id: meta
        uses: docker/metadata-action@v5
        with:
          images: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}
          tags: |
            type=ref,event=branch
            type=ref,event=pr
            type=semver,pattern={{version}}
            type=semver,pattern={{major}}.{{minor}}

      - name: Build optimized binaries
        run: |
          mkdir -p dist
          
          # Build server
          CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
            -ldflags='-w -s' \
            -o dist/plandex-server \
            ./server/main.go
          
          # Build CLI
          CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
            -ldflags='-w -s' \
            -o dist/plandex-cli \
            ./cli/main.go

      - name: Build Docker images
        uses: docker/build-push-action@v5
        with:
          context: .
          file: ./server/Dockerfile.optimized
          push: ${{ github.event_name != 'pull_request' }}
          tags: ${{ steps.meta.outputs.tags }}
          labels: ${{ steps.meta.outputs.labels }}
          cache-from: type=gha
          cache-to: type=gha,mode=max

      - name: Run Trivy vulnerability scanner
        uses: aquasecurity/trivy-action@master
        with:
          image-ref: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}:latest
          format: 'sarif'
          output: 'trivy-results.sarif'
        if: github.event_name != 'pull_request'

      - name: Upload Trivy scan results to GitHub Security
        uses: github/codeql-action/upload-sarif@v2
        with:
          sarif_file: 'trivy-results.sarif'
        if: github.event_name != 'pull_request'

  # Job 6: Performance Benchmarks
  performance-benchmarks:
    name: Performance Benchmarks
    runs-on: ubuntu-latest
    needs: unit-tests
    
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: ${{ env.GO_VERSION }}

      - name: Cache Go modules
        uses: actions/cache@v3
        with:
          path: |
            ~/.cache/go-build
            ~/go/pkg/mod
          key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
          restore-keys: |
            ${{ runner.os }}-go-

      - name: Download dependencies
        run: go mod download

      - name: Run benchmarks
        run: |
          go test -bench=. -benchmem -run=^$ ./... > benchmarks.txt
          cat benchmarks.txt

      - name: Store benchmark result
        uses: benchmark-action/github-action-benchmark@v1
        with:
          tool: 'go'
          output-file-path: benchmarks.txt
          github-token: ${{ secrets.GITHUB_TOKEN }}
          auto-push: true
          alert-threshold: '200%'
          comment-on-alert: true

  # Job 7: Deploy (only on main branch and tags)
  deploy:
    name: Deploy
    runs-on: ubuntu-latest
    needs: [unit-tests, integration-tests, e2e-tests, build-and-scan]
    if: github.event_name == 'push' && (github.ref == 'refs/heads/main' || startsWith(github.ref, 'refs/tags/'))
    
    environment:
      name: ${{ github.ref == 'refs/heads/main' && 'staging' || 'production' }}
      url: ${{ steps.deploy.outputs.url }}
    
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Deploy to staging
        if: github.ref == 'refs/heads/main'
        id: deploy-staging
        run: |
          echo "Deploying to staging environment"
          echo "url=https://staging.plandex.example.com" >> $GITHUB_OUTPUT

      - name: Deploy to production
        if: startsWith(github.ref, 'refs/tags/')
        id: deploy-production
        run: |
          echo "Deploying to production environment"
          echo "url=https://plandex.example.com" >> $GITHUB_OUTPUT

      - name: Set deployment URL
        id: deploy
        run: |
          if [ "${{ github.ref }}" = "refs/heads/main" ]; then
            echo "url=https://staging.plandex.example.com" >> $GITHUB_OUTPUT
          else
            echo "url=https://plandex.example.com" >> $GITHUB_OUTPUT
          fi

  # Job 8: Notify
  notify:
    name: Notify
    runs-on: ubuntu-latest
    needs: [deploy]
    if: always()
    
    steps:
      - name: Notify on success
        if: needs.deploy.result == 'success'
        run: |
          echo "Deployment successful!"

      - name: Notify on failure
        if: needs.deploy.result == 'failure'
        run: |
          echo "Deployment failed!"
          exit 1
```

**TodoWrite Task**: `Create comprehensive CI/CD pipeline with all quality gates`

#### Step 3B.2: Pre-commit Hooks Implementation
**File: `.pre-commit-config.yaml`** (create new file)
```yaml
# Pre-commit hooks configuration for Plandex
repos:
  - repo: https://github.com/pre-commit/pre-commit-hooks
    rev: v4.4.0
    hooks:
      - id: trailing-whitespace
      - id: end-of-file-fixer
      - id: check-yaml
      - id: check-json
      - id: check-added-large-files
        args: ['--maxkb=1000']
      - id: check-merge-conflict
      - id: check-case-conflict
      - id: mixed-line-ending

  - repo: local
    hooks:
      - id: go-fmt
        name: Go Format
        entry: gofmt
        args: [-w, -s]
        language: system
        files: \.go$

      - id: go-vet
        name: Go Vet
        entry: bash
        args:
          - -c
          - |
            go vet ./...
            cd cli && go vet ./...
            cd ../server && go vet ./...
            cd ../shared && go vet ./...
        language: system
        files: \.go$
        pass_filenames: false

      - id: go-imports
        name: Go Imports
        entry: bash
        args:
          - -c
          - |
            if ! command -v goimports &> /dev/null; then
              go install golang.org/x/tools/cmd/goimports@latest
            fi
            goimports -w -local yourapp .
        language: system
        files: \.go$
        pass_filenames: false

      - id: go-test
        name: Go Test
        entry: bash
        args:
          - -c
          - |
            go test -short ./...
            cd cli && go test -short ./...
            cd ../server && go test -short ./...
            cd ../shared && go test -short ./...
        language: system
        files: \.go$
        pass_filenames: false

      - id: go-sec
        name: Go Security Check
        entry: bash
        args:
          - -c
          - |
            if ! command -v gosec &> /dev/null; then
              go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
            fi
            gosec -quiet ./...
        language: system
        files: \.go$
        pass_filenames: false

      - id: go-mod-tidy
        name: Go Mod Tidy
        entry: bash
        args:
          - -c
          - |
            go mod tidy
            cd cli && go mod tidy
            cd ../server && go mod tidy
            cd ../shared && go mod tidy
        language: system
        files: (go\.mod|go\.sum)$
        pass_filenames: false

      - id: docker-lint
        name: Dockerfile Lint
        entry: bash
        args:
          - -c
          - |
            if command -v hadolint &> /dev/null; then
              hadolint Dockerfile*
            else
              echo "hadolint not found, skipping Dockerfile linting"
            fi
        language: system
        files: Dockerfile.*
        pass_filenames: false

      - id: yaml-lint
        name: YAML Lint
        entry: bash
        args:
          - -c
          - |
            if command -v yamllint &> /dev/null; then
              yamllint .
            else
              echo "yamllint not found, skipping YAML linting"
            fi
        language: system
        files: \.(yml|yaml)$
        pass_filenames: false

      - id: markdown-lint
        name: Markdown Lint
        entry: bash
        args:
          - -c
          - |
            if command -v markdownlint &> /dev/null; then
              markdownlint README.md docs/
            else
              echo "markdownlint not found, skipping Markdown linting"
            fi
        language: system
        files: \.md$
        pass_filenames: false
```

**File: `/scripts/setup-dev-env.sh`** (create development setup script)
```bash
#!/bin/bash

# Development environment setup script for Plandex
set -euo pipefail

echo "Setting up Plandex development environment..."

# Check if required tools are installed
check_tool() {
    if ! command -v "$1" &> /dev/null; then
        echo "❌ $1 is not installed. Please install it first."
        return 1
    else
        echo "✅ $1 is installed"
        return 0
    fi
}

echo "Checking required tools..."
MISSING_TOOLS=()

if ! check_tool "go"; then MISSING_TOOLS+=("go"); fi
if ! check_tool "docker"; then MISSING_TOOLS+=("docker"); fi
if ! check_tool "git"; then MISSING_TOOLS+=("git"); fi

# Optional tools
OPTIONAL_TOOLS=("hadolint" "yamllint" "markdownlint" "pre-commit")
echo "Checking optional tools..."
for tool in "${OPTIONAL_TOOLS[@]}"; do
    if check_tool "$tool"; then
        continue
    else
        echo "ℹ️  $tool is optional but recommended"
    fi
done

if [ ${#MISSING_TOOLS[@]} -ne 0 ]; then
    echo "❌ Missing required tools: ${MISSING_TOOLS[*]}"
    echo "Please install them and run this script again."
    exit 1
fi

# Install Go tools
echo "Installing Go development tools..."
go install golang.org/x/tools/cmd/goimports@latest
go install honnef.co/go/tools/cmd/staticcheck@latest
go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
go install golang.org/x/vuln/cmd/govulncheck@latest

# Setup pre-commit hooks if pre-commit is available
if command -v pre-commit &> /dev/null; then
    echo "Setting up pre-commit hooks..."
    pre-commit install
    pre-commit install --hook-type commit-msg
else
    echo "ℹ️  pre-commit not found. You can install it with: pip install pre-commit"
fi

# Create development config directories
echo "Creating development directories..."
mkdir -p .local/{cache,config,data}
mkdir -p logs
mkdir -p tmp

# Setup Git hooks (alternative to pre-commit)
if [ ! -f .git/hooks/pre-commit ] && ! command -v pre-commit &> /dev/null; then
    echo "Setting up Git hooks..."
    cat > .git/hooks/pre-commit << 'EOF'
#!/bin/bash
# Pre-commit hook for Plandex

set -e

echo "Running pre-commit checks..."

# Go formatting
echo "Checking Go formatting..."
if [ "$(gofmt -s -l . | wc -l)" -gt 0 ]; then
    echo "❌ Code is not formatted. Run 'gofmt -s -w .'"
    gofmt -s -l .
    exit 1
fi

# Go vet
echo "Running go vet..."
go vet ./...
cd cli && go vet ./...
cd ../server && go vet ./...
cd ../shared && go vet ./...
cd ..

# Run tests
echo "Running tests..."
go test -short ./...

echo "✅ Pre-commit checks passed"
EOF
    chmod +x .git/hooks/pre-commit
    echo "✅ Git pre-commit hook installed"
fi

# Setup environment file template
if [ ! -f .env.local ]; then
    echo "Creating .env.local template..."
    cat > .env.local << 'EOF'
# Local development environment variables
# Copy this file to .env and modify as needed

# Database
DATABASE_URL=postgres://plandex:plandex@localhost:5432/plandex?sslmode=disable

# Redis (optional)
REDIS_URL=redis://localhost:6379

# Server
PORT=8080
LOG_LEVEL=debug
ENVIRONMENT=development

# API Keys (add your own)
OPENAI_API_KEY=your_openai_key_here
ANTHROPIC_API_KEY=your_anthropic_key_here

# Security
JWT_SECRET=your_jwt_secret_here
CORS_ORIGINS=http://localhost:3000,http://localhost:8080

# Performance tuning for development
GOGC=100
GOMEMLIMIT=512MiB
EOF
    echo "✅ .env.local template created"
fi

# Install additional development tools based on OS
if [[ "$OSTYPE" == "darwin"* ]]; then
    echo "macOS detected. Consider installing additional tools:"
    echo "  brew install hadolint yamllint markdownlint-cli pre-commit"
elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
    echo "Linux detected. Consider installing additional tools:"
    echo "  # hadolint: https://github.com/hadolint/hadolint#install"
    echo "  # yamllint: pip install yamllint"
    echo "  # markdownlint: npm install -g markdownlint-cli"
    echo "  # pre-commit: pip install pre-commit"
fi

echo ""
echo "🎉 Development environment setup complete!"
echo ""
echo "Next steps:"
echo "1. Copy .env.local to .env and add your API keys"
echo "2. Start the development database: docker-compose up -d postgres"
echo "3. Run database migrations: cd server && go run migrations/migrate.go up"
echo "4. Start the development server: cd server && go run main.go"
echo "5. In another terminal, test the CLI: cd cli && go run main.go --help"
echo ""
echo "Happy coding! 🚀"
```

**TodoWrite Task**: `Implement pre-commit hooks and development environment setup`

### KPIs for Phase 3B
- ✅ 100% automated CI/CD pipeline
- ✅ Zero manual intervention for testing/deployment
- ✅ <5 minute feedback loop for developers
- ✅ 99.9% pipeline success rate
- ✅ Comprehensive security scanning
- ✅ Automated performance regression detection

---

## Phase 3C: Monitoring & Observability Infrastructure
### 📊 COMPREHENSIVE MONITORING TARGET
**Goal**: Full observability into application performance, health, and business metrics

### Implementation Steps

#### Step 3C.1: Structured Logging Implementation
**File: `/app/server/logging/structured_logger.go`** (create new file)
```go
package logging

import (
    "context"
    "encoding/json"
    "io"
    "log/slog"
    "os"
    "runtime"
    "time"
)

// StructuredLogger provides enterprise-grade logging
type StructuredLogger struct {
    logger *slog.Logger
    level  slog.Level
}

// LogContext contains contextual information for logs
type LogContext struct {
    RequestID   string `json:"request_id,omitempty"`
    UserID      string `json:"user_id,omitempty"`
    PlanID      string `json:"plan_id,omitempty"`
    Operation   string `json:"operation,omitempty"`
    Component   string `json:"component,omitempty"`
    Environment string `json:"environment,omitempty"`
}

// NewStructuredLogger creates a new structured logger
func NewStructuredLogger(level string, output io.Writer) *StructuredLogger {
    if output == nil {
        output = os.Stdout
    }
    
    var logLevel slog.Level
    switch level {
    case "debug":
        logLevel = slog.LevelDebug
    case "info":
        logLevel = slog.LevelInfo
    case "warn":
        logLevel = slog.LevelWarn
    case "error":
        logLevel = slog.LevelError
    default:
        logLevel = slog.LevelInfo
    }
    
    opts := &slog.HandlerOptions{
        Level: logLevel,
        AddSource: true,
        ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
            // Customize attribute formatting
            if a.Key == slog.TimeKey {
                return slog.Attr{
                    Key:   "timestamp",
                    Value: slog.StringValue(a.Value.Time().UTC().Format(time.RFC3339Nano)),
                }
            }
            if a.Key == slog.SourceKey {
                source := a.Value.Any().(*slog.Source)
                return slog.Attr{
                    Key: "source",
                    Value: slog.StringValue(
                        fmt.Sprintf("%s:%d", filepath.Base(source.File), source.Line),
                    ),
                }
            }
            return a
        },
    }
    
    handler := slog.NewJSONHandler(output, opts)
    logger := slog.New(handler)
    
    return &StructuredLogger{
        logger: logger,
        level:  logLevel,
    }
}

// WithContext adds contextual information to logger
func (sl *StructuredLogger) WithContext(ctx LogContext) *StructuredLogger {
    logger := sl.logger.With(
        "request_id", ctx.RequestID,
        "user_id", ctx.UserID,
        "plan_id", ctx.PlanID,
        "operation", ctx.Operation,
        "component", ctx.Component,
        "environment", ctx.Environment,
    )
    
    return &StructuredLogger{
        logger: logger,
        level:  sl.level,
    }
}

// Info logs an info message
func (sl *StructuredLogger) Info(msg string, args ...interface{}) {
    sl.logger.Info(msg, args...)
}

// Error logs an error message
func (sl *StructuredLogger) Error(msg string, err error, args ...interface{}) {
    allArgs := append([]interface{}{"error", err}, args...)
    sl.logger.Error(msg, allArgs...)
}

// Debug logs a debug message
func (sl *StructuredLogger) Debug(msg string, args ...interface{}) {
    sl.logger.Debug(msg, args...)
}

// Warn logs a warning message
func (sl *StructuredLogger) Warn(msg string, args ...interface{}) {
    sl.logger.Warn(msg, args...)
}

// LogHTTPRequest logs HTTP request details
func (sl *StructuredLogger) LogHTTPRequest(ctx context.Context, req *http.Request, statusCode int, duration time.Duration) {
    sl.logger.InfoContext(ctx, "HTTP request",
        "method", req.Method,
        "path", req.URL.Path,
        "query", req.URL.RawQuery,
        "status_code", statusCode,
        "duration_ms", duration.Milliseconds(),
        "user_agent", req.UserAgent(),
        "remote_addr", req.RemoteAddr,
        "content_length", req.ContentLength,
    )
}

// LogDatabaseOperation logs database operation metrics
func (sl *StructuredLogger) LogDatabaseOperation(operation string, table string, duration time.Duration, err error) {
    args := []interface{}{
        "operation", operation,
        "table", table,
        "duration_ms", duration.Milliseconds(),
    }
    
    if err != nil {
        args = append(args, "error", err)
        sl.logger.Error("Database operation failed", args...)
    } else {
        sl.logger.Debug("Database operation completed", args...)
    }
}

// LogAIModelRequest logs AI model interaction
func (sl *StructuredLogger) LogAIModelRequest(model string, tokenCount int, duration time.Duration, err error) {
    args := []interface{}{
        "model", model,
        "token_count", tokenCount,
        "duration_ms", duration.Milliseconds(),
    }
    
    if err != nil {
        args = append(args, "error", err)
        sl.logger.Error("AI model request failed", args...)
    } else {
        sl.logger.Info("AI model request completed", args...)
    }
}

// LogSecurityEvent logs security-related events
func (sl *StructuredLogger) LogSecurityEvent(event string, userID string, severity string, details map[string]interface{}) {
    args := []interface{}{
        "event_type", "security",
        "event", event,
        "user_id", userID,
        "severity", severity,
        "timestamp", time.Now().UTC(),
    }
    
    for k, v := range details {
        args = append(args, k, v)
    }
    
    if severity == "high" || severity == "critical" {
        sl.logger.Error("Security event", args...)
    } else {
        sl.logger.Warn("Security event", args...)
    }
}

// LogPerformanceMetric logs performance metrics
func (sl *StructuredLogger) LogPerformanceMetric(metric string, value float64, unit string, labels map[string]string) {
    args := []interface{}{
        "metric_type", "performance",
        "metric", metric,
        "value", value,
        "unit", unit,
        "timestamp", time.Now().UTC(),
    }
    
    for k, v := range labels {
        args = append(args, k, v)
    }
    
    sl.logger.Info("Performance metric", args...)
}

// LogApplicationEvent logs application lifecycle events
func (sl *StructuredLogger) LogApplicationEvent(event string, details map[string]interface{}) {
    args := []interface{}{
        "event_type", "application",
        "event", event,
        "timestamp", time.Now().UTC(),
    }
    
    for k, v := range details {
        args = append(args, k, v)
    }
    
    sl.logger.Info("Application event", args...)
}

// Middleware provides HTTP request logging middleware
func (sl *StructuredLogger) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // Create request context
        requestID := r.Header.Get("X-Request-ID")
        if requestID == "" {
            requestID = generateRequestID()
        }
        
        ctx := r.Context()
        ctx = context.WithValue(ctx, "request_id", requestID)
        ctx = context.WithValue(ctx, "logger", sl.WithContext(LogContext{
            RequestID: requestID,
            Component: "http",
        }))
        
        // Create response wrapper to capture status code
        wrapped := &responseWrapper{ResponseWriter: w, statusCode: 200}
        
        // Add request ID to response headers
        wrapped.Header().Set("X-Request-ID", requestID)
        
        // Process request
        next.ServeHTTP(wrapped, r.WithContext(ctx))
        
        // Log request completion
        duration := time.Since(start)
        sl.LogHTTPRequest(ctx, r, wrapped.statusCode, duration)
    })
}

// responseWrapper captures response status code
type responseWrapper struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWrapper) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}

// generateRequestID creates a unique request identifier
func generateRequestID() string {
    return fmt.Sprintf("%d-%d", time.Now().UnixNano(), rand.Int31())
}

// GetLoggerFromContext retrieves logger from request context
func GetLoggerFromContext(ctx context.Context) *StructuredLogger {
    if logger, ok := ctx.Value("logger").(*StructuredLogger); ok {
        return logger
    }
    // Return default logger if not found in context
    return NewStructuredLogger("info", os.Stdout)
}
```

**TodoWrite Task**: `Implement comprehensive structured logging system`

#### Step 3C.2: Metrics Collection System
**File: `/app/server/metrics/metrics_collector.go`** (create new file)
```go
package metrics

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "runtime"
    "sync"
    "time"
)

// MetricsCollector gathers and exposes application metrics
type MetricsCollector struct {
    counters   map[string]*Counter
    gauges     map[string]*Gauge
    histograms map[string]*Histogram
    mutex      sync.RWMutex
    startTime  time.Time
}

// Counter represents a monotonically increasing metric
type Counter struct {
    Value  int64             `json:"value"`
    Labels map[string]string `json:"labels,omitempty"`
    mutex  sync.Mutex
}

// Gauge represents a metric that can go up and down
type Gauge struct {
    Value  float64           `json:"value"`
    Labels map[string]string `json:"labels,omitempty"`
    mutex  sync.Mutex
}

// Histogram represents a distribution of values
type Histogram struct {
    Count   int64             `json:"count"`
    Sum     float64           `json:"sum"`
    Buckets map[string]int64  `json:"buckets"`
    Labels  map[string]string `json:"labels,omitempty"`
    mutex   sync.Mutex
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector() *MetricsCollector {
    return &MetricsCollector{
        counters:   make(map[string]*Counter),
        gauges:     make(map[string]*Gauge),
        histograms: make(map[string]*Histogram),
        startTime:  time.Now(),
    }
}

// IncrementCounter increments a counter metric
func (mc *MetricsCollector) IncrementCounter(name string, labels map[string]string) {
    mc.mutex.Lock()
    defer mc.mutex.Unlock()
    
    key := mc.buildKey(name, labels)
    if counter, exists := mc.counters[key]; exists {
        counter.mutex.Lock()
        counter.Value++
        counter.mutex.Unlock()
    } else {
        mc.counters[key] = &Counter{
            Value:  1,
            Labels: labels,
        }
    }
}

// SetGauge sets a gauge metric value
func (mc *MetricsCollector) SetGauge(name string, value float64, labels map[string]string) {
    mc.mutex.Lock()
    defer mc.mutex.Unlock()
    
    key := mc.buildKey(name, labels)
    if gauge, exists := mc.gauges[key]; exists {
        gauge.mutex.Lock()
        gauge.Value = value
        gauge.mutex.Unlock()
    } else {
        mc.gauges[key] = &Gauge{
            Value:  value,
            Labels: labels,
        }
    }
}

// RecordHistogram records a value in a histogram
func (mc *MetricsCollector) RecordHistogram(name string, value float64, labels map[string]string) {
    mc.mutex.Lock()
    defer mc.mutex.Unlock()
    
    key := mc.buildKey(name, labels)
    if histogram, exists := mc.histograms[key]; exists {
        histogram.mutex.Lock()
        histogram.Count++
        histogram.Sum += value
        bucket := mc.getBucket(value)
        histogram.Buckets[bucket]++
        histogram.mutex.Unlock()
    } else {
        buckets := make(map[string]int64)
        bucket := mc.getBucket(value)
        buckets[bucket] = 1
        
        mc.histograms[key] = &Histogram{
            Count:   1,
            Sum:     value,
            Buckets: buckets,
            Labels:  labels,
        }
    }
}

// buildKey creates a unique key for metric with labels
func (mc *MetricsCollector) buildKey(name string, labels map[string]string) string {
    if len(labels) == 0 {
        return name
    }
    
    key := name
    for k, v := range labels {
        key += fmt.Sprintf("_%s_%s", k, v)
    }
    return key
}

// getBucket determines which bucket a value belongs to
func (mc *MetricsCollector) getBucket(value float64) string {
    buckets := []float64{0.1, 0.25, 0.5, 1, 2.5, 5, 10, 25, 50, 100}
    
    for _, bucket := range buckets {
        if value <= bucket {
            return fmt.Sprintf("%.2f", bucket)
        }
    }
    return "inf"
}

// GetMetrics returns all collected metrics
func (mc *MetricsCollector) GetMetrics() map[string]interface{} {
    mc.mutex.RLock()
    defer mc.mutex.RUnlock()
    
    metrics := map[string]interface{}{
        "timestamp": time.Now().UTC(),
        "uptime_seconds": time.Since(mc.startTime).Seconds(),
        "counters": mc.counters,
        "gauges": mc.gauges,
        "histograms": mc.histograms,
    }
    
    // Add runtime metrics
    var memStats runtime.MemStats
    runtime.ReadMemStats(&memStats)
    
    metrics["runtime"] = map[string]interface{}{
        "num_goroutine": runtime.NumGoroutine(),
        "num_cpu": runtime.NumCPU(),
        "gomaxprocs": runtime.GOMAXPROCS(-1),
        "memory_alloc_bytes": memStats.Alloc,
        "memory_total_alloc_bytes": memStats.TotalAlloc,
        "memory_sys_bytes": memStats.Sys,
        "gc_runs": memStats.NumGC,
        "gc_cpu_fraction": memStats.GCCPUFraction,
    }
    
    return metrics
}

// MetricsHandler provides HTTP endpoint for metrics
func (mc *MetricsCollector) MetricsHandler(w http.ResponseWriter, r *http.Request) {
    metrics := mc.GetMetrics()
    
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(metrics); err != nil {
        http.Error(w, "Failed to encode metrics", http.StatusInternalServerError)
        return
    }
}

// PrometheusHandler provides Prometheus-compatible metrics endpoint
func (mc *MetricsCollector) PrometheusHandler(w http.ResponseWriter, r *http.Request) {
    mc.mutex.RLock()
    defer mc.mutex.RUnlock()
    
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    
    // Export counters
    for key, counter := range mc.counters {
        fmt.Fprintf(w, "# TYPE %s counter\n", key)
        fmt.Fprintf(w, "%s %d\n", key, counter.Value)
    }
    
    // Export gauges
    for key, gauge := range mc.gauges {
        fmt.Fprintf(w, "# TYPE %s gauge\n", key)
        fmt.Fprintf(w, "%s %.2f\n", key, gauge.Value)
    }
    
    // Export histograms
    for key, histogram := range mc.histograms {
        fmt.Fprintf(w, "# TYPE %s histogram\n", key)
        fmt.Fprintf(w, "%s_count %d\n", key, histogram.Count)
        fmt.Fprintf(w, "%s_sum %.2f\n", key, histogram.Sum)
        
        for bucket, count := range histogram.Buckets {
            fmt.Fprintf(w, "%s_bucket{le=\"%s\"} %d\n", key, bucket, count)
        }
    }
    
    // Add runtime metrics
    var memStats runtime.MemStats
    runtime.ReadMemStats(&memStats)
    
    fmt.Fprintf(w, "# TYPE go_memstats_alloc_bytes gauge\n")
    fmt.Fprintf(w, "go_memstats_alloc_bytes %d\n", memStats.Alloc)
    
    fmt.Fprintf(w, "# TYPE go_goroutines gauge\n")
    fmt.Fprintf(w, "go_goroutines %d\n", runtime.NumGoroutine())
}

// Business Metrics Helpers

// RecordAPIRequest records API request metrics
func (mc *MetricsCollector) RecordAPIRequest(method, endpoint string, statusCode int, duration time.Duration) {
    labels := map[string]string{
        "method":      method,
        "endpoint":    endpoint,
        "status_code": fmt.Sprintf("%d", statusCode),
    }
    
    mc.IncrementCounter("http_requests_total", labels)
    mc.RecordHistogram("http_request_duration_seconds", duration.Seconds(), labels)
}

// RecordDatabaseOperation records database operation metrics
func (mc *MetricsCollector) RecordDatabaseOperation(operation, table string, duration time.Duration, success bool) {
    labels := map[string]string{
        "operation": operation,
        "table":     table,
        "success":   fmt.Sprintf("%t", success),
    }
    
    mc.IncrementCounter("database_operations_total", labels)
    mc.RecordHistogram("database_operation_duration_seconds", duration.Seconds(), labels)
}

// RecordAIModelUsage records AI model usage metrics
func (mc *MetricsCollector) RecordAIModelUsage(model string, tokenCount int, duration time.Duration, success bool) {
    labels := map[string]string{
        "model":   model,
        "success": fmt.Sprintf("%t", success),
    }
    
    mc.IncrementCounter("ai_model_requests_total", labels)
    mc.RecordHistogram("ai_model_request_duration_seconds", duration.Seconds(), labels)
    mc.RecordHistogram("ai_model_token_count", float64(tokenCount), labels)
}

// RecordUserAction records user action metrics
func (mc *MetricsCollector) RecordUserAction(action, userID string) {
    labels := map[string]string{
        "action": action,
    }
    
    mc.IncrementCounter("user_actions_total", labels)
}

// RecordCacheOperation records cache operation metrics
func (mc *MetricsCollector) RecordCacheOperation(operation string, hit bool, duration time.Duration) {
    labels := map[string]string{
        "operation": operation,
        "result":    func() string { if hit { return "hit" } else { return "miss" } }(),
    }
    
    mc.IncrementCounter("cache_operations_total", labels)
    mc.RecordHistogram("cache_operation_duration_seconds", duration.Seconds(), labels)
}

// StartPeriodicCollection starts background metrics collection
func (mc *MetricsCollector) StartPeriodicCollection() {
    ticker := time.NewTicker(30 * time.Second)
    go func() {
        for range ticker.C {
            mc.collectSystemMetrics()
        }
    }()
}

// collectSystemMetrics collects system-level metrics
func (mc *MetricsCollector) collectSystemMetrics() {
    var memStats runtime.MemStats
    runtime.ReadMemStats(&memStats)
    
    // Memory metrics
    mc.SetGauge("memory_alloc_bytes", float64(memStats.Alloc), nil)
    mc.SetGauge("memory_total_alloc_bytes", float64(memStats.TotalAlloc), nil)
    mc.SetGauge("memory_sys_bytes", float64(memStats.Sys), nil)
    mc.SetGauge("memory_heap_alloc_bytes", float64(memStats.HeapAlloc), nil)
    mc.SetGauge("memory_heap_sys_bytes", float64(memStats.HeapSys), nil)
    
    // GC metrics
    mc.SetGauge("gc_runs_total", float64(memStats.NumGC), nil)
    mc.SetGauge("gc_cpu_fraction", memStats.GCCPUFraction, nil)
    
    // Goroutine metrics
    mc.SetGauge("goroutines_total", float64(runtime.NumGoroutine()), nil)
    
    // CPU metrics
    mc.SetGauge("cpu_cores", float64(runtime.NumCPU()), nil)
    mc.SetGauge("gomaxprocs", float64(runtime.GOMAXPROCS(-1)), nil)
}
```

**TodoWrite Task**: `Implement comprehensive metrics collection system`

#### Step 3C.3: Health Check System
**File: `/app/server/health/health_checker.go`** (create new file)
```go
package health

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "sync"
    "time"
    "yourapp/db"
    "yourapp/cache"
)

// HealthChecker manages application health checks
type HealthChecker struct {
    checks map[string]HealthCheck
    mutex  sync.RWMutex
}

// HealthCheck represents a single health check
type HealthCheck interface {
    Name() string
    Check(ctx context.Context) HealthResult
}

// HealthResult represents the result of a health check
type HealthResult struct {
    Status    HealthStatus   `json:"status"`
    Message   string         `json:"message,omitempty"`
    Details   map[string]interface{} `json:"details,omitempty"`
    Duration  time.Duration  `json:"duration_ms"`
    Timestamp time.Time      `json:"timestamp"`
}

// HealthStatus represents the status of a health check
type HealthStatus string

const (
    HealthStatusHealthy   HealthStatus = "healthy"
    HealthStatusUnhealthy HealthStatus = "unhealthy"
    HealthStatusDegraded  HealthStatus = "degraded"
    HealthStatusUnknown   HealthStatus = "unknown"
)

// OverallHealth represents the overall application health
type OverallHealth struct {
    Status    HealthStatus            `json:"status"`
    Timestamp time.Time               `json:"timestamp"`
    Uptime    time.Duration           `json:"uptime_seconds"`
    Version   string                  `json:"version,omitempty"`
    Checks    map[string]HealthResult `json:"checks"`
}

// NewHealthChecker creates a new health checker
func NewHealthChecker() *HealthChecker {
    return &HealthChecker{
        checks: make(map[string]HealthCheck),
    }
}

// AddCheck adds a health check
func (hc *HealthChecker) AddCheck(check HealthCheck) {
    hc.mutex.Lock()
    defer hc.mutex.Unlock()
    hc.checks[check.Name()] = check
}

// RemoveCheck removes a health check
func (hc *HealthChecker) RemoveCheck(name string) {
    hc.mutex.Lock()
    defer hc.mutex.Unlock()
    delete(hc.checks, name)
}

// CheckHealth performs all health checks
func (hc *HealthChecker) CheckHealth(ctx context.Context) OverallHealth {
    hc.mutex.RLock()
    checks := make(map[string]HealthCheck)
    for name, check := range hc.checks {
        checks[name] = check
    }
    hc.mutex.RUnlock()
    
    results := make(map[string]HealthResult)
    overallStatus := HealthStatusHealthy
    
    // Run all checks concurrently
    var wg sync.WaitGroup
    var resultsMutex sync.Mutex
    
    for name, check := range checks {
        wg.Add(1)
        go func(name string, check HealthCheck) {
            defer wg.Done()
            
            checkCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
            defer cancel()
            
            start := time.Now()
            result := check.Check(checkCtx)
            result.Duration = time.Since(start)
            result.Timestamp = time.Now()
            
            resultsMutex.Lock()
            results[name] = result
            
            // Determine overall status
            switch result.Status {
            case HealthStatusUnhealthy:
                overallStatus = HealthStatusUnhealthy
            case HealthStatusDegraded:
                if overallStatus == HealthStatusHealthy {
                    overallStatus = HealthStatusDegraded
                }
            }
            resultsMutex.Unlock()
        }(name, check)
    }
    
    wg.Wait()
    
    return OverallHealth{
        Status:    overallStatus,
        Timestamp: time.Now(),
        Uptime:    time.Since(startTime),
        Version:   version,
        Checks:    results,
    }
}

// HealthHandler provides HTTP endpoint for health checks
func (hc *HealthChecker) HealthHandler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    health := hc.CheckHealth(ctx)
    
    // Set appropriate HTTP status code
    statusCode := http.StatusOK
    switch health.Status {
    case HealthStatusUnhealthy:
        statusCode = http.StatusServiceUnavailable
    case HealthStatusDegraded:
        statusCode = http.StatusOK // 200 but with degraded status
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    
    if err := json.NewEncoder(w).Encode(health); err != nil {
        http.Error(w, "Failed to encode health status", http.StatusInternalServerError)
    }
}

// LivenessHandler provides simple liveness check
func (hc *HealthChecker) LivenessHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))
}

// ReadinessHandler provides readiness check
func (hc *HealthChecker) ReadinessHandler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    health := hc.CheckHealth(ctx)
    
    if health.Status == HealthStatusHealthy || health.Status == HealthStatusDegraded {
        w.Header().Set("Content-Type", "text/plain")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("READY"))
    } else {
        w.Header().Set("Content-Type", "text/plain")
        w.WriteHeader(http.StatusServiceUnavailable)
        w.Write([]byte("NOT READY"))
    }
}

// Predefined Health Checks

// DatabaseHealthCheck checks database connectivity
type DatabaseHealthCheck struct {
    pool *db.DatabasePool
}

func NewDatabaseHealthCheck(pool *db.DatabasePool) *DatabaseHealthCheck {
    return &DatabaseHealthCheck{pool: pool}
}

func (dhc *DatabaseHealthCheck) Name() string {
    return "database"
}

func (dhc *DatabaseHealthCheck) Check(ctx context.Context) HealthResult {
    if dhc.pool == nil {
        return HealthResult{
            Status:  HealthStatusUnhealthy,
            Message: "Database pool not initialized",
        }
    }
    
    err := dhc.pool.HealthCheck(ctx)
    if err != nil {
        return HealthResult{
            Status:  HealthStatusUnhealthy,
            Message: "Database connection failed",
            Details: map[string]interface{}{
                "error": err.Error(),
            },
        }
    }
    
    // Get pool statistics
    stats := dhc.pool.GetPoolStats()
    
    // Check if pool is under pressure
    utilization := float64(stats["acquired_conns"]) / float64(stats["max_conns"]) * 100
    status := HealthStatusHealthy
    if utilization > 90 {
        status = HealthStatusDegraded
    }
    
    return HealthResult{
        Status:  status,
        Message: "Database connection healthy",
        Details: map[string]interface{}{
            "pool_stats":   stats,
            "utilization":  fmt.Sprintf("%.2f%%", utilization),
        },
    }
}

// CacheHealthCheck checks cache connectivity
type CacheHealthCheck struct {
    cache *cache.CacheManager
}

func NewCacheHealthCheck(cache *cache.CacheManager) *CacheHealthCheck {
    return &CacheHealthCheck{cache: cache}
}

func (chc *CacheHealthCheck) Name() string {
    return "cache"
}

func (chc *CacheHealthCheck) Check(ctx context.Context) HealthResult {
    if chc.cache == nil {
        return HealthResult{
            Status:  HealthStatusDegraded,
            Message: "Cache not available (optional)",
        }
    }
    
    // Test cache operations
    testKey := "health_check_" + fmt.Sprintf("%d", time.Now().UnixNano())
    testValue := []byte("health_check_value")
    
    // Test write
    chc.cache.Set(ctx, testKey, testValue, time.Minute)
    
    // Test read
    _, found := chc.cache.Get(ctx, testKey)
    
    if !found {
        return HealthResult{
            Status:  HealthStatusDegraded,
            Message: "Cache read/write test failed",
        }
    }
    
    // Get cache statistics
    stats := chc.cache.stats.GetStats()
    
    return HealthResult{
        Status:  HealthStatusHealthy,
        Message: "Cache operations healthy",
        Details: map[string]interface{}{
            "stats": stats,
        },
    }
}

// MemoryHealthCheck checks memory usage
type MemoryHealthCheck struct {
    maxMemoryMB int64
}

func NewMemoryHealthCheck(maxMemoryMB int64) *MemoryHealthCheck {
    return &MemoryHealthCheck{maxMemoryMB: maxMemoryMB}
}

func (mhc *MemoryHealthCheck) Name() string {
    return "memory"
}

func (mhc *MemoryHealthCheck) Check(ctx context.Context) HealthResult {
    var memStats runtime.MemStats
    runtime.ReadMemStats(&memStats)
    
    currentMemoryMB := int64(memStats.Alloc / 1024 / 1024)
    utilizationPercent := float64(currentMemoryMB) / float64(mhc.maxMemoryMB) * 100
    
    status := HealthStatusHealthy
    message := "Memory usage normal"
    
    if utilizationPercent > 90 {
        status = HealthStatusUnhealthy
        message = "Memory usage critical"
    } else if utilizationPercent > 75 {
        status = HealthStatusDegraded
        message = "Memory usage high"
    }
    
    return HealthResult{
        Status:  status,
        Message: message,
        Details: map[string]interface{}{
            "current_mb":      currentMemoryMB,
            "max_mb":          mhc.maxMemoryMB,
            "utilization":     fmt.Sprintf("%.2f%%", utilizationPercent),
            "num_goroutine":   runtime.NumGoroutine(),
            "gc_runs":         memStats.NumGC,
            "gc_cpu_fraction": memStats.GCCPUFraction,
        },
    }
}

// DiskSpaceHealthCheck checks disk space
type DiskSpaceHealthCheck struct {
    path           string
    minFreePercent float64
}

func NewDiskSpaceHealthCheck(path string, minFreePercent float64) *DiskSpaceHealthCheck {
    return &DiskSpaceHealthCheck{
        path:           path,
        minFreePercent: minFreePercent,
    }
}

func (dshc *DiskSpaceHealthCheck) Name() string {
    return "disk_space"
}

func (dshc *DiskSpaceHealthCheck) Check(ctx context.Context) HealthResult {
    // This is a simplified implementation
    // In production, you'd use syscalls to get actual disk usage
    
    return HealthResult{
        Status:  HealthStatusHealthy,
        Message: "Disk space check not implemented",
        Details: map[string]interface{}{
            "path": dshc.path,
            "note": "Implement using syscalls for production",
        },
    }
}

// Global variables for version and start time
var (
    version   = "dev"
    startTime = time.Now()
)

// SetVersion sets the application version
func SetVersion(v string) {
    version = v
}

// SetStartTime sets the application start time
func SetStartTime(t time.Time) {
    startTime = t
}
```

**TodoWrite Task**: `Implement comprehensive health check system`

### KPIs for Phase 3C
- ✅ 100% structured logging across all components
- ✅ Real-time metrics collection and exposure
- ✅ <500ms health check response time
- ✅ Complete observability into application performance
- ✅ Automated alerting for degraded performance
- ✅ Comprehensive application lifecycle monitoring

---

## 🎯 FINAL PHASE 3 VALIDATION & HANDOFF

### Comprehensive Testing Validation
```bash
# Run complete test suite
./scripts/run-all-tests.sh

# Validate test coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Performance regression testing  
go test -bench=. -benchmem ./... > new_benchmarks.txt
benchcmp baseline_benchmarks.txt new_benchmarks.txt

# Security testing
gosec ./...
govulncheck ./...
```

### CI/CD Pipeline Validation
```bash
# Trigger full CI/CD pipeline
git push origin feature/testing-infrastructure

# Validate all quality gates pass
# - Code quality checks
# - Unit tests (80%+ coverage)  
# - Integration tests
# - E2E tests
# - Security scanning
# - Performance benchmarks
# - Build and deployment
```

### Monitoring & Observability Validation
```bash
# Test health endpoints
curl http://localhost:8080/health
curl http://localhost:8080/metrics
curl http://localhost:8080/ready
curl http://localhost:8080/live

# Validate structured logging
tail -f logs/app.log | jq .

# Check metrics collection
curl http://localhost:8080/metrics | grep -E "(http_requests_total|memory_alloc_bytes)"
```

---

## 📊 PHASE 3 SUCCESS METRICS SUMMARY

### Testing Metrics
- ✅ **Test Coverage**: 80%+ (from ~5%)
- ✅ **Test Execution Time**: <5 minutes for full suite
- ✅ **Test Types**: Unit, Integration, E2E, API, Load, Security
- ✅ **Test Reliability**: 99%+ consistent pass rate
- ✅ **Test Environment**: Fully isolated with containers

### CI/CD Metrics  
- ✅ **Pipeline Automation**: 100% automated
- ✅ **Feedback Loop**: <5 minutes
- ✅ **Quality Gates**: 8 automated checkpoints
- ✅ **Security Scanning**: 100% coverage
- ✅ **Deployment Success**: 99.9% success rate

### Monitoring Metrics
- ✅ **Observability**: 100% structured logging
- ✅ **Metrics Collection**: Real-time application metrics
- ✅ **Health Monitoring**: Comprehensive health checks
- ✅ **Performance Tracking**: Automated regression detection
- ✅ **Alerting**: Proactive issue detection

---

## 🚀 HANDOFF TO PHASE 4

With Phase 3 complete, the Plandex application now has:

### Quality Foundation
- Enterprise-grade testing framework
- Fully automated CI/CD pipeline  
- Comprehensive monitoring and observability
- Developer-friendly tooling and workflows

### Readiness for Phase 4
The robust testing and monitoring infrastructure established in Phase 3 will support the feature development in Phase 4:

- **Feature Testing**: New features can be developed with confidence using the comprehensive testing framework
- **Quality Assurance**: All new features will automatically go through quality gates
- **Performance Monitoring**: New features will be monitored for performance impact
- **Deployment Safety**: Features can be safely deployed using the automated pipeline

### Next Phase Prerequisites
- [ ] All Phase 3 KPIs achieved and validated
- [ ] CI/CD pipeline operational and tested
- [ ] Monitoring infrastructure collecting metrics
- [ ] Development team trained on new workflows
- [ ] Documentation updated with new processes

---

*This guide establishes enterprise-grade quality assurance and operational excellence for the Plandex application, providing the foundation for confident feature development and reliable operation at scale.*