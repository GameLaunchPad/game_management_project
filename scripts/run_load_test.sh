#!/usr/bin/env bash
# JMeter 压力测试运行脚本

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "=========================================="
echo "Game Service Load Test"
echo "=========================================="

# 获取项目根目录
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${PROJECT_ROOT}"

# JMeter 配置
JMETER_HOME="${JMETER_HOME:-}"
JMETER_TEST_PLAN="${PROJECT_ROOT}/jmeter/GameService_LoadTest.jmx"
JMETER_RESULTS="${PROJECT_ROOT}/jmeter/results.jtl"
JMETER_HTML_REPORT="${PROJECT_ROOT}/jmeter/html-report"
JMETER_LOG="${PROJECT_ROOT}/jmeter/jmeter.log"

# 测试参数（可通过环境变量覆盖）
# 性能目标：支持500个并发用户
THREADS="${THREADS:-500}"
RAMP_UP="${RAMP_UP:-60}"
DURATION="${DURATION:-120}"
SERVER_HOST="${SERVER_HOST:-localhost}"
SERVER_PORT="${SERVER_PORT:-8881}"

# 检查 JMeter 是否安装
check_jmeter() {
    if [ -n "$JMETER_HOME" ] && [ -f "$JMETER_HOME/bin/jmeter" ]; then
        JMETER_CMD="$JMETER_HOME/bin/jmeter"
        echo -e "${GREEN}✓ 使用 JMETER_HOME: ${JMETER_HOME}${NC}"
        return 0
    fi
    
    if command -v jmeter >/dev/null 2>&1; then
        JMETER_CMD="jmeter"
        echo -e "${GREEN}✓ 使用系统 JMeter${NC}"
        return 0
    fi
    
    echo -e "${RED}错误: JMeter 未安装或未在 PATH 中${NC}"
    echo "请安装 JMeter 或设置 JMETER_HOME 环境变量"
    echo ""
    echo "安装方法:"
    echo "1. 从官网下载: https://jmeter.apache.org/download_jmeter.cgi"
    echo "2. 设置环境变量: export JMETER_HOME=/path/to/apache-jmeter"
    echo "3. 或使用包管理器安装"
    exit 1
}

# 检查测试计划文件
check_test_plan() {
    if [ ! -f "$JMETER_TEST_PLAN" ]; then
        echo -e "${RED}错误: 测试计划文件不存在: ${JMETER_TEST_PLAN}${NC}"
        exit 1
    fi
    echo -e "${GREEN}✓ 测试计划文件存在${NC}"
}

# 检查服务是否运行
check_service() {
    echo ""
    echo "检查服务状态..."
    
    # 检查端口是否开放
    if command -v nc >/dev/null 2>&1; then
        if nc -z "$SERVER_HOST" "$SERVER_PORT" 2>/dev/null; then
            echo -e "${GREEN}✓ 服务运行在 ${SERVER_HOST}:${SERVER_PORT}${NC}"
        else
            echo -e "${YELLOW}警告: 无法连接到 ${SERVER_HOST}:${SERVER_PORT}${NC}"
            echo "请确保 game_platform_api 服务已启动"
            echo ""
            echo "启动服务:"
            echo "  cd game_platform_api"
            echo "  go run main.go"
            read -p "是否继续测试? (y/n) " -n 1 -r
            echo
            if [[ ! $REPLY =~ ^[Yy]$ ]]; then
                exit 1
            fi
        fi
    else
        echo -e "${YELLOW}警告: 无法检查服务状态（nc 命令不可用）${NC}"
        echo "请确保服务运行在 ${SERVER_HOST}:${SERVER_PORT}"
    fi
}

# 创建输出目录
create_output_dirs() {
    mkdir -p "$(dirname "$JMETER_RESULTS")"
    mkdir -p "$JMETER_HTML_REPORT"
    echo -e "${GREEN}✓ 输出目录已创建${NC}"
}

