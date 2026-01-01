<template>
  <div class="query-editor">
    <div class="editor-header">
      <div class="header-left">
        <span class="header-title">{{ $t('query.editor') }}</span>
      </div>
      <div class="header-controls">
        <div class="control-group">
          <label>{{ $t('query.precision') }}:</label>
          <select
            :value="selectedPrecision"
            @change="$emit('update:selectedPrecision', ($event.target as HTMLSelectElement).value)"
            class="control-select control-select-small"
            :disabled="disabled"
          >
            <option value="ns">ns (nanoseconds)</option>
            <option value="u">u (microseconds)</option>
            <option value="ms">ms (milliseconds)</option>
            <option value="s">s (seconds)</option>
            <option value="m">m (minutes)</option>
            <option value="h">h (hours)</option>
            <option value="rfc3339">RFC3339</option>
          </select>
        </div>
        <div class="control-group">
          <label>{{ $t('query.database') }}:</label>
          <select
            :value="selectedDatabase"
            @change="$emit('update:selectedDatabase', ($event.target as HTMLSelectElement).value)"
            class="control-select"
            :disabled="disabled"
          >
            <option v-for="db in databases" :key="db.name" :value="db.name">
              {{ db.name }}
            </option>
          </select>
        </div>
        <div class="control-group" v-if="selectedDatabase && retentionPolicies.length > 0">
          <label>{{ $t('query.retentionPolicy') }}:</label>
          <select
            :value="selectedRetentionPolicy"
            @change="$emit('update:selectedRetentionPolicy', ($event.target as HTMLSelectElement).value)"
            class="control-select"
            :disabled="disabled"
          >
            <option v-for="rp in retentionPolicies" :key="rp.name" :value="rp.name">
              {{ rp.name }}({{ rp.duration }})
            </option>
          </select>
        </div>
        <span class="char-count">{{ charCount }} {{ $t('query.characters') }}</span>
      </div>
    </div>
    <textarea 
      ref="textareaRef"
      :value="query"
      @input="$emit('update:query', ($event.target as HTMLTextAreaElement).value)"
      class="editor-textarea"
      :class="{ 'editor-disabled': disabled }"
      :placeholder="$t('query.placeholder')"
      spellcheck="false"
      :disabled="disabled"
    ></textarea>
  </div>
</template>

<script setup lang="ts">
import type { Database, RetentionPolicy } from '../types'
import { computed, ref } from 'vue'

const props = defineProps<{
  query: string
  charCount: number
  selectedDatabase: string
  selectedRetentionPolicy: string
  selectedPrecision: string
  databases: Database[]
  disabled?: boolean
}>()

defineEmits<{
  'update:query': [query: string]
  'update:selectedDatabase': [database: string]
  'update:selectedRetentionPolicy': [retentionPolicy: string]
  'update:selectedPrecision': [precision: string]
}>()

const retentionPolicies = computed((): RetentionPolicy[] => {
  const db = props.databases.find(d => d.name === props.selectedDatabase)
  return db?.retentionPolicies || []
})

const textareaRef = ref<HTMLTextAreaElement | null>(null)

const getQueryToExecute = (): string => {
  if (!textareaRef.value) return props.query

  const textarea = textareaRef.value
  const selectionStart = textarea.selectionStart
  const selectionEnd = textarea.selectionEnd

  if (selectionStart !== selectionEnd) {
    return props.query.substring(selectionStart, selectionEnd).trim()
  }

  const lines = props.query.split('\n')
  const textBeforeCursor = props.query.substring(0, selectionStart)
  const currentLineIndex = textBeforeCursor.split('\n').length - 1
  
  return lines[currentLineIndex]?.trim() || props.query
}

defineExpose({
  getQueryToExecute
})
</script>

<style scoped>
.query-editor {
  display: flex;
  flex-direction: column;
  height: 40%;
  border-bottom: 1px solid var(--border-color);
}

.editor-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: var(--bg-secondary);
  border-bottom: 1px solid var(--border-color);
  gap: 16px;
}

.header-left {
  display: flex;
  align-items: center;
}

.header-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
}

.header-controls {
  display: flex;
  align-items: center;
  gap: 16px;
  flex: 1;
  justify-content: flex-end;
}

.control-group {
  display: flex;
  align-items: center;
  gap: 8px;
}

.control-group label {
  font-size: 12px;
  color: var(--text-secondary);
  white-space: nowrap;
}

.control-select {
  padding: 6px 12px;
  background-color: var(--bg-primary);
  border: 1px solid var(--border-color);
  border-radius: 4px;
  color: var(--text-primary) !important;
  font-size: 12px;
  cursor: pointer;
  min-width: 150px;
  transition: all 0.2s;
  -webkit-appearance: none;
  -moz-appearance: none;
  appearance: none;
  background-image: url("data:image/svg+xml;charset=UTF-8,%3csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='%23e4e6eb' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3e%3cpolyline points='6 9 12 15 18 9'%3e%3c/polyline%3e%3c/svg%3e");
  background-repeat: no-repeat;
  background-position: right 6px center;
  background-size: 12px;
  padding-right: 28px;
}

:root[data-theme="light"] .control-select {
  background-image: url("data:image/svg+xml;charset=UTF-8,%3csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='%231f2937' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3e%3cpolyline points='6 9 12 15 18 9'%3e%3c/polyline%3e%3c/svg%3e");
}

.control-select option {
  background: var(--bg-primary);
  color: var(--text-primary);
}

.control-select-small {
  min-width: 120px;
}

.control-select:hover {
  border-color: var(--accent-color);
}

.control-select:focus {
  outline: none;
  border-color: var(--accent-color);
}

.char-count {
  font-size: 12px;
  color: var(--text-secondary);
  white-space: nowrap;
}

.editor-textarea {
  flex: 1;
  padding: 16px;
  background: var(--bg-primary);
  border: none;
  color: var(--text-primary);
  font-family: 'JetBrains Mono', 'LXGW WenKai Lite', 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 14px;
  line-height: 1.6;
  resize: none;
  outline: none;
}

.editor-textarea::placeholder {
  color: var(--text-secondary);
}

.editor-textarea:disabled,
.editor-disabled {
  opacity: 0.5;
  cursor: not-allowed;
  background: var(--bg-secondary);
}

.control-select:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>
