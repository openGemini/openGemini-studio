<template>
  <div class="app-container">
    <Toolbar 
      :is-connected="hasActiveConnection"
      :is-dark="theme === 'dark'"
      @connect="handleConnect"
      @disconnect="handleDisconnect"
      @execute-query="executeQuery"
      @export-results="exportResults"
      @toggle-theme="toggleTheme"
      @toggle-sidebar="sidebarCollapsed = !sidebarCollapsed"
      @toggle-history="historyVisible = !historyVisible"
      @toggle-settings="settingsVisible = true"
    />
    
    <div class="main-content">
      <ConnectionManager
        :connections="connections"
        :collapsed="sidebarCollapsed"
        @update:connections="connections = $event"
        @select-measurement="handleSelectMeasurement"
        @connect="handleConnectToDatabase"
        @disconnect="handleDisconnectFromDatabase"
      />
      
      <div class="editor-results-container">
        <QueryEditor
          ref="queryEditorRef"
          :query="query"
          :char-count="charCount"
          :selected-database="selectedDatabase"
          :selected-retention-policy="selectedRetentionPolicy"
          :selected-precision="selectedPrecision"
          :databases="availableDatabases"
          :disabled="!hasActiveConnection"
          @update:query="handleQueryUpdate"
          @update:selected-database="selectedDatabase = $event"
          @update:selected-retention-policy="selectedRetentionPolicy = $event"
          @update:selected-precision="selectedPrecision = $event"
        />
        
        <ResultsTable
          :results="results"
          :execution-time="executionTime"
          :error="error"
        />
      </div>

      <QueryHistory
        v-if="historyVisible"
        :history="queryHistory"
        @select="handleSelectHistoryQuery"
        @clear="queryHistory = []"
      />
    </div>

    <!-- 添加设置弹窗 -->
    <Settings
      :visible="settingsVisible"
      :settings="appSettings"
      @close="settingsVisible = false"
      @save="handleSaveSettings"
      @reset="handleResetSettings"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import Toolbar from './components/Toolbar.vue'
import ConnectionManager from './components/ConnectionManager.vue'
import QueryEditor from './components/QueryEditor.vue'
import ResultsTable from './components/ResultsTable.vue'
import QueryHistory from './components/QueryHistory.vue'
import Settings from './components/Settings.vue'
import { useTheme } from './composables/useTheme'
import { useSettings } from './composables/useSettings'
import type { SavedConnection, QueryHistoryItem, Database, AppSettings } from './types'

const { locale } = useI18n()
const themeComposable = useTheme()
const settingsComposable = useSettings()

const theme = themeComposable.theme
const toggleTheme = themeComposable.toggleTheme

const appSettings = settingsComposable.settings
const settingsVisible = ref(false)

const query = ref('SELECT * FROM measurement_name LIMIT 100')
const results = ref<any[]>([])
const executionTime = ref(0)
const error = ref('')
const sidebarCollapsed = ref(false)
const historyVisible = ref(false)
const queryHistory = ref<QueryHistoryItem[]>([])
const selectedDatabase = ref('')
const selectedRetentionPolicy = ref('')
const selectedPrecision = ref('ms')

const connections = ref<SavedConnection[]>([])

const charCount = computed(() => query.value.length)

const hasActiveConnection = computed(() => 
  connections.value.some(c => c.connected)
)

const availableDatabases = computed(() => {
  const allDatabases: Database[] = []
  connections.value.forEach(conn => {
    if (conn.connected) {
      allDatabases.push(...conn.databases)
    }
  })
  return allDatabases
})

const handleConnect = () => {
  // Legacy connect handler - not used anymore
}

const handleDisconnect = () => {
  connections.value.forEach(c => {
    c.connected = false
    c.databases = []
  })
  results.value = []
  error.value = ''
  selectedDatabase.value = ''
  selectedRetentionPolicy.value = ''
}

const handleConnectToDatabase = (_connection: SavedConnection) => {
  // Connection is already established in ConnectionManager
  // Database list is already populated via DialConnect
  // Just need to select a default database if none selected
  if (availableDatabases.value.length > 0 && !selectedDatabase.value) {
    selectedDatabase.value = availableDatabases.value[0].name
  }
}

const handleDisconnectFromDatabase = (connectionId: string) => {
  const connection = connections.value.find(c => c.id === connectionId)
  if (connection) {
    connection.connected = false
    connection.databases = []
    connection.expanded = false
  }
}

const handleSelectMeasurement = (data: { measurement: string; database: string; connection: SavedConnection }) => {
  query.value = `SELECT * FROM "${data.measurement}" LIMIT 100`
  selectedDatabase.value = data.database
}

