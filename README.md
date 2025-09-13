# 智能家居管理系统

基于Vue3 + Golang的前后端分离智能家居管理系统，集成Home Assistant数据。

## 项目结构

```
smart-home/
├── frontend/          # Vue3前端项目
├── backend/           # Golang后端项目
├── docker-compose.yml # Docker编排文件
└── README.md
```

## 功能特性

- 🏠 智能家居设备概览（参考HA设计）
- 📊 实时设备状态监控
- 🔌 设备控制面板
- 📱 响应式设计
- 🔗 Home Assistant API集成

## 快速开始

### 前端开发
```bash
cd frontend
npm install
npm run dev
```

### 后端开发
```bash
cd backend
go mod tidy
go run main.go
```

## 技术栈

### 前端
- Vue 3
- TypeScript
- Vite
- Element Plus
- Axios

### 后端
- Golang
- Gin框架
- GORM
- SQLite/MySQL