import { createApp } from 'vue'
import { createPinia } from 'pinia'
import TDesign from 'tdesign-vue-next'
import 'tdesign-vue-next/es/style/index.css'
import i18n from './locales'
import router from './router'
import App from './App.vue'
import './style.css'
import './styles/dark-theme.css'

const app = createApp(App)
app.use(createPinia())
app.use(TDesign)
app.use(i18n)
app.use(router)
app.mount('#app')
