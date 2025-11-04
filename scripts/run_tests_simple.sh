#!/usr/bin/env bash
# Game 服务单元测试脚本 - 适用于云效流水线
# 只运行 game 服务的单元测试，包含代码覆盖率检查

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "=========================================="
echo "Game 服务单元测试"
echo "=========================================="

# 获取项目根目录
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${PROJECT_ROOT}"

# 检查并设置 Go 环境
check_go_environment() {
    # 检查 Go 是否已安装
    if ! command -v go &> /dev/null; then
        echo -e "${RED}错误: 未找到 Go 命令${NC}"
        echo ""
        echo "请确保在流水线中已安装 Go 环境："
        echo "1. 在测试步骤之前添加 '安装go' 步骤"
        echo "2. 或者确保测试步骤继承了构建步骤的 Go 环境"
        echo ""
        echo "如果使用云效流水线，请检查："
        echo "- Go 安装步骤是否在测试步骤之前执行"
        echo "- 测试步骤是否与构建步骤在同一运行环境中"
        echo ""
        
        # 尝试查找常见的 Go 安装路径
        COMMON_GO_PATHS=(
            "/usr/local/go/bin"
            "/opt/go/bin"
            "$HOME/go/bin"
            "/usr/bin"
        )
        
        echo "尝试查找 Go 安装路径..."
        for path in "${COMMON_GO_PATHS[@]}"; do
            if [ -f "${path}/go" ]; then
                echo -e "${YELLOW}找到 Go: ${path}/go${NC}"
                echo "尝试添加到 PATH..."
                export PATH="${path}:${PATH}"
                if command -v go &> /dev/null; then
                    echo -e "${GREEN}✓ Go 环境已设置${NC}"
                    return 0
                fi
            fi
        done
        
        exit 1
    fi
    
    # 显示 Go 版本信息
    GO_VERSION=$(go version 2>/dev/null || echo "unknown")
    echo "Go 环境检查: ${GO_VERSION}"
    
    # 检查 GOPATH 和 GOROOT（如果有）
    if [ -n "$GOROOT" ]; then
        echo "GOROOT: ${GOROOT}"
    fi
    if [ -n "$GOPATH" ]; then
        echo "GOPATH: ${GOPATH}"
    fi
    
    echo ""
}

# 覆盖率阈值（百分比），可通过环境变量 COVERAGE_THRESHOLD 设置
COVERAGE_THRESHOLD="${COVERAGE_THRESHOLD:-0}"

# 检查 Go 环境
check_go_environment

# Game 服务路径
SERVICE="game"
SERVICE_PATH="${PROJECT_ROOT}/${SERVICE}"

# 检查服务目录是否存在
if [ ! -d "${SERVICE_PATH}" ]; then
    echo -e "${RED}错误: 服务目录 ${SERVICE} 不存在${NC}"
    exit 1
fi

echo ""
echo "----------------------------------------"
echo "测试服务: ${SERVICE}"
echo "----------------------------------------"

cd "${SERVICE_PATH}"

# 运行单元测试，排除集成测试文件
echo "运行单元测试..."
echo ""

# 运行测试并生成覆盖率报告
if go test -v -short -race -coverprofile=coverage.out \
    $(go list ./... | grep -v integration) \
    -run "^Test.*" -coverpkg=./... \
    $(find . -name "*_test.go" ! -name "*_it_test.go" | xargs dirname | sort -u | sed 's|^\.|./|' | tr '\n' ' ' 2>/dev/null || echo "./..."); then
    
    echo ""
    echo "----------------------------------------"
    echo "单元测试通过"
    echo "----------------------------------------"
    
    # 显示覆盖率报告
    if [ -f coverage.out ]; then
        echo ""
        echo "代码覆盖率报告:"
        echo "----------------------------------------"
        go tool cover -func=coverage.out
        
        echo ""
        echo "----------------------------------------"
        
        # 提取总覆盖率
        TOTAL_COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
        
        echo "总覆盖率: ${TOTAL_COVERAGE}%"
        echo "覆盖率阈值: ${COVERAGE_THRESHOLD}%"
        
        # 检查覆盖率阈值
        if [ -n "$COVERAGE_THRESHOLD" ] && [ "$COVERAGE_THRESHOLD" != "0" ]; then
            # 使用 awk 进行浮点数比较
            COVERAGE_CHECK=$(echo "$TOTAL_COVERAGE $COVERAGE_THRESHOLD" | awk '{if ($1 >= $2) print "PASS"; else print "FAIL"}')
            
            if [ "$COVERAGE_CHECK" = "FAIL" ]; then
                echo -e "${RED}错误: 代码覆盖率 ${TOTAL_COVERAGE}% 低于阈值 ${COVERAGE_THRESHOLD}%${NC}"
                exit 1
            else
                echo -e "${GREEN}✓ 代码覆盖率检查通过${NC}"
            fi
        fi
        
        # 生成 HTML 覆盖率报告（可选）
        if command -v go tool cover >/dev/null 2>&1; then
            go tool cover -html=coverage.out -o coverage.html
            echo ""
            echo "HTML 覆盖率报告已生成: ${SERVICE_PATH}/coverage.html"
        fi
    else
        echo -e "${YELLOW}警告: 未生成覆盖率报告文件${NC}"
    fi
else
    echo ""
    echo -e "${RED}错误: 单元测试失败${NC}"
    exit 1
fi

cd "${PROJECT_ROOT}"

echo ""
echo "=========================================="
echo "测试完成"
echo "=========================================="

