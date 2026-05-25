import request from '../utils/request'

export interface LoginParams {
  username: string
  password: string
}

export interface LoginResult {
  token: string
  user_info: {
    id: number
    username: string
    role: string
  }
}

export function login(params: LoginParams): Promise<LoginResult> {
  return request.post('/auth/login', params)
}

export function register(params: { username: string; password: string }) {
  return request.post('/auth/register', params)
}

export function getProfile() {
  return request.get('/auth/profile')
}
