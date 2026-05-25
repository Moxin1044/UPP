<template>
  <div class="login-wrapper">
    <canvas ref="canvasRef" class="tech-canvas"></canvas>

    <div class="bg-orb bg-orb-1"></div>
    <div class="bg-orb bg-orb-2"></div>

    <div class="login-card">
      <!-- 左侧品牌区 -->
      <div class="brand-panel">
        <div class="brand-glow brand-glow-1"></div>
        <div class="brand-glow brand-glow-2"></div>

        <div class="brand-inner">
          <div class="brand-logo">
            <svg viewBox="0 0 64 64" fill="none" xmlns="http://www.w3.org/2000/svg" class="logo-svg">
              <defs>
                <linearGradient id="logoGrad" x1="0" y1="0" x2="64" y2="64">
                  <stop stop-color="#00d4ff" />
                  <stop offset="1" stop-color="#7b61ff" />
                </linearGradient>
                <filter id="logoGlow">
                  <feGaussianBlur stdDeviation="3" result="blur" />
                  <feMerge><feMergeNode in="blur" /><feMergeNode in="SourceGraphic" /></feMerge>
                </filter>
              </defs>
              <rect x="2" y="2" width="60" height="60" rx="16" stroke="url(#logoGrad)" stroke-width="2" />
              <path d="M32 12L48 20V44L32 52L16 44V20L32 12Z" stroke="url(#logoGrad)" stroke-width="1.5" fill="none" filter="url(#logoGlow)" />
              <circle cx="32" cy="32" r="6" fill="url(#logoGrad)" opacity="0.8" />
              <path d="M32 12V52M16 20L48 44M48 20L16 44" stroke="url(#logoGrad)" stroke-width="0.5" opacity="0.4" />
            </svg>
          </div>

          <h1 class="brand-name">UPP Monitor</h1>
          <p class="brand-tagline">{{ $t('login.tagline') }}</p>

          <div class="brand-features">
            <div class="feature-item">
              <div class="feature-dot"></div>
              <span>{{ $t('login.feature1') }}</span>
            </div>
            <div class="feature-item">
              <div class="feature-dot"></div>
              <span>{{ $t('login.feature2') }}</span>
            </div>
            <div class="feature-item">
              <div class="feature-dot"></div>
              <span>{{ $t('login.feature3') }}</span>
            </div>
          </div>
        </div>

        <div class="brand-footer-text">
          <span>&copy; {{ new Date().getFullYear() }} UPP Monitor</span>
        </div>
      </div>

      <!-- 右侧表单区 -->
      <div class="form-panel">
        <div class="form-lang-switch">
          <t-select v-model="currentLocale" size="small" style="width: 100px" @change="onLocaleChange">
            <t-option value="zh" label="中文" />
            <t-option value="en" label="English" />
          </t-select>
        </div>

        <div class="form-header">
          <h2>{{ $t('login.title') }}</h2>
          <p>{{ $t('login.subtitle') }}</p>
        </div>

        <t-form ref="formRef" :data="formData" :rules="rules" label-align="top" @submit="onSubmit">
          <t-form-item :label="$t('login.username')" name="username">
            <t-input v-model="formData.username" :placeholder="$t('login.usernamePlaceholder')" clearable size="large">
              <template #prefix-icon><t-icon name="user" /></template>
            </t-input>
          </t-form-item>

          <t-form-item :label="$t('login.password')" name="password">
            <t-input v-model="formData.password" type="password" :placeholder="$t('login.passwordPlaceholder')" clearable size="large">
              <template #prefix-icon><t-icon name="lock-on" /></template>
            </t-input>
          </t-form-item>

          <t-form-item>
            <t-button theme="primary" type="submit" block size="large" :loading="loading" class="submit-btn">
              <span>{{ $t('login.submit') }}</span>
              <svg class="btn-arrow" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M5 12h14M12 5l7 7-7 7" />
              </svg>
            </t-button>
          </t-form-item>
        </t-form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'
import { useI18n } from 'vue-i18n'
import { MessagePlugin } from 'tdesign-vue-next'

