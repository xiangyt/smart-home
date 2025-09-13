# 🚀 智能家居系统 - 快速启动指南

## 📋 启动前检查

✅ **环境要求**
- Node.js 18+ 已安装
- Go 1.21+ 已安装
- Home Assistant 实例 (可选)

## 🔧 配置步骤

### 1. 配置 Home Assistant 连接

编辑 `backend/.env` 文件：

```bash
# 复制示例配置
cp backend/.env.example backend/.env

# 编辑配置文件
vim backend/.env
```

配置内容：
```env
HOME_ASSISTANT_URL=http://your-ha-ip:8123
HOME_ASSISTANT_TOKEN=your_long_lived_access_token
```

> 💡 **获取 HA Token**: 登录 HA → 用户设置 → 安全 → 长期访问令牌

### 2. 启动开发环境

#### 方式一：一键启动 (推荐)
```bash
./scripts/start-dev.sh
```

#### 方式二：分别启动

**后端服务:**
```bash
cd backend
go mod tidy
go run main.go
# 服务运行在 http://localhost:8080
```

**前端服务:**
```bash
cd frontend
npm install
npm run dev
# 服务运行在 http://localhost:3000
```

## 🌐 访问地址

- **前端界面**: http://localhost:3000
- **后端API**: http://localhost:8080
- **健康检查**: http://localhost:8080/api/health

## 🎯 功能测试

### 1. 概览页面
- 访问 http://localhost:3000/overview
- 查看设备状态卡片
- 测试设备控制功能

### 2. 设备管理
- 访问 http://localhost:3000/devices
- 查看设备列表
- 点击设备详情

### 3. API 测试
```bash
# 健康检查
curl http://localhost:8080/api/health

# 获取设备列表
curl http://localhost:8080/api/devices
```

## 🐳 Docker 部署

### 快速部署
```bash
# 1. 配置环境变量
echo "HOME_ASSISTANT_URL=http://your-ha-ip:8123" > .env
echo "HOME_ASSISTANT_TOKEN=your_token" >> .env

# 2. 启动服务
docker-compose up -d

# 3. 查看日志
docker-compose logs -f
```

### 停止服务
```bash
docker-compose down
```

## 🔍 故障排除

### 常见问题

**1. 后端启动失败**
```bash
# 检查Go版本
go version

# 重新安装依赖
cd backend && go mod tidy
```

**2. 前端启动失败**
```bash
# 检查Node版本
node --version

# 清理并重装
cd frontend
rm -rf node_modules package-lock.json
npm install
```

**3. 无法连接 Home Assistant**
- 检查 HA URL 是否正确
- 验证 Token 是否有效
- 确认网络连通性

**4. 设备列表为空**
- 确认 HA 中有设备
- 检查 Token 权限
- 查看后端日志

### 日志查看
```bash
# 后端日志
cd backend && go run main.go

# 前端日志
cd frontend && npm run dev

# Docker 日志
docker-compose logs backend
docker-compose logs frontend
```

## 📱 界面预览

### 概览页面特性
- 📊 设备状态统计卡片
- 💡 灯光设备控制面板
- 🔌 开关设备控制面板
- 📈 传感器数据展示

### 设备管理特性
- 📋 设备列表表格
- 🏷️ 设备类型标签
- ⚡ 实时状态更新
- 📝 设备详情弹窗

## 🎨 自定义配置

### 修改主题色彩
编辑 `frontend/src/App.vue`:
```css
:root {
  --primary-color: #409EFF;    /* 主色调 */
  --success-color: #67C23A;    /* 成功色 */
  --warning-color: #E6A23C;    /* 警告色 */
  --danger-color: #F56C6C;     /* 危险色 */
}
```

### 添加新设备类型
1. 后端: 修改 `backend/internal/services/device.go`
2. 前端: 更新 `frontend/src/stores/devices.ts`
3. 界面: 在 `frontend/src/views/Overview.vue` 添加卡片

## 📚 更多文档

- [详细安装指南](SETUP.md)
- [项目总结](PROJECT_SUMMARY.md)
- [API 文档](README.md)

## 🆘 获取帮助

遇到问题？
1. 查看 [故障排除](#故障排除) 部分
2. 检查项目 Issues
3. 创建新的 Issue 描述问题

---

**🎉 享受你的智能家居管理系统！**