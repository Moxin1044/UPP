<template>
  <div class="register-wrapper">
    <canvas ref="canvasRef" class="tech-canvas"></canvas>
    <div class="bg-orb bg-orb-1"></div>
    <div class="bg-orb bg-orb-2"></div>

    <div class="register-card">
      <div class="form-panel">
        <div class="form-lang-switch">
          <t-select v-model="currentLocale" size="small" style="width: 100px" @change="onLocaleChange">
            <t-option value="zh" label="中文" />
            <t-option value="en" label="English" />
          </t-select>
        </div>

        <div class="form-header">
          <svg viewBox="0 0 64 64" fill="none" xmlns="http://www.w3.org/2000/svg" class="form-logo">
            <defs>
              <linearGradient id="regLogoGrad" x1="0" y1="0" x2="64" y2="64">
                <stop stop-color="#00d4ff" />
                <stop offset="1" stop-color="#7b61ff" />
              </linearGradient>
            </defs>
            <rect x="2" y="2" width="60" height="60" rx="16" stroke="url(#regLogoGrad)" stroke-width="2" />
            <circle cx="32" cy="32" r="6" fill="url(#regLogoGrad)" opacity="0.8" />
          </svg>
          <h2>{{ $t('register.title') }}</h2>
          <p>{{ $t('register.subtitle') }}</p>
        </div>

        <t-form ref="formRef" :data="formData" :rules="rules" label-align="top" @submit="onSubmit">
          <t-form-item :label="$t('register.username')" name="username">
            <t-input v-model="formData.username" :placeholder="$t('register.usernamePlaceholder')" clearable size="large">
              <template #prefix-icon><t-icon name="user" /></template>
            </t-input>
          </t-form-item>
          <t-form-item :label="$t('register.password')" name="password">
            <t-input v-model="formData.password" type="password" :placeholder="$t('register.passwordPlaceholder')" clearable size="large">
              <template #prefix-icon><t-icon name="lock-on" /></template>
            </t-input>
          </t-form-item>
          <t-form-item :label="$t('register.confirmPassword')" name="confirmPassword">
            <t-input v-model="formData.confirmPassword" type="password" :placeholder="$t('register.confirmPasswordPlaceholder')" clearable size="large">
              <template #prefix-icon><t-icon name="secured" /></template>
            </t-input>
          </t-form-item>
          <t-form-item>
            <t-button theme="primary" type="submit" block size="large" :loading="loading" class="submit-btn">
              <span>{{ $t('register.submit') }}</span>
              <svg class="btn-arrow" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M5 12h14M12 5l7 7-7 7" />
              </svg>
            </t-button>
          </t-form-item>
        </t-form>

        <div class="form-footer">
          {{ $t('register.hasAccount') }}
          <t-link theme="primary" @click="$router.push('/login')">{{ $t('register.toLogin') }}</t-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { register } from '../api/auth'
import { useI18n } from 'vue-i18n'
import { MessagePlugin } from 'tdesign-vue-next'

const router = useRouter()
const { t, locale } = useI18n()

const loading = ref(false)
const currentLocale = ref(locale.value)
const canvasRef = ref<HTMLCanvasElement>()

const formData = reactive({ username: '', password: '', confirmPassword: '' })
const rules = {
  username: [{ required: true, message: () => t('register.usernameRequired') }, { min: 3, message: 'min 3' }],
  password: [{ required: true, message: () => t('register.passwordRequired') }, { min: 6, message: 'min 6' }],
  confirmPassword: [{ required: true, message: () => t('register.confirmPasswordRequired') }],
}

function onLocaleChange(val: string) {
  locale.value = val
  localStorage.setItem('locale', val)
}

async function onSubmit({ validateResult }: any) {
  if (validateResult !== true) return
  if (formData.password !== formData.confirmPassword) {
    MessagePlugin.error(t('register.passwordMismatch'))
    return
  }
  loading.value = true
  try {
    await register({ username: formData.username, password: formData.password })
    MessagePlugin.success(t('common.success'))
    router.push('/login')
  } catch (e: any) {
    MessagePlugin.error(e.message || t('common.failed'))
  } finally {
    loading.value = false
  }
}

// Particle animation (same as login)
let animId = 0
onMounted(() => { initCanvas() })
onUnmounted(() => { cancelAnimationFrame(animId) })

