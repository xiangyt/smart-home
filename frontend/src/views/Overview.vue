<template>
  <div class="overview-container">
    <!-- 状态卡片区域 -->
    <div class="status-cards">
      <el-row :gutter="20">
        <el-col :span="6">
          <el-card class="status-card">
            <div class="card-content">
              <div class="card-icon">
                <el-icon size="24" color="#67C23A"><House /></el-icon>
              </div>
              <div class="card-info">
                <div class="card-title">在线设备</div>
                <div class="card-value">{{ onlineDevicesCount }}</div>
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card class="status-card">
            <div class="card-content">
              <div class="card-icon">
                <el-icon size="24" color="#E6A23C"><Sunny /></el-icon>
              </div>
              <div class="card-info">
                <div class="card-title">灯光设备</div>
                <div class="card-value">{{ lightDevicesCount }}</div>
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card class="status-card">
            <div class="card-content">
              <div class="card-icon">
                <el-icon size="24" color="#409EFF"><Monitor /></el-icon>
              </div>
              <div class="card-info">
                <div class="card-title">传感器</div>
                <div class="card-value">{{ sensorDevicesCount }}</div>
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card class="status-card">
            <div class="card-content">
              <div class="card-icon">
                <el-icon size="24" color="#F56C6C"><Switch /></el-icon>
              </div>
              <div class="card-info">
                <div class="card-title">开关设备</div>
                <div class="card-value">{{ switchDevicesCount }}</div>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- 设备控制区域 -->
    <div class="device-sections">
      <el-row :gutter="20">
        <!-- 灯光控制 -->
        <el-col :span="12">
          <el-card class="device-section">
            <template #header>
              <div class="section-header">
                <el-icon><Sunny /></el-icon>
                <span>灯光控制</span>
              </div>
            </template>
            <div class="device-grid">
              <div 
                v-for="light in lightDevices" 
                :key="light.entity_id"
                class="device-item"
                :class="{ active: light.state === 'on' }"
                @click="toggleDevice(light)"
              >
                <div class="device-icon">
                  <el-icon><Sunny /></el-icon>
                </div>
                <div class="device-name">{{ light.friendly_name }}</div>
                <div class="device-state">{{ light.state === 'on' ? '开启' : '关闭' }}</div>
              </div>
            </div>
          </el-card>
        </el-col>

        <!-- 开关控制 -->
        <el-col :span="12">
          <el-card class="device-section">
            <template #header>
              <div class="section-header">
                <el-icon><Switch /></el-icon>
                <span>开关控制</span>
              </div>
            </template>
            <div class="device-grid">
              <div 
                v-for="switchDevice in switchDevices" 
                :key="switchDevice.entity_id"
                class="device-item"
                :class="{ active: switchDevice.state === 'on' }"
                @click="toggleDevice(switchDevice)"
              >
                <div class="device-icon">
                  <el-icon><Switch /></el-icon>
                </div>
                <div class="device-name">{{ switchDevice.friendly_name }}</div>
                <div class="device-state">{{ switchDevice.state === 'on' ? '开启' : '关闭' }}</div>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>

      <!-- 传感器数据 -->
      <el-row :gutter="20" style="margin-top: 20px;">
        <el-col :span="24">
          <el-card class="device-section">
            <template #header>
              <div class="section-header">
                <el-icon><Monitor /></el-icon>
                <span>传感器数据</span>
              </div>
            </template>
            <div class="sensor-grid">
              <div 
                v-for="sensor in sensorDevices" 
                :key="sensor.entity_id"
                class="sensor-item"
              >
                <div class="sensor-icon">
                  <el-icon><Monitor /></el-icon>
                </div>
                <div class="sensor-info">
                  <div class="sensor-name">{{ sensor.friendly_name }}</div>
                  <div class="sensor-value">
                    {{ sensor.state }}
                    <span v-if="sensor.unit_of_measurement">{{ sensor.unit_of_measurement }}</span>
                  </div>
                </div>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useDevicesStore } from '@/stores/devices'
import { House, Sunny, Monitor, Switch } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const devicesStore = useDevicesStore()

// 计算属性
const onlineDevicesCount = computed(() => devicesStore.getOnlineDevicesCount())
const lightDevices = computed(() => devicesStore.getDevicesByDomain('light'))
const switchDevices = computed(() => devicesStore.getDevicesByDomain('switch'))
const sensorDevices = computed(() => devicesStore.getDevicesByDomain('sensor'))

const lightDevicesCount = computed(() => lightDevices.value.length)
const switchDevicesCount = computed(() => switchDevices.value.length)
const sensorDevicesCount = computed(() => sensorDevices.value.length)

// 切换设备状态
const toggleDevice = async (device: any) => {
  try {
    const service = device.state === 'on' ? 'turn_off' : 'turn_on'
    await devicesStore.controlDevice(device.entity_id, service)
    ElMessage.success(`${device.friendly_name} ${service === 'turn_on' ? '已开启' : '已关闭'}`)
  } catch (error) {
    ElMessage.error('设备控制失败')
  }
}

// 页面加载时获取设备数据
onMounted(() => {
  devicesStore.fetchDevices()
})
</script>

<style scoped>
.overview-container {
  padding: 0;
}

.status-cards {
  margin-bottom: 20px;
}

.status-card {
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.card-content {
  display: flex;
  align-items: center;
  padding: 10px 0;
}

.card-icon {
  margin-right: 15px;
}

.card-info {
  flex: 1;
}

.card-title {
  font-size: 14px;
  color: #909399;
  margin-bottom: 5px;
}

.card-value {
  font-size: 24px;
  font-weight: bold;
  color: #303133;
}

.device-section {
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.section-header {
  display: flex;
  align-items: center;
  font-weight: bold;
  color: #303133;
}

.section-header .el-icon {
  margin-right: 8px;
}

.device-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
  gap: 15px;
}

.device-item {
  padding: 15px;
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s;
  background: #fff;
}

.device-item:hover {
  border-color: #409EFF;
  box-shadow: 0 2px 8px rgba(64, 158, 255, 0.2);
}

.device-item.active {
  background: #ecf5ff;
  border-color: #409EFF;
}

.device-icon {
  margin-bottom: 10px;
  color: #606266;
}

.device-item.active .device-icon {
  color: #409EFF;
}

.device-name {
  font-size: 14px;
  color: #303133;
  margin-bottom: 5px;
  font-weight: 500;
}

.device-state {
  font-size: 12px;
  color: #909399;
}

.device-item.active .device-state {
  color: #409EFF;
}

.sensor-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 15px;
}

.sensor-item {
  display: flex;
  align-items: center;
  padding: 15px;
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  background: #fff;
}

.sensor-icon {
  margin-right: 15px;
  color: #606266;
}

.sensor-info {
  flex: 1;
}

.sensor-name {
  font-size: 14px;
  color: #303133;
  margin-bottom: 5px;
  font-weight: 500;
}

.sensor-value {
  font-size: 16px;
  color: #409EFF;
  font-weight: bold;
}
</style>