#!/bin/bash

# 发生错误、使用未定义变量或管道命令失败时退出
set -euo pipefail

# GitHub Actions 通过环境变量指定服务器；本地继续兼容 ~/.ssh/config 中的 kr。
if [ -n "${SERVER_HOST:-}" ] && [ -n "${SERVER_USER:-}" ]; then
    DEPLOY_TARGET="${SERVER_USER}@${SERVER_HOST}"
elif [ -z "${SERVER_HOST:-}" ] && [ -z "${SERVER_USER:-}" ]; then
    DEPLOY_TARGET="${DEPLOY_TARGET:-kr}"
else
    echo "❌ SERVER_HOST 和 SERVER_USER 必须同时设置"
    exit 1
fi

# 获取脚本所在目录的上一级目录（项目根目录）
PROJECT_ROOT=$(cd "$(dirname "$0")/.." && pwd)
cd "$PROJECT_ROOT"

export TZ="Asia/Shanghai"

echo "🚀 开始构建项目..."
echo "🎯 部署目标: $DEPLOY_TARGET"
echo "🕒 构建时区: $TZ"

# 1. 检查并创建 public 目录
if [ ! -d "public" ]; then
    echo "📂 创建 public 目录..."
    mkdir -p public
fi

# 2. 生成 meta.json
echo "📄 生成 public/meta.json..."
CURRENT_TIME=$(date '+%Y-%m-%d %H:%M:%S')
echo "{\"deployTime\": \"$CURRENT_TIME\"}" > public/meta.json

# 3. 生成 Swagger 文档
echo "📚 生成 Swagger 文档..."
go run github.com/swaggo/swag/cmd/swag@v1.16.6 init -g ./cmd/main.go -o ./docs

# 4. 编译项目
echo "🔨 正在编译..."
# 设置环境变量进行交叉编译（如需在本机运行可去掉这些变量，这里保留原有逻辑）
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64

go build -o bluespot ./cmd/main.go

echo "✅ 构建成功！"
echo "📍 输出文件: $PROJECT_ROOT/bluespot"
echo "📅 部署时间: $CURRENT_TIME"

echo "📤 开始上传文件到服务器..."
rsync -avz --progress --partial ./bluespot "${DEPLOY_TARGET}:/opt/apps/bluespot-backend"
rsync -avz --progress --partial --exclude='video-export-flat' --exclude='video-export-grouped' --exclude='larges' --exclude='chunks' --exclude='uploads' ./public "${DEPLOY_TARGET}:/opt/apps/bluespot-backend"
rsync -avz --delete --progress --partial ./docs "${DEPLOY_TARGET}:/opt/apps/bluespot-backend"
rsync -avz --progress --partial ./internal/config/app.yml "${DEPLOY_TARGET}:/opt/apps/bluespot-backend"
if [ "$DEPLOY_TARGET" = "kr" ]; then
    rsync -avz --progress --partial ./internal/config/prod.yml "${DEPLOY_TARGET}:/opt/apps/bluespot-backend"
fi
rsync -avz --progress --partial ./internal/data "${DEPLOY_TARGET}:/opt/apps/bluespot-backend"
echo "✅ 上传完成！"

echo "🔄 重启远程 bluespot 服务..."
ssh "$DEPLOY_TARGET" 'systemctl restart bluespot'

echo "📋 查看远程 bluespot 服务状态..."
ssh "$DEPLOY_TARGET" 'systemctl status bluespot'
