<template>
  <div class="devices-container">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <div class="header-info">
          <h1>设备管理</h1>
          <p>管理和控制您的所有智能设备</p>
        </div>
        <div class="header-actions">
          <el-button type="primary" @click="refreshDevices">
            <el-icon><Refresh /></el-icon>
            刷新设备
          </el-button>
          <el-button type="success" @click="addDevice">
            <el-icon><Plus /></el-icon>
            添加设备
          </el-button>
        </div>
      </div>
    </div>

    <!-- 筛选和搜索 -->
    <div class="filter-section">
      <el-row :gutter="20">
        <el-col :span="8">
          <el-input
            v-model="searchQuery"
            placeholder="搜索设备名称..."
            clearable
            @input="handleSearch"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </el-col>
        <el-col :span="4">
          <el-select v-model="selectedDomain" placeholder="设备类型" clearable @change="handleFilter">
            <el-option label="全部类型" value="" />
            <el-option label="灯光" value="light" />
            <el-option label="开关" value="switch" />
            <el-option label="传感器" value="sensor" />
            <el-option label="空调" value="climate" />
            <el-option label="媒体播放器" value="media_player" />
          </el-select>
        </el-col>
        <el-col :span="4">
          <el-select v-model="selectedStatus" placeholder="设备状态" clearable @change="handleFilter">
            <el-option label="全部状态" value="" />
            <el-option label="在线" value="online" />
            <el-option label="离线" value="offline" />
            <el-option label="开启" value="on" />
            <el-option label="关闭" value="off" />
          </el-select>
        </el-col>
        <el-col :span="8">
          <div class="view-controls">
            <el-button-group>
              <el-button :type="viewMode === 'grid' ? 'primary' : 'default'" @click="viewMode = 'grid'">
                <el-icon><Grid /></el-icon>
                网格视图
              </el-button>
              <el-button :type="viewMode === 'list' ? 'primary' : 'default'" @click="viewMode = 'list'">
                <el-icon><List /></el-icon>
                列表视图
              </el-button>
            </el-button-group>
          </div>
        </el-col>
      </el-row>
    </div>

    <!-- 设备统计 -->
    <div class="stats-section">
      <el-row :gutter="20">
        <el-col :span="6">
          <div class="stat-card">
            <div class="stat-icon online">
              <el-icon size="24"><CircleCheck /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-number">{{ onlineDevicesCount }}</div>
              <div class="stat-label">在线设备</div>
            </div>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="stat-card">
            <div class="stat-icon offline">
              <el-icon size="24"><CircleClose /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-number">{{ offlineDevicesCount }}</div>
              <div class="stat-label">离线设备</div>
            </div>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="stat-card">
            <div class="stat-icon active">
              <el-icon size="24"><Switch /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-number">{{ activeDevicesCount }}</div>
              <div class="stat-label">活跃设备</div>
            </div>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="stat-card">
            <div class="stat-icon total">
              <el-icon size="24"><Monitor /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-number">{{ totalDevicesCount }}</div>
              <div class="stat-label">设备总数</div>
            </div>
          </div>
        </el-col>
      </el-row>
    </div>

    <!-- 设备列表/网格 -->
    <div class="devices-section">
      <!-- 网格视图 -->
      <div v-if="viewMode === 'grid'" class="devices-grid">
        <div 
          v-for="device in paginatedDevices" 
          :key="device.entity_id"
          class="device-card"
          @click="showDeviceDetail(device)"
        >
          <div class="device-status" :class="getDeviceStatusClass(device)"></div>
          <div class="device-header">
            <div class="device-icon">
              <el-icon size="28">
                <Sunny v-if="device.entity_id.includes('light')" />
                <Monitor v-else-if="device.entity_id.includes('sensor')" />
                <Switch v-else-if="device.entity_id.includes('switch')" />
                <VideoPlay v-else-if="device.entity_id.includes('media_player')" />
                <House v-else />
              </el-icon>
            </div>
            <div class="device-actions">
              <el-dropdown @command="handleDeviceAction">
                <el-button size="small" circle>
                  <el-icon><MoreFilled /></el-icon>
                </el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item :command="{action: 'toggle', device}">
                      {{ device.state === 'on' ? '关闭' : '开启' }}
                    </el-dropdown-item>
                    <el-dropdown-item :command="{action: 'detail', device}">查看详情</el-dropdown-item>
                    <el-dropdown-item :command="{action: 'edit', device}">编辑设备</el-dropdown-item>
                    <el-dropdown-item :command="{action: 'delete', device}" divided>删除设备</el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>
          </div>
          
          <div class="device-info">
            <div class="device-name">{{ device.friendly_name || device.entity_id }}</div>
            <div class="device-id">{{ device.entity_id }}</div>
            <div class="device-state" :class="getDeviceStateClass(device)">
              {{ formatDeviceState(device) }}
            </div>
          </div>
          
          <div class="device-footer">
            <div class="device-domain">{{ getDomainName(device.entity_id) }}</div>
            <div class="device-last-updated">
              {{ formatLastUpdated(device.last_updated) }}
            </div>
          </div>
        </div>
      </div>

      <!-- 列表视图 -->
      <div v-else class="devices-table">
        <el-table :data="paginatedDevices" style="width: 100%" @row-click="showDeviceDetail" v-loading="loading">
          <el-table-column width="60">
            <template #default="{ row }">
              <div class="table-status" :class="getDeviceStatusClass(row)"></div>
            </template>
          </el-table-column>
          
          <el-table-column label="设备" min-width="200">
            <template #default="{ row }">
              <div class="table-device-info">
                <div class="table-device-icon">
                  <el-icon size="20">
                    <Sunny v-if="row.entity_id.includes('light')" />
                    <Monitor v-else-if="row.entity_id.includes('sensor')" />
                    <Switch v-else-if="row.entity_id.includes('switch')" />
                    <VideoPlay v-else-if="row.entity_id.includes('media_player')" />
                    <House v-else />
                  </el-icon>
                </div>
                <div>
                  <div class="table-device-name">{{ row.friendly_name || row.entity_id }}</div>
                  <div class="table-device-id">{{ row.entity_id }}</div>
                </div>
              </div>
            </template>
          </el-table-column>
          
          <el-table-column label="类型" width="100">
            <template #default="{ row }">
              <el-tag size="small">{{ getDomainName(row.entity_id) }}</el-tag>
            </template>
          </el-table-column>
          
          <el-table-column label="状态" width="120">
            <template #default="{ row }">
              <el-tag :type="getDeviceStateType(row)" size="small">
                {{ formatDeviceState(row) }}
              </el-tag>
            </template>
          </el-table-column>
          
          <el-table-column label="最后更新" width="150">
            <template #default="{ row }">
              {{ formatLastUpdated(row.last_updated) }}
            </template>
          </el-table-column>
          
          <el-table-column label="操作" width="150">
            <template #default="{ row }">
              <el-button size="small" @click.stop="toggleDevice(row)" v-if="canControl(row)">
                {{ row.state === 'on' ? '关闭' : '开启' }}
              </el-button>
              <el-button size="small" type="info" @click.stop="showDeviceDetail(row)">
                详情
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <!-- 分页 -->
      <div class="pagination-section">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[12, 24, 48, 96]"
          :total="filteredDevices.length"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted } from 'vue'
