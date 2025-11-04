# 测试脚本说明

本目录包含用于自动化测试的脚本，适用于云效流水线。

## 脚本列表

### 1. `run_tests_simple.sh` (推荐用于流水线)

**Game 服务单元测试脚本**，专门为云效流水线优化。

**功能**:
- ✅ 只运行 game 服务的 handler 目录下的单元测试
- ✅ 排除集成测试（`*_it_test.go`）
- ✅ 只统计 handler 目录下代码的覆盖率（不包括测试文件和其他目录）
- ✅ 自动生成代码覆盖率报告
- ✅ 支持覆盖率阈值检查
- ✅ 生成 HTML 覆盖率报告
- ✅ **生成 JSONL 格式的测试报告**（用于云效测试报告）

**使用方法**:
```bash
# 基本使用（默认阈值为 0，不检查）
bash scripts/run_tests_simple.sh

# 设置覆盖率阈值（例如 50%）
COVERAGE_THRESHOLD=50 bash scripts/run_tests_simple.sh
```

### 2. `run_tests.sh`

功能完整的测试脚本，包含详细的输出和覆盖率检查（支持多个服务）。

## 在云效流水线中使用

### 快速配置步骤

1. **添加测试步骤**到流水线（在构建和部署之间）

2. **步骤类型**: 执行命令

3. **工作目录**: `${WORKSPACE}` 或项目根目录

4. **执行命令**:
   ```bash
   chmod +x scripts/run_tests_simple.sh
   bash scripts/run_tests_simple.sh
   ```

5. **环境变量**（可选）:
   - `COVERAGE_THRESHOLD`: 覆盖率阈值（百分比），例如 `50` 表示 50%

### 示例配置

#### 单元测试步骤（无覆盖率要求）
```
步骤名称: Game服务单元测试
执行命令:
  chmod +x scripts/run_tests_simple.sh
  bash scripts/run_tests_simple.sh
失败策略: 失败时停止
```

#### 单元测试步骤（带覆盖率要求）
```
步骤名称: Game服务单元测试（覆盖率检查）
环境变量:
  COVERAGE_THRESHOLD: 50
执行命令:
  chmod +x scripts/run_tests_simple.sh
  bash scripts/run_tests_simple.sh
失败策略: 失败时停止
```

## 输出说明

脚本会输出：
1. **测试结果**: 显示每个测试用例的执行结果
2. **测试摘要**: 显示通过、失败、跳过的测试数量
3. **覆盖率报告**: 只显示 handler 目录下每个函数的覆盖率
4. **Handler 目录覆盖率**: 显示 handler 目录的整体覆盖率百分比
5. **阈值检查**: 如果设置了阈值，会检查覆盖率是否达到要求
6. **HTML 报告**: 自动生成 `game/handler_coverage.html` 可视化报告
7. **JSONL 报告**: 自动生成 `game/test_report.jsonl` 测试报告（用于云效测试报告）

## JSONL 测试报告

脚本会自动生成 JSONL 格式的测试报告文件：`game/test_report.jsonl`

**JSONL 格式说明**：
- JSONL（JSON Lines）是一种文本格式，每行包含一个 JSON 对象
- Go 的 `go test -json` 会输出每行一个 JSON 对象，包含测试事件信息
- 每个 JSON 对象包含字段：`Time`、`Action`、`Package`、`Test`、`Output` 等

**在云效中使用**：
1. 在云效流水线的测试步骤中，将 `game/test_report.jsonl` 指定为测试报告文件
2. 云效会自动解析 JSONL 文件，显示测试结果
3. 可以在云效的测试报告页面查看详细的测试结果

**示例 JSONL 格式**：
```json
{"Time":"2025-01-01T10:00:00Z","Action":"run","Package":"github.com/.../game/handler","Test":"TestCreateGameDetail_Success"}
{"Time":"2025-01-01T10:00:01Z","Action":"pass","Package":"github.com/.../game/handler","Test":"TestCreateGameDetail_Success"}
```

## 注意事项

1. 确保脚本有执行权限（脚本中已包含 `chmod +x`）
2. 脚本只测试 `game/handler` 目录，不测试其他目录（dao, dal, service 等）
3. 只运行单元测试，不运行集成测试（`*_it_test.go` 文件会被排除）
4. 覆盖率统计只包括 handler 目录下的非测试代码（`*.go` 文件，不包括 `*_test.go`）
5. 覆盖率阈值设置为 0 时，不进行阈值检查
6. 覆盖率报告文件会生成在 `game/coverage.out` 和 `game/handler_coverage.html`
7. JSONL 测试报告文件会生成在 `game/test_report.jsonl`，可用于云效测试报告

