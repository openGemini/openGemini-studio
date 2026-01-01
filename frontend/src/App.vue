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
        @update-retention-policies="handleUpdateRetentionPolicies"
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
          :columns="resultColumns"
          :values="resultValues"
          :execution-time="executionTime"
          :error="error"
          :success="executionSuccess"
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
import { CloseConnect, ExecuteCommand, GetHistories } from '../wailsjs/go/main/App'
import { main } from '../wailsjs/go/models'

const { locale, t } = useI18n()
const themeComposable = useTheme()
const settingsComposable = useSettings()

const theme = themeComposable.theme
const toggleTheme = themeComposable.toggleTheme

const appSettings = settingsComposable.settings
const settingsVisible = ref(false)

const query = ref('SELECT * FROM measurement_name LIMIT 100')
const resultColumns = ref<string[]>([])
const resultValues = ref<any[][]>([])
const executionTime = ref(0)
const error = ref('')
const executionSuccess = ref(false)
const sidebarCollapsed = ref(false)
const historyVisible = ref(false)
const queryHistory = ref<QueryHistoryItem[]>([])
const selectedDatabase = ref('')
const selectedMeasurement = ref('')
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

const handleDisconnect = async () => {
  // Close all connections on backend
  for (const conn of connections.value) {
    if (conn.connected) {
      try {
        await CloseConnect(conn.id)
      } catch (error) {
        console.error(`Failed to close connection ${conn.id}:`, error)
      }
    }
  }

  // Update frontend state
  connections.value.forEach(c => {
    c.connected = false
    c.databases = []
    c.expanded = false
  })
  resultColumns.value = []
  resultValues.value = []
  executionSuccess.value = false
  error.value = ''
  selectedDatabase.value = ''
  selectedMeasurement.value = ''
  selectedRetentionPolicy.value = ''
}

const handleConnectToDatabase = (_connection: SavedConnection) => {
  // Connection is already established in ConnectionManager
  // Database list is already populated via DialConnect
  // Database selection is handled by watcher
}

const handleDisconnectFromDatabase = async (connectionId: string) => {
  const connection = connections.value.find(c => c.id === connectionId)
  if (connection) {
    try {
      // Call backend to close the connection
      await CloseConnect(connectionId)

      // Update frontend state
      connection.connected = false
      connection.databases = []
      connection.expanded = false
    } catch (error) {
      console.error(`Failed to close connection ${connectionId}:`, error)
      // Still update frontend state even if backend call fails
      connection.connected = false
      connection.databases = []
      connection.expanded = false
    }
  }
}

const handleSelectMeasurement = (data: { measurement: string; database: string; connection: SavedConnection }) => {
  query.value = `SELECT * FROM "${data.measurement}" LIMIT 100`
  selectedDatabase.value = data.database
  selectedMeasurement.value = data.measurement
}

const handleUpdateRetentionPolicies = (data: { policies: any[]; database: string }) => {
  // If this is the currently selected database, automatically select the first retention policy
  if (data.database === selectedDatabase.value) {
    if (data.policies.length > 0) {
      // Select the default policy if exists, otherwise select the first one
      const defaultPolicy = data.policies.find(p => p.isDefault)
      selectedRetentionPolicy.value = defaultPolicy ? defaultPolicy.name : data.policies[0].name
    } else {
      selectedRetentionPolicy.value = ''
    }
  }
}

const handleQueryUpdate = (newQuery: string) => {
  query.value = newQuery
}

const handleSelectHistoryQuery = (historyQuery: string) => {
  // Append to existing query instead of replacing
  if (query.value.trim()) {
    // If there's existing content, add a newline separator and append
    query.value = query.value.trim() + '\n\n' + historyQuery
  } else {
    // If empty, just set the query
    query.value = historyQuery
  }
}

const queryEditorRef = ref<InstanceType<typeof QueryEditor> | null>(null)