import { useDevicesStore } from '@/stores/devices'
import { 
  Refresh, Plus, Search, Grid, List, CircleCheck, CircleClose, 
  Switch, Monitor, Sunny, VideoPlay, House, MoreFilled 
} from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'

const devicesStore = useDevicesStore()

// 响应式数据
const loading = ref(false)
const searchQuery = ref('')
const selectedDomain = ref('')
const selectedStatus = ref('')
const viewMode = ref('grid')
const currentPage = ref(1)
const pageSize = ref(24)

// 计算属性
const devices = computed(() => devicesStore.devices || [])
const totalDevicesCount = computed(() => devices.value.length)
const onlineDevicesCount = computed(() => 
  devices.value.filter(device => device.state !== 'unavailable').length
)
const offlineDevicesCount = computed(() => 
  devices.value.filter(device => device.state === 'unavailable').length
)
const activeDevicesCount = computed(() => 
  devices.value.filter(device => device.state === 'on').length
)

// 筛选设备
const filteredDevices = computed(() => {
  let result = devices.value

  // 搜索筛选
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter(device => 
      (device.friendly_name || device.entity_id).toLowerCase().includes(query) ||
      device.entity_id.toLowerCase().includes(query)
    )
  }

  // 域名筛选
  if (selectedDomain.value) {
    result = result.filter(device => 
      device.entity_id.startsWith(selectedDomain.value + '.')
    )
  }

  // 状态筛选
  if (selectedStatus.value) {
    if (selectedStatus.value === 'online') {
      result = result.filter(device => device.state !== 'unavailable')
    } else if (selectedStatus.value === 'offline') {
      result = result.filter(device => device.state === 'unavailable')
    } else {
      result = result.filter(device => device.state === selectedStatus.value)
    }
  }

  return result
})

