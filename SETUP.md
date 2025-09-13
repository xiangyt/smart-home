# 智能家居管理系统 - 安装配置指南

## 项目概述

这是一个基于 Vue3 + Golang 的前后端分离智能家居管理系统，可以与 Home Assistant 集成，提供美观的设备管理界面。

## 功能特性

- 🏠 **设备概览**: 参考 Home Assistant 设计的概览页面
- 📊 **实时监控**: 实时显示设备状态和传感器数据
- 🎛️ **设备控制**: 支持灯光、开关等设备的远程控制
- 📱 **响应式设计**: 适配桌面和移动设备
- 🔗 **HA集成**: 通过 API 与 Home Assistant 无缝集成

## 技术栈

### 前端
- Vue 3 + TypeScript
- Vite 构建工具
- Element Plus UI 组件库
- Pinia 状态管理
- Axios HTTP 客户端

### 后端
- Golang + Gin 框架
- GORM ORM 框架
- SQLite 数据库
- Home Assistant REST API 集成

## 快速开始

### 1. 环境要求

- Node.js 18+
- Go 1.21+
- Home Assistant 实例 (可选，用于真实数据)

### 2. 克隆项目

```bash
git clone <your-repo-url>
cd smart-home
```

### 3. 配置 Home Assistant 连接

编辑 `backend/.env` 文件：

```env
# Home Assistant 配置
HOME_ASSISTANT_URL=http://your-ha-ip:8123
HOME_ASSISTANT_TOKEN=your_long_lived_access_token
```

#### 获取 Home Assistant Token:

1. 登录 Home Assistant
2. 点击左下角用户头像
3. 滚动到底部，点击 "创建令牌"
4. 输入令牌名称，点击 "确定"
5. 复制生成的令牌到 `.env` 文件

### 4. 启动开发环境

#### 方式一：使用启动脚本 (推荐)

```bash
./scripts/start-dev.sh
```

#### 方式二：手动启动

**启动后端:**
```bash
cd backend
go mod tidy
go run main.go
```

**启动前端:**
```bash
cd frontend
npm install
npm run dev
```

### 5. 访问应用

- 前端界面: http://localhost:3000
- 后端API: http://localhost:8080
- 健康检查: http://localhost:8080/api/health

## Docker 部署

### 1. 使用 Docker Compose

```bash
# 构建并启动
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

### 2. 环境变量配置

创建 `.env` 文件在项目根目录：

```env
HOME_ASSISTANT_URL=http://your-ha-ip:8123
HOME_ASSISTANT_TOKEN=your_token_here
```

## API 接口文档

### 设备相关接口

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | `/api/devices` | 获取所有设备列表 |
| GET | `/api/devices/{entity_id}/state` | 获取设备状态 |
| POST | `/api/devices/{entity_id}/control` | 控制设备 |
| GET | `/api/devices/{entity_id}/history` | 获取设备历史 |

### 控制设备示例

```bash
# 开启灯光
curl -X POST http://localhost:8080/api/devices/light.living_room/control \
  -H "Content-Type: application/json" \
  -d '{"service": "turn_on", "data": {}}'

# 关闭开关
curl -X POST http://localhost:8080/api/devices/switch.bedroom/control \
  -H "Content-Type: application/json" \
  -d '{"service": "turn_off", "data": {}}'
```

## 项目结构

```
smart-home/
├── frontend/                 # Vue3 前端项目
│   ├── src/
│   │   ├── components/      # 组件
│   │   ├── views/          # 页面
│   │   ├── stores/         # Pinia 状态管理
│   │   ├── api/            # API 接口
│   │   └── router/         # 路由配置
│   ├── package.json
│   └── vite.config.ts
├── backend/                 # Golang 后端项目
│   ├── internal/
│   │   ├── config/         # 配置管理
│   │   ├── database/       # 数据库
│   │   ├── handlers/       # HTTP 处理器
│   │   ├── models/         # 数据模型
│   │   └── services/       # 业务逻辑
│   ├── main.go
│   └── go.mod
├── scripts/                # 脚本文件
├── docker-compose.yml      # Docker 编排
└── README.md
```

## 开发指南

### 添加新设备类型

1. 在 `backend/internal/services/device.go` 中更新 `shouldSkipEntity` 函数
2. 在前端 `stores/devices.ts` 中添加新的设备类型处理
3. 在 `views/Overview.vue` 中添加新的设备卡片

### 自定义主题

编辑 `frontend/src/App.vue` 中的 CSS 变量：

```css
:root {
  --primary-color: #409EFF;
  --success-color: #67C23A;
  --warning-color: #E6A23C;
  --danger-color: #F56C6C;
}
```

## 故障排除

### 常见问题

1. **无法连接 Home Assistant**
   - 检查 HA URL 和 Token 是否正确
   - 确保 HA 实例可访问
   - 检查防火墙设置

2. **前端无法访问后端 API**
   - 检查后端是否在 8080 端口运行
   - 查看浏览器控制台错误信息
   - 检查 CORS 配置

3. **设备列表为空**
   - 确认 HA Token 有足够权限
   - 检查 HA 中是否有设备
   - 查看后端日志

### 日志查看

```bash
# 查看后端日志
cd backend && go run main.go

# 查看前端日志
cd frontend && npm run dev

# Docker 日志
docker-compose logs backend
docker-compose logs frontend
```

## 贡献指南

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 支持

如果你觉得这个项目有用，请给它一个 ⭐️！

有问题或建议？请创建一个 [Issue](../../issues)。