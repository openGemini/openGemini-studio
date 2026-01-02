<template>
  <div v-if="visible" class="settings-overlay" @click.self="$emit('close')">
    <div class="settings-modal">
      <div class="settings-header">
        <h2>{{ $t('settings.title') }}</h2>
        <button @click="$emit('close')" class="close-btn">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="18" y1="6" x2="6" y2="18"/>
            <line x1="6" y1="6" x2="18" y2="18"/>
          </svg>
        </button>
      </div>

      <div class="settings-content">
        <div class="setting-group">
          <label class="setting-label">{{ $t('settings.language') }}</label>
          <select v-model="localSettings.language" class="setting-select">
            <option value="en">{{ $t('settings.languageEnglish') }}</option>
            <option value="zh-CN">{{ $t('settings.languageChinese') }}</option>
          </select>
        </div>

        <div class="setting-group">
          <label class="setting-label">{{ $t('settings.themeMode') }}</label>
          <select v-model="localSettings.themeMode" class="setting-select">
            <option value="light">{{ $t('settings.themeLight') }}</option>
            <option value="dark">{{ $t('settings.themeDark') }}</option>
            <option value="system">{{ $t('settings.themeSystem') }}</option>
          </select>
        </div>

        <div class="setting-group">
          <label class="setting-label">{{ $t('settings.customFont') }}</label>
          <input
            v-model="localSettings.customFont"
            type="text"
            class="setting-input"
            :placeholder="$t('settings.customFontPlaceholder')"
          />
          <span class="setting-hint">{{ $t('settings.customFontHint') }}</span>
        </div>

        <div class="setting-group">
          <label class="setting-label">{{ $t('settings.maxHistoryCount') }}</label>
          <input
            v-model.number="localSettings.maxHistoryCount"
            type="number"
            min="10"
            max="500"
            class="setting-input"
          />
          <span class="setting-hint">{{ $t('settings.historyCountHint') }}</span>
        </div>

        <div class="setting-group">
          <label class="setting-label">{{ $t('settings.dataDirectory') }}</label>
          <input
            v-model="localSettings.dataDirectory"
            type="text"
            class="setting-input"
            placeholder="./data"
          />
          <span class="setting-hint">{{ $t('settings.dataDirectoryHint') }}</span>
        </div>

        <div class="setting-group">
          <label class="setting-label setting-checkbox-label">
            <input
              v-model="localSettings.debug"
              type="checkbox"
              class="setting-checkbox"
            />
            <span>{{ $t('settings.debug') }}</span>
          </label>
          <span class="setting-hint">{{ $t('settings.debugHint') }}</span>
        </div>
      </div>

      <div class="settings-footer">
        <button @click="handleReset" class="btn btn-secondary">{{ $t('settings.resetToDefaults') }}</button>
        <button @click="handleSave" class="btn btn-primary">{{ $t('settings.save') }}</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import type { AppSettings } from '../types'

const props = defineProps<{
  visible: boolean
  settings: AppSettings
}>()

const emit = defineEmits<{
  close: []
  save: [settings: AppSettings]
  reset: []
}>()

const localSettings = ref<AppSettings>({ ...props.settings })

watch(() => props.settings, (newSettings) => {
  localSettings.value = { ...newSettings }
}, { deep: true })

const handleSave = () => {
  emit('save', localSettings.value)
  emit('close')
}

const handleReset = () => {
  emit('reset')
  emit('close')
}
</script>

<style scoped>
.settings-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.settings-modal {
  background: var(--bg-secondary);
  border-radius: 8px;
  width: 90%;
  max-width: 600px;
  max-height: 80vh;
  display: flex;
  flex-direction: column;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
}

.settings-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid var(--border-color);
}

.settings-header h2 {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
  color: var(--text-primary);
}

.close-btn {
  background: transparent;
  border: none;
  color: var(--text-secondary);
  cursor: pointer;
  padding: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
  transition: all 0.2s;
}

.close-btn:hover {
  background: var(--bg-hover);
  color: var(--text-primary);
}

.settings-content {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
}

.setting-group {
  margin-bottom: 24px;
}

.setting-label {
  display: block;
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
  margin-bottom: 8px;
}

.setting-select,
.setting-input {
  width: 100%;
  padding: 10px 12px;
  background-color: var(--bg-tertiary);
  border: 1px solid var(--border-color);
  border-radius: 6px;
  color: var(--text-primary) !important;
  font-size: 14px;
  transition: all 0.2s;
  -webkit-appearance: none;
  -moz-appearance: none;
  appearance: none;
}

.setting-select {
  background-image: url("data:image/svg+xml;charset=UTF-8,%3csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='%23e4e6eb' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3e%3cpolyline points='6 9 12 15 18 9'%3e%3c/polyline%3e%3c/svg%3e");
  background-repeat: no-repeat;
  background-position: right 8px center;
  background-size: 16px;
  padding-right: 32px;
}

:root[data-theme="light"] .setting-select {
  background-image: url("data:image/svg+xml;charset=UTF-8,%3csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='%231f2937' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3e%3cpolyline points='6 9 12 15 18 9'%3e%3c/polyline%3e%3c/svg%3e");
}

.setting-select option {
  background: var(--bg-tertiary);
  color: var(--text-primary);
}

.setting-select:focus,
.setting-input:focus {
  outline: none;
  border-color: #3b82f6;
}

.setting-hint {
  display: block;
  font-size: 12px;
  color: var(--text-secondary);
  margin-top: 6px;
}

.settings-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 20px;
  border-top: 1px solid var(--border-color);
}

.btn {
  padding: 10px 20px;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-primary {
  background: #3b82f6;
  color: white;
}

.btn-primary:hover {
  background: #2563eb;
}

.btn-secondary {
  background: var(--bg-tertiary);
  color: var(--text-primary);
  border: 1px solid var(--border-color);
}

.btn-secondary:hover {
  background: var(--bg-hover);
}

.setting-checkbox-label {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  user-select: none;
}

.setting-checkbox {
  width: 18px;
  height: 18px;
  cursor: pointer;
  accent-color: #3b82f6;
}
</style>