// 分页设备
const paginatedDevices = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredDevices.value.slice(start, end)
})

// 方法
const refreshDevices = async () => {
  loading.value = true
  try {
    await devicesStore.fetchDevices()
    ElMessage.success('设备数据刷新成功')
  } catch (error) {
    ElMessage.error('刷新失败，请稍后重试')
  } finally {
    loading.value = false
  }
}

const addDevice = () => {
  ElMessage.info('添加设备功能开发中...')
}

const handleSearch = () => {
  currentPage.value = 1
}

const handleFilter = () => {
  currentPage.value = 1
}

const handleSizeChange = (size: number) => {
  pageSize.value = size
  currentPage.value = 1
}

const handleCurrentChange = (page: number) => {
  currentPage.value = page
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

const showDeviceDetail = (device: any) => {
  ElMessageBox.alert(
    `设备ID: ${device.entity_id}\n状态: ${device.state}\n域名: ${getDomainName(device.entity_id)}`,
    device.friendly_name || device.entity_id,
    { confirmButtonText: '确定' }
  )
}

const handleDeviceAction = ({ action, device }: { action: string, device: any }) => {
  switch (action) {
    case 'toggle':
      toggleDevice(device)
      break
    case 'detail':
      showDeviceDetail(device)
      break
    case 'edit':
      ElMessage.info('编辑设备功能开发中...')
      break
    case 'delete':
      ElMessage.info('删除设备功能开发中...')
      break
  }
}

// 工具函数
const getDomainName = (entityId: string) => {
  const domain = entityId.split('.')[0]
  const domainNames: Record<string, string> = {
    light: '灯光',
    switch: '开关',
    sensor: '传感器',
    binary_sensor: '二进制传感器',
    climate: '空调',
    media_player: '媒体播放器',
    cover: '窗帘',
    fan: '风扇',
    lock: '门锁',
    camera: '摄像头'
  }
  return domainNames[domain] || domain
}

const formatDeviceState = (device: any) => {
  const stateMap: Record<string, string> = {
    on: '开启',
    off: '关闭',
    unavailable: '离线',
    unknown: '未知'
  }
  return stateMap[device.state] || device.state
}

const formatLastUpdated = (timestamp: string) => {
  if (!timestamp) return '-'
  const date = new Date(timestamp)
  return date.toLocaleString('zh-CN')
}

const getDeviceStatusClass = (device: any) => {
  if (device.state === 'unavailable') return 'offline'
  if (device.state === 'on') return 'online'
  return 'idle'
}

const getDeviceStateClass = (device: any) => {
  if (device.state === 'on') return 'state-on'
  if (device.state === 'unavailable') return 'state-offline'
  return 'state-off'
}

const getDeviceStateType = (device: any) => {
  if (device.state === 'on') return 'success'
  if (device.state === 'unavailable') return 'danger'
  return 'info'
}

const canControl = (device: any) => {
  const controllableDomains = ['light', 'switch', 'fan', 'cover']
  const domain = device.entity_id.split('.')[0]
  return controllableDomains.includes(domain)
}

// 生命周期
onMounted(() => {
  devicesStore.fetchDevices()
})
</script>

<style scoped>
.devices-container {
  padding: 0;
  animation: fadeInUp 0.6s ease-out;
}

/* 页面头部 */
.page-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 16px;
  padding: 30px;
  margin-bottom: 30px;
  color: white;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-info h1 {
  margin: 0 0 8px 0;
  font-size: 28px;
  font-weight: 600;
}

.header-info p {
  margin: 0;
  font-size: 16px;
  opacity: 0.9;
}

.header-actions {
  display: flex;
  gap: 12px;
}

/* 筛选区域 */
.filter-section {
  background: white;
  border-radius: 16px;
  padding: 24px;
  margin-bottom: 24px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}

.view-controls {
  display: flex;
  justify-content: flex-end;
}

/* 统计区域 */
.stats-section {
  margin-bottom: 30px;
}

.stat-card {
  background: white;
  border-radius: 16px;
  padding: 24px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  display: flex;
  align-items: center;
  gap: 16px;
  transition: all 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.12);
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.stat-icon.online {
  background: linear-gradient(135deg, #67C23A, #85ce61);
}

.stat-icon.offline {
  background: linear-gradient(135deg, #F56C6C, #f78989);
}

.stat-icon.active {
  background: linear-gradient(135deg, #409EFF, #66b1ff);
}

.stat-icon.total {
  background: linear-gradient(135deg, #909399, #a6a9ad);
}

.stat-number {
  font-size: 24px;
  font-weight: 700;
  color: #303133;
  line-height: 1;
}

.stat-label {
  font-size: 14px;
  color: #909399;
  margin-top: 4px;
}

/* 设备区域 */
.devices-section {
  background: white;
  border-radius: 16px;
  padding: 24px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}

/* 网格视图 */
.devices-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 20px;
  margin-bottom: 30px;
}

.device-card {
  background: #f8f9fa;
  border: 2px solid transparent;
  border-radius: 16px;
  padding: 20px;
  cursor: pointer;
  transition: all 0.3s ease;
  position: relative;
}

.device-card:hover {
  background: white;
  border-color: #409EFF;
  box-shadow: 0 8px 25px rgba(64, 158, 255, 0.15);
  transform: translateY(-2px);
}

.device-status {
  position: absolute;
  top: 16px;
  right: 16px;
  width: 10px;
  height: 10px;
  border-radius: 50%;
  border: 2px solid white;
}

.device-status.online {
  background: #67C23A;
  box-shadow: 0 0 10px rgba(103, 194, 58, 0.5);
}

.device-status.offline {
  background: #F56C6C;
}

.device-status.idle {
  background: #909399;
}

.device-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 16px;
}

.device-icon {
  color: #606266;
}

.device-actions {
  opacity: 0;
  transition: opacity 0.3s ease;
}

.device-card:hover .device-actions {
  opacity: 1;
}

.device-info {
  margin-bottom: 16px;
}

.device-name {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 4px;
  line-height: 1.4;
}

.device-id {
  font-size: 12px;
  color: #909399;
  margin-bottom: 8px;
}

.device-state {
  font-size: 14px;
  font-weight: 500;
  padding: 4px 8px;
  border-radius: 6px;
  display: inline-block;
}

.device-state.state-on {
  color: #67C23A;
  background: rgba(103, 194, 58, 0.1);
}

.device-state.state-off {
  color: #909399;
  background: rgba(144, 147, 153, 0.1);
}

.device-state.state-offline {
  color: #F56C6C;
  background: rgba(245, 108, 108, 0.1);
}

.device-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 12px;
  color: #909399;
}

.device-domain {
  background: #f0f0f0;
  padding: 2px 8px;
  border-radius: 4px;
}

/* 列表视图 */
.devices-table {
  margin-bottom: 30px;
}

.table-status {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  margin: 0 auto;
}

.table-device-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.table-device-icon {
  color: #606266;
}

.table-device-name {
  font-weight: 500;
  color: #303133;
}

.table-device-id {
  font-size: 12px;
  color: #909399;
}

/* 分页 */
.pagination-section {
  display: flex;
  justify-content: center;
  padding-top: 20px;
  border-top: 1px solid #f0f0f0;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .header-content {
    flex-direction: column;
    gap: 20px;
    text-align: center;
  }
  
  .devices-grid {
    grid-template-columns: 1fr;
  }
  
  .filter-section .el-row {
    flex-direction: column;
    gap: 16px;
  }
  
  .filter-section .el-col {
    width: 100% !important;
  }
}

@media (max-width: 480px) {
  .page-header {
    padding: 20px;
  }
  
  .header-info h1 {
    font-size: 24px;
  }
  
  .devices-section {
    padding: 16px;
  }
  
  .device-card {
    padding: 16px;
  }
}
</style>