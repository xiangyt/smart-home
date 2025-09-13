<template>
  <div id="app">
    <el-container class="layout-container">
      <!-- 侧边栏 -->
      <el-aside :width="sidebarCollapsed ? '80px' : '280px'" class="sidebar" :class="{ collapsed: sidebarCollapsed }">
        <div class="logo">
          <div class="logo-content">
            <div class="logo-icon">
              <el-icon size="32"><House /></el-icon>
            </div>
            <div class="logo-text" v-show="!sidebarCollapsed">
              <h2>智能家居</h2>
            </div>
          </div>
        </div>
        <el-menu
          :default-active="$route.path"
          class="sidebar-menu"
          router
          background-color="transparent"
          text-color="#a0a8b8"
          active-text-color="#ffffff"
          :collapse="sidebarCollapsed"
        >
          <el-menu-item index="/overview" class="menu-item">
            <el-icon><House /></el-icon>
            <span>概览</span>
          </el-menu-item>
          <el-menu-item index="/devices" class="menu-item">
            <el-icon><Monitor /></el-icon>
            <span>设备管理</span>
          </el-menu-item>
          <el-menu-item index="/scenes" class="menu-item">
            <el-icon><Setting /></el-icon>
            <span>场景控制</span>
          </el-menu-item>
          <el-menu-item index="/automation" class="menu-item">
            <el-icon><Timer /></el-icon>
            <span>自动化</span>
          </el-menu-item>
        </el-menu>
        
        <!-- 收起按钮 - 右上角 -->
        <div class="collapse-btn" @click="toggleSidebar">
          <el-icon size="18">
            <Expand v-if="sidebarCollapsed" />
            <Fold v-else />
          </el-icon>
        </div>
        
        <!-- 侧边栏底部信息 -->
        <div class="sidebar-footer">
          <div class="system-status">
            <div class="status-item">
              <div class="status-dot" :class="systemStatus.class"></div>
              <span>{{ systemStatus.text }}</span>
            </div>
            <div class="status-item">
              <el-icon size="14"><Clock /></el-icon>
              <span>{{ currentTime }}</span>
            </div>
          </div>
        </div>
      </el-aside>

      <!-- 主内容区 -->
      <el-container>
        <el-header class="header">
          <div class="header-content">
            <div class="header-left">
              <h3>{{ pageTitle }}</h3>
              <p class="header-subtitle">{{ pageSubtitle }}</p>
            </div>
            <div class="header-actions">
              <el-avatar class="user-avatar" :size="40">
                <el-icon><User /></el-icon>
              </el-avatar>
            </div>
          </div>
        </el-header>
        
        <el-main class="main-content">
          <div class="content-wrapper">
            <router-view />
          </div>
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import { House, Monitor, Setting, Timer, Clock, User, Expand, Fold } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import api from './api'

const route = useRoute()
const currentTime = ref('')

// 侧边栏收起状态
const sidebarCollapsed = ref(false)

// 系统状态
const systemStatus = ref({
  online: false,
  class: 'offline',
  text: '系统离线'
})

// 切换侧边栏
const toggleSidebar = () => {
  sidebarCollapsed.value = !sidebarCollapsed.value
  // 保存状态到本地存储
  localStorage.setItem('sidebarCollapsed', sidebarCollapsed.value.toString())
}

const pageTitle = computed(() => {
  const titles: Record<string, string> = {
    '/overview': '系统概览',
    '/devices': '设备管理',
    '/scenes': '场景控制',
    '/automation': '自动化规则'
  }
  return titles[route.path] || '智能家居'
})

const pageSubtitle = computed(() => {
  const subtitles: Record<string, string> = {
    '/overview': '实时监控您的智能家居设备状态',
    '/devices': '管理和控制所有智能设备',
    '/scenes': '创建和管理智能场景',
    '/automation': '设置自动化规则和定时任务'
  }
  return subtitles[route.path] || '欢迎使用智能家居管理系统'
})

