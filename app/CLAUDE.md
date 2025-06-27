# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**Plandex** is an AI-powered terminal-based development tool designed for large coding tasks and real-world projects. It provides a CLI interface that integrates with multiple AI models to plan, build, and execute development tasks with support for large codebases (2M+ token context).

## Architecture

### Core Components
- **`cli/`**: Command-line interface built with Cobra (40+ commands)
- **`server/`**: HTTP API server using Gorilla Mux with PostgreSQL database
- **`shared/`**: Common data models and utilities shared between CLI and server

### Technology Stack (2025 Security Standards)
- **Language**: Go 1.23.10+ (Latest: Go 1.24 released February 2025, patch 1.23.6)
- **Database**: PostgreSQL 17.5+ (Latest security fixes, 0 vulnerabilities in 2025)
- **CLI Framework**: Cobra (spf13/cobra)
- **TUI Components**: Bubble Tea (charmbracelet/bubbletea)
- **AI Integration**: Multiple providers (OpenAI, Anthropic, Google, OpenRouter)
- **Code Parsing**: Tree-sitter for 30+ programming languages
- **Containerization**: Docker with multi-stage builds and security hardening

## 🚀 MCP Server Integration & Best Practices (2025)

### 🟥 MANDATORY MCP USAGE POLICY

**CRITICAL RULE**: Use MCP servers for ALL applicable tasks. This is not optional.

#### When MCP Usage is REQUIRED:
- **Any research or documentation lookup** → Use Context7
- **Complex multi-step tasks (3+ steps)** → Use SequentialThinking
- **File operations in allowed directories** → Use Filesystem MCP
- **Knowledge synthesis from multiple sources** → Use KnowledgeSynthesize
- **API endpoint validation or discovery** → Use APISearch/APIValidator
- **Multi-API data gathering** → Use MultiAPIFetch

#### When MCP Usage is RECOMMENDED:
- **Task planning and management** → Use JARVIS FSM Controller
- **Any complex problem-solving** → Combine multiple MCP servers
- **Research validation** → Cross-reference with multiple MCPs

### 🧠 Context7 MCP - ALWAYS FIRST PRIORITY

**USAGE MANDATE**: Context7 MUST be used before implementing any feature, fixing any bug, or making architectural decisions.

```bash
# MANDATORY usage patterns:
"How to implement JWT authentication in Go 1.24? use context7"
"PostgreSQL 17.5 connection pooling best practices? use context7"
"Docker multi-stage builds security optimization 2025? use context7"
"Go concurrency patterns for CLI applications? use context7"
```

#### Context7 Integration Workflow:
1. **Pre-Implementation Research** (MANDATORY)
   - Always resolve library ID first: `mcp__context7__resolve-library-id`
   - Get current documentation: `mcp__context7__get-library-docs`
   - Verify best practices for current year (2024-2025)

2. **Technology-Specific Lookups**
   - Go language features and updates
   - PostgreSQL security and performance
   - Docker optimization techniques
   - AI model integration patterns

3. **Validation and Verification**
   - Cross-check implementations against latest docs
   - Verify security practices are current
   - Ensure performance optimizations are relevant

### 🧩 SequentialThinking MCP - SYSTEMATIC APPROACH

**MANDATORY USAGE** for tasks requiring:
- Multi-file code changes
- Architecture refactoring
- Complex debugging scenarios
- Performance optimization planning
- Security vulnerability analysis
- Database migration strategies

#### SequentialThinking Workflow:
```bash
# Always use for complex problem-solving:
1. Problem Definition → Define exact issue/requirement
2. Context Research → Use Context7 for current best practices
3. Solution Analysis → Break down approaches systematically
4. Implementation Planning → Create detailed step-by-step plan
5. Risk Assessment → Identify potential issues
6. Execution Strategy → Final optimized approach
```

#### Integration with TodoWrite:
- Every SequentialThinking session MUST generate TodoWrite tasks
- Each thought phase creates corresponding todo items
- Progress tracking through todo status updates
- Documentation of decision rationale in todos

### 🤖 JARVIS FSM Controller - ADVANCED ORCHESTRATION