const executeQuery = async () => {
  // Get the actual query to execute from QueryEditor
  const queryToExecute = queryEditorRef.value?.getQueryToExecute() || query.value

  if (!queryToExecute.trim()) {
    error.value = 'No query to execute'
    return
  }

  // Get the active connection
  const activeConnection = connections.value.find(c => c.connected)
  if (!activeConnection) {
    error.value = 'No active connection. Please connect to a database first.'
    return
  }

  error.value = ''
  executionSuccess.value = false
  resultColumns.value = []
  resultValues.value = []

  try {
    const request: main.ExecuteRequest = {
      connect_name: activeConnection.id,
      database: selectedDatabase.value,
      retention_policy: selectedRetentionPolicy.value || '',
      measurement: selectedMeasurement.value || '',
      precision: selectedPrecision.value,
      command: queryToExecute
    }

    const response = await ExecuteCommand(request)
    // Use backend-measured execution time
    executionTime.value = response.execution_time || 0

    // Handle response based on NoContent flag
    if (response.no_content) {
      // Command executed successfully but has no results (like INSERT, DELETE, etc.)
      // Display as a table with Status and Message from backend
      executionSuccess.value = false
      resultColumns.value = [t('results.status'), t('results.message')]
      resultValues.value = [[t('results.success'), response.message || '']]
    } else {
      // Query returned results
      executionSuccess.value = false
      resultColumns.value = response.columns || []
      resultValues.value = response.values || []
    }

    // Reload query history from backend
    try {
      const histories = await GetHistories()
      queryHistory.value = histories.map(h => ({
        id: h.id,
        query: h.query,
        timestamp: new Date(h.timestamp),
        executionTime: h.execution_time,
        database: h.database || undefined,
        retentionPolicy: h.retention_policy || undefined,
        success: h.success,
        error: h.error || undefined
      }))
    } catch (err) {
      console.error('Failed to reload query history:', err)
    }
  } catch (e) {
    // For errors, execution time is set by backend
    executionTime.value = 0

    // Display error as a table with Status and Message columns
    // Extract error message from different error types
    let errorMessage: string
    if (e instanceof Error) {
      errorMessage = e.message
    } else if (typeof e === 'string') {
      errorMessage = e
    } else if (e && typeof e === 'object' && 'message' in e) {
      errorMessage = String((e as any).message)
    } else {
      errorMessage = String(e)
    }

    // Fallback to default message if empty
    if (!errorMessage || errorMessage.trim() === '') {
      errorMessage = 'Query execution failed'
    }

    resultColumns.value = [t('results.status'), t('results.message')]
    resultValues.value = [[t('results.failed'), errorMessage]]
    error.value = ''
    executionSuccess.value = false

    // Reload query history from backend (error is also saved by backend)
    try {
      const histories = await GetHistories()
      queryHistory.value = histories.map(h => ({
        id: h.id,
        query: h.query,
        timestamp: new Date(h.timestamp),
        executionTime: h.execution_time,
        database: h.database || undefined,
        retentionPolicy: h.retention_policy || undefined,
        success: h.success,
        error: h.error || undefined
      }))
    } catch (err) {
      console.error('Failed to reload query history:', err)
    }
  }
}

const exportResults = () => {
  if (resultValues.value.length === 0 || resultColumns.value.length === 0) return

  const csv = [
    resultColumns.value.join(','),
    ...resultValues.value.map(row => row.join(','))
  ].join('\n')

  const blob = new Blob([csv], { type: 'text/csv' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = 'query-results.csv'
  a.click()
  URL.revokeObjectURL(url)
}

// Watch availableDatabases and auto-select the first one if none is selected
watch(availableDatabases, (newDatabases) => {
  if (newDatabases.length > 0) {
    // If no database is selected or selected database no longer exists, select the first one
    if (!selectedDatabase.value || !newDatabases.find(db => db.name === selectedDatabase.value)) {
      selectedDatabase.value = newDatabases[0].name
    }
  } else {
    selectedDatabase.value = ''
    selectedRetentionPolicy.value = ''
  }
})

// Watch selectedDatabase and auto-select the first retention policy
watch(selectedDatabase, (newDatabase) => {
  if (newDatabase) {
    // Clear selected measurement when database changes
    selectedMeasurement.value = ''

    const db = availableDatabases.value.find(d => d.name === newDatabase)
    if (db && db.retentionPolicies && db.retentionPolicies.length > 0) {
      const defaultPolicy = db.retentionPolicies.find(p => p.isDefault)
      selectedRetentionPolicy.value = defaultPolicy ? defaultPolicy.name : db.retentionPolicies[0].name
    } else {
      selectedRetentionPolicy.value = ''
    }
  } else {
    selectedMeasurement.value = ''
    selectedRetentionPolicy.value = ''
  }
})

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

onMounted(async () => {
  themeComposable.initTheme()
  settingsComposable.initSettings()

  // Set initial locale
  locale.value = appSettings.value.language as 'en' | 'zh-CN'

  // Load query history from backend
  try {
    const histories = await GetHistories()
    queryHistory.value = histories.map(h => ({
      id: h.id,
      query: h.query,
      timestamp: new Date(h.timestamp),
      executionTime: h.execution_time,
      database: h.database || undefined,
      retentionPolicy: h.retention_policy || undefined,
      success: h.success,
      error: h.error || undefined
    }))
  } catch (err) {
    console.error('Failed to load query history:', err)
  }

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
