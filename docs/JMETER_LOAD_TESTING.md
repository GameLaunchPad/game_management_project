# JMeter 压力测试指南

## 概述

本文档说明如何使用 JMeter 对 game 服务进行压力测试。game 服务通过 `game_platform_api` HTTP API 网关暴露，运行在端口 8881。

## API 端点

### Game 服务接口

| 方法 | 路径 | 描述 | 状态 |
|------|------|------|------|
| GET | `/api/v1/games` | 获取游戏列表 | ✅ 已测试 |
| GET | `/api/v1/games/:id` | 获取游戏详情 | ✅ 已测试 |
| POST | `/api/v1/games` | 创建游戏 | ✅ 已测试 |
| PUT | `/api/v1/games/:id` | 更新游戏详情 | ✅ 已测试 |
| POST | `/api/v1/games/review` | 审核游戏版本 | ✅ 已测试 |
| DELETE | `/api/v1/games/:id/draft` | 删除游戏草稿 | ✅ 已测试 |

**服务地址**: `http://localhost:8881` (默认)

## 前置条件

### 1. 安装 JMeter

**Windows**:
1. 下载 JMeter: https://jmeter.apache.org/download_jmeter.cgi
2. 解压到任意目录（如 `C:\apache-jmeter-5.6`）
3. 运行 `bin/jmeter.bat`

**Linux/Mac**:
```bash
# 使用包管理器安装
# Ubuntu/Debian
sudo apt-get install jmeter

# macOS (使用 Homebrew)
brew install jmeter

# 或从官网下载
wget https://dlcdn.apache.org//jmeter/binaries/apache-jmeter-5.6.tgz
tar -xzf apache-jmeter-5.6.tgz
```

### 2. 启动服务

确保以下服务已启动：
```bash
# 1. 启动 game 服务（Thrift RPC）
cd game
./output/bin/game

# 2. 启动 game_platform_api（HTTP API网关）
cd game_platform_api
go run main.go
# 或使用构建后的二进制文件
```

### 3. 准备测试数据

确保数据库中有测试数据，或使用测试脚本创建。

## JMeter 测试计划

### 快速开始

1. **打开 JMeter 测试计划**
   - 文件位置: `jmeter/GameService_LoadTest.jmx`
   - 在 JMeter 中打开: `File > Open > GameService_LoadTest.jmx`

2. **配置测试参数**
   - 在 `Test Plan` 中设置线程数、循环次数等
   - 或使用 `User Defined Variables` 配置

3. **运行测试**
   - 点击绿色运行按钮（▶）
   - 或使用命令行: `jmeter -n -t jmeter/GameService_LoadTest.jmx -l jmeter/results.jtl`

### 测试计划结构

```
Test Plan
├── User Defined Variables (配置变量)
│   ├── server_host: localhost
│   ├── server_port: 8881
│   ├── threads: 100
│   ├── ramp_up: 10
│   └── duration: 60
│
├── HTTP Request Defaults (默认配置)
│   └── Server Name: ${server_host}:${server_port}
│
├── Thread Group (线程组)
│   ├── GetGameList (获取游戏列表) ✅
│   ├── GetGameDetail (获取游戏详情) ✅
│   ├── CreateGameDetail (创建游戏) ✅
│   ├── UpdateGameDetail (更新游戏) ✅ 新增
│   ├── ReviewGameVersion (审核游戏) ✅ 新增
│   └── DeleteGameDraft (删除草稿) ✅ 新增
│
├── Summary Report (结果统计)
└── View Results Tree (结果树)
```

## 测试场景

### 场景1: 基础负载测试（500并发）⭐ 默认场景

**目标**: 验证系统在500并发下的表现

**配置**:
- 线程数: 500（默认）
- 启动时间: 30秒（默认）
- 测试持续时间: 300秒（5分钟，默认）

**预期结果**:
- 响应时间 < 500ms (P95)
- 错误率 < 0.1%
- TPS > 2500
- 系统资源使用正常
- 通过率 > 99%

**运行方式**:
```bash
# 使用默认配置（500并发）
bash scripts/run_load_test.sh
```

### 场景2: 压力测试（逐步增加负载）

**目标**: 找到系统的性能瓶颈和最大容量

**配置**:
- 逐步增加线程数：200 → 500 → 800 → 1000
- 启动时间: 30秒
- 测试持续时间: 300秒（5分钟）

