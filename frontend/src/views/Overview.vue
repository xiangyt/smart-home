<template>
  <div class="overview-container">
    <!-- 欢迎横幅 -->
    <div class="welcome-banner">
      <div class="banner-content">
        <div class="banner-text">
          <h1>欢迎回来！</h1>
          <p>您的智能家居系统运行正常，共有 <strong>{{ totalDevicesCount }}</strong> 个设备在线</p>
        </div>
        <div class="banner-stats">
          <div class="quick-stat">
            <div class="stat-number">{{ onlineDevicesCount }}</div>
            <div class="stat-label">在线设备</div>
          </div>
          <div class="quick-stat">
            <div class="stat-number">{{ Math.round((onlineDevicesCount / totalDevicesCount) * 100) }}%</div>
            <div class="stat-label">在线率</div>
          </div>
        </div>
      </div>
    </div>

    <!-- 状态卡片区域 -->
    <div class="status-cards">
      <el-row :gutter="24">
        <el-col :xs="24" :sm="12" :md="6">
          <div class="status-card online-devices" @click="filterDevices('online')">
            <div class="card-background">
              <div class="card-icon">
                <el-icon size="32"><House /></el-icon>
              </div>
            </div>
            <div class="card-content">
              <div class="card-header">
                <div class="card-title">在线设备</div>
                <div class="card-trend up">
                  <el-icon size="12"><ArrowUp /></el-icon>
                  +2.5%
                </div>
              </div>
              <div class="card-value">{{ onlineDevicesCount }}</div>
              <div class="card-subtitle">设备正常运行</div>
            </div>
          </div>
        </el-col>
        
        <el-col :xs="24" :sm="12" :md="6">
          <div class="status-card light-devices" @click="filterDevices('light')">
            <div class="card-background">
              <div class="card-icon">
                <el-icon size="32"><Sunny /></el-icon>
              </div>
            </div>
            <div class="card-content">
              <div class="card-header">
                <div class="card-title">灯光设备</div>
                <div class="card-trend up">
                  <el-icon size="12"><ArrowUp /></el-icon>
                  +1.2%
                </div>
              </div>
              <div class="card-value">{{ lightDevicesCount }}</div>
              <div class="card-subtitle">照明系统</div>
            </div>
          </div>
        </el-col>
        
        <el-col :xs="24" :sm="12" :md="6">
          <div class="status-card sensor-devices" @click="filterDevices('sensor')">
            <div class="card-background">
              <div class="card-icon">
                <el-icon size="32"><Monitor /></el-icon>
              </div>
            </div>
            <div class="card-content">
              <div class="card-header">
                <div class="card-title">传感器</div>
                <div class="card-trend stable">
                  <el-icon size="12"><Minus /></el-icon>
                  0%
                </div>
              </div>
              <div class="card-value">{{ sensorDevicesCount }}</div>
              <div class="card-subtitle">环境监测</div>
            </div>
          </div>
        </el-col>
        
        <el-col :xs="24" :sm="12" :md="6">
          <div class="status-card switch-devices" @click="filterDevices('switch')">
            <div class="card-background">
              <div class="card-icon">
                <el-icon size="32"><Switch /></el-icon>
              </div>
            </div>
            <div class="card-content">
              <div class="card-header">
                <div class="card-title">开关设备</div>
                <div class="card-trend down">
                  <el-icon size="12"><ArrowDown /></el-icon>
                  -0.8%
                </div>
              </div>
              <div class="card-value">{{ switchDevicesCount }}</div>
              <div class="card-subtitle">控制开关</div>
            </div>
          </div>
        </el-col>
      </el-row>
    </div>

    <!-- 主要内容区域 -->
    <el-row :gutter="24" class="main-content-row">
      <!-- 设备列表 -->
      <el-col :xs="24" :lg="16">
        <div class="content-card device-list-card">
          <div class="card-header">
            <h3>设备状态</h3>
            <div class="header-actions">
              <el-select v-model="selectedFilter" placeholder="筛选设备" size="small" style="width: 120px">
                <el-option label="全部" value="all" />
                <el-option label="在线" value="online" />
                <el-option label="离线" value="offline" />
                <el-option label="灯光" value="light" />
                <el-option label="传感器" value="sensor" />
              </el-select>
              <el-button size="small" type="primary" @click="refreshDevices">
                <el-icon><Refresh /></el-icon>
              </el-button>
            </div>
          </div>
          
          <div class="device-grid">
            <div 
              v-for="device in filteredDevices.slice(0, 12)" 
              :key="device.entity_id"
              class="device-item"
              @click="toggleDevice(device)"
            >
              <div class="device-status" :class="{ 'online': device.state !== 'unavailable' }"></div>
              <div class="device-icon">
                <el-icon size="20">
                  <Sunny v-if="device.entity_id.includes('light')" />
                  <Monitor v-else-if="device.entity_id.includes('sensor')" />
                  <Switch v-else-if="device.entity_id.includes('switch')" />
                  <House v-else />
                </el-icon>
              </div>
              <div class="device-info">
                <div class="device-name">{{ device.friendly_name || device.entity_id }}</div>
                <div class="device-state">{{ device.state }}</div>
              </div>
            </div>
          </div>
          
          <div class="view-more" v-if="filteredDevices.length > 12">
            <el-button type="text" @click="$router.push('/devices')">
              查看全部 {{ filteredDevices.length }} 个设备
              <el-icon><ArrowRight /></el-icon>
            </el-button>
          </div>
        </div>
      </el-col>

      <!-- 系统信息 -->
      <el-col :xs="24" :lg="8">
        <div class="content-card system-info-card">
          <div class="card-header">
            <h3>系统信息</h3>
          </div>
          
          <div class="system-stats">
            <div class="stat-item">
              <div class="stat-icon">
                <el-icon size="24" color="#67C23A"><CircleCheck /></el-icon>
              </div>
              <div class="stat-content">
                <div class="stat-title">系统状态</div>
                <div class="stat-value">正常运行</div>
              </div>
            </div>
            
            <div class="stat-item">
              <div class="stat-icon">
                <el-icon size="24" color="#409EFF"><Clock /></el-icon>
              </div>
              <div class="stat-content">
                <div class="stat-title">运行时间</div>
                <div class="stat-value">7天 12小时</div>
              </div>
            </div>
            
            <div class="stat-item">
              <div class="stat-icon">
                <el-icon size="24" color="#E6A23C"><Warning /></el-icon>
              </div>
              <div class="stat-content">
                <div class="stat-title">待处理事件</div>
                <div class="stat-value">3个</div>
              </div>
            </div>
          </div>
          
          <div class="quick-actions">
            <h4>快速操作</h4>
            <div class="action-buttons">
              <el-button type="primary" size="small" class="action-btn">
                <el-icon><Sunny /></el-icon>
                全部开灯
              </el-button>
              <el-button type="info" size="small" class="action-btn">
                <el-icon><Moon /></el-icon>
                夜间模式
              </el-button>
              <el-button type="success" size="small" class="action-btn">
                <el-icon><House /></el-icon>
                回家模式
              </el-button>
              <el-button type="warning" size="small" class="action-btn">
                <el-icon><Lock /></el-icon>
                离家模式
              </el-button>
            </div>
          </div>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted } from 'vue'
