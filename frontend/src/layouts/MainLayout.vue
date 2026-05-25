<template>
  <t-layout class="main-layout">
    <t-aside class="sidebar" :collapsed="collapsed">
      <div class="logo">
        <h2 v-if="!collapsed">UUP</h2>
        <h2 v-else>U</h2>
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
          <t-button variant="text" @click="collapsed = !collapsed">
            <template #icon><t-icon name="view-list" /></template>
          </t-button>
        </div>
        <div class="header-right">
          <t-select v-model="locale" size="small" style="width: 100px" @change="onLocaleChange">
            <t-option value="zh" label="中文" />
            <t-option value="en" label="English" />
          </t-select>
          <t-dropdown :options="userMenuOptions" @click="onUserMenuClick">
            <t-button variant="text">
              {{ userStore.userInfo?.username || 'User' }}
              <template #suffix><t-icon name="chevron-down" /></template>
            </t-button>
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
}
.sidebar {
  background: #001529;
  transition: all 0.3s;
}
.logo {
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
}
.logo h2 {
  margin: 0;
  font-size: 20px;
}
.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  background: #fff;
  border-bottom: 1px solid #e7e7e7;
  height: 64px;
}
.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}
.content {
  padding: 24px;
  background: #f0f2f5;
  overflow-y: auto;
}
:deep(.t-menu) {
  background: transparent;
  color: rgba(255, 255, 255, 0.65);
}
:deep(.t-menu-item.t-is-active) {
  color: #fff;
  background: #1890ff;
}
</style>
