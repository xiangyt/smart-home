#!/bin/bash

# 智能家居项目开发环境启动脚本

echo "🏠 启动智能家居开发环境..."

# 检查是否安装了必要的工具
check_requirements() {
    echo "📋 检查环境要求..."
    
    # 检查Node.js
    if ! command -v node &> /dev/null; then
        echo "❌ Node.js 未安装，请先安装 Node.js 18+"
        exit 1
    fi
    
    # 检查Go
    if ! command -v go &> /dev/null; then
        echo "❌ Go 未安装，请先安装 Go 1.21+"
        exit 1
    fi
    
    echo "✅ 环境检查通过"
}

# 启动后端
start_backend() {
    echo "🚀 启动后端服务..."
    cd backend
    
    # 检查.env文件
    if [ ! -f .env ]; then
        echo "📝 创建.env文件..."
        cp .env.example .env
        echo "⚠️  请编辑 backend/.env 文件，配置你的 Home Assistant URL 和 Token"
    fi
    
    # 安装Go依赖
    echo "📦 安装Go依赖..."
    go mod tidy
    
    # 启动后端服务
    echo "🔧 启动后端服务 (端口: 8080)..."
    go run main.go &
    BACKEND_PID=$!
    
    cd ..
}

# 启动前端
start_frontend() {
    echo "🎨 启动前端服务..."
    cd frontend
    
    # 安装npm依赖
    echo "📦 安装npm依赖..."
    npm install
    
    # 启动前端服务
    echo "🔧 启动前端服务 (端口: 3000)..."
    npm run dev &
    FRONTEND_PID=$!
    
    cd ..
}

# 清理函数
cleanup() {
    echo "🛑 停止服务..."
    if [ ! -z "$BACKEND_PID" ]; then
        kill $BACKEND_PID 2>/dev/null
    fi
    if [ ! -z "$FRONTEND_PID" ]; then
        kill $FRONTEND_PID 2>/dev/null
    fi
    exit 0
}

# 设置信号处理
trap cleanup SIGINT SIGTERM

# 主函数
main() {
    check_requirements
    start_backend
    sleep 3  # 等待后端启动
    start_frontend
    
    echo ""
    echo "🎉 智能家居系统启动完成!"
    echo "📱 前端地址: http://localhost:3000"
    echo "🔧 后端API: http://localhost:8080"
    echo "📚 API文档: http://localhost:8080/api/health"
    echo ""
    echo "⚠️  请确保已配置 Home Assistant 连接信息"
    echo "📝 配置文件: backend/.env"
    echo ""
    echo "按 Ctrl+C 停止服务"
    
    # 等待用户中断
    wait
}

# 运行主函数
main