<template>
  <div class="scenes-container">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <div class="header-info">
          <h1>场景控制</h1>
          <p>创建和管理智能场景，一键控制多个设备</p>
        </div>
        <div class="header-actions">
          <el-button type="primary" @click="createScene">
            <el-icon><Plus /></el-icon>
            创建场景
          </el-button>
        </div>
      </div>
    </div>

    <!-- 快速场景 -->
    <div class="quick-scenes">
      <h3>快速场景</h3>
      <el-row :gutter="20">
        <el-col :xs="24" :sm="12" :md="6" v-for="scene in quickScenes" :key="scene.id">
          <div class="scene-card quick-scene" @click="activateScene(scene)">
            <div class="scene-background" :style="{ background: scene.gradient }">
              <div class="scene-icon">
                <el-icon size="32" :color="scene.iconColor">
                  <component :is="scene.icon" />
                </el-icon>
              </div>
            </div>
            <div class="scene-content">
              <div class="scene-name">{{ scene.name }}</div>
              <div class="scene-description">{{ scene.description }}</div>
              <div class="scene-devices">{{ scene.deviceCount }} 个设备</div>
            </div>
            <div class="scene-status" :class="{ active: scene.active }">
              {{ scene.active ? '已激活' : '未激活' }}
            </div>
          </div>
        </el-col>
      </el-row>
    </div>

    <!-- 自定义场景 -->
    <div class="custom-scenes">
      <div class="section-header">
        <h3>自定义场景</h3>
        <el-button type="text" @click="showAllScenes">
          查看全部
          <el-icon><ArrowRight /></el-icon>
        </el-button>
      </div>
      
      <div class="scenes-grid">
        <div 
          v-for="scene in customScenes" 
          :key="scene.id"
          class="scene-card custom-scene"
          @click="activateScene(scene)"
        >
          <div class="scene-header">
            <div class="scene-icon">
              <el-icon size="24">
                <component :is="scene.icon" />
              </el-icon>
            </div>
            <div class="scene-actions">
              <el-dropdown @command="handleSceneAction">
                <el-button size="small" circle>
                  <el-icon><MoreFilled /></el-icon>
                </el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item :command="{action: 'activate', scene}">激活场景</el-dropdown-item>
                    <el-dropdown-item :command="{action: 'edit', scene}">编辑场景</el-dropdown-item>
                    <el-dropdown-item :command="{action: 'duplicate', scene}">复制场景</el-dropdown-item>
                    <el-dropdown-item :command="{action: 'delete', scene}" divided>删除场景</el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>
          </div>
          
          <div class="scene-info">
            <div class="scene-name">{{ scene.name }}</div>
            <div class="scene-description">{{ scene.description }}</div>
          </div>
          
          <div class="scene-devices-list">
            <div class="device-item" v-for="device in scene.devices.slice(0, 3)" :key="device">
              <el-icon size="14"><Monitor /></el-icon>
              <span>{{ device }}</span>
            </div>
            <div v-if="scene.devices.length > 3" class="more-devices">
              +{{ scene.devices.length - 3 }} 更多
            </div>
          </div>
          
          <div class="scene-footer">
            <div class="scene-status" :class="{ active: scene.active }">
              <div class="status-dot"></div>
              {{ scene.active ? '已激活' : '未激活' }}
            </div>
            <div class="scene-last-used">
              {{ formatLastUsed(scene.lastUsed) }}
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 场景统计 -->
    <div class="scene-stats">
      <el-row :gutter="20">
        <el-col :span="8">
          <div class="stat-card">
            <div class="stat-icon">
              <el-icon size="24" color="#409EFF"><Setting /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-number">{{ totalScenes }}</div>
              <div class="stat-label">总场景数</div>
            </div>
          </div>
        </el-col>
        <el-col :span="8">
          <div class="stat-card">
            <div class="stat-icon">
              <el-icon size="24" color="#67C23A"><CircleCheck /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-number">{{ activeScenes }}</div>
              <div class="stat-label">活跃场景</div>
            </div>
          </div>
        </el-col>
        <el-col :span="8">
          <div class="stat-card">
            <div class="stat-icon">
              <el-icon size="24" color="#E6A23C"><Clock /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-number">{{ todayUsage }}</div>
              <div class="stat-label">今日使用</div>
            </div>
          </div>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { 
  Plus, ArrowRight, MoreFilled, Monitor, Setting, CircleCheck, Clock,
  Sunny, Moon, House, VideoPlay, Switch
} from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

// 快速场景数据
const quickScenes = ref([
  {
    id: 1,
    name: '回家模式',
    description: '开启客厅灯光，调节温度',
    deviceCount: 8,
    active: false,
    icon: House,
    iconColor: '#ffffff',
    gradient: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)'
  },
  {
    id: 2,
    name: '离家模式',
    description: '关闭所有灯光和电器',
    deviceCount: 12,
    active: false,
    icon: Switch,
    iconColor: '#ffffff',
    gradient: 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)'
  },
  {
    id: 3,
    name: '睡眠模式',
    description: '调暗灯光，关闭娱乐设备',
    deviceCount: 6,
    active: true,
    icon: Moon,
    iconColor: '#ffffff',
    gradient: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)'
  },
  {
    id: 4,
    name: '观影模式',
    description: '调暗灯光，开启投影仪',
    deviceCount: 5,
    active: false,
    icon: VideoPlay,
    iconColor: '#ffffff',
    gradient: 'linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)'
  }
])

