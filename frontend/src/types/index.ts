export interface ConnectionConfig {
  id: string
  name: string
  // Legacy fields (optional for backward compatibility)
  host?: string
  port?: number
  // New fields
  protocol?: 'http' | 'https'
  address?: string
  enableAuth?: boolean
  username: string
  password: string
  database: string
  caCert?: string
  clientCert?: string
  clientKey?: string
  insecureTls?: boolean
  insecureHostname?: boolean
}

export interface Measurement {
  name: string
  fields: string[]
}

export interface RetentionPolicy {
  name: string
  duration: string
  replication: number
  isDefault: boolean
}

export interface Database {
  name: string
  expanded: boolean
  measurements: Measurement[]
  retentionPolicies: RetentionPolicy[]
}

export interface QueryResult {
  [key: string]: string | number
}

export interface TreeItem {
  type: "database" | "measurement" | "field"
  name: string
  parent?: string
}

export interface SavedConnection extends ConnectionConfig {
  databases: Database[]
  expanded: boolean
  connected: boolean
}

export interface QueryHistoryItem {
  id: string
  query: string
  timestamp: Date
  executionTime: number
  database?: string
  retentionPolicy?: string
  success: boolean
  error?: string
}

export type Theme = "light" | "dark"

export type ThemeMode = "light" | "dark" | "system"

export interface AppSettings {
  language: string
  themeMode: ThemeMode
  customFont: string
  maxHistoryCount: number
  dataDirectory: string
}
