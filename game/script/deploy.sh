#!/bin/bash

# ========== 配置区 ==========
APP_NAME="kitex-server"
APP_HOME="/home/admin/kitex-app"
LOG_FILE="$APP_HOME/app.log"

# ========== 函数定义 ==========
stop_app() {
    echo "正在停止旧服务..."
    PID=$(pgrep -f "$APP_HOME/$APP_NAME")
    if [ -n "$PID" ]; then
        echo "找到进程 PID: $PID, 正在停止..."
        kill -15 $PID
        # 等待进程优雅关闭
        for i in {1..10}; do
            if ! kill -0 $PID 2>/dev/null; then
                echo "进程已停止"
                return 0
            fi
            sleep 1
        done
        # 强制杀死
        kill -9 $PID 2>/dev/null
        echo "进程已强制停止"
    else
        echo "未找到运行中的进程"
    fi
}

start_app() {
    echo "正在启动新服务..."
    cd $APP_HOME
    chmod +x $APP_NAME

    # 后台启动，日志重定向
    nohup ./$APP_NAME > $LOG_FILE 2>&1 &

    NEW_PID=$!
    echo "服务已启动，PID: $NEW_PID"

    # 检查启动是否成功
    sleep 2
    if kill -0 $NEW_PID 2>/dev/null; then
        echo "✅ 服务启动成功！"
        return 0
    else
        echo "❌ 服务启动失败，请查看日志: $LOG_FILE"
        return 1
    fi
}

# ========== 主流程 ==========
echo "========================================="
echo "开始部署 Kitex 服务"
echo "时间: $(date '+%Y-%m-%d %H:%M:%S')"
echo "========================================="

# 创建目录
mkdir -p $APP_HOME

# 停止旧服务
stop_app

# 启动新服务
start_app

echo "========================================="
echo "部署完成"
echo "========================================="
