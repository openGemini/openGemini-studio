import { ref } from "vue"
import type { AppSettings } from "../types"
import { GetSetting, UpdateSetting } from "../../wailsjs/go/main/App"
import { main } from "../../wailsjs/go/models"

const DEFAULT_SETTINGS: AppSettings = {
  language: "en",
  themeMode: "light",
  customFont: "",
  maxHistoryCount: 50,
  dataDirectory: "./data",
  debug: false,
}

// Data transformation utilities
const toBackendSettings = (settings: AppSettings): main.AppSetting => {
  return {
    language: settings.language,
    theme_mode: settings.themeMode,
    custom_font: settings.customFont,
    max_history_count: settings.maxHistoryCount,
    data_dir: settings.dataDirectory,
    debug: settings.debug,
  }
}

const fromBackendSettings = (backendSettings: main.AppSetting): AppSettings => {
  return {
    language: backendSettings.language,
    themeMode: backendSettings.theme_mode as "light" | "dark" | "system",
    customFont: backendSettings.custom_font || "",
    maxHistoryCount: backendSettings.max_history_count,
    dataDirectory: backendSettings.data_dir,
    debug: backendSettings.debug || false,
  }
}

let settingsInstance: ReturnType<typeof createSettings> | null = null

function createSettings() {
  const settings = ref<AppSettings>({ ...DEFAULT_SETTINGS })
  const isInitialized = ref(false)

  // 初始化设置
  const initSettings = async () => {
    if (!isInitialized.value) {
      try {
        // Load settings from backend
        const backendSettings = await GetSetting()
        if (backendSettings) {
          settings.value = fromBackendSettings(backendSettings)
        }
      } catch (error) {
        console.error("Failed to load settings from backend:", error)
        // Fall back to default settings
        settings.value = { ...DEFAULT_SETTINGS }
      }
      isInitialized.value = true
    }
  }

  // 立即初始化
  initSettings()

  // 更新设置
  const updateSettings = async (newSettings: Partial<AppSettings>) => {
    settings.value = { ...settings.value, ...newSettings }
    try {
      // Save to backend
      await UpdateSetting(toBackendSettings(settings.value))
    } catch (error) {
      console.error("Failed to save settings to backend:", error)
    }
  }

  // 重置设置
  const resetSettings = async () => {
    settings.value = { ...DEFAULT_SETTINGS }
    try {
      // Save default settings to backend
      await UpdateSetting(toBackendSettings(settings.value))
    } catch (error) {
      console.error("Failed to save default settings to backend:", error)
    }
  }

  return {
    settings,
    updateSettings,
    resetSettings,
    initSettings,
  }
}

export function useSettings() {
  if (!settingsInstance) {
    settingsInstance = createSettings()
  }
  return settingsInstance
}
