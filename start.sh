#!/bin/bash

# 记录进程ID的文件
PID_FILE=".snake.pid"
# 停止标志文件
STOP_FILE=".stop"

# 检查是否已经在运行
if [ -f "$PID_FILE" ]; then
    echo "服务器似乎已经在运行中。如果确定没有运行，请删除 $PID_FILE 文件后重试。"
    exit 1
fi

# 启动服务器的函数
start_server() {
    ./bin/server -port 12080 &
    echo $! > "$PID_FILE"
}

# 确保在脚本退出时清理PID文件
cleanup() {
    rm -f "$PID_FILE"
    exit 0
}

trap cleanup SIGINT SIGTERM

# 编译服务器
make build

# 持续监控并在需要时重启服务器
while true; do
    if [ ! -f "$PID_FILE" ]; then
        echo "启动服务器..."
        start_server
    fi
    
    # 检查进程是否存在
    if ! kill -0 $(cat "$PID_FILE") 2>/dev/null; then
        # 检查是否存在停止标志文件
        if [ -f "$STOP_FILE" ]; then
            echo "检测到停止标志，服务器将不会重启。"
            rm -f "$STOP_FILE"
            rm -f "$PID_FILE"
            exit 0
        fi
        
        echo "服务器进程已终止，正在重启..."
        rm -f "$PID_FILE"
        start_server
    fi
    
    sleep 5
done