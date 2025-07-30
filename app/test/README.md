# Plandex Testing Framework

This directory contains the comprehensive testing framework for Plandex, designed to achieve 80%+ test coverage across all modules with optimized performance for MacBook 2012 environments.

## Directory Structure

```
test/
├── unit/          # Unit tests for individual functions and methods
├── integration/   # Integration tests for component interactions
├── e2e/          # End-to-end tests for complete workflows
├── performance/  # Performance benchmarks and regression tests
├── security/     # Security testing and vulnerability assessments
├── utils/        # Test utilities and helpers
├── fixtures/     # Test data and mock responses
├── config/       # Test configuration files
└── scripts/      # Test execution and reporting scripts
```

## Test Categories

### Unit Tests (test/unit/)
- Individual function testing
- Isolated component behavior
- Mock dependencies
- Fast execution (<100ms per test)

### Integration Tests (test/integration/)
- Database interactions
- API endpoint testing
- Inter-module communication
- Moderate execution time (100ms-1s per test)

### End-to-End Tests (test/e2e/)
- Complete CLI workflows
- Full application scenarios
- Real environment testing
- Slower execution (1s+ per test)

### Performance Tests (test/performance/)
- Benchmarking critical paths
- Memory usage profiling
- Concurrent load testing
- MacBook 2012 optimization validation

### Security Tests (test/security/)
- Input validation testing
- Authentication/authorization
- Vulnerability scanning
- Security policy enforcement

## Test Execution

### Local Development
```bash
# Run all tests
make test

# Run specific test categories
make test-unit
make test-integration
make test-e2e
make test-performance
make test-security

# Run with coverage
make test-coverage

# MacBook 2012 optimized execution
make test-macbook2012
```

### CI/CD Pipeline
```bash
# Quality gates
make test-quality-gates

# Coverage reporting
make test-coverage-report

# Performance regression detection
make test-performance-regression
```

## MacBook 2012 Optimizations

### Resource Constraints
- Memory: 4-8GB RAM limitation
- CPU: 2.5-2.9GHz dual/quad-core
- Parallel test execution: Max 2-4 processes
- Test data size limits: <100MB per test suite

### Environment Variables
```bash
export PLANDEX_TEST_MEMORY_LIMIT=512MB
export PLANDEX_TEST_CPU_LIMIT=2
export PLANDEX_TEST_PARALLEL_JOBS=2
export PLANDEX_TEST_TIMEOUT=30s
```

## Quality Gates

### Coverage Requirements
- Overall coverage: 80%+
- Critical modules: 90%+
- New code: 95%+

### Performance Benchmarks
- Unit tests: <100ms average
- Integration tests: <1s average
- E2E tests: <30s average
- Memory usage: <512MB peak

### Security Standards
- Zero critical vulnerabilities
- Input validation: 100% coverage
- Authentication: Complete test coverage
- Authorization: Role-based testing

## Test Data Management

### Fixtures
- Standardized test data sets
- Mock API responses
- Database seed data
- File system test structures

### Test Isolation
- Clean state between tests
- Containerized test environments
- Database transaction rollbacks
- Temporary file cleanup

## Reporting and Metrics

### Coverage Reports
- HTML coverage reports
- Module-level breakdowns
- Trend analysis
- Quality gate status

### Performance Metrics
- Execution time tracking
- Memory usage profiling
- Benchmark comparisons
- Regression detection

### Security Assessment
- Vulnerability scan results
- Security policy compliance
- Authentication/authorization coverage
- Input validation verification