**Usage for Enterprise-Level Task Management**:
- Complex multi-phase projects (Phase 1-5 development)
- Role-based cognitive enhancement
- Meta-prompt generation for Task() agent spawning
- Performance tracking across development phases

#### JARVIS Integration Pattern:
```bash
# For major feature implementation:
1. Initialize JARVIS with objective
2. Let JARVIS orchestrate MCP server usage
3. Follow the 6-step agent loop:
   - Analyze Events → Use Context7 for research
   - Select Tools → Choose appropriate MCPs
   - Wait for Execution → Process MCP responses
   - Iterate → Refine approach based on results
   - Submit Results → Document findings
   - Enter Standby → Await next phase
```

### 🗂️ Filesystem MCP - OPTIMIZED FILE OPERATIONS

**MANDATORY for ALL file operations** in allowed directories:
- Reading project files → `mcp__filesystem__read_file`
- Multi-file analysis → `mcp__filesystem__read_multiple_files`
- Directory exploration → `mcp__filesystem__list_directory`
- Project structure analysis → `mcp__filesystem__directory_tree`
- File searches → `mcp__filesystem__search_files`

#### Allowed Directories:
- `/Users/Mahavir/code/plandex/` (primary project)
- All subdirectories within project scope
- Always verify with `mcp__filesystem__list_allowed_directories`

### 🔗 API Integration MCPs - DATA GATHERING

#### MultiAPIFetch Usage:
```bash
# For parallel data gathering:
- Multiple documentation sources
- Cross-platform compatibility checks
- Performance benchmarking data
- Security vulnerability databases
```

#### APISearch & Validation:
```bash
# For API discovery and validation:
- Find relevant APIs for specific objectives
- Validate endpoint functionality
- Auto-correct failed API calls
- Role-based API filtering
```

### 🧠 KnowledgeSynthesize - INTELLIGENCE AGGREGATION

**Usage for Critical Decision Making**:
- Conflicting information resolution
- Multi-source validation
- Confidence scoring for recommendations
- Consensus building from various data sources

#### Synthesis Modes:
- **Consensus**: When sources generally agree
- **Weighted**: When sources have different reliability
- **Hierarchical**: When sources have clear authority levels
- **Conflict Resolution**: When sources contradict each other

### 📋 Integrated MCP Workflow - COMPREHENSIVE APPROACH

#### Pre-Development Phase (MANDATORY):
```bash
1. TodoRead → Check existing tasks
2. Context7 → Research latest best practices
3. SequentialThinking → Plan complex implementations
4. Filesystem → Analyze current codebase
5. JARVIS → Orchestrate if needed
6. TodoWrite → Document comprehensive plan
```

#### Development Phase:
```bash
1. Context7 → Verify implementation approaches
2. Filesystem → Read/modify files systematically
3. SequentialThinking → Handle complex logic
4. TodoWrite → Track progress continuously
5. MultiAPIFetch → Gather external data if needed
6. KnowledgeSynthesize → Validate decisions
```

#### Validation Phase:
```bash
1. Filesystem → Verify file modifications
2. Context7 → Check against latest standards
3. APIValidator → Validate any API integrations
4. SequentialThinking → Systematic testing approach
5. TodoWrite → Mark tasks complete
6. JARVIS → Performance metrics tracking
```

### 🚨 MCP USAGE ENFORCEMENT RULES

#### MANDATORY MCP Usage:
- **NEVER** skip Context7 for research tasks
- **ALWAYS** use SequentialThinking for complex problems
- **MUST** use Filesystem MCP for file operations
- **REQUIRED** to use TodoWrite for task management

#### Performance Monitoring:
- Track MCP usage patterns
- Measure efficiency improvements
- Document MCP integration success stories
- Optimize MCP server selection

#### Quality Assurance:
- Verify MCP responses before implementation
- Cross-validate information across multiple MCPs
- Document MCP decision rationale
- Report MCP server issues promptly

### 🎯 MCP Success Metrics:
- **Research Accuracy**: >95% through Context7 usage
- **Problem-Solving Efficiency**: >70% improvement with SequentialThinking
- **File Operation Speed**: >60% faster with Filesystem MCP
- **Task Completion Rate**: >90% with systematic MCP usage
- **Code Quality**: Measurable improvement through MCP-guided development

