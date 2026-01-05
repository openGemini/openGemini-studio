<template>
  <div class="toolbar">
    <div class="toolbar-left">
      <button @click="$emit('toggleSidebar')" class="btn btn-icon" :title="$t('toolbar.toggleSidebar')">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <rect x="3" y="3" width="18" height="18" rx="2"/>
          <line x1="9" y1="3" x2="9" y2="21"/>
        </svg>
      </button>
      <div v-if="isConnected" class="connection-status">
        <span class="status-dot"></span>
        <span class="status-text">{{ $t('common.connected') }}</span>
      </div>
    </div>

    <div class="toolbar-right">
      <button @click="$emit('toggleHistory')" class="btn btn-icon" :title="$t('toolbar.history')">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="10"/>
          <polyline points="12 6 12 12 16 14"/>
        </svg>
      </button>
      <button @click="$emit('toggleTheme')" class="btn btn-icon" :title="$t('toolbar.toggleTheme')">
        <svg v-if="isDark" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="5"/>
          <line x1="12" y1="1" x2="12" y2="3"/>
          <line x1="12" y1="21" x2="12" y2="23"/>
          <line x1="4.22" y1="4.22" x2="5.64" y2="5.64"/>
          <line x1="18.36" y1="18.36" x2="19.78" y2="19.78"/>
          <line x1="1" y1="12" x2="3" y2="12"/>
          <line x1="21" y1="12" x2="23" y2="12"/>
          <line x1="4.22" y1="19.78" x2="5.64" y2="18.36"/>
          <line x1="18.36" y1="5.64" x2="19.78" y2="4.22"/>
        </svg>
        <svg v-else width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/>
        </svg>
      </button>
      <button @click="$emit('toggleSettings')" class="btn btn-icon" :title="$t('toolbar.settings')">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="3"/>
          <path d="M12 1v6m0 6v6M5.64 5.64l4.24 4.24m4.24 4.24l4.24 4.24M1 12h6m6 0h6M5.64 18.36l4.24-4.24m4.24-4.24l4.24-4.24"/>
        </svg>
      </button>
      <button @click="$emit('toggleAbout')" class="btn btn-icon" :title="$t('toolbar.about')">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <circle cx="12" cy="12" r="10"/>
          <path d="M9.09 9a3 3 0 0 1 5.83 1c0 2-3 3-3 3"/>
          <line x1="12" y1="17" x2="12.01" y2="17"/>
        </svg>
      </button>

      <template v-if="isConnected">
        <button @click="$emit('executeQuery')" class="btn btn-success">
          {{ $t('toolbar.execute') }}
        </button>
        <button @click="$emit('exportResults')" class="btn btn-secondary">
          {{ $t('toolbar.export') }}
        </button>
        <button @click="$emit('disconnect')" class="btn btn-danger">
          {{ $t('toolbar.disconnect') }}
        </button>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  isConnected: boolean
  isDark: boolean
}>()

defineEmits<{
  connect: []
  disconnect: []
  executeQuery: []
  exportResults: []
  toggleTheme: []
  toggleSidebar: []
  toggleHistory: []
  toggleSettings: []
  toggleAbout: []
}>()
</script>

<style scoped>
.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 20px;
  background: var(--bg-secondary);
  border-bottom: 1px solid var(--border-color);
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 20px;
}

.title {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
}

.connection-status {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 4px 12px;
  background: #1a3a1a;
  border-radius: 12px;
}

.status-dot {
  width: 8px;
  height: 8px;
  background: #4ade80;
  border-radius: 50%;
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

.status-text {
  font-size: 13px;
  color: #4ade80;
}

.toolbar-right {
  display: flex;
  gap: 10px;
}

.btn {
  padding: 8px 16px;
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

.btn-success {
  background: #10b981;
  color: white;
}

.btn-success:hover {
  background: #059669;
}

.btn-secondary {
  background: #6b7280;
  color: white;
}

.btn-secondary:hover {
  background: #4b5563;
}

.btn-danger {
  background: #ef4444;
  color: white;
}

.btn-danger:hover {
  background: #dc2626;
}

.btn-icon {
  padding: 8px 12px;
  background: var(--bg-tertiary);
  color: var(--text-primary);
  font-size: 16px;
}

.btn-icon:hover {
  background: var(--bg-hover);
}
</style>
