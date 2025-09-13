<template>
  <div class="devices-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>设备管理</span>
          <el-button type="primary" @click="refreshDevices">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </div>
      </template>
      
      <el-table :data="devices" v-loading="loading" style="width: 100%">
        <el-table-column prop="friendly_name" label="设备名称" width="200" />
        <el-table-column prop="entity_id" label="实体ID" width="250" />
        <el-table-column prop="domain" label="类型" width="100">
          <template #default="scope">
            <el-tag :type="getDomainTagType(scope.row.domain)">
              {{ getDomainLabel(scope.row.domain) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="state" label="状态" width="120">
          <template #default="scope">
            <el-tag :type="getStateTagType(scope.row.state)">
              {{ getStateLabel(scope.row.state) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="last_updated" label="最后更新" width="180">
          <template #default="scope">
            {{ formatTime(scope.row.last_updated) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200">
          <template #default="scope">
            <el-button 
              v-if="canControl(scope.row)"
              size="small" 
              type="primary"
              @click="toggleDevice(scope.row)"
            >
              {{ scope.row.state === 'on' ? '关闭' : '开启' }}
            </el-button>
            <el-button size="small" @click="viewDetails(scope.row)">
              详情
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 设备详情对话框 -->
    <el-dialog v-model="detailsVisible" title="设备详情" width="600px">
      <div v-if="selectedDevice">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="设备名称">
            {{ selectedDevice.friendly_name }}
          </el-descriptions-item>
          <el-descriptions-item label="实体ID">
            {{ selectedDevice.entity_id }}
          </el-descriptions-item>
          <el-descriptions-item label="域名">
            {{ selectedDevice.domain }}
          </el-descriptions-item>
          <el-descriptions-item label="状态">
            {{ selectedDevice.state }}
          </el-descriptions-item>
          <el-descriptions-item label="设备类别" v-if="selectedDevice.device_class">
            {{ selectedDevice.device_class }}
          </el-descriptions-item>
          <el-descriptions-item label="单位" v-if="selectedDevice.unit_of_measurement">
            {{ selectedDevice.unit_of_measurement }}
          </el-descriptions-item>
          <el-descriptions-item label="最后更新">
            {{ formatTime(selectedDevice.last_updated) }}
          </el-descriptions-item>
          <el-descriptions-item label="最后改变">
            {{ formatTime(selectedDevice.last_changed) }}
          </el-descriptions-item>
        </el-descriptions>
        
        <div style="margin-top: 20px;">
          <h4>属性信息</h4>
          <el-table :data="attributesList" size="small">
            <el-table-column prop="key" label="属性" width="200" />
            <el-table-column prop="value" label="值" />
          </el-table>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useDevicesStore, type Device } from '@/stores/devices'
import { Refresh } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const devicesStore = useDevicesStore()
const detailsVisible = ref(false)
const selectedDevice = ref<Device | null>(null)

const devices = computed(() => devicesStore.devices)
const loading = computed(() => devicesStore.loading)

// 刷新设备列表
const refreshDevices = () => {
  devicesStore.fetchDevices()
}

// 获取域名标签类型
const getDomainTagType = (domain: string) => {
  const typeMap: Record<string, string> = {
    light: 'warning',
    switch: 'success',
    sensor: 'info',
    binary_sensor: 'info',
    climate: 'primary'
  }
  return typeMap[domain] || ''
}

// 获取域名标签文本
const getDomainLabel = (domain: string) => {
  const labelMap: Record<string, string> = {
    light: '灯光',
    switch: '开关',
    sensor: '传感器',
    binary_sensor: '二进制传感器',
    climate: '空调'
  }
  return labelMap[domain] || domain
}

// 获取状态标签类型
const getStateTagType = (state: string) => {
  if (state === 'on') return 'success'
  if (state === 'off') return 'info'
  if (state === 'unavailable') return 'danger'
  return ''
}

// 获取状态标签文本
const getStateLabel = (state: string) => {
  const labelMap: Record<string, string> = {
    on: '开启',
    off: '关闭',
    unavailable: '不可用',
    unknown: '未知'
  }
  return labelMap[state] || state
}

// 判断设备是否可控制
const canControl = (device: Device) => {
  return ['light', 'switch'].includes(device.domain)
}

// 切换设备状态
const toggleDevice = async (device: Device) => {
  try {
    const service = device.state === 'on' ? 'turn_off' : 'turn_on'
    await devicesStore.controlDevice(device.entity_id, service)
    ElMessage.success(`${device.friendly_name} ${service === 'turn_on' ? '已开启' : '已关闭'}`)
  } catch (error) {
    ElMessage.error('设备控制失败')
  }
}

// 查看设备详情
const viewDetails = (device: Device) => {
  selectedDevice.value = device
  detailsVisible.value = true
}

// 格式化时间
const formatTime = (timeStr: string) => {
  return new Date(timeStr).toLocaleString('zh-CN')
}

// 属性列表
const attributesList = computed(() => {
  if (!selectedDevice.value?.attributes) return []
  return Object.entries(selectedDevice.value.attributes).map(([key, value]) => ({
    key,
    value: typeof value === 'object' ? JSON.stringify(value) : String(value)
  }))
})

onMounted(() => {
  devicesStore.fetchDevices()
})
</script>

<style scoped>
.devices-container {
  padding: 0;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>