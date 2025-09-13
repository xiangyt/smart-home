<template>
  <div class="automation-container">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <div class="header-info">
          <h1>自动化规则</h1>
          <p>设置智能自动化规则，让您的家居更加智能</p>
        </div>
        <div class="header-actions">
          <el-button type="primary" @click="createAutomation">
            <el-icon><Plus /></el-icon>
            创建规则
          </el-button>
        </div>
      </div>
    </div>

    <!-- 规则统计 -->
    <div class="stats-section">
      <el-row :gutter="20">
        <el-col :span="6">
          <div class="stat-card">
            <div class="stat-icon active">
              <el-icon size="24"><CircleCheck /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-number">{{ activeRules }}</div>
              <div class="stat-label">活跃规则</div>
            </div>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="stat-card">
            <div class="stat-icon total">
              <el-icon size="24"><Setting /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-number">{{ totalRules }}</div>
              <div class="stat-label">总规则数</div>
            </div>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="stat-card">
            <div class="stat-icon executed">
              <el-icon size="24"><Clock /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-number">{{ todayExecutions }}</div>
              <div class="stat-label">今日执行</div>
            </div>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="stat-card">
            <div class="stat-icon success">
              <el-icon size="24"><SuccessFilled /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-number">{{ successRate }}%</div>
              <div class="stat-label">成功率</div>
            </div>
          </div>
        </el-col>
      </el-row>
    </div>

    <!-- 快速规则模板 -->
    <div class="templates-section">
      <h3>快速创建</h3>
      <el-row :gutter="20">
        <el-col :xs="24" :sm="12" :md="8" v-for="template in ruleTemplates" :key="template.id">
          <div class="template-card" @click="useTemplate(template)">
            <div class="template-icon" :style="{ background: template.color }">
              <el-icon size="24" color="white">
                <component :is="template.icon" />
              </el-icon>
            </div>
            <div class="template-content">
              <div class="template-name">{{ template.name }}</div>
              <div class="template-description">{{ template.description }}</div>
            </div>
            <div class="template-arrow">
              <el-icon><ArrowRight /></el-icon>
            </div>
          </div>
        </el-col>
      </el-row>
    </div>

    <!-- 自动化规则列表 -->
    <div class="rules-section">
      <div class="section-header">
        <h3>自动化规则</h3>
        <div class="header-controls">
          <el-select v-model="filterStatus" placeholder="筛选状态" size="small" style="width: 120px">
            <el-option label="全部" value="" />
            <el-option label="启用" value="enabled" />
            <el-option label="禁用" value="disabled" />
          </el-select>
          <el-select v-model="filterType" placeholder="规则类型" size="small" style="width: 120px">
            <el-option label="全部类型" value="" />
            <el-option label="时间触发" value="time" />
            <el-option label="设备触发" value="device" />
            <el-option label="传感器触发" value="sensor" />
          </el-select>
        </div>
      </div>

      <div class="rules-grid">
        <div 
          v-for="rule in filteredRules" 
          :key="rule.id"
          class="rule-card"
        >
          <div class="rule-header">
            <div class="rule-status">
              <el-switch 
                v-model="rule.enabled" 
                @change="toggleRule(rule)"
                :loading="rule.loading"
              />
            </div>
            <div class="rule-actions">
              <el-dropdown @command="handleRuleAction">
                <el-button size="small" circle>
                  <el-icon><MoreFilled /></el-icon>
                </el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item :command="{action: 'edit', rule}">编辑规则</el-dropdown-item>
                    <el-dropdown-item :command="{action: 'duplicate', rule}">复制规则</el-dropdown-item>
                    <el-dropdown-item :command="{action: 'test', rule}">测试执行</el-dropdown-item>
                    <el-dropdown-item :command="{action: 'logs', rule}">查看日志</el-dropdown-item>
                    <el-dropdown-item :command="{action: 'delete', rule}" divided>删除规则</el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>
          </div>

          <div class="rule-content">
            <div class="rule-name">{{ rule.name }}</div>
            <div class="rule-description">{{ rule.description }}</div>
            
            <div class="rule-flow">
              <div class="flow-item trigger">
                <div class="flow-icon">
                  <el-icon size="16"><Lightning /></el-icon>
                </div>
                <div class="flow-text">
                  <div class="flow-label">触发条件</div>
                  <div class="flow-value">{{ rule.trigger }}</div>
                </div>
              </div>
              
              <div class="flow-arrow">
                <el-icon><ArrowRight /></el-icon>
              </div>
              
              <div class="flow-item condition" v-if="rule.condition">
                <div class="flow-icon">
                  <el-icon size="16"><Filter /></el-icon>
                </div>
                <div class="flow-text">
                  <div class="flow-label">条件</div>
                  <div class="flow-value">{{ rule.condition }}</div>
                </div>
              </div>
              
              <div class="flow-arrow" v-if="rule.condition">
                <el-icon><ArrowRight /></el-icon>
              </div>
              
              <div class="flow-item action">
                <div class="flow-icon">
                  <el-icon size="16"><Setting /></el-icon>
                </div>
                <div class="flow-text">
                  <div class="flow-label">执行动作</div>
                  <div class="flow-value">{{ rule.action }}</div>
                </div>
              </div>
            </div>
          </div>

          <div class="rule-footer">
            <div class="rule-stats">
              <div class="stat-item">
                <el-icon size="14"><Clock /></el-icon>
                <span>最后执行: {{ formatLastExecution(rule.lastExecution) }}</span>
              </div>
              <div class="stat-item">
                <el-icon size="14"><DataAnalysis /></el-icon>
                <span>执行次数: {{ rule.executionCount }}</span>
              </div>
            </div>
            <div class="rule-type">
              <el-tag size="small" :type="getRuleTypeColor(rule.type)">
                {{ getRuleTypeName(rule.type) }}
              </el-tag>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { 
  Plus, CircleCheck, Setting, Clock, SuccessFilled, ArrowRight, 
  MoreFilled, Lightning, Filter, DataAnalysis, Sunny, Moon, 
  Timer, Thermometer, House, Switch
} from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

