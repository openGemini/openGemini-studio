<template>
  <div 
    v-if="!collapsed"
    class="connection-manager" 
    :style="{ width: sidebarWidth + 'px' }"
  >
    <div class="manager-header">
      <h3>{{ $t('connection.title') }}</h3>
      <!-- Replace + button with SVG icon -->
      <button @click="showAddConnection = true" class="btn-add" :title="$t('connection.addConnection')">
        <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
          <line x1="12" y1="5" x2="12" y2="19" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
          <line x1="5" y1="12" x2="19" y2="12" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
        </svg>
      </button>
    </div>

    <div class="connections-list">
      <div v-for="conn in connections" :key="conn.id" class="connection-item">
        <div 
          class="connection-header" 
          @click="conn.connected ? toggleConnection(conn) : null"
          @contextmenu.prevent="showContextMenu($event, conn)"
          :class="{ 'not-connected': !conn.connected }"
        >
          <!-- Replace expand/collapse arrow with SVG chevron -->
          <svg v-if="conn.connected && !conn.expanded" class="icon-chevron" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M9 18l6-6-6-6" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
          <svg v-else-if="conn.connected && conn.expanded" class="icon-chevron" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M6 9l6 6 6-6" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
          <span v-else class="icon-placeholder"></span>

          <svg v-if="conn.connected" class="status-icon connected" viewBox="0 0 16 16" fill="none" xmlns="http://www.w3.org/2000/svg">
            <circle cx="8" cy="8" r="6" stroke="currentColor" stroke-width="2"/>
            <circle cx="8" cy="8" r="3" fill="currentColor"/>
          </svg>
          <svg v-else class="status-icon disconnected" viewBox="0 0 16 16" fill="none" xmlns="http://www.w3.org/2000/svg">
            <circle cx="8" cy="8" r="6" stroke="currentColor" stroke-width="2"/>
          </svg>
          <span class="name">{{ conn.name }}</span>
          <div class="connection-actions" @click.stop>
            <button
              v-if="!conn.connected"
              @click="connectToDatabase(conn)"
              class="btn-action"
              :title="$t('common.connect')"
            >
              <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M13 2L3 14h8l-1 8 10-12h-8l1-8z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              </svg>
            </button>
            <button
              v-else
              @click="disconnectFromDatabase(conn)"
              class="btn-action"
              :title="$t('common.disconnect')"
            >
              <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <line x1="18" y1="6" x2="6" y2="18" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
                <line x1="6" y1="6" x2="18" y2="18" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
              </svg>
            </button>
          </div>
        </div>

        <div v-if="conn.expanded && conn.connected" class="connection-content">
          <div v-for="db in conn.databases" :key="db.name" class="database-item">
            <div class="database-header" @dblclick.prevent="loadDatabaseMetadata(db, conn)">
              <!-- Replace database expand/collapse with SVG chevron -->
              <svg v-if="!db.expanded" class="icon-chevron indent" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M9 18l6-6-6-6" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              </svg>
              <svg v-else class="icon-chevron indent" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M6 9l6 6 6-6" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              </svg>
              <!-- Replace database emoji with SVG icon -->
              <svg class="icon-db" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <ellipse cx="12" cy="5" rx="9" ry="3" stroke="currentColor" stroke-width="2"/>
                <path d="M3 5v14c0 1.66 4 3 9 3s9-1.34 9-3V5" stroke="currentColor" stroke-width="2"/>
                <path d="M3 12c0 1.66 4 3 9 3s9-1.34 9-3" stroke="currentColor" stroke-width="2"/>
              </svg>
              <span class="name">{{ db.name }}</span>
            </div>

            <div v-if="db.expanded" class="measurements-list">
              <div 
                v-for="measurement in db.measurements" 
                :key="measurement.name"
                class="measurement-item"
                @click="selectMeasurement(measurement.name, db.name, conn)"
              >
                <!-- Replace table emoji with SVG icon -->
                <svg class="icon-table double-indent" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <rect x="3" y="3" width="18" height="18" rx="2" stroke="currentColor" stroke-width="2"/>
                  <line x1="3" y1="9" x2="21" y2="9" stroke="currentColor" stroke-width="2"/>
                  <line x1="3" y1="15" x2="21" y2="15" stroke="currentColor" stroke-width="2"/>
                  <line x1="12" y1="9" x2="12" y2="21" stroke="currentColor" stroke-width="2"/>
                </svg>
                <span class="name">{{ measurement.name }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Context Menu -->
    <div 
      v-if="contextMenu.visible" 
      class="context-menu"
      :style="{ top: contextMenu.y + 'px', left: contextMenu.x + 'px' }"
      @click="hideContextMenu"
    >
      <!-- Replace emoji icons in context menu with SVG icons -->
      <div class="context-menu-item" @click="editConnection(contextMenu.connection!)">
        <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
          <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>
        <span>{{ $t('common.edit') }}</span>
      </div>
      <div class="context-menu-item" @click="deleteConnection(contextMenu.connection!.id)">
        <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
          <path d="M3 6h18M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>
        <span>{{ $t('common.delete') }}</span>
      </div>
      <div class="context-menu-item" @click="refreshConnection(contextMenu.connection!)">
        <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
          <path d="M21.5 2v6h-6M2.5 22v-6h6M2 11.5a10 10 0 0 1 18.8-4.3M22 12.5a10 10 0 0 1-18.8 4.2" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>
        <span>{{ $t('connection.refresh') }}</span>
      </div>
      <div
        v-if="contextMenu.connection?.connected"
        class="context-menu-item"
        @click="disconnectFromDatabase(contextMenu.connection!)"
      >
        <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
          <line x1="18" y1="6" x2="6" y2="18" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
          <line x1="6" y1="6" x2="18" y2="18" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
        </svg>
        <span>{{ $t('common.disconnect') }}</span>
      </div>
      <div
        v-else
        class="context-menu-item"
        @click="connectToDatabase(contextMenu.connection!)"
      >
        <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
          <path d="M13 2L3 14h8l-1 8 10-12h-8l1-8z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>
        <span>{{ $t('common.connect') }}</span>
      </div>
    </div>

    <!-- Add Connection Modal -->
    <div v-if="showAddConnection" class="modal-overlay" @click="showAddConnection = false">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>{{ $t('connection.newConnection') }}</h3>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label>{{ $t('connection.connectionName') }}</label>
            <input v-model="newConnection.name" type="text" placeholder="My Database" />
          </div>
          <div class="form-group">
            <label>
              {{ $t('connection.address') }}
              <span class="help-icon" :title="$t('connection.addressHelp')">
                <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg" width="14" height="14">
                  <circle cx="12" cy="12" r="9" stroke="currentColor" stroke-width="2"/>
                  <path d="M12 16v-4M12 8h.01" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
                </svg>
              </span>
            </label>
            <div class="address-input-group">
              <select v-model="newConnection.protocol" class="protocol-select">
                <option value="http">HTTP</option>
                <option value="https">HTTPS</option>
              </select>
              <input v-model="newConnection.address" type="text" placeholder="localhost:8086" class="address-input" />
            </div>
          </div>
          <div v-if="newConnection.protocol === 'https'" class="https-options">
            <div class="form-group">
              <label>{{ $t('connection.caCertificate') }}</label>
              <div class="file-input-group">
                <input v-model="newConnection.caCert" type="text" placeholder="/path/to/ca.crt" />
                <button type="button" class="btn-browse" @click="selectCaCertFile('new')">
                  {{ $t('connection.browse') }}
                </button>
              </div>
            </div>
            <div class="form-group">
              <label>{{ $t('connection.clientCertificate') }}</label>
              <div class="file-input-group">
                <input v-model="newConnection.clientCert" type="text" placeholder="/path/to/client.crt" />
                <button type="button" class="btn-browse" @click="selectClientCertFile('new')">
                  {{ $t('connection.browse') }}
                </button>
              </div>
            </div>
            <div class="form-group">
              <label>{{ $t('connection.clientKey') }}</label>
              <div class="file-input-group">
                <input v-model="newConnection.clientKey" type="text" placeholder="/path/to/client.key" />
                <button type="button" class="btn-browse" @click="selectClientKeyFile('new')">
                  {{ $t('connection.browse') }}
                </button>
              </div>
            </div>
            <div class="form-group checkbox-group">
              <label>
                <input v-model="newConnection.insecureTls" type="checkbox" />
                {{ $t('connection.insecureTls') }}
              </label>
            </div>
            <div class="form-group checkbox-group">
              <label>
                <input v-model="newConnection.insecureHostname" type="checkbox" />
                {{ $t('connection.insecureHostname') }}
              </label>
            </div>
          </div>
          <div class="form-group">
            <label>
              {{ $t('connection.authentication') }}
              <button type="button" class="auth-toggle" @click="newConnection.enableAuth = !newConnection.enableAuth">
                {{ newConnection.enableAuth ? $t('connection.disable') : $t('connection.enable') }}
              </button>
            </label>
            <div v-if="newConnection.enableAuth" class="auth-fields">
              <div class="form-group">
                <label>{{ $t('connection.username') }}</label>
                <input v-model="newConnection.username" type="text" placeholder="admin" />
              </div>
              <div class="form-group">
                <label>{{ $t('connection.password') }}</label>
                <input v-model="newConnection.password" type="password" placeholder="••••••••" />
              </div>
            </div>
          </div>
          <div class="form-group">
            <label>
              {{ $t('connection.sshTunnel') }}
              <button type="button" class="auth-toggle" @click="newConnection.enableSSH = !newConnection.enableSSH">
                {{ newConnection.enableSSH ? $t('connection.disable') : $t('connection.enable') }}
              </button>
            </label>
            <div v-if="newConnection.enableSSH" class="ssh-fields">
              <div class="form-group">
                <label>{{ $t('connection.sshHost') }}</label>
                <input v-model="newConnection.sshHost" type="text" placeholder="ssh.example.com" />
              </div>
              <div class="form-group">
                <label>{{ $t('connection.sshPort') }}</label>
                <input v-model.number="newConnection.sshPort" type="number" placeholder="22" />
              </div>
              <div class="form-group">
                <label>{{ $t('connection.sshUsername') }}</label>
                <input v-model="newConnection.sshUsername" type="text" placeholder="user" />
              </div>
              <div class="form-group">
                <label>{{ $t('connection.sshPassword') }}</label>
                <input v-model="newConnection.sshPassword" type="password" placeholder="••••••••" />
              </div>
              <div class="form-group">
                <label>{{ $t('connection.sshKeyPath') }}</label>
                <div class="file-input-group">
                  <input v-model="newConnection.sshKeyPath" type="text" placeholder="/path/to/id_rsa" />
                  <button type="button" class="btn-browse" @click="selectSSHKeyFile('new')">
                    {{ $t('connection.browse') }}
                  </button>
                </div>
              </div>
              <div class="form-group">
                <label>{{ $t('connection.sshKeyPassphrase') }}</label>
                <input v-model="newConnection.sshKeyPassphrase" type="password" placeholder="••••••••" />
              </div>
            </div>
          </div>
        </div>
        <div class="modal-footer">
          <button @click="showAddConnection = false" class="btn-cancel">{{ $t('common.cancel') }}</button>
          <button @click="addConnection" class="btn-save">{{ $t('common.save') }}</button>
        </div>
      </div>
    </div>

    <!-- Edit Connection Modal -->
    <div v-if="showEditConnection" class="modal-overlay" @click="showEditConnection = false">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>{{ $t('connection.editConnection') }}</h3>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label>{{ $t('connection.connectionName') }}</label>
            <input v-model="editingConnection.name" type="text" placeholder="My Database" />
          </div>
          <div class="form-group">
            <label>
              {{ $t('connection.address') }}
              <span class="help-icon" :title="$t('connection.addressHelp')">
                <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg" width="14" height="14">
                  <circle cx="12" cy="12" r="9" stroke="currentColor" stroke-width="2"/>
                  <path d="M12 16v-4M12 8h.01" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
                </svg>
              </span>
            </label>
            <div class="address-input-group">
              <select v-model="editingConnection.protocol" class="protocol-select">
                <option value="http">HTTP</option>
                <option value="https">HTTPS</option>
              </select>
              <input v-model="editingConnection.address" type="text" placeholder="localhost:8086" class="address-input" />
            </div>
          </div>
          <div v-if="editingConnection.protocol === 'https'" class="https-options">
            <div class="form-group">
              <label>{{ $t('connection.caCertificate') }}</label>
              <div class="file-input-group">
                <input v-model="editingConnection.caCert" type="text" placeholder="/path/to/ca.crt" />
                <button type="button" class="btn-browse" @click="selectCaCertFile('edit')">
                  {{ $t('connection.browse') }}
                </button>
              </div>
            </div>
            <div class="form-group">
              <label>{{ $t('connection.clientCertificate') }}</label>
              <div class="file-input-group">
                <input v-model="editingConnection.clientCert" type="text" placeholder="/path/to/client.crt" />
                <button type="button" class="btn-browse" @click="selectClientCertFile('edit')">
                  {{ $t('connection.browse') }}
                </button>
              </div>
            </div>
            <div class="form-group">
              <label>{{ $t('connection.clientKey') }}</label>
              <div class="file-input-group">
                <input v-model="editingConnection.clientKey" type="text" placeholder="/path/to/client.key" />
                <button type="button" class="btn-browse" @click="selectClientKeyFile('edit')">
                  {{ $t('connection.browse') }}
                </button>
              </div>
            </div>
            <div class="form-group checkbox-group">
              <label>
                <input v-model="editingConnection.insecureTls" type="checkbox" />
                {{ $t('connection.insecureTls') }}
              </label>
            </div>
            <div class="form-group checkbox-group">
              <label>
                <input v-model="editingConnection.insecureHostname" type="checkbox" />
                {{ $t('connection.insecureHostname') }}
              </label>
            </div>
          </div>
          <div class="form-group">
            <label>
              {{ $t('connection.authentication') }}
              <button type="button" class="auth-toggle" @click="editingConnection.enableAuth = !editingConnection.enableAuth">
                {{ editingConnection.enableAuth ? $t('connection.disable') : $t('connection.enable') }}
              </button>
            </label>
            <div v-if="editingConnection.enableAuth" class="auth-fields">
              <div class="form-group">
                <label>{{ $t('connection.username') }}</label>
                <input v-model="editingConnection.username" type="text" placeholder="admin" />
              </div>
              <div class="form-group">
                <label>{{ $t('connection.password') }}</label>
                <input v-model="editingConnection.password" type="password" placeholder="••••••••" />
              </div>
            </div>
          </div>
          <div class="form-group">
            <label>
              {{ $t('connection.sshTunnel') }}
              <button type="button" class="auth-toggle" @click="editingConnection.enableSSH = !editingConnection.enableSSH">
                {{ editingConnection.enableSSH ? $t('connection.disable') : $t('connection.enable') }}
              </button>
            </label>
            <div v-if="editingConnection.enableSSH" class="ssh-fields">
              <div class="form-group">
                <label>{{ $t('connection.sshHost') }}</label>
                <input v-model="editingConnection.sshHost" type="text" placeholder="ssh.example.com" />
              </div>
              <div class="form-group">
                <label>{{ $t('connection.sshPort') }}</label>
                <input v-model.number="editingConnection.sshPort" type="number" placeholder="22" />
              </div>
              <div class="form-group">
                <label>{{ $t('connection.sshUsername') }}</label>
                <input v-model="editingConnection.sshUsername" type="text" placeholder="user" />
              </div>
              <div class="form-group">
                <label>{{ $t('connection.sshPassword') }}</label>
                <input v-model="editingConnection.sshPassword" type="password" placeholder="••••••••" />
              </div>
              <div class="form-group">
                <label>{{ $t('connection.sshKeyPath') }}</label>
                <div class="file-input-group">
                  <input v-model="editingConnection.sshKeyPath" type="text" placeholder="/path/to/id_rsa" />
                  <button type="button" class="btn-browse" @click="selectSSHKeyFile('edit')">
                    {{ $t('connection.browse') }}
                  </button>
                </div>
              </div>
              <div class="form-group">
                <label>{{ $t('connection.sshKeyPassphrase') }}</label>
                <input v-model="editingConnection.sshKeyPassphrase" type="password" placeholder="••••••••" />
              </div>
            </div>
          </div>
        </div>
        <div class="modal-footer">
          <button @click="showEditConnection = false" class="btn-cancel">{{ $t('common.cancel') }}</button>
          <button @click="saveEditConnection" class="btn-save">{{ $t('common.save') }}</button>
        </div>
      </div>
    </div>

    <!-- Error Dialog -->
    <div v-if="errorDialog.visible" class="modal-overlay" @click="closeErrorDialog">
      <div class="modal-content error-dialog" @click.stop>
        <div class="error-header">
          <h3>{{ $t('common.error') }}</h3>
          <button @click="closeErrorDialog" class="btn-close" :title="$t('common.close')">
            <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <line x1="18" y1="6" x2="6" y2="18" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
              <line x1="6" y1="6" x2="18" y2="18" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
            </svg>
          </button>
        </div>
        <div class="error-content">
          <svg class="error-icon" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2"/>
            <line x1="12" y1="8" x2="12" y2="12" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
            <line x1="12" y1="16" x2="12.01" y2="16" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
          </svg>
          <p class="error-message">{{ errorDialog.message }}</p>
        </div>
        <div class="modal-footer">
          <button @click="closeErrorDialog" class="btn-confirm">{{ $t('common.ok') }}</button>
        </div>
      </div>
    </div>

    <!-- Confirm Dialog -->
    <div v-if="confirmDialog.visible" class="modal-overlay" @click="closeConfirmDialog">
      <div class="modal-content confirm-dialog" @click.stop>
        <div class="confirm-header">
          <h3>{{ $t('common.confirm') }}</h3>
          <button @click="closeConfirmDialog" class="btn-close" :title="$t('common.close')">
            <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <line x1="18" y1="6" x2="6" y2="18" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
              <line x1="6" y1="6" x2="18" y2="18" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
            </svg>
          </button>
        </div>
        <div class="confirm-content">
          <svg class="confirm-icon" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2"/>
            <line x1="12" y1="8" x2="12" y2="12" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
            <line x1="12" y1="16" x2="12.01" y2="16" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
          </svg>
          <p class="confirm-message">{{ confirmDialog.message }}</p>
        </div>
        <div class="modal-footer">
          <button @click="closeConfirmDialog" class="btn-cancel">{{ $t('common.cancel') }}</button>
          <button @click="handleConfirm" class="btn-confirm">{{ $t('common.confirm') }}</button>
        </div>
      </div>
    </div>

    <!-- Resize Handle -->
    <div
      class="resize-handle"
      @mousedown="startResize"
    ></div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import type { SavedConnection, Database, ConnectionConfig } from '../types'
import { AddConnect, DeleteConnect, ListConnects, UpdateConnect, DialConnect, OpenFileDialog, GetDatabaseMetadata, CloseConnect } from '../../wailsjs/go/main/App'
import { main } from '../../wailsjs/go/models'

// Data transformation utilities
const toBackendConfig = (config: ConnectionConfig): main.ConnectConfig => {
  return {
    name: config.name,
    address: config.address || '',
    http_schema: config.protocol || 'http',
    ca_certificate: config.caCert || '',
    client_certificate: config.clientCert || '',
    client_key: config.clientKey || '',
    insecure_tls: config.insecureTls || false,
    insecure_hostname: config.insecureHostname || false,
    enable_auth: config.enableAuth || false,
    username: config.username || '',
    password: config.password || '',
    enable_ssh: config.enableSSH || false,
    ssh_host: config.sshHost || '',
    ssh_port: config.sshPort || 22,
    ssh_username: config.sshUsername || '',
    ssh_password: config.sshPassword || '',
    ssh_key_path: config.sshKeyPath || '',
    ssh_key_passphrase: config.sshKeyPassphrase || ''
  }
}

const fromBackendConfig = (config: main.ConnectConfig): SavedConnection => {
  return {
    id: config.name, // Using name as ID
    name: config.name,
    address: config.address,
    protocol: (config.http_schema === 'https' ? 'https' : 'http') as 'http' | 'https',
    caCert: config.ca_certificate,
    clientCert: config.client_certificate,
    clientKey: config.client_key,
    insecureTls: config.insecure_tls,
    insecureHostname: config.insecure_hostname,
    enableAuth: config.enable_auth,
    username: config.username,
    password: config.password,
    enableSSH: config.enable_ssh,
    sshHost: config.ssh_host,
    sshPort: config.ssh_port,
    sshUsername: config.ssh_username,
    sshPassword: config.ssh_password,
    sshKeyPath: config.ssh_key_path,
    sshKeyPassphrase: config.ssh_key_passphrase,
    database: '',
    databases: [],
    expanded: false,
    connected: false
  }
}

const props = defineProps<{
  connections: SavedConnection[]
  collapsed: boolean
}>()

const emit = defineEmits<{
  'update:connections': [connections: SavedConnection[]]
  'select-measurement': [data: { measurement: string; database: string; connection: SavedConnection }]
  'connect': [connection: SavedConnection]
  'disconnect': [connectionId: string]
  'update-retention-policies': [data: { policies: any[]; database: string }]
}>()

const sidebarWidth = ref(300)
const showAddConnection = ref(false)
const showEditConnection = ref(false)
const editingConnection = ref<SavedConnection>({
  id: '',
  name: '',
  protocol: 'http',
  address: '',
  enableAuth: false,
  username: '',
  password: '',
  database: '',
  caCert: '',
  clientCert: '',
  clientKey: '',
  insecureTls: false,
  insecureHostname: false,
  enableSSH: false,
  sshHost: '',
  sshPort: 22,
  sshUsername: '',
  sshPassword: '',
  sshKeyPath: '',
  sshKeyPassphrase: '',
  expanded: false,
  connected: false,
  databases: []
})

const contextMenu = ref<{
  visible: boolean
  x: number
  y: number
  connection: SavedConnection | null
}>({
  visible: false,
  x: 0,
  y: 0,
  connection: null
})

const newConnection = ref<ConnectionConfig>({
  id: '',
  name: 'localhost',
  protocol: 'http',
  address: 'localhost:8086',
  enableAuth: false,
  username: '',
  password: '',
  database: '',
  caCert: '',
  clientCert: '',
  clientKey: '',
  insecureTls: false,
  insecureHostname: false,
  enableSSH: false,
  sshHost: '',
  sshPort: 22,
  sshUsername: '',
  sshPassword: '',
  sshKeyPath: '',
  sshKeyPassphrase: ''
})

// Error dialog state
const errorDialog = ref({
  visible: false,
  message: ''
})

// Confirm dialog state
const confirmDialog = ref({
  visible: false,
  message: '',
  onConfirm: () => {}
})

// Error handling utility
const showError = (error: unknown, context: string) => {
  const errorMessage = error instanceof Error ? error.message : String(error)
  console.error(`${context}:`, error)
  errorDialog.value = {
    visible: true,
    message: `${context}: ${errorMessage}`
  }
}

const closeErrorDialog = () => {
  errorDialog.value.visible = false
}

// Confirm dialog utilities
const showConfirm = (message: string, onConfirm: () => void) => {
  confirmDialog.value = {
    visible: true,
    message,
    onConfirm
  }
}

const closeConfirmDialog = () => {
  confirmDialog.value.visible = false
}

const handleConfirm = () => {
  confirmDialog.value.onConfirm()
  closeConfirmDialog()
}

const toggleConnection = (conn: SavedConnection) => {
  if (conn.connected) {
    conn.expanded = !conn.expanded
  }
}

const loadDatabaseMetadata = async (db: Database, conn: SavedConnection) => {
  try {
    // Toggle expanded state
    db.expanded = !db.expanded

    // If expanding and measurements not loaded yet, fetch metadata
    if (db.expanded && db.measurements.length === 0) {
      const metadata = await GetDatabaseMetadata(conn.id, db.name)

      // Convert measurements from string[] to Measurement[] format
      db.measurements = metadata.measurements.map((name: string) => ({
        name,
        fields: []
      }))

      // Convert retention policies from backend format to frontend format
      db.retentionPolicies = metadata.retention_policies.map(policy => ({
        name: policy.name,
        duration: policy.duration,
        replication: 1,  // Default value
        isDefault: policy.name === 'autogen'  // Default policy is usually 'autogen'
      }))

      // Emit retention policies to parent component
      emit('update-retention-policies', {
        policies: db.retentionPolicies,
        database: db.name
      })
    }
  } catch (error) {
    showError(error, `Failed to load metadata for database: ${db.name}`)
  }
}

const selectMeasurement = (measurement: string, database: string, connection: SavedConnection) => {
  emit('select-measurement', { measurement, database, connection })
}

// Load connections from backend on mount
onMounted(async () => {
  try {
    const backendConnections = await ListConnects()
    const loadedConnections = (backendConnections || []).map(fromBackendConfig)
    emit('update:connections', loadedConnections)
  } catch (error) {
    showError(error, 'Failed to load connections')
  }
})

const addConnection = async () => {
  try {
    // Create backend config from form data
    const backendConfig = toBackendConfig(newConnection.value)

    // Save to backend
    await AddConnect(backendConfig)

    // Create frontend connection object
    const connection: SavedConnection = fromBackendConfig(backendConfig)

    // Update local state
    emit('update:connections', [...props.connections, connection])

    // Reset form
    newConnection.value = {
      id: '',
      name: 'localhost',
      protocol: 'http',
      address: 'localhost:8086',
      enableAuth: false,
      username: '',
      password: '',
      database: '',
      caCert: '',
      clientCert: '',
      clientKey: '',
      insecureTls: false,
      insecureHostname: false,
      enableSSH: false,
      sshHost: '',
      sshPort: 22,
      sshUsername: '',
      sshPassword: '',
      sshKeyPath: '',
      sshKeyPassphrase: ''
    }
    showAddConnection.value = false
  } catch (error) {
    showError(error, 'Failed to add connection')
  }
}

const deleteConnection = async (id: string) => {
  showConfirm('Are you sure you want to delete this connection?', async () => {
    try {
      // Delete from backend using connection name (id is the name)
      await DeleteConnect(id)

      // Update local state
      emit('update:connections', props.connections.filter(c => c.id !== id))
    } catch (error) {
      showError(error, 'Failed to delete connection')
    }
  })
}

const connectToDatabase = async (conn: SavedConnection) => {
  try {
    // Disconnect all other connections before connecting to this one
    // This ensures only one connection is active at a time
    for (const existingConn of props.connections) {
      if (existingConn.id !== conn.id && existingConn.connected) {
        try {
          // Call backend to close the connection
          await CloseConnect(existingConn.id)

          // Update frontend state
          existingConn.connected = false
          existingConn.expanded = false
          existingConn.databases = []
        } catch (error) {
          console.error(`Failed to close connection ${existingConn.id}:`, error)
        }
      }
    }

    // Call DialConnect to establish connection and get database list
    const databaseNames = await DialConnect(conn.id)

    // Convert database names to Database objects
    const databases: Database[] = databaseNames.map(name => ({
      name,
      expanded: false,
      measurements: [],
      retentionPolicies: []
    }))

    // Update connection status and databases
    conn.connected = true
    conn.expanded = true
    conn.databases = databases

    // Emit connect event to parent component
    emit('connect', conn)
  } catch (error) {
    // Show error dialog if connection failed
    const errorMessage = error instanceof Error ? error.message : String(error)
    showError(error, `Failed to connect to database: ${errorMessage}`)
  }
}

const disconnectFromDatabase = (conn: SavedConnection) => {
  emit('disconnect', conn.id)
}

const showContextMenu = (event: MouseEvent, conn: SavedConnection) => {
  contextMenu.value = {
    visible: true,
    x: event.clientX,
    y: event.clientY,
    connection: conn
  }
  
  // Close context menu when clicking outside
  setTimeout(() => {
    document.addEventListener('click', hideContextMenu, { once: true })
  }, 0)
}

const hideContextMenu = () => {
  contextMenu.value.visible = false
}

const editConnection = (conn: SavedConnection) => {
  // Convert legacy format to new format if needed
  const converted = { ...conn }

  // Ensure protocol has a default value
  if (!converted.protocol) {
    converted.protocol = 'http'
  }

  // Generate address from host/port if address is not set
  if (!converted.address && converted.host && converted.port) {
    converted.address = `${converted.host}:${converted.port}`
  } else if (!converted.address) {
    converted.address = ''
  }

  // Ensure enableAuth has a default value
  if (converted.enableAuth === undefined) {
    converted.enableAuth = !!(converted.username && converted.password)
  }

  // Ensure optional fields have empty string defaults
  if (!converted.caCert) converted.caCert = ''
  if (!converted.clientCert) converted.clientCert = ''
  if (!converted.clientKey) converted.clientKey = ''
  if (converted.insecureTls === undefined) converted.insecureTls = false
  if (converted.insecureHostname === undefined) converted.insecureHostname = false
  if (converted.enableSSH === undefined) converted.enableSSH = false
  if (!converted.sshHost) converted.sshHost = ''
  if (!converted.sshPort) converted.sshPort = 22
  if (!converted.sshUsername) converted.sshUsername = ''
  if (!converted.sshPassword) converted.sshPassword = ''
  if (!converted.sshKeyPath) converted.sshKeyPath = ''
  if (!converted.sshKeyPassphrase) converted.sshKeyPassphrase = ''

  editingConnection.value = converted
  showEditConnection.value = true
}

const saveEditConnection = async () => {
  try {
    const oldName = editingConnection.value.id

    // Convert to backend format
    const backendConfig = toBackendConfig(editingConnection.value)

    // Update connection using UpdateConnect API
    await UpdateConnect(oldName, backendConfig)

    // Update local state
    const index = props.connections.findIndex(c => c.id === oldName)
    if (index !== -1) {
      const updated = [...props.connections]
      updated[index] = fromBackendConfig(backendConfig)
      emit('update:connections', updated)
    }

    showEditConnection.value = false
  } catch (error) {
    showError(error, 'Failed to update connection')
  }
}

const refreshConnection = async (conn: SavedConnection) => {
  if (conn.connected) {
    try {
      // Call DialConnect to re-establish connection and get updated database list
      const databaseNames = await DialConnect(conn.id)

      // Convert database names to Database objects
      const databases: Database[] = databaseNames.map(name => ({
        name,
        expanded: false,
        measurements: [],
        retentionPolicies: []
      }))

      // Update connection's database list without triggering disconnect/connect events
      conn.databases = databases
      conn.connected = true
      conn.expanded = true
    } catch (error) {
      // Show error dialog if refresh failed
      const errorMessage = error instanceof Error ? error.message : String(error)
      showError(error, `Failed to refresh connection: ${errorMessage}`)
    }
  }
}

// File selection functions
const selectCaCertFile = async (type: 'new' | 'edit') => {
  try {
    const filePath = await OpenFileDialog()
    if (filePath) {
      if (type === 'new') {
        newConnection.value.caCert = filePath
      } else {
        editingConnection.value.caCert = filePath
      }
    }
  } catch (error) {
    console.error('Failed to select CA certificate file:', error)
  }
}

const selectClientCertFile = async (type: 'new' | 'edit') => {
  try {
    const filePath = await OpenFileDialog()
    if (filePath) {
      if (type === 'new') {
        newConnection.value.clientCert = filePath
      } else {
        editingConnection.value.clientCert = filePath
      }
    }
  } catch (error) {
    console.error('Failed to select client certificate file:', error)
  }
}

const selectClientKeyFile = async (type: 'new' | 'edit') => {
  try {
    const filePath = await OpenFileDialog()
    if (filePath) {
      if (type === 'new') {
        newConnection.value.clientKey = filePath
      } else {
        editingConnection.value.clientKey = filePath
      }
    }
  } catch (error) {
    console.error('Failed to select client key file:', error)
  }
}

const selectSSHKeyFile = async (type: 'new' | 'edit') => {
  try {
    const filePath = await OpenFileDialog()
    if (filePath) {
      if (type === 'new') {
        newConnection.value.sshKeyPath = filePath
      } else {
        editingConnection.value.sshKeyPath = filePath
      }
    }
  } catch (error) {
    console.error('Failed to select SSH key file:', error)
  }
}

// Resize functionality
const startResize = (e: MouseEvent) => {
  e.preventDefault()
  const startX = e.clientX
  const startWidth = sidebarWidth.value

  const onMouseMove = (e: MouseEvent) => {
    const newWidth = startWidth + (e.clientX - startX)
    sidebarWidth.value = Math.max(200, Math.min(600, newWidth))
  }

  const onMouseUp = () => {
    document.removeEventListener('mousemove', onMouseMove)
    document.removeEventListener('mouseup', onMouseUp)
  }

  document.addEventListener('mousemove', onMouseMove)
  document.addEventListener('mouseup', onMouseUp)
}
</script>

<style scoped>
.connection-manager {
  background: var(--bg-secondary);
  border-right: 1px solid var(--border-color);
  display: flex;
  flex-direction: column;
  position: relative;
  min-width: 200px;
  max-width: 600px;
}

.manager-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  border-bottom: 1px solid var(--border-color);
}