# 运行测试
run_test() {
    echo ""
    echo "=========================================="
    echo "开始压力测试"
    echo "=========================================="
    echo "测试参数:"
    echo "  - 服务器: ${SERVER_HOST}:${SERVER_PORT}"
    echo "  - 线程数: ${THREADS}"
    echo "  - 启动时间: ${RAMP_UP}秒"
    echo "  - 持续时间: ${DURATION}秒"
    echo ""
    
    # 运行 JMeter（无GUI模式）
    # 注意：-J 参数设置的是 JMeter 属性，在测试计划中通过 ${__P(property,default)} 访问
    # 显示实际使用的参数值（用于调试）
    echo "JMeter 属性:"
    echo "  - server_host: ${SERVER_HOST}"
    echo "  - server_port: ${SERVER_PORT}"
    echo "  - threads: ${THREADS}"
    echo "  - ramp_up: ${RAMP_UP}"
    echo "  - duration: ${DURATION}"
    echo ""
    
    "${JMETER_CMD}" -n -t "${JMETER_TEST_PLAN}" \
        -l "${JMETER_RESULTS}" \
        -j "${JMETER_LOG}" \
        -Jserver_host="${SERVER_HOST}" \
        -Jserver_port="${SERVER_PORT}" \
        -Jthreads="${THREADS}" \
        -Jramp_up="${RAMP_UP}" \
        -Jduration="${DURATION}" \
        -e -o "${JMETER_HTML_REPORT}" 2>&1 | tee "${JMETER_LOG}.output"
    
    # 检查是否有测试数据
    if [ -f "${JMETER_RESULTS}" ]; then
        # 检查结果文件是否有数据（除了表头）
        RESULT_LINES=$(wc -l < "${JMETER_RESULTS}" 2>/dev/null || echo "0")
        if [ "$RESULT_LINES" -le "1" ]; then
            echo -e "${YELLOW}警告: 测试没有产生任何数据${NC}"
            echo "可能的原因："
            echo "  1. 线程数设置为 0"
            echo "  2. 测试持续时间太短"
            echo "  3. 服务未启动或无法连接"
            echo ""
            echo "建议："
            echo "  1. 检查服务是否运行: netstat -an | grep ${SERVER_PORT}"
            echo "  2. 增加线程数: THREADS=10 bash scripts/run_load_test.sh"
            echo "  3. 增加测试持续时间: DURATION=120 bash scripts/run_load_test.sh"
            echo ""
            echo "如果测试确实没有运行，HTML报告将无法生成"
        fi
    fi
    
    TEST_EXIT_CODE=$?
    
    # 检查测试结果
    if [ $TEST_EXIT_CODE -eq 0 ]; then
        echo ""
        echo -e "${GREEN}✓ 测试执行完成${NC}"
        
        # 检查是否有测试数据
        if [ -f "${JMETER_RESULTS}" ]; then
            RESULT_LINES=$(wc -l < "${JMETER_RESULTS}" 2>/dev/null || echo "0")
            if [ "$RESULT_LINES" -gt "1" ]; then
                # 有测试数据，检查HTML报告
                if [ -f "${JMETER_HTML_REPORT}/index.html" ]; then
                    echo -e "${GREEN}✓ HTML报告已生成${NC}"
                else
                    echo -e "${YELLOW}警告: HTML报告未生成，尝试手动生成...${NC}"
                    # 尝试手动生成HTML报告
                    if "${JMETER_CMD}" -g "${JMETER_RESULTS}" -o "${JMETER_HTML_REPORT}" 2>/dev/null; then
                        echo -e "${GREEN}✓ HTML报告已手动生成${NC}"
                    else
                        echo -e "${YELLOW}警告: 无法生成HTML报告，但测试数据已保存${NC}"
                    fi
                fi
            else
                echo -e "${YELLOW}警告: 测试没有产生任何数据${NC}"
                echo "HTML报告无法生成（需要测试数据）"
            fi
        fi
        
        echo ""
        echo "结果文件:"
        echo "  - 结果数据: ${JMETER_RESULTS}"
        if [ -f "${JMETER_HTML_REPORT}/index.html" ]; then
            echo "  - HTML报告: ${JMETER_HTML_REPORT}/index.html"
            echo ""
            echo "查看HTML报告:"
            echo "  open ${JMETER_HTML_REPORT}/index.html"
        else
            echo "  - HTML报告: 未生成（需要测试数据）"
        fi
        echo "  - 日志文件: ${JMETER_LOG}"
    else
        echo -e "${RED}错误: 测试执行失败${NC}"
        echo "查看日志: ${JMETER_LOG}"
        
        # 即使测试失败，如果结果文件存在，也尝试生成HTML报告
        if [ -f "${JMETER_RESULTS}" ]; then
            RESULT_LINES=$(wc -l < "${JMETER_RESULTS}" 2>/dev/null || echo "0")
            if [ "$RESULT_LINES" -gt "1" ]; then
                echo ""
                echo "尝试从已有结果生成HTML报告..."
                if "${JMETER_CMD}" -g "${JMETER_RESULTS}" -o "${JMETER_HTML_REPORT}" 2>/dev/null; then
                    echo -e "${GREEN}✓ HTML报告已生成${NC}"
                    echo "  - HTML报告: ${JMETER_HTML_REPORT}/index.html"
                fi
            fi
        fi
        
        exit 1
    fi
}

# 主函数
main() {
    check_jmeter
    check_test_plan
    check_service
    create_output_dirs
    run_test
}

# 运行主函数
main