// 响应式数据
const filterStatus = ref('')
const filterType = ref('')

// 规则模板
const ruleTemplates = ref([
  {
    id: 1,
    name: '定时开关灯',
    description: '根据时间自动开关灯光',
    icon: Sunny,
    color: 'linear-gradient(135deg, #FFA726, #FF7043)'
  },
  {
    id: 2,
    name: '夜间模式',
    description: '晚上自动调暗所有灯光',
    icon: Moon,
    color: 'linear-gradient(135deg, #42A5F5, #1E88E5)'
  },
  {
    id: 3,
    name: '温度控制',
    description: '根据温度自动调节空调',
    icon: Thermometer,
    color: 'linear-gradient(135deg, #66BB6A, #43A047)'
  },
  {
    id: 4,
    name: '离家模式',
    description: '检测到离家时关闭所有设备',
    icon: House,
    color: 'linear-gradient(135deg, #AB47BC, #8E24AA)'
  },
  {
    id: 5,
    name: '安全监控',
    description: '检测到异常时发送通知',
    icon: Switch,
    color: 'linear-gradient(135deg, #EF5350, #E53935)'
  },
  {
    id: 6,
    name: '定时任务',
    description: '设置定时执行的任务',
    icon: Timer,
    color: 'linear-gradient(135deg, #26A69A, #00897B)'
  }
])

// 自动化规则数据
const automationRules = ref([
  {
    id: 1,
    name: '早晨唤醒',
    description: '工作日早上7点自动开启卧室灯光',
    enabled: true,
    loading: false,
    trigger: '时间 07:00 (工作日)',
    condition: '当前时间为工作日',
    action: '开启卧室灯光 (50%亮度)',
    type: 'time',
    lastExecution: '2024-01-15 07:00:00',
    executionCount: 45
  },
  {
    id: 2,
    name: '离家安全',
    description: '检测到所有人离家时自动关闭电器',
    enabled: true,
    loading: false,
    trigger: '所有人离家',
    condition: '离家超过5分钟',
    action: '关闭所有灯光和电器',
    type: 'device',
    lastExecution: '2024-01-14 18:30:00',
    executionCount: 12
  },
  {
    id: 3,
    name: '温度调节',
    description: '室内温度过高时自动开启空调',
    enabled: false,
    loading: false,
    trigger: '温度传感器 > 26°C',
    condition: '当前时间 08:00-22:00',
    action: '开启空调制冷模式',
    type: 'sensor',
    lastExecution: '2024-01-13 15:20:00',
    executionCount: 8
  },
  {
    id: 4,
    name: '夜间模式',
    description: '晚上10点后自动调暗所有灯光',
    enabled: true,
    loading: false,
    trigger: '时间 22:00',
    condition: null,
    action: '调暗所有灯光至20%',
    type: 'time',
    lastExecution: '2024-01-14 22:00:00',
    executionCount: 30
  }
])

// 计算属性
const totalRules = computed(() => automationRules.value.length)
const activeRules = computed(() => automationRules.value.filter(rule => rule.enabled).length)
const todayExecutions = computed(() => 23) // 模拟数据
const successRate = computed(() => 98) // 模拟数据

const filteredRules = computed(() => {
  let rules = automationRules.value

  if (filterStatus.value) {
    if (filterStatus.value === 'enabled') {
      rules = rules.filter(rule => rule.enabled)
    } else if (filterStatus.value === 'disabled') {
      rules = rules.filter(rule => !rule.enabled)
    }
  }

  if (filterType.value) {
    rules = rules.filter(rule => rule.type === filterType.value)
  }

  return rules
})

// 方法
const createAutomation = () => {
  ElMessage.info('创建自动化规则功能开发中...')
}

const useTemplate = (template: any) => {
  ElMessage.success(`使用模板: ${template.name}`)
}

const toggleRule = (rule: any) => {
  rule.loading = true
  setTimeout(() => {
    rule.loading = false
    ElMessage.success(`规则 "${rule.name}" 已${rule.enabled ? '启用' : '禁用'}`)
  }, 1000)
}