.manager-header h3 {
  margin: 0;
  font-size: 16px;
  color: var(--text-primary);
}

.btn-add {
  width: 28px;
  height: 28px;
  border: none;
  background: var(--accent-color);
  color: white;
  border-radius: 4px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background 0.2s;
}

/* Add styles for SVG icon in add button */
.btn-add svg {
  width: 18px;
  height: 18px;
}

.connections-list {
  flex: 1;
  overflow-y: auto;
  padding: 8px;
}

.connection-item {
  margin-bottom: 4px;
  /* Add border to separate connection items */
  border: 1px solid var(--border-color);
  border-radius: 6px;
  overflow: hidden;
}

.connection-header {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px;
  cursor: pointer;
  transition: background 0.2s;
  user-select: none;
  -webkit-user-select: none;
  -moz-user-select: none;
  -ms-user-select: none;
}

.connection-header:hover {
  background: var(--bg-hover);
}

.connection-header .name {
  flex: 1;
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
}

.connection-actions {
  display: flex;
  gap: 4px;
}

.btn-action {
  padding: 4px;
  border: none;
  background: transparent;
  cursor: pointer;
  border-radius: 3px;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
  /* Add width/height for SVG icons */
  width: 24px;
  height: 24px;
  color: var(--text-secondary);
}