// 更新时间
const updateTime = () => {
  const now = new Date()
  currentTime.value = now.toLocaleTimeString('zh-CN', { 
    hour12: false,
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 检查系统状态
const checkSystemStatus = async () => {
  try {
    // 使用专门的健康检查接口
    await api.get('/health', { timeout: 3000 })
    systemStatus.value = {
      online: true,
      class: 'online',
      text: '系统在线'
    }
  } catch (error) {
    systemStatus.value = {
      online: false,
      class: 'offline',
      text: '系统离线'
    }
  }
}



let timeInterval: NodeJS.Timeout
let statusInterval: NodeJS.Timeout

onMounted(() => {
  updateTime()
  timeInterval = setInterval(updateTime, 1000)
  
  // 立即检查一次系统状态
  checkSystemStatus()
  // 每5秒检查一次系统状态
  statusInterval = setInterval(checkSystemStatus, 5000)
  
  // 从本地存储恢复侧边栏状态
  const savedState = localStorage.getItem('sidebarCollapsed')
  if (savedState !== null) {
    sidebarCollapsed.value = savedState === 'true'
  }
})

onUnmounted(() => {
  if (timeInterval) {
    clearInterval(timeInterval)
  }
  if (statusInterval) {
    clearInterval(statusInterval)
  }
})
</script>

<style scoped>
.layout-container {
  height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.sidebar {
  background: linear-gradient(180deg, #2c3e50 0%, #34495e 100%);
  color: white;
  box-shadow: 2px 0 10px rgba(0, 0, 0, 0.1);
  position: relative;
  overflow: hidden;
  transition: width 0.8s cubic-bezier(0.4, 0, 0.2, 1);
}

.sidebar.collapsed {
  width: 80px !important;
}

.sidebar.collapsed .logo-text {
  opacity: 0;
  transform: translateX(-20px);
  transition: opacity 0.2s ease, transform 0.3s ease;
}

.sidebar.collapsed .system-status {
  opacity: 0;
  transition: opacity 0.2s ease;
}

.sidebar:not(.collapsed) .logo-text {
  opacity: 1;
  transform: translateX(0);
  transition: opacity 0.6s ease 0.4s, transform 0.6s ease 0.4s;
}

.sidebar:not(.collapsed) .system-status {
  opacity: 1;
  transition: opacity 0.6s ease 0.4s;
}

.sidebar.collapsed .sidebar-menu {
  padding: 20px 0;
}

.sidebar.collapsed .el-menu-item {
  padding: 0 !important;
  justify-content: center;
  margin: 8px 8px;
  width: calc(100% - 16px);
  height: 48px !important;
  min-height: 48px !important;
  max-height: 48px !important;
  display: flex !important;
  align-items: center !important;
  line-height: 48px !important;
  transition: all 0.8s cubic-bezier(0.4, 0, 0.2, 1) !important;
}

.sidebar.collapsed .el-menu-item .el-icon {
  margin-right: 0 !important;
  font-size: 18px;
  line-height: 18px;
  height: 18px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: margin-right 0.8s cubic-bezier(0.4, 0, 0.2, 1), font-size 0.8s cubic-bezier(0.4, 0, 0.2, 1);
}

.sidebar.collapsed .el-menu-item span {
  opacity: 0;
  transform: translateX(-10px);
  transition: opacity 0.1s ease, transform 0.1s ease;
  pointer-events: none;
}

.sidebar:not(.collapsed) .el-menu-item span {
  opacity: 1;
  transform: translateX(0);
  transition: opacity 0.4s ease 0.8s, transform 0.4s ease 0.8s;
}

.sidebar.collapsed .logo {
  padding: 20px 10px;
  height: 80px;
  min-height: 80px;
}

.sidebar.collapsed .logo-content {
  justify-content: center;
}

.sidebar::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: url('data:image/svg+xml,<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><defs><pattern id="grain" width="100" height="100" patternUnits="userSpaceOnUse"><circle cx="25" cy="25" r="1" fill="rgba(255,255,255,0.02)"/><circle cx="75" cy="75" r="1" fill="rgba(255,255,255,0.02)"/></pattern></defs><rect width="100" height="100" fill="url(%23grain)"/></svg>');
  pointer-events: none;
}

.logo {
  padding: 30px 20px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  position: relative;
  z-index: 1;
  height: 80px;
  min-height: 80px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: padding 0.8s cubic-bezier(0.4, 0, 0.2, 1);
}

.logo-content {
  display: flex;
  align-items: center;
  gap: 12px;
  transition: justify-content 0.8s cubic-bezier(0.4, 0, 0.2, 1);
}

.sidebar.collapsed .logo {
  padding: 20px 10px;
}

.sidebar.collapsed .logo-icon {
  margin-bottom: 0;
}

.logo-icon {
  color: #3498db;
  flex-shrink: 0;
}

.logo h2 {
  margin: 0;
  color: #ffffff;
  font-size: 20px;
  font-weight: 600;
  letter-spacing: 1px;
  white-space: nowrap;
}

.sidebar-menu {
  border: none;
  padding: 20px 0;
  position: relative;
  z-index: 1;
}

.el-menu-item .el-icon {
  font-size: 18px;
  line-height: 18px;
  height: 18px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: margin-right 0.8s cubic-bezier(0.4, 0, 0.2, 1), font-size 0.8s cubic-bezier(0.4, 0, 0.2, 1);
}

.sidebar.collapsed .sidebar-menu {
  padding: 20px 0;
}

.sidebar.collapsed .el-menu-item {
  padding: 0 !important;
  justify-content: center;
  margin: 8px 8px;
  width: 64px;
  height: 48px;
  display: flex;
  align-items: center;
}

.sidebar.collapsed .el-menu-item .el-icon {
  margin-right: 0 !important;
  font-size: 20px;
}

.sidebar.collapsed .el-menu-item span {
  display: none;
}

.menu-item {
  position: relative;
  overflow: hidden;
  border-radius: 12px;
  transition: none;
}

.sidebar:not(.collapsed) .menu-item {
  margin: 8px 16px;
  height: 48px;
  display: flex;
  align-items: center;
  transition: transform 0.5s cubic-bezier(0.4, 0, 0.2, 1), margin 0.8s cubic-bezier(0.4, 0, 0.2, 1);
}

.sidebar.collapsed .menu-item {
  margin: 8px 8px;
  width: 64px;
  height: 48px !important;
  display: flex !important;
  align-items: center !important;
  transition: transform 0.5s cubic-bezier(0.4, 0, 0.2, 1), margin 0.8s cubic-bezier(0.4, 0, 0.2, 1), width 0.8s cubic-bezier(0.4, 0, 0.2, 1);
}

.sidebar:not(.collapsed) .menu-item:hover:not(.is-active) {
  background: rgba(255, 255, 255, 0.1) !important;
  transform: translateX(5px);
}

.sidebar.collapsed .menu-item:hover:not(.is-active) {
  background: rgba(255, 255, 255, 0.1) !important;
  transform: scale(1.1);
}

.sidebar:not(.collapsed) .menu-item.is-active {
  background: linear-gradient(135deg, #3498db, #2980b9) !important;
  box-shadow: 0 4px 15px rgba(52, 152, 219, 0.3);
  transition: width 0.8s cubic-bezier(0.4, 0, 0.2, 1), margin 0.8s cubic-bezier(0.4, 0, 0.2, 1), padding 0.8s cubic-bezier(0.4, 0, 0.2, 1), background 0.3s ease, box-shadow 0.3s ease;
}

.sidebar.collapsed .menu-item.is-active {
  background: linear-gradient(135deg, #3498db, #2980b9) !important;
  box-shadow: 0 4px 15px rgba(52, 152, 219, 0.3);
  transition: width 0.8s cubic-bezier(0.4, 0, 0.2, 1), margin 0.8s cubic-bezier(0.4, 0, 0.2, 1), padding 0.8s cubic-bezier(0.4, 0, 0.2, 1), background 0.3s ease, box-shadow 0.3s ease;
}

.sidebar.collapsed .menu-item.is-active:hover {
  background: linear-gradient(135deg, #2980b9, #1f5f8b) !important;
  transform: scale(1.05);
}

.menu-item.is-active::before {
  content: '';
  position: absolute;
  left: -16px;
  top: 50%;
  transform: translateY(-50%);
  width: 4px;
  height: 60%;
  background: #ffffff;
  border-radius: 0 2px 2px 0;
  opacity: 1;
  transition: opacity 0.2s ease, transform 0.3s ease;
}

.sidebar.collapsed .menu-item.is-active::before {
  opacity: 0;
  transform: translateY(-50%) translateX(-10px);
}

.sidebar:not(.collapsed) .menu-item.is-active::before {
  opacity: 1;
  transform: translateY(-50%) translateX(0);
}

.sidebar-footer {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 20px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
  background: rgba(0, 0, 0, 0.1);
}

.collapse-btn {
  position: fixed;
  top: 50%;
  left: calc(280px - 15px);
  transform: translateY(-50%);
  display: flex !important;
  align-items: center;
  justify-content: center;
  width: 30px;
  height: 30px;
  background: rgba(255, 255, 255, 0.95);
  border: 2px solid #3498db;
  border-radius: 50%;
  cursor: pointer;
  transition: left 0.8s cubic-bezier(0.4, 0, 0.2, 1), background 0.5s ease, transform 0.5s ease, box-shadow 0.5s ease;
  color: #3498db;
  z-index: 1000;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
  opacity: 1;
  visibility: visible;
}

.collapse-btn:hover {
  background: #3498db;
  color: #ffffff;
  transform: translateY(-50%) scale(1.1);
  box-shadow: 0 4px 12px rgba(52, 152, 219, 0.3);
}

/* 收起状态下的按钮样式 */
.sidebar.collapsed .collapse-btn {
  left: calc(80px - 15px);
}

.sidebar.collapsed .collapse-btn:hover {
  background: #3498db;
  color: #ffffff;
  transform: translateY(-50%) scale(1.1);
}

.system-status {
  display: flex;
  flex-direction: column;
  gap: 8px;
  height: 44px; /* 减少高度，去掉多余空白 */
  transition: opacity 0.3s ease;
}

.status-item {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #a0a8b8;
  font-size: 12px;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.3);
}

.status-dot.online {
  background: #27ae60;
  box-shadow: 0 0 10px rgba(39, 174, 96, 0.5);
  animation: pulse 2s infinite;
}

.status-dot.offline {
  background: #e74c3c;
  box-shadow: 0 0 10px rgba(231, 76, 60, 0.5);
}

@keyframes pulse {
  0% { box-shadow: 0 0 10px rgba(39, 174, 96, 0.5); }
  50% { box-shadow: 0 0 20px rgba(39, 174, 96, 0.8); }
  100% { box-shadow: 0 0 10px rgba(39, 174, 96, 0.5); }
}

.header {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
  padding: 0 30px;
  display: flex;
  align-items: center;
  box-shadow: 0 2px 20px rgba(0, 0, 0, 0.05);
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.header-left h3 {
  margin: 0 0 5px 0;
  color: #2c3e50;
  font-size: 24px;
  font-weight: 600;
}

.header-subtitle {
  margin: 0;
  color: #7f8c8d;
  font-size: 14px;
  font-weight: 400;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 15px;
}

.refresh-btn {
  background: linear-gradient(135deg, #3498db, #2980b9);
  border: none;
  border-radius: 8px;
  padding: 10px 20px;
  font-weight: 500;
  transition: all 0.3s ease;
  box-shadow: 0 4px 15px rgba(52, 152, 219, 0.2);
}

.refresh-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(52, 152, 219, 0.3);
}

.user-avatar {
  background: linear-gradient(135deg, #e74c3c, #c0392b);
  cursor: pointer;
  transition: all 0.3s ease;
}

.user-avatar:hover {
  transform: scale(1.1);
  box-shadow: 0 4px 15px rgba(231, 76, 60, 0.3);
}

.main-content {
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
  padding: 0;
  position: relative;
  overflow: hidden;
}

.content-wrapper {
  padding: 30px;
  min-height: 100%;
  position: relative;
  z-index: 1;
}

.main-content::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: url('data:image/svg+xml,<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><defs><pattern id="dots" width="20" height="20" patternUnits="userSpaceOnUse"><circle cx="10" cy="10" r="1" fill="rgba(0,0,0,0.02)"/></pattern></defs><rect width="100" height="100" fill="url(%23dots)"/></svg>');
  pointer-events: none;
}
</style>

<style>
* {
  box-sizing: border-box;
}

body {
  margin: 0;
  font-family: 'Inter', 'Helvetica Neue', Helvetica, 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', '微软雅黑', Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  background: #f8f9fa;
}

#app {
  height: 100vh;
  overflow: hidden;
}

/* 自定义滚动条 */
::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}

::-webkit-scrollbar-track {
  background: rgba(0, 0, 0, 0.05);
  border-radius: 3px;
}

::-webkit-scrollbar-thumb {
  background: rgba(0, 0, 0, 0.2);
  border-radius: 3px;
  transition: all 0.3s ease;
}

::-webkit-scrollbar-thumb:hover {
  background: rgba(0, 0, 0, 0.3);
}

/* Element Plus 组件样式覆盖 */
.el-menu-item {
  font-weight: 500;
  letter-spacing: 0.5px;
  transition: all 0.8s cubic-bezier(0.4, 0, 0.2, 1);
}

.sidebar:not(.collapsed) .el-menu-item {
  width: calc(100% - 32px);
  padding: 0 20px;
  margin: 4px 16px;
  transition: all 0.8s cubic-bezier(0.4, 0, 0.2, 1);
}

.el-menu-item i {
  margin-right: 8px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .sidebar {
    width: 200px !important;
  }
  
  .header-content {
    padding: 0 15px;
  }
  
  .header-left h3 {
    font-size: 20px;
  }
  
  .header-subtitle {
    display: none;
  }
  
  .content-wrapper {
    padding: 20px 15px;
  }
}

@media (max-width: 480px) {
  .sidebar {
    width: 100% !important;
    position: fixed;
    z-index: 1000;
    transform: translateX(-100%);
    transition: transform 0.3s ease;
  }
  
  .sidebar.mobile-open {
    transform: translateX(0);
  }
}

/* 加载动画 */
@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.fade-in-up {
  animation: fadeInUp 0.6s ease-out;
}

/* 卡片悬停效果 */
.hover-card {
  transition: all 0.3s ease;
  cursor: pointer;
}

.hover-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
}
</style>