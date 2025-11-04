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

# 只测试 handler 目录
HANDLER_DIR="./handler"

# 检查 handler 目录是否存在
if [ ! -d "${HANDLER_DIR}" ]; then
    echo -e "${RED}错误: handler 目录不存在${NC}"
    exit 1
fi

echo "测试范围: ${HANDLER_DIR}"
echo "运行单元测试..."
echo ""

# JSONL 测试报告文件路径
JSONL_REPORT="${SERVICE_PATH}/test_report.jsonl"

# 运行 handler 目录下的单元测试（排除集成测试 *_it_test.go）
# -coverpkg=./handler 只统计 handler 目录下代码的覆盖率（不包括测试文件）
# -json 参数会输出 JSONL 格式的测试报告（每行一个 JSON 对象）
# 注意：go test 会自动排除 *_test.go 文件在覆盖率统计中
# 注意：-json 会将测试输出重定向到 JSONL 格式，但覆盖率信息会单独输出到文件

# 运行测试并生成 JSONL 报告
# -json 参数会输出 JSONL 格式（每行一个 JSON 对象）
# 同时我们需要在控制台显示一些基本信息
if go test -short -race -coverprofile=coverage.out \
    -coverpkg=./handler \
    -run "^Test.*" \
    -timeout 5m \
    -json \
    ./handler > "${JSONL_REPORT}" 2>&1; then
    
    # 检查测试是否真的成功（通过 JSONL 文件中的状态判断）
    TEST_EXIT_CODE=0
    if grep -q '"Action":"fail"' "${JSONL_REPORT}"; then
        TEST_EXIT_CODE=1
    fi
    
    # 显示测试摘要（从 JSONL 中提取关键信息）
    echo ""
    echo "测试摘要:"
    echo "----------------------------------------"
    
    # 从 JSONL 中提取测试结果统计
    # go test -json 输出格式：每行一个 JSON 对象，包含 Action、Test、Package 等字段
    # 统计每个 Action 的数量
    PASSED=$(grep -c '"Action":"pass"' "${JSONL_REPORT}" 2>/dev/null || echo "0")
    FAILED=$(grep -c '"Action":"fail"' "${JSONL_REPORT}" 2>/dev/null || echo "0")
    SKIP=$(grep -c '"Action":"skip"' "${JSONL_REPORT}" 2>/dev/null || echo "0")
    
    # 提取测试包和测试函数信息
    echo "通过的测试: ${PASSED}"
    if [ "$FAILED" != "0" ]; then
        echo -e "${RED}失败的测试: ${FAILED}${NC}"
        # 显示失败的测试详情
        echo ""
        echo "失败的测试详情:"
        grep '"Action":"fail"' "${JSONL_REPORT}" | grep -o '"Test":"[^"]*"' | sed 's/"Test":"//;s/"$//' | while read -r test_name; do
            echo -e "  ${RED}✗ ${test_name}${NC}"
        done
    else
        echo "失败的测试: ${FAILED}"
    fi
    if [ "$SKIP" != "0" ]; then
        echo "跳过的测试: ${SKIP}"
    fi
    
    echo ""
    echo -e "${GREEN}✓ JSONL 测试报告已生成: ${JSONL_REPORT}${NC}"
    
    # 如果测试失败，退出
    if [ "$TEST_EXIT_CODE" != "0" ]; then
        echo ""
        echo -e "${RED}错误: 单元测试失败${NC}"
        exit 1
    fi
    
    echo ""
    echo "----------------------------------------"
    echo "单元测试通过"
    echo "----------------------------------------"
    
    # 显示覆盖率报告（只显示 handler 目录的覆盖率）
    if [ -f coverage.out ]; then
        echo ""
        echo "代码覆盖率报告 (仅 handler 目录):"
        echo "----------------------------------------"
        
        # 使用 grep 过滤，只显示 handler 目录的覆盖率
        go tool cover -func=coverage.out | grep -E "(handler/|total)" | grep -v "_test.go"
        
        echo ""
        echo "----------------------------------------"
        
        # 提取 handler 目录的总覆盖率
        # 由于使用了 -coverpkg=./handler，total 行应该只包含 handler 包的覆盖率
        # 但为了更准确，我们过滤掉测试文件相关的行
        HANDLER_COVERAGE=$(go tool cover -func=coverage.out | grep "^total:" | awk '{print $3}' | sed 's/%//')
        
        # 如果无法从 total 行提取，尝试计算 handler 目录下所有非测试文件的覆盖率平均值
        if [ -z "$HANDLER_COVERAGE" ] || [ "$HANDLER_COVERAGE" = "0" ]; then
            # 计算 handler 目录下所有非测试文件的覆盖率平均值
            HANDLER_COVERAGE=$(go tool cover -func=coverage.out | grep "handler/" | grep -v "_test.go" | awk '{
                match($3, /([0-9]+\.[0-9]+)%/, arr);
                if (arr[1] != "") {
                    sum += arr[1];
                    count++;
                }
            } END {
                if (count > 0) {
                    printf "%.2f", sum / count;
                } else {
                    print "0";
                }
            }')
        fi
        
        if [ -n "$HANDLER_COVERAGE" ] && [ "$HANDLER_COVERAGE" != "0" ]; then
            echo "Handler 目录覆盖率: ${HANDLER_COVERAGE}%"
            echo "覆盖率阈值: ${COVERAGE_THRESHOLD}%"
            
            # 检查覆盖率阈值
            if [ -n "$COVERAGE_THRESHOLD" ] && [ "$COVERAGE_THRESHOLD" != "0" ]; then
                # 使用 awk 进行浮点数比较
                COVERAGE_CHECK=$(echo "$HANDLER_COVERAGE $COVERAGE_THRESHOLD" | awk '{if ($1 >= $2) print "PASS"; else print "FAIL"}')
                
                if [ "$COVERAGE_CHECK" = "FAIL" ]; then
                    echo -e "${RED}错误: Handler 目录代码覆盖率 ${HANDLER_COVERAGE}% 低于阈值 ${COVERAGE_THRESHOLD}%${NC}"
                    exit 1
                else
                    echo -e "${GREEN}✓ 代码覆盖率检查通过${NC}"
                fi
            fi
        else
            echo -e "${YELLOW}警告: 无法计算 handler 目录覆盖率${NC}"
            echo "显示完整覆盖率报告:"
            go tool cover -func=coverage.out | tail -5
        fi
        
        # 生成 HTML 覆盖率报告（可选）
        if command -v go tool cover >/dev/null 2>&1; then
            go tool cover -html=coverage.out -o handler_coverage.html
            echo ""
            echo "HTML 覆盖率报告已生成: ${SERVICE_PATH}/handler_coverage.html"
        fi
        
        # 生成云效可解析的覆盖率报告（JSON格式）
        generate_coverage_json() {
            local coverage_file="$1"
            local json_file="${SERVICE_PATH}/coverage_report.json"
            
            # 提取覆盖率数据
            local total_coverage="$HANDLER_COVERAGE"
            local coverage_data=$(go tool cover -func="$coverage_file" | grep "handler/" | grep -v "_test.go")
            
            # 计算各个函数的覆盖率
            local function_array=""
            local first=true
            while IFS= read -r line; do
                if [ -n "$line" ] && echo "$line" | grep -q "handler/"; then
                    local func_full=$(echo "$line" | awk '{print $1}')
                    local func_name=$(echo "$func_full" | sed 's/.*\///')
                    local func_coverage=$(echo "$line" | awk '{print $3}' | sed 's/%//')
                    
                    if [ -z "$func_coverage" ] || [ "$func_coverage" = "0" ]; then
                        continue
                    fi
                    
                    if [ "$first" = true ]; then
                        first=false
                    else
                        function_array="${function_array},"
                    fi
                    
                    # 转义函数名中的特殊字符
                    func_name=$(echo "$func_name" | sed 's/"/\\"/g')
                    function_array="${function_array}{\"function\":\"${func_name}\",\"coverage\":${func_coverage}}"
                fi
            done <<< "$coverage_data"
            
            # 生成JSON报告
            cat > "$json_file" <<EOF
{
  "total_coverage": ${total_coverage:-0},
  "coverage_threshold": ${COVERAGE_THRESHOLD:-0},
  "package": "game/handler",
  "timestamp": "$(date -u +"%Y-%m-%dT%H:%M:%SZ" 2>/dev/null || date +"%Y-%m-%dT%H:%M:%SZ")",
  "functions": [${function_array}]
}
EOF
            
            echo "覆盖率JSON报告已生成: ${json_file}"
        }
        
        # 将覆盖率信息添加到JSONL测试报告中（作为测试摘要）
        append_coverage_to_jsonl() {
            local jsonl_file="${JSONL_REPORT}"
            local total_coverage="$HANDLER_COVERAGE"
            
            if [ -f "$jsonl_file" ] && [ -n "$total_coverage" ] && [ "$total_coverage" != "0" ]; then
                # 在JSONL文件末尾添加覆盖率摘要信息
                cat >> "$jsonl_file" <<EOF
{"Time":"$(date -u +"%Y-%m-%dT%H:%M:%SZ" 2>/dev/null || date +"%Y-%m-%dT%H:%M:%SZ")","Action":"coverage","Package":"github.com/GameLaunchPad/game_management_project/game/handler","Coverage":${total_coverage},"Threshold":${COVERAGE_THRESHOLD:-0}}
EOF
                echo "覆盖率信息已添加到JSONL测试报告"
            fi
        }
        
        # 生成覆盖率JSON报告
        if [ -f coverage.out ] && [ -n "$HANDLER_COVERAGE" ] && [ "$HANDLER_COVERAGE" != "0" ]; then
            generate_coverage_json coverage.out
            # 同时将覆盖率信息添加到JSONL报告中
            append_coverage_to_jsonl
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

