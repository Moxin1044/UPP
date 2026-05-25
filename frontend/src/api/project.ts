import request from '../utils/request'

export interface Project {
  id: number
  user_id: number
  name: string
  description: string
  created_at: string
}

export function getProjects(): Promise<Project[]> {
  return request.get('/projects')
}

export function getProject(id: number): Promise<Project> {
  return request.get(`/projects/${id}`)
}

export function createProject(data: { name: string; description: string }): Promise<Project> {
  return request.post('/projects', data)
}

export function updateProject(id: number, data: { name: string; description: string }): Promise<Project> {
  return request.put(`/projects/${id}`, data)
}

export function deleteProject(id: number) {
  return request.delete(`/projects/${id}`)
}
