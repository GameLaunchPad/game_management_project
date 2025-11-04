# JMeter 压力测试快速开始

## 5分钟快速开始

### 第一步：安装 JMeter

**Windows**:
1. 下载: https://jmeter.apache.org/download_jmeter.cgi
2. 解压到任意目录（如 `C:\apache-jmeter-5.6`）
3. 设置环境变量（可选）:
   ```cmd
   set JMETER_HOME=C:\apache-jmeter-5.6
   set PATH=%PATH%;%JMETER_HOME%\bin
   ```

**Linux/Mac**:
```bash
# 使用包管理器
# Ubuntu/Debian
sudo apt-get install jmeter

# macOS
brew install jmeter

# 或从官网下载
wget https://dlcdn.apache.org//jmeter/binaries/apache-jmeter-5.6.tgz
tar -xzf apache-jmeter-5.6.tgz
export JMETER_HOME=$(pwd)/apache-jmeter-5.6
```

### 第二步：启动服务

确保以下服务已启动：

```bash
# 1. 启动 game 服务（Thrift RPC）
cd game
./output/bin/game

# 2. 启动 game_platform_api（HTTP API网关，端口 8881）
cd game_platform_api
go run main.go
```

### 第三步：运行测试

**方法1: 使用脚本（推荐）**

```bash
# 基本运行（使用默认参数）
bash scripts/run_load_test.sh

# 自定义参数
THREADS=200 RAMP_UP=30 DURATION=300 bash scripts/run_load_test.sh
```

**方法2: 使用 JMeter GUI**

```bash
# 打开 JMeter
jmeter

# 在 JMeter 中：
# 1. File > Open > jmeter/GameService_LoadTest.jmx
# 2. 配置测试参数（Thread Group）
# 3. 点击运行按钮（▶）
```

**方法3: 命令行运行**

```bash
# 无GUI模式
jmeter -n -t jmeter/GameService_LoadTest.jmx \
  -l jmeter/results.jtl \
  -e -o jmeter/html-report \
  -Jthreads=100 \
  -Jramp_up=10 \
  -Jduration=60
```

## 查看结果

### 1. HTML 报告（推荐）

```bash
# 打开HTML报告
open jmeter/html-report/index.html

# Windows
start jmeter/html-report/index.html
```

### 2. 在 JMeter GUI 中查看

1. 打开 JMeter
2. File > Open > `jmeter/results.jtl`
3. 查看 `Summary Report` 或 `View Results Tree`

## 测试参数说明

| 参数 | 说明 | 默认值 | 示例 | 性能目标 |
|------|------|--------|------|----------|
| `THREADS` | 并发线程数 | **500** | 1000 | 支持500并发用户 |
| `RAMP_UP` | 启动时间（秒） | **30** | 60 | 30秒内启动500线程 |
| `DURATION` | 测试持续时间（秒） | **300** | 600 | 持续运行5分钟 |
| `SERVER_HOST` | 服务器地址 | localhost | 192.168.1.100 | - |
| `SERVER_PORT` | 服务器端口 | 8881 | 8881 | - |

**注意**: 默认配置已针对500并发用户优化（确保通过率>99%）。如需调整，可通过环境变量覆盖。

## 测试场景示例

### 场景1: 基础负载测试（500并发）⭐ 默认

```bash
# 使用默认配置（500并发，30秒启动，5分钟持续）
bash scripts/run_load_test.sh
```

**目标**: 验证系统在500并发下的表现（性能目标，通过率>99%）

### 场景2: 压力测试（逐步增加负载）

```bash
# 逐步增加负载，找到系统瓶颈
THREADS=200 DURATION=300 bash scripts/run_load_test.sh
THREADS=500 DURATION=300 bash scripts/run_load_test.sh
THREADS=800 DURATION=300 bash scripts/run_load_test.sh
THREADS=1000 DURATION=300 bash scripts/run_load_test.sh
```

**目标**: 找到系统的性能瓶颈和最大容量（确保通过率>99%）

### 场景3: 峰值测试（突发流量）

```bash
THREADS=2000 RAMP_UP=10 DURATION=120 bash scripts/run_load_test.sh
```

**目标**: 测试系统在突发流量下的表现和恢复能力

### 场景4: 耐力测试（长时间运行）

```bash
# 测试系统在长时间高负载下的稳定性
THREADS=500 DURATION=3600 bash scripts/run_load_test.sh  # 1小时
```

**目标**: 验证系统在长时间高负载下的稳定性

## API 端点

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | `/api/v1/games?page_num=1&page_size=10` | 获取游戏列表 |
| GET | `/api/v1/games/:id` | 获取游戏详情 |
| POST | `/api/v1/games` | 创建游戏 |
| PUT | `/api/v1/games/:id` | 更新游戏 |
| POST | `/api/v1/games/review` | 审核游戏 |
| DELETE | `/api/v1/games/:id/draft` | 删除草稿 |

## 常见问题

### Q: 连接被拒绝

**A**: 确保服务已启动
```bash
# 检查端口
netstat -an | grep 8881

# 检查服务日志
tail -f game_platform_api/logs/app.log
```

### Q: JMeter 命令找不到

**A**: 设置环境变量或使用完整路径
```bash
export JMETER_HOME=/path/to/apache-jmeter
$JMETER_HOME/bin/jmeter -n -t ...
```

### Q: 如何查看实时结果

**A**: 使用 JMeter GUI 或查看 HTML 报告
```bash
# 生成HTML报告
jmeter -g jmeter/results.jtl -o jmeter/html-report
open jmeter/html-report/index.html
```

## 下一步

1. 查看详细文档: [JMETER_LOAD_TESTING.md](./JMETER_LOAD_TESTING.md)
2. 自定义测试计划: 编辑 `jmeter/GameService_LoadTest.jmx`
3. 分析性能瓶颈: 查看HTML报告中的性能指标