.btn-action:hover {
  background: var(--bg-hover);
  color: var(--text-primary);
}

/* Style SVG icons in buttons */
.btn-action svg {
  width: 16px;
  height: 16px;
}

/* Add styles for SVG status icons */
.status-icon {
  width: 16px;
  height: 16px;
  flex-shrink: 0;
}

.status-icon.connected {
  color: #22c55e;
}

.status-icon.disconnected {
  color: #94a3b8;
}

.database-item {
  margin-bottom: 2px;
}

.database-header {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 8px;
  cursor: pointer;
  border-radius: 4px;
  transition: background 0.2s;
  user-select: none;
  -webkit-user-select: none;
  -moz-user-select: none;
  -ms-user-select: none;
}

.database-header:hover {
  background: var(--bg-hover);
}

.database-header .name {
  font-size: 13px;
  color: var(--text-primary);
}

.measurements-list {
  margin-left: 12px;
}

.measurement-item {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 8px;
  cursor: pointer;
  border-radius: 4px;
  transition: background 0.2s;
  user-select: none;
  -webkit-user-select: none;
  -moz-user-select: none;
  -ms-user-select: none;
}

.measurement-item:hover {
  background: var(--bg-hover);
}

.measurement-item .name {
  font-size: 13px;
  color: var(--text-secondary);
}

