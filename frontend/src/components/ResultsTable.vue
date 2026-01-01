<template>
  <div class="results-table">
    <div class="results-header">
      <span class="header-title">{{ $t('results.title') }}</span>
      <span v-if="executionTime > 0" class="execution-time">
        {{ $t('results.executedIn', { time: executionTime.toFixed(2) }) }}
      </span>
    </div>

    <div v-if="error" class="error-message">
      ❌ {{ error }}
    </div>

    <div v-else-if="success" class="success-message">
      ✅ {{ $t('results.commandSuccess') }}
    </div>

    <div v-else-if="!columns || columns.length === 0" class="empty-state">
      {{ $t('results.emptyState') }}
    </div>

    <div v-else class="table-container">
      <table ref="tableRef">
        <thead>
          <tr>
            <th
              v-for="(column, index) in columns"
              :key="column"
              :style="{ width: columnWidths[index] + 'px' }"
              class="resizable-th"
            >
              <span class="column-content">{{ column }}</span>
              <div
                v-if="index < columns.length - 1"
                class="column-resizer"
                @mousedown="startResize($event, index)"
              ></div>
            </th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(row, rowIndex) in values" :key="rowIndex">
            <td
              v-for="(cell, cellIndex) in row"
              :key="cellIndex"
              :style="{ width: columnWidths[cellIndex] + 'px' }"
            >
              {{ formatCell(cell) }}
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'

const props = defineProps<{
  columns: string[]
  values: any[][]
  executionTime: number
  error: string
  success: boolean
}>()

const tableRef = ref<HTMLTableElement | null>(null)
const columnWidths = ref<number[]>([])
const resizingColumn = ref<number | null>(null)
const startX = ref(0)
const startWidth = ref(0)

// Initialize column widths
watch(() => props.columns, (newColumns) => {
  if (newColumns && newColumns.length > 0) {
    // Initialize with default width of 150px for each column
    columnWidths.value = newColumns.map(() => 150)
  }
}, { immediate: true })

const formatCell = (value: any): string => {
  if (value === null || value === undefined) {
    return ''
  }
  if (typeof value === 'object') {
    return JSON.stringify(value)
  }
  return String(value)
}

const startResize = (event: MouseEvent, columnIndex: number) => {
  resizingColumn.value = columnIndex
  startX.value = event.pageX
  startWidth.value = columnWidths.value[columnIndex]

  document.addEventListener('mousemove', handleMouseMove)
  document.addEventListener('mouseup', stopResize)
  event.preventDefault()
}

const handleMouseMove = (event: MouseEvent) => {
  if (resizingColumn.value === null) return

  const diff = event.pageX - startX.value
  const newWidth = Math.max(50, startWidth.value + diff) // Minimum width of 50px

  columnWidths.value[resizingColumn.value] = newWidth
}

const stopResize = () => {
  resizingColumn.value = null
  document.removeEventListener('mousemove', handleMouseMove)
  document.removeEventListener('mouseup', stopResize)
}
</script>

<style scoped>
.results-table {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.results-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: var(--bg-secondary);
  border-bottom: 1px solid var(--border-color);
}

.header-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
}

.execution-time {
  font-size: 12px;
  color: #10b981;
}

.error-message {
  padding: 20px;
  color: #ef4444;
  text-align: center;
}

.success-message {
  padding: 20px;
  color: #10b981;
  text-align: center;
  font-size: 14px;
}

.empty-state {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-secondary);
  font-size: 14px;
}

.table-container {
  flex: 1;
  overflow: auto;
  background: var(--bg-primary);
  position: relative;
}

table {
  border-collapse: collapse;
  table-layout: fixed;
  min-width: 100%;
}

thead {
  position: sticky;
  top: 0;
  background: var(--bg-secondary);
  z-index: 2;
}

.resizable-th {
  padding: 12px 16px;
  text-align: left;
  font-size: 13px;
  font-weight: 600;
  color: var(--text-secondary);
  background: var(--bg-secondary);
  border-bottom: 1px solid var(--border-color);
  position: relative;
  user-select: none;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.resizable-th:last-child {
  border-right: none;
}

.column-content {
  display: inline-block;
  max-width: 100%;
  overflow: hidden;
  text-overflow: ellipsis;
}

.column-resizer {
  position: absolute;
  top: 0;
  right: 0;
  width: 5px;
  height: 100%;
  cursor: col-resize;
  user-select: none;
  border-right: 1px solid var(--border-color);
  background: transparent;
  z-index: 10;
}

.column-resizer:hover {
  border-right: 2px solid var(--accent-color);
}

.column-resizer:active {
  border-right: 2px solid var(--accent-color);
}

td {
  padding: 10px 16px;
  font-size: 13px;
  color: var(--text-primary);
  border-bottom: 1px solid var(--border-color);
  border-right: 1px solid var(--border-color);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

td:last-child {
  border-right: none;
}

tbody tr:hover {
  background: var(--bg-hover);
}

/* 自定义滚动条样式 */
.table-container::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

.table-container::-webkit-scrollbar-track {
  background: var(--bg-secondary);
}

.table-container::-webkit-scrollbar-thumb {
  background: var(--text-secondary);
  border-radius: 4px;
}

.table-container::-webkit-scrollbar-thumb:hover {
  background: var(--text-primary);
}
</style>