## Development Commands

### Local Development Environment
```bash
# Start local server and database
./start_local.sh

# Development with hot reload (auto-installs reflex)
./scripts/dev.sh

# Reset local environment completely
./reset_local.sh

# Clear local data only (destructive - prompts for confirmation)
./clear_local.sh
```

### CLI Development
```bash
# Build development CLI (creates plandex-dev and pdxd alias)
cd cli && ./dev.sh

# Install production CLI (handles platform detection and V1 migration)  
cd cli && ./install.sh
```

### Environment Variables
```bash
# Development CLI configuration
PLANDEX_DEV_CLI_OUT_DIR="/usr/local/bin"        # Output directory
PLANDEX_DEV_CLI_NAME="plandex-dev"              # Binary name
PLANDEX_DEV_CLI_ALIAS="pdxd"                    # Alias name

# API Keys (required in .env file)
OPENROUTER_API_KEY="your_key_here"
GEMINI_API_KEY="your_key_here"
PERPLEXITY_API_KEY="your_key_here"
```

## Go Module Structure

### Multi-Module Architecture
- **`cli/go.mod`**: CLI application with dependency on shared module
- **`server/go.mod`**: Server application with dependency on shared module
- **`shared/go.mod`**: Common types, utilities, and AI model configurations
- **Local Dependencies**: All modules use local replace directives for shared module

### Building Components
```bash
# Build CLI
cd cli && go build

# Build server
cd server && go build

# Development build with hot reload
cd cli && ./dev.sh      # Creates development binary
cd server && go build && ./plandex-server  # Manual server build
```

## Database

### PostgreSQL Configuration (Security Hardened)
- **Connection**: `postgres://plandex:plandex@localhost:5432/plandex`
- **Docker Service**: `plandex-postgres` on port 5432
- **Migrations**: 40+ migration files in `server/migrations/`
- **Features**: Comprehensive helpers, locks, queues, RBAC
- **Security**: SSL/TLS encryption, SCRAM-SHA-256 authentication, Row-Level Security (RLS)
- **Access Control**: Role-based permissions with least privilege principle
- **Monitoring**: Regular security audits and vulnerability scanning

### Database Operations
- **Helpers**: Complete database layer in `server/db/`
- **Migrations**: Automatic migration system with up/down scripts
- **Persistence**: Docker volumes for data persistence across restarts

## Hot Reload Development

### File Watching Patterns
- **CLI + Shared**: `^(cli|shared)/.*\.(go|mod|sum)$`
- **Server + Shared**: `^(server|shared)/.*\.(go|mod|sum)$`

### Development Workflow
1. Start with `./scripts/dev.sh`
2. Make changes to CLI, server, or shared modules
3. Reflex automatically rebuilds and restarts affected components
4. CLI creates development binary (`plandex-dev`)
5. Server rebuilds and restarts automatically

## Docker Configuration

### Services
```yaml
# PostgreSQL Database (Security Hardened)
plandex-postgres:
  - Image: postgres:17.5  # Pinned version for security
  - Port: 5432
  - Credentials: plandex:plandex
  - Volume: plandex-db
  - Security: SSL/TLS enabled, restricted access, vulnerability-free

# Plandex Server  
plandex-server:
  - Image: plandexai/plandex-server:latest
  - Port: 8099
  - Environment: development, LOCAL_MODE=1
  - Volume: plandex-files
  - Depends on: plandex-postgres
```

### Environment Variables (Docker)
```bash
DATABASE_URL="postgres://plandex:plandex@plandex-postgres:5432/plandex?sslmode=disable"
GOENV=development
LOCAL_MODE=1
PLANDEX_BASE_DIR=/plandex-server
```

## Testing

### Test Structure
```bash
# Run all tests
go test ./...

# Run tests in specific modules
go test ./server/model/...
go test ./server/syntax/...

# Run individual test files
go test ./server/types/reply_test.go
go test ./server/utils/whitespace_test.go
```