import { useDevicesStore } from '@/stores/devices'
import { useRouter } from 'vue-router'
import { 
  House, Sunny, Monitor, Switch, ArrowUp, ArrowDown, Minus, 
  Refresh, ArrowRight, CircleCheck, Clock, Warning, Moon, Lock 
} from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const devicesStore = useDevicesStore()
const router = useRouter()

// 响应式数据
const selectedFilter = ref('all')

// 计算属性
const totalDevicesCount = computed(() => devicesStore.devices.length)
const onlineDevicesCount = computed(() => devicesStore.getOnlineDevicesCount())
const lightDevices = computed(() => devicesStore.getDevicesByDomain('light'))
const switchDevices = computed(() => devicesStore.getDevicesByDomain('switch'))
const sensorDevices = computed(() => devicesStore.getDevicesByDomain('sensor'))

const lightDevicesCount = computed(() => lightDevices.value.length)
const switchDevicesCount = computed(() => switchDevices.value.length)
const sensorDevicesCount = computed(() => sensorDevices.value.length)

// 筛选设备
const filteredDevices = computed(() => {
  if (!Array.isArray(devicesStore.devices)) return []
  
  switch (selectedFilter.value) {
    case 'online':
      return devicesStore.devices.filter(device => device.state !== 'unavailable')
    case 'offline':
      return devicesStore.devices.filter(device => device.state === 'unavailable')
    case 'light':
      return lightDevices.value
    case 'sensor':
      return sensorDevices.value
    default:
      return devicesStore.devices
  }
})

// 方法
const filterDevices = (type: string) => {
  selectedFilter.value = type
}

const refreshDevices = async () => {
  try {
    await devicesStore.fetchDevices()
    ElMessage.success('设备数据刷新成功')
  } catch (error) {
    ElMessage.error('刷新失败，请稍后重试')
  }
}

const toggleDevice = async (device: any) => {
  try {
    const service = device.state === 'on' ? 'turn_off' : 'turn_on'
    await devicesStore.controlDevice(device.entity_id, service)
    ElMessage.success(`${device.friendly_name || device.entity_id} ${service === 'turn_on' ? '已开启' : '已关闭'}`)
  } catch (error) {
    ElMessage.error('设备控制失败')
  }
}

// 生命周期
onMounted(() => {
  devicesStore.fetchDevices()
})
</script>

<style scoped>
.overview-container {
  padding: 0;
  animation: fadeInUp 0.6s ease-out;
}

/* 欢迎横幅 */
.welcome-banner {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 16px;
  padding: 30px;
  margin-bottom: 30px;
  color: white;
  position: relative;
  overflow: hidden;
}

