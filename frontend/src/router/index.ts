import { createRouter, createWebHistory } from 'vue-router'
import Overview from '@/views/Overview.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      redirect: '/overview'
    },
    {
      path: '/overview',
      name: 'Overview',
      component: Overview
    },
    {
      path: '/devices',
      name: 'Devices',
      component: () => import('@/views/Devices.vue')
    },
    {
      path: '/scenes',
      name: 'Scenes',
      component: () => import('@/views/Scenes.vue')
    },
    {
      path: '/automation',
      name: 'Automation',
      component: () => import('@/views/Automation.vue')
    }
  ]
})

export default router