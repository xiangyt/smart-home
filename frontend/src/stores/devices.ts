import { defineStore } from 'pinia'
import { ref } from 'vue'
import { deviceApi } from '@/api/devices'

export interface Device {
  entity_id: string
  friendly_name: string
  state: string
  attributes: Record<string, any>
  domain: string
  device_class?: string
  unit_of_measurement?: string
  last_changed: string
  last_updated: string
}

export const useDevicesStore = defineStore('devices', () => {
  const devices = ref<Device[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  // 获取所有设备
  const fetchDevices = async () => {
    loading.value = true
    error.value = null
    try {
      const response = await deviceApi.getDevices()
      // 处理后端返回的数据格式 {data: [...], count: number, cached: boolean}
      if (response.data && Array.isArray(response.data.data)) {
        devices.value = response.data.data
      } else if (Array.isArray(response.data)) {
        devices.value = response.data
      } else {
        devices.value = []
        console.warn('Unexpected data format:', response.data)
      }
    } catch (err) {
      error.value = '获取设备列表失败'
      console.error('Failed to fetch devices:', err)
      devices.value = []
    } finally {
      loading.value = false
    }
  }

  // 控制设备
  const controlDevice = async (entityId: string, service: string, data?: any) => {
    try {
      await deviceApi.controlDevice(entityId, service, data)
      // 重新获取设备状态
      await fetchDevices()
    } catch (err) {
      error.value = '设备控制失败'
      console.error('Failed to control device:', err)
    }
  }

  // 按域名分组设备
  const getDevicesByDomain = (domain: string) => {
    if (!Array.isArray(devices.value)) {
      return []
    }
    return devices.value.filter(device => device.domain === domain)
  }

  // 获取在线设备数量
  const getOnlineDevicesCount = () => {
    if (!Array.isArray(devices.value)) {
      return 0
    }
    return devices.value.filter(device => 
      device.state !== 'unavailable' && device.state !== 'unknown'
    ).length
  }

  return {
    devices,
    loading,
    error,
    fetchDevices,
    controlDevice,
    getDevicesByDomain,
    getOnlineDevicesCount
  }
})