.icon {
  font-size: 14px;
  flex-shrink: 0;
}

.icon.indent {
  margin-left: 12px;
}

.icon.double-indent {
  margin-left: 24px;
}

.icon.icon-placeholder {
  width: 16px;
  height: 16px;
  display: inline-block;
}

/* Add styles for chevron icons */
.icon-chevron {
  width: 16px;
  height: 16px;
  flex-shrink: 0;
  color: var(--text-secondary);
}

.icon-chevron.indent {
  margin-left: 12px;
}

/* Add styles for database and table icons */
.icon-db {
  width: 16px;
  height: 16px;
  flex-shrink: 0;
  color: var(--text-secondary);
}

.icon-table {
  width: 14px;
  height: 14px;
  flex-shrink: 0;
  color: var(--text-secondary);
}

.icon-table.double-indent {
  margin-left: 24px;
}

/* Context Menu Styles */
.context-menu {
  position: fixed;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 6px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
  z-index: 1000;
  min-width: 160px;
  padding: 4px 0;
}

.context-menu-item {
  padding: 8px 16px;
  cursor: pointer;
  color: var(--text-primary);
  font-size: 13px;
  display: flex;
  align-items: center;
  gap: 8px;
  transition: background 0.2s;
}

.context-menu-item svg {
  width: 16px;
  height: 16px;
  flex-shrink: 0;
  color: var(--text-secondary);
}