const router = useRouter()
const userStore = useUserStore()
const { t, locale } = useI18n()

const loading = ref(false)
const canvasRef = ref<HTMLCanvasElement>()
const currentLocale = ref(locale.value)

const formData = reactive({ username: '', password: '' })
const rules = {
  username: [{ required: true, message: () => t('login.usernameRequired') }],
  password: [{ required: true, message: () => t('login.passwordRequired') }],
}

function onLocaleChange(val: string) {
  locale.value = val
  localStorage.setItem('locale', val)
}

async function onSubmit({ validateResult }: any) {
  if (validateResult !== true) return
  loading.value = true
  try {
    await userStore.login(formData.username, formData.password)
    MessagePlugin.success(t('login.success'))
    router.push('/')
  } catch (e: any) {
    MessagePlugin.error(e.message || t('login.failed'))
  } finally {
    loading.value = false
  }
}

// ===== 科技粒子背景动画 =====
let animId = 0

onMounted(() => {
  initCanvas()
})

onUnmounted(() => {
  cancelAnimationFrame(animId)
})

const initCanvas = () => {
  const canvas = canvasRef.value
  if (!canvas) return
  const ctx = canvas.getContext('2d')
  if (!ctx) return

  const resize = () => {
    canvas.width = window.innerWidth
    canvas.height = window.innerHeight
  }
  resize()
  window.addEventListener('resize', resize)

  interface Particle {
    x: number; y: number; vx: number; vy: number; r: number; alpha: number
  }

  const particles: Particle[] = []
  const count = 80

  for (let i = 0; i < count; i++) {
    particles.push({
      x: Math.random() * canvas.width,
      y: Math.random() * canvas.height,
      vx: (Math.random() - 0.5) * 0.6,
      vy: (Math.random() - 0.5) * 0.6,
      r: Math.random() * 2 + 0.5,
      alpha: Math.random() * 0.5 + 0.15,
    })
  }

  const draw = () => {
    const w = canvas.width
    const h = canvas.height
    ctx.clearRect(0, 0, w, h)

    for (const p of particles) {
      p.x += p.vx
      p.y += p.vy
      if (p.x < 0) p.x = w
      if (p.x > w) p.x = 0
      if (p.y < 0) p.y = h
      if (p.y > h) p.y = 0

      ctx.beginPath()
      ctx.arc(p.x, p.y, p.r, 0, Math.PI * 2)
      ctx.fillStyle = `rgba(0, 212, 255, ${p.alpha})`
      ctx.fill()
    }

    for (let i = 0; i < particles.length; i++) {
      for (let j = i + 1; j < particles.length; j++) {
        const dx = particles[i].x - particles[j].x
        const dy = particles[i].y - particles[j].y
        const dist = Math.sqrt(dx * dx + dy * dy)
        if (dist < 140) {
          ctx.beginPath()
          ctx.moveTo(particles[i].x, particles[i].y)
          ctx.lineTo(particles[j].x, particles[j].y)
          ctx.strokeStyle = `rgba(0, 180, 255, ${0.08 * (1 - dist / 140)})`
          ctx.lineWidth = 0.5
          ctx.stroke()
        }
      }
    }

    animId = requestAnimationFrame(draw)
  }

  draw()
}
</script>

