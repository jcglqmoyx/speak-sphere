import { useDark, useToggle } from '@vueuse/core'

// 主题状态管理
export const useThemeStore = () => {
  // 使用 localStorage 持久化存储主题偏好
  const isDark = useDark({
    storageKey: 'speak-sphere-theme',
    valueDark: 'dark',
    valueLight: 'light'
  })
  
  const toggleDark = useToggle(isDark)
  
  // 初始化主题（避免页面加载时的闪烁）
  const initTheme = () => {
    const savedTheme = localStorage.getItem('speak-sphere-theme')
    if (!savedTheme) {
      // 如果没有保存的主题，使用系统偏好但不要立即应用
      const systemIsDark = window.matchMedia('(prefers-color-scheme: dark)').matches
      isDark.value = systemIsDark
    }
  }
  
  return {
    isDark,
    toggleDark,
    initTheme
  }
}

// 创建全局主题实例
let themeInstance = null

export const getThemeInstance = () => {
  if (!themeInstance) {
    themeInstance = useThemeStore()
    themeInstance.initTheme()
  }
  return themeInstance
}
