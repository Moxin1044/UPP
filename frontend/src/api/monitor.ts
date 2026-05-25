import request from '../utils/request'

export interface Monitor {
  id: number
  user_id: number
  project_id: number
  name: string
  type: string
  target: string
  interval: number
  timeout: number
  enabled: boolean
  created_at: string
}

export interface MonitorResult {
  id: number
  monitor_id: number
  status: string
  latency: number
  message: string
  created_at: string
}

export interface MonitorStats {
  total_checks: number
  up_checks: number
  down_checks: number
  uptime: number
  avg_latency: number
}

export function getMonitors(): Promise<Monitor[]> {
  return request.get('/monitors')
}

export function getMonitorsByProject(projectId: number): Promise<Monitor[]> {
  return request.get(`/projects/${projectId}/monitors`)
}

export function getMonitor(id: number): Promise<Monitor> {
  return request.get(`/monitors/${id}`)
}

export function createMonitor(data: Partial<Monitor>): Promise<Monitor> {
  return request.post('/monitors', data)
}

export function updateMonitor(id: number, data: Partial<Monitor>): Promise<Monitor> {
  return request.put(`/monitors/${id}`, data)
}

export function toggleMonitor(id: number, enabled: boolean): Promise<Monitor> {
  return request.patch(`/monitors/${id}/toggle`, { enabled })
}

export function deleteMonitor(id: number) {
  return request.delete(`/monitors/${id}`)
}

export function getMonitorResults(id: number, limit = 100): Promise<MonitorResult[]> {
  return request.get(`/monitors/${id}/results`, { params: { limit } })
}

export function getMonitorStats(id: number): Promise<MonitorStats> {
  return request.get(`/monitors/${id}/stats`)
}
