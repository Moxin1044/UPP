import { defineStore } from 'pinia'
import { ref } from 'vue'
import { login as loginApi, getProfile } from '../api/auth'

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const userInfo = ref<{ id: number; username: string; role: string } | null>(null)

  async function login(username: string, password: string, remember = false) {
    const res = await loginApi({ username, password })
    token.value = res.token
    userInfo.value = res.user_info
    if (remember) {
      localStorage.setItem('token', res.token)
    } else {
      sessionStorage.setItem('token', res.token)
    }
    localStorage.setItem('token', res.token)
  }

  async function fetchProfile() {
    try {
      const profile = await getProfile()
      userInfo.value = profile as any
    } catch {
      logout()
    }
  }

  function logout() {
    token.value = ''
    userInfo.value = null
    localStorage.removeItem('token')
    sessionStorage.removeItem('token')
  }

  const isAdmin = () => userInfo.value?.role === 'admin'

  return { token, userInfo, login, fetchProfile, logout, isAdmin }
})