const handleRuleAction = ({ action, rule }: { action: string, rule: any }) => {
  switch (action) {
    case 'edit':
      ElMessage.info('编辑规则功能开发中...')
      break
    case 'duplicate':
      ElMessage.info('复制规则功能开发中...')
      break
    case 'test':
      ElMessage.success(`规则 "${rule.name}" 测试执行成功`)
      break
    case 'logs':
      ElMessage.info('查看日志功能开发中...')
      break
    case 'delete':
      ElMessage.info('删除规则功能开发中...')
      break
  }
}

const formatLastExecution = (timestamp: string) => {
  if (!timestamp) return '从未执行'
  const date = new Date(timestamp)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const hours = Math.floor(diff / (1000 * 60 * 60))
  
  if (hours < 1) return '刚刚'
  if (hours < 24) return `${hours}小时前`
  const days = Math.floor(hours / 24)
  return `${days}天前`
}

const getRuleTypeName = (type: string) => {
  const typeNames: Record<string, string> = {
    time: '时间触发',
    device: '设备触发',
    sensor: '传感器触发'
  }
  return typeNames[type] || type
}

const getRuleTypeColor = (type: string) => {
  const typeColors: Record<string, string> = {
    time: 'primary',
    device: 'success',
    sensor: 'warning'
  }
  return typeColors[type] || 'info'
}
</script>

<style scoped>
.automation-container {
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

/* 统计区域 */
.stats-section {
  margin-bottom: 40px;
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

.stat-icon.active {
  background: linear-gradient(135deg, #67C23A, #85ce61);
}

.stat-icon.total {
  background: linear-gradient(135deg, #409EFF, #66b1ff);
}

.stat-icon.executed {
  background: linear-gradient(135deg, #E6A23C, #eebe77);
}

.stat-icon.success {
  background: linear-gradient(135deg, #67C23A, #85ce61);
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

/* 模板区域 */
.templates-section {
  margin-bottom: 40px;
}

.templates-section h3 {
  margin: 0 0 20px 0;
  font-size: 20px;
  font-weight: 600;
  color: #303133;
}

.template-card {
  background: white;
  border-radius: 16px;
  padding: 20px;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 16px;
}

.template-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.12);
}

.template-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.template-content {
  flex: 1;
}

.template-name {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 4px;
}

.template-description {
  font-size: 14px;
  color: #909399;
}

.template-arrow {
  color: #909399;
  transition: all 0.3s ease;
}

.template-card:hover .template-arrow {
  color: #409EFF;
  transform: translateX(4px);
}

/* 规则区域 */
.rules-section {
  background: white;
  border-radius: 16px;
  padding: 24px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  padding-bottom: 16px;
  border-bottom: 1px solid #f0f0f0;
}

.section-header h3 {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
  color: #303133;
}

.header-controls {
  display: flex;
  gap: 12px;
}

.rules-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(400px, 1fr));
  gap: 20px;
}

.rule-card {
  background: #f8f9fa;
  border: 2px solid transparent;
  border-radius: 16px;
  padding: 20px;
  transition: all 0.3s ease;
}

.rule-card:hover {
  background: white;
  border-color: #409EFF;
  box-shadow: 0 4px 15px rgba(64, 158, 255, 0.1);
}

.rule-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.rule-actions {
  opacity: 0;
  transition: opacity 0.3s ease;
}

.rule-card:hover .rule-actions {
  opacity: 1;
}

.rule-name {
  font-size: 18px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 8px;
}

.rule-description {
  font-size: 14px;
  color: #909399;
  margin-bottom: 20px;
  line-height: 1.4;
}

.rule-flow {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 20px;
  flex-wrap: wrap;
}

.flow-item {
  display: flex;
  align-items: center;
  gap: 8px;
  background: white;
  padding: 12px;
  border-radius: 8px;
  border: 1px solid #e4e7ed;
  min-width: 0;
  flex: 1;
}

.flow-icon {
  color: #606266;
  flex-shrink: 0;
}

.flow-text {
  min-width: 0;
}

.flow-label {
  font-size: 12px;
  color: #909399;
  margin-bottom: 2px;
}

.flow-value {
  font-size: 13px;
  color: #303133;
  font-weight: 500;
  word-break: break-all;
}

.flow-arrow {
  color: #909399;
  flex-shrink: 0;
}

.rule-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 16px;
  border-top: 1px solid #f0f0f0;
}

.rule-stats {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: #909399;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .header-content {
    flex-direction: column;
    gap: 20px;
    text-align: center;
  }
  
  .rules-grid {
    grid-template-columns: 1fr;
  }
  
  .section-header {
    flex-direction: column;
    gap: 16px;
    align-items: flex-start;
  }
  
  .rule-flow {
    flex-direction: column;
    align-items: stretch;
  }
  
  .flow-arrow {
    transform: rotate(90deg);
    align-self: center;
  }
}

@media (max-width: 480px) {
  .page-header {
    padding: 20px;
  }
  
  .header-info h1 {
    font-size: 24px;
  }
  
  .rules-section {
    padding: 16px;
  }
  
  .rule-card {
    padding: 16px;
  }
  
  .template-card {
    padding: 16px;
  }
}
</style>