import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getDashboardStats, type DashboardStats } from '../api/dashboard'

export const useDashboardStore = defineStore('dashboard', () => {
  const stats = ref<DashboardStats | null>(null)
  const loading = ref(false)

  async function fetchStats() {
    loading.value = true
    try {
      stats.value = await getDashboardStats()
    } finally {
      loading.value = false
    }
  }

  return { stats, loading, fetchStats }
})
