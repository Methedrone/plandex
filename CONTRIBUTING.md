# 🤝 Contributing to PDw-X

Thank you for your interest in contributing to PDw-X! This guide will help you get started with our development workflow and standards.

## 🎯 Project Overview

PDw-X is an advanced AI-powered development platform built on Go with enterprise-grade features including:
- Multi-agent AI coordination
- Parallel model execution  
- Advanced security hardening
- Performance optimization for resource-constrained environments
- Comprehensive monitoring and analytics

## 🚀 Quick Start

### Prerequisites
- **Go 1.23.10+** (latest security patches)
- **Git**
- **Docker** (optional, for containerized development)
- **PostgreSQL 17.5+** (for database features)

### Development Environment Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/methedrone/PDw-X.git
   cd PDw-X/app
   ```

2. **Install development tools**
   ```bash
   # Install Go linting tools
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
   go install golang.org/x/vuln/cmd/govulncheck@latest
   
   # Install build dependencies
   go mod download
   ```

3. **Run build validation**
   ```bash
   # Comprehensive validation (MacBook 2012 optimized)
   ./scripts/comprehensive-build-validation.sh
   
   # Quick build validation
   ./scripts/phase1-quality-gates.sh
   ```

4. **Start development environment**
   ```bash
   # Start local development
   ./start_local.sh
   
   # Start with hot reload
   ./scripts/dev.sh
   ```

## 🏗️ Project Architecture

### Module Structure
```
app/
├── cli/           # Command-line interface (Cobra-based)
├── server/        # HTTP API server (Gorilla Mux + PostgreSQL)
├── shared/        # Common data models and utilities
└── docs/          # Docusaurus documentation site
```

### Key Technologies
- **Language**: Go 1.23.10+ with security patches
- **Database**: PostgreSQL 17.5 with security hardening
- **CLI Framework**: Cobra (spf13/cobra)
- **TUI Components**: Bubble Tea (charmbracelet/bubbletea)
- **AI Integration**: Multiple providers (OpenAI, Anthropic, Google, OpenRouter)
- **Containerization**: Docker with multi-stage optimization

## 📋 Development Workflow

### 1. Issue Creation and Assignment
- Check existing issues before creating new ones
- Use issue templates for bug reports and feature requests
- Assign yourself to issues you plan to work on
- Add appropriate labels (bug, enhancement, security, performance, etc.)

### 2. Branch Strategy
```bash
# Create feature branch
git checkout -b feature/your-feature-name

# Create bugfix branch  
git checkout -b bugfix/issue-description

# Create security branch
git checkout -b security/vulnerability-fix
```

### 3. Development Process

#### Code Standards
- **Go Standards**: Follow [Effective Go](https://golang.org/doc/effective_go.html) practices
- **Error Handling**: Comprehensive error handling with structured responses
- **Security**: Input validation, secure logging, credential protection
- **Performance**: Memory efficiency, database optimization, caching strategies
- **Testing**: Aim for 80%+ test coverage with meaningful tests

#### Required Checks Before Commit
```bash
# 1. Run comprehensive validation
./scripts/comprehensive-build-validation.sh

# 2. Format code
gofmt -w .
go mod tidy

# 3. Run static analysis
golangci-lint run

# 4. Run security scan
govulncheck ./...

# 5. Run tests
go test ./...
```

### 4. Commit Guidelines

#### Commit Message Format
Follow [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `perf`: Performance improvements
- `test`: Adding or updating tests
- `build`: Build system changes
- `ci`: CI/CD changes
- `chore`: Maintenance tasks
- `security`: Security fixes or improvements

**Examples:**
```bash
feat(cli): add parallel model execution support

fix(server): resolve memory leak in connection pooling

security(auth): implement rate limiting for API endpoints

perf(database): optimize query performance with indexing

docs(api): update API reference with new endpoints
```

### 5. Pull Request Process

#### Before Submitting
- [ ] All tests pass locally
- [ ] Code follows project standards
- [ ] Documentation is updated
- [ ] Security implications considered
- [ ] Performance impact assessed

#### PR Title and Description
- Use descriptive titles following conventional commit format
- Include comprehensive description of changes
- Reference related issues using `Fixes #123` or `Closes #456`
- Add screenshots for UI changes
- Include performance benchmarks for performance changes

#### PR Template
```markdown
## Description
Brief description of changes made.

## Type of Change
- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Documentation update
- [ ] Performance improvement
- [ ] Security enhancement

## Testing
- [ ] Unit tests added/updated
- [ ] Integration tests added/updated
- [ ] Manual testing completed
- [ ] Performance testing completed

## Security Considerations
- [ ] Input validation implemented
- [ ] No sensitive data exposed
- [ ] Authentication/authorization properly handled
- [ ] Security scanning completed

## Performance Impact
- [ ] No performance regression
- [ ] Performance improvements documented
- [ ] Memory usage optimized
- [ ] Database queries optimized

## Documentation
- [ ] Code comments updated
- [ ] API documentation updated
- [ ] User documentation updated
- [ ] Changelog updated

## Checklist
- [ ] Code follows project style guidelines
- [ ] Self-review completed
- [ ] Breaking changes documented
- [ ] Tests added for new functionality
- [ ] All CI checks pass
```

## 🧪 Testing Standards