.welcome-banner::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: url('data:image/svg+xml,<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><defs><pattern id="dots" width="20" height="20" patternUnits="userSpaceOnUse"><circle cx="10" cy="10" r="1" fill="rgba(255,255,255,0.1)"/></pattern></defs><rect width="100" height="100" fill="url(%23dots)"/></svg>');
  pointer-events: none;
}

.banner-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  position: relative;
  z-index: 1;
}

.banner-text h1 {
  margin: 0 0 10px 0;
  font-size: 28px;
  font-weight: 600;
}

.banner-text p {
  margin: 0;
  font-size: 16px;
  opacity: 0.9;
}

.banner-stats {
  display: flex;
  gap: 30px;
}

.quick-stat {
  text-align: center;
}

.stat-number {
  font-size: 32px;
  font-weight: 700;
  line-height: 1;
}

.stat-label {
  font-size: 14px;
  opacity: 0.8;
  margin-top: 5px;
}

/* 状态卡片 */
.status-cards {
  margin-bottom: 30px;
}

.status-card {
  background: white;
  border-radius: 16px;
  padding: 24px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  cursor: pointer;
  transition: all 0.3s ease;
  position: relative;
  overflow: hidden;
  margin-bottom: 20px;
}

.status-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.12);
}

.card-background {
  position: absolute;
  top: -20px;
  right: -20px;
  width: 80px;
  height: 80px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0.1;
}

.online-devices .card-background {
  background: #67C23A;
}

.light-devices .card-background {
  background: #E6A23C;
}

.sensor-devices .card-background {
  background: #409EFF;
}

.switch-devices .card-background {
  background: #F56C6C;
}

.card-background .card-icon {
  color: white;
}

.card-content {
  position: relative;
  z-index: 1;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.card-title {
  font-size: 14px;
  color: #909399;
  font-weight: 500;
}

.card-trend {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  font-weight: 600;
  padding: 4px 8px;
  border-radius: 12px;
}

.card-trend.up {
  color: #67C23A;
  background: rgba(103, 194, 58, 0.1);
}

.card-trend.down {
  color: #F56C6C;
  background: rgba(245, 108, 108, 0.1);
}

.card-trend.stable {
  color: #909399;
  background: rgba(144, 147, 153, 0.1);
}

.card-value {
  font-size: 32px;
  font-weight: 700;
  color: #303133;
  line-height: 1;
  margin-bottom: 8px;
}

.card-subtitle {
  font-size: 12px;
  color: #909399;
}

/* 主要内容区域 */
.main-content-row {
  margin-bottom: 30px;
}

.content-card {
  background: white;
  border-radius: 16px;
  padding: 24px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 15px;
  border-bottom: 1px solid #f0f0f0;
}

.card-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}

.header-actions {
  display: flex;
  gap: 10px;
  align-items: center;
}

/* 设备网格 */
.device-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 16px;
  margin-bottom: 20px;
}

.device-item {
  background: #f8f9fa;
  border: 2px solid transparent;
  border-radius: 12px;
  padding: 16px;
  cursor: pointer;
  transition: all 0.3s ease;
  position: relative;
}

.device-item:hover {
  background: white;
  border-color: #409EFF;
  box-shadow: 0 4px 15px rgba(64, 158, 255, 0.2);
}

.device-status {
  position: absolute;
  top: 12px;
  right: 12px;
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #ddd;
}

.device-status.online {
  background: #67C23A;
  box-shadow: 0 0 10px rgba(103, 194, 58, 0.5);
}

.device-icon {
  margin-bottom: 12px;
  color: #606266;
}

.device-info {
  text-align: left;
}

.device-name {
  font-size: 14px;
  font-weight: 500;
  color: #303133;
  margin-bottom: 4px;
  line-height: 1.4;
}

.device-state {
  font-size: 12px;
  color: #909399;
}

.view-more {
  text-align: center;
  padding-top: 15px;
  border-top: 1px solid #f0f0f0;
}

/* 系统信息 */
.system-stats {
  margin-bottom: 25px;
}

.stat-item {
  display: flex;
  align-items: center;
  padding: 15px 0;
  border-bottom: 1px solid #f0f0f0;
}

.stat-item:last-child {
  border-bottom: none;
}

.stat-icon {
  margin-right: 15px;
}

.stat-content {
  flex: 1;
}

.stat-title {
  font-size: 14px;
  color: #909399;
  margin-bottom: 4px;
}

.stat-value {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

/* 快速操作 */
.quick-actions h4 {
  margin: 0 0 15px 0;
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.action-buttons {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
}

.action-btn {
  border-radius: 8px;
  font-weight: 500;
  transition: all 0.3s ease;
}

.action-btn:hover {
  transform: translateY(-2px);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .banner-content {
    flex-direction: column;
    text-align: center;
    gap: 20px;
  }
  
  .banner-stats {
    justify-content: center;
  }
  
  .device-grid {
    grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
  }
  
  .action-buttons {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 480px) {
  .welcome-banner {
    padding: 20px;
  }
  
  .banner-text h1 {
    font-size: 24px;
  }
  
  .content-card {
    padding: 16px;
  }
  
  .device-grid {
    grid-template-columns: 1fr;
  }
}
</style>