const initCanvas = () => {
  const canvas = canvasRef.value
  if (!canvas) return
  const ctx = canvas.getContext('2d')
  if (!ctx) return
  const resize = () => { canvas.width = window.innerWidth; canvas.height = window.innerHeight }
  resize()
  window.addEventListener('resize', resize)
  interface P { x: number; y: number; vx: number; vy: number; r: number; alpha: number }
  const particles: P[] = []
  for (let i = 0; i < 80; i++) {
    particles.push({ x: Math.random() * canvas.width, y: Math.random() * canvas.height, vx: (Math.random() - 0.5) * 0.6, vy: (Math.random() - 0.5) * 0.6, r: Math.random() * 2 + 0.5, alpha: Math.random() * 0.5 + 0.15 })
  }
  const draw = () => {
    const w = canvas.width, h = canvas.height
    ctx.clearRect(0, 0, w, h)
    for (const p of particles) {
      p.x += p.vx; p.y += p.vy
      if (p.x < 0) p.x = w; if (p.x > w) p.x = 0; if (p.y < 0) p.y = h; if (p.y > h) p.y = 0
      ctx.beginPath(); ctx.arc(p.x, p.y, p.r, 0, Math.PI * 2); ctx.fillStyle = `rgba(0, 212, 255, ${p.alpha})`; ctx.fill()
    }
    for (let i = 0; i < particles.length; i++) {
      for (let j = i + 1; j < particles.length; j++) {
        const dx = particles[i].x - particles[j].x, dy = particles[i].y - particles[j].y, dist = Math.sqrt(dx * dx + dy * dy)
        if (dist < 140) { ctx.beginPath(); ctx.moveTo(particles[i].x, particles[i].y); ctx.lineTo(particles[j].x, particles[j].y); ctx.strokeStyle = `rgba(0, 180, 255, ${0.08 * (1 - dist / 140)})`; ctx.lineWidth = 0.5; ctx.stroke() }
      }
    }
    animId = requestAnimationFrame(draw)
  }
  draw()
}
</script>

<style scoped>
.register-wrapper {
  height: 100vh; width: 100vw;
  display: flex; align-items: center; justify-content: center;
  background: radial-gradient(ellipse at 30% 50%, #0a1628 0%, #060e1a 50%, #020810 100%);
  overflow: hidden; position: relative;
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
}
.tech-canvas { position: absolute; inset: 0; z-index: 0; pointer-events: none; }

.bg-orb { position: absolute; border-radius: 50%; pointer-events: none; z-index: 0; filter: blur(60px); opacity: 0.18; }
.bg-orb-1 { width: 500px; height: 500px; top: -15%; left: -8%; background: radial-gradient(circle, #00d4ff 0%, #0066ff 40%, transparent 70%); animation: orbFloat1 12s ease-in-out infinite; }
.bg-orb-2 { width: 400px; height: 400px; bottom: -12%; right: -5%; background: radial-gradient(circle, #7b61ff 0%, #a855f7 40%, transparent 70%); animation: orbFloat2 14s ease-in-out infinite; }

@keyframes orbFloat1 { 0% { transform: translate(0,0) scale(1); } 50% { transform: translate(60px,-40px) scale(1.25); } 100% { transform: translate(0,0) scale(1); } }
@keyframes orbFloat2 { 0% { transform: translate(0,0) scale(1); } 50% { transform: translate(-50px,30px) scale(1.3); } 100% { transform: translate(0,0) scale(1); } }

.register-card {
  position: relative; z-index: 1;
  width: 440px;
  border-radius: 20px; overflow: hidden;
  box-shadow: 0 0 0 1px rgba(255,255,255,0.06), 0 40px 80px rgba(0,0,0,0.5), 0 0 120px rgba(0,180,255,0.08);
  backdrop-filter: blur(4px);
}

.form-panel {
  padding: 48px 44px;
  background: rgba(14, 23, 41, 0.92);
  backdrop-filter: blur(20px);
  position: relative;
}

.form-lang-switch { position: absolute; top: 16px; right: 16px; z-index: 1; }

.form-header { text-align: center; margin-bottom: 32px; }
.form-logo { width: 48px; height: 48px; margin-bottom: 12px; }
.form-header h2 { font-size: 22px; font-weight: 700; color: #fff; margin: 0 0 4px; }
.form-header p { font-size: 13px; color: rgba(255,255,255,0.35); margin: 0; }

/* TDesign dark overrides */
:deep(.t-form__item) { margin-bottom: 20px; }
:deep(.t-form__label) { color: rgba(255,255,255,0.6) !important; font-size: 13px; font-weight: 500; margin-bottom: 6px; }
:deep(.t-input) { background: rgba(255,255,255,0.04) !important; border-color: rgba(255,255,255,0.1) !important; border-radius: 8px !important; }
:deep(.t-input:hover) { border-color: rgba(0,212,255,0.35) !important; background: rgba(255,255,255,0.06) !important; }
:deep(.t-input.t-is-focused) { border-color: #00d4ff !important; box-shadow: 0 0 0 2px rgba(0,212,255,0.12) !important; }
:deep(.t-input input) { color: #fff !important; caret-color: #00d4ff; }
:deep(.t-input input::placeholder) { color: rgba(255,255,255,0.2) !important; }
:deep(.t-input__prefix) { color: rgba(255,255,255,0.3); }

.submit-btn {
  width: 100%; height: 46px; border-radius: 8px; font-size: 15px; font-weight: 600;
  background: linear-gradient(135deg, #00a8ff, #0066ff) !important; border: none !important;
  position: relative; overflow: hidden; transition: all 0.3s;
}
.submit-btn .btn-arrow { width: 18px; height: 18px; margin-left: 8px; transition: transform 0.3s; }
.submit-btn:hover { transform: translateY(-1px); box-shadow: 0 6px 25px rgba(0,150,255,0.4); }
.submit-btn:hover .btn-arrow { transform: translateX(3px); }

.form-footer { text-align: center; margin-top: 16px; color: rgba(255,255,255,0.35); font-size: 13px; }
.form-footer :deep(.t-link) { color: #00d4ff !important; }

@media (max-width: 560px) {
  .register-card { width: 94vw; border-radius: 12px; }
  .form-panel { padding: 28px 20px; }
}
</style>
