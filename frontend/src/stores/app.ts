// 全局应用状态管理（Pinia Store）
// 管理主题、配置、全局通知等应用级状态

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { GetConfig, SaveConfig, GetTheme, SetTheme } from '../../wailsjs/go/config/Config'

// 应用配置接口定义
interface AppConfig {
  theme: string
  language: string
  fontSize: number
  density: string
  transparency: number
  alwaysOnTop: boolean
  minToTray: boolean
  autoBackup: boolean
}

// Toast 通知接口
export interface Toast {
  id: string
  type: 'success' | 'error' | 'warning' | 'info'
  message: string
  duration?: number
}

export const useAppStore = defineStore('app', () => {
  // ============ 状态定义 ============

  // 应用配置
  const config = ref<AppConfig>({
    theme: 'dark',
    language: 'zh-CN',
    fontSize: 14,
    density: 'normal',
    transparency: 0.95,
    alwaysOnTop: false,
    minToTray: true,
    autoBackup: true,
  })

  // 当前主题
  const theme = computed(() => config.value.theme)

  // 全局搜索显示状态
  const searchVisible = ref(false)

  // Toast 通知列表
  const toasts = ref<Toast[]>([])

  // 侧边栏折叠状态
  const sidebarCollapsed = ref(false)

  // ============ Actions ============

  // 从后端加载配置
  async function loadConfig() {
    try {
      const cfg = await GetConfig()
      if (cfg) {
        config.value = cfg as AppConfig
        applyTheme(config.value.theme)
      }
    } catch (err) {
      console.error('加载配置失败:', err)
    }
  }

  // 保存配置到后端
  async function saveConfig(newConfig: Partial<AppConfig>) {
    config.value = { ...config.value, ...newConfig }
    try {
      await SaveConfig(JSON.stringify(config.value))
      applyTheme(config.value.theme)
    } catch (err) {
      showToast('error', '保存配置失败: ' + String(err))
    }
  }

  // 切换主题
  async function switchTheme(themeId: string) {
    config.value.theme = themeId
    applyTheme(themeId)
    try {
      await SetTheme(themeId)
    } catch (err) {
      console.error('切换主题失败:', err)
    }
  }

  // 应用主题到 DOM
  function applyTheme(themeId: string) {
    const root = document.documentElement
    // 移除所有主题类
    root.classList.remove('dark', 'light', 'theme-blue', 'theme-green', 'theme-purple', 'theme-orange', 'theme-mica')

    if (themeId === 'auto') {
      const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
      root.classList.add(prefersDark ? 'dark' : 'light')
    } else if (themeId === 'light') {
      root.classList.add('light')
    } else {
      root.classList.add('dark')
      if (themeId !== 'dark') {
        root.classList.add(`theme-${themeId}`)
      }
    }
  }

  // 显示 Toast 通知
  function showToast(type: Toast['type'], message: string, duration = 3000) {
    const id = Date.now().toString() + Math.random().toString(36).slice(2)
    const toast: Toast = { id, type, message, duration }
    toasts.value.push(toast)

    // 自动移除
    if (duration > 0) {
      setTimeout(() => removeToast(id), duration)
    }
  }

  // 移除 Toast
  function removeToast(id: string) {
    const index = toasts.value.findIndex(t => t.id === id)
    if (index !== -1) toasts.value.splice(index, 1)
  }

  // 切换侧边栏折叠
  function toggleSidebar() {
    sidebarCollapsed.value = !sidebarCollapsed.value
  }

  return {
    config,
    theme,
    searchVisible,
    toasts,
    sidebarCollapsed,
    loadConfig,
    saveConfig,
    switchTheme,
    showToast,
    removeToast,
    toggleSidebar,
  }
})
