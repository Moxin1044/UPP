import request from '../utils/request'

export interface DashboardStats {
  total_monitors: number
  up_monitors: number
  down_monitors: number
  avg_latency: number
}

export function getDashboardStats(): Promise<DashboardStats> {
  return request.get('/dashboard/stats')
}

export function getProjectDashboard(projectId: number) {
  return request.get(`/dashboard/project/${projectId}`)
}

export function getSettings(): Promise<Record<string, string>> {
  return request.get('/settings')
}

export function setSetting(key: string, value: string) {
  return request.post('/settings', { key, value })
}
