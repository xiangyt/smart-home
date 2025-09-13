import axios from 'axios'

const api = axios.create({
  baseURL: 'http://localhost:8080/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  }
})

// 请求拦截器
api.interceptors.request.use(
  (config) => {
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
api.interceptors.response.use(
  (response) => {
    return response
  },
  (error) => {
    console.error('API Error:', error)
    return Promise.reject(error)
  }
)

export const deviceApi = {
  // 获取所有设备
  getDevices: () => api.get('/devices'),
  
  // 获取设备状态
  getDeviceState: (entityId: string) => api.get(`/devices/${entityId}/state`),
  
  // 控制设备
  controlDevice: (entityId: string, service: string, data?: any) => 
    api.post(`/devices/${entityId}/control`, { service, data }),
  
  // 获取设备历史数据
  getDeviceHistory: (entityId: string, startTime?: string, endTime?: string) => 
    api.get(`/devices/${entityId}/history`, { params: { start_time: startTime, end_time: endTime } })
}

export default api