const handleQueryUpdate = (newQuery: string) => {
  query.value = newQuery
}

const handleSelectHistoryQuery = (historyQuery: string) => {
  query.value = historyQuery
}

const queryEditorRef = ref<InstanceType<typeof QueryEditor> | null>(null)

const executeQuery = () => {
  // Get the actual query to execute from QueryEditor
  const queryToExecute = queryEditorRef.value?.getQueryToExecute() || query.value
  
  if (!queryToExecute.trim()) {
    error.value = 'No query to execute'
    return
  }
  
  const startTime = performance.now()
  error.value = ''
  
  try {
    results.value = [
      { time: '2024-01-15T10:30:00Z', host: 'server-01', cpu: '45.2', memory: '2048', status: 'active' },
      { time: '2024-01-15T10:31:00Z', host: 'server-01', cpu: '48.7', memory: '2056', status: 'active' },
      { time: '2024-01-15T10:32:00Z', host: 'server-02', cpu: '52.1', memory: '3072', status: 'active' },
      { time: '2024-01-15T10:33:00Z', host: 'server-02', cpu: '49.8', memory: '3068', status: 'active' },
      { time: '2024-01-15T10:34:00Z', host: 'server-03', cpu: '41.3', memory: '1536', status: 'active' }
    ]
    
    const endTime = performance.now()
    executionTime.value = endTime - startTime
    
    queryHistory.value.unshift({
      id: Date.now().toString(),
      query: queryToExecute,
      timestamp: new Date(),
      executionTime: executionTime.value,
      database: selectedDatabase.value || undefined,
      retentionPolicy: selectedRetentionPolicy.value || undefined,
      success: true
    })
    
    if (queryHistory.value.length > appSettings.value.maxHistoryCount) {
      queryHistory.value = queryHistory.value.slice(0, appSettings.value.maxHistoryCount)
    }
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Query execution failed'
    
    queryHistory.value.unshift({
      id: Date.now().toString(),
      query: queryToExecute,
      timestamp: new Date(),
      executionTime: 0,
      database: selectedDatabase.value || undefined,
      retentionPolicy: selectedRetentionPolicy.value || undefined,
      success: false,
      error: error.value
    })
  }
}

const exportResults = () => {
  if (results.value.length === 0) return
  
  const headers = Object.keys(results.value[0])
  const csv = [
    headers.join(','),
    ...results.value.map(row => headers.map(h => row[h]).join(','))
  ].join('\n')
  
  const blob = new Blob([csv], { type: 'text/csv' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = 'query-results.csv'
  a.click()
  URL.revokeObjectURL(url)
}

watch(() => appSettings.value.themeMode, (mode) => {
  if (mode === 'system') {
    // 跟随系统主题
    if (typeof window !== 'undefined' && window.matchMedia) {
      const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
      theme.value = prefersDark ? 'dark' : 'light'
    }
  } else {
    theme.value = mode
  }
}, { immediate: true })

// Watch language changes
watch(() => appSettings.value.language, (newLanguage) => {
  locale.value = newLanguage as 'en' | 'zh-CN'
}, { immediate: true })

// Watch custom font changes
watch(() => appSettings.value.customFont, (newFont) => {
  if (typeof document !== 'undefined') {
    const defaultFonts = '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif'
    if (newFont && newFont.trim()) {
      document.body.style.fontFamily = `${newFont}, ${defaultFonts}`
    } else {
      document.body.style.fontFamily = defaultFonts
    }
  }
}, { immediate: true })

onMounted(() => {
  themeComposable.initTheme()
  settingsComposable.initSettings()

  // Set initial locale
  locale.value = appSettings.value.language as 'en' | 'zh-CN'

  if (typeof window !== 'undefined' && window.matchMedia) {
    const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
    const handleChange = (e: MediaQueryListEvent) => {
      if (appSettings.value.themeMode === 'system') {
        theme.value = e.matches ? 'dark' : 'light'
      }
    }
    mediaQuery.addEventListener('change', handleChange)
  }
})

const handleSaveSettings = (newSettings: AppSettings) => {
  settingsComposable.updateSettings(newSettings)
  
  // 如果最大历史记录数改变了，裁剪历史记录
  if (queryHistory.value.length > newSettings.maxHistoryCount) {
    queryHistory.value = queryHistory.value.slice(0, newSettings.maxHistoryCount)
  }
}

const handleResetSettings = () => {
  settingsComposable.resetSettings()
}
</script>

<style scoped>
.app-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background: var(--bg-primary);
}

.main-content {
  display: flex;
  flex: 1;
  overflow: hidden;
}

.editor-results-container {
  display: flex;
  flex-direction: column;
  flex: 1;
  overflow: hidden;
}
</style>
