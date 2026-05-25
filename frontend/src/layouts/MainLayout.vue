<template>
  <t-layout class="main-layout">
    <t-aside class="sidebar" :collapsed="collapsed">
      <div class="logo">
        <svg v-if="!collapsed" viewBox="0 0 64 64" fill="none" xmlns="http://www.w3.org/2000/svg" class="logo-svg">
          <defs>
            <linearGradient id="sideLogoGrad" x1="0" y1="0" x2="64" y2="64">
              <stop stop-color="#00d4ff" />
              <stop offset="1" stop-color="#7b61ff" />
            </linearGradient>
          </defs>
          <rect x="2" y="2" width="60" height="60" rx="16" stroke="url(#sideLogoGrad)" stroke-width="2" />
          <path d="M32 12L48 20V44L32 52L16 44V20L32 12Z" stroke="url(#sideLogoGrad)" stroke-width="1.5" fill="none" />
          <circle cx="32" cy="32" r="6" fill="url(#sideLogoGrad)" opacity="0.8" />
        </svg>
        <span v-if="!collapsed" class="logo-text">UPP</span>
        <svg v-else viewBox="0 0 64 64" fill="none" xmlns="http://www.w3.org/2000/svg" class="logo-svg-small">
          <defs>
            <linearGradient id="sideLogoGradS" x1="0" y1="0" x2="64" y2="64">
              <stop stop-color="#00d4ff" />
              <stop offset="1" stop-color="#7b61ff" />
            </linearGradient>
          </defs>
          <circle cx="32" cy="32" r="20" fill="url(#sideLogoGradS)" opacity="0.8" />
        </svg>
      </div>
      <t-menu :value="activeMenu" :collapsed="collapsed" @change="onMenuChange">
        <t-menu-item value="Dashboard">
          <template #icon><t-icon name="dashboard" /></template>
          {{ $t('nav.dashboard') }}
        </t-menu-item>
        <t-menu-item value="Projects">
          <template #icon><t-icon name="folder" /></template>
          {{ $t('nav.projects') }}
        </t-menu-item>
        <t-menu-item value="Monitors">
          <template #icon><t-icon name="root-list" /></template>
          {{ $t('nav.monitors') }}
        </t-menu-item>
        <t-menu-item value="Notifications">
          <template #icon><t-icon name="notification" /></template>
          {{ $t('nav.notifications') }}
        </t-menu-item>
        <t-menu-item v-if="isAdmin" value="Users">
          <template #icon><t-icon name="user-circle" /></template>
          {{ $t('nav.users') }}
        </t-menu-item>
        <t-menu-item v-if="isAdmin" value="Settings">
          <template #icon><t-icon name="setting" /></template>
          {{ $t('nav.settings') }}
        </t-menu-item>
      </t-menu>
    </t-aside>
    <t-layout>
      <t-header class="header">
        <div class="header-left">
          <t-button variant="text" @click="collapsed = !collapsed" class="collapse-btn">
            <template #icon><t-icon name="view-list" /></template>
          </t-button>
        </div>
        <div class="header-right">
          <t-select v-model="currentLocale" size="small" class="lang-select" @change="onLocaleChange">
            <t-option value="zh" label="中文" />
            <t-option value="en" label="English" />
          </t-select>
          <t-dropdown :options="userMenuOptions" @click="onUserMenuClick">
            <div class="user-badge">
              <t-icon name="user-circle" size="20px" />
              <span>{{ userStore.userInfo?.username || 'User' }}</span>
              <t-icon name="chevron-down" size="16px" />
            </div>
          </t-dropdown>
        </div>
      </t-header>
      <t-content class="content">
        <router-view />
      </t-content>
    </t-layout>
  </t-layout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '../stores/user'
import { useI18n } from 'vue-i18n'
import { MessagePlugin } from 'tdesign-vue-next'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const { locale } = useI18n()

const collapsed = ref(false)
const currentLocale = ref(locale.value)

const activeMenu = computed(() => {
  const name = route.name as string
  if (name?.startsWith('Project')) return 'Projects'
  if (name?.startsWith('Monitor')) return 'Monitors'
  return name || 'Dashboard'
})

const isAdmin = computed(() => userStore.isAdmin())

const userMenuOptions = [
  { content: '退出登录', value: 'logout' },
]

function onMenuChange(value: string) {
  router.push({ name: value })
}

function onLocaleChange(value: string) {
  locale.value = value
  localStorage.setItem('locale', value)
}

function onUserMenuClick(data: { value: string }) {
  if (data.value === 'logout') {
    userStore.logout()
    router.push('/login')
    MessagePlugin.success('已退出登录')
  }
}

onMounted(() => {
  userStore.fetchProfile()
})
</script>

<style scoped>
.main-layout {
  height: 100vh;
  background: var(--upp-bg-primary);
}

.sidebar {
  background: linear-gradient(180deg, #060e1a 0%, #0a1628 100%) !important;
  border-right: 1px solid var(--upp-border) !important;
  transition: all 0.3s;
}

.logo {
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  border-bottom: 1px solid var(--upp-border);
  padding: 0 16px;
}

.logo-svg {
  width: 32px;
  height: 32px;
  flex-shrink: 0;
}

.logo-svg-small {
  width: 28px;
  height: 28px;
}

.logo-text {
  font-size: 18px;
  font-weight: 800;
  letter-spacing: 2px;
  background: linear-gradient(135deg, #fff 0%, #a0d2ff 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  background: rgba(10, 22, 40, 0.8) !important;
  border-bottom: 1px solid var(--upp-border) !important;
  height: 64px;
  backdrop-filter: blur(12px);
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.collapse-btn {
  color: var(--upp-text-secondary) !important;
}

.lang-select {
  width: 100px;
}

.user-badge {
  display: flex;
  align-items: center;
  gap: 6px;
  color: var(--upp-text-secondary);
  cursor: pointer;
  padding: 6px 12px;
  border-radius: 8px;
  transition: all 0.2s;
}
.user-badge:hover {
  background: rgba(0, 212, 255, 0.06);
  color: var(--upp-accent);
}

.content {
  padding: 24px;
  background: var(--upp-bg-primary) !important;
  overflow-y: auto;
}

/* Menu overrides */
:deep(.t-menu) {
  background: transparent !important;
  color: rgba(255, 255, 255, 0.5) !important;
}
:deep(.t-menu-item) {
  color: rgba(255, 255, 255, 0.5) !important;
  border-radius: 8px !important;
  margin: 2px 8px !important;
  transition: all 0.2s;
}
:deep(.t-menu-item:hover) {
  color: rgba(255, 255, 255, 0.8) !important;
  background: rgba(0, 212, 255, 0.06) !important;
}
:deep(.t-menu-item.t-is-active) {
  color: #fff !important;
  background: linear-gradient(135deg, rgba(0, 168, 255, 0.2), rgba(0, 102, 255, 0.2)) !important;
  box-shadow: 0 0 12px rgba(0, 180, 255, 0.1);
}
:deep(.t-menu-item .t-icon) {
  color: inherit !important;
}
</style>
