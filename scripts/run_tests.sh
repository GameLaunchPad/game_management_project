#!/usr/bin/env bash
# 自动化测试脚本
# 用于在云效流水线中运行单元测试和集成测试

set -e  # 遇到错误立即退出

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 打印函数
print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 测试模式：unit（单元测试）、integration（集成测试）、all（全部）
TEST_MODE="${TEST_MODE:-all}"

# 测试覆盖率阈值（百分比）
COVERAGE_THRESHOLD="${COVERAGE_THRESHOLD:-0}"

# 测试结果统计
TOTAL_TESTS=0
FAILED_TESTS=0

# 运行单个服务的单元测试
run_unit_tests() {
    local service=$1
    local service_path=$2
    
    print_info "运行 ${service} 服务的单元测试..."
    
    cd "${service_path}" || exit 1
    
    # 运行单元测试（排除集成测试 *_it_test.go）
    if go test -v -short -race -coverprofile=coverage.out -covermode=atomic \
        $(go list ./... | grep -v "/integration") \
        -run "^Test.*" -coverpkg=./... \
        $(find . -name "*_test.go" ! -name "*_it_test.go" | sed 's|^\./||' | sed 's|/[^/]*$||' | sort -u | sed 's|^|./|' | tr '\n' ' '); then
        print_info "${service} 单元测试通过"
        
        # 生成覆盖率报告
        if [ -f coverage.out ]; then
            coverage=$(go tool cover -func=coverage.out | grep total | awk '{print $3}')
            print_info "${service} 单元测试覆盖率: ${coverage}"
            
            # 检查覆盖率阈值
            coverage_num=$(echo "$coverage" | sed 's/%//')
            if (( $(echo "$coverage_num < $COVERAGE_THRESHOLD" | bc -l) )); then
                print_warn "${service} 单元测试覆盖率 ${coverage} 低于阈值 ${COVERAGE_THRESHOLD}%"
            fi
        fi
    else
        print_error "${service} 单元测试失败"
        FAILED_TESTS=$((FAILED_TESTS + 1))
        return 1
    fi
    
    cd - > /dev/null || exit 1
}

# 运行单个服务的集成测试
run_integration_tests() {
    local service=$1
    local service_path=$2
    
    print_info "运行 ${service} 服务的集成测试..."
    
    # 检查数据库连接
    if [ -z "$TEST_DB_DSN" ]; then
        print_warn "TEST_DB_DSN 未设置，使用默认数据库配置"
        print_warn "如需使用自定义数据库，请设置环境变量 TEST_DB_DSN"
    fi
    
    cd "${service_path}" || exit 1
    
    # 运行集成测试（只运行 *_it_test.go）
    if go test -v -timeout 10m \
        -run "^Test.*Suite" \
        $(find . -name "*_it_test.go" | sed 's|^\./||' | sed 's|/[^/]*$||' | sort -u | sed 's|^|./|' | tr '\n' ' '); then
        print_info "${service} 集成测试通过"
    else
        print_error "${service} 集成测试失败"
        FAILED_TESTS=$((FAILED_TESTS + 1))
        return 1
    fi
    
    cd - > /dev/null || exit 1
}

# 主函数
main() {
    print_info "开始运行自动化测试..."
    print_info "测试模式: ${TEST_MODE}"
    
    # 获取项目根目录
    PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
    cd "${PROJECT_ROOT}" || exit 1
    
    # 定义服务列表
    SERVICES=("game" "cp_center")
    
    # 根据测试模式运行测试
    for service in "${SERVICES[@]}"; do
        service_path="${PROJECT_ROOT}/${service}"
        
        if [ ! -d "${service_path}" ]; then
            print_warn "服务目录 ${service} 不存在，跳过"
            continue
        fi
        
        case "${TEST_MODE}" in
            unit)
                run_unit_tests "${service}" "${service_path}" || true
                ;;
            integration)
                run_integration_tests "${service}" "${service_path}" || true
                ;;
            all)
                run_unit_tests "${service}" "${service_path}" || true
                run_integration_tests "${service}" "${service_path}" || true
                ;;
            *)
                print_error "未知的测试模式: ${TEST_MODE}"
                print_error "支持的模式: unit, integration, all"
                exit 1
                ;;
        esac
    done
    
    # 输出测试总结
    echo ""
    print_info "================================"
    print_info "测试总结"
    print_info "================================"
    
    if [ $FAILED_TESTS -eq 0 ]; then
        print_info "所有测试通过！"
        exit 0
    else
        print_error "有 ${FAILED_TESTS} 个测试失败"
        exit 1
    fi
}

# 执行主函数
main

