#!/bin/bash

# 智能家居项目构建脚本

echo "🏗️  构建智能家居项目..."

# 构建前端
build_frontend() {
    echo "🎨 构建前端..."
    cd frontend
    
    # 安装依赖
    npm install
    
    # 构建生产版本
    npm run build
    
    echo "✅ 前端构建完成"
    cd ..
}

# 构建后端
build_backend() {
    echo "🔧 构建后端..."
    cd backend
    
    # 安装Go依赖
    go mod tidy
    
    # 构建二进制文件
    CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o smart-home-backend .
    
    echo "✅ 后端构建完成"
    cd ..
}

# 主函数
main() {
    build_backend
    build_frontend
    
    echo ""
    echo "🎉 项目构建完成!"
    echo "📁 前端构建文件: frontend/dist/"
    echo "📁 后端二进制文件: backend/smart-home-backend"
    echo ""
    echo "🚀 使用 Docker 部署:"
    echo "   docker-compose up -d"
}

# 运行主函数
main