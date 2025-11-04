# 云效流水线自动化测试快速开始

## 3分钟快速配置

### 第一步：添加测试步骤

在云效流水线中，在**构建**和**部署**之间添加一个新的**测试阶段**。

### 第二步：配置单元测试

1. 在测试阶段中添加一个**执行命令**步骤
2. 步骤名称：`Game服务单元测试`
3. 工作目录：`${WORKSPACE}` 或留空（使用项目根目录）
4. **执行命令**：
   ```bash
   chmod +x scripts/run_tests_simple.sh
   bash scripts/run_tests_simple.sh
   ```
5. **环境变量**（可选，设置覆盖率阈值）：
   ```
   COVERAGE_THRESHOLD=50
   ```
   > 说明：如果设置了 `COVERAGE_THRESHOLD`，测试会检查覆盖率是否达到阈值，未达到会失败。
6. 失败策略：**失败时停止**（测试失败时阻止部署）

### 第三步：保存并运行

保存流水线配置，然后触发一次构建，验证测试步骤是否正常工作。

## 功能说明

### 当前测试配置

- ✅ **只测试 game 服务**（不测试 cp_center 和 game_platform_api）
- ✅ **只运行单元测试**（不运行集成测试）
- ✅ **自动生成代码覆盖率报告**
- ✅ **支持覆盖率阈值检查**（可选）
- ✅ **生成 HTML 覆盖率报告**

## 环境变量配置

### 可选的环境变量

- `COVERAGE_THRESHOLD`: 覆盖率阈值（百分比）
  - 默认值：`0`（不检查）
  - 示例：`50` 表示要求覆盖率至少 50%
  - 如果覆盖率低于阈值，测试会失败

## 常见问题

### Q: 测试会运行哪些内容？

A: 
- ✅ 只运行 `game` 服务的单元测试
- ✅ 排除集成测试文件（`*_it_test.go`）
- ✅ 不测试 `cp_center` 和 `game_platform_api` 服务

### Q: 如何设置覆盖率要求？

A: 在流水线步骤的环境变量中设置 `COVERAGE_THRESHOLD`，例如：
```
COVERAGE_THRESHOLD=50
```
这表示要求代码覆盖率至少达到 50%，否则测试会失败。

### Q: 测试失败怎么办？

A: 
1. 查看流水线日志，找到失败的具体测试用例
2. 如果是覆盖率检查失败，查看覆盖率报告，了解哪些代码没有被测试覆盖
3. 检查环境变量是否正确配置
4. 可以在本地运行相同的命令进行调试：
   ```bash
   cd game
   go test -v -short -race -coverprofile=coverage.out ./...
   go tool cover -func=coverage.out
   ```

### Q: 如何查看详细的覆盖率报告？

A: 
1. 脚本会在控制台输出覆盖率报告
2. 脚本会自动生成 HTML 报告：`game/coverage.html`
3. 可以在流水线构建产物中下载查看

### Q: 可以同时测试多个服务吗？

A: 当前脚本只测试 `game` 服务。如果需要测试其他服务，可以使用 `run_tests.sh` 脚本。

## 下一步

- 查看详细配置文档：[CI_CD_PIPELINE_SETUP.md](./CI_CD_PIPELINE_SETUP.md)
- 查看脚本说明：[../scripts/README.md](../scripts/README.md)