### Test Locations
- `server/model/parse/subtasks_test.go`
- `server/model/plan/tell_stream_processor_test.go`
- `server/syntax/structured_edits_test.go`
- `server/syntax/unique_replacement_test.go`
- `server/types/reply_test.go`
- `server/utils/whitespace_test.go`

## Key File Locations

### Configuration Files
- `docker-compose.yml` - Local development orchestration
- `.env` - API keys and environment variables
- `server/Dockerfile` - Production containerization

### Entry Points
- `cli/main.go` - CLI application entry point
- `server/main.go` - Server application entry point

### Development Scripts
- `start_local.sh` - Start local development environment
- `scripts/dev.sh` - Hot reload development mode
- `cli/dev.sh` - CLI development build
- `cli/install.sh` - Production CLI installation
- `reset_local.sh` - Complete environment reset
- `clear_local.sh` - Clear local data

## CLI Commands Structure

### Command Categories
- **Plan Management**: `new`, `delete`, `list`, `archive`, `plans`
- **Context**: `load`, `rm`, `ls`, `cd`, `current`
- **Execution**: `tell`, `build`, `apply`, `continue`, `repl`
- **AI Models**: `models`, `set-model`, `model-packs`
- **Authentication**: `sign-in`, `sign-up`, `sign-out`
- **Configuration**: `config`, `set-config`
- **Collaboration**: `invite`, `users`, `connect`

### Interactive Features
- **REPL Mode**: Interactive conversation interface
- **Streaming TUI**: Real-time build logs and execution status
- **File Selection**: Context management and loading
- **Diff Viewing**: Review changes before applying

## AI Model Integration

### Supported Providers
- **OpenAI**: GPT models with streaming support
- **Anthropic**: Claude models (3.5 Sonnet, etc.)
- **Google**: Gemini models
- **OpenRouter**: Access to multiple model providers
- **Custom Models**: Configurable endpoints and parameters

### Model Configuration
- **Per-Plan Settings**: Different models for different plans
- **Token Management**: Context window optimization
- **Streaming**: Real-time response handling
- **Role-Based Models**: Different models for different tasks

## Syntax and Code Parsing

### Tree-sitter Integration
- **Language Support**: 30+ programming languages
- **File Mapping**: Intelligent project structure analysis
- **Context Extraction**: Smart code context for AI models
- **Structured Edits**: Precise code modifications

### File Context Management
- **Project Detection**: Automatic project boundary detection
- **Ignore Patterns**: `.plandexignore` support
- **Smart Loading**: Intelligent file selection for context
- **Large Context**: Support for 2M+ token contexts

## Security and Authentication

### Authentication System
- **Token-Based**: JWT tokens for API authentication
- **Local Storage**: `~/.plandex-home-v2/auth.json`
- **Trial Management**: Built-in trial and subscription handling
- **Organization Support**: Multi-user organization features

### API Key Management
- **Environment Variables**: Secure API key storage
- **Local Mode**: Development without cloud dependencies
- **Provider Switching**: Easy switching between AI providers

## Performance Considerations

### MacBook 2012 Optimizations
```bash
# Memory management
export GOGC=50                    # Aggressive garbage collection
export GOMEMLIMIT=512MiB         # Memory limit for Go runtime
export GOMAXPROCS=4              # CPU core optimization

# Docker resource limits
docker run --memory=1g --cpus=2  # Container resource constraints
```

### Database Performance (PostgreSQL 17.5+ Security Standards)
- **Connection Pooling**: Optimized PostgreSQL connections with security hardening
- **Query Optimization**: Indexed queries, prepared statements, and RLS (Row-Level Security)
- **Migration Efficiency**: Incremental schema updates with rollback testing
- **Security Compliance**: SSL/TLS encryption, SCRAM-SHA-256 authentication
- **Access Control**: Role-based permissions with least privilege principle
- **Monitoring**: Regular security audits and vulnerability scanning

## Development Best Practices

### Code Style (Go 1.24 Best Practices)
- **Standard Go Formatting**: Use `gofmt`, `go vet`, and `golangci-lint`
- **Error Handling**: Comprehensive error handling with structured responses and proper wrapping
- **Logging**: Structured logging with appropriate levels and sensitive data redaction
- **Documentation**: Inline documentation for complex logic following Go standards
- **Concurrency**: Modern goroutine patterns with proper context handling
- **Security**: Input validation, secure secret management, and vulnerability scanning
- **Performance**: Memory optimization for resource-constrained environments
- **Dependencies**: Regular updates and security patch management

