# 本地运行测试脚本指南

## 快速开始

### 方法1: 直接运行脚本（推荐）

```bash
# 在项目根目录下运行
bash scripts/run_tests_simple.sh
```

### 方法2: 设置权限后运行

```bash
# 设置执行权限
chmod +x scripts/run_tests_simple.sh

# 运行脚本
./scripts/run_tests_simple.sh
```

### 方法3: 使用环境变量配置

```bash
# 设置覆盖率阈值（可选）
COVERAGE_THRESHOLD=80 bash scripts/run_tests_simple.sh

# 禁用race detection（如果CGO不可用）
ENABLE_RACE=false bash scripts/run_tests_simple.sh

# 启用CGO和race detection（如果需要）
CGO_ENABLED=1 bash scripts/run_tests_simple.sh
```

## 前置条件

1. **Go环境**: 确保已安装Go（1.20+）
   ```bash
   go version
   ```

2. **依赖**: 确保已下载依赖
   ```bash
   cd game
   go mod download
   ```

3. **工作目录**: 在项目根目录下运行脚本
   ```bash
   # 项目根目录结构
   game_management_project/
   ├── scripts/
   │   └── run_tests_simple.sh
   ├── game/
   │   ├── go.mod
   │   └── handler/
   └── ...
   ```

## 脚本执行流程

1. **检查Go环境**: 自动检测Go是否安装
2. **切换到game目录**: 自动切换到game服务目录
3. **验证模块**: 检查go.mod文件是否存在
4. **运行测试**: 执行handler目录下的单元测试
5. **生成报告**: 
   - 生成JSONL测试报告
   - 生成HTML覆盖率报告（包含总覆盖率）
   - 生成JSON覆盖率报告

## 输出文件

运行成功后，会在 `game/` 目录下生成以下文件：

1. **coverage.out**: 覆盖率数据文件
2. **handler_coverage.html**: HTML覆盖率报告（已包含总覆盖率信息）
3. **test_report.jsonl**: JSONL格式测试报告
4. **coverage_report.json**: JSON格式覆盖率报告

## 常见问题

### 问题1: `go: command not found`

**解决方案**:
- 确保Go已安装
- 检查PATH环境变量
- 在Windows上，确保Go安装路径在PATH中

### 问题2: `go.mod file not found`

**解决方案**:
- 确保在项目根目录下运行脚本
- 检查 `game/go.mod` 文件是否存在

### 问题3: 测试失败

**解决方案**:
- 查看测试输出中的错误信息
- 检查测试文件是否正确
- 确保依赖已下载：`cd game && go mod download`

### 问题4: 覆盖率报告为空

**解决方案**:
- 检查 `coverage.out` 文件是否存在
- 确保测试实际运行了
- 查看脚本输出中的警告信息

### 问题5: `-race requires cgo` 错误

**原因**: race detection 需要 CGO，但 CGO 不可用

**解决方案**:
- **方案1（推荐）**: 脚本会自动检测并禁用 race detection，无需手动操作
- **方案2**: 手动禁用 race detection
  ```bash
  ENABLE_RACE=false bash scripts/run_tests_simple.sh
  ```
- **方案3**: 启用 CGO（如果可能）
  ```bash
  CGO_ENABLED=1 bash scripts/run_tests_simple.sh
  ```

**注意**: race detection 是可选的，禁用它不会影响测试和覆盖率统计

## 查看结果

### 查看HTML覆盖率报告

```bash
# Linux/Mac
open game/handler_coverage.html

# Windows
start game/handler_coverage.html

# 或者直接在浏览器中打开
# file:///path/to/game_management_project/game/handler_coverage.html
```

### 查看测试报告

```bash
# 查看JSONL测试报告
cat game/test_report.jsonl | head -20

# 查看JSON覆盖率报告
cat game/coverage_report.json | python -m json.tool
```

### 查看覆盖率统计

```bash
cd game
go tool cover -func=coverage.out | grep "handler/" | grep -v "_test.go"
```

## 调试模式

如果需要查看更详细的输出，可以：

```bash
# 使用bash -x查看详细执行过程
bash -x scripts/run_tests_simple.sh

# 或者修改脚本，在开头添加
# set -x  # 显示每个命令
```

## 示例输出

成功运行后，你会看到类似以下的输出：

```
==========================================
Game 服务单元测试
==========================================

当前工作目录: /path/to/game_management_project/game
Go模块信息:
  - 模块路径: github.com/GameLaunchPad/game_management_project/game
验证handler包...
  - Handler包路径: github.com/GameLaunchPad/game_management_project/game/handler
测试范围: ./handler
覆盖率目标: 95%

运行单元测试...

✓ 使用 -coverpkg 参数成功

测试摘要:
----------------------------------------
通过的测试: 25
失败的测试: 0
跳过的测试: 0

✓ JSONL 测试报告已生成: /path/to/game/handler_coverage.html

单元测试通过
----------------------------------------

代码覆盖率报告 (仅 handler 目录):
----------------------------------------
github.com/GameLaunchPad/game_management_project/game/handler/create_game_detail.go:CreateGameDetail  100.0%
...

当前覆盖率: 95.50%
覆盖率目标: 95%
✓ 代码覆盖率检查通过 (超出目标 0.50%)

HTML 覆盖率报告已生成（已添加总覆盖率信息）: /path/to/game/handler_coverage.html
```

## 下一步

1. 查看HTML覆盖率报告，了解详细的代码覆盖情况
2. 如果覆盖率低于目标，参考 `docs/TEST_CASE_IMPROVEMENT_GUIDE.md` 添加测试用例
3. 在云效流水线中配置测试步骤，参考 `docs/QUICK_START_TESTING.md`

