#!/bin/bash

echo "================================================"
echo "CDK 充值系统 - 本地开发启动"
echo "================================================"
echo ""

# 检查环境变量文件
ENV_FILE="${ENV_FILE:-.env}"
if [ ! -f "$ENV_FILE" ] && [ -f .env.local ]; then
    ENV_FILE=.env.local
fi
if [ ! -f "$ENV_FILE" ]; then
    echo "✗ 未找到 $ENV_FILE 或 .env.local 文件"
    exit 1
fi
export ENV_FILE="$(cd "$(dirname "$ENV_FILE")" && pwd)/$(basename "$ENV_FILE")"

echo "✓ 环境文件已找到"
echo ""

DB_DRIVER=$(sed -n 's/^DB_DRIVER=//p' "$ENV_FILE" | tail -n 1)
if [ "$DB_DRIVER" = "postgres" ] || [ "$DB_DRIVER" = "postgresql" ] || [ "$DB_DRIVER" = "pgx" ]; then
    echo "检查 PostgreSQL 连接..."
    set -a
    . "$ENV_FILE"
    set +a
    psql -U "${DB_USER:-postgres}" -h "${DB_HOST:-localhost}" -p "${DB_PORT:-5432}" -d "${DB_NAME:-cdk_recharge}" -c "SELECT NOW();" >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo "✗ PostgreSQL 连接失败"
        echo "  请确保 PostgreSQL 已启动并检查 $ENV_FILE"
        exit 1
    fi
    echo "✓ PostgreSQL 连接成功"
else
    echo "✓ 使用 SQLite 数据库"
fi

# 检查 Redis
redis-cli ping >/dev/null 2>&1
if [ $? -ne 0 ]; then
    echo "✗ Redis 连接失败 (可选，不影响基本功能)"
else
    echo "✓ Redis 连接成功"
fi

echo ""
echo "启动服务..."
echo ""

# 启动后端
cd backend
echo "启动后端服务 (端口 8080)..."
go run ./cmd/server/main.go &
BACKEND_PID=$!
echo "  后端 PID: $BACKEND_PID"

sleep 2

# 启动前端
cd ../frontend
echo "启动前端服务 (端口 5173)..."
npm run dev &
FRONTEND_PID=$!
echo "  前端 PID: $FRONTEND_PID"

echo ""
echo "================================================"
echo "✓ 服务已启动！"
echo "================================================"
echo ""
echo "前端:    http://localhost:5173"
echo "后端:    http://localhost:8080"
echo "API:     http://localhost:8080/api/v1"
echo "健康检查: curl http://localhost:8080/health"
echo ""
echo "按 Ctrl+C 停止所有服务..."
echo ""

# 等待中断信号
trap "kill $BACKEND_PID $FRONTEND_PID; exit" SIGINT SIGTERM

# 保持脚本运行
wait