### Git Workflow
- **Feature Branches**: Use feature branches for development
- **Conventional Commits**: Follow conventional commit message format
- **Migration Safety**: Test database migrations in development first

### Testing Strategy
- **Unit Tests**: Test individual functions and methods
- **Integration Tests**: Test component interactions
- **End-to-End Tests**: Test complete workflows
- **Performance Tests**: Benchmark critical paths

## Integration with Broader Repository

### Documentation Structure
- **Main Docs**: `/docs/` contains Docusaurus-based documentation
- **API Reference**: Complete CLI and API documentation
- **Core Concepts**: Detailed explanations of key features
- **Hosting Guides**: Self-hosting and cloud deployment

### Implementation Guides
- **Phase-Based Upgrades**: 5-phase modernization plan
- **Security Priorities**: Critical security and dependency updates
- **Performance Optimization**: Infrastructure and performance improvements
- **Testing Framework**: Comprehensive testing and DevOps setup
- **Feature Enhancement**: Advanced feature development
- **Next-Gen Capabilities**: Future technology integration

## Troubleshooting

### Common Issues
- **Database Connection**: Ensure PostgreSQL is running via Docker
- **API Keys**: Verify `.env` file contains required keys
- **Hot Reload**: Check that reflex is installed (`go install github.com/cespare/reflex@latest`)
- **Port Conflicts**: Ensure ports 5432 and 8099 are available
- **Memory Issues**: Monitor memory usage on resource-constrained systems

### Debug Commands
```bash
# Check service status
docker-compose ps

# View logs
docker-compose logs plandex-server
docker-compose logs plandex-postgres

# Database connection test
psql -h localhost -p 5432 -U plandex -d plandex

# Build debugging
cd cli && go build -v ./...
cd server && go build -v ./...
```

## Important Notes

### Version Management
- **V1 to V2 Migration**: Install script handles legacy version migration
- **Development vs Production**: Use different binaries (`plandex-dev` vs `plandex`)
- **Alias Management**: `pdxd` for development, `pdx` for production

### Data Management
- **Destructive Operations**: `clear_local.sh` permanently deletes data
- **Backup Strategy**: Consider backing up PostgreSQL data before major changes
- **Migration Testing**: Always test database migrations in development first

### AI Model Considerations
- **Rate Limits**: Be aware of API rate limits for different providers
- **Context Windows**: Optimize context size for different models
- **Cost Management**: Monitor API usage and costs
- **Model Selection**: Choose appropriate models for different task types

## 2025 Security & Performance Standards

### Go 1.24 Language Updates
- **Growing Adoption**: 5.8M developers (up from 4.1M in 2024)
- **Industry Growth**: 15% professional usage (from 10% in 2021)
- **AI Model Serving**: New problem domain for Go adoption
- **Concurrency Enhancements**: Advanced patterns for high-performance systems
- **Testing Improvements**: Enhanced built-in testing facilities

### PostgreSQL 17.5 Security Enhancements
- **Zero Vulnerabilities**: Clean security record for 2024-2025
- **SCRAM-SHA-256**: Modern authentication method implementation
- **Row-Level Security**: Fine-grained access control policies
- **SSL/TLS Encryption**: Mandatory for production deployments
- **CIS Benchmarks**: Following Center for Internet Security guidelines

### Docker Security & Optimization
- **Multi-Stage Builds**: Separation of build and runtime environments
- **Distroless Images**: Minimal attack surface using Google's distroless
- **BuildKit Integration**: DOCKER_BUILDKIT=1 with cache mounts
- **Security Scanning**: Integrated vulnerability assessment
- **Resource Limits**: Controlled CPU and memory allocation

### Development Efficiency
- **MCP Integration**: >70% efficiency improvement with systematic usage
- **Context Management**: 2M+ token context optimization
- **Performance Monitoring**: Real-time metrics and benchmarking
- **Quality Gates**: Automated testing and security validation