<style scoped>
.login-wrapper {
  height: 100vh;
  width: 100vw;
  display: flex;
  align-items: center;
  justify-content: center;
  background: radial-gradient(ellipse at 30% 50%, #0a1628 0%, #060e1a 50%, #020810 100%);
  overflow: hidden;
  position: relative;
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
}

.tech-canvas {
  position: absolute;
  inset: 0;
  z-index: 0;
  pointer-events: none;
}

.login-card {
  position: relative;
  z-index: 1;
  display: flex;
  width: 960px;
  height: 540px;
  border-radius: 20px;
  overflow: hidden;
  box-shadow:
    0 0 0 1px rgba(255, 255, 255, 0.06),
    0 40px 80px rgba(0, 0, 0, 0.5),
    0 0 120px rgba(0, 180, 255, 0.08);
  backdrop-filter: blur(4px);
}

/* 背景装饰圆形 */
.bg-orb {
  position: absolute;
  border-radius: 50%;
  pointer-events: none;
  z-index: 0;
  filter: blur(60px);
  opacity: 0.18;
}
.bg-orb-1 {
  width: 500px; height: 500px;
  top: -15%; left: -8%;
  background: radial-gradient(circle, #00d4ff 0%, #0066ff 40%, transparent 70%);
  animation: orbFloat1 12s ease-in-out infinite;
}
.bg-orb-2 {
  width: 400px; height: 400px;
  bottom: -12%; right: -5%;
  background: radial-gradient(circle, #7b61ff 0%, #a855f7 40%, transparent 70%);
  animation: orbFloat2 14s ease-in-out infinite;
}

@keyframes orbFloat1 {
  0%   { transform: translate(0, 0) scale(1) rotate(0deg); }
  25%  { transform: translate(60px, -40px) scale(1.25) rotate(90deg); }
  50%  { transform: translate(-30px, 30px) scale(0.9) rotate(180deg); }
  75%  { transform: translate(-50px, -20px) scale(1.15) rotate(270deg); }
  100% { transform: translate(0, 0) scale(1) rotate(360deg); }
}
@keyframes orbFloat2 {
  0%   { transform: translate(0, 0) scale(1) rotate(0deg); }
  25%  { transform: translate(-50px, 30px) scale(1.3) rotate(-90deg); }
  50%  { transform: translate(40px, -20px) scale(0.85) rotate(-180deg); }
  75%  { transform: translate(20px, -40px) scale(1.2) rotate(-270deg); }
  100% { transform: translate(0, 0) scale(1) rotate(-360deg); }
}

/* 左侧品牌面板 */
.brand-panel {
  flex: 1;
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 48px 40px;
  background: linear-gradient(160deg, rgba(10, 25, 55, 0.95) 0%, rgba(13, 38, 80, 0.95) 50%, rgba(8, 20, 45, 0.95) 100%);
  overflow: hidden;
}

.brand-glow {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  opacity: 0.25;
  pointer-events: none;
}
.brand-glow-1 {
  width: 300px; height: 300px;
  top: -80px; right: -60px;
  background: radial-gradient(circle, #00d4ff 0%, transparent 70%);
  animation: glowPulse 6s ease-in-out infinite;
}
.brand-glow-2 {
  width: 250px; height: 250px;
  bottom: -60px; left: -40px;
  background: radial-gradient(circle, #7b61ff 0%, transparent 70%);
  animation: glowPulse 6s ease-in-out infinite reverse;
}

@keyframes glowPulse {
  0%, 100% { opacity: 0.2; transform: scale(1); }
  50% { opacity: 0.35; transform: scale(1.15); }
}

.brand-inner {
  position: relative;
  z-index: 1;
  text-align: center;
}

.brand-logo {
  margin-bottom: 20px;
}
.logo-svg {
  width: 64px;
  height: 64px;
}

.brand-name {
  font-size: 26px;
  font-weight: 800;
  color: #fff;
  margin: 0 0 4px;
  letter-spacing: 2px;
  background: linear-gradient(135deg, #fff 0%, #a0d2ff 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.brand-tagline {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.35);
  letter-spacing: 4px;
  margin: 0 0 36px;
  font-weight: 500;
}

.brand-features {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 0 10px;
}

.feature-item {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 13px;
  color: rgba(255, 255, 255, 0.55);
}
.feature-dot {
  width: 5px;
  height: 5px;
  border-radius: 50%;
  background: #00d4ff;
  flex-shrink: 0;
  box-shadow: 0 0 6px rgba(0, 212, 255, 0.5);
  animation: dotBlink 2s ease-in-out infinite;
}
.feature-item:nth-child(2) .feature-dot { animation-delay: 0.3s; }
.feature-item:nth-child(3) .feature-dot { animation-delay: 0.6s; }

@keyframes dotBlink {
  0%, 100% { opacity: 0.4; }
  50% { opacity: 1; }
}

.brand-footer-text {
  position: absolute;
  bottom: 20px;
  font-size: 11px;
  color: rgba(255, 255, 255, 0.18);
  text-align: center;
}

/* 右侧表单面板 */
.form-panel {
  width: 440px;
  display: flex;
  flex-direction: column;
  align-items: stretch;
  justify-content: center;
  padding: 48px 52px;
  background: rgba(14, 23, 41, 0.92);
  backdrop-filter: blur(20px);
  border-left: 1px solid rgba(255, 255, 255, 0.06);
  position: relative;
}

.form-lang-switch {
  position: absolute;
  top: 16px;
  right: 16px;
  z-index: 1;
}

.form-header {
  text-align: center;
  margin-bottom: 32px;
}
.form-header h2 {
  font-size: 22px;
  font-weight: 700;
  color: #fff;
  margin: 0 0 4px;
}
.form-header p {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.35);
  margin: 0;
}

/* 覆盖 TDesign 暗色样式 */
:deep(.t-form__item) {
  margin-bottom: 20px;
}
:deep(.t-form__label) {
  color: rgba(255, 255, 255, 0.6) !important;
  font-size: 13px;
  font-weight: 500;
  margin-bottom: 6px;
}
:deep(.t-input) {
  background: rgba(255, 255, 255, 0.04) !important;
  border-color: rgba(255, 255, 255, 0.1) !important;
  border-radius: 8px !important;
  transition: all 0.3s ease;
}
:deep(.t-input:hover) {
  border-color: rgba(0, 212, 255, 0.35) !important;
  background: rgba(255, 255, 255, 0.06) !important;
}
:deep(.t-input.t-is-focused) {
  border-color: #00d4ff !important;
  box-shadow: 0 0 0 2px rgba(0, 212, 255, 0.12) !important;
  background: rgba(255, 255, 255, 0.06) !important;
}
:deep(.t-input input) {
  color: #fff !important;
  caret-color: #00d4ff;
}
:deep(.t-input input::placeholder) {
  color: rgba(255, 255, 255, 0.2) !important;
}
:deep(.t-input__prefix) {
  color: rgba(255, 255, 255, 0.3);
}
:deep(.t-input__suffix) {
  color: rgba(255, 255, 255, 0.2);
}

.submit-btn {
  width: 100%;
  height: 46px;
  border-radius: 8px;
  font-size: 15px;
  font-weight: 600;
  letter-spacing: 1px;
  margin-top: 4px;
  position: relative;
  overflow: hidden;
  background: linear-gradient(135deg, #00a8ff, #0066ff) !important;
  border: none !important;
  transition: all 0.3s ease;
}
.submit-btn .btn-arrow {
  width: 18px;
  height: 18px;
  margin-left: 8px;
  transition: transform 0.3s ease;
}
.submit-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 6px 25px rgba(0, 150, 255, 0.4);
  background: linear-gradient(135deg, #00b8ff, #0070ff) !important;
}
.submit-btn:hover .btn-arrow {
  transform: translateX(3px);
}
.submit-btn:active {
  transform: translateY(0);
}
.submit-btn::after {
  content: '';
  position: absolute;
  top: 0; left: -100%;
  width: 100%; height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255,255,255,0.15), transparent);
  transition: left 0.6s;
}
.submit-btn:hover::after {
  left: 100%;
}

/* 响应式 */
@media (max-width: 1000px) {
  .login-card {
    width: 94vw;
    height: auto;
    min-height: 560px;
    flex-direction: column;
  }
  .brand-panel {
    padding: 36px 28px;
  }
  .brand-features { display: none; }
  .brand-tagline { margin-bottom: 24px; }
  .brand-footer-text { display: none; }
  .form-panel {
    width: 100%;
    padding: 36px 32px;
    border-left: none;
    border-top: 1px solid rgba(255, 255, 255, 0.06);
  }
}
@media (max-width: 560px) {
  .login-card { border-radius: 12px; }
  .brand-panel { padding: 28px 20px; }
  .form-panel { padding: 28px 20px; }
  .brand-name { font-size: 22px; }
}
</style>
