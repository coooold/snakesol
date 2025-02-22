#!/bin/bash

# 记录进程ID的文件
PID_FILE=".snake.pid"
# 停止标志文件
STOP_FILE=".stop"

# 检查PID文件是否存在
if [ ! -f "$PID_FILE" ]; then
    echo "服务器似乎没有运行。"
    exit 0
fi

# 读取PID并终止进程
PID=$(cat "$PID_FILE")
if kill -0 $PID 2>/dev/null; then
    echo "正在停止服务器进程..."
    kill $PID
    rm -f "$PID_FILE"
    # 创建停止标志文件
    touch "$STOP_FILE"
    echo "服务器已停止。"
else
    echo "服务器进程已不存在。"
    rm -f "$PID_FILE"
    # 创建停止标志文件
    touch "$STOP_FILE"
fi