// 自定义场景数据
const customScenes = ref([
  {
    id: 5,
    name: '早晨唤醒',
    description: '逐渐调亮灯光，播放轻音乐',
    devices: ['客厅灯', '卧室灯', '音响', '窗帘'],
    active: false,
    icon: Sunny,
    lastUsed: '2024-01-15 07:30:00'
  },
  {
    id: 6,
    name: '工作专注',
    description: '调节灯光亮度，关闭娱乐设备',
    devices: ['书房灯', '台灯', '空调'],
    active: false,
    icon: Monitor,
    lastUsed: '2024-01-14 14:20:00'
  },
  {
    id: 7,
    name: '聚会模式',
    description: '开启彩色灯光，调节音响',
    devices: ['客厅灯', '餐厅灯', '音响', '彩灯', '空调', '投影仪'],
    active: false,
    icon: Setting,
    lastUsed: '2024-01-13 19:45:00'
  }
])

// 计算属性
const totalScenes = computed(() => quickScenes.value.length + customScenes.value.length)
const activeScenes = computed(() => 
  [...quickScenes.value, ...customScenes.value].filter(scene => scene.active).length
)
const todayUsage = computed(() => 15) // 模拟数据

// 方法
const createScene = () => {
  ElMessage.info('创建场景功能开发中...')
}

const activateScene = (scene: any) => {
  // 切换场景状态
  scene.active = !scene.active
  ElMessage.success(`${scene.name} ${scene.active ? '已激活' : '已停用'}`)
}

const showAllScenes = () => {
  ElMessage.info('查看全部场景功能开发中...')
}

const handleSceneAction = ({ action, scene }: { action: string, scene: any }) => {
  switch (action) {
    case 'activate':
      activateScene(scene)
      break
    case 'edit':
      ElMessage.info('编辑场景功能开发中...')
      break
    case 'duplicate':
      ElMessage.info('复制场景功能开发中...')
      break
    case 'delete':
      ElMessage.info('删除场景功能开发中...')
      break
  }
}

const formatLastUsed = (timestamp: string) => {
  const date = new Date(timestamp)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))
  
  if (days === 0) return '今天'
  if (days === 1) return '昨天'
  if (days < 7) return `${days}天前`
  return date.toLocaleDateString('zh-CN')
}
</script>

<style scoped>
.scenes-container {
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

/* 快速场景 */
.quick-scenes {
  margin-bottom: 40px;
}

.quick-scenes h3 {
  margin: 0 0 20px 0;
  font-size: 20px;
  font-weight: 600;
  color: #303133;
}

.scene-card {
  background: white;
  border-radius: 16px;
  padding: 20px;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  margin-bottom: 20px;
  position: relative;
  overflow: hidden;
}

.scene-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.15);
}

.quick-scene {
  text-align: center;
  min-height: 180px;
}

.scene-background {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 80px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.scene-content {
  position: relative;
  z-index: 1;
  padding-top: 60px;
}

.scene-name {
  font-size: 18px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 8px;
}

.scene-description {
  font-size: 14px;
  color: #909399;
  margin-bottom: 12px;
  line-height: 1.4;
}

.scene-devices {
  font-size: 12px;
  color: #606266;
}

.scene-status {
  position: absolute;
  top: 16px;
  right: 16px;
  font-size: 12px;
  padding: 4px 8px;
  border-radius: 12px;
  background: rgba(0, 0, 0, 0.1);
  color: #909399;
  z-index: 2;
}

.scene-status.active {
  background: rgba(103, 194, 58, 0.2);
  color: #67C23A;
}

/* 自定义场景 */
.custom-scenes {
  margin-bottom: 40px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.section-header h3 {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
  color: #303133;
}

.scenes-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
}

.custom-scene {
  min-height: 200px;
}

.scene-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.scene-icon {
  color: #606266;
}

.scene-actions {
  opacity: 0;
  transition: opacity 0.3s ease;
}

.scene-card:hover .scene-actions {
  opacity: 1;
}

.scene-info {
  margin-bottom: 16px;
}

.scene-devices-list {
  margin-bottom: 16px;
}

.device-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: #909399;
  margin-bottom: 4px;
}

.more-devices {
  font-size: 12px;
  color: #409EFF;
  margin-top: 4px;
}

.scene-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 12px;
  border-top: 1px solid #f0f0f0;
}

.scene-footer .scene-status {
  position: static;
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  padding: 0;
  background: none;
}

.status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #ddd;
}

.scene-status.active .status-dot {
  background: #67C23A;
}

.scene-last-used {
  font-size: 12px;
  color: #909399;
}

/* 场景统计 */
.scene-stats {
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
  background: #f8f9fa;
  display: flex;
  align-items: center;
  justify-content: center;
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

/* 响应式设计 */
@media (max-width: 768px) {
  .header-content {
    flex-direction: column;
    gap: 20px;
    text-align: center;
  }
  
  .scenes-grid {
    grid-template-columns: 1fr;
  }
  
  .section-header {
    flex-direction: column;
    gap: 16px;
    align-items: flex-start;
  }
}

@media (max-width: 480px) {
  .page-header {
    padding: 20px;
  }
  
  .header-info h1 {
    font-size: 24px;
  }
  
  .scene-card {
    padding: 16px;
  }
  
  .quick-scene {
    min-height: 160px;
  }
}
</style>