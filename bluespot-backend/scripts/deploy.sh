#!/bin/bash

# 设置错误时退出
set -e

# 获取脚本所在目录的上一级目录（项目根目录）
PROJECT_ROOT=$(cd "$(dirname "$0")/.." && pwd)
cd "$PROJECT_ROOT"

echo "🚀 开始构建项目..."

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
rsync -avz --progress --partial ./bluespot kr:/opt/apps/bluespot-backend
rsync -avz --progress --partial --exclude='video-export-flat' --exclude='video-export-grouped' --exclude='larges' --exclude='chunks' --exclude='uploads' ./public kr:/opt/apps/bluespot-backend
rsync -avz --delete --progress --partial ./docs kr:/opt/apps/bluespot-backend
rsync -avz --progress --partial ./internal/config/app.yml kr:/opt/apps/bluespot-backend
rsync -avz --progress --partial ./internal/config/prod.yml kr:/opt/apps/bluespot-backend
rsync -avz --progress --partial ./internal/data kr:/opt/apps/bluespot-backend
echo "✅ 上传完成！"

echo "🔄 重启远程 bluespot 服务..."
ssh kr 'systemctl restart bluespot'

echo "📋 查看远程 bluespot 服务状态..."
ssh kr 'systemctl status bluespot'
