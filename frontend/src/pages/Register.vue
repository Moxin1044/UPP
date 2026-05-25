<template>
  <div class="login-page">
    <t-card class="login-card">
      <h2>{{ $t('register.title') }}</h2>
      <t-form ref="formRef" :data="formData" :rules="rules" @submit="onSubmit">
        <t-form-item :label="$t('register.username')" name="username">
          <t-input v-model="formData.username" :placeholder="$t('register.username')" />
        </t-form-item>
        <t-form-item :label="$t('register.password')" name="password">
          <t-input v-model="formData.password" type="password" :placeholder="$t('register.password')" />
        </t-form-item>
        <t-form-item :label="$t('register.confirmPassword')" name="confirmPassword">
          <t-input v-model="formData.confirmPassword" type="password" :placeholder="$t('register.confirmPassword')" />
        </t-form-item>
        <t-form-item>
          <t-button theme="primary" type="submit" block :loading="loading">
            {{ $t('register.submit') }}
          </t-button>
        </t-form-item>
      </t-form>
      <div class="login-footer">
        {{ $t('register.hasAccount') }}
        <t-link theme="primary" @click="$router.push('/login')">{{ $t('register.toLogin') }}</t-link>
      </div>
    </t-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { register } from '../api/auth'
import { useI18n } from 'vue-i18n'
import { MessagePlugin } from 'tdesign-vue-next'

const router = useRouter()
const { t } = useI18n()

const loading = ref(false)
const formData = reactive({ username: '', password: '', confirmPassword: '' })
const rules = {
  username: [{ required: true, message: 'required' }, { min: 3, message: 'min 3' }],
  password: [{ required: true, message: 'required' }, { min: 6, message: 'min 6' }],
  confirmPassword: [{ required: true, message: 'required' }],
}

async function onSubmit({ validateResult }: any) {
  if (validateResult !== true) return
  if (formData.password !== formData.confirmPassword) {
    MessagePlugin.error('Passwords do not match')
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
