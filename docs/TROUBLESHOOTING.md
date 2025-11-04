# 故障排除指南

## 常见问题

### 问题1: `go.mod file not found` 错误

**错误信息**:
```
cover: no required module provides package github.com/GameLaunchPad/game_management_project/game/handler: go.mod file not found
```

**原因**:
- 工作目录不正确
- `-coverpkg` 参数使用了错误的包路径格式
- Go模块路径解析问题

**解决方案**:

1. **确保在正确的目录下运行**
   - 脚本会自动切换到 `game` 目录
   - 确保 `game/go.mod` 文件存在

2. **检查Go模块配置**
   ```bash
   cd game
   go list -m
   go list ./handler
   ```

3. **如果问题持续，脚本会自动降级**
   - 脚本会先尝试使用 `-coverpkg=./handler`
   - 如果失败，会自动重试不使用 `-coverpkg` 参数
   - 覆盖率报告会过滤为仅handler目录

**在云效流水线中**:
- 确保工作目录设置为 `${WORKSPACE}` 或项目根目录
- 脚本会自动切换到正确的模块目录

### 问题2: 覆盖率统计包含其他包

**原因**:
- `-coverpkg` 参数失败，脚本降级为不使用该参数

**解决方案**:
- 脚本会自动过滤覆盖率报告，只显示handler目录
- 查看覆盖率报告时，只关注 `handler/` 开头的行

### 问题3: 测试执行失败

**检查步骤**:
1. 查看JSONL报告中的错误信息
2. 检查Go环境是否正确配置
3. 检查依赖是否已下载

**调试命令**:
```bash
cd game
go list ./...
go test -v ./handler
```