**预期结果**:
- 观察响应时间曲线
- 记录最大TPS
- 识别性能瓶颈
- 找到系统的最佳并发数

**运行方式**:
```bash
# 逐步增加负载，找到系统瓶颈
THREADS=200 DURATION=300 bash scripts/run_load_test.sh
THREADS=500 DURATION=300 bash scripts/run_load_test.sh
THREADS=800 DURATION=300 bash scripts/run_load_test.sh
THREADS=1000 DURATION=300 bash scripts/run_load_test.sh
```

### 场景3: 峰值测试（突发流量）

**目标**: 测试系统在突发流量下的表现和恢复能力

**配置**:
- 线程数: 2000
- 启动时间: 10秒（快速启动）
- 测试持续时间: 120秒（2分钟）

**预期结果**:
- 系统不应崩溃
- 错误率应控制在合理范围
- 观察系统恢复能力

**运行方式**:
```bash
THREADS=2000 RAMP_UP=10 DURATION=120 bash scripts/run_load_test.sh
```

### 场景4: 耐力测试（长时间运行）

**目标**: 验证系统在长时间高负载下的稳定性

**配置**:
- 线程数: 500
- 启动时间: 30秒
- 测试持续时间: 3600秒（1小时）

**预期结果**:
- 系统稳定运行
- 无内存泄漏
- 无连接池耗尽
- 性能指标稳定

**运行方式**:
```bash
THREADS=500 DURATION=3600 bash scripts/run_load_test.sh
```

## 测试用例说明

### 1. GetGameList (获取游戏列表)

**请求**:
```http
GET /api/v1/games?page_num=1&page_size=10
Host: localhost:8881
```

**JMeter 配置**:
- Method: GET
- Path: `/api/v1/games`
- Parameters:
  - `page_num`: 1
  - `page_size`: 10

### 2. GetGameDetail (获取游戏详情)

**请求**:
```http
GET /api/v1/games/12345
Host: localhost:8881
```

**JMeter 配置**:
- Method: GET
- Path: `/api/v1/games/${game_id}`
- 使用 CSV 文件读取 game_id

### 3. CreateGameDetail (创建游戏)

**请求**:
```http
POST /api/v1/games
Host: localhost:8881
Content-Type: application/json

{
  "game_detail": {
    "cp_id": 1001,
    "game_version": {
      "game_name": "Test Game",
      "game_icon": "https://example.com/icon.png",
      "header_image": "https://example.com/header.png",
      "game_introduction": "Test game description",
      "package_name": "com.test.game",
      "download_url": "https://example.com/download"
    }
  },
  "submit_mode": 1
}
```

**JMeter 配置**:
- Method: POST
- Path: `/api/v1/games`
- Body Data: JSON格式（使用变量）

### 4. UpdateGameDetail (更新游戏)

**请求**:
```http
PUT /api/v1/games/12345
Host: localhost:8881
Content-Type: application/json

{
  "game_detail": {
    "game_id": 12345,
    "cp_id": 1001,
    "game_version": {
      "game_name": "Updated Game Name",
      ...
    }
  },
  "submit_mode": 1
}
```

### 4. UpdateGameDetail (更新游戏)

**请求**:
```http
PUT /api/v1/games/12345
Host: localhost:8881
Content-Type: application/json

{
  "game_detail": {
    "game_id": "12345",
    "cp_id": 1001,
    "game_version": {
      "game_name": "Updated Game Name",
      "game_icon": "https://example.com/icon.png",
      "header_image": "https://example.com/header.png",
      "game_introduction": "Updated game description",
      "package_name": "com.test.game",
      "download_url": "https://example.com/download"
    }
  },
  "submit_mode": 1
}
```

**JMeter 配置**:
- Method: PUT
- Path: `/api/v1/games/12345`
- Body Data: JSON格式（使用变量）

### 5. ReviewGameVersion (审核游戏)

**请求**:
```http
POST /api/v1/games/review
Host: localhost:8881
Content-Type: application/json

{
  "game_id": "12345",
  "game_version_id": "67890",
  "review_result": 1,
  "review_remark": {
    "remark": "Approved",
    "operator": "admin",
    "review_time": 1733395200000,
    "meta": ""
  }
}
```

**JMeter 配置**:
- Method: POST
- Path: `/api/v1/games/review`
- Body Data: JSON格式（审核数据）

### 6. DeleteGameDraft (删除草稿)

**请求**:
```http
DELETE /api/v1/games/12345/draft
Host: localhost:8881
```

