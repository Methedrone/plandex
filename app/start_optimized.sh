#!/bin/bash
# Plandex Optimized Startup Script
# Automatically detects environment and applies appropriate performance settings

set -e

# Function to detect system resources
detect_system() {
    echo "🔍 Detecting system resources..."
    
    # Detect CPU cores
    if command -v nproc >/dev/null 2>&1; then
        CPU_CORES=$(nproc)
    elif command -v sysctl >/dev/null 2>&1; then
        CPU_CORES=$(sysctl -n hw.ncpu 2>/dev/null || echo "4")
    else
        CPU_CORES=4
    fi
    
    # Detect available memory (rough estimation)
    if command -v free >/dev/null 2>&1; then
        MEMORY_GB=$(free -g | awk '/^Mem:/{print $2}')
    elif command -v sysctl >/dev/null 2>&1; then
        MEMORY_BYTES=$(sysctl -n hw.memsize 2>/dev/null || echo "8589934592")
        MEMORY_GB=$((MEMORY_BYTES / 1024 / 1024 / 1024))
    else
        MEMORY_GB=8
    fi
    
    echo "  CPU Cores: $CPU_CORES"
    echo "  Memory: ${MEMORY_GB}GB"
}

# Function to apply MacBook 2012 optimizations
apply_macbook2012_settings() {
    echo "🍎 Applying MacBook 2012 optimizations..."
    
    export GOGC=50
    export GOMEMLIMIT=512MiB
    export GOMAXPROCS=4
    export GOENV=development
    export CONSTRAINED_ENV=true
    export PLANDEX_PERF_DEBUG=true
    
    echo "  GOGC: $GOGC (aggressive GC)"
    echo "  GOMEMLIMIT: $GOMEMLIMIT"
    echo "  GOMAXPROCS: $GOMAXPROCS"
}

# Function to apply production settings
apply_production_settings() {
    echo "🚀 Applying production optimizations..."
    
    export GOGC=100
    export GOMEMLIMIT=2GiB
    export GOENV=production
    export CONSTRAINED_ENV=false
    export PLANDEX_PERF_DEBUG=false
    
    echo "  GOGC: $GOGC (standard GC)"
    echo "  GOMEMLIMIT: $GOMEMLIMIT"
    echo "  GOMAXPROCS: all cores"
}

# Function to start the server
start_server() {
    echo "🚀 Starting Plandex server with optimized settings..."
    
    # Make sure we're in the right directory
    cd "$(dirname "$0")"
    
    # Check if server binary exists
    if [ ! -f "server/plandex-server" ]; then
        echo "📦 Building server..."
        cd server
        go build -o plandex-server .
        cd ..
    fi
    
    # Start the server
    if [ "$CONSTRAINED_ENV" = "true" ]; then
        echo "🔋 Starting in resource-efficient mode..."
    else
        echo "⚡ Starting in high-performance mode..."
    fi
    
    exec ./server/plandex-server
}

# Main execution
main() {
    echo "🎯 Plandex Performance Optimization Startup"
    echo "==========================================="
    
    # Parse command line arguments
    FORCE_MODE=""
    for arg in "$@"; do
        case $arg in
            --macbook2012)
                FORCE_MODE="macbook2012"
                shift
                ;;
            --production)
                FORCE_MODE="production"
                shift
                ;;
            --help|-h)
                echo "Usage: $0 [--macbook2012|--production]"
                echo ""
                echo "Options:"
                echo "  --macbook2012    Force MacBook 2012 optimizations"
                echo "  --production     Force production optimizations"
                echo "  --help, -h       Show this help message"
                echo ""
                echo "Without options, system resources are auto-detected."
                exit 0
                ;;
        esac
    done
    
    # Detect system resources
    detect_system
    
    # Apply appropriate settings
    if [ "$FORCE_MODE" = "macbook2012" ]; then
        apply_macbook2012_settings
    elif [ "$FORCE_MODE" = "production" ]; then
        apply_production_settings
    elif [ "$CPU_CORES" -le 4 ] && [ "$MEMORY_GB" -le 8 ]; then
        echo "🔍 Detected resource-constrained environment"
        apply_macbook2012_settings
    else
        echo "🔍 Detected high-resource environment"
        apply_production_settings
    fi
    
    echo ""
    
    # Start the server
    start_server
}

# Run main function with all arguments
main "$@"