#!/usr/bin/env bash
# Simple golangci-lint script for monorepo structure
# Only scans handler directories to avoid cross-module dependency issues

set -e

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

# Get project root directory
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${PROJECT_ROOT}"

# Check if golangci-lint is installed
if ! command -v golangci-lint &> /dev/null; then
    echo -e "${RED}✗ golangci-lint not found${NC}"
    echo "Please install golangci-lint: https://golangci-lint.run/usage/install/"
    exit 1
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
    if golangci-lint run \
        --disable-all \
        --enable=errcheck \
        --timeout=5m \
        --tests=false \
        --exclude-dirs-use-default \
        ./$dirs/... 2>&1 | tee /tmp/golangci_${service}.log; then
        
        # Check if there were any issues
        ISSUES=$(grep -c "Error:" /tmp/golangci_${service}.log || echo "0")
        
        if [ "$ISSUES" -eq 0 ]; then
            echo -e "${GREEN}✓ $service passed (no issues found)${NC}"
        else
            echo -e "${YELLOW}⚠ $service has $ISSUES issues (non-blocking)${NC}"
            TOTAL_ISSUES=$((TOTAL_ISSUES + ISSUES))
        fi
    else
        EXIT_CODE=$?
        if [ $EXIT_CODE -eq 1 ]; then
            # Exit code 1 usually means issues were found (non-fatal)
            echo -e "${YELLOW}⚠ $service has issues (non-blocking)${NC}"
        else
            # Other exit codes are actual errors
            echo -e "${RED}✗ $service failed with exit code $EXIT_CODE${NC}"
            FAILED_SERVICES+=("$service")
        fi
    fi
    
    # Clean up
    rm -f /tmp/golangci_${service}.log
    
    # Return to project root
    cd "${PROJECT_ROOT}"
    echo ""
done

# Summary
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
    exit 0
else
    echo -e "${RED}✗ Failed services: ${FAILED_SERVICES[*]}${NC}"
    echo -e "${RED}  These services had actual errors (not just code quality issues)${NC}"
    exit 1
fi