**JMeter 配置**:
- Method: DELETE
- Path: `/api/v1/games/12345/draft`

## 命令行运行

### 基本命令

```bash
# 运行测试（无GUI模式）
jmeter -n -t jmeter/GameService_LoadTest.jmx -l jmeter/results.jtl

# 生成HTML报告
jmeter -n -t jmeter/GameService_LoadTest.jmx -l jmeter/results.jtl -e -o jmeter/html-report

# 指定日志文件
jmeter -n -t jmeter/GameService_LoadTest.jmx -l jmeter/results.jtl -j jmeter/jmeter.log

# 指定测试属性
jmeter -n -t jmeter/GameService_LoadTest.jmx -l jmeter/results.jtl -Jthreads=200 -Jduration=300
```

### 使用测试脚本

```bash
# 运行压力测试脚本
bash scripts/run_load_test.sh

# 或使用 Python 脚本
python scripts/jmeter_runner.py
```

## 结果分析

### 1. 查看实时结果

在 JMeter GUI 中：
- 查看 `Summary Report` 获取统计信息
- 查看 `View Results Tree` 查看详细请求/响应
- 查看 `Response Times Over Time` 图表

### 2. 生成HTML报告

```bash
# 从结果文件生成HTML报告
jmeter -g jmeter/results.jtl -o jmeter/html-report

# 在浏览器中打开
open jmeter/html-report/index.html
```

### 3. 关键指标

- **响应时间**: P50, P95, P99
- **吞吐量**: TPS (Transactions Per Second)
- **错误率**: 失败请求百分比
- **并发用户数**: 同时在线用户数

## 性能基准

### 性能需求：支持 500 个并发用户

系统需要能够支持 **500 个并发用户**同时访问 game 服务。

### 推荐基准值

| 指标 | 目标值 | 警告值 | 说明 |
|------|--------|--------|------|
| **并发用户数** | 500 | < 300 | 系统应支持500个并发用户 |
| **P95响应时间** | < 500ms | > 1000ms | 95%的请求应在500ms内完成 |
| **P99响应时间** | < 1000ms | > 2000ms | 99%的请求应在1秒内完成 |
| **平均响应时间** | < 200ms | > 500ms | 平均响应时间应在200ms以内 |
| **TPS（吞吐量）** | > 2500 | < 1500 | 系统应能处理至少2500 TPS |
| **错误率** | < 0.1% | > 1% | 500并发下，错误率应控制在0.1%以内 |
| **CPU使用率** | < 70% | > 90% | CPU使用率不应超过70% |
| **内存使用率** | < 80% | > 90% | 内存使用率不应超过80% |

### 默认测试配置

| 参数 | 默认值 | 说明 |
|------|--------|------|
| **并发线程数** | 500 | 模拟500个并发用户 |
| **启动时间** | 30秒 | 在30秒内逐渐增加到500个线程 |
| **测试持续时间** | 300秒（5分钟） | 持续运行5分钟，观察系统稳定性 |

详细性能目标请参考：[性能目标和基准](./JMETER_PERFORMANCE_TARGETS.md)

## 故障排除

### 问题1: 连接被拒绝

**原因**: 服务未启动或端口错误

**解决**:
```bash
# 检查服务是否运行
netstat -an | grep 8881

# 检查服务日志
tail -f game_platform_api/logs/app.log
```

### 问题2: 响应时间过长

**原因**: 数据库连接池不足、查询慢、资源不足

**解决**:
- 检查数据库连接池配置
- 优化慢查询
- 增加服务器资源

### 问题3: 大量错误

**原因**: 参数验证失败、数据不存在

**解决**:
- 检查测试数据是否正确
- 查看服务日志获取详细错误信息
- 调整测试用例参数

## 最佳实践

1. **循序渐进**: 从小负载开始，逐步增加
2. **监控资源**: 同时监控CPU、内存、网络、数据库
3. **预热阶段**: 测试前先进行预热，让系统达到稳定状态
4. **多次测试**: 进行多次测试取平均值
5. **记录基准**: 建立性能基准，便于对比
6. **环境隔离**: 使用独立的测试环境，不影响生产

## 相关文档

- [JMeter 官方文档](https://jmeter.apache.org/usermanual/)
- [性能测试最佳实践](./docs/PERFORMANCE_TESTING_BEST_PRACTICES.md)
- [测试数据准备脚本](./scripts/prepare_test_data.sh)

