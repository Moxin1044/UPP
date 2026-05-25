import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      name: 'Login',
      component: () => import('../pages/Login.vue'),
      meta: { public: true },
    },
    {
      path: '/register',
      name: 'Register',
      component: () => import('../pages/Register.vue'),
      meta: { public: true },
    },
    {
      path: '/',
      component: () => import('../layouts/MainLayout.vue'),
      children: [
        { path: '', name: 'Dashboard', component: () => import('../pages/Dashboard.vue') },
        { path: 'projects', name: 'Projects', component: () => import('../pages/Projects.vue') },
        { path: 'projects/:id', name: 'ProjectDetail', component: () => import('../pages/ProjectDetail.vue') },
        { path: 'monitors', name: 'Monitors', component: () => import('../pages/Monitors.vue') },
        { path: 'monitors/:id', name: 'MonitorDetail', component: () => import('../pages/MonitorDetail.vue') },
        { path: 'notifications', name: 'Notifications', component: () => import('../pages/Notifications.vue') },
        { path: 'users', name: 'Users', component: () => import('../pages/Users.vue'), meta: { admin: true } },
        { path: 'settings', name: 'Settings', component: () => import('../pages/Settings.vue'), meta: { admin: true } },
      ],
    },
  ],
})

router.afterEach((to) => {
  const title = to.name ? `${String(to.name)} - UPP Monitor` : 'UPP Monitor'
  document.title = title
})

router.beforeEach((to, _from, next) => {
  const token = localStorage.getItem('token')
  if (to.meta.public) {
    if (token && (to.name === 'Login' || to.name === 'Register')) {
      next({ name: 'Dashboard' })
    } else {
      next()
    }
  } else if (!token) {
    next({ name: 'Login' })
  } else {
    next()
  }
})

export default router