.context-menu-item:hover {
  background: var(--bg-hover);
}

.context-menu-item:hover svg {
  color: var(--text-primary);
}

/* Resize Handle */
.resize-handle {
  position: absolute;
  right: 0;
  top: 0;
  bottom: 0;
  width: 4px;
  cursor: col-resize;
  background: transparent;
  transition: background 0.2s;
}

.resize-handle:hover {
  background: var(--accent-color);
}

/* Modal Overlay and Content Styles */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2000;
}

.modal-content {
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  min-width: 400px;
  max-width: 500px;
  width: 90%;
  max-height: 90vh;
  display: flex;
  flex-direction: column;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.4);
  overflow: hidden;
}

.modal-header {
  padding: 20px 24px;
  border-bottom: 1px solid var(--border-color);
  flex-shrink: 0;
}

.modal-header h3 {
  margin: 0;
  font-size: 18px;
  color: var(--text-primary);
}

.modal-body {
  padding: 20px 24px;
  overflow-y: auto;
  flex: 1;
  min-height: 0;
}

.modal-footer {
  padding: 16px 24px;
  border-top: 1px solid var(--border-color);
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  flex-shrink: 0;
  background: var(--bg-secondary);
}

.form-group {
  margin-bottom: 16px;
}

.form-group label {
  display: block;
  margin-bottom: 6px;
  font-size: 13px;
  font-weight: 500;
  color: var(--text-primary);
}

