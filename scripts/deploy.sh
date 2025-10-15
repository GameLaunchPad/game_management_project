#!/bin/bash
# 服务部署脚本模板
# 使用方式：在构建时替换变量后使用

APP_NAME="{{APP_NAME}}"
APP_PORT="{{APP_PORT}}"
APP_HOME="/home/admin/services/$APP_NAME"

echo "========================================="
echo "部署服务: $APP_NAME"
echo "端口: $APP_PORT"
echo "时间: $(date '+%Y-%m-%d %H:%M:%S')"
echo "========================================="

# 创建目录
mkdir -p $APP_HOME

# 停止旧进程
PID=$(lsof -ti:$APP_PORT 2>/dev/null)
if [ -n "$PID" ]; then
    echo "检测到旧进程 PID: $PID，正在停止..."
    kill -15 $PID

    # 等待进程优雅关闭（最多10秒）
    for i in {1..10}; do
        if ! kill -0 $PID 2>/dev/null; then
            echo "✅ 旧进程已停止"
            break
        fi
        sleep 1
    done

    # 强制杀死（如果还在运行）
    if kill -0 $PID 2>/dev/null; then
        echo "⚠️  强制停止进程"
        kill -9 $PID 2>/dev/null
    fi
fi

# 启动新服务
cd $APP_HOME
chmod +x $APP_NAME

nohup ./$APP_NAME > app.log 2>&1 &
NEW_PID=$!

# 验证启动
sleep 3
if kill -0 $NEW_PID 2>/dev/null; then
    echo "✅ 服务启动成功！PID: $NEW_PID"
    echo "日志文件: $APP_HOME/app.log"

    # 检查端口监听
    if lsof -ti:$APP_PORT >/dev/null 2>&1; then
        echo "✅ 端口 $APP_PORT 监听正常"
    else
        echo "⚠️  警告：端口 $APP_PORT 未监听"
    fi
else
    echo "❌ 服务启动失败！"
    echo "查看日志: tail -n 50 $APP_HOME/app.log"
    exit 1
fi

echo "========================================="