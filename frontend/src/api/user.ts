import request from '../utils/request'

export interface User {
  id: number
  username: string
  role: string
  created_at: string
}

export function getUsers(): Promise<User[]> {
  return request.get('/users')
}

export function createUser(data: { username: string; password: string; role: string }): Promise<User> {
  return request.post('/users', data)
}

export function updateUserPassword(id: number, password: string) {
  return request.put(`/users/${id}/password`, { password })
}

export function deleteUser(id: number) {
  return request.delete(`/users/${id}`)
}
