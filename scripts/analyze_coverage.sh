#!/usr/bin/env bash
# 覆盖率分析工具
# 用于分析测试覆盖率，找出未覆盖的代码

set -e

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
COVERAGE_FILE="${PROJECT_ROOT}/game/coverage.out"

if [ ! -f "$COVERAGE_FILE" ]; then
    echo "错误: 覆盖率文件不存在: $COVERAGE_FILE"
    echo "请先运行测试脚本生成覆盖率文件"
    exit 1
fi

echo "=========================================="
echo "覆盖率分析报告"
echo "=========================================="
echo ""

echo "1. 未覆盖的函数（覆盖率为0%）："
echo "----------------------------------------"
go tool cover -func="$COVERAGE_FILE" | grep "handler/" | grep -E '\s+0\.0%' || echo "  无（所有函数都有一定覆盖率）"
echo ""

echo "2. 覆盖率较低的函数（<50%）："
echo "----------------------------------------"
go tool cover -func="$COVERAGE_FILE" | grep "handler/" | awk '$3 < 50 {print}' || echo "  无"
echo ""

echo "3. 覆盖率中等的函数（50%-80%）："
echo "----------------------------------------"
go tool cover -func="$COVERAGE_FILE" | grep "handler/" | awk '$3 >= 50 && $3 < 80 {print}' || echo "  无"
echo ""

echo "4. 覆盖率较高的函数（>=80%）："
echo "----------------------------------------"
go tool cover -func="$COVERAGE_FILE" | grep "handler/" | awk '$3 >= 80 {print}' | head -20
echo ""

echo "5. Handler 目录覆盖率："
echo "----------------------------------------"
# 计算handler目录的平均覆盖率（不直接使用total，因为可能包含其他包）
HANDLER_COV=$(go tool cover -func="$COVERAGE_FILE" 2>/dev/null | grep "handler/" | grep -v "_test.go" | awk '{
    if (match($3, /([0-9]+\.[0-9]+)%/, arr)) {
        sum += arr[1];
        count++;
    } else if (match($3, /([0-9]+)%/, arr)) {
        sum += arr[1];
        count++;
    }
} END {
    if (count > 0) {
        printf "%.2f%%", sum / count;
    } else {
        print "0%";
    }
}' || echo "无法计算")

echo "Handler 目录平均覆盖率: ${HANDLER_COV}"

echo ""
echo "6. 总覆盖率（所有包）："
echo "----------------------------------------"
go tool cover -func="$COVERAGE_FILE" 2>/dev/null | grep "^total:" || echo "  无法获取"
echo ""

echo "建议："
echo "- 查看 game/handler_coverage.html 了解详细的未覆盖代码行"
echo "- 为未覆盖的代码添加测试用例"
echo "- 重点关注边界情况和错误处理分支"

