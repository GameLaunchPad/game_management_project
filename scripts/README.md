# 测试脚本说明

本目录包含用于自动化测试的脚本，适用于云效流水线。

## 脚本列表

### 1. `run_tests_simple.sh` (推荐用于流水线)

**Game 服务单元测试脚本**，专门为云效流水线优化。

**功能**:
- ✅ 只运行 game 服务的单元测试
- ✅ 排除集成测试（`*_it_test.go`）
- ✅ 自动生成代码覆盖率报告
- ✅ 支持覆盖率阈值检查
- ✅ 生成 HTML 覆盖率报告

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
2. **覆盖率报告**: 显示每个包和函数的覆盖率
3. **总覆盖率**: 显示整体覆盖率百分比
4. **阈值检查**: 如果设置了阈值，会检查覆盖率是否达到要求
5. **HTML 报告**: 自动生成 `game/coverage.html` 可视化报告

## 注意事项

1. 确保脚本有执行权限（脚本中已包含 `chmod +x`）
2. 脚本只测试 `game` 服务，不测试 `cp_center` 或 `game_platform_api`
3. 只运行单元测试，不运行集成测试（`*_it_test.go` 文件会被排除）
4. 覆盖率阈值设置为 0 时，不进行阈值检查
5. 覆盖率报告文件会生成在 `game/coverage.out` 和 `game/coverage.html`

