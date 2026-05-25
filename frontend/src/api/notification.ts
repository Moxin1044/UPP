import request from '../utils/request'

export interface Notification {
  id: number
  user_id: number
  type: string
  config: string
  enabled: boolean
}

export function getNotifications(): Promise<Notification[]> {
  return request.get('/notifications')
}

export function createNotification(data: { type: string; config: Record<string, any>; enabled?: boolean }): Promise<Notification> {
  return request.post('/notifications', data)
}

export function updateNotification(id: number, data: { type: string; config: Record<string, any>; enabled?: boolean }): Promise<Notification> {
  return request.put(`/notifications/${id}`, data)
}

export function deleteNotification(id: number) {
  return request.delete(`/notifications/${id}`)
}

export function testNotification(id: number) {
  return request.post(`/notifications/${id}/test`)
}