.form-group input {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  background: var(--bg-primary);
  color: var(--text-primary);
  font-size: 14px;
  box-sizing: border-box;
}

.form-group input:focus {
  outline: none;
  border-color: var(--accent-color);
}

.file-input-group {
  display: flex;
  gap: 8px;
  align-items: center;
}

.file-input-group input {
  flex: 1;
}

.btn-browse {
  padding: 8px 16px;
  background: var(--bg-tertiary);
  color: var(--text-primary);
  border: 1px solid var(--border-color);
  border-radius: 4px;
  font-size: 14px;
  cursor: pointer;
  white-space: nowrap;
  transition: all 0.2s;
}

.btn-browse:hover {
  background: var(--bg-hover);
  border-color: var(--accent-color);
}

.btn-cancel,
.btn-save {
  padding: 8px 16px;
  border: none;
  border-radius: 4px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-cancel {
  background: var(--bg-hover);
  color: var(--text-primary);
}

.btn-cancel:hover {
  background: var(--border-color);
}

.btn-save {
  background: var(--accent-color);
  color: white;
}

.btn-save:hover {
  opacity: 0.9;
}

/* New connection form styles */
.help-icon {
  display: inline-flex;
  align-items: center;
  margin-left: 4px;
  color: var(--text-secondary);
  cursor: help;
}

.help-icon svg {
  vertical-align: middle;
}

.address-input-group {
  display: flex;
  gap: 8px;
}

.protocol-select {
  padding: 8px 12px;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  background-color: var(--bg-primary);
  color: var(--text-primary) !important;
  font-size: 14px;
  min-width: 80px;
  -webkit-appearance: none;
  -moz-appearance: none;
  appearance: none;
  background-image: url("data:image/svg+xml;charset=UTF-8,%3csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='%23e4e6eb' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3e%3cpolyline points='6 9 12 15 18 9'%3e%3c/polyline%3e%3c/svg%3e");
  background-repeat: no-repeat;
  background-position: right 6px center;
  background-size: 12px;
  padding-right: 24px;
}

:root[data-theme="light"] .protocol-select {
  background-image: url("data:image/svg+xml;charset=UTF-8,%3csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='%231f2937' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3e%3cpolyline points='6 9 12 15 18 9'%3e%3c/polyline%3e%3c/svg%3e");
}

.protocol-select option {
  background: var(--bg-primary);
  color: var(--text-primary);
}

.address-input {
  flex: 1;
}

.https-options {
  margin-top: 12px;
  padding: 12px;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  background: var(--bg-hover);
}

.auth-toggle {
  margin-left: 8px;
  padding: 4px 8px;
  font-size: 12px;
  background: var(--accent-color);
  color: white;
  border: none;
  border-radius: 3px;
  cursor: pointer;
  transition: opacity 0.2s;
}

.auth-toggle:hover {
  opacity: 0.9;
}

.auth-fields {
  margin-top: 12px;
  padding: 12px;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  background: var(--bg-hover);
}

.ssh-fields {
  margin-top: 12px;
  padding: 12px;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  background: var(--bg-hover);
}

.checkbox-group label {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}

.checkbox-group input[type="checkbox"] {
  width: auto;
  margin: 0;
}

/* Disable pointer events when not connected */
.connection-header.not-connected {
  cursor: default;
}

.connection-header.not-connected:hover {
  background: transparent;
}

/* Error Dialog Styles */
.error-dialog {
  min-width: 450px;
  max-width: 550px;
}

.error-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.error-header h3 {
  margin: 0;
  font-size: 18px;
  color: #ef4444;
  display: flex;
  align-items: center;
  gap: 8px;
}

.btn-close {
  width: 28px;
  height: 28px;
  border: none;
  background: transparent;
  color: var(--text-secondary);
  border-radius: 4px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
  padding: 0;
}

.btn-close:hover {
  background: var(--bg-hover);
  color: var(--text-primary);
}

.btn-close svg {
  width: 18px;
  height: 18px;
}

.error-content {
  display: flex;
  gap: 16px;
  margin-bottom: 24px;
  padding: 16px;
  background: rgba(239, 68, 68, 0.1);
  border: 1px solid rgba(239, 68, 68, 0.2);
  border-radius: 6px;
}

.error-icon {
  width: 24px;
  height: 24px;
  flex-shrink: 0;
  color: #ef4444;
}

.error-message {
  flex: 1;
  margin: 0;
  font-size: 14px;
  line-height: 1.6;
  color: var(--text-primary);
  word-break: break-word;
}

.btn-confirm {
  padding: 8px 24px;
  border: none;
  border-radius: 4px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
  background: var(--accent-color);
  color: white;
  min-width: 80px;
}

.btn-confirm:hover {
  opacity: 0.9;
}

/* Confirm Dialog Styles */
.confirm-dialog {
  min-width: 450px;
  max-width: 550px;
}

.confirm-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.confirm-header h3 {
  margin: 0;
  font-size: 18px;
  color: var(--accent-color);
  display: flex;
  align-items: center;
  gap: 8px;
}

.confirm-content {
  display: flex;
  gap: 16px;
  margin-bottom: 24px;
  padding: 16px;
  background: rgba(59, 130, 246, 0.1);
  border: 1px solid rgba(59, 130, 246, 0.2);
  border-radius: 6px;
}

.confirm-icon {
  width: 24px;
  height: 24px;
  flex-shrink: 0;
  color: var(--accent-color);
}

.confirm-message {
  flex: 1;
  margin: 0;
  font-size: 14px;
  line-height: 1.6;
  color: var(--text-primary);
  word-break: break-word;
}

.btn-cancel {
  padding: 8px 24px;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
  background: transparent;
  color: var(--text-primary);
  min-width: 80px;
}

.btn-cancel:hover {
  background: var(--hover-bg);
}

</style>
