<template>
  <div class="login-page">
    <t-card class="login-card">
      <h2>{{ $t('login.title') }}</h2>
      <t-form ref="formRef" :data="formData" :rules="rules" @submit="onSubmit">
        <t-form-item :label="$t('login.username')" name="username">
          <t-input v-model="formData.username" :placeholder="$t('login.username')" />
        </t-form-item>
        <t-form-item :label="$t('login.password')" name="password">
          <t-input v-model="formData.password" type="password" :placeholder="$t('login.password')" />
        </t-form-item>
        <t-form-item>
          <t-checkbox v-model="remember">{{ $t('login.remember') }}</t-checkbox>
        </t-form-item>
        <t-form-item>
          <t-button theme="primary" type="submit" block :loading="loading">
            {{ $t('login.submit') }}
          </t-button>
        </t-form-item>
      </t-form>
      <div class="login-footer">
        {{ $t('login.noAccount') }}
        <t-link theme="primary" @click="$router.push('/register')">{{ $t('login.register') }}</t-link>
      </div>
    </t-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'
import { useI18n } from 'vue-i18n'
import { MessagePlugin } from 'tdesign-vue-next'

const router = useRouter()
const userStore = useUserStore()
const { t } = useI18n()

const loading = ref(false)
const remember = ref(false)
const formData = reactive({ username: '', password: '' })
const rules = {
  username: [{ required: true, message: t('login.username') }],
  password: [{ required: true, message: t('login.password') }],
}

async function onSubmit({ validateResult }: any) {
  if (validateResult !== true) return
  loading.value = true
  try {
    await userStore.login(formData.username, formData.password, remember.value)
    MessagePlugin.success(t('common.success'))
    router.push('/')
  } catch (e: any) {
    MessagePlugin.error(e.message || t('common.failed'))
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}
.login-card {
  width: 400px;
  padding: 24px;
}
.login-card h2 {
  text-align: center;
  margin-bottom: 24px;
}
.login-footer {
  text-align: center;
  margin-top: 16px;
}
</style>