### Test Categories
1. **Unit Tests**: Test individual functions and methods
2. **Integration Tests**: Test component interactions
3. **End-to-End Tests**: Test complete workflows
4. **Performance Tests**: Benchmark critical paths
5. **Security Tests**: Vulnerability and penetration testing

### Test File Organization
```
module/
├── handler.go
├── handler_test.go
├── integration/
│   └── handler_integration_test.go
└── testdata/
    └── test_fixtures.json
```

### Test Requirements
- Tests must be deterministic and repeatable
- Use table-driven tests for multiple scenarios
- Mock external dependencies appropriately
- Include negative test cases
- Test error conditions thoroughly

### Running Tests
```bash
# Unit tests
go test ./...

# Integration tests
go test -tags=integration ./...

# Performance benchmarks
go test -bench=. -benchmem ./...

# Race condition detection
go test -race ./...

# Coverage analysis
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 🔒 Security Guidelines

### Security First Approach
- **Input Validation**: Validate all inputs at boundaries
- **Authentication**: Implement proper authentication mechanisms
- **Authorization**: Follow principle of least privilege
- **Encryption**: Use TLS for all communications
- **Logging**: Secure logging with sensitive data redaction
- **Dependencies**: Regular security updates and vulnerability scanning

### Security Checklist
- [ ] Input sanitization implemented
- [ ] SQL injection prevention
- [ ] XSS protection
- [ ] CSRF protection
- [ ] Rate limiting implemented
- [ ] Error messages don't leak sensitive information
- [ ] Secrets management properly implemented
- [ ] Audit logging for security events

### Vulnerability Reporting
For security vulnerabilities, please email [security@pdw-x.dev](mailto:security@pdw-x.dev) instead of creating a public issue.

## ⚡ Performance Guidelines

### Performance Requirements
- **Memory Efficiency**: Optimize for 4-8GB RAM environments (MacBook 2012 compatibility)
- **CPU Usage**: Conservative CPU utilization patterns
- **Database Performance**: Query optimization and connection pooling
- **Caching**: Implement intelligent caching strategies
- **Monitoring**: Performance metrics and alerting

### Performance Optimization Checklist
- [ ] Memory allocations minimized
- [ ] Object pooling implemented where appropriate
- [ ] Database queries optimized
- [ ] Caching implemented for frequently accessed data
- [ ] Goroutine leaks prevented
- [ ] Resource cleanup properly implemented

### Benchmarking Standards
```bash
# CPU profiling
go test -cpuprofile=cpu.prof -bench=.

# Memory profiling  
go test -memprofile=mem.prof -bench=.

# Trace analysis
go test -trace=trace.out -bench=.
```

## 🏗️ Build and Deployment

### Build Optimization
The project includes several build modes optimized for different scenarios:

```bash
# Development build with debugging
./build-local-cli.sh --mode debug

# Production build with optimization
./build-local-cli.sh --mode release

# Maximum optimization for older hardware
./build-local-cli.sh --mode optimized
```

### Docker Development
```bash
# Start development environment
docker-compose up -d

# Build optimized images
docker build --target=production -t pdw-x:latest .

# Run with resource constraints (MacBook 2012)
docker run --memory=1g --cpus=2 pdw-x:latest
```

## 📚 Documentation Standards

### Code Documentation
- Use clear, descriptive function and variable names
- Add godoc comments for public APIs
- Include examples in documentation when helpful
- Document complex algorithms and business logic

### User Documentation
- Update relevant documentation with feature changes
- Include code examples and usage patterns
- Maintain API reference accuracy
- Update troubleshooting guides when fixing common issues

### Documentation Structure
```
docs/
├── core-concepts/     # Fundamental concepts
├── api-reference/     # Complete API documentation
├── deployment/        # Deployment guides
├── security/          # Security best practices
└── troubleshooting/   # Common issues and solutions
```

## 🤝 Community Guidelines

### Code of Conduct
We follow the [Contributor Covenant Code of Conduct](https://www.contributor-covenant.org/version/2/1/code_of_conduct/). Please read and follow these guidelines to ensure a welcoming environment for all contributors.

### Communication Channels
- **GitHub Issues**: Bug reports and feature requests
- **GitHub Discussions**: General questions and community discussion
- **Discord**: Real-time community chat
- **Email**: security@pdw-x.dev for security issues

### Getting Help
- Check existing documentation first
- Search closed issues for similar problems
- Ask questions in GitHub Discussions
- Join our Discord community for real-time help

## 🎖️ Recognition

We appreciate all contributions to PDw-X! Contributors are recognized in:
- Release notes for significant contributions
- Contributors section in README
- Annual contributor acknowledgments
- Special recognition for security improvements

## 📈 Roadmap Alignment

Current development priorities:
1. **Phase 3**: Testing & DevOps infrastructure
2. **Phase 4**: Feature enhancement and IDE integration
3. **Phase 5**: Next-gen capabilities and enterprise features

See [MASTER_IMPLEMENTATION_COORDINATOR_GUIDE.md](MASTER_IMPLEMENTATION_COORDINATOR_GUIDE.md) for detailed roadmap.

---

## 🚀 Ready to Contribute?

1. **Fork the repository**
2. **Create a feature branch**
3. **Make your changes** following these guidelines
4. **Test thoroughly** using our validation scripts
5. **Submit a pull request** with a clear description

Thank you for contributing to PDw-X and helping build the future of AI-powered development tools!

---

*For questions about contributing, please reach out through GitHub Discussions or join our Discord community.*