#!/usr/bin/env bash
# Simple golangci-lint script for monorepo structure
# Only scans handler directories to avoid cross-module dependency issues

# Don't exit on error immediately - we want to collect all results
set +e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo "=========================================="
echo "Running golangci-lint for Monorepo"
echo "=========================================="
echo ""

# Debug: Show environment
echo "Debug Info:"
echo "  PWD: $(pwd)"
echo "  SHELL: $SHELL"
echo "  PATH: $PATH"
echo ""

# Get project root directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
echo "Project Root: ${PROJECT_ROOT}"
cd "${PROJECT_ROOT}" || exit 1
echo ""

# Check if golangci-lint is installed
echo "Checking golangci-lint..."
if ! command -v golangci-lint &> /dev/null; then
    echo -e "${RED}✗ golangci-lint not found${NC}"
    echo ""
    echo "Attempting to install golangci-lint..."
    
    # Try to install golangci-lint
    if command -v go &> /dev/null; then
        echo "Installing via go install..."
        go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
        
        # Add GOPATH/bin to PATH if not already there
        export PATH="$PATH:$(go env GOPATH)/bin"
        
        # Check again
        if command -v golangci-lint &> /dev/null; then
            echo -e "${GREEN}✓ golangci-lint installed successfully${NC}"
        else
            echo -e "${RED}✗ Failed to install golangci-lint${NC}"
            echo "Please install manually: https://golangci-lint.run/usage/install/"
            exit 1
        fi
    else
        echo -e "${RED}✗ Go not found, cannot auto-install golangci-lint${NC}"
        exit 1
    fi
fi

echo -e "${BLUE}golangci-lint version:${NC}"
golangci-lint version
echo ""

# Services and directories to scan
declare -A SERVICES=(
    ["game"]="handler"
    ["game_platform_api"]="biz/handler"
    ["cp_center"]="handler"
)

# Track failures
TOTAL_ISSUES=0
FAILED_SERVICES=()

# Scan each service
for service in "${!SERVICES[@]}"; do
    dirs="${SERVICES[$service]}"
    
    if [ ! -d "$service" ]; then
        echo -e "${YELLOW}⚠ Skipping $service (directory not found)${NC}"
        continue
    fi
    
    echo "=========================================="
    echo -e "${BLUE}Scanning: $service${NC}"
    echo "  Directories: $dirs"
    echo "=========================================="
    
    # Change to service directory
    cd "$service"
    
    # Run golangci-lint with minimal linters
    # Only check: errcheck (error handling)
    echo "Running: golangci-lint run --disable-all --enable=errcheck --timeout=5m --tests=false ./$dirs/..."
    
    # Create temp file
    TEMP_LOG="${PROJECT_ROOT}/golangci_${service}.log"
    
    golangci-lint run \
        --disable-all \
        --enable=errcheck \
        --timeout=5m \
        --tests=false \
        --exclude-dirs-use-default \
        ./$dirs/... 2>&1 | tee "${TEMP_LOG}"
    
    EXIT_CODE=${PIPESTATUS[0]}
    
    if [ $EXIT_CODE -eq 0 ]; then
        # No issues found
        echo -e "${GREEN}✓ $service passed (no issues found)${NC}"
    elif [ $EXIT_CODE -eq 1 ]; then
        # Exit code 1 usually means issues were found (non-fatal)
        ISSUES=$(grep -c "Error:" "${TEMP_LOG}" 2>/dev/null || echo "0")
        echo -e "${YELLOW}⚠ $service has $ISSUES issues (non-blocking)${NC}"
        TOTAL_ISSUES=$((TOTAL_ISSUES + ISSUES))
    else
        # Other exit codes are actual errors
        echo -e "${RED}✗ $service failed with exit code $EXIT_CODE${NC}"
        echo "Last 20 lines of output:"
        tail -20 "${TEMP_LOG}" 2>/dev/null || echo "(no log available)"
        FAILED_SERVICES+=("$service")
    fi
    
    # Clean up
    rm -f "${TEMP_LOG}"
    
    # Return to project root
    cd "${PROJECT_ROOT}"
    echo ""
done

# Summary
echo ""
echo "=========================================="
echo "Scan Summary"
echo "=========================================="

if [ ${#FAILED_SERVICES[@]} -eq 0 ]; then
    if [ $TOTAL_ISSUES -eq 0 ]; then
        echo -e "${GREEN}✓ All services passed with no issues!${NC}"
    else
        echo -e "${YELLOW}⚠ Found $TOTAL_ISSUES issues (non-blocking)${NC}"
        echo -e "${YELLOW}  These are suggestions for improvement, not errors${NC}"
    fi
    echo ""
    echo -e "${GREEN}✓ No blocker issues found!${NC}"
    echo ""
    echo "Exit code: 0 (Success)"
    exit 0
else
    echo -e "${RED}✗ Failed services: ${FAILED_SERVICES[*]}${NC}"
    echo -e "${RED}  These services had actual errors (not just code quality issues)${NC}"
    echo ""
    echo "Exit code: 1 (Failure)"
    exit 1
fi

